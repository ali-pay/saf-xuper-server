package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
	"xupercc/xkernel"
)

func AccountAcl(c *gin.Context) {

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

	//所有地址的权限都是1
	ask := make(map[string]float32)
	for _, v := range req.Address {
		ask[v] = 1
	}

	acl := xkernel.InitAcl(acc, req.Node, req.BcName, req.ContractAccount)
	//给服务费用的地址
	acl.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	acl.Cfg.EndorseServiceHost = req.Node

	txid, err := acl.AclDoit(xkernel.ACCOUNT, req.ContractName, req.MethodName, ask)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "设置账户权限失败",
			"error": err.Error(),
		})
		log.Printf("set account acl fail, err: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "设置成功",
		"resp": controllers.Result{
			Txid: txid,
		},
	})
}
