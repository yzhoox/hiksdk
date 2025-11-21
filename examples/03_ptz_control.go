package main

import (
	"fmt"
	"time"

	"github.com/samsaralc/hiksdk/core"
)

// PTZ云台控制示例
func main() {
	fmt.Println("========================================")
	fmt.Println("海康威视 SDK - PTZ控制示例")
	fmt.Println("========================================")

	// 设备连接信息
	deviceInfo := core.DeviceInfo{
		IP:       "192.168.1.64",
		Port:     8000,
		UserName: "admin",
		Password: "asdf234.",
	}

	// 创建设备并登录
	dev := core.NewHKDevice(deviceInfo)
	loginId, err := dev.LoginV30()
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		return
	}
	fmt.Printf("登录成功 (ID: %d)\n", loginId)
	defer dev.Logout()

	// 选择通道
	channel := 1

	// ==================== 方向控制 ====================
	fmt.Println("\n[1] 方向控制示例:")

	// 向右转动
	fmt.Println("  向右转动...")
	dev.PTZControlWithSpeed_Other(channel, core.PAN_RIGHT, 0, 7) // 速度4
	time.Sleep(2 * time.Second)
	dev.PTZControlWithSpeed_Other(channel, core.PAN_RIGHT, 1, 7) // 停止

	// 向上转动
	fmt.Println("  向上转动...")
	dev.PTZControlWithSpeed_Other(channel, core.TILT_UP, 0, 7)
	time.Sleep(2 * time.Second)
	dev.PTZControlWithSpeed_Other(channel, core.TILT_UP, 1, 7)

	// ==================== 变焦控制 ====================
	fmt.Println("\n[2] 变焦控制示例:")

	// 放大
	fmt.Println("  放大...")
	dev.PTZControl_Other(channel, core.ZOOM_IN, 0)
	time.Sleep(2 * time.Second)
	dev.PTZControl_Other(channel, core.ZOOM_IN, 1)

	// 缩小
	fmt.Println("  缩小...")
	dev.PTZControl_Other(channel, core.ZOOM_OUT, 0)
	time.Sleep(2 * time.Second)
	dev.PTZControl_Other(channel, core.ZOOM_OUT, 1)

	// ==================== 预置点操作 ====================
	fmt.Println("\n[3] 预置点操作:")

	// 设置预置点1
	fmt.Println("  设置预置点1...")
	dev.PTZPreset_Other(channel, core.SET_PRESET, 1)

	// 移动到其他位置
	fmt.Println("  移动到其他位置...")
	dev.PTZControlWithSpeed_Other(channel, core.PAN_LEFT, 0, 4)
	time.Sleep(3 * time.Second)
	dev.PTZControlWithSpeed_Other(channel, core.PAN_LEFT, 1, 4)

	// 调用预置点1
	fmt.Println("  调用预置点1...")
	dev.PTZPreset_Other(channel, core.GOTO_PRESET, 1)
	time.Sleep(3 * time.Second)

	fmt.Println("\n示例完成!")
}
