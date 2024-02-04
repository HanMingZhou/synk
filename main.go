package main

import (
	"filesend/server"
	"os"
	"os/exec"
	"os/signal"
)

func main() {

	// 启动gin服务
	go func() {
		server.Run()
	}()

	// 启动chrome浏览器
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// 执行chrome 127.0.0.1:27149
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:27149/static/index.html")
	cmd.Start()

	// 监听中断信号
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)

	// 等待中断信号
	select {
	case <-chSignal: // 阻塞等待信号
		cmd.Process.Kill()
	}

}
