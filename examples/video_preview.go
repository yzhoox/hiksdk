package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/samsaralc/hiksdk/pkg"
)

// 视频预览示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 视频预览示例")
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

	// 获取设备信息
	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		fmt.Println("✗ 设备没有可用通道")
		os.Exit(1)
	}

	fmt.Printf("✓ 设备有 %d 个通道\n", info.ByChanNum)

	// 选择要预览的通道
	channelId := 1
	fmt.Printf("\n使用通道 %d 进行视频预览\n", channelId)

	// 创建接收器
	receiver := &pkg.Receiver{}
	err = receiver.Start()
	if err != nil {
		fmt.Printf("✗ 启动接收器失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ 接收器启动成功")

	// 启动实时预览
	fmt.Println("\n启动实时视频预览...")
	startTime := time.Now()
	realHandle, err := dev.RealPlay_V40(channelId, receiver)
	if err != nil {
		fmt.Printf("✗ 启动预览失败: %v\n", err)
		os.Exit(1)
	}

	if realHandle < 0 {
		fmt.Printf("✗ 预览句柄无效: %d\n", realHandle)
		os.Exit(1)
	}

	connectTime := time.Since(startTime)
	fmt.Printf("✓ 视频预览已启动 (耗时: %v)\n", connectTime)
	fmt.Printf("  - 预览句柄: %d\n", realHandle)

	// 确保停止预览
	defer func() {
		fmt.Println("\n停止视频预览...")
		dev.StopRealPlay()
		fmt.Println("✓ 视频预览已停止")
	}()

	// 监听系统信号以便优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("\n========================================")
	fmt.Println("视频预览运行中...")
	fmt.Println("========================================")
	fmt.Println("\n提示:")
	fmt.Println("  - 视频数据正在接收中")
	fmt.Println("  - 按 Ctrl+C 可以停止预览并退出")
	fmt.Println("----------------------------------------\n")

	// 统计信息
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	previewStartTime := time.Now()
	dataReceivedCount := 0

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n收到中断信号，正在退出...")
			return

		case data, ok := <-receiver.PSMouth:
			if !ok {
				fmt.Println("\n✗ 接收通道已关闭")
				return
			}
			dataReceivedCount++

			// 每收到100个数据包打印一次
			if dataReceivedCount%100 == 0 {
				elapsed := time.Since(previewStartTime)
				rate := float64(dataReceivedCount) / elapsed.Seconds()
				fmt.Printf("\r[运行时间: %v] 已接收 %d 个数据包 (%.2f pkt/s, 最新包: %d bytes)    ",
					elapsed.Round(time.Second), dataReceivedCount, rate, len(data))
			}

		case <-ticker.C:
			if dataReceivedCount == 0 {
				fmt.Printf("\r⚠ 警告: 未收到任何视频数据，请检查网络和设备状态                    ")
			} else {
				elapsed := time.Since(previewStartTime)
				avgRate := float64(dataReceivedCount) / elapsed.Seconds()
				fmt.Printf("\r[运行时间: %v] 总计: %d 包, 平均: %.2f pkt/s    ",
					elapsed.Round(time.Second), dataReceivedCount, avgRate)
			}
		}
	}
}
