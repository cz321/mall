package api

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
)

func SendSms(ctx *gin.Context) {
	//表单验证
	senSmsForm := &forms.SendSmsForm{}
	err := ctx.ShouldBindJSON(senSmsForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	mobile := senSmsForm.Mobile
	smsCode := generateSmsCode(4)

	config := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential(global.ServerConfig.AliSmsInfo.AccessKeyId, global.ServerConfig.AliSmsInfo.AccessKeySecret)
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := dysmsapi.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		panic(err)
	}

	request := dysmsapi.CreateSendSmsRequest()

	request.Scheme = "https"
	request.SignName = "阿里云短信测试"
	request.TemplateCode = "SMS_154950909"
	request.PhoneNumbers = mobile

	request.TemplateParam = "{\"code\":" + smsCode + "}"

	_, err = client.SendSms(request)
	if err != nil {
		zap.S().Error(err.Error())
	}
	//fmt.Printf("response is %#v\n", response)

	//将验证码和手机号保存在数据库
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	//rdb.Set(context.Background(),mobile,smsCode,2 * time.Minute)
	rdb.Set(context.Background(), mobile, smsCode, time.Duration(global.ServerConfig.AliSmsInfo.Expire)*time.Second) //设置不过期

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "验证码发送成功",
	})
}

//生成width长度的短信新验证码
func generateSmsCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().Unix())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
