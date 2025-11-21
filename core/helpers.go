package core

/*
#include <stdio.h>
#include <stdlib.h>
#include "hiksdk_wrapper.h"
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"

	"github.com/samsaralc/hiksdk/consts"
	"golang.org/x/text/encoding/simplifiedchinese"
)

// strcpy 将Go字符串安全地复制到C字符数组
// 自动处理字符串长度截断和空终止符，避免缓冲区溢出
// 参数：
//   - dst: 目标C字符数组的指针
//   - src: 源Go字符串
//   - dstLen: 目标数组的长度（包括空终止符）
//
// 注意：
//   - 如果源字符串长度超过 dstLen-1，将被截断
//   - 自动添加空终止符
func strcpy(dst unsafe.Pointer, src string, dstLen int) {
	if dstLen <= 0 {
		return // 无效的目标长度
	}

	srcBytes := []byte(src)
	copyLen := len(srcBytes)

	// 确保不会溢出，留出空间给空终止符
	if copyLen >= dstLen {
		copyLen = dstLen - 1
	}

	// 复制数据
	dstSlice := (*[1 << 30]byte)(dst)
	copy(dstSlice[:copyLen], srcBytes[:copyLen])

	// 添加空终止符
	dstSlice[copyLen] = 0
}

// GBKToUTF8 将GBK编码转换为UTF-8
// 用于处理设备返回的中文字符串
// 参数：
//   - b: GBK编码的字节数组
//
// 返回值：
//   - string: UTF-8编码的字符串
//   - error: 转换错误
func GBKToUTF8(b []byte) (string, error) {
	r, err := simplifiedchinese.GBK.NewDecoder().Bytes(b)
	return string(r), err
}

// UTF8ToGBK 将UTF-8编码转换为GBK
// 用于向设备发送中文字符串
// 参数：
//   - s: UTF-8编码的字符串
//
// 返回值：
//   - []byte: GBK编码的字节数组
//   - error: 转换错误
func UTF8ToGBK(s string) ([]byte, error) {
	return simplifiedchinese.GBK.NewEncoder().Bytes([]byte(s))
}

// LoginWithRetry 带重试的登录
// 在网络不稳定时自动重试登录
// 参数：
//   - maxRetries: 最大重试次数
//   - retryDelay: 重试间隔时间
//
// 返回值：
//   - int: 登录ID
//   - error: 错误信息
func (device *HKDevice) LoginWithRetry(maxRetries int, retryDelay time.Duration) (int, error) {
	var loginId int
	var err error

	for i := 0; i < maxRetries; i++ {
		loginId, err = device.LoginV40()
		if err == nil {
			return loginId, nil
		}

		// 检查是否是连接失败错误（错误码 7）
		if hkErr, ok := err.(*HKError); ok && hkErr.Code == 7 {
			fmt.Printf("⚠ 登录失败 (尝试 %d/%d): %s\n", i+1, maxRetries, hkErr.Msg)
			if i < maxRetries-1 {
				fmt.Printf("  等待 %v 后重试...\n", retryDelay)
				time.Sleep(retryDelay)
				continue
			}
		}

		// 其他错误直接返回
		return -1, err
	}

	return -1, fmt.Errorf("登录失败：已重试 %d 次", maxRetries)
}

// IsConnected 检查设备连接状态
// 通过发送测试请求来检测连接是否正常
// 返回值：
//   - bool: true表示已连接且正常，false表示未连接或连接异常
func (device *HKDevice) IsConnected() bool {
	// 首先检查登录状态
	if device.loginId < 0 {
		return false
	}

	// 通过获取设备时间来检测连接（这是一个轻量级操作）
	var deviceTime C.NET_DVR_TIME
	result := C.NET_DVR_GetDVRConfig(
		C.LONG(device.loginId),
		C.DWORD(consts.NET_DVR_GET_TIMECFG),
		C.LONG(0),
		unsafe.Pointer(&deviceTime),
		C.DWORD(unsafe.Sizeof(deviceTime)),
		nil,
	)

	return result == C.TRUE
}

// GetDeviceIP 获取设备IP地址
func (device *HKDevice) GetDeviceIP() string {
	return device.ip
}

// GetDevicePort 获取设备端口
func (device *HKDevice) GetDevicePort() int {
	return device.port
}

// GetRealHandle 获取实时预览句柄
// 返回值：
//   - int: 预览句柄，-1表示未启动预览
func (device *HKDevice) GetRealHandle() int {
	return device.lRealHandle
}

// GetSDKVersion 获取SDK版本信息
// 返回SDK的版本号和编译时间
// 返回值：
//   - version: 版本号
//   - buildTime: 编译时间
func GetSDKVersion() (version uint32, buildTime uint32) {
	sdkVersion := C.NET_DVR_GetSDKVersion()
	version = uint32(sdkVersion & 0xFFFF0000 >> 16)
	buildTime = uint32(sdkVersion & 0x0000FFFF)
	return
}

// GetSDKBuildVersion 获取SDK详细版本信息
// 返回更详细的SDK版本字符串
// 返回值：
//   - string: 版本信息字符串
func GetSDKBuildVersion() string {
	version, buildTime := GetSDKVersion()
	majorVersion := (version >> 8) & 0xFF
	minorVersion := version & 0xFF
	return fmt.Sprintf("V%d.%d Build %d", majorVersion, minorVersion, buildTime)
}

// SetLogConfig 配置SDK日志
// 设置SDK的日志级别和输出文件
// 参数：
//   - level: 日志级别（0=关闭，1=错误，2=警告，3=信息，4=调试）
//   - logDir: 日志目录
//   - autoDelete: 是否自动删除旧日志
//
// 返回值：
//   - error: 错误信息
func SetLogConfig(level int, logDir string, autoDelete bool) error {
	cLogDir := C.CString(logDir)
	defer C.free(unsafe.Pointer(cLogDir))

	deleteFlag := C.BOOL(0)
	if autoDelete {
		deleteFlag = C.BOOL(1)
	}

	result := C.NET_DVR_SetLogToFile(
		C.DWORD(level),
		cLogDir,
		deleteFlag,
	)

	if result != C.TRUE {
		return fmt.Errorf("设置日志配置失败")
	}

	return nil
}
