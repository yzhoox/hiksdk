package core

/*
#include <stdio.h>
#include <stdlib.h>
#include "hiksdk_wrapper.h"

// 声明Go回调函数
extern void AlarmCallBack(LONG command, NET_DVR_ALARMER *alarm, char *info, DWORD len, void *user);
*/
import "C"
import (
	"fmt"
	"log"
	"unsafe"

	"github.com/samsaralc/hiksdk/consts"
)

// AlarmCallBack 报警回调函数
// 由C代码调用，接收设备的报警信息
// 这是一个全局回调函数，会被所有设备共享
//
//export AlarmCallBack
func AlarmCallBack(command C.LONG, alarm *C.NET_DVR_ALARMER, info *C.char, len C.DWORD, user unsafe.Pointer) {
	// 安全检查
	if alarm == nil {
		log.Println("警告: 收到空报警信息")
		return
	}

	// 获取报警设备信息
	var deviceIP string
	if alarm.byDeviceIPValid == 1 {
		// 提取设备IP（需要从DWORD转换为IP字符串）
		ip := alarm.dwDeviceIP
		deviceIP = fmt.Sprintf("%d.%d.%d.%d",
			byte(ip&0xFF),
			byte((ip>>8)&0xFF),
			byte((ip>>16)&0xFF),
			byte((ip>>24)&0xFF))
	}

	// 获取设备序列号
	var serialNo string
	if alarm.bySerialValid == 1 {
		// 手动复制序列号字节
		serialBytes := make([]byte, consts.MAX_SERIALNO_LEN)
		for i := 0; i < consts.MAX_SERIALNO_LEN; i++ {
			serialBytes[i] = byte(alarm.sSerialNumber[i])
		}
		// 去除空字符
		for i, b := range serialBytes {
			if b == 0 {
				serialNo = string(serialBytes[:i])
				break
			}
		}
		if serialNo == "" {
			serialNo = string(serialBytes)
		}
	}

	// 记录报警信息
	log.Printf("收到报警 - 类型: 0x%X, 设备IP: %s, 序列号: %s", command, deviceIP, serialNo)

	// 根据命令类型处理不同的报警
	switch int(command) {
	case consts.COMM_ALARM_RULE:
		log.Println("  → 行为分析报警（规则检测）")
	case consts.COMM_ALARM_V30:
		log.Println("  → 标准报警（移动侦测/视频丢失/遮挡/IO信号等）")
	case consts.COMM_ALARM_V40:
		log.Println("  → 扩展报警V40（增强型报警信息）")
	default:
		log.Printf("  → 其他报警类型: 0x%X\n", command)
	}

	// TODO: 这里可以添加自定义的报警处理逻辑
	// 例如：发送通知、记录数据库、触发其他动作等
}

// SetAlarmCallBack 设置报警回调函数
// 必须在 StartListenAlarmMsg 之前调用
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) SetAlarmCallBack() error {
	// 设置报警消息回调函数
	// 注意：这里使用全局函数AlarmCallBack，所以无法直接访问device实例
	// 如果需要访问device，可以通过user参数传递
	C.NET_DVR_SetDVRMessageCallBack_V30(
		(*[0]byte)(C.AlarmCallBack),
		unsafe.Pointer(device),
	)

	log.Println("报警回调函数设置成功")
	return nil
}

// StartListenAlarmMsg 启动报警监听
// 建立报警上传通道，开始接收设备的报警信息
// 如移动侦测、遮挡报警、信号丢失、硬盘满、硬盘故障等
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) StartListenAlarmMsg() error {
	// 设置报警布防参数
	var setupParam C.NET_DVR_SETUPALARM_PARAM
	setupParam.dwSize = C.DWORD(unsafe.Sizeof(setupParam))
	setupParam.byLevel = 1         // 布防等级
	setupParam.byAlarmInfoType = 1 // 上传报警信息类型：0-老报警信息，1-新报警信息

	// 建立报警上传通道
	device.alarmHandle = int(C.NET_DVR_SetupAlarmChan_V41(
		C.LONG(device.loginId),
		&setupParam,
	))

	if device.alarmHandle < 0 {
		return device.HKErr("建立报警通道失败")
	}

	log.Printf("报警监听启动成功，句柄: %d", device.alarmHandle)
	return nil
}

// StopListenAlarmMsg 停止报警监听
// 撤销报警上传通道，停止接收设备的报警信息
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) StopListenAlarmMsg() error {
	if device.alarmHandle >= 0 {
		if C.NET_DVR_CloseAlarmChan_V30(C.LONG(device.alarmHandle)) != C.TRUE {
			return device.HKErr("关闭报警通道失败")
		}

		device.alarmHandle = -1
		log.Println("报警监听已停止")
	}
	return nil
}

// SetAlarmOut 控制报警输出
// 控制设备的报警输出端口，如声光报警器
// 参数：
//   - alarmOutPort: 报警输出端口号，从1开始
//   - alarmOutStatic: 报警输出状态，0=停止输出，1=开始输出
//
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) SetAlarmOut(alarmOutPort, alarmOutStatic int) error {
	if C.NET_DVR_SetAlarmOut(
		C.LONG(device.loginId),
		C.LONG(alarmOutPort),
		C.DWORD(alarmOutStatic),
	) != C.TRUE {
		return device.HKErr(fmt.Sprintf("设置报警输出端口 %d 失败", alarmOutPort))
	}

	status := "停止"
	if alarmOutStatic == 1 {
		status = "触发"
	}
	log.Printf("报警输出端口 %d %s", alarmOutPort, status)
	return nil
}

// GetAlarmOut 获取报警输出状态
// 查询设备报警输出端口的当前状态
// 参数：
//   - alarmOutPort: 报警输出端口号，从1开始
//
// 返回值：
//   - int: 报警输出状态，0=停止输出，1=正在输出
//   - error: 错误信息，成功时为nil
func (device *HKDevice) GetAlarmOut(alarmOutPort int) (int, error) {
	var alarmOutState C.DWORD

	if C.NET_DVR_GetAlarmOut(
		C.LONG(device.loginId),
		C.LONG(alarmOutPort),
		&alarmOutState,
	) != C.TRUE {
		return -1, device.HKErr(fmt.Sprintf("获取报警输出端口 %d 状态失败", alarmOutPort))
	}

	return int(alarmOutState), nil
}

// AlarmHostConfig 报警主机配置
type AlarmHostConfig struct {
	Enable      bool   // 是否启用
	IPAddress   string // 报警主机IP
	Port        int    // 报警主机端口
	Protocol    int    // 协议类型
	UserName    string // 用户名
	Password    string // 密码
	AlarmChanNo int    // 报警通道号
}

// SetAlarmHost 配置报警主机
// 配置第三方报警主机的连接参数
// 参数：
//   - config: 报警主机配置参数
//
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) SetAlarmHost(config *AlarmHostConfig) error {
	// 可以根据需要实现报警主机配置
	// 使用 NET_DVR_SetDVRConfig 配置相关参数
	return nil
}
