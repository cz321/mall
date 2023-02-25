package validator

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

func ValidateMobile(fl validator.FieldLevel)  bool{
	mobile := fl.Field().String()
	//使用正则表达式判断是否合法
	//ok,_ := regexp.Match(`^1([38][0-9]14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`,[]byte(mobile))
	ok,_ := regexp.Match(`^1([358][0-9]|4[579]|66|7[0135678]|9[89])[0-9]{8}$`,[]byte(mobile))
	if !ok {
		return false
	}
	return true
}
