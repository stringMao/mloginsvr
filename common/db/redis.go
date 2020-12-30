package db

import (
	"fmt"
	"mloginsvr/common/config"
	"mloginsvr/common/log"

	"github.com/gomodule/redigo/redis"
)

var redisPool *redis.Pool

//InitRedis ..
func InitRedis() {
	var err error
	host, err := config.AppCfg.GetValue("redis", "host")
	if err != nil {
		log.Logger.Fatalln("read app.ini of redis-host is err:", err)
	}
	port, err := config.AppCfg.Int("redis", "port")
	if err != nil {
		log.Logger.Fatalln("read app.ini of redis-port is err:", err)
	}
	// username, err := config.AppCfg.GetValue("redis", "username")
	// if err != nil {
	// 	log.Logger.Fatalln("read app.ini of redis-username is err:", err)
	// }
	// password, err := config.AppCfg.GetValue("redis", "password")
	// if err != nil {
	// 	log.Logger.Fatalln("read app.ini of redis-password is err:", err)
	// }
	poolSize, err := config.AppCfg.Int("redis", "poolsize")
	if err != nil {
		log.Logger.Fatalln("read app.ini of redis-poolsize is err:", err)
	}

	redisPool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
		if err != nil {
			log.Logger.Fatalln("redis connet is err:", err)
			return nil, err
		}
		return c, nil
	}, poolSize)
	//redisPool.MaxActive = 100

	//test redis connect
	GetRedis()
	log.Logger.Info("redis init success")
	// pool := &redis.Pool{
	// 	MaxActive:   100,                              //  最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）
	// 	MaxIdle:     10,                               // 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
	// 	IdleTimeout: time.Duration(100) * time.Second, // 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用
	// 	Wait:        true,                             // 当超过最大连接数 是报错还是等待， true 等待 false 报错
	// 	Dial: func() (redis.Conn, error) {
	// 		conn, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	// 		if err != nil {
	// 			log.Logger.Fatalln("redis connet is err:", err)
	// 			return nil, err
	// 		}
	// 		return conn, nil
	// 	},
	// }
}

//GetRedis ..
func GetRedis() redis.Conn {
	return redisPool.Get()
}
