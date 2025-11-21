package main

import (
	"fmt"
	"os"

	"github.com/samsaralc/hiksdk/core"
)

// 两种登录方式示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - 两种登录方式示例")
	fmt.Println("========================================")

	// 设备连接信息
	deviceInfo := core.DeviceInfo{
		IP:       "192.168.1.64", // 替换为你的设备IP
		Port:     8000,           // 替换为你的端口
		UserName: "admin",        // 替换为你的用户名
		Password: "password",     // 替换为你的密码
	}

	fmt.Println("\n设备连接信息:")
	fmt.Printf("  - IP地址: %s\n", deviceInfo.IP)
	fmt.Printf("  - 端口: %d\n", deviceInfo.Port)
	fmt.Printf("  - 用户名: %s\n", deviceInfo.UserName)

	// ==================== 方式1: LoginV40 (V40推荐) ====================
	fmt.Println("\n========================================")
	fmt.Println("方式1: 使用 NET_DVR_Login_V40（推荐）")
	fmt.Println("========================================")

	fmt.Println("\n[1] 创建设备实例...")
	dev1 := core.NewHKDevice(deviceInfo)
	fmt.Println("✓ 设备实例创建成功（SDK已自动初始化）")

	// 使用V40登录
	fmt.Println("\n[2] 使用Login()方法登录（V40）...")
	loginId1, err := dev1.LoginV40()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)

		// 如果是HKError，显示详细信息
		if hkErr, ok := err.(*core.HKError); ok {
			fmt.Printf("\n错误详情:\n")
			fmt.Printf("  错误代码: %d\n", hkErr.Code)
			fmt.Printf("  错误消息: %s\n", hkErr.Msg)
			fmt.Printf("  操作: %s\n", hkErr.Operation)
			fmt.Printf("  设备IP: %s\n", hkErr.IP)
		}

		fmt.Println("\n可能的原因:")
		fmt.Println("  1. 设备不在线或网络不可达")
		fmt.Println("  2. 用户名或密码错误")
		fmt.Println("  3. 设备端口配置错误")
	} else {
		fmt.Printf("✓ 登录成功 (登录ID: %d)\n", loginId1)

		// 获取设备信息
		fmt.Println("\n[3] 获取设备信息...")
		info, err := dev1.GetDeviceInfo()
		if err == nil {
			fmt.Printf("  设备名称: %s\n", info.DeviceName)
			fmt.Printf("  序列号: %s\n", info.DeviceID)
			fmt.Printf("  通道数: %d\n", info.ByChanNum)
		}

		// 登出设备
		fmt.Println("\n[4] 登出设备...")
		if err := dev1.Logout(); err != nil {
			fmt.Printf("✗ 登出失败: %v\n", err)
		} else {
			fmt.Println("✓ 登出成功")
		}
	}

	// ==================== 方式2: LoginV30 (兼容旧设备) ====================
	fmt.Println("\n========================================")
	fmt.Println("方式2: 使用 NET_DVR_Login_V30（兼容旧设备）")
	fmt.Println("========================================")

	fmt.Println("\n[1] 创建设备实例...")
	dev2 := core.NewHKDevice(deviceInfo)
	fmt.Println("✓ 设备实例创建成功")

	// 使用V30登录
	fmt.Println("\n[2] 使用LoginV30()方法登录...")
	loginId2, err := dev2.LoginV30()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)

		// 显示错误详情
		if hkErr, ok := err.(*core.HKError); ok {
			fmt.Printf("错误详情: %s\n", hkErr.JSON())
		}
	} else {
		fmt.Printf("✓ 登录成功 (登录ID: %d)\n", loginId2)

		// 获取设备信息
		fmt.Println("\n[3] 获取设备信息...")
		info, err := dev2.GetDeviceInfo()
		if err == nil {
			fmt.Printf("  设备名称: %s\n", info.DeviceName)
			fmt.Printf("  序列号: %s\n", info.DeviceID)
			fmt.Printf("  通道数: %d\n", info.ByChanNum)
		}

		// 登出设备
		fmt.Println("\n[4] 登出设备...")
		if err := dev2.Logout(); err != nil {
			fmt.Printf("✗ 登出失败: %v\n", err)
		} else {
			fmt.Println("✓ 登出成功")
		}
	}

	// ==================== 对比说明 ====================
	fmt.Println("\n========================================")
	fmt.Println("两种登录方式对比")
	fmt.Println("========================================")

	fmt.Println("\nLoginV40() [NET_DVR_Login_V40]:")
	fmt.Println("  ✓ 推荐使用")
	fmt.Println("  ✓ 支持更多功能")
	fmt.Println("  ✓ 更好的性能")
	fmt.Println("  ✓ 支持同步/异步登录")
	fmt.Println("  ✓ 设备信息更详细")
	fmt.Println("  ✓ 适用于新设备")

	fmt.Println("\nLoginV30() [NET_DVR_Login_V30]:")
	fmt.Println("  ✓ 兼容旧设备")
	fmt.Println("  ✓ 简单直接")
	fmt.Println("  ✓ 适用于老版本设备")
	fmt.Println("  ✗ 功能相对较少")
	fmt.Println("  ✗ 只支持同步登录")

	fmt.Println("\n💡 建议:")
	fmt.Println("  1. 优先使用 LoginV40() 方法（V40）")
	fmt.Println("  2. 如果 LoginV40() 失败，可尝试 LoginV30()")
	fmt.Println("  3. 对于确定是旧设备的，直接使用 LoginV30()")

	fmt.Println("\n💡 注意事项:")
	fmt.Println("  - 设备最多支持32个注册用户名")
	fmt.Println("  - 同时最多允许128个用户注册")
	fmt.Println("  - SDK支持2048个注册，UserID取值范围0~2047")
	fmt.Println("  - 登录后务必调用 Logout() 释放资源")

	// 程序结束时清理SDK
	defer core.Cleanup()
	os.Exit(0)
}
