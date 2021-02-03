package signup

/* 注册相关接口 注册，获取验证码，修改密码，修改昵称 */

import (
	"fmt"
	"math/rand"
	"mloginsvr/common/config"
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"mloginsvr/global"
	"mloginsvr/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type askRegisterAcc struct {
	Phone  string `json:"phone" binding:"required"` //手机号作为账号名
	Passwd string `json:"passwd" binding:"required,min=1"`
	//Nickname string `json:"nickname"`
	Smscode string `json:"smscode" binding:"required"`
}

//RegisterAccount 注册账号=================================
func RegisterAccount(c *gin.Context) {
	var data askRegisterAcc
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("RegisterAccount read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	//验证手机号合法
	if global.VerifyMobileFormat(data.Phone) == false {
		c.JSON(http.StatusOK, global.GetResultData(global.CodePhoneErr, "账号错误", nil))
		return
	}

	//验证码检查
	str := fmt.Sprintf("smscode_1_%s", data.Phone)
	re, err := redis.String(db.GetRedis().Do("GET", str))
	if err != nil {
		log.Logger.Debug(err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSmsVerifyFail, "验证码不存在", nil))
		return
	}
	if re != data.Smscode {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSmsVerifyFail, "验证码错误", nil))
		return
	}
	//密码格式检查(客户端上传加密后的密码)
	if data.Passwd == "" || len(data.Passwd) != 32 {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeRegiterOfPasswdErr, "密码不合规", nil))
		return
	}

	//手机号做账号，账号查重
	acc := new(models.Account)
	if acc.ExistByUsername(data.Phone) {
		//账号已存在
		c.JSON(http.StatusOK, global.GetResultData(global.CodeAccOfPhoneRepeat, "该手机号已注册账号", nil))
		return
	}

	//生成账号
	acc.Username = data.Phone
	acc.Passwd = data.Passwd
	acc.Nickname = "游客"
	acc.Createtime = time.Now()
	acc.Accounttype = 0 //账号用户
	acc.Status = 0
	acc.Phone = data.Phone
	if acc.Insert() {
		//生成默认昵称
		acc.Nickname = fmt.Sprintf("游客_%d", acc.Userid)
		acc.UpdateNickname()
	}
	c.JSON(http.StatusOK, global.GetResultSucData(nil)) //通知注册成功
}

type askSMS struct {
	Phone    string `json:"phone" binding:"required"`
	CodeType int    `json:"codetype" binding:"required"` //申请验证码的业务类型 1注册验证码 2重置密码的验证码
}
type replyApplySmsCode struct {
}

//ApplySMSVerificationCode 申请短息验证码===============
func ApplySMSVerificationCode(c *gin.Context) {
	var data askSMS
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("ApplySMSVerificationCode read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	//验证手机号合法
	if global.VerifyMobileFormat(data.Phone) == false {
		c.JSON(http.StatusOK, global.GetResultData(global.CodePhoneErr, "请输入正确的手机号", nil))
		return
	}
	str := fmt.Sprintf("smscode_%d_%s", data.CodeType, data.Phone)
	strlock := fmt.Sprintf("smslock_%s", data.Phone)
	//检查是否已经发送过了，且还未过期
	keyexit, _ := redis.Bool(db.GetRedis().Do("EXISTS", strlock))
	if keyexit == true {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSMSOften, "短信请求太频繁，请稍后再试", nil))
		return
	}

	//获取随机验证码
	rndCode := fmt.Sprintf("%06v", rand.Int31n(1000000))
	log.Logger.Debug("验证码：", rndCode)

	//验证码保存到redis
	_, err := db.GetRedis().Do("SETEX", str, 5*60, rndCode)
	if err != nil {
		log.Logger.Errorln("redis write [smscode] err:", err)
		return
	}
	db.GetRedis().Do("SETEX", strlock, 60, "短信已发送")

	//发送短信验证码=====
	//alibaba.SendSms(data.Phone, rndCode)

	c.JSON(http.StatusOK, global.GetResultSucData(nil)) //通知短信已发送
}

type askModifyNickname struct {
	Userid   int64  `json:"userid" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

//ModifyNickname 修改昵称=============================
func ModifyNickname(c *gin.Context) {
	var data askModifyNickname
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("ModifyNickname read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)

	if data.Nickname == "" {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeNicknameErr, "昵称不能为空", nil))
		return
	}
	if len(data.Nickname) > 12 {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeNicknameErr, "昵称太长了", nil))
		return
	}

	//昵称查重
	acc := new(models.Account)
	if acc.ExistByUsername(data.Nickname) {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeNicknameErr, "昵称已存在", nil))
		return
	}

	//昵称敏感词检查
	if config.HaveSenWords(data.Nickname) {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeNicknameErr, "昵称不可用", nil))
		return
	}

	//完成修改
	acc.Userid = data.Userid
	acc.Nickname = data.Nickname
	if acc.UpdateNickname() {
		c.JSON(http.StatusOK, global.GetResultSucData(nil))
		return
	}
}

type askResetPasswd struct {
	Phone     string `json:"phone" binding:"required"`
	PasswdNew string `json:"passwd" binding:"required"`  //新密码
	Smscode   string `json:"smscode" binding:"required"` //短信验证码
}

//LostPasswd 忘记密码(通过手机号验证码重置密码)====================
func LostPasswd(c *gin.Context) {
	var data askResetPasswd
	if err := c.ShouldBindJSON(&data); err != nil {
		log.Logger.Errorln("LostPasswd read data err:", err)
		c.JSON(http.StatusOK, global.GetResultData(global.CodeProtoErr, "协议解析错误", nil))
		return
	}
	log.Logger.Debugf("%+v", data)
	//新密码格式检查
	if data.PasswdNew == "" || len(data.PasswdNew) != 32 {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeRegiterOfPasswdErr, "密码不合规", nil))
		return
	}

	//验证码检查
	str := fmt.Sprintf("smscode_2_%s", data.Phone)
	re, err := redis.String(db.GetRedis().Do("GET", str))
	if err != nil {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSmsVerifyFail, "验证码不存在", nil))
		return
	}
	if re != data.Smscode {
		c.JSON(http.StatusOK, global.GetResultData(global.CodeSmsVerifyFail, "验证码错误", nil))
		return
	}
	//修改数据库
	acc := new(models.Account)
	acc.Username = data.Phone
	acc.Passwd = data.PasswdNew

	if !acc.ResetPasswdByUsername() {
		//修改数据库密码失败
		c.JSON(http.StatusOK, global.GetResultData(global.CodeDBExecErr, "db异常", nil))
		return
	}
	c.JSON(http.StatusOK, global.GetResultSucData(nil))
	return
}
