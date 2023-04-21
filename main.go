package main

import (
	"embed"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

//go:embed frontend/dist/*
var FS embed.FS

func main() {
	port := "27149"
	go func() { // 开启一个gin协程，防止阻塞调起chrome
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		router.StaticFS("/static", http.FS(staticFiles))
		router.POST("/api/v1/texts", TextsController)
		router.GET("/api/v1/addresses", AddressesController)
		router.GET("/uploads/:path", UploadsController)
		router.GET("/api/v1/qrcodes", QrcodesController)
		router.POST("/api/v1/files", FilesController)
		router.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/static") {
				reader, err := staticFiles.Open("index.html")
				if err != nil {
					log.Fatal("err", err)
				}
				defer reader.Close()
				// 获取文件大小
				stat, err := reader.Stat()
				if err != nil {
					log.Fatal("err", err)
				}
				c.DataFromReader(http.StatusOK, stat.Size(), "text/html;charset=utf-8", reader, nil)
			} else {
				c.Status(http.StatusNotFound)
			}
		})
		router.Run(":" + port)
	}()
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+port+"/static")
	cmd.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}

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

// AddressesController 地址处理
func AddressesController(c *gin.Context) {
	addrs, _ := net.InterfaceAddrs()
	var result []string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "ok", "addresses": result})
}

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

// QrcodesController 二维码生成
func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256) // 生成二维码
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png) //返回二维码图片
		fmt.Println("content", content)
	} else {
		c.Status(http.StatusBadRequest)
	}
}

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
