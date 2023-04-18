package main

import (
	"os"
	"os/exec"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	go func() { // 开启一个gin协程，防止阻塞调起chrome
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("a test page"))
		})
		router.Run(":8080")
	}()
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:8080/")
	cmd.Start()
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	select {
	case <-chSignal:
		cmd.Process.Kill()
	}
}
