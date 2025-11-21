package main

import (
	"fmt"

	"github.com/samsaralc/hiksdk/core"
)

// 获取设备信息示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 获取设备信息示例")
	fmt.Println("========================================")

	// 设备连接信息
	deviceInfo := core.DeviceInfo{
		IP:       "192.168.1.64",
		Port:     8000,
		UserName: "admin",
		Password: "password",
	}

	// 创建设备实例
	dev := core.NewHKDevice(deviceInfo)

	// 登录设备
	fmt.Println("\n[1] 登录设备...")
	loginId, err := dev.LoginV40()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
		return
	}
	fmt.Printf("✓ 登录成功 (ID: %d)\n", loginId)
	defer dev.Logout()

	// 获取设备基本信息
	fmt.Println("\n[2] 获取设备信息...")
	info, err := dev.GetDeviceInfo()
	if err != nil {
		fmt.Printf("✗ 获取失败: %v\n", err)
		return
	}

	fmt.Println("\n设备基本信息:")
	fmt.Printf("  设备名称: %s\n", info.DeviceName)
	fmt.Printf("  序列号: %s\n", info.DeviceID)
	fmt.Printf("  IP地址: %s\n", info.IP)
	fmt.Printf("  端口: %d\n", info.Port)
	fmt.Printf("  通道数: %d\n", info.ByChanNum)

	// 获取通道名称
	fmt.Println("\n[3] 获取通道名称...")
	channelNames, err := dev.GetChannelName()
	if err == nil && len(channelNames) > 0 {
		fmt.Println("\n通道列表:")
		for ch, name := range channelNames {
			fmt.Printf("  通道 %d: %s\n", ch, name)
		}
	} else {
		fmt.Println("  无法获取通道名称")
	}

	fmt.Println("\n示例完成!")
}
