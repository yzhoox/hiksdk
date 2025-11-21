package main

import (
	"fmt"
	"time"

	"github.com/samsaralc/hiksdk/core"
)

// 视频预览示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 视频预览示例")
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

	// 创建接收器
	receiver := &core.Receiver{}
	if err := receiver.Start(); err != nil {
		fmt.Printf("启动接收器失败: %v\n", err)
		return
	}
	defer receiver.Stop()

	// 启动实时预览
	fmt.Println("\n启动实时预览...")
	channelId := 1 // 通道1
	previewHandle, err := dev.RealPlay_V40(channelId, receiver)
	if err != nil {
		fmt.Printf("启动预览失败: %v\n", err)
		return
	}
	fmt.Printf("预览启动成功 (句柄: %d)\n", previewHandle)

	// 接收视频数据
	fmt.Println("\n接收视频流数据...")
	fmt.Println("按 Ctrl+C 退出")

	// 统计信息
	packetCount := 0
	totalSize := 0
	startTime := time.Now()

	// 接收数据
	go func() {
		for data := range receiver.PSMouth {
			packetCount++
			totalSize += len(data)

			// 每100个包打印一次统计
			if packetCount%100 == 0 {
				duration := time.Since(startTime)
				rate := float64(totalSize) / duration.Seconds() / 1024 / 1024
				fmt.Printf("\r数据包: %d, 总大小: %.2f MB, 速率: %.2f MB/s",
					packetCount, float64(totalSize)/1024/1024, rate)
			}
		}
	}()

	// 运行10秒后退出
	time.Sleep(10 * time.Second)

	// 停止预览
	fmt.Println("\n\n停止预览...")
	dev.StopRealPlay()

	// 打印最终统计
	duration := time.Since(startTime)
	fmt.Println("\n统计信息:")
	fmt.Printf("  运行时间: %.1f 秒\n", duration.Seconds())
	fmt.Printf("  数据包数: %d\n", packetCount)
	fmt.Printf("  数据总量: %.2f MB\n", float64(totalSize)/1024/1024)
	fmt.Printf("  平均速率: %.2f MB/s\n", float64(totalSize)/duration.Seconds()/1024/1024)

	fmt.Println("\n示例完成!")
}
