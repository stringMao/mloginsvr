package signin

import (
	"encoding/json"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type hallSvrverInfo struct {
	Address string
	Svrname string
}

//askAccountLogin 客户端账号登入的请求信息
type askAccountLogin struct {
	Username string `json:"username" binding:"required,min=5"`
	Passwd   string `json:"passwd" binding:"required,min=1"`
	Channel  int    `json:"channel" binding:"required"` //位运算 (001)1安卓机  (010)2IOS机  (100)4PC机
}

//replyAccLogin 用户账号登入成功后的返回信息
type replyAccLogin struct {
	Userid  int64            `json:"userid"`
	Token   string           `json:"token"`
	Svrlist []hallSvrverInfo `json:"svrlist"`
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

	acc := new(models.Account)
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
		token := new(models.Token)
		token.Userid = acc.Userid
		//生成token
		token.Token = global.GetLoginToken(acc.Userid)
		log.Logger.Debug("create token:" + token.Token)

		if !token.InsertOrUpdate() {
			c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "db error", nil))
			return
		}
		log.Logger.Debugln("token 写入db完成")

		//token放进redis====================================
		u := global.RedisUserInfo{Gender: 0, CompleteRealname: false, Age: 0}
		u.Token = token.Token
		u.Nickname = acc.Nickname
		u.Accounttype = acc.Accounttype
		//u.Indentity = acc.Identity
		//获得实名信息
		var userinfo models.UserRealInfo
		if userinfo.GetByUserid(acc.Userid) {
			u.CompleteRealname = true
			u.Gender = userinfo.Gender
			u.Age = global.GetCitizenAge([]byte(userinfo.Identity), false)
		}

		var key = fmt.Sprintf("token_%d", token.Userid)
		if val, err := json.Marshal(u); err == nil {
			db.GetRedis().Do("SET", key, string(val), "EX", 60)
		}

		// args := global.StructToRedis(str, u)
		// _, err := db.GetRedis().Do("HMSET", args...)
		// if err != nil {
		// 	log.Logger.Errorln("redis write token err:", err)
		// } else {
		// 	db.GetRedis().Do("EXPIRE", str, 60*30)
		// }
		//===========================================================

		var reply replyAccLogin
		reply.Userid = token.Userid
		reply.Token = token.Token
		//渠道分发及负载均衡策略-下发大厅服务器ip+端口
		reply.Svrlist = hallLoadBalanced(data.Channel)

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

	//先在Redis寻找==================================
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
			//验证失败
			c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
			return
		}
	}

	//mysql里查验===========================
	token := new(models.Token)
	token.GetByUserid(data.Userid)
	if token.Token == data.Token {
		//acc := new(models.Account)
		var acc models.Account
		acc.GetByUserid(data.Userid)

		//验证通过
		reply := replyCheckLoginToken{Gender: 0, CompleteRealname: false, Age: 0}
		reply.Userid = data.Userid
		reply.Nickname = acc.Nickname
		reply.Accounttype = acc.Accounttype

		var info models.UserRealInfo
		if info.GetByUserid(data.Userid) {
			reply.Gender = info.Gender
			reply.CompleteRealname = true
			reply.Age = global.GetCitizenAge([]byte(info.Identity), false)
		}

		c.JSON(http.StatusOK, global.GetResultSucData(reply))
		return
	}
	//验证失败
	c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
	return

}

//hallLoadBalanced 大厅渠道分发及负载策略
func hallLoadBalanced(channel int) []hallSvrverInfo {
	var result []hallSvrverInfo
	//渠道分发
	for _, v := range global.HallList {
		if (channel&v.Channel) > 0 && v.Status == 0 {
			result = append(result, hallSvrverInfo{v.Address, v.Servername})
		}
	}
	// 负载均衡
	// if len(result) > 1 {

	// }

	return result
}
