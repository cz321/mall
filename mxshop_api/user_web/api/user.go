package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_api/user_web/forms"
	"mxshop_api/user_web/global"
	"mxshop_api/user_web/global/response"
	"mxshop_api/user_web/middlewares"
	"mxshop_api/user_web/models"
	"mxshop_api/user_web/proto"
)

//获取用户列表
func GetUserList(ctx *gin.Context) {
	//取出jwt时候传入的用户身份信息
	claims, _ := ctx.Get("claims")
	curUser := claims.(*models.CustomClaims)
	zap.S().Infof("用户%d访问",curUser.ID)

	pn, _ := strconv.Atoi(ctx.DefaultQuery("pn", "0"))
	pSize, _ := strconv.Atoi(ctx.DefaultQuery("pSize", "0"))
	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pn),
		PSize: uint32(pSize),
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【查询用户列表失败】")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	result := make([]interface{}, 0)
	for _, v := range rsp.Data {
		user := response.UserResponse{
			ID:       v.Id,
			NickName: v.NickName,
			Gender:   v.Gender,
			Birthday: response.JsonTime(time.Unix(int64(v.Birthday), 0)),
			Mobile:   v.Mobile,
		}
		result = append(result, user)
	}
	ctx.JSON(http.StatusOK, result)
}

//登录
func PasswordLogin(ctx *gin.Context) {
	//表单验证
	passwordLoginForm := &forms.PasswordLoginForm{}
	err := ctx.ShouldBindJSON(passwordLoginForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	//验证码验证
	ok := store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, false)
	if !ok {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"captcha":"验证码错误",
		})
		return
	}

	//登录
	rsp, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"mobile": "用户不存在",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"mobile": "登录失败",
				})
			}
			return
		}
		return
	}

	//核验密码
	passRsp, passErr := global.UserSrvClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
		Password:          passwordLoginForm.Password,
		EncryptedPassword: rsp.Password,
	})
	if passErr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"password": "登录失败",
		})
		return
	}

	if passRsp.Success {
		//生成token
		j := middlewares.NewJWT()
		claims := models.CustomClaims{
			ID:          uint32(rsp.Id),
			NickName:    rsp.NickName,
			AuthorityId: uint32(rsp.Role),
			StandardClaims: jwt.StandardClaims{
				NotBefore: time.Now().Unix(),               //签名生效时间
				ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
				Issuer:    "imooc",
			},
		}
		token,err := j.CreateToken(claims)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"msg": "生成token失败",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"id": rsp.Id,
			"nickname": rsp.NickName,
			"token": token,
			"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
		})
	}else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "密码错误",
		})
	}
}

//注册
func Register(ctx *gin.Context) {
	//表单验证
	registerForm := &forms.RegisterFrom{}
	err := ctx.ShouldBindJSON(registerForm)
	if err != nil {
		HandleValidatorError(ctx, err)
		return
	}
	//验证码校验
	//将验证码和手机号保存在数据库
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	v,err := rdb.Get(context.Background(), registerForm.Mobile).Result()

	if err != nil {
		if err == redis.Nil {
			ctx.JSON(http.StatusBadRequest,gin.H{
				"msg":"key不存在",
			})
		}else {
			ctx.JSON(http.StatusInternalServerError,gin.H{
				"msg":"内部错误",
			})
		}
		return
	}

	if v != registerForm.Code {
		ctx.JSON(http.StatusBadRequest,gin.H{
			"msg":"验证码错误",
		})
		return
	}

	rsp, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		Mobile:   registerForm.Mobile,
		NickName: registerForm.Mobile,
		Password: registerForm.Password,
	})

	if err != nil {
		zap.S().Errorf("[NewUserClient] 新建 【新建用户失败】 失败 ：%s",err.Error())
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	//生成token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint32(rsp.Id),
		NickName:    rsp.NickName,
		AuthorityId: uint32(rsp.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),               //签名生效时间
			ExpiresAt: time.Now().Unix() + 60*60*24*30, //30天过期
			Issuer:    "imooc",
		},
	}
	token,err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.Id,
		"nickname": rsp.NickName,
		"token": token,
		"expired_at": (time.Now().Unix() + 60*60*24*30)*1000,
	})
}

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fields {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

//将 grpc 的 code 转换成 http 的状态码
func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				ctx.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				ctx.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其他错误" + e.Message(),
				})
			}
			return
		}
	}
}

//处理 Validator 错误
func HandleValidatorError(ctx *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	ctx.JSON(http.StatusBadRequest, gin.H{
		"err": removeTopStruct(errs.Translate(global.Trans)),
	})
}
