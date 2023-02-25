package handler

import (
	"context"
	"crypto/sha512"
	"fmt"
	"strings"
	"time"

	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/model"
	"mxshop_srvs/user_srv/proto"
)

type UserServer struct{}

//获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	zap.S().Info("获取用户列表")
	var users []model.User
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp := &proto.UserListResponse{
		Total: int32(result.RowsAffected),
	}

	//分页
	global.DB.Scopes(Paginate(int(req.Pn),int(req.PSize))).Find(&users)

	for _,user := range users {
		userInfoResponse := ModelToResponse(user)
		rsp.Data = append(rsp.Data,userInfoResponse)
	}
	return rsp,nil
}

//通过mobile查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)

	if result.RowsAffected == 0 {
		return nil,status.Error(codes.NotFound,"用户不存在")
	}
	if result.Error != nil {
		return nil,result.Error
	}

	return  ModelToResponse(user),nil
}

//通过id查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user,req.Id)

	if result.RowsAffected == 0 {
		return nil,status.Error(codes.NotFound,"用户不存在")
	}
	if result.Error != nil {
		return nil,result.Error
	}

	return  ModelToResponse(user),nil
}

//添加用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)

	if result.RowsAffected == 1 {
		return nil,status.Error(codes.AlreadyExists,"用户已存在")
	}
	user.Mobile = req.Mobile
	user.NickName = req.NickName

	user.Role = 1		//默认是普通用户

	//加密
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt,encodePwd := password.Encode(req.Password,options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodePwd)

	result = global.DB.Create(&user)
	if result.Error != nil {
		return nil,status.Error(codes.Internal,result.Error.Error())
	}

	return  ModelToResponse(user),nil
}
//更新
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user,req.Id)

	if result.RowsAffected == 0 {
		return nil,status.Error(codes.NotFound,"用户不存在")
	}
	if result.Error != nil {
		return nil,result.Error
	}
	birthday := time.Unix(int64(req.Birthday), 0)
	user.NickName = req.NickName
	user.Birthday = &birthday
	user.Gender = req.Gender
	result = global.DB.Save(&user)
	if result.Error != nil {
		return nil,status.Error(codes.Internal,result.Error.Error())
	}

	return  &empty.Empty{},nil
}
//密码校验
func (s *UserServer) CheckPassword(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	pwdInfo := strings.Split(req.EncryptedPassword,"$")
	check := password.Verify(req.Password,pwdInfo[2],pwdInfo[3],options)
	return &proto.CheckResponse{Success: check},nil
}

//分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func ModelToResponse(user model.User) *proto.UserInfoResponse {
	userInfoResponse := &proto.UserInfoResponse{
		Id :user.ID,
		Password: user.Password,
		NickName: user.NickName,
		Gender: user.Gender,
		Role: int32(user.Role),
		Mobile: user.Mobile,
	}
	if user.Birthday != nil {
		userInfoResponse.Birthday = uint64(user.Birthday.Unix())
	}
	return userInfoResponse
}















