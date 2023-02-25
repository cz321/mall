package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/userop_srv/proto"
)

var (
	userFavClient proto.UserFavClient
	MessageClient proto.MessageClient
	AddressClient proto.AddressClient
	conn        *grpc.ClientConn
)

func Init() {
	var err error
	conn, err = grpc.Dial("192.168.1.129:7777", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		panic(err)
	}

	userFavClient = proto.NewUserFavClient(conn)
	//messageClient = proto.NewMessageClient(conn)
	//addressClient = proto.NewAddressClient(conn)
}
func testCreateAddress() {
	_, err := AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       4,
		Province:     "陕西",
		City:         "西安",
		District:     "长安区",
		Address:      "东大街道",
		SignerName:   "cz",
		SignerMobile: "18209285903",
	})
	if err != nil {
		panic(err)
	}
}


func main() {
	Init()
	defer conn.Close()

	testCreateAddress()
}


