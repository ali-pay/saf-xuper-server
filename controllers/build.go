package controllers

import (
	"fmt"
	//"log"
	"os/exec"
	"path/filepath"
	"strings"

	"xupercc/conf"
	log "xupercc/utils"
)

//ccfile: 带后缀的是文件，不带的是包
//output: xx.wasm
func BuildCC(ccfile string) error {
	command := `./xdev/single.sh ` + ccfile //以main的启动目录判断
	cmd := exec.Command("/bin/bash", "-c", command)

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if !strings.Contains(string(output), "LD wasm") {
		return fmt.Errorf("build desc: cmd: %s output: %s", cmd.String(), string(output))
	}

	return nil
}

//output: xx.go.wasm
func BuildGo(gofile string) error {

	wasmfile := gofile + ".wasm" //go的合约加上后缀
	in := filepath.Join(conf.Code.CodePath, gofile)
	out := filepath.Join(conf.Code.WasmPath, wasmfile)
	cmdStr := fmt.Sprintf("GOOS=js GOARCH=wasm go build -o %s %s", out, in)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)

	log.Println(cmd.String())

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if string(output) != "" {
		return fmt.Errorf("build desc: cmd: %s output: %s", cmd.String(), string(output))
	}

	return nil
}
