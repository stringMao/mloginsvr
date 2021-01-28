package router

import (
	"strings"

	"github.com/gin-gonic/gin"

	"mloginsvr/common/config"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/logic/signin"
	"mloginsvr/logic/signup"
	"mloginsvr/middle"
)

var router = gin.Default()

//Init ..
func Init() {
	//路由分组
	//大厅服务器路由
	svrRG := router.Group(strings.ToLower("/mloginsvr/hall"))
	svrRG.Use(middle.Respone(), middle.AuthCheckSign(global.HallSignKey), middle.Authentication()) //中间件设置
	svrRGRouter(svrRG)

	//客户端路由
	clientRG := router.Group(strings.ToLower("/mloginsvr/client"))
	clientRG.Use(middle.Respone(), middle.AuthCheckSign(global.ClientSignKey)) //中间件设置
	clientRGRouter(clientRG)

	//管理员路由

}

//Start webapi启动
func Start() {
	port, err := config.AppCfg.GetValue("router", "port")
	if err != nil {
		log.Logger.Fatalln("read app.ini of router-port is err:", err)
	}

	err = router.Run(":" + port)
	if err != nil {
		log.Logger.Fatalln("router start is err:", err)
	}

	//router.RunTLS(conf.JSONConf.Port, crtPath, keyPath)
}

//svrRGRouter 大厅服务器访问接口
func svrRGRouter(group *gin.RouterGroup) {
	//登入token验证
	group.POST(strings.ToLower("/checktoken"), signin.CheckLoginToken)
	//修改昵称
	group.POST(strings.ToLower("/modifynickname"), signup.ModifyNickname)
}

//registerLoginRouter 账号登入及token验证等相关路由
func clientRGRouter(group *gin.RouterGroup) {
	//客户端账号登入
	group.POST(strings.ToLower("/signin"), signin.AccountLogin)
	//账号注册
	group.POST(strings.ToLower("/signup"), signup.RegisterAccount)
	//请求短信验证码
	group.POST(strings.ToLower("/applysms"), signup.ApplySMSVerificationCode)
	//忘记密码重置密码
	group.POST(strings.ToLower("/resetpasswd"), signup.LostPasswd)

}
