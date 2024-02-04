package main

import (
	"embed"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
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
	go func() {
		gin.SetMode(gin.DebugMode)
		r := gin.Default()

		staticFiles, _ := fs.Sub(FS, "frontend/dist")
		r.StaticFS("/static", http.FS(staticFiles))
		r.POST("/api/v1/files", FilesController)
		r.GET("/api/v1/qrcodes", QrcodesController)
		r.GET("/uploads/:path", UploadController)
		r.GET("/api/v1/addresses", AddressesController)
		r.POST("/api/v1/texts", TextsController)
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

		r.Run(":" + port)
	}()

	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// 执行chrome 127.0.0.1:8080
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+port+"/static/index.html")
	cmd.Start()

	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	select {
	case <-chSignal: // 阻塞等待信号
		cmd.Process.Kill()
	}

}

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

func QrcodesController(c *gin.Context) {
	if content := c.Query("content"); content != "" {
		png, err := qrcode.Encode(content, qrcode.Medium, 256)
		if err != nil {
			log.Fatal(err)
		}
		c.Data(http.StatusOK, "image/png", png)
	} else {
		c.Status(http.StatusBadRequest)
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

// AddressesController
func AddressesController(c *gin.Context) {
	addrs, err := net.InterfaceAddrs() // 获取所有的ip地址
	if err != nil {
		log.Fatal(err)
	}
	var result []string
	for _, address := range addrs { // 遍历所有的ip
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				result = append(result, ipnet.IP.String())
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"addresses": result})
}

// TextController()
func TextsController(c *gin.Context) {
	var json struct {
		Raw string `json:"raw"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		exe, err := os.Executable() // 获取当前文件的路径
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exe) // 获取当前文件的目录
		if err != nil {
			log.Fatal(err)
		}

		filename := uuid.New().String()          // 生成一个文件名
		uploads := filepath.Join(dir, "uploads") // 拼接 uploads 的绝对路径
		err = os.MkdirAll(uploads, os.ModePerm)  // 创建 uploads 目录
		if err != nil {
			log.Fatal(err)
		}

		fullpath := path.Join("uploads", filename+".txt")                        // 拼接文件的绝对路径(不含目录)
		err = os.WriteFile(filepath.Join(dir, fullpath), []byte(json.Raw), 0644) // 写入文件
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, gin.H{"url": "/" + fullpath}) // 返回文件的绝对路径（不含 exe 所在目录）

	}
}
