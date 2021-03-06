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

func blessing() {
	////////////////////////////////////////////////////////////////////
	//                          _ooOoo_                               //
	//                         o8888888o                              //
	//                         88" . "88                              //
	//                         (| ^_^ |)                              //
	//                         O\  =  /O                              //
	//                      ____/`---'\____                           //
	//                    .'  \\|     |//  `.                         //
	//                   /  \\|||  :  |||//  \                        //
	//                  /  _||||| -:- |||||-  \                       //
	//                  |   | \\\  -  /// |   |                       //
	//                  | \_|  ''\---/''  |   |                       //
	//                  \  .-\__  `-`  ___/-. /                       //
	//                ___`. .'  /--.--\  `. . ___                     //
	//              ."" '<  `.___\_<|>_/___.'  >'"".                  //
	//            | | :  `- \`.;`\ _ /`;.`/ - ` : | |                 //
	//            \  \ `-.   \_ __\ /__ _/   .-` /  /                 //
	//      ========`-.____`-.___\_____/___.-`____.-'========         //
	//                           `=---='                              //
	//      ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^        //
	//      佛祖保佑       永不宕机     永无BUG      性能无敌            //
	////////////////////////////////////////////////////////////////////
}

func init() {
	//先祷告
	blessing()
	// 设置随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	sysinit()

	//路由====
	router.Start()
}

//
func sysinit() {
	//1.log的init
	//2.读取系统配置app.ini
	config.Init()
	//3.日志参数重置===============
	lv, err := config.AppCfg.GetValue("log", "level")
	if err != nil {
		log.Logger.Fatalln("read app.ini of log-lvevl is err:", err)
	}
	log.Init(lv)
	//4.db连接
	db.InitMysql()
	db.InitRedis()

	//5.
	global.Init()

	//4.路由
	router.Init()

}
