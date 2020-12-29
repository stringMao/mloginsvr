package main

import (
	"msigninser/common/dbmanager"
	log "msigninser/common/logmanager"
	"msigninser/router"
)

var n = 0

func main() {
	log.Logger.Infoln("server start=======================")

	//数据库
	dbmanager.InitMysql()
	dbmanager.InitRedis()

	//路由==
	router.Start()
}
