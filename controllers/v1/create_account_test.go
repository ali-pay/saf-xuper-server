package v1

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/jason-cn-dev/xupercc/controllers"
)

func TestCreateAccount(t *testing.T) {
	req := &controllers.Req{
		Node: "127.0.0.1:37101",
	}

	//请求参数
	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

/**
{
    "node": "127.0.0.1:37101"
}
*/
