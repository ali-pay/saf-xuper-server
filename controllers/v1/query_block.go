package v1

import (
	"context"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jason-cn-dev/xuper-sdk-go/pb"
	"google.golang.org/grpc"

	"github.com/jason-cn-dev/xupercc/controllers"
	log "github.com/jason-cn-dev/xupercc/utils"
)

func QueryBlock_D(c *gin.Context) {
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

	var action string
	var reply *pb.Block

	if req.BlockID != "" {
		action = "根据区块ID查询"
		rawBlockid, err := hex.DecodeString(req.BlockID)
		if err != nil {

		}
		blockIDPB := &pb.BlockID{
			Bcname:      req.BcName,
			Blockid:     rawBlockid,
			NeedContent: true,
		}
		reply, err = client.GetBlock(ctx, blockIDPB)

	} else if req.BlockHeight >= 0 {
		action = "根据区块高度查询"
		blockHeightPB := &pb.BlockHeight{
			Bcname: req.BcName,
			Height: req.BlockHeight,
		}
		reply, err = client.GetBlockByHeight(ctx, blockHeightPB)
	} else {
		err = errors.New("参数无效，区块id或区块高度不能为空")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": err.Error(),
		})
		log.Printf("query block fail, err: %s", err.Error())
		return
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "查询失败",
			"error": reply.Header.Error.String(),
		})
		log.Printf("query block fail, err: %s", reply.Header.Error.String())
		return
	}

	if reply.Block == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "查询失败，找不到该区块",
		})
		log.Printf("block not found")
		return
	}

	iblock := log.FromInternalBlockPB(reply.Block)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  action + "成功",
		"resp": iblock,
	})
}

func GetChainBlock(node, bcname, id string, height int64) (*log.InternalBlock, error) {
	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	var reply *pb.Block

	if id != "" {
		rawBlockid, err := hex.DecodeString(id)
		if err != nil {
			return nil, errors.New("区块id无效")
		}
		blockIDPB := &pb.BlockID{
			Bcname:      bcname,
			Blockid:     rawBlockid,
			NeedContent: true,
		}
		reply, err = client.GetBlock(ctx, blockIDPB)

	} else if height >= 0 {
		blockHeightPB := &pb.BlockHeight{
			Bcname: bcname,
			Height: height,
		}
		reply, err = client.GetBlockByHeight(ctx, blockHeightPB)
	} else {
		err = errors.New("参数无效，区块id或区块高度不能为空")
	}

	if err != nil {
		log.Printf("query block fail, err: %s", err.Error())
		return nil, err
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		log.Printf("query block fail, err: %s", reply.Header.Error.String())
		return nil, errors.New("查询失败")
	}

	if reply.Block == nil {
		log.Printf("block not found")
		return nil, errors.New("查询失败，找不到该区块")
	}

	return log.SimpleBlock(reply.Block), nil
}

func QueryBlock(c *gin.Context) {
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

	iblock, err := GetChainBlock(req.Node, req.BcName, req.BlockID, req.BlockHeight)
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
		"resp": iblock,
	})
}
