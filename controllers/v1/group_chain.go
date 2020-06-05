package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/contract"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
)

func GroupChain(c *gin.Context) {

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

	acc, err := account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "助记词无效",
		})
		log.Printf("mnemonic can not retrieve account, err: %s", err.Error())
		return
	}

	req.ContractName = "group_chain"

	switch req.Method {
	case "list":
		req.MethodName = "listChain"
		req.Query = true
	case "del":
		req.MethodName = "delChain"
	case "add":
		req.MethodName = "addChain"
	}

	wasmContract := contract.InitWasmContract(acc, req.Node, req.BcName, req.ContractName, req.ContractAccount)
	//给服务费用的地址
	wasmContract.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	wasmContract.Cfg.EndorseServiceHost = req.Node

	if req.Query {
		query(c, req, wasmContract)
	} else {
		invoke(c, req, wasmContract, "listChain")
	}
}
