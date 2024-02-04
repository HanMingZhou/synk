package main

import (
	"filesend/config"
	"filesend/server"
	"os"
	"os/exec"
	"os/signal"
)

func main() {

	// 启动gin服务
	go server.Run()

	// 启动chrome浏览器
	cmd := startBrowser()

	// 监听中断信号
	chSignal := ListneTosignalHandler()

	// 等待中断信号
	select {
	case <-chSignal: // 阻塞等待信号
		cmd.Process.Kill()
	}

}

func startBrowser() *exec.Cmd {
	// 先写死路径，后面在照着 lorca改
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// 执行chrome 127.0.0.1:27149
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.GetPort()+"/static/index.html")
	cmd.Start()
	return cmd
}

func ListneTosignalHandler() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
