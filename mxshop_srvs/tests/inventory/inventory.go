package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/inventory_srv/proto"
)

var (
	brandClient proto.InventoryClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.1.129:5555", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	brandClient = proto.NewInventoryClient(conn)
}
func testSetInv() {
	_, err := brandClient.SetInv(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
		Num:     1001,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("设置库存成功")
}
func testInvDetail() {
	rsp, err := brandClient.InvDetail(context.Background(), &proto.GoodsInvInfo{
		GoodsId: 421,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}
func testSell(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := brandClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 421, Num: 1},
			{GoodsId: 422, Num: 1},
		},
	})
	if err != nil {
		fmt.Println("库存扣减失败",err)
		return
	}
	fmt.Println("库存扣减成功")
}
func testSell1(wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := brandClient.Sell(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{
			{GoodsId: 422, Num: 1},
			{GoodsId: 421, Num: 1},
		},
	})
	if err != nil {
		fmt.Println("库存扣减失败",err)
		return
	}
	fmt.Println("库存扣减成功")
}
func testReback() {
	_, err := brandClient.Reback(context.Background(), &proto.SellInfo{
		GoodsInfo: []*proto.GoodsInvInfo{{GoodsId: 1000, Num: 1}, {GoodsId: 422, Num: 1}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("库存归还成功")
}

func main() {
	Init()
	defer conn.Close()

	//testInvDetail()
	//testSetInv()
	//testInvDetail()

	//testSell()
	//testReback()

	var wg sync.WaitGroup
	wg.Add(20)
	now := time.Now()
	for i := 0; i < 10; i++ {
		go testSell1(&wg)
		go testSell(&wg)
	}
	wg.Wait()
	fmt.Println(time.Since(now).Microseconds())
}


