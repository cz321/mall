package address

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mxshop_api/userop_web/api"
	"mxshop_api/userop_web/forms"
	"mxshop_api/userop_web/global"
	"mxshop_api/userop_web/models"
	"mxshop_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	id, _ := ctx.Get("userId")
	claims, _ := ctx.Get("claims")
	model := claims.(*models.CustomClaims)

	var userId int32
	if model.AuthorityId == 1 {
		userId = int32(id.(uint32))
	}

	request := proto.AddressRequest{
		UserId: userId,
	}
	rsp, err := global.AddressClient.GetAddressList(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	var data []gin.H

	for _, v := range rsp.Data {
		data = append(data, gin.H{
			"id":            v.Id,
			"user_id":       v.UserId,
			"province":      v.Province,
			"city":          v.City,
			"district":      v.District,
			"address":       v.Address,
			"signer_name":   v.SignerName,
			"signer_mobile": v.SignerMobile,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"total": rsp.Total,
		"data":  data,
	})
}

func New(ctx *gin.Context) {
	var addressForm forms.AddressForm
	err := ctx.ShouldBindJSON(&addressForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")

	request := proto.AddressRequest{
		UserId:       int32(userId.(uint32)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	}
	rsp, err := global.AddressClient.CreateAddress(context.Background(), &request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
	})
}

func Update(ctx *gin.Context) {
	var addressForm forms.AddressForm
	err := ctx.ShouldBindJSON(&addressForm)
	if err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.Atoi(id)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, _ := ctx.Get("userId")

	_, err = global.AddressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id: int32(i),
		UserId: int32(userId.(uint32)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
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
		ctx.Status(http.StatusBadRequest)
		return
	}
	userId, _ := ctx.Get("userId")

	_, err = global.AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{
		Id: int32(i),
		UserId: int32(userId.(uint32)),
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}
