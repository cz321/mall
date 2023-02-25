package user_fav

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	userId, _ := ctx.Get("userId")
	rsp, err := global.UserFavSrvClient.GetFavList(context.Background(), &proto.UserFavRequest{
		UserId: int32(userId.(uint32)),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	var ids []int32
	for _,v := range rsp.Data {
		ids = append(ids, v.GoodsId)
	}

	if len(ids) == 0 {
		ctx.JSON(http.StatusOK,gin.H{
			"total":0,
			"data":gin.H{},
		})
	}

	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: ids,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	var data []gin.H

	for _,v := range rsp.Data {
		for _,good := range goods.Data {
			if good.Id == v.GoodsId {
				data = append(data,gin.H{
					"id":good.Id,
					"name":good.Name,
					"shop_price":good.ShopPrice,
				})
			}
		}
	}
	ctx.JSON(http.StatusOK,gin.H{
		"total":rsp.Total,
		"data":data,
	})
}

func New(ctx *gin.Context) {
	var userFavForm forms.UserFavForm
	err := ctx.ShouldBindJSON(&userFavForm)
	if err != nil {
		zap.S().Error("form表单出错")
		api.HandleValidatorError(ctx, err)
		return
	}
	userId, _ := ctx.Get("userId")

	//查看goods是否存在
	_,err = global.GoodsSrvClient.GetGoodsDetail(context.Background(),&proto.GoodInfoRequest{Id: userFavForm.GoodsId})
	if err != nil {
		zap.S().Error("goods不存在")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	_, err = global.UserFavSrvClient.AddUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: userFavForm.GoodsId,
	})
	if err != nil {
		zap.S().Error("添加失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, _ := ctx.Get("userId")

	_, err = global.UserFavSrvClient.GetUserFavDetail(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, _ := ctx.Get("userId")

	_, err = global.UserFavSrvClient.DeleteUserFav(context.Background(), &proto.UserFavRequest{
		UserId:  int32(userId.(uint32)),
		GoodsId: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}