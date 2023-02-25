package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/olivere/elastic/v7"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/goods_srv/model"
)


func main() {
	//dsn := "root:root@tcp(localhost:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	Logger: logger.Default.LogMode(logger.Info),
	//	NamingStrategy: schema.NamingStrategy{
	//		SingularTable: true,
	//	},
	//})
	//if err != nil {
	//	panic(err)
	//}
	////迁移
	//_ = db.AutoMigrate(&model.Category{},&model.Brands{},&model.GoodsCategoryBrand{},&model.Banner{},&model.Goods{})

	mysql2Es()
}

func mysql2Es()  {
	dsn := "root:root@tcp(localhost:3306)/mxshop_goods_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}

	url := "http://192.168.10.105:9200"
	logger := log.New(os.Stdout,"es:",log.LstdFlags)
	client, _ := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false),elastic.SetTraceLog(logger))

	var goods []model.Goods
	db.Find(&goods)
	for _,g := range goods {
		esGoods := model.EsGoods{
			ID:          g.ID,
			CategoryId:  g.CategoryId,
			BrandsId:    g.BrandsId,
			OnSale:      g.OnSale,
			ShipFree:    g.ShipFree,
			IsNew:       g.IsNew,
			IsHot:       g.IsHot,
			Name:        g.Name,
			ClickNum:    g.ClickNum,
			SoldNum:     g.SoldNum,
			FavNum:      g.FavNum,
			MarketPrice: g.MarketPrice,
			ShopPrice:   g.ShopPrice,
			GoodsBrief:  g.GoodsBrief,
		}
		_, err = client.Index().Index(esGoods.GetIndexName()).BodyJson(esGoods).Id(strconv.Itoa(int(g.ID))).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}