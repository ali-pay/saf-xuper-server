package v0

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jason-cn-dev/xuper-sdk-go/account"
	"github.com/jason-cn-dev/xuper-sdk-go/contract"

	"github.com/jason-cn-dev/xupercc/conf"
	"github.com/jason-cn-dev/xupercc/controllers"
)

func Deploy(c *gin.Context) {

	req := new(controllers.Req)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数无效",
		})
		log.Printf("param invalid, err: %s", err.Error())
		return
	}

	req.Args = c.PostFormMap("args")
	req.Runtime = c.PostForm("runtime")
	req.ContractFile = c.PostForm("wasmfile")

	//获取身份
	acc, err := account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "助记词无效",
		})
		log.Printf("mnemonic can not retrieve account, err: %s", err.Error())
		return
	}

	wasmContract := contract.InitWasmContract(acc, req.Node, req.BcName, req.ContractName, req.ContractAccount)

	//升级合约
	if req.Upgrade {
		log.Printf("升级合约，未实现")
		return
	}

	//部署合约
	txid, err := wasmContract.DeployWasmContract(req.Args, req.ContractFile, req.Runtime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "部署失败",
			"error": err.Error(),
		})
		log.Printf("deploy contract fail, err: %s", err.Error())
		return
	}

	log.Printf("deploy contract success, txid: %s", txid)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "部署成功",
		"resp": Result{txid},
	})
}

type Result struct {
	Txid string `json:"txid"`
}
