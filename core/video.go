package core

/*
#include <stdio.h>
#include <stdlib.h>
#include "hiksdk_wrapper.h"

// 声明Go回调函数，供C代码调用
extern void GoRealDataCallback(LONG lRealHandle, DWORD dwDataType, BYTE *pBuffer, DWORD dwBufSize, uintptr_t handle);
*/
import "C"
import (
	"fmt"
	"log"
	"runtime/cgo"
	"unsafe"
)

// GoRealDataCallback 实时数据回调函数
// 由C代码调用，接收视频流数据
//
//export GoRealDataCallback
func GoRealDataCallback(lRealHandle C.LONG, dwDataType C.DWORD, pBuffer *C.BYTE, dwBufSize C.DWORD, handle C.uintptr_t) {
	// 从 cgo.Handle 中还原 Go 指针
	h := cgo.Handle(handle)
	receiver, ok := h.Value().(*Receiver)
	if !ok {
		fmt.Println("Error: invalid receiver handle")
		return
	}
	if receiver == nil {
		fmt.Println("Error: receiver is nil in callback")
		return
	}

	size := int(dwBufSize)
	if size > 0 && pBuffer != nil {
		// 将C指针转换为Go的byte切片
		buffer := (*[1 << 30]C.BYTE)(unsafe.Pointer(pBuffer))[:size:size]
		// 高效转换[]C.BYTE为[]byte
		goBuffer := (*[1 << 30]byte)(unsafe.Pointer(&buffer[0]))[:len(buffer):len(buffer)]
		receiver.ReadPSData(goBuffer)
	}
}

// RealPlay_V40 启动实时视频预览
// 使用 NET_DVR_RealPlay_V40 接口启动指定通道的实时视频流
// 参数：
//   - channelId: 通道号，从1开始
//   - receiver: 数据接收器，用于接收视频流数据
//
// 返回值：
//   - int: 预览句柄，大于0表示成功，-1表示失败
//   - error: 错误信息，成功时为nil
func (device *HKDevice) RealPlay_V40(channelId int, receiver *Receiver) (int, error) {
	// 参数验证
	if device.loginId < 0 {
		return -1, fmt.Errorf("设备未登录，无法启动预览")
	}
	if channelId < 1 || channelId > 256 {
		return -1, fmt.Errorf("通道号无效: %d（有效范围: 1-256）", channelId)
	}
	if receiver == nil {
		return -1, fmt.Errorf("接收器不能为空")
	}

	// 检查是否已经在预览中
	if device.lRealHandle >= 0 {
		log.Printf("警告: 设备已有活动的预览会话（句柄: %d），将先停止旧会话", device.lRealHandle)
		device.StopRealPlay()
	}

	// 配置预览参数
	previewInfo := C.NET_DVR_PREVIEWINFO{}
	previewInfo.lChannel = C.LONG(channelId)  // 通道号
	previewInfo.dwStreamType = C.DWORD(0)     // 码流类型：0-主码流，1-子码流
	previewInfo.dwLinkMode = C.DWORD(0)       // 连接方式：0-TCP，1-UDP
	previewInfo.bBlocked = C.BOOL(0)          // 0-非阻塞，1-阻塞
	previewInfo.byProtoType = C.BYTE(0)       // 传输协议：0-私有协议，1-RTSP
	previewInfo.dwDisplayBufNum = C.DWORD(15) // 播放缓冲区大小

	// 为receiver创建cgo.Handle
	device.receiverHandle = cgo.NewHandle(receiver)

	// 启动实时预览
	// 使用C的包装函数来设置回调
	device.lRealHandle = int(C.NET_DVR_RealPlay_V40_WithCallback(
		C.LONG(device.loginId),
		&previewInfo,
		C.uintptr_t(device.receiverHandle),
	))

	if device.lRealHandle < 0 {
		// 清理Handle
		device.receiverHandle.Delete()
		device.receiverHandle = 0
		return -1, device.HKErr(fmt.Sprintf("启动通道%d的实时预览失败", channelId))
	}

	log.Printf("✓ 实时预览启动成功 - 通道: %d, 句柄: %d", channelId, device.lRealHandle)
	return device.lRealHandle, nil
}

// StopRealPlay 停止实时视频预览
// 停止 RealPlay_V40 启动的视频预览，释放相关资源
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) StopRealPlay() error {
	if device.lRealHandle < 0 {
		// 已经停止，不是错误
		return nil
	}

	// 停止预览
	result := C.NET_DVR_StopRealPlay(C.LONG(device.lRealHandle))

	// 无论成功失败，都清理资源
	oldHandle := device.lRealHandle
	device.lRealHandle = -1

	// 清理receiver handle
	if device.receiverHandle != 0 {
		device.receiverHandle.Delete()
		device.receiverHandle = 0
	}

	if result != C.TRUE {
		return device.HKErr(fmt.Sprintf("停止预览失败（句柄: %d）", oldHandle))
	}

	log.Printf("✓ 实时预览已停止（句柄: %d）", oldHandle)
	return nil
}

// SaveRealData 保存实时流数据到文件
// 将实时视频流保存为文件（如MP4、AVI等）
// 参数：
//   - fileName: 保存的文件名
//
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) SaveRealData(fileName string) error {
	if device.lRealHandle < 0 {
		return fmt.Errorf("未启动实时预览")
	}

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	if C.NET_DVR_SaveRealData(C.LONG(device.lRealHandle), cFileName) != C.TRUE {
		return device.HKErr("保存实时流失败")
	}

	log.Printf("开始保存实时流到文件: %s", fileName)
	return nil
}

// StopSaveRealData 停止保存实时流数据
// 停止 SaveRealData 开始的文件保存操作
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) StopSaveRealData() error {
	if device.lRealHandle < 0 {
		return fmt.Errorf("未启动实时预览")
	}

	if C.NET_DVR_StopSaveRealData(C.LONG(device.lRealHandle)) != C.TRUE {
		return device.HKErr("停止保存实时流失败")
	}

	log.Println("停止保存实时流")
	return nil
}

// CapturePicture 抓图
// 从实时视频流中抓取一帧保存为图片
// 参数：
//   - fileName: 图片文件名（支持bmp、jpg格式）
//
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) CapturePicture(fileName string) error {
	if device.lRealHandle < 0 {
		return fmt.Errorf("未启动实时预览")
	}

	cFileName := C.CString(fileName)
	defer C.free(unsafe.Pointer(cFileName))

	if C.NET_DVR_CapturePicture(C.LONG(device.lRealHandle), cFileName) != C.TRUE {
		return device.HKErr("抓图失败")
	}

	log.Printf("抓图成功: %s", fileName)
	return nil
}
