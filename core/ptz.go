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

// PTZCruise 巡航操作（使用实时预览句柄）
// 需要先调用 RealPlay_V40 启动预览
// 参数：
//   - dwPTZCruiseCmd: 巡航命令，如 RUN_SEQ(37)、STOP_SEQ(38) 等
//   - byCruiseRoute: 巡航路径，最多支持 32 条路径（序号从 1 开始）
//   - byCruisePoint: 巡航点，最多支持 32 个点（序号从 1 开始）
//   - wInput: 输入参数，不同命令对应不同值（预置点/时间/速度）
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZCruise(dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput int) (bool, error) {
	if device.lRealHandle == 0 {
		return false, errors.New("巡航操作失败: 需要先调用 RealPlay_V40 启动预览")
	}

	if C.NET_DVR_PTZCruise(
		C.LONG(device.lRealHandle),
		C.DWORD(dwPTZCruiseCmd),
		C.BYTE(byCruiseRoute),
		C.BYTE(byCruisePoint),
		C.WORD(wInput),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("巡航操作失败 (Cmd:%d Route:%d Point:%d Input:%d)",
			dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput))
	}

	log.Printf("巡航操作成功 (Cmd:%d Route:%d Point:%d Input:%d)",
		dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput)
	return true, nil
}

// PTZCruise_Other 巡航操作（指定通道）
// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
// 参数：
//   - lChannel: 通道号，从1开始
//   - dwPTZCruiseCmd: 巡航命令，如 RUN_SEQ(37)、STOP_SEQ(38) 等
//   - byCruiseRoute: 巡航路径，最多支持 32 条路径（序号从 1 开始）
//   - byCruisePoint: 巡航点，最多支持 32 个点（序号从 1 开始）
//   - wInput: 输入参数，不同命令对应不同值（预置点/时间/速度）
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZCruise_Other(lChannel, dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput int) (bool, error) {
	if C.NET_DVR_PTZCruise_Other(
		C.LONG(device.loginId),
		C.LONG(lChannel),
		C.DWORD(dwPTZCruiseCmd),
		C.BYTE(byCruiseRoute),
		C.BYTE(byCruisePoint),
		C.WORD(wInput),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("通道 %d 巡航操作失败 (Cmd:%d Route:%d Point:%d Input:%d)",
			lChannel, dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput))
	}

	log.Printf("通道 %d 巡航操作成功 (Cmd:%d Route:%d Point:%d Input:%d)",
		lChannel, dwPTZCruiseCmd, byCruiseRoute, byCruisePoint, wInput)
	return true, nil
}

// PTZTrack 轨迹操作（使用实时预览句柄）
// 需要先调用 RealPlay_V40 启动预览
// 参数：
//   - dwPTZTrackCmd: 轨迹命令，如 STA_MEM_CRUISE(34)、RUN_CRUISE(36) 等
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZTrack(dwPTZTrackCmd int) (bool, error) {
	if device.lRealHandle == 0 {
		return false, errors.New("轨迹操作失败: 需要先调用 RealPlay_V40 启动预览")
	}

	if C.NET_DVR_PTZTrack(
		C.LONG(device.lRealHandle),
		C.DWORD(dwPTZTrackCmd),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("轨迹操作失败 (Cmd:%d)", dwPTZTrackCmd))
	}

	log.Printf("轨迹操作成功 (Cmd:%d)", dwPTZTrackCmd)
	return true, nil
}

// PTZTrack_Other 轨迹操作（指定通道）
// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
// 参数：
//   - lChannel: 通道号，从1开始
//   - dwPTZTrackCmd: 轨迹命令，如 STA_MEM_CRUISE(34)、RUN_CRUISE(36) 等
//
// 返回值：
//   - bool: 操作是否成功
//   - error: 错误信息，成功时为nil
func (device *HKDevice) PTZTrack_Other(lChannel, dwPTZTrackCmd int) (bool, error) {
	if C.NET_DVR_PTZTrack_Other(
		C.LONG(device.loginId),
		C.LONG(lChannel),
		C.DWORD(dwPTZTrackCmd),
	) != C.TRUE {
		return false, device.HKErr(fmt.Sprintf("通道 %d 轨迹操作失败 (Cmd:%d)", lChannel, dwPTZTrackCmd))
	}

	log.Printf("通道 %d 轨迹操作成功 (Cmd:%d)", lChannel, dwPTZTrackCmd)
	return true, nil
}
