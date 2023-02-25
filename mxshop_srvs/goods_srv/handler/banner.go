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

//轮播图

func (s *GoodsServer)BannerList(context.Context, *empty.Empty) (*proto.BannerListResponse, error) {
	var banners []model.Banner
	result := global.DB.Find(&banners)
	if result.Error != nil {
		return nil,result.Error
	}
	var total int64
	global.DB.Model(&model.Banner{}).Count(&total)

	response := &proto.BannerListResponse{}
	response.Total = int32(total)

	var bannerResponses []*proto.BannerResponse
	for _,banner := range banners {
		bannerResponses = append(bannerResponses,&proto.BannerResponse{
			Id: banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url: banner.Url,
		})
	}
	response.Data = bannerResponses
	return response,nil
}

func (s *GoodsServer)CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := &model.Banner{
		Image: req.Image,
		Index: req.Index,
		Url: req.Url,
	}
	global.DB.Save(banner)
	return &proto.BannerResponse{
		Id: banner.ID,
		Image: banner.Image,
		Url: banner.Url,
		Index: banner.Index,
	},nil
}
func (s *GoodsServer)DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	result := global.DB.Delete(&model.Banner{},req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"轮播图不存在")
	}
	return &empty.Empty{},nil
}

func (s *GoodsServer)UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*empty.Empty, error) {
	banner := &model.Banner{}

	result := global.DB.First(banner,req.Id)
	if result.RowsAffected == 0 {
		return nil,status.Errorf(codes.NotFound,"轮播图不存在")
	}
	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}
	global.DB.Save(banner)
	return &empty.Empty{},nil
}