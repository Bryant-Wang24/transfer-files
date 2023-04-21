package controllers

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TextsController 文本处理
func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}
	err := c.ShouldBindJSON(&json)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	exe, err := os.Executable() // 获取当前可执行文件的路径
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)                 // 获取当前可执行文件的目录
	filename := uuid.New().String()          // 生成一个文件名
	uploads := filepath.Join(dir, "uploads") // 拼接uploads目录的绝对路径
	err = os.MkdirAll(uploads, os.ModePerm)  // 创建uploads目录
	if err != nil {
		log.Fatal(err)
	}
	fullpath := filepath.Join("uploads", filename+".txt")                        // 拼接文件的绝对路径（不含exe 所在目录）
	err = ioutil.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 将json.Raw写入文件
	if err != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "url": "/" + fullpath}) // 返回文件的绝对路径（ 不含exe 所在目录）
}
