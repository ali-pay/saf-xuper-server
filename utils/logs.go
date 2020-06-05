package utils

import (
	"os"
	"path"
	"sync"

	"github.com/astaxie/beego/logs"

	"xupercc/conf"
)

//使用beego的log模块
func initLogs() {
	dir := path.Base(conf.Log.FilePath)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logFilePath := conf.Log.FilePath
	logFileName := conf.Log.RunTimeFile
	logfile := path.Join(logFilePath, logFileName)
	configs := `{"filename":"` + logfile + `","color":true,"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`

	err = logs.SetLogger(logs.AdapterFile, configs)
	if err != nil {
		panic(err)
	}
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(4)
}

//改成全局变量，避免被回收
var Once sync.Once

func Printf(format string, v ...interface{}) {

	Once.Do(initLogs)

	switch len(v) {
	case 0:
		logs.Info(format)
	case 1:
		logs.Info(format, v[0])
	case 2:
		logs.Info(format, v[0], v[1])
	case 3:
		logs.Info(format, v[0], v[1], v[2])
	case 4:
		logs.Info(format, v[0], v[1], v[2], v[3])
	case 5:
		logs.Info(format, v[0], v[1], v[2], v[3], v[4])
	default:
		logs.Info(format, v)
	}

}

func Println(v ...interface{}) {
	Once.Do(initLogs)
	logs.Info(v)
}
