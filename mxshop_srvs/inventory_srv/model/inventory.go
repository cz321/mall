package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Inventory struct {
	BaseModel
	Goods   int32 `gorm:"type:int;index" json:"goods"`
	Stocks  int32 `gorm:"type:int" json:"stocks"`
	Version int32 `gorm:"type:int" json:"version"`
}

//type InventoryHistory struct {
//	User int32
//	Goods int32
//	Nums int32
//	OrderSn int32
//	Status int32 //1.表示库存是预扣减,2.表示已经支付
//}

type StockSellDetail struct {
	OrderSn string `gorm:"type:varchar(200);index:idx_order_sn,unique;"`
	Status int32 `gorm:"type:varchar(200)"` //1 表示已扣减 2. 表示已归还
	Detail GoodsDetailList `gorm:"type:varchar(200)"`
}

func (StockSellDetail) TableName() string {
	return "stockselldetail"
}

type GoodsDetail struct {
	Goods int32
	Num int32
}
type GoodsDetailList []GoodsDetail

func (g GoodsDetailList) Value() (driver.Value, error){
	return json.Marshal(g)
}

// 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (g *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &g)
}