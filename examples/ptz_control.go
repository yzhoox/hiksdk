package main

import (
	"fmt"
	"os"
	"time"

	"github.com/samsaralc/hiksdk/pkg"
)

// PTZ云台控制示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - PTZ 云台控制示例")
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
	_, err := dev.Login()
	if err != nil {
		fmt.Printf("✗ 登录失败: %v\n", err)
		os.Exit(1)
	}
	defer dev.Logout()

	fmt.Println("✓ 登录成功")

	// 获取通道数
	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		fmt.Println("✗ 设备没有可用通道")
		os.Exit(1)
	}

	channelId := 1 // 使用通道1

	// PTZ 命令常量
	const (
		TILT_UP     = 21 // 云台上仰
		TILT_DOWN   = 22 // 云台下俯
		PAN_LEFT    = 23 // 云台左转
		PAN_RIGHT   = 24 // 云台右转
		ZOOM_IN     = 11 // 焦距变大
		ZOOM_OUT    = 12 // 焦距变小
		SET_PRESET  = 8  // 设置预置点
		GOTO_PRESET = 39 // 转到预置点
	)

	// 1. 云台右转
	fmt.Println("\n[示例1] 云台右转...")
	success, err := dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 4)
	if err != nil {
		fmt.Printf("✗ 云台右转失败: %v\n", err)
	} else if success {
		fmt.Println("✓ 云台右转开始")
		time.Sleep(2 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 1, 4) // 停止
		fmt.Println("✓ 云台已停止")
	}

	time.Sleep(500 * time.Millisecond)

	// 2. 云台左转
	fmt.Println("\n[示例2] 云台左转...")
	success, err = dev.PTZControlWithSpeed_Other(channelId, PAN_LEFT, 0, 4)
	if err != nil {
		fmt.Printf("✗ 云台左转失败: %v\n", err)
	} else if success {
		fmt.Println("✓ 云台左转开始")
		time.Sleep(2 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, PAN_LEFT, 1, 4) // 停止
		fmt.Println("✓ 云台已停止")
	}

	time.Sleep(500 * time.Millisecond)

	// 3. 云台上仰
	fmt.Println("\n[示例3] 云台上仰...")
	success, err = dev.PTZControlWithSpeed_Other(channelId, TILT_UP, 0, 4)
	if err != nil {
		fmt.Printf("✗ 云台上仰失败: %v\n", err)
	} else if success {
		fmt.Println("✓ 云台上仰开始")
		time.Sleep(2 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, TILT_UP, 1, 4) // 停止
		fmt.Println("✓ 云台已停止")
	}

	time.Sleep(500 * time.Millisecond)

	// 4. 变焦控制
	fmt.Println("\n[示例4] 焦距放大...")
	success, err = dev.PTZControl_Other(channelId, ZOOM_IN, 0)
	if err != nil {
		fmt.Printf("✗ 焦距放大失败: %v\n", err)
	} else if success {
		fmt.Println("✓ 焦距放大开始")
		time.Sleep(1 * time.Second)
		dev.PTZControl_Other(channelId, ZOOM_IN, 1) // 停止
		fmt.Println("✓ 焦距已停止")
	}

	time.Sleep(500 * time.Millisecond)

	// 5. 设置预置点
	fmt.Println("\n[示例5] 设置预置点...")
	presetId := 1
	success, err = dev.PTZControl_Other(channelId, SET_PRESET, presetId)
	if err != nil {
		fmt.Printf("✗ 设置预置点失败: %v\n", err)
	} else if success {
		fmt.Printf("✓ 预置点 %d 设置成功\n", presetId)

		// 移动云台到其他位置
		fmt.Println("\n[示例6] 移动云台到其他位置...")
		dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 3)
		time.Sleep(3 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 1, 3)

		time.Sleep(1 * time.Second)

		// 转到预置点
		fmt.Printf("\n[示例7] 转到预置点 %d...\n", presetId)
		success, err = dev.PTZControl_Other(channelId, GOTO_PRESET, presetId)
		if err != nil {
			fmt.Printf("✗ 转到预置点失败: %v\n", err)
		} else if success {
			fmt.Printf("✓ 正在转到预置点 %d\n", presetId)
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Println("\n========================================")
	fmt.Println("PTZ 控制示例完成!")
	fmt.Println("========================================")
	fmt.Println("\n💡 提示:")
	fmt.Println("  - 云台速度范围: 0-7")
	fmt.Println("  - dwStop=0 开始动作，dwStop=1 停止动作")
	fmt.Println("  - 预置点ID范围: 通常为 1-300")
	fmt.Println("  - 某些命令需要设备硬件支持")
}
