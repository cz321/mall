package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/goods_srv/proto"
)
var (
	brandClient proto.GoodsClient
	conn *grpc.ClientConn
)

func Init() {
	var err error
	conn,err = grpc.Dial("192.168.1.129:14128",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	brandClient = proto.NewGoodsClient(conn)
}
func testGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _,brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()
	defer conn.Close()

	testGetBrandList()
}