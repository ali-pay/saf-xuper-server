package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jason-cn-dev/xupercc/controllers"
)

func TestTransfer(t *testing.T) {
	req := &controllers.Req{
		Node:     "127.0.0.1:37101",
		BcName:   "xuper",
		Mnemonic: "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		To:       "ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt",
		Amount:   100,
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
    "to": "ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt",
    "amount": 100
}
*/
