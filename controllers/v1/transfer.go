package v1

import (
	//"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/transfer"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
)

func Transfer(c *gin.Context) {

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

	//转账
	trans := transfer.InitTrans(acc, req.Node, req.BcName)
	//给服务费用的地址
	trans.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	trans.Cfg.EndorseServiceHost = req.Node

	amount := strconv.FormatInt(req.Amount, 10)
	fee := strconv.FormatInt(req.Fee, 10)
	txid, err := trans.Transfer(req.To, amount, fee, req.Desc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "转账失败",
			"error": err.Error(),
		})
		log.Printf("transfer fail, err: %s", err.Error())
		return
	}
	log.Printf("transfer success, txid: %s", txid)

	//查询余额
	balance, err := trans.GetBalance()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": err.Error(),
		})
		log.Printf("get balance fail, err: %s", err.Error())
		return
	}
	log.Printf("get balance success, balance: %s", balance)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "转账成功",
		"resp": controllers.Result{
			Txid:           txid,
			AccountBalance: balance,
		},
	})
}
