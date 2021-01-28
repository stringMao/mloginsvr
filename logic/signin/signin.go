package signin

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

//redis 储存的用户信息
type userinfo struct {
	Token       string
	Nickname    string
	Accounttype int
}
type hallSvrverInfo struct {
	Address string
	Svrname string
}

//askAccountLogin 客户端账号登入的请求信息
type askAccountLogin struct {
	Username string `json:"username"  binding:"required"`
	Passwd   string `json:"passwd"  binding:"required"`
	Channel  int    `json:"channel"  binding:"required"` //位运算 (001)1安卓机  (010)2IOS机  (100)4PC机
}

//replyAccLogin 用户账号登入成功后的返回信息
type replyAccLogin struct {
	Userid  int64            `json:"userid"`
	Token   string           `json:"token"`
	Svrlist []hallSvrverInfo `json:"svrlist"`
}

//AccountLogin 账号密码登入
func AccountLogin(c *gin.Context) {
	data := &askAccountLogin{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("AccountLogin read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	acc := new(models.Account)
	acc.GetByUsername(data.Username)
	log.Logger.Debugf("%+v", acc)
	if acc.Username == "" {
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
		token := new(models.Token)
		token.Userid = acc.Userid
		//生成token
		buf := fmt.Sprintf("userid=%d&time=%d&key=%s", acc.Userid, time.Now().Unix(), global.LoginTokenKey)
		token.Token = fmt.Sprintf("%x", md5.Sum([]byte(buf)))
		log.Logger.Debug("create token:" + token.Token)

		if token.InsertOrUpdate() {
			log.Logger.Debugln("token 写入db完成")

			//token放进redis====================================
			u := userinfo{}
			u.Token = token.Token
			u.Nickname = acc.Nickname
			u.Accounttype = acc.Accounttype
			//u.Status = acc.Status

			var str = fmt.Sprintf("token_%d", token.Userid)
			// db.GetRedis().Do("SETEX", str, 60, token.Token)
			args := global.StructToRedis(str, u)
			_, err := db.GetRedis().Do("HMSET", args...)
			if err != nil {
				log.Logger.Errorln("redis write token err:", err)
			} else {
				db.GetRedis().Do("EXPIRE", str, 60)
			}

			var reply replyAccLogin
			reply.Userid = token.Userid
			reply.Token = token.Token
			//渠道分发及负载均衡策略-下发大厅服务器ip+端口
			for _, v := range global.HallList {
				if (data.Channel & v.Channel) > 0 {
					reply.Svrlist = append(reply.Svrlist, hallSvrverInfo{v.Address, v.Servername})
				}
			}

			c.JSON(http.StatusOK, global.GetResultSucData(reply))
			return
		}
		c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "db error", nil))
		return
	}

	c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "密码错误", nil))
	return
}

type askCheckLoginToken struct {
	Userid int64  `json:"userid"`
	Token  string `json:"token"`
	Passwd string `json:"passwd"` //服务器身份认证
}
type replyCheckLoginToken struct {
	Userid      int64  `json:"userid"`
	Nickname    string `json:"nickname"`
	Accounttype int    `json:"accounttype"`
}

//CheckLoginToken 验证登入token
func CheckLoginToken(c *gin.Context) {
	data := &askCheckLoginToken{}
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("CheckLoginToken read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	//先在Redis寻找==================================
	var str = fmt.Sprintf("token_%d", data.Userid)
	if re, err := redis.String(db.GetRedis().Do("HGETALL", str)); err == nil {
		u := new(userinfo)
		if err2 := json.Unmarshal([]byte(re), u); err2 == nil {
			//redis.ScanStruct(re, u)
			log.Logger.Debugf("%+v", u)
			if u.Token == data.Token {
				//验证通过
				var reply replyCheckLoginToken
				reply.Userid = data.Userid
				reply.Nickname = u.Nickname
				reply.Accounttype = u.Accounttype
				c.JSON(http.StatusOK, global.GetResultSucData(reply))
				return
			}
			//验证失败
			c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
			return
		}
	}

	//mysql里查验===========================
	token := new(models.Token)
	token.GetByUserid(data.Userid)
	if token.Token == data.Token {
		acc := new(models.Account)
		acc.GetByUserid(data.Userid)
		//验证通过
		var reply replyCheckLoginToken
		reply.Userid = data.Userid
		reply.Nickname = acc.Nickname
		reply.Accounttype = acc.Accounttype
		c.JSON(http.StatusOK, global.GetResultSucData(reply))
		return
	}
	//验证失败
	c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
	return

}
