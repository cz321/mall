package main

import (
	"context"
	"fmt"

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
	conn,err = grpc.Dial("192.168.1.129:9980",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	GoodsClient = proto.NewGoodsClient(conn)
}

func testCategoryBrandList() {
	rsp, err := GoodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{

	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}
func testGetCategoryBrandList() {
	rsp, err := GoodsClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 135475,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}

func main() {
	Init()
	defer conn.Close()

	//testCategoryBrandList()

	testGetCategoryBrandList()
}
