package v1

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jason-cn-dev/xuper-sdk-go/account"

	"github.com/jason-cn-dev/xupercc/conf"
	"github.com/jason-cn-dev/xupercc/controllers"
	log "github.com/jason-cn-dev/xupercc/utils"
	"github.com/jason-cn-dev/xupercc/xkernel"
)

func CreateContractAccount(c *gin.Context) {

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

	//有时候随机数会是0开头，fmt会截断它，所以使用for来跳过不符合长度的数
	for len(req.ContractAccount) != 16 {
		//req.ContractAccount = "1234567812345678"
		req.ContractAccount = fmt.Sprintf("%08v", rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000000000))
	}
	ca := xkernel.InitAcl(acc, req.Node, req.BcName, req.ContractAccount)
	//给服务费用的地址
	ca.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	ca.Cfg.EndorseServiceHost = req.Node

	gas, acl, txid, err := ca.CreateContractAccount()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "创建失败",
			"error": err.Error(),
		})
		log.Printf("create contract account fail, err: %s", err.Error())
		return
	}

	//log.Printf("create contract account success, txid: %s", txid)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"resp": controllers.Result{
			Txid:            txid,
			AccountAcl:      acl,
			GasUsed:         gas,
			ContractAccount: "XC" + req.ContractAccount + "@" + req.BcName,
		},
	})
}
