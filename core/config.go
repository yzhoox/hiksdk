package core

/*
#include <stdio.h>
#include <stdlib.h>
#include "hiksdk_wrapper.h"
*/
import "C"
import (
	"fmt"
	"log"
	"strings"
	"unsafe"

	"github.com/samsaralc/hiksdk/consts"
)

// GetDeviceInfo 获取设备详细信息
// 使用 NET_DVR_GetDVRConfig 接口获取设备的配置信息
// 返回值：
//   - *DeviceInfo: 设备信息结构体，包含设备名称、序列号、通道数等
//   - error: 错误信息，成功时为nil
func (device *HKDevice) GetDeviceInfo() (*DeviceInfo, error) {
	// 参数验证
	if device.loginId < 0 {
		return nil, fmt.Errorf("设备未登录")
	}

	var deviceInfo C.NET_DVR_DEVICECFG
	var bytesReturned C.DWORD
	deviceInfo.dwSize = C.DWORD(unsafe.Sizeof(deviceInfo))

	// 调用SDK获取设备配置
	result := C.NET_DVR_GetDVRConfig(
		C.LONG(device.loginId),
		C.DWORD(consts.NET_DVR_GET_DEVICECFG),
		C.LONG(0),
		unsafe.Pointer(&deviceInfo),
		C.DWORD(unsafe.Sizeof(deviceInfo)),
		&bytesReturned,
	)

	if result != C.TRUE {
		return nil, device.HKErr("获取设备信息失败")
	}

	// 转换设备名称
	dvrNameBytes := make([]byte, len(deviceInfo.sDVRName))
	for i := range deviceInfo.sDVRName {
		dvrNameBytes[i] = byte(deviceInfo.sDVRName[i])
	}
	sDVRName := string(dvrNameBytes)

	// 转换序列号
	serialNumberBytes := make([]byte, len(deviceInfo.sSerialNumber))
	for i := range deviceInfo.sSerialNumber {
		serialNumberBytes[i] = byte(deviceInfo.sSerialNumber[i])
	}
	sSerialNumber := string(serialNumberBytes)

	// 清理字符串
	sDVRName = strings.TrimRight(sDVRName, "\x00")
	sDVRName = strings.TrimSpace(sDVRName)
	sSerialNumber = strings.TrimRight(sSerialNumber, "\x00")
	sSerialNumber = strings.TrimSpace(sSerialNumber)

	// GBK转UTF8
	sDVRName, _ = GBKToUTF8([]byte(sDVRName))

	// 更新设备通道数
	device.byChanNum = int(deviceInfo.byChanNum)

	info := &DeviceInfo{
		IP:         device.ip,
		Port:       device.port,
		UserName:   device.username,
		Password:   device.password,
		DeviceID:   sSerialNumber,
		DeviceName: sDVRName,
		ByChanNum:  int(deviceInfo.byChanNum),
	}

	log.Printf("✓ 获取设备信息成功 - 名称: %s, 序列号: %s, 通道数: %d", sDVRName, sSerialNumber, info.ByChanNum)
	return info, nil
}

// GetChannelName 获取所有通道的名称
// 遍历设备的所有通道，获取每个通道的名称
// 返回值：
//   - map[int]string: 通道ID到通道名称的映射
//   - error: 错误信息，成功时为nil
func (device *HKDevice) GetChannelName() (map[int]string, error) {
	// 参数验证
	if device.loginId < 0 {
		return nil, fmt.Errorf("设备未登录")
	}
	if device.byChanNum == 0 {
		log.Println("警告: 通道数为0，尝试获取设备信息")
		// 尝试获取设备信息以更新通道数
		if _, err := device.GetDeviceInfo(); err != nil {
			return nil, err
		}
	}

	channelNames := make(map[int]string, device.byChanNum)

	successCount := 0
	for i := 1; i <= device.byChanNum; i++ {
		var channelInfo C.NET_DVR_PICCFG
		var bytesReturned C.DWORD
		channelInfo.dwSize = C.DWORD(unsafe.Sizeof(channelInfo))
		var sDVRName string

		// 获取通道配置
		result := C.NET_DVR_GetDVRConfig(
			C.LONG(device.loginId),
			C.DWORD(consts.NET_DVR_GET_PICCFG),
			C.LONG(i),
			unsafe.Pointer(&channelInfo),
			C.DWORD(unsafe.Sizeof(channelInfo)),
			&bytesReturned,
		)

		if result != C.TRUE {
			// 获取失败时使用默认名称
			sDVRName = fmt.Sprintf("Camera_%d", i)
			log.Printf("警告: 获取通道%d名称失败，使用默认名称", i)
		} else {
			// 转换通道名称
			chanNameBytes := make([]byte, len(channelInfo.sChanName))
			for j := range channelInfo.sChanName {
				chanNameBytes[j] = byte(channelInfo.sChanName[j])
			}
			sDVRName = string(chanNameBytes)

			// 清理字符串
			sDVRName = strings.TrimRight(sDVRName, "\x00")
			sDVRName = strings.TrimSpace(sDVRName)

			// GBK转UTF8
			if len(sDVRName) > 0 {
				if utf8Name, err := GBKToUTF8([]byte(sDVRName)); err == nil {
					sDVRName = utf8Name
				}
			}

			// 如果转换后为空，使用默认名称
			if sDVRName == "" {
				sDVRName = fmt.Sprintf("Camera_%d", i)
			}

			successCount++
		}

		channelNames[i] = sDVRName
	}

	log.Printf("✓ 获取通道名称成功 - 成功: %d/%d", successCount, device.byChanNum)
	return channelNames, nil
}

// GetChannelPTZ 获取指定通道的PTZ位置信息
// 获取云台当前的水平角度、垂直角度和变焦倍数
// 参数：
//   - channel: 通道号，从1开始
func (device *HKDevice) GetChannelPTZ(channel int) {
	var ptzPos C.NET_DVR_PTZPOS
	var ptzScope C.NET_DVR_PTZSCOPE
	var bytesReturned C.DWORD

	// 获取PTZ位置
	if C.NET_DVR_GetDVRConfig(
		C.LONG(device.loginId),
		C.DWORD(consts.NET_DVR_GET_PTZPOS),
		C.LONG(channel),
		unsafe.Pointer(&ptzPos),
		C.DWORD(unsafe.Sizeof(ptzPos)),
		&bytesReturned,
	) != C.TRUE {
		// 获取失败，但不返回错误
	}

	// 获取PTZ范围
	if C.NET_DVR_GetDVRConfig(
		C.LONG(device.loginId),
		C.DWORD(consts.NET_DVR_GET_PTZSCOPE),
		C.LONG(channel),
		unsafe.Pointer(&ptzScope),
		C.DWORD(unsafe.Sizeof(ptzScope)),
		&bytesReturned,
	) != C.TRUE {
		// 获取失败，但不返回错误
	}

	// 可以在这里处理或返回PTZ信息
}

// SetupNewChannels 批量配置新通道（内部函数）
// 用于初始化设备的多个通道配置
func (device *HKDevice) SetupNewChannels(startChannel, endChannel int, configs []ChannelConfig) error {
	for i := startChannel; i <= endChannel && i-startChannel < len(configs); i++ {
		config := configs[i-startChannel]

		// 设置通道名称
		var piccfg C.NET_DVR_PICCFG
		piccfg.dwSize = C.DWORD(unsafe.Sizeof(piccfg))

		// 转换通道名称为GBK
		nameBytes, _ := UTF8ToGBK(config.Name)
		for j := 0; j < len(nameBytes) && j < len(piccfg.sChanName); j++ {
			piccfg.sChanName[j] = C.BYTE(nameBytes[j])
		}

		// 设置其他参数
		piccfg.dwShowChanName = C.DWORD(1) // 显示通道名称

		// 应用配置
		result := C.NET_DVR_SetDVRConfig(
			C.LONG(device.loginId),
			C.DWORD(consts.NET_DVR_SET_PICCFG),
			C.LONG(i),
			unsafe.Pointer(&piccfg),
			C.DWORD(unsafe.Sizeof(piccfg)),
		)

		if result != C.TRUE {
			return device.HKErr(fmt.Sprintf("设置通道%d配置失败", i))
		}
	}

	return nil
}

// ChannelConfig 通道配置结构体
type ChannelConfig struct {
	Name       string // 通道名称
	ShowName   bool   // 是否显示名称
	Brightness uint8  // 亮度
	Contrast   uint8  // 对比度
	Saturation uint8  // 饱和度
	Hue        uint8  // 色调
}
