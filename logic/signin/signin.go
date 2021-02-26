package signin

import (
	"encoding/json"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/logic/util"
	"mloginsvr/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

//askAccountLogin 客户端账号登入的请求信息
type askAccountLogin struct {
	Username string `json:"username" binding:"required,min=5"`
	Passwd   string `json:"passwd" binding:"required,min=1"`
	Channel  int    `json:"channel" binding:"required"` //位运算 (001)1安卓机  (010)2IOS机  (100)4PC机
}

//replyAccLogin 用户账号登入成功后的返回信息
type replyAccLogin struct {
	Userid  int64                   `json:"userid"`
	Token   string                  `json:"token"`
	Svrlist []global.HallSvrverInfo `json:"svrlist"`
}

//AccountLogin 账号密码登入
func AccountLogin(c *gin.Context) {
	var data askAccountLogin
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("AccountLogin read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	var acc models.Account
	if !acc.GetByUsername(data.Username) {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "账号不存在", nil))
		return
	}

	if strings.EqualFold(acc.Passwd, data.Passwd) {
		log.Logger.Debugln("登入密码验证成功")

		//检查账户状态
		if acc.Status != 0 {
			c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "账号状态异常", nil))
			return
		}

		//生成token
		token := global.GetLoginToken(acc.Userid)
		log.Logger.Debug("create token:" + token)
		//token放进redis
		util.InsertTokenToRedis(token, acc)

		var reply replyAccLogin
		reply.Userid = acc.Userid
		reply.Token = token
		//渠道分发及负载均衡策略-下发大厅服务器ip+端口
		reply.Svrlist = util.HallLoadBalanced(data.Channel)

		c.JSON(http.StatusOK, global.GetResultSucData(reply))
		return

	}

	c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "密码错误", nil))
	return
}

type askCheckLoginToken struct {
	Userid int64  `json:"userid" binding:"required"`
	Token  string `json:"token" binding:"required"`
}
type replyCheckLoginToken struct {
	Userid           int64  `json:"userid"`
	Nickname         string `json:"nickname"`
	Accounttype      int    `json:"accounttype"`
	Gender           int    `json:"gender"`
	CompleteRealname bool   `json:"completerealname"` //是否完成实名认证
	Age              int    `json:"age"`              //是否成年
}

//CheckLoginToken 验证登入token
func CheckLoginToken(c *gin.Context) {
	var data askCheckLoginToken
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("CheckLoginToken read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	//在Redis寻找==================================
	var str = fmt.Sprintf("token_%d", data.Userid)
	if re, err := redis.String(db.GetRedis().Do("GET", str)); err == nil {
		var u global.RedisUserInfo
		if json.Unmarshal([]byte(re), &u) == nil {
			//redis.ScanStruct(re, u)
			log.Logger.Debugf("%+v", u)
			if u.Token == data.Token {
				//验证通过
				var reply replyCheckLoginToken
				reply.Userid = data.Userid
				reply.Nickname = u.Nickname
				reply.Accounttype = u.Accounttype
				reply.Gender = u.Gender
				reply.CompleteRealname = u.CompleteRealname
				reply.Age = u.Age
				c.JSON(http.StatusOK, global.GetResultSucData(reply))
				return
			}
		}
	}

	//验证失败
	c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
	return

}
