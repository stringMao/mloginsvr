package main

import (
	"math/rand"
	"mloginsvr/common/config"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/router"
	"time"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	sysinit()

	//路由==
	router.Start()
}

func sysinit() {
	//1.先开启日志
	log.Init()
	//2.读取系统配置app.ini
	config.Init()
	//3.db
	db.InitMysql()
	db.InitRedis()
	//4.路由
	router.Init()

}
