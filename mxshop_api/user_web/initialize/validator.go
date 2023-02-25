package initialize

import (
	"fmt"
	"mxshop_api/user_web/global"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_teanslations "github.com/go-playground/validator/v10/translations/en"
	zh_teanslations "github.com/go-playground/validator/v10/translations/zh"
	"go.uber.org/zap"
)

func InitTrans(locale string) (err error) {
	//修改gin中的validator引擎属性
	if v,ok := binding.Validator.Engine().(*validator.Validate);ok {
		//注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField)string{
			name := strings.SplitN(fld.Tag.Get("json"),",",2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhT := zh.New()
		enT := en.New()
		uni := ut.New(enT,zhT,enT)	//第一个参数为备用语言环境
		global.Trans,ok = uni.GetTranslator(locale)
		if !ok {
			zap.S().Errorf("uni.GGetTranslator(%s)",locale)
			return fmt.Errorf("uni.GGetTranslator(%s)",locale)
		}
		switch locale {
		case "en":
			en_teanslations.RegisterDefaultTranslations(v,global.Trans)
		case "zh":
			zh_teanslations.RegisterDefaultTranslations(v,global.Trans)
		default:
			en_teanslations.RegisterDefaultTranslations(v,global.Trans)
		}
		return nil
	}
	return nil
}