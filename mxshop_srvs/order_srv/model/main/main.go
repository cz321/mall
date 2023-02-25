package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/order_srv/model"
)


func main() {
	dsn := "root:root@tcp(localhost:3306)/mxshop_order_srv?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	//迁移
	_ = db.AutoMigrate(&model.ShoppingCart{},&model.OrderInfo{},&model.OrderGoods{})
}
