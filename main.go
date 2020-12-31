package main

import (
	"math/rand"
	"mloginsvr/common/config"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/router"
	"time"
)

func init() {
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	sysinit()

	//路由====
	router.Start()
}

func sysinit() {
	//读取系统配置app.ini
	config.Init()
	//开启日志
	log.Init()
	//
	global.Init()
	//3.db
	db.InitMysql()
	db.InitRedis()
	//4.路由
	router.Init()

}
