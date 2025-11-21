package pkg

import (
	"os"
	"testing"
	"time"
)

// 测试设备登录功能
func TestDeviceLogin(t *testing.T) {
	// 跳过测试如果没有设置环境变量
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
	loginId, err := dev.Login()

	if err != nil {
		t.Fatalf("登录失败: %v", err)
	}

	if loginId <= 0 {
		t.Fatalf("登录ID无效: %d", loginId)
	}

	t.Logf("✓ 登录成功，登录ID: %d", loginId)

	// 测试登出
	err = dev.Logout()
	if err != nil {
		t.Errorf("登出失败: %v", err)
	}
}

// 测试获取设备信息
func TestGetDeviceInfo(t *testing.T) {
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
	if err != nil {
		t.Fatalf("获取设备信息失败: %v", err)
	}

	t.Logf("设备名称: %s", info.DeviceName)
	t.Logf("设备序列号: %s", info.DeviceID)
	t.Logf("通道数量: %d", info.ByChanNum)

	if info.ByChanNum == 0 {
		t.Error("设备通道数为0")
	}
}

// 测试获取通道名称
func TestGetChannelName(t *testing.T) {
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

	channels, err := dev.GetChannelName()
	if err != nil {
		t.Fatalf("获取通道名称失败: %v", err)
	}

	t.Logf("获取到 %d 个通道", len(channels))
	for id, name := range channels {
		t.Logf("通道 %d: %s", id, name)
	}
}

// 测试重复登录
func TestRepeatedLogin(t *testing.T) {
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

	// 第一次登录
	_, err := dev.Login()
	if err != nil {
		t.Fatalf("第一次登录失败: %v", err)
	}

	dev.Logout()
	time.Sleep(500 * time.Millisecond)

	// 第二次登录
	_, err = dev.LoginV4()
	if err != nil {
		t.Fatalf("第二次登录失败: %v", err)
	}
	defer dev.Logout()

	t.Log("✓ 重复登录测试成功")
}

// 基准测试：登录性能
func BenchmarkDeviceLogin(b *testing.B) {
	ip := os.Getenv("HIK_IP")
	if ip == "" {
		b.Skip("跳过测试: 请设置环境变量 HIK_IP, HIK_PORT, HIK_USER, HIK_PASSWORD")
	}

	deviceInfo := DeviceInfo{
		IP:       ip,
		Port:     8000,
		UserName: os.Getenv("HIK_USER"),
		Password: os.Getenv("HIK_PASSWORD"),
	}

	InitHikSDK()
	defer HKExit()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dev := NewHKDevice(deviceInfo)
		_, err := dev.Login()
		if err != nil {
			b.Fatalf("登录失败: %v", err)
		}
		dev.Logout()
		time.Sleep(100 * time.Millisecond) // 避免频繁连接
	}
}
