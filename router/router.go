package router

import (
	"strings"

	"github.com/gin-gonic/gin"

	"mloginsvr/common/confmanager"
	log "mloginsvr/common/logmanager"
	"mloginsvr/logic/signin"
	"mloginsvr/middle"
)

var router = gin.Default()

func init() {
	//账号登入路由组============================================
	rgLogin := router.Group(strings.ToLower("/mloginsvr"))
	rgLogin.Use(middle.Respone()) //中间件设置
	registerLoginRouter(rgLogin)
}

//Start webapi启动
func Start() {
	port, err := confmanager.AppCfg.GetValue("router", "port")
	if err != nil {
		log.Logger.Fatalln("read app.ini of router-port is err:", err)
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Logger.Fatalln("router start is err:", err)
	}

	//router.RunTLS(conf.JSONConf.Port, crtPath, keyPath)
}

//registerLoginRouter 账号登入及token验证等相关路由
func registerLoginRouter(group *gin.RouterGroup) {
	group.POST(strings.ToLower("/signin"), signin.AccountLogin)
	group.POST(strings.ToLower("/checktoken"), signin.CheckLoginToken)
}
