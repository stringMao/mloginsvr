package confmanager

import (
	log "msigninser/common/logmanager"
	"os"
	"path/filepath"

	"github.com/Unknwon/goconfig"
)

//AppCfg 系统配置的全局变量
var AppCfg *goconfig.ConfigFile

//系统配置读取
func init() {
	var err error
	dirPath := filepath.Dir(os.Args[0])
	confPath, err := filepath.Abs(dirPath + "/config/app.ini")
	if err != nil {
		log.Logger.Fatalln("[app.ini]文件未找到：", err)
	}
	AppCfg, err = goconfig.LoadConfigFile(confPath)
	if err != nil {
		log.Logger.Fatalln("app.ini read err:", err)
	}

	//更新日志设置
	lv, err := AppCfg.GetValue("log", "level")
	if err != nil {
		log.Logger.Fatalln("read app.ini of log-lvevl is err:", err)
	}
	if l := log.SetLogLevel(lv); l == false {
		log.Logger.Fatalln("app.ini of log-lvevl is err:")
	}
}
