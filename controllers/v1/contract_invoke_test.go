package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"xupercc/controllers"
)

func TestContractInvoke(t *testing.T) {
	req := &controllers.Req{
		Node:         "127.0.0.1:37101",
		BcName:       "xuper",
		Mnemonic:     "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractName: "group_chain",
		MethodName:   "addNode",
		Args: map[string]string{
			"bcname":  "wtf",
			"ip":      "ip",
			"address": "address",
		},
		Query: false,
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**查询
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "contract_name": "group_chain",
    "method_name": "listNode",
    "args": {
        "bcname": "wtf"
    },
    "query": true
}
*/

/**调用
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "contract_name": "group_chain",
    "method_name": "addNode",
    "args": {
        "address": "address",
        "bcname": "wtf",
        "ip": "ip"
    }
}
*/
