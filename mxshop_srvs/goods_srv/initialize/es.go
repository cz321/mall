package initialize

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
)

func InitEs() {
	url := fmt.Sprintf("http://%s:%d",global.ServerConfig.EsInfo.Host,global.ServerConfig.EsInfo.Port)

	logger := log.New(os.Stdout,"es:",log.LstdFlags)

	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false),elastic.SetTraceLog(logger))
	if err != nil {
		zap.S().Error("ES初始化错误",err)
		panic(err)
	}

	isExit, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		zap.S().Error("ES初始化错误",err)
		panic(err)
	}
	if !isExit {
		_,err = global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			zap.S().Error("ES初始化错误",err)
			panic(err)
		}
	}
}