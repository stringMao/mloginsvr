package main

import (
	"mloginsvr/common/dbmanager"
	log "mloginsvr/common/logmanager"
	"mloginsvr/router"
)

func main() {
	log.Logger.Infoln("server start=======================")

	//数据库
	dbmanager.InitMysql()
	dbmanager.InitRedis()

	//路由==
	router.Start()
}
