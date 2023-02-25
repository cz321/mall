package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/user_srv/proto"
)
var (
	userClient proto.UserClient
	conn *grpc.ClientConn
)

func Init() {
	var err error
	conn,err = grpc.Dial("localhost:50051",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	userClient = proto.NewUserClient(conn)
}
func testGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _,user := range rsp.Data {
		fmt.Println(user.Mobile,user.NickName,user.Password)
		checkRsp, err := userClient.CheckPassword(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.Password,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func main() {
	Init()
	defer conn.Close()

	testGetUserList()
}