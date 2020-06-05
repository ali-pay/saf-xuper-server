package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
)

func CreateAccount(c *gin.Context) {

	acc, err := account.CreateAccount(conf.Req.Strength, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  500,
			"msg":   "系统内部错误",
			"error": err.Error(),
		})
		log.Printf("create account fail, err: %s", err.Error())
		return
	}

	//log.Printf("create account success, addr: %s, mnemonic: %s", acc.Address, acc.Mnemonic)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"resp": controllers.Result{
			Mnemonic: acc.Mnemonic,
			Address:  acc.Address,
		}})
}
