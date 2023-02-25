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
	conn,err = grpc.Dial("192.168.1.129:12139",grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	GoodsClient = proto.NewGoodsClient(conn)
}

func testGoodsList() {
	rsp, err := GoodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		TopCategory: 130361,
		KeyWords: "鱼",
		PriceMin: 90,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _,good := range rsp.Data {
		fmt.Println(good.Name,good.ShopPrice)
	}
}

func testBatchGetGoods() {
	rsp, err := GoodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: []int32{421,422,423},
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _,good := range rsp.Data {
		fmt.Println(good.Name)
	}
}
func testGetGoodsDetail() {
	rsp, err := GoodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{
		Id: 421,
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func main() {
	Init()
	defer conn.Close()

	//testGoodsList()

	//testBatchGetGoods()

	testGetGoodsDetail()
}

