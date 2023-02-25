package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	"mxshop_api/goods_web/api"
	"mxshop_api/goods_web/forms"
	"mxshop_api/goods_web/global"
	"mxshop_api/goods_web/proto"
)

func List(ctx *gin.Context) {
	rsp, err := global.GoodsSrvClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(rsp.JsonData), &data)
	if err != nil {
		zap.S().Error("[list] 查询失败", err.Error())
	}
	ctx.JSON(http.StatusOK, data)
}

func New(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	err := ctx.ShouldBindJSON(categoryForm)
	if err != nil {
		api.HandleValidatorError(ctx,err)
		return
	}
	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,
		ParentCategory: categoryForm.ParentCategory,
		Level:          categoryForm.Level,
		IsTab:          *categoryForm.IsTab,
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
	}

	ctx.JSON(http.StatusOK,gin.H{
		"id":rsp.Id,
		"name":rsp.Name,
		"parent":rsp.ParentCategory,
		"level":rsp.Level,
		"is_tab":rsp.IsTab,
	})
}

func Update(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	err := ctx.ShouldBindJSON(categoryForm)
	if err != nil {
		api.HandleValidatorError(ctx,err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),
		Name: categoryForm.Name,
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab
	}
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
	}

	ctx.Status(http.StatusOK)

}

func Detail(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	rsp, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}

	subCategories := make([]interface{},0)
	for _,v := range rsp.SubCategorys {
		subCategories = append(subCategories, gin.H{
			"id":v.Id,
			"name":v.Name,
			"level":v.Level,
			"parent_category":v.ParentCategory,
			"is_tab":v.IsTab,
		})
	}

	ctx.JSON(http.StatusOK,gin.H{
		"id":              rsp.Info.Id,
		"name":            rsp.Info.Name,
		"level":           rsp.Info.Level,
		"parent_category": rsp.Info.ParentCategory,
		"is_tab":          rsp.Info.IsTab,
		"sub_category":    subCategories,
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i,err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{
		Id: int32(i),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
	}

	ctx.Status(http.StatusOK)
}
