package models

import (
	"mloginsvr/common/db"
	"mloginsvr/common/log"
	"time"
)

//Account 用户信息表
type Account struct {
	Userid     int64  `xorm:"pk autoincr notnull"`
	Username   string `xorm:"unique notnull"`
	Passwd     string
	Thirdid    string
	Nickname   string
	Nickname2  string
	Createtime time.Time
	Updatetime time.Time
	Phone      string
	Channelid  int
	Status     int
}

//TableName ..
func (*Account) TableName() string {
	return "accounts"
}

// //GetByUsername ..
// func (a *Account) GetByUsername(usernaem string) *Account {
// 	temp := new(Account)
// 	has, err := dbmanager.MasterDB.Where("username=?", usernaem).Get(temp)
// 	if err != nil {
// 		log.Logger.Errorln("Account GetByUsername is err:", err)
// 		log.Logger.Errorln("Account GetByUsername is has:", has)
// 	}
// 	return temp
// }

//GetByUserid ..
func (a *Account) GetByUserid(id int64) {
	//temp := new(Account)
	has, err := db.MasterDB.Where("userid=?", id).Get(a)
	if err != nil {
		log.WithFields(log.Fields{
			"has":    has,
			"err":    err,
			"userid": id,
		}).Error("Account [GetByUserid] is err")
	}
	return
}

//GetByUsername ..
func (a *Account) GetByUsername(username string) {
	//temp := new(Account)
	has, err := db.MasterDB.Where("username=?", username).Get(a)
	if err != nil {
		log.WithFields(log.Fields{
			"has":      has,
			"err":      err,
			"username": username,
		}).Error("Account [GetByUsername] is err")
	}
	return
}

//Insert 插入一条数据
func (a *Account) Insert() bool {
	affected, err := db.MasterDB.Insert(a)
	if err != nil {
		log.Logger.Errorln("Account [Insert] is err:", err)
	}
	if affected == 1 {
		return true
	}
	return false
}

//ExistByUsername 用户名查重
func (a *Account) ExistByUsername(username string) bool {
	has, err := db.MasterDB.Where("username = ?", username).Exist(&Account{})
	if err != nil {
		log.WithFields(log.Fields{
			"has":      has,
			"err":      err,
			"username": username,
		}).Error("Account [ExistByUsername] is err")
	}
	return has
}

//UpdateNickname ..
func (a *Account) UpdateNickname() bool {
	affected, err := db.MasterDB.Where("userid = ?", a.Userid).Cols("nickname").Update(a)
	if err != nil || affected != 1 {
		log.WithFields(log.Fields{
			"affected": affected,
			"err":      err,
			"id":       a.Userid,
		}).Error("Token.go [UpdateNickname] is err")
		return false
	}
	return true
}
