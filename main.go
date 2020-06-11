package main

import (
	"github.com/jason-cn-dev/xupercc/conf"
	"github.com/jason-cn-dev/xupercc/routers"
)

func main() {
	// 从配置文件读取配置
	conf.Init()

	// 装载路由
	r := routers.NewRouter()
	r.Run(":" + conf.Server.HttpPort)
}
