# 项目结构说明

本文档详细说明 HikSDK 项目的目录结构和文件组织。

## 目录结构

```
hiksdk/
├── pkg/                      # SDK 核心包
│   ├── HKDevice.go          # 设备管理核心实现
│   ├── device.go            # 设备接口定义
│   ├── audio.go             # 音频处理
│   ├── video.go             # 视频处理
│   ├── transceiver.go       # 数据收发器
│   ├── device_test.go       # 设备管理测试
│   └── ptz_test.go          # PTZ 控制测试
│
├── examples/                 # 示例代码目录
│   ├── basic_usage.go       # 基础使用示例
│   ├── ptz_control.go       # PTZ 控制示例
│   ├── video_preview.go     # 视频预览示例
│   └── alarm_listen.go      # 报警监听示例
│
├── include/                  # C 头文件
│   └── HCNetSDK.h           # 海康 SDK 头文件
│
├── lib/                      # 动态链接库
│   ├── Windows/             # Windows 平台库
│   │   ├── HCNetSDK.dll
│   │   ├── HCCore.dll
│   │   ├── PlayCtrl.dll
│   │   └── HCNetSDKCom/     # 组件库
│   └── Linux/               # Linux 平台库
│       ├── libhcnetsdk.so
│       ├── libHCCore.so
│       ├── libPlayCtrl.so
│       └── HCNetSDKCom/     # 组件库
│
├── docs/                     # 文档目录
│   ├── PROJECT_STRUCTURE.md # 项目结构说明（本文件）
│   └── HANDLE_EXPLANATION.md # 句柄概念说明
│
├── go.mod                    # Go 模块定义
├── go.sum                    # Go 依赖校验
├── .gitignore               # Git 忽略文件
│
├── README.md                 # 项目主文档
├── CONTRIBUTING.md          # 贡献指南
└── LICENSE                  # 许可证文件
```

## 核心模块说明

### pkg/ - SDK 核心包

#### HKDevice.go
设备管理的核心实现文件，包含：
- 设备结构体定义
- 登录/登出实现
- PTZ 控制实现
- 视频预览实现
- 报警监听实现
- CGO 调用海康 C SDK

**主要类型：**
- `HKDevice` - 设备对象
- 实现了 `Device` 接口的所有方法

**主要函数：**
- `InitHikSDK()` - 初始化 SDK
- `HKExit()` - 释放 SDK
- `NewHKDevice()` - 创建设备实例
- `Login()` / `LoginV4()` - 登录
- `Logout()` - 登出
- `PTZControlWithSpeed_Other()` - PTZ 控制
- `RealPlay_V40()` - 视频预览
- `SetAlarmCallBack()` - 设置报警回调

#### device.go
设备接口定义，包含：
- `Device` 接口 - 定义设备的所有操作
- `DeviceInfo` 结构体 - 设备信息

#### audio.go
音频相关功能（预留）

#### video.go
视频处理相关功能，包含：
- PS 流处理
- 视频数据解析

#### transceiver.go
数据收发器实现，包含：
- `Receiver` - 数据接收器
- `Sender` - 数据发送器
- 数据通道管理

#### device_test.go
设备管理相关的单元测试：
- `TestDeviceLogin` - 登录测试
- `TestGetDeviceInfo` - 获取设备信息测试
- `TestGetChannelName` - 获取通道名称测试
- `TestRepeatedLogin` - 重复登录测试
- `BenchmarkDeviceLogin` - 登录性能测试

#### ptz_test.go
PTZ 控制相关的单元测试：
- `TestPTZControl` - PTZ 控制测试
- `TestPTZZoom` - 变焦测试
- `TestPTZPreset` - 预置点测试
- `TestGetPTZPosition` - 获取 PTZ 位置测试

### examples/ - 示例代码

#### basic_usage.go
基础使用示例，展示：
- SDK 初始化
- 设备登录/登出
- 获取设备信息
- 获取通道信息

#### ptz_control.go
PTZ 控制示例，展示：
- 云台移动控制
- 变焦控制
- 预置点设置和使用
- 各种 PTZ 命令的使用

#### video_preview.go
视频预览示例，展示：
- 创建接收器
- 启动视频预览
- 接收和处理视频数据
- 统计视频流信息

#### alarm_listen.go
报警监听示例，展示：
- 设置报警回调
- 启动报警监听
- 等待报警事件
- 优雅退出

### include/ - C 头文件

包含海康威视 SDK 的 C 头文件，用于 CGO 调用。

### lib/ - 动态链接库

包含海康威视 SDK 的动态链接库文件：
- Windows 平台：DLL 文件
- Linux 平台：SO 文件

**注意：** 这些库文件体积较大，通常不包含在 Git 仓库中。

## 文档说明

### README.md
项目的主要文档，包含：
- 项目简介
- 功能特性
- 安装说明
- 快速开始
- API 文档
- 常见问题
- 贡献指南

### QUICKSTART.md
快速开始指南，适合新用户：
- 最小示例
- 常见用法
- 快速上手

### CHANGELOG.md
版本更新日志，记录：
- 版本历史
- 新增功能
- Bug 修复
- 破坏性更改

### CONTRIBUTING.md
贡献指南，包含：
- 如何报告问题
- 如何提交代码
- 代码规范
- 开发流程

### LICENSE
MIT 许可证文件

## 配置文件

### go.mod
Go 模块定义文件，管理依赖包
```go
module github.com/samsaralc/hiksdk
go 1.25
require golang.org/x/text v0.31.0
```

### .gitignore
Git 忽略规则，重要配置：
- 忽略一般的 .dll 和 .so 文件
- 但 lib/ 目录下的库文件不被忽略（使用 !lib/**/*.dll 规则）
- 确保海康 SDK 库文件被 Git 托管

## 开发流程

### 1. 开发新功能

1. 在 `pkg/` 中实现功能
2. 在 `pkg/*_test.go` 中添加测试
3. 在 `examples/` 中添加示例（如果需要）
4. 更新 `README.md` 和其他文档
5. 运行 `go fmt ./...` 格式化代码
6. 运行 `go vet ./...` 检查代码
7. 运行 `go test ./pkg/...` 测试功能

### 2. 测试

```bash
# 设置环境变量
export HIK_IP="192.168.1.64"
export HIK_USER="admin"
export HIK_PASSWORD="password"

# 运行测试
make test

# 代码检查
make lint
```

### 3. 运行示例

```bash
# 直接运行示例（修改代码中的 IP、用户名、密码）
go run examples/basic_usage.go
go run examples/ptz_control.go

# 或编译后运行
go build -o bin/basic_usage examples/basic_usage.go
./bin/basic_usage  # Linux/macOS
# 或
.\bin\basic_usage.exe  # Windows
```

### 4. 清理构建文件

```bash
# 清理编译生成的文件
go clean

# 删除 bin 目录
rm -rf bin/  # Linux/macOS
# 或
Remove-Item -Recurse bin\  # Windows
```

## 依赖管理

项目使用 Go Modules 管理依赖：

```bash
# 添加依赖
go get package-name

# 更新依赖
go get -u

# 整理依赖
go mod tidy
```

## 平台差异

### Windows
- 使用 `lib/Windows/` 下的 DLL 文件
- 需要 MinGW-w64 进行 CGO 编译

### Linux
- 使用 `lib/Linux/` 下的 SO 文件
- 需要 gcc 进行 CGO 编译

## 扩展开发

### 添加新的设备功能

1. 在 `pkg/device.go` 的 `Device` 接口中添加方法定义
2. 在 `pkg/HKDevice.go` 中实现该方法
3. 添加对应的测试
4. 更新文档

### 添加新的示例

1. 在 `examples/` 中创建新文件
2. 实现示例代码
3. 在 README.md 中添加说明
4. 在 Makefile 中添加构建目标（如果需要）

## 注意事项

1. **CGO 依赖**：本项目使用 CGO 调用 C 库，编译时需要 CGO 编译器
2. **动态库路径**：运行时需要确保动态库在系统路径中
3. **设备限制**：某些功能受设备型号和固件版本限制
4. **线程安全**：注意 CGO 调用的线程安全性

## 相关资源

- [Go Modules 文档](https://golang.org/ref/mod)
- [CGO 文档](https://golang.org/cmd/cgo/)
- [海康威视开发者平台](https://open.hikvision.com/)

---

最后更新：2025-11-20

