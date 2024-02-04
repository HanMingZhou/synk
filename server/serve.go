package server

import (
	"embed"
	c "filesend/server/controller"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var FS embed.FS

func Run() {

	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	staticFiles, _ := fs.Sub(FS, "frontend/dist")
	r.StaticFS("/static", http.FS(staticFiles))
	r.POST("/api/v1/files", c.FilesController)
	r.GET("/api/v1/qrcodes", c.QrcodesController)
	r.GET("/uploads/:path", c.UploadController)
	r.GET("/api/v1/addresses", c.AddressesController)
	r.POST("/api/v1/texts", c.TextsController)
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/static/") {
			reader, err := staticFiles.Open("index.html")
			if err != nil {
				log.Fatal(err)
			}
			defer reader.Close()
			stat, err := reader.Stat()
			if err != nil {
				log.Fatal(err)
			}
			c.DataFromReader(http.StatusOK, stat.Size(), "text/html", reader, nil)
		} else {
			c.Status(http.StatusNotFound)
		}
	})

	r.Run(":27149")
}
