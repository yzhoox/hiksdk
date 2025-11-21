# 海康威视 SDK Go 示例代码

本目录包含海康威视 SDK 的示例代码，展示了各种功能的使用方法。

## 示例列表

### 01_login_methods.go - 两种登录方式
展示如何使用两种不同的登录方法登录设备。

**主要功能：**
- Login() - 使用 NET_DVR_Login_V40（推荐）
- LoginV30() - 使用 NET_DVR_Login_V30（兼容旧设备）
- 获取设备信息
- 设备登出

### 02_device_info.go - 获取设备信息
展示如何获取设备的详细信息和通道名称。

**主要功能：**
- 获取设备基本信息（名称、序列号、IP等）
- 获取通道数量
- 获取各通道名称

### 03_ptz_control.go - PTZ云台控制
展示如何控制PTZ云台和预置点。

**主要功能：**
- 方向控制（上下左右）
- 变焦控制（放大缩小）
- 预置点操作（设置、调用）

### 04_video_preview.go - 视频预览
展示如何接收实时视频流。

**主要功能：**
- 启动实时预览
- 接收PS流数据
- 数据统计

### 05_alarm_listen.go - 报警监听
展示如何监听设备报警事件。

**主要功能：**
- 设置报警回调
- 启动报警监听
- 接收报警消息

## 使用方法

### 1. 修改设备信息

每个示例中都需要修改设备连接信息：

```go
deviceInfo := core.DeviceInfo{
    IP:       "192.168.1.64", // 替换为你的设备IP
    Port:     8000,            // 替换为你的端口
    UserName: "admin",         // 替换为你的用户名
    Password: "password",      // 替换为你的密码
}
```

### 2. 运行示例

```bash
cd examples
go run 01_login_methods.go
```

## 登录方式选择

### Login() - 推荐使用
```go
dev := core.NewHKDevice(deviceInfo)
loginId, err := dev.Login()  // 使用 NET_DVR_Login_V40
```
- ✅ 支持更多功能
- ✅ 更好的性能
- ✅ 适用于新设备
- ✅ 支持同步/异步模式

### LoginV30() - 兼容旧设备
```go
dev := core.NewHKDevice(deviceInfo)
loginId, err := dev.LoginV30()  // 使用 NET_DVR_Login_V30
```
- ✅ 兼容旧设备
- ✅ 简单直接
- ⚠️ 功能相对较少
- ⚠️ 只支持同步模式

## 注意事项

1. **SDK自动初始化** - 第一次创建设备时SDK会自动初始化
2. **设备限制** - 设备最多支持32个注册用户名，同时最多128个用户注册
3. **SDK限制** - SDK支持2048个注册，UserID取值范围0~2047
4. **资源管理** - 登录后必须调用 Logout() 释放资源
5. **错误处理** - 所有函数返回 HKError 类型，包含详细错误信息

## 最简示例

```go
package main

import (
    "fmt"
    "github.com/samsaralc/hiksdk/core"
)

func main() {
    // 创建设备（SDK自动初始化）
    dev := core.NewHKDevice(core.DeviceInfo{
        IP:       "192.168.1.64",
        Port:     8000,
        UserName: "admin",
        Password: "password",
    })
    
    // 登录（推荐V40）
    if _, err := dev.Login(); err != nil {
        // 如果V40失败，尝试V30
        if _, err := dev.LoginV30(); err != nil {
            fmt.Printf("登录失败: %v\n", err)
            return
        }
    }
    defer dev.Logout()
    
    fmt.Println("登录成功!")
}
```

## 动态IP登录

如果设备使用动态IP，可以通过解析服务器获取：

```go
// 使用动态IP登录（V40版本）
loginId, err := dev.LoginWithDynamicIP(
    "ipserver.com",  // 解析服务器地址
    7071,            // IPServer端口(7071)或hiDDNS端口(80)
    "device-name",   // 设备名称
    "serial-number", // 设备序列号
)

// 或使用V30版本
loginId, err := dev.LoginV30WithDynamicIP(
    "ipserver.com",
    7071,
    "device-name",
    "serial-number",
)
```