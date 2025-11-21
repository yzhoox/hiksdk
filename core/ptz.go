package core

/*
#include <stdio.h>
#include <stdlib.h>
#include "hiksdk_wrapper.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"log"
)

// PTZControlWithSpeed PTZ云台控制（带速度参数）
// 使用实时播放句柄控制，需要先调用 RealPlay_V40 启动预览
// 参数：
//   - dwPTZCommand: PTZ命令，如 PAN_RIGHT(24)、TILT_UP(21) 等
//   - dwStop: 停止标志，0=开始动作，1=停止动作
//   - dwSpeed: 速度值，范围0-7，数值越大速度越快
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZControlWithSpeed(dwPTZCommand, dwStop, dwSpeed int) (bool, error) {
	if device.lRealHandle == 0 {
		return false, errors.New("PTZ控制失败: 需要先调用 RealPlay_V40 启动预览")
	}
	
	if C.NET_DVR_PTZControlWithSpeed(
		C.LONG(device.lRealHandle),
		C.DWORD(dwPTZCommand),
		C.DWORD(dwStop),
		C.DWORD(dwSpeed),
	) != C.TRUE {
		return false, device.HKErr("PTZ控制失败")
	}
	
	log.Println("PTZ控制成功")
	return true, nil
}

// PTZControlWithSpeed_Other PTZ云台控制（指定通道，带速度参数）
// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
// 参数：
//   - lChannel: 通道号，从1开始
//   - dwPTZCommand: PTZ命令，如 PAN_RIGHT(24)、TILT_UP(21) 等
//   - dwStop: 停止标志，0=开始动作，1=停止动作
//   - dwSpeed: 速度值，范围0-7，数值越大速度越快
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZControlWithSpeed_Other(lChannel, dwPTZCommand, dwStop, dwSpeed int) (bool, error) {
	if C.NET_DVR_PTZControlWithSpeed_Other(
		C.LONG(device.loginId),
		C.LONG(lChannel),
		C.DWORD(dwPTZCommand),
		C.DWORD(dwStop),
		C.DWORD(dwSpeed),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("PTZ控制通道 %d 失败", lChannel))
	}
	
	log.Printf("PTZ控制通道 %d 成功", lChannel)
	return true, nil
}

// PTZControl PTZ云台控制（无速度参数）
// 使用实时播放句柄控制，需要先调用 RealPlay_V40 启动预览
// 参数：
//   - dwPTZCommand: PTZ命令，如 ZOOM_IN(11)、FOCUS_NEAR(13) 等
//   - dwStop: 停止标志，0=开始动作，1=停止动作
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZControl(dwPTZCommand, dwStop int) (bool, error) {
	if device.lRealHandle == 0 {
		return false, errors.New("PTZ控制失败: 需要先调用 RealPlay_V40 启动预览")
	}
	
	if C.NET_DVR_PTZControl(
		C.LONG(device.lRealHandle),
		C.DWORD(dwPTZCommand),
		C.DWORD(dwStop),
	) != C.TRUE {
		return false, device.HKErr("PTZ控制失败")
	}
	
	log.Println("PTZ控制成功")
	return true, nil
}

// PTZControl_Other PTZ云台控制（指定通道，无速度参数）
// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
// 参数：
//   - lChannel: 通道号，从1开始
//   - dwPTZCommand: PTZ命令，如 ZOOM_IN(11)、FOCUS_NEAR(13) 等
//   - dwStop: 停止标志，0=开始动作，1=停止动作
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZControl_Other(lChannel, dwPTZCommand, dwStop int) (bool, error) {
	if C.NET_DVR_PTZControl_Other(
		C.LONG(device.loginId),
		C.LONG(lChannel),
		C.DWORD(dwPTZCommand),
		C.DWORD(dwStop),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("PTZ控制通道 %d 失败", lChannel))
	}
	
	log.Printf("PTZ控制通道 %d 成功", lChannel)
	return true, nil
}

// PTZPreset 预置点操作（使用实时预览句柄）
// 需要先调用 RealPlay_V40 启动预览
// 参数：
//   - dwPTZCommand: 预置点命令，如 SET_PRESET(8)、GOTO_PRESET(39)、CLE_PRESET(9)
//   - presetIndex: 预置点编号（1-255）
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZPreset(dwPTZCommand, presetIndex int) (bool, error) {
	if device.lRealHandle == 0 {
		return false, errors.New("预置点操作失败: 需要先调用 RealPlay_V40 启动预览")
	}
	
	if C.NET_DVR_PTZPreset(
		C.LONG(device.lRealHandle),
		C.DWORD(dwPTZCommand),
		C.DWORD(presetIndex),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("预置点 %d 操作失败", presetIndex))
	}
	
	log.Printf("预置点 %d 操作成功", presetIndex)
	return true, nil
}

// PTZPreset_Other 预置点操作（指定通道）
// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
// 参数：
//   - lChannel: 通道号，从1开始
//   - dwPTZCommand: 预置点命令，如 SET_PRESET(8)、GOTO_PRESET(39)、CLE_PRESET(9)
//   - presetIndex: 预置点编号（1-255）
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZPreset_Other(lChannel, dwPTZCommand, presetIndex int) (bool, error) {
	if C.NET_DVR_PTZPreset_Other(
		C.LONG(device.loginId),
		C.LONG(lChannel),
		C.DWORD(dwPTZCommand),
		C.DWORD(presetIndex),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("通道 %d 预置点 %d 操作失败", lChannel, presetIndex))
	}
	
	log.Printf("通道 %d 预置点 %d 操作成功", lChannel, presetIndex)
	return true, nil
}

// PTZCruise 巡航操作
// 控制云台按预定路径自动巡航
// 参数：
//   - lChannel: 通道号，从1开始
//   - cruiseRoute: 巡航路线号
//   - cruisePoint: 巡航点号
//   - value: 预置点值或时间值
//
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZCruise(lChannel, cruiseRoute, cruisePoint, value int) error {
	// 可以根据需要实现巡航功能
	// 使用 NET_DVR_PTZCruise 接口
	return nil
}

// PTZTrack 轨迹操作
// 控制云台按记录的轨迹运动
// 参数：
//   - lChannel: 通道号，从1开始
//   - command: 轨迹命令（开始记录、停止记录、运行轨迹）
//
// 返回值：
//   - error: 错误信息，成功时为nil  
func (device *HKDevice) PTZTrack(lChannel, command int) error {
	// 可以根据需要实现轨迹功能
	// 使用 NET_DVR_PTZTrack_Other 接口
	return nil
}
