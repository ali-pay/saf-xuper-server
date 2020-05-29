package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"xupercc/controllers"
)

func TestMethodAcl(t *testing.T) {
	req := controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		Mnemonic:        "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractAccount: "XC1234567812345678@xuper",
		ContractName:    "group_chain",
		MethodName:      "listNode",
		//Address: []string{"ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt","dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"},
		Address: []string{"ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt"},
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
    "contract_account": "XC1234567812345678@xuper",
    "contract_name": "group_chain",
    "method_name": "listNode",
    "address": [
        "ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt"
    ]
}
*/
