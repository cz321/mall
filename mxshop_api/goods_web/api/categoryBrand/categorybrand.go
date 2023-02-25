package categoryBrand

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
	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	resData := make([]interface{}, 0)
	for _, v := range rsp.Data {
		resData = append(resData, gin.H{
			"id":   v.Id,
			"category": gin.H{
				"id": v.Category.Id,
				"name": v.Category.Name,
			},
			"brand": gin.H{
				"id": v.Brand.Id,
				"name": v.Brand.Name,
				"logo": v.Brand.Logo,
			},
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data" : resData,
	})
}

func New(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	err := ctx.ShouldBindJSON(&categoryBrandForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(ctx *gin.Context) {
	categoryBrandForm := forms.CategoryBrandForm{}
	err := ctx.ShouldBindJSON(&categoryBrandForm)
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

	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
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

	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.Status(http.StatusOK)
}

func GetCategoryBrandList(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id:int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	res := make([]interface{}, 0)
	for _, v := range rsp.Data {
		res = append(res, gin.H{
			"id":   v.Id,
			"name": v.Name,
			"logo": v.Logo,
		})
	}

	ctx.JSON(http.StatusOK, res)
}