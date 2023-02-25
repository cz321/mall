package handler

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/model"
	"mxshop_srvs/userop_srv/proto"
)

func (*UserOpServer) GetAddressList(ctx context.Context, req *proto.AddressRequest) (*proto.AddressListResponse, error) {
	var addresses []model.Address
	result := global.DB.Where(&model.Address{User:req.UserId}).Find(&addresses)

	var data []*proto.AddressResponse
	for _,address := range addresses {
		data = append(data,&proto.AddressResponse{
			Id:           address.ID,
			UserId:       address.User,
			Province:     address.Province,
			City:         address.City,
			District:     address.District,
			Address:      address.Address,
			SignerName:   address.SignerName,
			SignerMobile: address.SignerMobile,
		})
	}
	return &proto.AddressListResponse{
		Total: int32(result.RowsAffected),
		Data:  data,
	},nil
}

func (*UserOpServer) CreateAddress(ctx context.Context, req *proto.AddressRequest) (*proto.AddressResponse, error) {
	var address = model.Address{
		User:         req.UserId,
		Province:     req.Province,
		City:         req.City,
		District:     req.District,
		Address:      req.Address,
		SignerName:   req.SignerName,
		SignerMobile: req.SignerMobile,
	}
	result := global.DB.Save(&address)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.AlreadyExists, "收货地址已存在")
	}
	return &proto.AddressResponse{Id: address.ID},nil
}

func (*UserOpServer) DeleteAddress(ctx context.Context, req*proto.AddressRequest) (*empty.Empty, error) {
	result := global.DB.Delete(&model.Address{}, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	return &empty.Empty{},nil
}

func (*UserOpServer) UpdateAddress(ctx context.Context, req*proto.AddressRequest) (*empty.Empty, error) {
	var address model.Address
	result := global.DB.Where("id = ? and user = ?",req.Id,req.UserId).First(&address)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "收货地址不存在")
	}
	if req.Province != "" {
		address.Province = req.Province
	}
	if req.City != "" {
		address.Province = req.City
	}
	if req.District != "" {
		address.Province = req.District
	}
	if req.Address != "" {
		address.Province = req.Address
	}
	if req.SignerMobile != "" {
		address.Province = req.SignerMobile
	}
	if req.SignerName != "" {
		address.Province = req.SignerName
	}
	global.DB.Save(&address)
	return &empty.Empty{},nil
}