package global

import (
	"mloginsvr/common/log"
	"reflect"
	"regexp"
)

//StructToRedis ..将结构对象转化成Redis的hash参数
func StructToRedis(key string, src interface{}) []interface{} {
	args := []interface{}{key}
	// 获取结构体实例的反射类型对象
	typeOfSrc := reflect.TypeOf(src)
	valueofSrc := reflect.ValueOf(src)
	// 遍历结构体所有成员
	for i := 0; i < typeOfSrc.NumField(); i++ {
		log.Logger.Debugf("StructToRedis name: %v  value: '%v'", typeOfSrc.Field(i).Name, valueofSrc.Field(i))
		args = append(args, typeOfSrc.Field(i).Name, valueofSrc.Field(i))
	}
	return args
}

//VerifyMobileFormat 手机号合法性检测
func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
