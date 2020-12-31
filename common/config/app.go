package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/Unknwon/goconfig"
)

//AppCfg 系统配置的全局变量
var AppCfg *goconfig.ConfigFile

//Init 系统配置读取
func Init() {
	var err error
	dirPath := filepath.Dir(os.Args[0])
	confPath, err := filepath.Abs(dirPath + "/app.ini")
	if err != nil {
		log.Fatalln("[app.ini]文件未找到：", err)
	}
	AppCfg, err = goconfig.LoadConfigFile(confPath)
	if err != nil {
		log.Fatalln("app.ini read err:", err)
	}
}
