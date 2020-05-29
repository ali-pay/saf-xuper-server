package v1

import (
	"encoding/json"
	"fmt"
	"testing"
	"xupercc/controllers"
)

func TestJsonNotCode(t *testing.T) {
	req := controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		Mnemonic:        "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractAccount: "XC1234567812345678@xuper",
		ContractName:    "name",
		ContractCode:    "code",
		Args:            map[string]string{"key": "vaule"},
		Runtime:         "c",
		Upgrade:         false,
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
    "contract_name": "name",
    "contract_code": "code",
    "args": {
        "key": "vaule"
    },
    "runtime": "c"
}
*/

func TestJsonCC(t *testing.T) {
	req := controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		Mnemonic:        "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractAccount: "XC1234567812345678@xuper",
		ContractName:    "ccCounter",
		ContractCode:    ccCounter,
		Args:            map[string]string{"creator": "xuper"},
		Runtime:         "c",
		Upgrade:         false,
	}

	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

func TestJsonGO(t *testing.T) {
	req := controllers.Req{
		Node:            "127.0.0.1:37101",
		BcName:          "xuper",
		Mnemonic:        "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南",
		ContractAccount: "XC1234567812345678@xuper",
		ContractName:    "goCounter",
		ContractCode:    goCounter,
		Args:            map[string]string{"creator": "xuper"},
		Runtime:         "go",
		Upgrade:         false,
	}

	jsonByte, _ := json.Marshal(req)
	fmt.Println(string(jsonByte))
}

var ccCounter = `#include "xchain/xchain.h"

struct Counter : public xchain::Contract {};

DEFINE_METHOD(Counter, initialize) {
    xchain::Context* ctx = self.context();
    const std::string& creator = ctx->arg("creator");
    if (creator.empty()) {
        ctx->error("missing creator");
        return;
    }
    ctx->put_object("creator", creator);
    ctx->ok("initialize succeed");
}

DEFINE_METHOD(Counter, increase) {
    xchain::Context* ctx = self.context();
    const std::string& key = ctx->arg("key");
    std::string value;
    ctx->get_object(key, &value);
    int cnt = 0;
    cnt = atoi(value.c_str());
    char buf[32];
    snprintf(buf, 32, "%d", cnt + 1);
    ctx->put_object(key, buf);
    ctx->ok(buf);
}

DEFINE_METHOD(Counter, get) {
    xchain::Context* ctx = self.context();
    const std::string& key = ctx->arg("key");
    std::string value;
    if (ctx->get_object(key, &value)) {
        ctx->ok(value);
    } else {
        ctx->error("key not found");
    }
}
`

var goCounter = `package main

import (
	"strconv"

	"github.com/xuperchain/xuperchain/core/contractsdk/go/code"
	"github.com/xuperchain/xuperchain/core/contractsdk/go/driver"
)

type counter struct{}

func (c *counter) Initialize(ctx code.Context) code.Response {
	creator, ok := ctx.Args()["creator"]
	if !ok {
		return code.Errors("missing creator")
	}
	err := ctx.PutObject([]byte("creator"), creator)
	if err != nil {
		return code.Error(err)
	}
	return code.OK(nil)
}

func (c *counter) Increase(ctx code.Context) code.Response {
	key, ok := ctx.Args()["key"]
	if !ok {
		return code.Errors("missing key")
	}
	value, err := ctx.GetObject(key)
	cnt := 0
	if err == nil {
		cnt, _ = strconv.Atoi(string(value))
	}

	cntstr := strconv.Itoa(cnt + 1)

	err = ctx.PutObject(key, []byte(cntstr))
	if err != nil {
		return code.Error(err)
	}
	return code.OK([]byte(cntstr))
}

func (c *counter) Get(ctx code.Context) code.Response {
	key, ok := ctx.Args()["key"]
	if !ok {
		return code.Errors("missing key")
	}
	value, err := ctx.GetObject(key)
	if err != nil {
		return code.Error(err)
	}
	return code.OK(value)
}

func main() {
	driver.Serve(new(counter))
}
`
