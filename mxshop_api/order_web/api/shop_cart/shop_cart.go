package shop_cart

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mxshop_api/order_web/api"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/proto"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: int32(userId.(uint32)),
	})
	if err != nil {
		zap.S().Error("[List] 查询 [购物车列表] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ids := make([]int32, 0)
	for _, v := range rsp.Data {
		ids = append(ids, v.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"total": 0,
		})
		return
	}
	//请求商品服务，获得商品信息
	goodsRsp, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		zap.S().Error("[List] 查询 [购物车列表] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	var data []gin.H
	for _, v := range rsp.Data {
		for _, goods := range goodsRsp.Data {
			if goods.Id == v.GoodsId {
				data = append(data, gin.H{
					"id":          v.Id,
					"nums":        v.Id,
					"checked":     v.Checked,
					"goods_id":    goods.Id,
					"goods_name":  goods.Name,
					"goods_image": goods.Images,
					"goods_price": goods.ShopPrice,
				})
			}
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  data,
	})
}

//添加商品到购物车
func New(ctx *gin.Context) {
	shopCartItemForm := forms.ShopCartItemForm{}
	err := ctx.ShouldBindJSON(&shopCartItemForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}
	//检查商品是否存在
	_, err = global.GoodsSrvClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: shopCartItemForm.GoodsId},
	)
	if err != nil {
		zap.S().Error("[New] 查询 [商品信息] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	inventoryDetail, err := global.InventorySrvClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: shopCartItemForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("[New] 查询 [库存信息] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	if inventoryDetail.Num < shopCartItemForm.Nums {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"nums": "库存不足",
		})
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.OrderSrvClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: shopCartItemForm.GoodsId,
		Nums:    shopCartItemForm.Nums,
	})
	if err != nil {
		zap.S().Error("[New] 添加 [添加到购物车] 失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	shopCartItemUpdateForm := forms.ShopCartItemUpdateForm{}
	err = ctx.ShouldBindJSON(&shopCartItemUpdateForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")

	request := proto.CartItemRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: int32(goodsId),
		Nums: shopCartItemUpdateForm.Nums,
	}
	if shopCartItemUpdateForm.Checked != nil {
		request.Checked = *shopCartItemUpdateForm.Checked
	}
	_, err = global.OrderSrvClient.UpdateCartItem(context.Background(), &request)
	if err != nil {
		zap.S().Error("[Update] 更新 [更新购物车记录] 失败")
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	goodsId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	userId, _ := ctx.Get("userId")

	_, err = global.OrderSrvClient.DeleteCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: int32(goodsId),
	})
	if err != nil {
		zap.S().Error("[Delete] 删除 [删除购物车记录] 失败")
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
