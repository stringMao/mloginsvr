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
func (a *Account) GetByUsername(usernaem string) {
	//temp := new(Account)
	has, err := db.MasterDB.Where("username=?", usernaem).Get(a)
	if err != nil {
		log.WithFields(log.Fields{
			"has":      has,
			"err":      err,
			"usernaem": usernaem,
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
