package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// GetUploadsDir 获取上传文件的目录
func GetUploadsDir() (uploads string) {
	exe, err := os.Executable() // 获取当前可执行文件的路径
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe) // 获取当前可执行文件的目录
	uploads = filepath.Join(dir, "uploads")
	return
}

// UploadsController 文件下载
func UploadsController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-Disposition", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename="+path)
		c.Header("Content-Type", "application/octet-stream")
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}
