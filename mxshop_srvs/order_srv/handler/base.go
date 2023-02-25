package handler

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

//分页
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//订单号的生成
func GenerateOrderSn(userId int32) string{
	/*
	规则: 年月日时分秒 + userId + 2位随机数
	 */
	now := time.Now()
	rand.Seed(now.UnixNano())
	fmt.Println(now)
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Nanosecond(),
		userId,
		rand.Intn(90) + 10,
	)
	return orderSn
}