package wechat

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	appid  = "sssssssssss"
	secret = "sssssssssssss"
)

type askGetAccessTokenByCode struct {
	Code string `json:"code"`
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

	url := "https://api.weixin.qq.com/sns/oauth2/access_token"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Logger.Errorln(err)
		return
	}

	q := req.URL.Query()
	q.Add("appid", appid)
	q.Add("secret", secret)
	q.Add("code", data.Code)
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()
	//处理返回结果
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Logger.Errorln(err)
		return
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Logger.Errorln(err)
		return
	}
	retMap := make(map[string]string)
	if json.Unmarshal(body, &retMap) != nil {
		log.Logger.Errorln("wechat.go GetAccessTokenByCode json.Unmarshal err")
		return
	}

	/*错误返回样例：
	{
	  "errcode": 40029,
	  "errmsg": "invalid code"
	}*/
	if errcode, ok := retMap["errcode"]; ok { //有错误码
		log.Logger.Errorf("wechat.go GetAccessTokenByCode errcode:%s,errmsg:%s ", errcode, retMap["errmsg"])
		c.JSON(http.StatusOK, global.GetResultData(global.CodeLoginCodeErr, "code错误", nil))
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
	//检查账户是否存在
	var acc models.Account
	username := fmt.Sprintf("wechat_%s", retMap["openid"])
	if !acc.GetByUsername(username) {
		//账号不存在则插入
		acc.Username = username
		acc.Passwd = fmt.Sprintf("%x", md5.Sum([]byte(time.Now().Format("060102150405"))))
		acc.Thirdid = retMap["openid"]
		acc.Nickname = fmt.Sprintf("wechat_defalut")
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

	//获得token======================
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

}
