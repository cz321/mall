package main

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/goods_srv/proto"
)
var (
	GoodsClient proto.GoodsClient
	conn *grpc.ClientConn
)

func Init() {
	var err error
	conn,err = grpc.Dial("192.168.1.129:3916",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	GoodsClient = proto.NewGoodsClient(conn)
}
func testGetAllCategorysList() {
	rsp, err := GoodsClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.JsonData)
	//for _,category := range rsp.Data {
	//	fmt.Println(category.Name)
	//}
}
func testGetSubCategory() {
	rsp, err := GoodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: 135487,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Info)
	fmt.Println(rsp.SubCategorys)
}

func main() {
	Init()
	defer conn.Close()

	//testGetAllCategorysList()

	testGetSubCategory()
}