package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mxshop_api/order_web/api"
	"mxshop_api/order_web/api/pay"
	"mxshop_api/order_web/forms"
	"mxshop_api/order_web/global"
	"mxshop_api/order_web/models"
	"mxshop_api/order_web/proto"
)

//获得订单列表
func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")

	pages, _ := strconv.Atoi(ctx.DefaultQuery("p", "0"))
	perNum, _ := strconv.Atoi(ctx.DefaultQuery("pnum", "0"))

	model := claims.(*models.CustomClaims)
	request := proto.OrderFilterRequest{
		Pages:       int32(pages),
		PagePerNums: int32(perNum),
	}
	if model.AuthorityId == 1 {
		//是普通用户
		request.UserId = int32(userId.(uint32))
	}

	rsp, err := global.OrderSrvClient.OrderList(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		zap.S().Error("[List] 查询 [订单列表] 失败")
		return
	}

	var data []gin.H

	for _, order := range rsp.Data {
		data = append(data, gin.H{
			"id":order.Id,
			"status":order.Status,
			"pay_type":order.PayType,
			"user":order.UserId,
			"post":order.Post,
			"total":order.Total,
			"address":order.Address,
			"name":order.Name,
			"mobile":order.Mobile,
			"orderSn":order.OrderSn,
			"addTime":order.AddTime,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  data,
	})
}

//新建订单
func New(ctx *gin.Context) {
	createOrderForm := forms.CreateOrderForm{}
	err := ctx.ShouldBindJSON(&createOrderForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")

	rsp, err := global.OrderSrvClient.CreateOrder(context.WithValue(context.Background(),"ginCtx",ctx), &proto.OrderRequest{
		UserId:  int32(userId.(uint32)),
		Address: createOrderForm.Address,
		Name:    createOrderForm.Name,
		Mobile:  createOrderForm.Mobile,
		Post:    createOrderForm.Post,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		zap.S().Error("[New] 新建 [新建订单] 失败")
		return
	}

	//返回支付url
	url,err := pay.GeneratePayUrl(rsp.OrderSn,rsp.Total)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,"返回支付url失败")
		zap.S().Error("[New] 新建 [支付url] 失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
		"pay_url":url,
	})
}

//订单详情
func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	claims, _ := ctx.Get("claims")
	userId, _ := ctx.Get("userId")

	orderId, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "url格式出错",
		})
		return
	}

	request := proto.OrderRequest{
		Id:      int32(orderId),
		UserId:  0,
	}
	model := claims.(*models.CustomClaims)
	if model.AuthorityId == 1 {
		//是普通用户
		request.UserId = int32(userId.(uint32))
	}
	rsp,err := global.OrderSrvClient.OrderDetail(context.Background(),&request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		zap.S().Error("[Detail] 查询 [订单详情] 失败")
		return
	}

	var goods []interface{}
	for _,v := range rsp.Goods {
		goods  = append(goods,gin.H{
			"id":v.GoodsId,
			"name":v.GoodsName,
			"price":v.GoodsPrice,
			"image":v.GoodsImage,
			"nums":v.Nums,
		})
	}
	//返回支付url
	url,err := pay.GeneratePayUrl(rsp.OrderInfo.OrderSn,rsp.OrderInfo.Total)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,"返回支付url失败")
		zap.S().Error("[Detail] 新建 [支付url] 失败")
	}

	ctx.JSON(http.StatusOK,gin.H{
		"id":rsp.OrderInfo.Id,
		"status":rsp.OrderInfo.Status,
		"userId":rsp.OrderInfo.UserId,
		"post":rsp.OrderInfo.Post,
		"total":rsp.OrderInfo.Total,
		"address":rsp.OrderInfo.Address,
		"addtime":rsp.OrderInfo.AddTime,
		"name":rsp.OrderInfo.Name,
		"mobile":rsp.OrderInfo.Mobile,
		"pay_type":rsp.OrderInfo.PayType,
		"order_sn":rsp.OrderInfo.OrderSn,
		"goods":goods,
		"pay_url":url,
	})
}
