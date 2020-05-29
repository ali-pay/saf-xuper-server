package v0

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"xupercc/conf"
	"xupercc/utils"
)

func Download(c *gin.Context) {
	filename := c.Param("filename")
	wasmfile := conf.Code.WasmPath + filename

	//检查是否存在
	isExist, err := utils.PathExists(wasmfile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统内部错误",
		})
		log.Printf("find file fail, err: %s", err.Error())
		return
	}

	if !isExist {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该合约文件不存在，请编译后再下载",
		})
		log.Printf("file not exist, file: %s", wasmfile)
		return
	}

	log.Printf("starting download, file: %s", wasmfile)

	//返回编译后的文件
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(wasmfile))
	c.File(wasmfile)
}
