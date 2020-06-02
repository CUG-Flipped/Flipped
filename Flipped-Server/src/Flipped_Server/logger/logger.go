package logger

import (
	"Flipped_Server/initialSetting"
	"github.com/sirupsen/logrus"
	"os"
)

var(
	Logger = logrus.New()
	logFilePath string
)

func initSettingsLogger() {
	logFilePath = initialSetting.LoggerConfig["filePath"].(string)
}

func InitLog() {
	initSettingsLogger()
	// 设置日志格式为json格式
	Logger.Formatter = &logrus.JSONFormatter{}
	Logger.SetReportCaller(true)
	//file, err := os.OpenFile("./server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	file, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		Logger.SetOutput(file)
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		panic(err)
	}
}

func SetToLogger(level logrus.Level, function string, msg string, otherMsg string){
	switch level {
		case logrus.DebugLevel:{
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Debug(otherMsg)
		}
		case logrus.InfoLevel:{
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Info(otherMsg)
		}
		case logrus.WarnLevel:{
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Warning(otherMsg)
		}
		case logrus.ErrorLevel:{
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Error(otherMsg)
		}
		case logrus.FatalLevel:{
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Fatal(otherMsg)
		}
		default: {
			Logger.WithFields(logrus.Fields{
				"function": function,
				"msg": msg,
			}).Debug(otherMsg)
		}
	}
}