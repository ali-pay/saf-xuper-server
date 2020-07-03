package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jason-cn-dev/xuper-sdk-go/pb"
	"google.golang.org/grpc"

	"github.com/jason-cn-dev/xupercc/controllers"
	log "github.com/jason-cn-dev/xupercc/utils"
)

func Status_Old(c *gin.Context) {

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

	status := log.FromSystemStatusPB(reply.GetSystemsStatus())

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"resp": status,
	})

}

func GetChainStatus(node, bcname string) (*log.SystemStatus, error) {
	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("can not connect to node, err: %s", err.Error())
		return nil, errors.New("查询失败，无法连接到该节点")
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
	if bcname != "" {
		bcStatusPB := &pb.BCStatus{Bcname: bcname}
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
		log.Printf("query node status fail, err: %s", err.Error())
		return nil, errors.New("查询失败")
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		msg := "查询失败"
		if reply.Header.Error.String() == "CONNECT_REFUSE" {
			msg = "该链不存在"
		}
		log.Printf("query node status fail, err: %s", reply.Header.Error.String())
		return nil, errors.New(msg)
	}

	return log.FromSystemStatusPB(reply.GetSystemsStatus()), nil
}

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

	status, err := GetChainStatus(req.Node, req.BcName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"resp": status,
	})
}
