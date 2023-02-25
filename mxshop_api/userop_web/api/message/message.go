package message

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/models"
	"mxshop_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	id,_ := ctx.Get("userId")
	claims,_ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)
	var userId int32
	if model.AuthorityId == 1 {
		userId = int32(id.(uint32))
	}

	request := proto.MessageRequest{
		UserId:      userId,
	}
	rsp, err := global.MessageSrvClient.MessageList(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}
	var data []gin.H

	for _,v := range rsp.Data {
		data = append(data,gin.H{
			"id":v.Id,
			"user_id":v.UserId,
			"type":v.MessageType,
			"subject":v.Subject,
			"message":v.Message,
			"file":v.File,
		})
	}
	ctx.JSON(http.StatusOK,gin.H{
		"total":rsp.Total,
		"data": data,
	})
}
func New(ctx *gin.Context) {
	var messageForm forms.MessageForm
	err := ctx.ShouldBindJSON(&messageForm)
	if err != nil {
		api.HandleValidatorError(ctx,err)
	}

	userId,_ := ctx.Get("userId")

	request := proto.MessageRequest{
		UserId:      int32(userId.(uint32)),
		MessageType: messageForm.MessageType,
		Subject:     messageForm.Subject,
		Message:     messageForm.Message,
		File:        messageForm.File,
	}
	rsp, err := global.MessageSrvClient.CreateMessage(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err,ctx)
		return
	}
	ctx.JSON(http.StatusOK,gin.H{
		"id":rsp.Id,
	})
}