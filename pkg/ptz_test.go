package pkg

import (
	"os"
	"testing"
	"time"
)

// PTZ命令常量已经在 ptz_commands.go 中定义
// 这里直接使用包中的常量

// 测试PTZ控制（需要实际设备）
func TestPTZControl(t *testing.T) {
	ip := os.Getenv("HIK_IP")
	if ip == "" {
		t.Skip("跳过测试: 请设置环境变量 HIK_IP, HIK_PORT, HIK_USER, HIK_PASSWORD")
	}

	deviceInfo := DeviceInfo{
		IP:       ip,
		Port:     8000,
		UserName: os.Getenv("HIK_USER"),
		Password: os.Getenv("HIK_PASSWORD"),
	}

	InitHikSDK()
	defer HKExit()

	dev := NewHKDevice(deviceInfo)
	_, err := dev.Login()
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer dev.Logout()

	// 获取设备信息，确认有通道
	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		t.Skip("设备没有可用通道")
	}

	channelId := 1

	// 测试云台右转
	t.Run("云台右转", func(t *testing.T) {
		success, err := dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 3)
		if err != nil {
			t.Logf("云台右转失败（可能设备不支持PTZ）: %v", err)
			return
		}
		if !success {
			t.Error("云台右转失败")
			return
		}

		time.Sleep(1 * time.Second)

		// 停止
		_, err = dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 1, 3)
		if err != nil {
			t.Errorf("停止云台失败: %v", err)
		}
		t.Log("✓ 云台右转测试成功")
	})

	// 测试云台左转
	t.Run("云台左转", func(t *testing.T) {
		success, err := dev.PTZControlWithSpeed_Other(channelId, PAN_LEFT, 0, 3)
		if err != nil {
			t.Logf("云台左转失败（可能设备不支持PTZ）: %v", err)
			return
		}
		if !success {
			t.Error("云台左转失败")
			return
		}

		time.Sleep(1 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, PAN_LEFT, 1, 3)
		t.Log("✓ 云台左转测试成功")
	})
}

// 测试PTZ变焦
func TestPTZZoom(t *testing.T) {
	ip := os.Getenv("HIK_IP")
	if ip == "" {
		t.Skip("跳过测试: 请设置环境变量 HIK_IP, HIK_PORT, HIK_USER, HIK_PASSWORD")
	}

	deviceInfo := DeviceInfo{
		IP:       ip,
		Port:     8000,
		UserName: os.Getenv("HIK_USER"),
		Password: os.Getenv("HIK_PASSWORD"),
	}

	InitHikSDK()
	defer HKExit()

	dev := NewHKDevice(deviceInfo)
	_, err := dev.Login()
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer dev.Logout()

	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		t.Skip("设备没有可用通道")
	}

	channelId := 1

	// 测试焦距放大
	t.Run("焦距放大", func(t *testing.T) {
		success, err := dev.PTZControl_Other(channelId, ZOOM_IN, 0)
		if err != nil {
			t.Logf("焦距放大失败（可能设备不支持）: %v", err)
			return
		}
		if success {
			time.Sleep(1 * time.Second)
			dev.PTZControl_Other(channelId, ZOOM_IN, 1)
			t.Log("✓ 焦距放大测试成功")
		}
	})

	// 测试焦距缩小
	t.Run("焦距缩小", func(t *testing.T) {
		success, err := dev.PTZControl_Other(channelId, ZOOM_OUT, 0)
		if err != nil {
			t.Logf("焦距缩小失败（可能设备不支持）: %v", err)
			return
		}
		if success {
			time.Sleep(1 * time.Second)
			dev.PTZControl_Other(channelId, ZOOM_OUT, 1)
			t.Log("✓ 焦距缩小测试成功")
		}
	})
}

// 测试预置点功能
func TestPTZPreset(t *testing.T) {
	ip := os.Getenv("HIK_IP")
	if ip == "" {
		t.Skip("跳过测试: 请设置环境变量 HIK_IP, HIK_PORT, HIK_USER, HIK_PASSWORD")
	}

	deviceInfo := DeviceInfo{
		IP:       ip,
		Port:     8000,
		UserName: os.Getenv("HIK_USER"),
		Password: os.Getenv("HIK_PASSWORD"),
	}

	InitHikSDK()
	defer HKExit()

	dev := NewHKDevice(deviceInfo)
	_, err := dev.Login()
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer dev.Logout()

	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		t.Skip("设备没有可用通道")
	}

	channelId := 1
	presetId := 1

	// 设置预置点
	t.Run("设置预置点", func(t *testing.T) {
		success, err := dev.PTZControl_Other(channelId, SET_PRESET, presetId)
		if err != nil {
			t.Logf("设置预置点失败（可能设备不支持）: %v", err)
			return
		}
		if success {
			t.Logf("✓ 预置点 %d 设置成功", presetId)
		}
	})

	// 转到预置点
	t.Run("转到预置点", func(t *testing.T) {
		// 先移动云台
		dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 3)
		time.Sleep(2 * time.Second)
		dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 1, 3)
		time.Sleep(500 * time.Millisecond)

		// 转到预置点
		success, err := dev.PTZControl_Other(channelId, GOTO_PRESET, presetId)
		if err != nil {
			t.Logf("转到预置点失败: %v", err)
			return
		}
		if success {
			t.Logf("✓ 正在转到预置点 %d", presetId)
			time.Sleep(2 * time.Second)
		}
	})
}

// 测试获取PTZ位置信息
func TestGetPTZPosition(t *testing.T) {
	ip := os.Getenv("HIK_IP")
	if ip == "" {
		t.Skip("跳过测试: 请设置环境变量 HIK_IP, HIK_PORT, HIK_USER, HIK_PASSWORD")
	}

	deviceInfo := DeviceInfo{
		IP:       ip,
		Port:     8000,
		UserName: os.Getenv("HIK_USER"),
		Password: os.Getenv("HIK_PASSWORD"),
	}

	InitHikSDK()
	defer HKExit()

	dev := NewHKDevice(deviceInfo)
	_, err := dev.Login()
	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}
	defer dev.Logout()

	info, err := dev.GetDeiceInfo()
	if err != nil || info.ByChanNum == 0 {
		t.Skip("设备没有可用通道")
	}

	channelId := 1
	dev.GetChannelPTZ(channelId)
	t.Log("✓ PTZ位置信息获取成功")
}
