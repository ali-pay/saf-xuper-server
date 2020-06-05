package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"xupercc/controllers"
)

func TestQueryTx(t *testing.T) {
	req := &controllers.Req{
		Node:     "127.0.0.1:37101",
		BcName:   "xuper",
		Mnemonic: "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		Txid:     "2f749eef25f6d9dabf197acc08e747793b663f36f4e67f820ab67e5fa65c9762",
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "txid": "2f749eef25f6d9dabf197acc08e747793b663f36f4e67f820ab67e5fa65c9762"
}
*/

/**
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "txid":"2f749eef25f6d9dabf197acc08e747793b663f36f4e67f820ab67e5fa65c9762"
}
*/
