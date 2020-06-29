package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/jason-cn-dev/xupercc/controllers/v1"
	"github.com/jason-cn-dev/xupercc/middlewares"
)

func NewRouter() *gin.Engine {

	//记录到文件
	//dir := path.Base(conf.Log.FilePath)
	//err := os.MkdirAll(dir, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//logFilePath := conf.Log.FilePath
	//logFileName := conf.Log.RouterFile
	//logfile := path.Join(logFilePath, logFileName)
	//f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	//if err != nil {
	//	panic(err)
	//}
	//
	//gin.DefaultWriter = io.MultiWriter(f) //输出到文件
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout) //输出到文件和控制台

	r := gin.Default() //不要使用logger
	//r := gin.New()
	//r.Use(gin.Recovery())

	// 中间件
	r.Use(middlewares.Cors())
	r.Use(middlewares.Logs())

	// 路由
	//r.POST("upload", v0.Upload)
	//r.GET("download/:filename", v0.Download)

	rv1 := r.Group("v1")
	rv1.POST("contract_deploy", v1.ContractDeploy)
	rv1.POST("contract_invoke", v1.ContractInvoke)
	rv1.POST("create_account", v1.CreateAccount)
	rv1.POST("create_contract_account", v1.CreateContractAccount)
	rv1.POST("balance", v1.Balance)
	rv1.POST("transfer", v1.Transfer)
	rv1.POST("query_tx", v1.QueryTx)
	rv1.POST("method_acl", v1.MethodAcl)
	rv1.POST("account_acl", v1.AccountAcl)
	rv1.POST("status", v1.Status)
	rv1.POST("group_chain", v1.GroupChain)
	rv1.POST("group_node", v1.GroupNode)
	rv1.POST("query_acl", v1.QueryAcl)
	rv1.POST("query_block", v1.QueryBlock)
	rv1.POST("query_list", v1.QueryLists)

	return r
}
