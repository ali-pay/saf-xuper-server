package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jason-cn-dev/xupercc/controllers"
)

func TestQueryAcl(t *testing.T) {

	req := controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		ContractAccount: "XC1234567812345678@xuper",
		//ContractName: "group_chain",
		//MethodName:   "listNode",
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**查询合约账户
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "contract_account": "XC1234567812345678@xuper"
}
*/

/**查询合约方法
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "contract_name": "group_chain",
    "method_name": "listNode"
}
*/
