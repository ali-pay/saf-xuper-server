package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jason-cn-dev/xuper-sdk-go/pb"
	"google.golang.org/grpc"

	"github.com/jason-cn-dev/xupercc/controllers"
	log "github.com/jason-cn-dev/xupercc/utils"
)

func QueryAcl(c *gin.Context) {
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

	aclStatus := &pb.AclStatus{
		Bcname:       req.BcName,
		AccountName:  req.ContractAccount,
		ContractName: req.ContractName,
		MethodName:   req.MethodName,
	}

	if len(aclStatus.AccountName) == 0 && len(aclStatus.ContractName) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		log.Printf("param invalid, contractAccount or contractName is empty")
		return
	}

	var action string
	if aclStatus.AccountName != "" {
		action = "查询合约账户acl"
	}

	if aclStatus.ContractName != "" && aclStatus.MethodName != "" {
		action = "查询合约方法acl"
	}
	//fmt.Println(action)

	conn, err := grpc.Dial(req.Node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	reply, err := client.QueryACL(ctx, aclStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": err.Error(),
		})
		log.Printf("query acl fail, err: %s", err.Error())
		return
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": reply.Header.Error.String(),
		})
		log.Printf("query acl fail, err: %s", reply.Header.Error.String())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  action + "成功",
		"resp": reply.GetAcl(),
	})
}
