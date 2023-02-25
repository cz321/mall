package handler

import (
	"context"
	"encoding/json"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

//获取该用户的购物车列表
func (*OrderServer) CartItemList(_ context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart

	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	}

	var data []*proto.ShopCartInfoResponse
	for _, shopCart := range shopCarts {
		data = append(data, &proto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.User,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}

	return &proto.CartItemListResponse{
		Total: int32(result.RowsAffected),
		Data:  data,
	}, nil
}

//商品加入到购物车
func (*OrderServer) CreateCartItem(_ context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
	//1.购物车没有这种商品
	//2.购物车有这种商品
	var shopCart model.ShoppingCart

	result := global.DB.Where(&model.ShoppingCart{User: req.Id, Goods: req.GoodsId}).First(&shopCart)
	if result.RowsAffected == 0 {
		//购物车没有这种商品
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	} else if result.RowsAffected == 1 {
		//购物车有这种商品,跟新记录
		shopCart.Nums += req.Nums
	}

	global.DB.Save(&shopCart)
	return &proto.ShopCartInfoResponse{
		Id: shopCart.ID,
	}, nil
}

//更新购物车某个记录
func (*OrderServer) UpdateCartItem(_ context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	var shopCart model.ShoppingCart

	result := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).First(&shopCart)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "记录不存在")
	}

	//更新数量
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	//更新选中状态
	shopCart.Checked = req.Checked

	global.DB.Save(&shopCart)

	return &empty.Empty{}, nil
}

//删除购物车某个记录
func (*OrderServer) DeleteCartItem(_ context.Context, req *proto.CartItemRequest) (*empty.Empty, error) {
	result := global.DB.Where("goods = ? and user = ?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{})
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "记录不存在")
	}

	return &empty.Empty{}, nil
}

type OrderListener struct{
	Code codes.Code
	Detail string
	ID int32
	OrderAmount float32
	Ctx context.Context
}

func NewOrderListener(ctx context.Context) *OrderListener {
	return &OrderListener{
		Ctx: ctx,
	}
}

//运行本地事务
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	parentSpan := opentracing.SpanFromContext(o.Ctx)

	/*
		1. 从购物车获取选中的商品
		2. 商品金额查询 ->商品服务
		3. 库存扣减 ->库存服务
		4. 生成订单表 ->订单商品信息表
		5. 从购物车中删除已购买的记录
	*/

	//1. 从购物车获取选中的商品
	shopCartSpan := opentracing.GlobalTracer().StartSpan("select_shopcart",opentracing.ChildOf(parentSpan.Context()))
	var shopingCarts []model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopingCarts)
	if result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "没有选中结算的商品"
		return primitive.RollbackMessageState
	}
	shopCartSpan.Finish()

	var goodsIds []int32
	goodsIdToNum := make(map[int32]int32)
	for _, shopingCart := range shopingCarts {
		goodsIds = append(goodsIds, shopingCart.Goods)
		goodsIdToNum[shopingCart.Goods] = shopingCart.Nums
	}

	//2. 商品金额查询 ->商品服务
	queryGoodsSpan := opentracing.GlobalTracer().StartSpan("query_goods",opentracing.ChildOf(parentSpan.Context()))
	goodsListResponse, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		o.Code = codes.InvalidArgument
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	queryGoodsSpan.Finish()

	//订单合计价格
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfos []*proto.GoodsInvInfo
	for _, goodsInfo := range goodsListResponse.Data {
		orderAmount += goodsInfo.ShopPrice * float32(goodsIdToNum[goodsInfo.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      goodsInfo.Id,
			GoodsName:  goodsInfo.Name,
			GoodsImage: goodsInfo.GoodsFrontImage,
			GoodsPrice: goodsInfo.ShopPrice,
			Nums:       goodsIdToNum[goodsInfo.Id],
		})
		goodsInvInfos = append(goodsInvInfos, &proto.GoodsInvInfo{
			GoodsId: goodsInfo.Id,
			Num:     goodsIdToNum[goodsInfo.Id],
		})
	}

	//3. 库存扣减 ->库存服务
	invSellSpan := opentracing.GlobalTracer().StartSpan("inv_sell",opentracing.ChildOf(parentSpan.Context()))
	_, err = global.InventorySrvClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: goodsInvInfos,
		OrderSn: orderInfo.OrderSn,
	})
	if err != nil {
		o.Code = codes.InvalidArgument
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	invSellSpan.Finish()
	/*
		到这里，库存微服务已经扣减成功了
		1.如果本地事务出现问题或网络出现故障，应该进行commit提交，让半消息执行发送进行归还库存，确保最终一致性
		2.如果本地事务执行没有问题，那么进行rollback，半消息不会发送，库存最终就不会归还
	 */

	//4. 生成订单表 ->订单商品信息表
	orderInfo.OrderMount = orderAmount

	//开启事务
	tx := global.DB.Begin()

	//保存订单
	result = global.DB.Save(&orderInfo)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.InvalidArgument
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID

	//补充订单外键
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}

	//批量插入
	result = tx.CreateInBatches(orderGoods, 100)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.InvalidArgument
		o.Detail = "批量插入订单失败"
		return primitive.CommitMessageState
	}

	//5. 从购物车中删除已购买的记录
	result = tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&shopingCarts)
	if result.Error != nil || result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.InvalidArgument
		o.Detail = "删除购物的记录失败"
		return primitive.CommitMessageState
	}


	//发送延时消息
	timeOutMes := primitive.NewMessage("order_timeout",msg.Body)
	timeOutMes.WithDelayTimeLevel(3)
	_, err = global.Producer.SendSync(context.Background(), timeOutMes)

	if err != nil {
		zap.S().Error("发送延时消息失败",err.Error())
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延时消息失败"
		return primitive.CommitMessageState
	}

	zap.S().Info("发送延时消息成功")

	//事务提交
	tx.Commit()

	o.Code = codes.OK

	return primitive.RollbackMessageState
}

//本地事务回查
func (*OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	zap.S().Info("本地事务回查")

	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	result := global.DB.Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo)
	if result.RowsAffected == 0 {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

//新建订单
func (*OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {

	orderListener := NewOrderListener(ctx)

	orderSn := GenerateOrderSn(req.UserId)

	p, _ := rocketmq.NewTransactionProducer(
		orderListener,
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"192.168.10.105:9876"})),
		producer.WithInstanceName("order_create" + orderSn),
		producer.WithRetry(1),
	)

	//启动producer
	err := p.Start()
	if err != nil {
		zap.S().Error("启动producer失败", err)
		return nil,status.Error(codes.Internal,"启动producer失败")
	}

	//订单信息
	orderInfo := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      orderSn,
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
	}

	orderInfoJsonSting, _ := json.Marshal(orderInfo)

	//开启一个order_reback主题，让库存服务去消费
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", orderInfoJsonSting))

	if err != nil {
		zap.S().Error("发送消息失败: ", err)
		return nil,status.Error(codes.Internal,"发送消息失败")
	} else {
		zap.S().Info("发送成功: ", res.String())
	}

	if orderListener.Code != codes.OK{
		zap.S().Error("新建订单失败: ", orderListener.Code)
		return nil,status.Error(codes.Internal,"新建订单失败")
	}

	return &proto.OrderInfoResponse{
		Id:      orderListener.ID,
		UserId:  req.Id,
		OrderSn: orderInfo.OrderSn,
		Total:   orderListener.OrderAmount,
	}, nil
}

//订单列表
func (*OrderServer) OrderList(_ context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	var total int64
	global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: req.UserId}).Count(&total)

	var orderInfos []*model.OrderInfo
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orderInfos)

	var data []*proto.OrderInfoResponse

	for _, orderInfo := range orderInfos {
		data = append(data, &proto.OrderInfoResponse{
			Id:      orderInfo.ID,
			UserId:  orderInfo.User,
			OrderSn: orderInfo.OrderSn,
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			Post:    orderInfo.Post,
			Total:   orderInfo.OrderMount,
			Address: orderInfo.Address,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SingerMobile,
			AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &proto.OrderListResponse{
		Total: int32(total),
		Data:  data,
	}, nil
}

//订单详情
func (*OrderServer) OrderDetail(_ context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	var orderInfo model.OrderInfo
	result := global.DB.Where(&model.OrderInfo{BaseModel: model.BaseModel{ID: req.Id}, User: req.UserId}).First(&orderInfo)
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.InvalidArgument, "订单不存在")
	}

	var orderGoods []model.OrderGoods
	global.DB.Where(&model.OrderGoods{Order: req.Id}).Find(&orderGoods)

	var goods []*proto.OrderItemResponse
	for _, orderGood := range orderGoods {
		goods = append(goods, &proto.OrderItemResponse{
			Id:         orderGood.ID,
			OrderId:    orderGood.Order,
			GoodsId:    orderGood.Goods,
			GoodsName:  orderGood.GoodsName,
			GoodsImage: orderGood.GoodsImage,
			GoodsPrice: orderGood.GoodsPrice,
			Nums:       orderGood.Nums,
		})
	}

	return &proto.OrderInfoDetailResponse{
		OrderInfo: &proto.OrderInfoResponse{
			Id:      orderInfo.ID,
			UserId:  orderInfo.User,
			OrderSn: orderInfo.OrderSn,
			PayType: orderInfo.PayType,
			Status:  orderInfo.Status,
			Post:    orderInfo.Post,
			Total:   orderInfo.OrderMount,
			Address: orderInfo.Address,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SingerMobile,
			AddTime: orderInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		},
		Goods: goods,
	}, nil
}

//更新订单状态
func (*OrderServer) UpdateOrderStatus(_ context.Context, req *proto.OrderStatus) (*empty.Empty, error) {
	result := global.DB.Model(&model.OrderInfo{Status: req.Status}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status)
	if result.Error != nil || result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "订单更新失败")
	}
	return &empty.Empty{}, nil
}


//订单延时消息处理
func OrderTimeout(_ context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	zap.S().Info("******OrderTimeout订单延时消息处理1********")
	for i := range msgs {
		//确保避免重复归还的问题，幂等性
		var orderInfo OrderInfo
		err := json.Unmarshal(msgs[i].Body, &orderInfo)
		if err != nil {
			zap.S().Error("解析json失败",msgs[i].Body)
			return consumer.ConsumeSuccess, nil
		}

		zap.S().Info("******OrderTimeout订单延时消息处理2********")

		//查询订单支付状态如果未支付，归还库存

		var order model.OrderInfo
		result := global.DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&order)
		if result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		if order.Status != "TRADE_SUCCESS" {
			//归还库存
			tx := global.DB.Begin()
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)

			_, err = global.Producer.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))

			if err != nil {
				tx.Rollback()
				zap.S().Error("发送超时归还库存消息失败",err.Error())
				return consumer.ConsumeRetryLater,nil
			}

			tx.Commit()
		}
	}
	return consumer.ConsumeSuccess, nil
}
