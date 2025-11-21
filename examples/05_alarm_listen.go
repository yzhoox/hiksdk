package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/samsaralc/hiksdk/core"
)

// 报警监听示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 报警监听示例")
	fmt.Println("========================================")

	// 设备连接信息
	deviceInfo := core.DeviceInfo{
		IP:       "192.168.1.64",
		Port:     8000,
		UserName: "admin",
		Password: "password",
	}

	// 创建设备并登录
	dev := core.NewHKDevice(deviceInfo)
	loginId, err := dev.LoginV40()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Printf("登录成功 (ID: %d)\n", loginId)
	defer dev.Logout()

	// 设置报警回调
	fmt.Println("\n设置报警回调...")
	if err := dev.SetAlarmCallBack(); err != nil {
		fmt.Printf("设置回调失败: %v\n", err)
		return
	}

	// 启动报警监听
	fmt.Println("启动报警监听...")
	if err := dev.StartListenAlarmMsg(); err != nil {
		fmt.Printf("启动监听失败: %v\n", err)
		return
	}
	defer dev.StopListenAlarmMsg()

	fmt.Println("\n监听中... 按 Ctrl+C 退出")
	fmt.Println("等待接收报警消息（移动侦测、遮挡报警等）")

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n停止监听...")
	fmt.Println("示例完成!")
}
