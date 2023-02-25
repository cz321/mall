package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/goods_srv/proto"
)

//品牌和轮播图

//获得品牌列表
func (s *GoodsServer)BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	var brands []model.Brands
	//分页查询
	result := global.DB.Scopes(Paginate(int(req.Pages),int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return nil,result.Error
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)

	response := &proto.BrandListResponse{}
	response.Total = int32(total)

	var brandResponses []*proto.BrandInfoResponse
	for _,brand := range brands {
		brandResponses = append(brandResponses,&proto.BrandInfoResponse{
			Id: brand.ID,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	response.Data = brandResponses
	return response,nil
}

//新建品牌
func (s *GoodsServer)CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	result := global.DB.First(&model.Brands{})
	if result.RowsAffected != 0 {
		return nil,status.Errorf(codes.InvalidArgument,"品牌已存在")
	}

	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Save(brand)
	return &proto.BrandInfoResponse{
		Id:   brand.ID,
		Name: req.Name,
		Logo: req.Logo,
	},nil
}

//删除品牌
func (s *GoodsServer)DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	result := global.DB.Delete(&model.Brands{},req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"品牌不存在")
	}
	return &empty.Empty{},nil
}

//更新品牌
func (s *GoodsServer)UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*empty.Empty, error) {
	brands := &model.Brands{}

	result := global.DB.First(brands,req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"品牌不存在")
	}
	if req.Name != "" {
		brands.Name = req.Name
	}
	if req.Logo != "" {
		brands.Name = req.Logo
	}
	global.DB.Save(brands)
	return &empty.Empty{},nil
}


