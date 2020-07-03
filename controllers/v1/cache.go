package v1

import (
	"sync"

	"github.com/jason-cn-dev/xupercc/conf"
	"github.com/jason-cn-dev/xupercc/utils"
)

var (
	size  = conf.Cache.Size         //缓存个数
	lists = map[string]*QueryList{} //key:bcname value:data
)

type QueryList struct {
	node   string
	bcname string
	blocks []*utils.InternalBlock
	txs    []*utils.Transaction
	lock   *sync.RWMutex
}

func NewQueryList(node, bcname string) *QueryList {
	return &QueryList{
		node:   node,
		bcname: bcname,
		blocks: make([]*utils.InternalBlock, size),
		txs:    make([]*utils.Transaction, size),
		lock:   new(sync.RWMutex),
	}
}

func (this *QueryList) AddBlock(blocks []*utils.InternalBlock) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.blocks = append(blocks, this.blocks[:size-len(blocks)]...)
}

func (this *QueryList) GetBlocks() []*utils.InternalBlock {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.blocks
}

//判断最新区块的高度是否高于缓存中的区块，返回相差的高度（需要获取的区块个数）
func (this *QueryList) IsNew(height int64) int64 {
	block := this.blocks[0]
	if block == nil {
		return int64(size)
	}
	diff := height - block.Height
	if diff > int64(size) {
		return int64(size)
	}
	return diff
}

func (this *QueryList) AddTx(txs []*utils.Transaction) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.txs = append(txs, this.txs[:size-len(txs)]...)

}

func (this *QueryList) GetTxs() []*utils.Transaction {
	this.lock.RLock()
	defer this.lock.RUnlock()
	return this.txs
}
