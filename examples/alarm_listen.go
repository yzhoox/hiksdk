package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samsaralc/hiksdk/pkg"
)

// 报警监听示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 报警监听示例")
	fmt.Println("========================================")

	// 初始化 SDK
	pkg.InitHikSDK()
	defer pkg.HKExit()

	// 设备连接信息
	deviceInfo := pkg.DeviceInfo{
		IP:       "192.168.1.64", // 替换为你的设备IP
		Port:     8000,
		UserName: "admin",
		Password: "password",
	}

	// 登录设备
	dev := pkg.NewHKDevice(deviceInfo)
	loginId, err := dev.Login()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
		os.Exit(1)
	}
	defer dev.Logout()

	fmt.Printf("✓ 登录成功 (登录ID: %d)\n", loginId)

	// 设置报警回调函数
	fmt.Println("\n设置报警回调函数...")
	err = dev.SetAlarmCallBack()
	if err != nil {
		fmt.Printf("✗ 设置报警回调失败: %v\n", err)
		fmt.Println("注意: 某些设备可能不支持报警回调功能")
	} else {
		fmt.Println("✓ 报警回调函数设置成功")
	}

	// 启动报警监听
	fmt.Println("\n启动报警监听...")
	err = dev.StartListenAlarmMsg()
	if err != nil {
		fmt.Printf("✗ 启动报警监听失败: %v\n", err)
		fmt.Println("\n可能的原因:")
		fmt.Println("  1. 设备不支持该功能")
		fmt.Println("  2. 设备固件版本过低")
		fmt.Println("  3. 报警输入未配置")
		os.Exit(1)
	}
	fmt.Println("✓ 报警监听已启动")

	// 确保停止监听
	defer func() {
		fmt.Println("\n停止报警监听...")
		err := dev.StopListenAlarmMsg()
		if err != nil {
			fmt.Printf("✗ 停止报警监听失败: %v\n", err)
		} else {
			fmt.Println("✓ 报警监听已停止")
		}
	}()

	// 监听系统信号以便优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("\n========================================")
	fmt.Println("报警监听已启动，等待报警事件...")
	fmt.Println("========================================")
	fmt.Println("\n提示:")
	fmt.Println("  - 程序将持续监听报警事件")
	fmt.Println("  - 触发设备报警输入以测试报警回调")
	fmt.Println("  - 按 Ctrl+C 可以退出程序")
	fmt.Println("\n支持的报警类型:")
	fmt.Println("  - 移动侦测报警")
	fmt.Println("  - 遮挡报警")
	fmt.Println("  - 音频异常报警")
	fmt.Println("  - 硬盘满报警")
	fmt.Println("  - 硬盘故障报警")
	fmt.Println("  - 视频信号丢失报警")
	fmt.Println("  - 输入/输出报警")
	fmt.Println("----------------------------------------\n")

	// 定时打印状态
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n收到中断信号，正在退出...")
			return
		case <-ticker.C:
			elapsed := time.Since(startTime)
			fmt.Printf("\r[运行时间: %v] 等待报警事件... (按 Ctrl+C 退出)    ",
				elapsed.Round(time.Second))
		}
	}
}
