package validator

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"regexp"
)

// 自定义手机验证器
func ValidateMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	zap.S().Infof("我是mobile参数：%v", mobile)
	// 用正则表达式判断是否合法
	ok, _ := regexp.MatchString(`^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`, mobile)
	if !ok {
		return false
	}
	return true

}
