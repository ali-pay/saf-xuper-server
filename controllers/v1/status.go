package v1

import (
	"context"
	//"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuperchain/core/pb"
	"google.golang.org/grpc"

	"xupercc/controllers"
	log "xupercc/utils"
)

func Status(c *gin.Context) {

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
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败，无法连接到该节点",
			"error": err.Error(),
		})
		log.Printf("can not connect to node, err: %s", err.Error())
		return
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	reply := &pb.SystemsStatusReply{
		SystemsStatus: &pb.SystemsStatus{
			BcsStatus: make([]*pb.BCStatus, 0),
		},
	}

	//查询单条链
	if req.BcName != "" {
		bcStatusPB := &pb.BCStatus{Bcname: req.BcName}
		var bcStatus *pb.BCStatus
		bcStatus, err = client.GetBlockChainStatus(ctx, bcStatusPB)
		reply.SystemsStatus.BcsStatus = append(reply.SystemsStatus.BcsStatus, bcStatus)
		if bcStatus != nil {
			reply.Header = bcStatus.Header
		}

	} else {
		//查询所有链
		reply, err = client.GetSystemStatus(ctx, &pb.CommonIn{})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": err.Error(),
		})
		log.Printf("query node status fail, err: %s", err.Error())
		return
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		msg := "查询失败"
		if reply.Header.Error.String() == "CONNECT_REFUSE" {
			msg = "该链不存在"
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   msg,
			"error": reply.Header.Error.String(),
		})
		log.Printf("query node status fail, err: %s", reply.Header.Error.String())
		return
	}

	status := controllers.FromSystemStatusPB(reply.GetSystemsStatus())

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"resp": status,
	})

}
