package global

import (
	"mloginsvr/common/log"
	"reflect"
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
