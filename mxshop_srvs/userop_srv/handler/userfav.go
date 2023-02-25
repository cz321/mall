package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/model"
	"mxshop_srvs/userop_srv/proto"
)

func (*UserOpServer) GetFavList(ctx context.Context, req *proto.UserFavRequest) (*proto.UserFavListResponse, error) {
	var userFavs []model.UserFav
	//1.查询用户的收藏记录
	//2.查询某件商品被哪些用户收藏了
	result := global.DB.Where(&model.UserFav{User: req.UserId,Goods: req.GoodsId}).Find(&userFavs)

	var data []*proto.UserFavResponse
	for _,userFav := range userFavs {
		data = append(data,&proto.UserFavResponse{
			UserId:  userFav.User,
			GoodsId: userFav.Goods,
		})
	}
	return &proto.UserFavListResponse{
		Total: int32(result.RowsAffected),
		Data:  data,
	},nil
}

func (*UserOpServer) AddUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	result := global.DB.Save(&model.UserFav{
		User:  req.UserId,
		Goods: req.GoodsId,
	})

	if result.RowsAffected == 0 {
		return nil,status.Error(codes.AlreadyExists,"记录已存在")
	}
	return &empty.Empty{},nil
}

func (*UserOpServer) DeleteUserFav(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	if result := global.DB.Unscoped().Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.UserFav{}); result.RowsAffected == 0{
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &emptypb.Empty{}, nil
}

func (*UserOpServer) GetUserFavDetail(ctx context.Context, req *proto.UserFavRequest) (*empty.Empty, error) {
	var row int64
	global.DB.Model(&model.UserFav{}).Where("goods=? and user=?", req.GoodsId, req.UserId).Count(&row)
	if row == 0 {
		return nil, status.Errorf(codes.NotFound, "记录不存在")
	}
	return &emptypb.Empty{}, nil
}