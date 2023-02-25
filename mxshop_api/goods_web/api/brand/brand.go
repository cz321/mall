package brand

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
)

func List(ctx *gin.Context) {
	pn := ctx.DefaultQuery("pn","0")
	pSize := ctx.DefaultQuery("psize","10")

	pnInt,err := strconv.Atoi(pn)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	pSizeInt,err := strconv.Atoi(pSize)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	resData := make([]interface{},0)
	for _,v := range rsp.Data[pnInt:pnInt*pSizeInt+pSizeInt] {
		resData = append(resData,gin.H{
			"id":v.Id,
			"name":v.Name,
			"logo":v.Logo,
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"total":rsp.Total,
		"data":resData,
	})
}

func New(ctx *gin.Context) {
	brandFrom := forms.BrandForm{}
	err := ctx.ShouldBindJSON(&brandFrom)
	if err != nil {
		api.HandleValidatorError(ctx,err)
		return
	}

	rsp,err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandFrom.Logo,
		Logo: brandFrom.Name,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	ctx.JSON(http.StatusOK,gin.H{
		"id":rsp.Id,
		"name":rsp.Name,
		"logo":rsp.Logo,
	})
}

func Update(ctx *gin.Context) {
	brandFrom := forms.BrandForm{}
	err := ctx.ShouldBindJSON(&brandFrom)
	if err != nil {
		api.HandleValidatorError(ctx,err)
		return
	}

	id := ctx.Param("id")
	i,err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_,err = global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Id: int32(i),
		Name: brandFrom.Logo,
		Logo: brandFrom.Name,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i,err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_,err = global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	ctx.Status(http.StatusOK)
}