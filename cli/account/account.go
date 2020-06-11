package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	accountsdk "github.com/jason-cn-dev/xuper-sdk-go/account"
	p2p_base "github.com/xuperchain/xuperchain/core/p2p/base"
)

//创建账号
func main() {
	output := "output/%s-%s/%s/"
	t := time.Now().Format("20060102150405")

	acc, err := accountsdk.CreateAccount(3, 1)
	if err != nil {
		panic(err)
	}

	//fmt.Printf("Address: %s\nMnemonic: %s\nPublicKey: %s\nPrivateKey: %s\n",
	//	acc.Address,
	//	acc.Mnemonic,
	//	acc.PublicKey,
	//	acc.PrivateKey)

	keys := fmt.Sprintf(output, t, acc.Address, "keys")
	err = os.MkdirAll(keys, os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(keys+"address", []byte(acc.Address), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(keys+"mnemonic", []byte(acc.Mnemonic), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(keys+"public.key", []byte(acc.PublicKey), os.ModePerm)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(keys+"private.key", []byte(acc.PrivateKey), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("create account success, output:", keys)


	netkeys := fmt.Sprintf(output, t, acc.Address, "netkeys")
	err = p2p_base.GenerateKeyPairWithPath(netkeys)
	if err != nil {
		panic(err)
	}

	peerid, err := p2p_base.GetPeerIDFromPath(netkeys)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(netkeys+"peerid", []byte(peerid), os.ModePerm)
	if err != nil {
		panic(err)
	}

	fmt.Println("create netUrl success, output:", netkeys, "peerid:", peerid)
}
