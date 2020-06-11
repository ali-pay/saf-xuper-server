package v0

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jason-cn-dev/xupercc/conf"
	"github.com/jason-cn-dev/xupercc/controllers"
	"github.com/jason-cn-dev/xupercc/utils"
)

var fileMaxSize int64 = 2 << 20 // 2 MiB

func Upload(c *gin.Context) {
	_, download := c.GetPostForm("download") //带上这个字段代表需要下载编译后的文件

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "没有获取到上传的文件",
		})
		log.Printf("not found upload file, err: %s", err.Error())
		return
	}

	//格式检查
	if !strings.HasSuffix(file.Filename, ".cc") && !strings.HasSuffix(file.Filename, ".go") {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件格式无效，当前只支持c++和go",
		})
		log.Printf("file format invalid, filename: %s", file.Filename)
		return
	}

	//大小检查
	if file.Size > fileMaxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件太大，合约文件应当小于2M",
		})
		log.Printf("file size too big, filesize: %d", file.Size)
		return
	}

	//名称检查
	filename := filepath.Base(file.Filename)
	r := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9.]{5,50}$")
	if !r.MatchString(filename) || strings.Count(filename, ".") > 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "文件名无效，只能包含字母和数字，且为字母开头，长度为5-50",
		})
		log.Printf("file name invalid, filename: %s", file.Filename)
		return
	}

	//存放目录
	err = os.MkdirAll(conf.Code.CodePath, 0755)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
		})
		log.Printf("mkdir fail, err: %s", err.Error())
		return
	}

	//保存文件
	codefile := conf.Code.CodePath + filename //文件存放位置
	err = c.SaveUploadedFile(file, codefile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
		})
		log.Printf("save file fail, err: %s", err.Error())
		return
	}

	log.Printf("file uploaded, file: %s", codefile)

	//编译合约
	var wasmfile string
	if strings.HasSuffix(filename, ".go") {
		err = controllers.BuildGo(filename)
		wasmfile = conf.Code.WasmPath + filename + ".wasm"
		c.Request.PostForm.Add("runtime", "go") //记录合约类型

	} else {
		err = controllers.BuildCC(filename)
		filename = strings.ReplaceAll(filename, ".cc", "") //去除".cc"
		wasmfile = conf.Code.WasmPath + filename + ".wasm"
		c.Request.PostForm.Add("runtime", "c")
	}

	c.Request.PostForm.Add("wasmfile", wasmfile) //记录文件地址

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  400,
			"msg":   "合约编译失败，请检查是否有误",
			"error": err.Error(),
		})
		log.Printf("build file fail, err: %s", err.Error())
		return
	}

	log.Printf("file builded, file: %s", wasmfile)

	//返回编译后的文件
	if download {
		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(wasmfile))
		c.File(wasmfile)
		return
	}

	//返回下载链接
	ip, err := utils.GetLocalIP()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
		})
		log.Printf("get local ip, err: %s", err.Error())
		return
	}

	url := conf.Server.Protocol + "://%s:%s/download/%s"
	filename = filepath.Base(wasmfile)
	downloadUrl := fmt.Sprintf(url, ip.String(), conf.Server.HttpPort, filename)
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "合约编译成功，请下载部署",
		"download": downloadUrl,
	})

	//部署合约
	//Deploy(c)

}
