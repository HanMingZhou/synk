package controller

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// UploadController
func UploadController(c *gin.Context) {
	if path := c.Param("path"); path != "" {
		target := filepath.Join(GetUploadsDir(), path)
		c.Header("Content-type", "application/octet-stream")
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Transfer-Encoding", "binary")
		c.Header("Content-Disposition", "attachment; filename"+path)
		c.File(target)
	} else {
		c.Status(http.StatusNotFound)
	}
}

func GetUploadsDir() (uploads string) {
	exe, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Dir(exe)
	uploads = filepath.Join(dir, "uploads")
	return uploads
}
