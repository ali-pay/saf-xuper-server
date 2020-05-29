package v1

import (
	"io/ioutil"
	//"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/xuperchain/xuper-sdk-go/account"

	"xupercc/conf"
	"xupercc/controllers"
	log "xupercc/utils"
	"xupercc/xkernel"
)

var codeMaxSize = 2 << 20 // 2 MiB

func ContractDeploy(c *gin.Context) {

	req := new(controllers.Req)
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数无效",
		})
		log.Printf("param invalid, err: %s", err.Error())
		return
	}

	_, err = account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "助记词无效",
		})
		log.Printf("mnemonic can not retrieve account, err: %s", err.Error())
		return
	}

	//代码检查
	if !check(c, req) {
		return
	}

	//保存文件
	if !save(c, req) {
		return
	}

	//编译合约
	if !build(c, req) {
		return
	}

	//部署合约
	deploy(c, req)

}

//代码检查
func check(c *gin.Context, req *controllers.Req) bool {

	//格式检查
	if req.Runtime != "c" && req.Runtime != "go" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "合约格式无效，当前只支持c和go",
		})
		log.Printf("runtime invalid, current code runtime: %s", req.Runtime)
		return false
	}

	//名称检查
	ccname := filepath.Base(req.ContractName)
	r := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{4,16}$`)
	if !r.MatchString(ccname) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "合约名无效，只能包含字母/数字/下划线，且为字母开头，长度为4-16",
		})
		log.Printf("name invalid, current contract name: %s", req.ContractName)
		return false
	}

	//大小检查
	if req.ContractCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "没有获取到合约代码",
		})
		log.Printf("not found contract code")
		return false
	}

	if len(req.ContractCode) > codeMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "合约太大了，代码应当小于2M",
		})
		log.Printf("file size too big, current code size: %d", len(req.ContractCode))
		return false
	}

	return true
}

//保存文件
func save(c *gin.Context, req *controllers.Req) bool {
	//存放目录
	err := os.MkdirAll(conf.Code.CodePath, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
			"error": err.Error(),
		})
		log.Printf("mkdir fail, err: %s", err.Error())
		return false
	}

	//保存文件
	var filename string //文件保存名称
	if req.Runtime == "go" {
		filename = req.ContractName + ".go"
	} else {
		filename = req.ContractName + ".cc"
	}
	codefile := filepath.Join(conf.Code.CodePath, filename) //文件存放位置
	//codefile := conf.Code.CodePath + filename //担心CodePath会有问题，所有用filepath.Join
	code := []byte(req.ContractCode)
	err = ioutil.WriteFile(codefile, code, os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
			"error": err.Error(),
		})
		log.Printf("contract save fail, err: %s", err.Error())
		return false
	}

	log.Printf("contract saved, file: %s", codefile)
	return true
}

//编译合约
func build(c *gin.Context, req *controllers.Req) bool {
	var wasmfile string
	var err error
	if req.Runtime == "go" {
		filename := req.ContractName + ".go"
		err = controllers.BuildGo(filename)
		//wasmfile = conf.Code.WasmPath + req.ContractName + ".go.wasm"
		wasmfile = filepath.Join(conf.Code.WasmPath, req.ContractName+".go.wasm")
	} else {
		filename := req.ContractName + ".cc"
		err = controllers.BuildCC(filename)
		//wasmfile = conf.Code.WasmPath + req.ContractName + ".wasm"
		wasmfile = filepath.Join(conf.Code.WasmPath, req.ContractName+".wasm")
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "合约编译失败，请检查是否有误",
			"error": err.Error(),
		})
		log.Printf("build file fail, err: %s", err.Error())
		return false
	}

	log.Printf("file builded, file: %s", wasmfile)
	return true
}

//部署合约
func deploy(c *gin.Context, req *controllers.Req) {
	//获取身份
	acc, err := account.RetrieveAccount(req.Mnemonic, conf.Req.Language)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "助记词无效",
		})
		log.Printf("mnemonic can not retrieve account, err: %s", err.Error())
		return
	}

	//部署合约
	contract := xkernel.InitContract(acc, req.Node, req.BcName, req.ContractAccount, req.ContractName)
	//给服务费用的地址
	contract.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr = acc.Address
	//服务地址
	contract.Cfg.EndorseServiceHost = req.Node

	var action string
	if req.Upgrade {
		action = xkernel.UPGEADE
	} else {
		action = xkernel.DEPLOY
	}

	var wasmfile string
	if req.Runtime == "go" {
		//wasmfile = conf.Code.WasmPath + req.ContractName + ".go.wasm"
		wasmfile = filepath.Join(conf.Code.WasmPath, req.ContractName+".go.wasm")
	} else {
		//wasmfile = conf.Code.WasmPath + req.ContractName + ".wasm"
		wasmfile = filepath.Join(conf.Code.WasmPath, req.ContractName+".wasm")
	}
	txid, err := contract.ContractDoit(action, wasmfile, req.Runtime, req.Args)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "部署失败",
			"error": err.Error(),
		})
		log.Printf("deploy contract fail, err: %s", err.Error())
		return
	}

	log.Printf("deploy contract success, txid: %s", txid)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "部署成功",
		"resp": controllers.Result{Txid: txid},
	})
}
