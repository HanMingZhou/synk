package main

import (
	"filesend/config"
	"filesend/server"
	"os"
	"os/exec"
	"os/signal"
)

func main() {
	// 利用channel构造两个goroutine之间的通信
	// chChromeDie 监听chrome是否关闭
	// chBackDie 监听后端是否关闭
	chChromeDie := make(chan struct{}) // 通常使用结构体类型的channel
	chBackendDie := make(chan struct{})

	// 启动gin服务
	go server.Run()

	// 启动chrome浏览器
	go startBrowser(chChromeDie, chBackendDie)

	// 监听ctrl+c中断信号
	chSignal := ListneTosignalHandler()

	// 等待中断信号
	for {
		select {
		case <-chSignal: // 阻塞等待信号
			chBackendDie <- struct{}{}
		case <-chChromeDie: // 等待chrome关闭后的channel通信
			os.Exit(0)
		}
	}

}
func startBrowser(chChromeDie chan struct{}, chBackendDie chan struct{}) {
	// 先写死路径，后面在照着 lorca改
	chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// 执行chrome 127.0.0.1:27149
	cmd := exec.Command(chromePath, "--app=http://127.0.0.1:"+config.GetPort()+"/static/index.html")
	cmd.Start()
	go func() {
		<-chBackendDie
		cmd.Process.Kill()
	}()
	go func() {
		cmd.Wait()
		chChromeDie <- struct{}{}
	}()
}
func ListneTosignalHandler() chan os.Signal {
	chSignal := make(chan os.Signal, 1)
	signal.Notify(chSignal, os.Interrupt)
	return chSignal
}
