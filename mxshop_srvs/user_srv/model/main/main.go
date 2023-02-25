package main

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"mxshop_srvs/user_srv/model"
)


func main()  {
	dsn := "root:root@tcp(localhost:3306)/mxshop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"
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
	//_ = db.AutoMigrate(&model.User{})

	//添加十条数据
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt,encodePwd := password.Encode("admin123",options)
	pwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodePwd)
	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("bobby%d",i),
			Mobile: fmt.Sprintf("1820928590%d",i),
			Password: pwd,
		}
		db.Save(&user)
	}
}
func jiayan(code string) {
	options := &password.Options{SaltLen: 16, Iterations: 100, KeyLen: 32, HashFunction: sha512.New}
	salt,encodePwd := password.Encode("code",options)
	pwd := fmt.Sprintf("$pbkdf2-sha512$%s$%s",salt,encodePwd)
	fmt.Println(salt)
	fmt.Println(pwd)
	//校验
	pwdInfo := strings.Split(pwd,"$")
	check := password.Verify(code,pwdInfo[2],pwdInfo[3],options)
	fmt.Println(check)
}
func genMd5(code string) string {
	Md5 := md5.New()
	_, _ = io.WriteString(Md5, code)
	return hex.EncodeToString(Md5.Sum(nil))
}
