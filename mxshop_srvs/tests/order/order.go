package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/order_srv/proto"
)

var (
	orderClient proto.OrderClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.1.129:4444", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	orderClient = proto.NewOrderClient(conn)
}

func testCreateCartItem(userId, goodsId, nums int32) {
	rsp, err := orderClient.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  userId,
		GoodsId: goodsId,
		Nums:    nums,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func testCartItemList(userId int32) {
	rsp, err := orderClient.CartItemList(context.Background(), &proto.UserInfo{
		Id: userId,
	})
	if err != nil {
		panic(err)
	}
	for _, v := range rsp.Data {
		fmt.Println(v.UserId, v.GoodsId, v.GoodsId, v.Checked)
	}
}

func testUpdateCartItem(id int32) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("update success")
}

func testCreateOrder() {
	rsp, err := orderClient.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  2,
		Address: "东大村",
		Name:    "cz",
		Mobile:  "18209285903",
		Post:    "i am post",
	})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Id)
}

func testOrderDetail(orderId int32) {
	rsp,err := orderClient.OrderDetail(context.Background(), &proto.OrderRequest{
		Id: orderId,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.OrderInfo)

	for _,good := range rsp.Goods {
		fmt.Println(good)
	}
}

func testOrderList() {
	rsp,err := orderClient.OrderList(context.Background(),&proto.OrderFilterRequest{UserId: 1})

	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)

	for _,v := range rsp.Data {
		fmt.Println(v.OrderSn)
	}
}
func main() {
	Init()
	defer conn.Close()

	//testCreateCartItem(2,422,2)

	//testCartItemList(2)
	//testUpdateCartItem(1)
	//testCartItemList(2)

	//testCreateOrder()

	//testOrderDetail(2)

	testOrderList()
}

