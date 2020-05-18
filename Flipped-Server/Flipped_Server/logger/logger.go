package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var(
	Logger = logrus.New()
)

func InitLog() {
	// 设置日志格式为json格式
	Logger.Formatter = &logrus.JSONFormatter{}
	file, err := os.OpenFile("./log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		panic(err)
	}
}