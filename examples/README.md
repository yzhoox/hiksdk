# 海康威视 SDK Go 示例代码

本目录包含海康威视 SDK 的测试示例，展示了各种功能的使用方法。

## 快速开始

### 1. 修改设备信息

在测试文件中修改设备连接信息：

```go
cred := &auth.Credentials{
    IP:       "192.168.1.64", // 替换为你的设备IP
    Port:     8000,            // 替换为你的端口
    Username: "admin",         // 替换为你的用户名
    Password: "password",      // 替换为你的密码
}
```

### 2. 运行示例

```bash
# 运行所有测试示例
go test -v

# 运行特定测试
go test -v -run TestLoginMethods
go test -v -run TestPTZControl
go test -v -run TestAlarmListen
go test -v -run TestCruiseTrack
go test -v -run TestPTZAdvanced
```

## 示例列表

| 测试文件 | 功能说明 |
|---------|---------|
| `login_test.go` | 两种登录方式（V40推荐、V30兼容） |
| `ptz_control_test.go` | PTZ云台控制（方向、变焦、预置点） |
| `alarm_listen_test.go` | 报警监听（移动侦测、遮挡报警等） |
| `cruise_track_test.go` | 巡航与轨迹（自动巡航路径、轨迹录制回放） |
| `ptz_advanced_test.go` | PTZ高级控制（自动扫描、辅助设备） |
| `error_handling_test.go` | 错误处理（详细的错误码和错误描述） |

## 最简示例

### 登录示例

```go
package main

import (
    "fmt"
    "github.com/samsaralc/hiksdk/core/auth"
)

func main() {
    // 创建设备连接凭据
    cred := &auth.Credentials{
        IP:       "192.168.1.64",
        Port:     8000,
        Username: "admin",
        Password: "password",
    }
    
    // 登录设备（推荐V40）
    session, err := auth.LoginV40(cred)
    if err != nil {
        // 如果V40失败，尝试V30（兼容旧设备）
        session, err = auth.LoginV30(cred)
        if err != nil {
            fmt.Printf("登录失败: %v\n", err)
            return
        }
    }
    defer auth.Logout(session.LoginID)
    defer auth.Cleanup()
    
    fmt.Printf("登录成功! 设备序列号: %s\n", session.SerialNumber)
}
```

### PTZ控制示例

```go
package main

import (
    "time"
    "github.com/samsaralc/hiksdk/core/auth"
    "github.com/samsaralc/hiksdk/core/ptz"
)

func main() {
    // 登录设备
    cred := &auth.Credentials{IP: "192.168.1.64", Port: 8000, Username: "admin", Password: "password"}
    session, _ := auth.LoginV30(cred)
    defer auth.Logout(session.LoginID)
    defer auth.Cleanup()
    
    // 创建PTZ控制器（统一控制云台、相机、辅助设备）
    ctrl := ptz.NewController(session.LoginID, 1)
    
    // 方式1：自动控制时长（简单）
    ctrl.Right(5, 2*time.Second)  // 右转2秒后自动停止
    ctrl.Up(5, 2*time.Second)     // 上仰2秒后自动停止
    
    // 方式2：手动开始/停止（灵活）
    ctrl.StartLeft(5)              // 开始左转
    time.Sleep(3*time.Second)      // 自己控制时长
    ctrl.StopLeft()                // 停止左转
    
    // 相机控制
    ctrl.ZoomIn(1*time.Second)     // 焦距放大1秒
    
    // 辅助设备
    ctrl.LightOn()                 // 开启灯光
    time.Sleep(2*time.Second)
    ctrl.LightOff()                // 关闭灯光
}
```

查看 `examples` 目录中的测试文件了解更多功能使用方法。

## 注意事项

1. **SDK自动初始化** - 第一次调用登录时SDK会自动初始化
2. **设备限制** - 设备最多支持32个注册用户名，同时最多128个用户注册
3. **SDK限制** - SDK支持2048个注册，UserID取值范围0~2047
4. **资源管理** - 登录后必须调用 `Logout()` 释放资源
5. **错误处理** - 所有函数返回 `error` 类型，包含详细错误信息
6. **测试环境** - 示例中的设备信息需要根据实际环境修改

## 登录方式选择

### LoginV40() - 推荐使用
```go
session, err := auth.LoginV40(cred)
```
- ✅ 支持更多功能
- ✅ 更好的性能
- ✅ 适用于新设备
- ✅ 支持同步/异步模式

### LoginV30() - 兼容旧设备
```go
session, err := auth.LoginV30(cred)
```
- ✅ 兼容旧设备
- ✅ 简单直接
- ⚠️ 功能相对较少
- ⚠️ 只支持同步模式

## 动态IP登录

如果设备使用动态IP，可以通过解析服务器获取：

```go
// 使用动态IP登录（V40版本）
session, err := auth.LoginV40WithDynamicIP(
    cred,
    "ipserver.com",  // 解析服务器地址
    7071,            // IPServer端口(7071)或hiDDNS端口(80)
    "device-name",   // 设备名称
    "serial-number", // 设备序列号
)

// 或使用V30版本
session, err := auth.LoginV30WithDynamicIP(
    cred,
    "ipserver.com",
    7071,
    "device-name",
    "serial-number",
)
```

## 错误处理

所有函数都返回 `*core.HKError` 错误对象，包含错误码、错误描述、操作名称等详细信息：

```go
session, err := auth.LoginV30(cred)
if err != nil {
    // 错误信息格式：登录设备(V30)失败，错误码: 1, 用户名或密码错误
    fmt.Printf("错误: %v\n", err)
    
    // 可以获取详细信息
    if hkErr, ok := err.(*core.HKError); ok {
        fmt.Printf("错误码: %d\n", hkErr.Code)
        fmt.Printf("错误描述: %s\n", hkErr.Msg)
        fmt.Printf("操作: %s\n", hkErr.Operation)
        fmt.Printf("JSON格式: %s\n", hkErr.JSON())
    }
}
```

### 常见错误码

| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 1 | 用户名或密码错误 | 检查设备凭据是否正确 |
| 7 | 连接设备失败 | 检查设备是否在线、网络是否通畅、IP和端口是否正确 |
| 4 | 通道号错误 | 检查通道号是否在有效范围内 |
| 17 | 参数错误 | 检查函数参数是否符合要求 |
| 23 | 设备不支持该功能 | 更换支持该功能的设备或使用其他功能 |
| 30 | 用户数达到最大 | 断开其他连接或等待空闲 |

## 更多示例

查看 `examples` 目录中的测试文件了解更多功能使用方法。
