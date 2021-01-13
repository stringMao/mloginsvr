package config

import (
	"io/ioutil"
	"mloginsvr/common/log"
	"os"
	"path/filepath"
	"strings"
)

//SenWords 敏感词库
var SenWords map[string]bool = make(map[string]bool)

//LoadSensitiveWords 加载敏感词库
func LoadSensitiveWords() {
	pathstr, err := AppCfg.GetValue("server", "wordsfile")
	if err != nil {
		log.Logger.Fatalln("read app.ini of server-wordsfile is err:", err)
	}
	secpath := strings.Split(pathstr, ",")
	for _, v := range secpath {
		str := getWordsFileContent(v)
		sec := strings.Split(str, ",")
		for _, v2 := range sec {
			SenWords[v2] = true
		}
	}

	log.Logger.Debugf("敏感词库:%+v", SenWords)
}

func getWordsFileContent(filename string) string {
	var err error
	dirPath := filepath.Dir(os.Args[0])
	filepath, err := filepath.Abs(dirPath + "/src/words/" + filename)
	if err != nil {
		log.Logger.Fatalf("[%s]文件路径描述错误：", filename, err)
	}
	f, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Logger.Fatalf("读取文件[%s]失败:%#v", filename, err)
		return ""
	}
	str := string(f)
	// 去除空格
	str = strings.Replace(str, " ", "", -1)
	// 去除换行符
	str = strings.Replace(str, "\n", "", -1)
	return str
}

//HaveSenWords 敏感词检查
func HaveSenWords(str string) bool {
	if SenWords == nil || len(SenWords) == 0 {
		return false
	}

	for i := 0; i < len(str); i++ {
		for j := i + 1; j <= len(str); j++ {
			subStr := str[i:j]
			if _, found := SenWords[subStr]; found {
				return true
			}
		}
	}
	return false
}
