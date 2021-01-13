package log

import (
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

//Logger 全局日志对象
var Logger = logrus.New()

func init() {
	//Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.AddHook(newLfsHook("logs/"))
}

//Init ..
func Init(lv string) {
	//更新日志设置
	logrusLogLevel, err := logrus.ParseLevel(lv)
	if err != nil {
		Logger.Fatalln("app.ini of log-lvevl is err:", err)
	}
	Logger.SetLevel(logrusLogLevel) //设置等级
}

// 设置日志文件切割及软链接
func newLfsHook(filepath string) logrus.Hook {
	var err error
	//===debuglog======================================
	logpath := filepath + "debug/"
	writerDebug, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),       // 生成软链，指向最新日志文件
		rotatelogs.WithRotationTime(time.Hour), //设置日志分割的时间，这里设置为一小时分割一次
		//WithMaxAge和WithRotationCount二者只能设置一个，
		rotatelogs.WithMaxAge(time.Hour*12), // 文件最大保存时间
	)
	if err != nil {
		logrus.Errorf("writerDebug logger error. %+v", errors.WithStack(err))
	}

	//====infolog===================================
	logpath = filepath + "info/"
	writerInfo, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerInfo logger error. %+v", errors.WithStack(err))
	}

	//====Errlog===================================
	logpath = filepath + "error/"
	writerErr, err := rotatelogs.New(
		logpath+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithRotationTime(time.Hour*24),
		rotatelogs.WithMaxAge(time.Hour*24*30),
	)
	if err != nil {
		logrus.Errorf("writerErr logger error. %+v", errors.WithStack(err))
	}

	/*
		logrus 拥有六种日志级别：debug、info、warn、error、fatal 和 panic，
		logrus.Debug(“Useful debugging information.”)
		logrus.Info(“Something noteworthy happened!”)
		logrus.Warn(“You should probably take a look at this.”)
		logrus.Error(“Something failed but I'm not quitting.”)
		logrus.Fatal(“Bye.”) //log之后会调用os.Exit(1)
		logrus.Panic(“I'm bailing.”) //log之后会panic()
	*/
	//设置默认等级
	logrusLogLevel, _ := logrus.ParseLevel("debug")
	Logger.SetLevel(logrusLogLevel) //设置等级

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writerDebug,
		logrus.InfoLevel:  writerInfo,
		logrus.WarnLevel:  writerErr,
		logrus.ErrorLevel: writerErr,
		logrus.FatalLevel: writerErr,
		logrus.PanicLevel: writerErr,
	}, &logrus.TextFormatter{DisableColors: true})

	return lfsHook
}

//Fields ..
type Fields map[string]interface{}

//WithFields 重写此函数，便于使用
func WithFields(fields Fields) *logrus.Entry {
	return Logger.WithFields(logrus.Fields(fields))
}

//WithField 重写此函数，便于使用
func WithField(key string, value interface{}) *logrus.Entry {
	return Logger.WithField(key, value)
}
