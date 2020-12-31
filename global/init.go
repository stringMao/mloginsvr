package global

import (
	"mloginsvr/common/config"
	"mloginsvr/common/log"
)

//Init global初始化
func Init() {
	var err error
	GameHallToken, err = config.AppCfg.GetValue("server", "gamehalltoken")
	if err != nil {
		log.Logger.Fatalln("read app.ini of server-gamehalltoken is err:", err)
	}
	UserTokenKey, err = config.AppCfg.GetValue("server", "usertokenkey")
	if err != nil {
		log.Logger.Fatalln("read app.ini of server-usertokenkey is err:", err)
	}

	log.Logger.Info("global init success")
}

//GameHallToken 大厅服务器身份令牌
var GameHallToken = "default"

//UserTokenKey 生成用户登入token的key
var UserTokenKey = "default"
