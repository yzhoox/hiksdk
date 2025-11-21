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
	"strings"
	"unsafe"
)

// LoginV40 登录设备（推荐使用）
// 使用 NET_DVR_Login_V40 接口进行同步登录
// V40版本支持更多功能，性能更好
// 返回值：
//   - int: 登录ID，大于0表示成功，-1表示失败
//   - error: 错误信息，成功时为nil
func (device *HKDevice) LoginV40() (int, error) {
	var deviceInfoV40 C.NET_DVR_DEVICEINFO_V40
	var userLoginInfo C.NET_DVR_USER_LOGIN_INFO

	// 设置登录参数
	strcpy(unsafe.Pointer(&userLoginInfo.sDeviceAddress[0]), device.ip, len(userLoginInfo.sDeviceAddress))
	userLoginInfo.wPort = C.WORD(device.port)
	strcpy(unsafe.Pointer(&userLoginInfo.sUserName[0]), device.username, len(userLoginInfo.sUserName))
	strcpy(unsafe.Pointer(&userLoginInfo.sPassword[0]), device.password, len(userLoginInfo.sPassword))

	// 设置为同步登录模式（0=同步，1=异步）
	userLoginInfo.byUseAsynLogin = 0

	// 调用NET_DVR_Login_V40函数
	device.loginId = int(C.NET_DVR_Login_V40(&userLoginInfo, (*C.NET_DVR_DEVICEINFO_V40)(unsafe.Pointer(&deviceInfoV40))))

	// 提取设备序列号
	serialNumberBytes := make([]byte, len(deviceInfoV40.struDeviceV30.sSerialNumber))
	for i := range deviceInfoV40.struDeviceV30.sSerialNumber {
		serialNumberBytes[i] = byte(deviceInfoV40.struDeviceV30.sSerialNumber[i])
	}
	serialNumber := strings.Trim(string(serialNumberBytes), "\x00")

	if device.loginId < 0 {
		return -1, device.HKErr("登录失败(V40)")
	}

	// 保存通道数
	device.byChanNum = int(deviceInfoV40.struDeviceV30.byChanNum)

	log.Printf("登录成功(V40) - 用户ID: %d, 设备序列号: %s, 通道数: %d",
		device.loginId, serialNumber, device.byChanNum)
	return device.loginId, nil
}

// LoginV30 登录设备（兼容旧设备）
// 使用 NET_DVR_Login_V30 接口进行登录
// 适用于较旧的设备或需要兼容性的场景
// 返回值：
//   - int: 登录ID，大于0表示成功，-1表示失败
//   - error: 错误信息，成功时为nil
func (device *HKDevice) LoginV30() (int, error) {
	var deviceInfoV30 C.NET_DVR_DEVICEINFO_V30

	// 转换参数
	ip := C.CString(device.ip)
	usr := C.CString(device.username)
	passwd := C.CString(device.password)
	defer func() {
		C.free(unsafe.Pointer(ip))
		C.free(unsafe.Pointer(usr))
		C.free(unsafe.Pointer(passwd))
	}()

	// 调用NET_DVR_Login_V30函数
	device.loginId = int(C.NET_DVR_Login_V30(
		ip,
		C.WORD(device.port),
		usr,
		passwd,
		(*C.NET_DVR_DEVICEINFO_V30)(unsafe.Pointer(&deviceInfoV30)),
	))

	// 提取设备序列号
	serialNumberBytes := make([]byte, len(deviceInfoV30.sSerialNumber))
	for i := range deviceInfoV30.sSerialNumber {
		serialNumberBytes[i] = byte(deviceInfoV30.sSerialNumber[i])
	}
	serialNumber := string(serialNumberBytes)
	serialNumber = strings.Trim(serialNumber, "\x00")

	if device.loginId < 0 {
		return -1, device.HKErr("登录失败(V30)")
	}

	// 保存通道数
	device.byChanNum = int(deviceInfoV30.byChanNum)

	log.Printf("登录成功(V30) - 用户ID: %d, 设备序列号: %s, 通道数: %d",
		device.loginId, serialNumber, device.byChanNum)
	return device.loginId, nil
}

// LoginAsync 异步登录设备
// 使用 NET_DVR_Login_V40 接口进行异步登录
// 登录结果通过回调函数返回
// 参数：
//   - callback: 登录结果回调函数
//
// 返回值：
//   - error: 调用错误信息，nil表示成功发起异步登录
func (device *HKDevice) LoginAsync(callback func(loginID int, err error)) error {
	var userLoginInfo C.NET_DVR_USER_LOGIN_INFO

	// 设置登录参数
	strcpy(unsafe.Pointer(&userLoginInfo.sDeviceAddress[0]), device.ip, len(userLoginInfo.sDeviceAddress))
	userLoginInfo.wPort = C.WORD(device.port)
	strcpy(unsafe.Pointer(&userLoginInfo.sUserName[0]), device.username, len(userLoginInfo.sUserName))
	strcpy(unsafe.Pointer(&userLoginInfo.sPassword[0]), device.password, len(userLoginInfo.sPassword))

	// 设置为异步登录模式（1=异步）
	userLoginInfo.byUseAsynLogin = 1

	// 注意：异步登录需要设置回调函数
	// 这里简化处理，实际使用需要更复杂的CGO回调机制
	// userLoginInfo.fLoginResultCallBack = ...
	// userLoginInfo.pUser = ...

	// 暂时返回未实现
	return fmt.Errorf("异步登录暂未实现，请使用同步登录")
}

// Logout 登出设备
// 使用 NET_DVR_Logout 接口断开与设备的连接
// 返回值：
//   - error: 错误信息，成功时为nil
func (device *HKDevice) Logout() error {
	// NET_DVR_Logout 返回 TRUE(1) 表示成功，FALSE(0) 表示失败
	result := C.NET_DVR_Logout(C.LONG(device.loginId))
	if result == 0 {
		return device.HKErr("登出失败")
	}

	device.loginId = -1
	log.Println("登出成功")
	return nil
}

// ResolveDynamicIP 通过解析服务器获取设备的动态IP地址和端口
// 支持通过设备名称或序列号从IPServer/hiDDNS服务器解析设备当前IP
// 参数：
//   - serverIP: 解析服务器的IP地址或域名
//   - serverPort: 解析服务器端口（IPServer: 7071, hiDDNS: 80）
//   - dvrName: 设备名称或域名（可选，与序列号至少提供一个）
//   - serialNumber: 设备序列号（可选，与设备名称至少提供一个）
//
// 返回值：
//   - ip: 解析得到的设备IP地址
//   - port: 解析得到的设备端口号
//   - error: 错误信息
func ResolveDynamicIP(serverIP string, serverPort uint16, dvrName string, serialNumber string) (string, uint32, error) {
	// 验证输入参数
	if serverIP == "" {
		return "", 0, errors.New("解析服务器IP不能为空")
	}
	if dvrName == "" && serialNumber == "" {
		return "", 0, errors.New("设备名称和序列号不能同时为空")
	}

	// 转换参数为C类型
	cServerIP := C.CString(serverIP)
	defer C.free(unsafe.Pointer(cServerIP))

	// 准备设备名称参数
	var cDVRName *C.BYTE
	var dvrNameLen C.WORD
	if dvrName != "" {
		dvrNameBytes := []byte(dvrName)
		cDVRName = (*C.BYTE)(unsafe.Pointer(&dvrNameBytes[0]))
		dvrNameLen = C.WORD(len(dvrNameBytes))
	}

	// 准备序列号参数
	var cSerialNumber *C.BYTE
	var serialLen C.WORD
	if serialNumber != "" {
		serialBytes := []byte(serialNumber)
		cSerialNumber = (*C.BYTE)(unsafe.Pointer(&serialBytes[0]))
		serialLen = C.WORD(len(serialBytes))
	}

	// 准备输出缓冲区
	ipBuffer := make([]byte, 128)
	cGetIP := (*C.char)(unsafe.Pointer(&ipBuffer[0]))
	var dwPort C.DWORD

	// 调用SDK函数
	result := C.NET_DVR_GetDVRIPByResolveSvr_EX(
		cServerIP,
		C.WORD(serverPort),
		cDVRName,
		dvrNameLen,
		cSerialNumber,
		serialLen,
		cGetIP,
		&dwPort,
	)

	if result == 0 {
		errorCode := int(C.NET_DVR_GetLastError())
		return "", 0, fmt.Errorf("解析设备IP失败，错误码: %d", errorCode)
	}

	// 转换结果
	resolvedIP := C.GoString(cGetIP)
	resolvedPort := uint32(dwPort)

	return resolvedIP, resolvedPort, nil
}

// LoginWithDynamicIP 使用动态IP解析后登录设备（V40版本）
// 先通过解析服务器获取设备当前IP，然后进行登录
// 参数：
//   - serverIP: 解析服务器地址
//   - serverPort: 解析服务器端口
//   - dvrName: 设备名称
//   - serialNumber: 设备序列号
//
// 返回值：
//   - int: 登录ID
//   - error: 错误信息
func (device *HKDevice) LoginWithDynamicIP(serverIP string, serverPort uint16, dvrName string, serialNumber string) (int, error) {
	resolvedIP, resolvedPort, err := ResolveDynamicIP(serverIP, serverPort, dvrName, serialNumber)
	if err != nil {
		return -1, err
	}

	device.ip = resolvedIP
	device.port = int(resolvedPort)
	return device.LoginV40()
}

// LoginV30WithDynamicIP 使用动态IP解析后登录设备（V30版本）
// 先通过解析服务器获取设备当前IP，然后使用V30接口登录
// 参数：
//   - serverIP: 解析服务器地址
//   - serverPort: 解析服务器端口
//   - dvrName: 设备名称
//   - serialNumber: 设备序列号
//
// 返回值：
//   - int: 登录ID
//   - error: 错误信息
func (device *HKDevice) LoginV30WithDynamicIP(serverIP string, serverPort uint16, dvrName string, serialNumber string) (int, error) {
	resolvedIP, resolvedPort, err := ResolveDynamicIP(serverIP, serverPort, dvrName, serialNumber)
	if err != nil {
		return -1, err
	}

	device.ip = resolvedIP
	device.port = int(resolvedPort)
	return device.LoginV30()
}
