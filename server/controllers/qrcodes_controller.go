package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

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
