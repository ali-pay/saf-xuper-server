package test

import (
	"testing"

	"github.com/astaxie/beego/logs"
)

func TestLogs(t *testing.T) {
	//an official log.Logger
	l := logs.GetLogger()
	logs.SetLogger(logs.AdapterFile,`{
"filename":"app.log",
"level":7,
"color":true,
"separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]
}`)
	logs.EnableFuncCallDepth(true)

	l.Println("this is a message of http")
	//an official log.Logger with prefix ORM
	logs.GetLogger("ORM").Println("this is a message of orm")
	logs.Debug("my book is bought in the year of ", 2016)
	logs.Info("this %s cat is %v years old", "yellow", 3)
	logs.Warn("json is a type of kv like", map[string]int{"key": 2016})
	logs.Error(1024, "is a very", "good game")
	logs.Critical("oh,crash")
}
