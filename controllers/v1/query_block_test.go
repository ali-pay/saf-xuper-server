package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"xupercc/controllers"
)

func TestQueryBlock(t *testing.T) {
	req := &controllers.Req{
		Node:   "127.0.0.1:37101",
		BcName: "xuper",
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**根据区块id
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "block_id": "6aa48e13afce1d3e68c4f5ac16297513cea3b4b941ceccb02ddeeb28b0a40b74"
}
*/

/**根据区块高度
{
    "node": "127.0.0.1:37101",
    "bc_name": "xuper",
    "block_height": 0
}
*/
