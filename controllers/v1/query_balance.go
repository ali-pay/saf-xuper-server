package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/transfer"
	"google.golang.org/grpc"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
)

func BalanceSDK(c *gin.Context) {
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

	trans := transfer.InitTrans(acc, req.Node, req.BcName)
	//给服务费用的地址
	trans.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	trans.Cfg.EndorseServiceHost = req.Node

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

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"resp": controllers.Result{
			AccountBalance: balance,
		},
	})
}

func Balance(c *gin.Context) {

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

	conn, err := grpc.Dial(req.Node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	//余额
	//addrstatus := &pb.AddressStatus{
	//	Address: req.Account,
	//	Bcs: []*pb.TokenDetail{
	//		{Bcname: req.BcName},
	//	},
	//}
	//是否获取冻结金额
	//fGetBalance := client.GetBalance
	//if req.Frozen {
	//	fGetBalance = client.GetFrozenBalance
	//}
	//reply, err := fGetBalance(ctx, addrstatus)

	//余额详情（包含了冻结金额）
	tfds := []*pb.TokenFrozenDetails{{Bcname: req.BcName}}
	addStatus := &pb.AddressBalanceStatus{
		Address: req.Account,
		Tfds:    tfds,
	}
	reply, err := client.GetBalanceDetail(ctx, addStatus)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": err.Error(),
		})
		log.Printf("query balance fail, err: %s", err.Error())
		return
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": reply.Header.Error.String(),
		})
		log.Printf("query balance fail, err: %s", reply.Header.Error.String())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		//"resp": reply.Bcs[0].Balance,
		"resp": reply.Tfds[0].Tfd,
	})
}
