package main

import (
	"fmt"
	"os"

	"github.com/samsaralc/hiksdk/pkg"
)

// 基础使用示例：登录、获取设备信息、登出
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 基础使用示例")
	fmt.Println("========================================")

	// 初始化 SDK（必须）
	fmt.Println("\n[1] 初始化 SDK...")
	pkg.InitHikSDK()
	defer func() {
		fmt.Println("\n[5] 释放 SDK 资源...")
		pkg.HKExit()
		fmt.Println("✓ 完成")
	}()
	fmt.Println("✓ SDK 初始化成功")

	// 设备连接信息
	deviceInfo := pkg.DeviceInfo{
		IP:       "192.168.1.64", // 替换为你的设备IP
		Port:     8000,           // 替换为你的端口
		UserName: "admin",        // 替换为你的用户名
		Password: "password",     // 替换为你的密码
	}

	fmt.Println("\n[2] 创建设备实例...")
	dev := pkg.NewHKDevice(deviceInfo)
	fmt.Println("✓ 设备实例创建成功")

	// 登录设备
	fmt.Println("\n[3] 登录设备...")
	loginId, err := dev.Login()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ 登录成功 (登录ID: %d)\n", loginId)

	defer func() {
		fmt.Println("\n[4] 登出设备...")
		if err := dev.Logout(); err != nil {
			fmt.Printf("✗ 登出失败: %v\n", err)
		} else {
			fmt.Println("✓ 登出成功")
		}
	}()

	// 获取设备信息
	fmt.Println("\n[4] 获取设备信息...")
	info, err := dev.GetDeiceInfo()
	if err != nil {
		fmt.Printf("✗ 获取设备信息失败: %v\n", err)
		return
	}

	fmt.Println("\n设备详细信息:")
	fmt.Printf("  - 设备名称: %s\n", info.DeviceName)
	fmt.Printf("  - 设备序列号: %s\n", info.DeviceID)
	fmt.Printf("  - 通道数量: %d\n", info.ByChanNum)
	fmt.Printf("  - IP 地址: %s\n", info.IP)
	fmt.Printf("  - 端口: %d\n", info.Port)
	fmt.Printf("  - 用户名: %s\n", info.UserName)

	// 获取通道名称
	fmt.Println("\n[5] 获取通道信息...")
	channels, err := dev.GetChannelName()
	if err != nil {
		fmt.Printf("✗ 获取通道信息失败: %v\n", err)
		return
	}

	fmt.Printf("\n通道列表 (共 %d 个):\n", len(channels))
	for id, name := range channels {
		fmt.Printf("  - 通道 %d: %s\n", id, name)
	}

	fmt.Println("\n========================================")
	fmt.Println("示例完成!")
	fmt.Println("========================================")
}
