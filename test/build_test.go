package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestBuildCC(t *testing.T) {
	command := `../xdev/single.sh ` + "erc20.cc"
	cmd := exec.Command("/bin/bash", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))

	if !strings.Contains(string(output), "LD wasm") {
		//编译失败
	}
}

func TestBuildPkg(t *testing.T) {
	command := `../xdev/single.sh ` + "xrc01"
	cmd := exec.Command("/bin/bash", "-c", command)

	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
		return
	}
	fmt.Printf("Execute Shell:%s finished with output:\n%s", command, string(output))

	if !strings.Contains(string(output), "LD wasm") {
		//编译失败
		//清理文件
		os.Remove("../contract_wasm/xrc01.wasm")
	}
}

func TestBuildGo(t *testing.T) {
	file := "erc20.go"
	wasm := file + ".wasm"
	in := "../contract_code/"
	out := "../contract_wasm/"
	cmd := exec.Command("/bin/bash", "-c", "GOOS=js GOARCH=wasm go build -o "+out+wasm+" "+in+file)

	fmt.Println(cmd.String())
	output, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("Execute Shell:%s finished with output:\n%s", cmd.String(), string(output))
}
