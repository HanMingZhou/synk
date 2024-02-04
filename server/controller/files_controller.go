package controller

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FilesController(c *gin.Context) {
	file, err := c.FormFile("raw")
	if err != nil {
		log.Fatal(err)
	}
	exe, _ := os.Executable() // 当前exe执行的路径
	dir := filepath.Dir(exe)  // 当前exe执行的目录

	filename := uuid.New().String()          // 生成文件的新名字
	uploads := filepath.Join(dir, "uploads") // 新路径：dir/uploads
	err = os.MkdirAll(uploads, os.ModePerm)  // 创建 dir/uploads目录
	if err != nil {
		log.Fatal(err)
	}
	fullpath := path.Join("uploads", filename+filepath.Ext(file.Filename)) // 文件的full路径
	fileErr := c.SaveUploadedFile(file, filepath.Join(dir, fullpath))      // 保存上传的文件到具体的路径
	if fileErr != nil {
		log.Fatal(fileErr)
	}
	c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath})

}
