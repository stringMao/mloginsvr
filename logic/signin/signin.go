package signin

import (
	"crypto/md5"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type askAccountLogin struct {
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
}
type userinfo struct {
	Token     string
	Nickname  string
	Channelid int
	Status    int
}

//AccountLogin 账号密码登入
func AccountLogin(c *gin.Context) {
	data := &askAccountLogin{}
	if err := c.BindJSON(&data); err != nil {
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

	if acc.Passwd == data.Passwd {
		log.Logger.Debugln("登入验证成功")
		token := new(models.Token)
		token.Userid = acc.Userid
		//生成token
		buf := fmt.Sprintf("userid=%d&time=%d&key=%s", acc.Userid, time.Now().Unix(), global.UserTokenKey)
		token.Token = fmt.Sprintf("%x", md5.Sum([]byte(buf)))
		log.Logger.Debug(token.Token)

		re := token.InsertOrUpdate()
		if re {
			log.Logger.Debugln("token 写入db完成")

			//token放进redis====================================
			u := userinfo{}
			u.Token = token.Token
			if acc.Nickname2 != "" {
				u.Nickname = acc.Nickname2
			} else {
				u.Nickname = acc.Nickname
			}
			u.Channelid = acc.Channelid
			u.Status = acc.Status

			var str = fmt.Sprintf("token_%d", token.Userid)
			// db.GetRedis().Do("SETEX", str, 60, token.Token)
			args := global.StructToRedis(str, u)
			_, err := db.GetRedis().Do("HMSET", args...)
			if err != nil {
				log.Logger.Errorln("redis write token err:", err)
			} else {
				db.GetRedis().Do("EXPIRE", str, 60)
			}
			c.JSON(http.StatusOK, global.GetResultSucData(gin.H{"userid": token.Userid, "token": token.Token}))
			return
		} else {
			c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginFail, "db error", nil))
			return
		}
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
	Userid    int64  `json:"userid"`
	Nickname  string `json:"nickname"`
	Channelid int    `json:"channelid"`
	Status    int    `json:"status"`
}

//CheckLoginToken 验证登入token
func CheckLoginToken(c *gin.Context) {
	data := &askCheckLoginToken{}
	if err := c.BindJSON(&data); err != nil {
		log.Logger.Errorln("CheckLoginToken read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	//身份认证
	if data.Passwd != global.HallToken {
		log.Logger.Errorln("ModifyNickname server pw is wrong")
		return
	}

	//先在Redis寻找==================================
	var str = fmt.Sprintf("token_%d", data.Userid)
	re, err := redis.Values(db.GetRedis().Do("HGETALL", str))
	if err == nil && len(re) != 0 {
		// for _, v := range re {
		// 	log.Logger.Debugf("%s", v.([]byte))
		// }
		u := new(userinfo)
		redis.ScanStruct(re, u)
		log.Logger.Debugf("%+v", u)
		if u.Token == data.Token {
			//验证通过
			var reply replyCheckLoginToken
			reply.Userid = data.Userid
			reply.Nickname = u.Nickname
			reply.Channelid = u.Channelid
			reply.Status = u.Status
			c.JSON(http.StatusOK, global.GetResultSucData(reply))
			return
		} else {
			//验证失败
			c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
			return
		}
	} else {
		log.Logger.Debugln("redis get token is err:", err)
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
		reply.Channelid = acc.Channelid
		reply.Status = acc.Status
		c.JSON(http.StatusOK, global.GetResultSucData(reply))
		return
	}
	//验证失败
	c.JSON(http.StatusOK, global.GetResultData(global.CodeCheckTokenFail, "token 错误", nil))
	return

}
