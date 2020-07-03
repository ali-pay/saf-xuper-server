package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/jason-cn-dev/xupercc/controllers"
	log "github.com/jason-cn-dev/xupercc/utils"
)

func QueryLists(c *gin.Context) {

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

	//缓存链数据
	ql, ok := lists[req.BcName]
	if !ok {
		ql = NewQueryList(req.Node, req.BcName)
		lists[req.BcName] = ql
	}

	//获取链状态
	status, err := GetChainStatus(req.Node, req.BcName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  err.Error(),
		})
		return
	}

	//先将最新的区块缓存起来,因为要逆序插入缓存列表
	var blocks []*log.InternalBlock

	//过滤不匹配的链
	var chain *log.ChainStatus
	for _, c := range status.ChainStatus {
		if c.Name != req.BcName {
			continue
		}
		chain = &c
	}
	if chain == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该链不存在",
		})
		return
	}

	//获取高度
	height := chain.Block.Height
	//计算需要获取的父区块数
	diff := ql.IsNew(height)
	//log.Println("diff:",diff)
	//判断是否需要获取父区块
	if diff != 0 {
		//缓存第一个区块
		//blocks = append(blocks, chain.Block)
		//获取父区块 因为已经获取了最新的区块了，所以这里判断要大于1才去获取更多父区块
		//问题：由于获取链数据的区块是没有包含交易详情的，所以这里重新获取一次，判断该为大于0，并将高度递减放到后面
		for i := diff; i > 0; i-- {
			//log.Println("获取区块",height)
			block, err := GetChainBlock(req.Node, req.BcName, "", height)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"code": 500,
					"msg":  err.Error(),
				})
				return
			}
			//加入临时缓存
			blocks = append(blocks, block)
			//高度递减
			height--
		}

		//插入区块缓存列表
		ql.AddBlock(blocks)

		//插入交易缓存列表
		count := 0
		var txs []*log.Transaction
	out:
		for _, block := range blocks {
			for _, tx := range block.Transactions {
				count++
				if count > size {
					break out
				}
				txs = append(txs, tx)
			}
		}
		ql.AddTx(txs)
	}

	//构造应答数据
	resp := []interface{}{ql.GetBlocks(), ql.GetTxs()}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "查询成功",
		"resp": resp,
	})
}
