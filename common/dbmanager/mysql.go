package dbmanager

import (
	"fmt"
	"mloginsvr/common/confmanager"
	log "mloginsvr/common/logmanager"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

//MasterDB ..
var MasterDB *xorm.Engine

//InitMysql 数据库连接
func InitMysql() {
	var err error
	host, err := confmanager.AppCfg.GetValue("mysql", "host")
	if err != nil {
		log.Logger.Fatalln("read app.ini of mysql-host is err:", err)
	}
	port, err := confmanager.AppCfg.GetValue("mysql", "port")
	if err != nil {
		log.Logger.Fatalln("read app.ini of mysql-port is err:", err)
	}
	username, err := confmanager.AppCfg.GetValue("mysql", "username")
	if err != nil {
		log.Logger.Fatalln("read app.ini of mysql-username is err:", err)
	}
	password, err := confmanager.AppCfg.GetValue("mysql", "password")
	if err != nil {
		log.Logger.Fatalln("read app.ini of mysql-password is err:", err)
	}
	dbname, err := confmanager.AppCfg.GetValue("mysql", "dbname")
	if err != nil {
		log.Logger.Fatalln("read app.ini of mysql-dbname is err:", err)
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", username, password, host, port, dbname)
	MasterDB, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		log.Logger.Fatalln("mysql connet is err:", err)
	}
	//engine.SetMapper(names.SameMapper{})//"xorm.io/xorm/names"

	MasterDB.SetMaxIdleConns(10)  //连接池的空闲数大小
	MasterDB.SetMaxOpenConns(100) //最大打开连接数

	log.Logger.Info("mysql init success")
}

// type User struct {
// 	Id      int64
// 	Name    string `xorm:"name varchar(254)"`
// 	Salt    string
// 	AgeId   int
// 	Passwd  string    `xorm:"varchar(200)"`
// 	Created time.Time `xorm:"created1"`
// 	Updated time.Time `xorm:"'updated'"`
// }

// func Test() {
// 	err := MasterDB.Sync2(new(User))
// 	if err != nil {
// 		log.Logger.Error(err)
// 	}
// }
