package handler

import (
	"context"
	"encoding/json"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

//商品分类
func (s *GoodsServer)GetAllCategorysList(context.Context, *empty.Empty) (*proto.CategoryListResponse, error) {
	/*
	[
		"id",xxx,
		"name","",
		"level":1,
		"is_tab";false;
		"parent":13xxx,
		"sub_category":[
			"id",xxx,
			"name","",
			"level":1,
			"is_tab";false;
			"parent":13xxx,
		]
	]
	 */
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	bytes, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{
		JsonData: string(bytes),
	},nil
}
//获取子分类
func (s *GoodsServer)GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	var category model.Category
	result := global.DB.First(&category)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"商品分类不存在")
	}
	//preloads := ""
	//if category.Level == 1 {
	//	preloads = "SubCategory.SubCategory"
	//}else if category.Level == 2 {
	//	preloads = "SubCategory"
	//}
	var subCategorys []model.Category
	global.DB.Where(&model.Category{ParentCategoryId: req.Id}).Find(&subCategorys)

	subCategoryListRsp := &proto.SubCategoryListResponse{
		Info: &proto.CategoryInfoResponse{
			Id: category.ID,
			Name: category.Name,
			Level: category.Level,
			IsTab: category.IsTab,
			ParentCategory: category.ParentCategoryId,
		},
	}
	for _,subCategory := range subCategorys {
		subCategoryListRsp.SubCategorys = append(subCategoryListRsp.SubCategorys,&proto.CategoryInfoResponse{
			Id: subCategory.ID,
			Name: subCategory.Name,
			Level: subCategory.Level,
			IsTab: subCategory.IsTab,
			ParentCategory: subCategory.ParentCategoryId,
		})
	}
	return subCategoryListRsp,nil

}
func (s *GoodsServer)CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{
		Name:  req.Name,
		Level: req.Level,
	}
	if req.Level != 1 {
		category.ParentCategoryId = req.ParentCategory
	}
	category.IsTab = req.IsTab

	global.DB.Save(&category)

	return &proto.CategoryInfoResponse{
		Id: int32(category.ID),
	},nil
}
func (s *GoodsServer)DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*empty.Empty, error) {
	result := global.DB.Delete(&model.Category{}, req.Id)
	if result.RowsAffected == 0 {
		return &empty.Empty{},status.Error(codes.NotFound,"商品分类不存在")
	}
	return &empty.Empty{},nil
}
func (s *GoodsServer)UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*empty.Empty, error) {
	var category model.Category
	result := global.DB.Delete(&model.Category{}, req.Id)
	if result.RowsAffected == 0 {
		return &empty.Empty{},status.Error(codes.NotFound,"商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name
	}
	if req.Level != 0 {
		category.Level = req.Level
	}
	if req.IsTab {
		category.IsTab = req.IsTab
	}
	if req.ParentCategory != 0  {
		category.ParentCategoryId = req.ParentCategory
	}

	global.DB.Save(&category)
	return &empty.Empty{},nil
}