package test

import (
	"testing"

	"github.com/jason-cn-dev/xupercc/conf"
	log "github.com/jason-cn-dev/xupercc/utils"
)

//测试调用代码所在文件和行号的深度
func TestUtilLogs(t *testing.T) {

	conf.Init()
	log.Printf("hello","world","!")
}
/**
0 log.go:276
1 log.go:443
2 log.go:629
3 logs.go:43
4 utillogs_test.go:13
5 testing.go:909
6 asm_amd64.s:1357
7 ???:0
-1 extern.go:182

 */