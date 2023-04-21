package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FilesController 文件处理
func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw")
	if err != nil {
		log.Fatal("err", err)
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
	fullpath := filepath.Join("uploads", filename+filepath.Ext(file.Filename)) // 拼接文件的绝对路径（不含exe 所在目录）
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))          // 保存文件
	fmt.Print("fullpath：", fullpath)
	fmt.Print("fileErr：", filepath.Join(dir, fullpath))
	if fileErr != nil {
		log.Fatal(err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "url": "/" + fullpath}) // 返回文件的绝对路径（ 不含exe 所在目录）
}
