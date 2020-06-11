package test

import (
	"testing"

	"github.com/jason-cn-dev/xuper-sdk-go/account"
	"github.com/jason-cn-dev/xuper-sdk-go/contract"
)

func TestAccount(t *testing.T) {
	acc, err := account.RetrieveAccount("致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(acc.Address) //ZUjrEbucZYBxF6U7YJKCuSJYbBQewAMWt
}

var (
	node            = "127.0.0.1:37101"
	bcname          = "xuper"
	contractName    = "counter"
	contractAccount = "XC1234567812345678@xuper"
	args            = map[string]string{"creator": "xchain"}
	contractFile    = "../contract_wasm/counter.wasm"
	runtime         = "c"
	mnemonic        = "致 端 全 刘 积 旁 扰 蔬 伪 欢 近 南"
	language        = 1
)

func TestDeploy(t *testing.T) {
	acc, _ := account.RetrieveAccount(mnemonic, language)
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)
	txid, err := wasmContract.DeployWasmContract(args, contractFile, runtime)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txid)
}
