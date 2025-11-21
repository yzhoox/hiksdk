package core

/*
#cgo CFLAGS: -I../include

// Linux 平台的链接配置
// 链接顺序：被依赖的库放在后面，依赖别人的库放在前面
#cgo linux LDFLAGS: -L../lib/Linux -lhcnetsdk -lhpr -lHCCore

// Windows 平台的链接配置
// 链接顺序：被依赖的库放在后面，依赖别人的库放在前面
#cgo windows LDFLAGS: -L../lib/Windows -lHCNetSDK -lHCCore

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include "hiksdk_wrapper.h"
*/
import "C"
import (
	"log"
	"runtime/cgo"
	"sync"
)

var (
	// sdkMutex 保护SDK初始化和清理操作的互斥锁
	sdkMutex sync.Mutex
	// sdkInitialized 标记SDK是否已初始化
	// true 表示已经调用过 NET_DVR_Init 并成功
	// false 表示未初始化或已调用 NET_DVR_Cleanup
	sdkInitialized bool
)

// HKDevice 海康设备结构体
// 封装了海康威视设备的所有操作
type HKDevice struct {
	// 连接信息
	ip       string // 设备IP地址
	port     int    // 设备端口（默认8000）
	username string // 登录用户名
	password string // 登录密码

	// 状态信息
	loginId        int        // 登录ID，-1表示未登录
	alarmHandle    int        // 报警句柄
	lRealHandle    int        // 实时预览句柄
	byChanNum      int        // 通道数量
	receiverHandle cgo.Handle // 数据接收器句柄
}

// DeviceInfo 设备信息结构体
// 包含设备的基本信息和连接参数
type DeviceInfo struct {
	IP         string // 设备IP地址
	Port       int    // 端口号（默认8000）
	UserName   string // 登录用户名
	Password   string // 登录密码
	DeviceID   string // 设备序列号（设备唯一标识）
	DeviceName string // 设备名称（用户自定义）
	ByChanNum  int    // 通道数量
}

// initSDK 初始化SDK
// 使用互斥锁确保线程安全
// 特性：
//   - 第一次调用时执行 NET_DVR_Init
//   - 调用 Cleanup 后再次调用，会重新初始化
//   - 多次调用 NewHKDevice 只会初始化一次
func initSDK() {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()

	// 如果已经初始化，直接返回
	if sdkInitialized {
		return
	}

	// 执行初始化
	result := C.NET_DVR_Init()
	if result != C.TRUE {
		log.Println("✗ SDK初始化失败")
		return
	}

	// 设置连接超时参数
	C.NET_DVR_SetConnectTime(2000, 5) // 连接超时2秒，重试5次
	// 设置重连参数
	C.NET_DVR_SetReconnect(10000, 1) // 重连间隔10秒，启用重连

	sdkInitialized = true
	log.Println("✓ 海康SDK初始化成功")
}

// NewHKDevice 创建海康设备实例
// SDK会在第一次创建设备时自动初始化
// 参数：
//   - info: 设备信息，包含IP、端口、用户名、密码等
//
// 返回：
//   - *HKDevice: 设备实例指针
func NewHKDevice(info DeviceInfo) *HKDevice {
	// 确保SDK已初始化
	initSDK()

	return &HKDevice{
		ip:       info.IP,
		port:     info.Port,
		username: info.UserName,
		password: info.Password,
		loginId:  -1, // 初始化为未登录状态
	}
}

// Cleanup 清理SDK资源
// 注意：
//   - 调用此函数后，所有通过 SDK 建立的连接都将失效
//   - 通常在程序退出时调用（例如：defer core.Cleanup()）
//   - 调用 Cleanup 之后，可以再次调用 NewHKDevice，SDK 会自动重新初始化
func Cleanup() {
	sdkMutex.Lock()
	defer sdkMutex.Unlock()

	if !sdkInitialized {
		return // 已经清理或未初始化
	}

	// 执行清理
	result := C.NET_DVR_Cleanup()
	if result != C.TRUE {
		log.Println("✗ SDK清理失败")
		return
	}

	sdkInitialized = false
	log.Println("✓ 海康SDK已清理（可重新初始化）")
}

// SetSDKLog 配置SDK日志
// 参数：
//   - level: 日志级别（0-关闭, 1-错误, 2-警告, 3-信息, 4-调试）
//   - logDir: 日志存储目录（空字符串表示默认）
//   - autoDelete: 是否自动删除超量日志
func SetSDKLog(level int, logDir string, autoDelete bool) error {
	// 确保SDK已初始化
	initSDK()
	return SetLogConfig(level, logDir, autoDelete)
}

// IsLoggedIn 检查设备是否已登录
func (device *HKDevice) IsLoggedIn() bool {
	return device.loginId >= 0
}

// GetLoginID 获取登录ID
func (device *HKDevice) GetLoginID() int {
	return device.loginId
}

// GetChannelCount 获取通道数量
func (device *HKDevice) GetChannelCount() int {
	return device.byChanNum
}
