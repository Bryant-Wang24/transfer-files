package server

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"

	c "example.com/m/server/controllers"
	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run(Port string) {
	gin.SetMode(gin.DebugMode)
	router := gin.Default()
	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	router.StaticFS("/static", http.FS(staticFiles))
	router.POST("/api/v1/texts", c.TextsController)
	router.GET("/api/v1/addresses", c.AddressesController)
	router.GET("/uploads/:path", c.UploadsController)
	router.GET("/api/v1/qrcodes", c.QrcodesController)
	router.POST("/api/v1/files", c.FilesController)
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
	router.Run(":" + Port)
}
