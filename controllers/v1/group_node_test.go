package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"xupercc/controllers"
)

func TestGroupNode(t *testing.T) {
	req := &controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		Mnemonic:        "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractAccount: "XC1234567812345678@xuper",
		Method:          "list",
		Args: map[string]string{
			"bcname":  "xuper",
			"ip":      "ip1",
			"address": "add1",
		},
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**节点列表
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "contract_account": "XC1234567812345678@xuper",
    "args": {
        "bcname": "xuper"
    },
    "method": "list"
}
*/

/**添加节点
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "contract_account": "XC1234567812345678@xuper",
    "args": {
        "bcname": "wtf",
        "ip":"ip1",
        "address":"add1"
    },
    "method": "add"
}
*/

/**删除节点
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "mnemonic": "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
    "contract_account": "XC1234567812345678@xuper",
    "args": {
        "bcname": "wtf",
        "ip":"ip1"
    },
    "method": "del"
}
*/
