package handler

import (
	"context"

	"mxshop_srvs/userop_srv/global"
	"mxshop_srvs/userop_srv/model"
	"mxshop_srvs/userop_srv/proto"
)

func (*UserOpServer) MessageList(ctx context.Context, req *proto.MessageRequest) (*proto.MessageListResponse, error) {
	var messages []model.LeavingMessages
	result := global.DB.Where(&model.LeavingMessages{User: req.UserId}).Find(&messages)

	var data []*proto.MessageResponse
	for _,message := range messages {
		data = append(data,&proto.MessageResponse{
			Id:          message.ID,
			UserId:      message.User,
			MessageType: message.MessageType,
			Subject:     message.Subject,
			Message:     message.Message,
			File:        message.File,
		})
	}

	return &proto.MessageListResponse{
		Total: int32(result.RowsAffected),
		Data:  data,
	},nil
}

func (*UserOpServer) CreateMessage(ctx context.Context, req *proto.MessageRequest) (*proto.MessageResponse, error) {
	var message = model.LeavingMessages{
		User:        req.UserId,
		MessageType: req.MessageType,
		Subject:     req.Subject,
		Message:     req.Message,
		File:        req.File,
	}

	global.DB.Save(&message)

	return &proto.MessageResponse{Id: message.ID},nil
}
