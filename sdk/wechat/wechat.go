package wechat

import (
	"crypto/md5"
	"fmt"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/logic/util"
	"mloginsvr/models"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

const (
	appid  = "xxx"
	secret = "xxx"
)

//replyWechatLogin 微信登入成功后的返回信息
type replyWechatLogin struct {
	Openid      string                  `json:"openid"`
	AccessToken string                  `json:"accesstoken"`
	Userid      int64                   `json:"userid"`
	Token       string                  `json:"token"`
	Svrlist     []global.HallSvrverInfo `json:"svrlist"`
}

type askGetAccessTokenByCode struct {
	Code    string `json:"code" binding:"required"`
	Channel int    `json:"channel" binding:"required"` //位运算 (001)1安卓机  (010)2IOS机  (100)4PC机
}

//GetAccessTokenByCode 通过 code 获取 access_token
func GetAccessTokenByCode(c *gin.Context) {
	var data askGetAccessTokenByCode
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("GetAccessTokenByCode read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	params := make(global.Params)
	params.SetString("appid", appid)
	params.SetString("secret", secret)
	params.SetString("code", data.Code)
	params.SetString("grant_type", "authorization_code")

	var ret refreshAccessTokenStruct
	err := global.NewClient().GetWithoutCert("https://api.weixin.qq.com/sns/oauth2/access_token", params, &ret)
	if err != nil {
		log.Logger.Errorln("wechat.go GetAccessTokenByCode GetWithoutCert err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeAskFail, "wechat接口访问失败", nil))
		return
	}
	/*错误返回样例：
	{
	  "errcode": 40029,
	  "errmsg": "invalid code"
	}*/
	if ret.Errcode != 0 { //有错误码
		log.Logger.Errorf("wechat.go GetAccessTokenByCode errcode:%d,errmsg:%s ", ret.Errcode, ret.Errmsg)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeWechatLoginCodeErr, "code错误", nil))
		return
	}
	/*正确的返回：
	{
	  "access_token": "ACCESS_TOKEN",
	  "expires_in": 7200,
	  "refresh_token": "REFRESH_TOKEN",
	  "openid": "OPENID",   //授权用户唯一标识
	  "scope": "SCOPE"
	}*/
	//获取数据

	var we models.WechatUserInfo
	we.Openid = ret.Openid
	we.Accesstoken = ret.AccessToken
	we.Refreshtoken = ret.RefreshToken

	//检查账户是否存在
	var acc models.Account
	username := fmt.Sprintf("wechat_%s", ret.Openid)
	if !acc.GetByUsername(username) { //账号不存在
		//==第一登入获取微信的账号信息
		getUserWechatInfo(&we)
		//插入新账号
		acc.Username = username
		acc.Passwd = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().Format("060102150405"))))
		acc.Thirdid = ret.Openid
		if we.Nickname != "" {
			acc.Nickname = we.Nickname
		} else {
			acc.Nickname = fmt.Sprintf("wechat_defalut")
		}
		acc.Createtime = time.Now()
		acc.Accounttype = global.AccTypeWechat
		acc.Status = 0
		//acc.Phone = data.Phone
		if !acc.Insert() {
			log.Logger.Errorln("wechat.go GetAccessTokenByCode acc.Insert err")
			c.JSON(http.StatusOK, global.GetResultData(global.CodeDBExecErr, "db错误", nil))
			return
		}
		// 注册完成再取一次
		if !acc.GetByUsername(username) {
			log.Logger.Errorln("wechat.go GetAccessTokenByCode acc.GetByUsername err")
			c.JSON(http.StatusOK, global.GetResultData(global.CodeDBExecErr, "db错误", nil))
			return
		}
	}

	token := global.GetLoginToken(acc.Userid) //生成token
	log.Logger.Debug("create token:" + token)
	util.InsertTokenToRedis(token, acc) //token放进redis

	//=========================================================
	//将access_token有效期信息保存在redis
	db.GetRedis().Do("SET", fmt.Sprintf("Wechat_actoken_%s", ret.Openid), ret.AccessToken, "EX", ret.ExpriesIn-5)
	//将refresh_token 刷新信息保存在sql
	we.InsertOrUpdate()

	//返回登入token给客户端
	var reply replyWechatLogin
	reply.Userid = acc.Userid
	reply.Token = token
	reply.AccessToken = ret.AccessToken
	reply.Openid = ret.Openid
	reply.Svrlist = util.HallLoadBalanced(data.Channel)
	c.JSON(http.StatusOK, global.GetResultSucData(reply))
}

type askLoginByAccessToken struct {
	Openid      string `json:"openid" binding:"required"`
	AccessToken string `json:"accesstoken" binding:"required"`
	Channel     int    `json:"channel" binding:"required"` //位运算 (001)1安卓机  (010)2IOS机  (100)4PC机
}

//LoginByAccessToken 通过accesstoken登入
func LoginByAccessToken(c *gin.Context) {
	var data askLoginByAccessToken
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("LoginByAccessToken read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}

	bRefresh := true
	//现在redis检查accesstoken是否过期
	str := fmt.Sprintf("Wechat_actoken_%s", data.Openid)
	if re, err := redis.String(db.GetRedis().Do("GET", str)); err == nil {
		if data.AccessToken != re {
			c.JSON(http.StatusOK, global.GetResultData(global.CodeWechatAccessTokenErr, "请重新code登入", nil))
			return
		}
		//accesstoken 没有过期
		if checkAccessToken(re, data.Openid) {
			bRefresh = false
		}
	}

	var reply replyWechatLogin
	reply.Openid = data.Openid
	reply.AccessToken = data.AccessToken

	//如果过期-直接尝试刷新accesstoken
	if bRefresh {
		bsucc := false
		var we models.WechatUserInfo
		if we.Get(data.Openid) {
			t, o := refreshAccessToken(we.Refreshtoken)
			if t != "" && o != "" && o == data.Openid {
				reply.AccessToken = t
				bsucc = true
			}
		}
		if !bsucc { //刷新失败
			c.JSON(http.StatusOK, global.GetResultData(global.CodeWechatAccessTokenErr, "请重新code登入", nil))
			return
		}
	}

	//用access_token验证登入
	if !checkAccessToken(reply.AccessToken, reply.Openid) {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeWechatAccessTokenErr, "请重新code登入", nil))
		return
	}

	//获取账号信息
	username := fmt.Sprintf("wechat_%s", reply.Openid)
	var acc models.Account
	if !acc.GetByUsername(username) {
		log.Logger.Errorf("微信登入 账号不存在:%s", username)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeWechatAccessTokenErr, "请重新code登入", nil))
		return
	}
	reply.Userid = acc.Userid

	//生成token
	token := global.GetLoginToken(acc.Userid)
	log.Logger.Debug("create token:" + token)
	util.InsertTokenToRedis(token, acc) //token放进redis
	reply.Token = token
	reply.Svrlist = util.HallLoadBalanced(data.Channel)

	c.JSON(http.StatusOK, global.GetResultSucData(reply))
	return
}

type checkAccessTokenStruct struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

//checkAccessToken 检验授权凭证（access_token）是否有效
func checkAccessToken(actoken, openid string) bool {
	params := make(global.Params)
	params.SetString("access_token", actoken)
	params.SetString("openid", openid)

	var ret checkAccessTokenStruct
	err := global.NewClient().GetWithoutCert("https://api.weixin.qq.com/sns/auth", params, &ret)
	if err != nil {
		log.Logger.Errorln("wechat.go checkAccessToken GetWithoutCert err:", err)
		return false
	}
	/*
			错误返回样例：
			{
				"errcode": 40030,
				"errmsg": "invalid openid"
			}
			正确的 Json 返回结果
			{
		  		"errcode": 0,
		  		"errmsg": "ok"
			}
	*/
	if ret.Errcode == 0 && ret.Errmsg == "ok" {
		return true
	}

	log.Logger.Debugf("wechat.go checkAccessToken result code:%d", ret.Errcode)
	return false
}

type getUserWechatInfoStruct struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Openid     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	Headimgurl string `json:"headimgurl"`
	//Privilege  string `json:"privilege"`
	Unionid string `json:"unionid"`
}

//getUserWechatInfo 获取用户微信信息
func getUserWechatInfo(res *models.WechatUserInfo) {
	params := make(global.Params)
	params.SetString("access_token", res.Accesstoken)
	params.SetString("openid", res.Openid)

	var ret getUserWechatInfoStruct
	err := global.NewClient().GetWithoutCert("https://api.weixin.qq.com/sns/userinfo", params, &ret)
	if err != nil {
		log.Logger.Errorln("wechat.go getUserWechatInfo GetWithoutCert err:", err)
		return
	}
	if ret.Errcode != 0 {
		log.Logger.Errorf("wechat.go getUserWechatInfo errcode:%d,errmsg:%s ", ret.Errcode, ret.Errmsg)
		return
	}
	res.Nickname = ret.Nickname
	res.Sex = ret.Sex
	res.Province = ret.Province
	res.City = ret.City
	res.Country = ret.Country
	res.Headimgurl = ret.Headimgurl
	res.Privilege = ""
	res.Unionid = ret.Unionid
}

type refreshAccessTokenStruct struct {
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpriesIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

//refreshAccessToken 刷新access_token
func refreshAccessToken(refreshtoken string) (actoken, wxopenid string) {
	params := make(global.Params)
	params.SetString("appid", appid)
	params.SetString("grant_type", "refresh_token")
	params.SetString("refresh_token", refreshtoken)
	var ret refreshAccessTokenStruct
	err := global.NewClient().GetWithoutCert("https://api.weixin.qq.com/sns/oauth2/refresh_token", params, &ret)
	if err != nil {
		log.Logger.Errorln("wechat.go refreshAccessToken GetWithoutCert err:", err)
		return "", ""
	}
	/*错误返回样例：
	{
		"errcode": 40030,
		"errmsg": "invalid refresh_token"
	}*/
	if ret.Errcode != 0 { //有错误码
		log.Logger.Errorf("wechat.go refreshAccessToken errcode:%d,errmsg:%s ", ret.Errcode, ret.Errmsg)
		return "", ""
	}
	/*正确的返回：
	{
		"access_token": "ACCESS_TOKEN",
		"expires_in": 7200,
		"refresh_token": "REFRESH_TOKEN",
		"openid": "OPENID",
		"scope": "SCOPE"
	}*/
	//将access_token有效期信息保存在redis
	db.GetRedis().Do("SET", fmt.Sprintf("Wechat_actoken_%s", ret.Openid), ret.AccessToken, "EX", ret.ExpriesIn-5)

	//将refresh_token 刷新信息保存在sql
	var we models.WechatUserInfo
	we.Openid = ret.Openid
	we.Accesstoken = ret.AccessToken
	we.Refreshtoken = ret.RefreshToken
	we.InsertOrUpdate()

	return ret.AccessToken, ret.Openid
}
