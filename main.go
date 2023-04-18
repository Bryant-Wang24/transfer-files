package main

import (
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() { // 开启一个gin协程，防止阻塞调起chrome
		gin.SetMode(gin.ReleaseMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("a test page"))
		})
		router.Run(":8080")
	}()
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/")
	cmd.Start()
	select {}
}
