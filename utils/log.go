package utils

import (
	"log"
	"os"
	"path"
	"xupercc/conf"
)

//使用官方的log模块，加入日志保存
var logger *log.Logger

func initLog() {
	dir := path.Base(conf.Log.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFilePath := conf.Log.FilePath
	logFileName := conf.Log.RunTimeFile
	logfile := path.Join(logFilePath, logFileName)
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logger = log.New(f, "", log.Llongfile)
	logger.SetFlags(log.LstdFlags)
}

//func Printf(format string, v ...interface{}) {
//	if logger == nil {
//		initLog()
//	}
//
//	switch len(v) {
//	case 0:
//		logger.Println(format)
//	case 1:
//		logger.Printf(format, v[0])
//	case 2:
//		logger.Printf(format, v[0], v[1])
//	case 3:
//		logger.Printf(format, v[0], v[1], v[2])
//	case 4:
//		logger.Printf(format, v[0], v[1], v[2], v[3])
//	case 5:
//		logger.Printf(format, v[0], v[1], v[2], v[3], v[4])
//	default:
//		logger.Printf(format, v)
//	}
//
//}

//func Println(v ...interface{}) {
//	if logger == nil {
//		initLog()
//	}
//	logger.Println(v)
//}
