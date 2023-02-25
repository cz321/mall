package banner

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
)

//轮播图

func List(ctx *gin.Context) {
	rsp, err := global.GoodsSrvClient.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	res := make([]interface{}, 0)

	for _, v := range rsp.Data {
		res = append(res, gin.H{
			"id":    v.Id,
			"index": v.Index,
			"image": v.Image,
			"url":   v.Url,
		})
	}
	ctx.JSON(http.StatusOK, res)
}

func New(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	err := ctx.ShouldBindJSON(&bannerForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: bannerForm.Index,
		Image: bannerForm.Image,
		Url:   bannerForm.Url,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":    rsp.Id,
		"index": rsp.Index,
		"image": rsp.Image,
		"url":   rsp.Url,
	})
}

func Update(ctx *gin.Context) {
	bannerForm := forms.BannerForm{}
	err := ctx.ShouldBindJSON(&bannerForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),
		Index: bannerForm.Index,
		Url:   bannerForm.Url,
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
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"msg":"删除成功",
	})
}
