package models

import (
	"mloginsvr/common/db"
	"mloginsvr/common/log"
)

//WechatUserInfo ..
type WechatUserInfo struct {
	Openid       string
	Accesstoken  string
	Refreshtoken string
	Nickname     string //普通用户昵称
	Sex          int    //普通用户性别，1 为男性，2 为女性
	Province     string //普通用户个人资料填写的省份
	City         string //普通用户个人资料填写的城市
	Country      string //国家，如中国为 CN
	Headimgurl   string //用户头像，最后一个数值代表正方形头像大小（有 0、46、64、96、132 数值可选，0 代表 640*640 正方形头像），用户没有头像时该项为空
	Privilege    string //用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	Unionid      string //	用户统一标识。针对一个微信开放平台帐号下的应用，同一用户的 unionid 是唯一的
}

//TableName ..
func (*WechatUserInfo) TableName() string {
	return "wechatUserInfo"
}

//Get ..
func (t *WechatUserInfo) Get(openid string) bool {
	has, err := db.MasterDB.Where("openid=?", openid).Get(t)
	if err != nil {
		log.WithFields(log.Fields{
			"has":    has,
			"err":    err,
			"openid": openid,
		}).Error("WechatUserInfo.go [Get] is err")
		return false
	}
	return has
}

//GetRefreshToken 获取refresh_token
func (t *WechatUserInfo) GetRefreshToken(openid string) string {
	has, err := db.MasterDB.Where("openid='?'", openid).Cols("refreshtoken").Get(t)
	if err != nil {
		log.WithFields(log.Fields{
			"has":    has,
			"err":    err,
			"openid": openid,
		}).Error("WechatUserInfo.go [Get] is err")
		return ""
	}
	if has {
		return t.Refreshtoken
	}
	return ""
}

//InsertOrUpdate 更新或插入token
func (t *WechatUserInfo) InsertOrUpdate() bool {
	has, err := db.MasterDB.Where("openid = ?", t.Openid).Exist(&Token{})
	if err != nil {
		log.WithFields(log.Fields{
			"has":         has,
			"err":         err,
			"accesstoken": t.Accesstoken,
		}).Error("WechatUserInfo.go [InsertOrUpdate]-1 is err")
		return false
	}
	if has {
		//存在则更新
		affected, err := db.MasterDB.Where("openid = ?", t.Openid).Update(t)
		if err != nil {
			log.WithFields(log.Fields{
				"affected":    affected,
				"err":         err,
				"accesstoken": t.Accesstoken,
			}).Error("WechatUserInfo.go [InsertOrUpdate]-2 is err")
			return false
		}
	} else {
		//不存在，则插入
		affected, err := db.MasterDB.Insert(t)
		if err != nil {
			log.WithFields(log.Fields{
				"affected":    affected,
				"err":         err,
				"accesstoken": t.Accesstoken,
			}).Error("WechatUserInfo.go [InsertOrUpdate]-3 is err")
			return false
		}
	}

	return true
}
