package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redsync/redsync/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

//设置库存
func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*empty.Empty, error) {
	var inventory model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inventory)

	inventory.Goods = req.GoodsId
	inventory.Stocks = req.Num

	global.DB.Save(&inventory)
	return &empty.Empty{}, nil
}

//库存详情
func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inventory model.Inventory
	result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inventory)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "无库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inventory.Goods,
		Num:     inventory.Stocks,
	}, nil
}

//库存预扣减
//mysql实现分布式锁
func (*InventoryServer) Sell_(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	//手动事务
	tx := global.DB.Begin()

	//排序避免死锁
	sort.Slice(req.GoodsInfo, func(i, j int) bool {
		return req.GoodsInfo[i].GoodsId < req.GoodsInfo[j].GoodsId
	})

	for _, goodsInfo := range req.GoodsInfo {
		var inventory model.Inventory

		/****************mysql悲观锁***********************/
		//查询索引是行锁，否则升级为表锁，均为写锁
		//result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodsInfo.GoodsId}).First(&inventory)
		//if result.RowsAffected == 0 {
		//	//回滚事务
		//	tx.Rollback()
		//	return nil,status.Error(codes.InvalidArgument,"无库存信息")
		//}
		//if inventory.Stocks < goodsInfo.Num {
		//	//回滚事务
		//	tx.Rollback()
		//	return nil,status.Error(codes.ResourceExhausted,"库存不足")
		//}
		//inventory.Stocks -= goodsInfo.Num
		//tx.Save(&inventory)
		/****************mysql悲观锁***********************/

		/****************mysql乐观锁***********************/
		for {
			result := global.DB.Where(&model.Inventory{Goods: goodsInfo.GoodsId}).First(&inventory)
			if result.RowsAffected == 0 {
				//回滚事务
				tx.Rollback()
				return nil, status.Error(codes.InvalidArgument, "无库存信息")
			}
			if inventory.Stocks < goodsInfo.Num {
				//回滚事务
				tx.Rollback()
				return nil, status.Error(codes.ResourceExhausted, "库存不足")
			}
			inventory.Stocks -= goodsInfo.Num

			//update inventory set stocks = stocks-1,version = version+1 where goods=goods and version=version
			result = tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods = ? and version = ?", goodsInfo.GoodsId, inventory.Version).Updates(model.Inventory{
				Stocks:  inventory.Stocks,
				Version: inventory.Version + 1,
			})
			if result.RowsAffected == 0 {
				zap.S().Info("库存扣减失败")
			} else {
				break
			}
		}
		/****************mysql乐观锁***********************/
	}
	//提交事务
	tx.Commit()
	return &empty.Empty{}, nil
}

//库存预扣减
//redis实现
func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {

	tx := global.DB.Begin()

	var mutexes []*redsync.Mutex

	defer func() {
		for _,mutex := range mutexes {
			if ok, err := mutex.Unlock(); !ok || err != nil {
				zap.S().Error("释放redis分布式锁异常")
			}
		}
	}()

	//排序避免死锁
	sort.Slice(req.GoodsInfo, func(i, j int) bool {
		return req.GoodsInfo[i].GoodsId < req.GoodsInfo[j].GoodsId
	})


	var details []model.GoodsDetail
	for _, goodInfo := range req.GoodsInfo {
		details = append(details,model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})

		var inv model.Inventory

		mutex := global.Redsync.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))


		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		//将获得的锁加入队列
		mutexes = append(mutexes,mutex)

		if result := global.DB.Where(&model.Inventory{Goods:goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}

	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
		Detail: details,
	}

	result := tx.Create(&sellDetail)
	if result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.InvalidArgument, "保存库存扣减历史信息失败")
	}

	tx.Commit() // 需要自己手动提交操作
	return &emptypb.Empty{}, nil
}

//库存归还
func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*empty.Empty, error) {
	//case 1 : 订单超时归还
	//case 2 : 订单创建失败
	//case 3 : 手动归还
	//手动事务
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		var inventory model.Inventory
		result := global.DB.Where(&model.Inventory{Goods: goodsInfo.GoodsId}).First(&inventory)
		if result.RowsAffected == 0 {
			//回滚事务
			tx.Rollback()
			return nil, status.Error(codes.InvalidArgument, "无库存信息")
		}
		inventory.Stocks += goodsInfo.Num
		tx.Save(&inventory)
	}
	//提交事务
	tx.Commit()
	return &empty.Empty{}, nil
}

//自动归还库存
func AutoReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	zap.S().Info("************自动归还库存***************")
	for i := range msgs {
		//确保避免重复归还的问题，幂等性
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Error("解析json失败",msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		zap.S().Info("******AutoReback归还库存********")
		//去将inv的库存加回去
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).First(&sellDetail)
		if result.RowsAffected == 0 {
			//已经归还过了
			return consumer.ConsumeSuccess, nil
		}
		for _,orderGoods := range sellDetail.Detail {
			result :=  tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods:orderGoods.Goods}).Update("stocks", gorm.Expr("stocks+?", orderGoods.Num))
			if result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}

		//2表示已归还，防止多次归还
		sellDetail.Status = 2

		result = tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).Update("status",2)
		if result.RowsAffected == 0 {
			return consumer.ConsumeRetryLater, nil
		}

		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
