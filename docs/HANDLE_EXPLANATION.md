# 海康 SDK 句柄说明

## 📚 loginId 和 lRealHandle 的概念

在海康威视 SDK 中，有两个重要的"句柄"（Handle）概念，它们代表不同层次的连接状态。

---

## 1️⃣ **loginId（登录句柄）**

### 定义
`loginId` 是设备**登录成功后**返回的唯一标识符，代表与设备的**会话连接**。

### 获取方式
```go
// 调用登录函数后获得
device.loginId = int(C.NET_DVR_Login_V30(ip, port, username, password, &deviceInfo))

// 或者
device.loginId = int(C.NET_DVR_Login_V40(...))
```

### 生命周期
```
用户调用 Login() 
    ↓
SDK 连接设备
    ↓
返回 loginId (例如: 1, 2, 3...)
    ↓
用户可以使用 loginId 进行各种设备操作
    ↓
用户调用 Logout()
    ↓
loginId 失效
```

### 用途（loginId 能做什么）
使用 `loginId` 可以进行**设备级别**的操作：

| 功能类别 | API 示例 | 说明 |
|---------|---------|------|
| **设备配置** | `NET_DVR_GetDVRConfig(loginId, ...)` | 获取设备配置 |
| **设备信息** | `NET_DVR_GetDeviceConfig(loginId, ...)` | 获取设备信息 |
| **通道控制** | `NET_DVR_PTZControl_Other(loginId, channel, ...)` | 控制指定通道PTZ |
| **报警监听** | `NET_DVR_SetupAlarmChan_V41(loginId, ...)` | 设置报警通道 |
| **视频预览** | `NET_DVR_RealPlay_V40(loginId, ...)` | 启动视频预览 |
| **回放控制** | `NET_DVR_PlayBackByTime(loginId, ...)` | 按时间回放 |

### 代码示例
```go
// 登录获取 loginId
dev := pkg.NewHKDevice(deviceInfo)
loginId, err := dev.Login()
// loginId = 1（假设返回1）

// 使用 loginId 控制PTZ（通过通道号）
dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 4)
// 内部调用: NET_DVR_PTZControlWithSpeed_Other(loginId, channelId, ...)

// 使用 loginId 获取设备信息
info, _ := dev.GetDeviceInfo()
// 内部调用: NET_DVR_GetDVRConfig(loginId, ...)

// 登出后 loginId 失效
dev.Logout()
```

---

## 2️⃣ **lRealHandle（实时预览句柄）**

### 定义
`lRealHandle` 是**启动实时视频预览后**返回的句柄，代表一个**视频流会话**。

### 获取方式
```go
// 调用实时预览函数后获得
device.lRealHandle = int(C.NET_DVR_RealPlay_V40(loginId, &previewInfo, callback, userData))
```

### 生命周期
```
用户已登录（有 loginId）
    ↓
用户调用 RealPlay_V40(channelId, receiver)
    ↓
SDK 建立视频流连接
    ↓
返回 lRealHandle (例如: 100, 101, 102...)
    ↓
用户可以使用 lRealHandle 控制这个视频流
    ↓
用户调用 StopRealPlay()
    ↓
lRealHandle 失效
```

### 用途（lRealHandle 能做什么）
使用 `lRealHandle` 可以进行**视频流级别**的操作：

| 功能类别 | API 示例 | 说明 |
|---------|---------|------|
| **PTZ 控制** | `NET_DVR_PTZControlWithSpeed(lRealHandle, ...)` | 控制当前预览通道的PTZ |
| **停止预览** | `NET_DVR_StopRealPlay(lRealHandle)` | 停止视频预览 |
| **录像控制** | `NET_DVR_SaveRealData(lRealHandle, filename)` | 保存视频到文件 |
| **抓图** | `NET_DVR_CapturePicture(lRealHandle, ...)` | 从视频流抓图 |
| **音频** | `NET_DVR_OpenSound(lRealHandle)` | 打开声音 |
| **视频参数** | `NET_DVR_GetVideoEffect(lRealHandle, ...)` | 获取视频效果 |

### 代码示例
```go
// 假设已经登录（loginId = 1）
receiver := &pkg.Receiver{}
receiver.Start()

// 启动预览获取 lRealHandle
lRealHandle, err := dev.RealPlay_V40(1, receiver)
// lRealHandle = 100（假设返回100）

// 使用 lRealHandle 控制PTZ
dev.PTZControlWithSpeed(PAN_RIGHT, 0, 4)
// 内部调用: NET_DVR_PTZControlWithSpeed(lRealHandle, ...)

// 停止预览后 lRealHandle 失效
dev.StopRealPlay()
```

---

## 🔄 **两者的关系**

```
设备
  ↓
用户登录 → loginId（设备会话）
  ↓
  ├─ 启动预览通道1 → lRealHandle1（视频流1）
  ├─ 启动预览通道2 → lRealHandle2（视频流2）
  └─ 启动预览通道3 → lRealHandle3（视频流3）
```

### 层级关系
```
1 个 loginId（设备连接）
    ↓ 可以创建
多个 lRealHandle（多个视频流）
```

### 依赖关系
- **loginId 是前提**：必须先登录获得 loginId
- **lRealHandle 依赖 loginId**：使用 loginId 来创建 lRealHandle
- **独立生命周期**：停止视频预览不影响登录状态

---

## ⚖️ **对比总结**

| 特性 | loginId | lRealHandle |
|-----|---------|-------------|
| **中文名** | 登录句柄/用户句柄 | 实时预览句柄/播放句柄 |
| **代表** | 设备连接会话 | 视频流会话 |
| **获取** | `NET_DVR_Login_V30/V40` | `NET_DVR_RealPlay_V40` |
| **数量** | 每个设备1个 | 可以有多个（多通道） |
| **作用域** | 设备级别操作 | 视频流级别操作 |
| **生命周期** | 登录→登出 | 预览→停止预览 |
| **必需性** | 必须有 | 可选（不预览可以不要） |

---

## 🎯 **PTZ 控制的两种方式**

这就是为什么 PTZ 控制有两套 API：

### 方式 1：使用 loginId（推荐）✅
```go
// 不需要预览，直接控制
dev.PTZControlWithSpeed_Other(channelId, PAN_RIGHT, 0, 4)
// 参数: loginId, channelId, command, stop, speed
```

**优点：**
- ✅ 不需要启动视频预览
- ✅ 可以控制任意通道
- ✅ 更灵活
- ✅ 资源占用少

### 方式 2：使用 lRealHandle
```go
// 必须先启动预览
receiver := &pkg.Receiver{}
receiver.Start()
dev.RealPlay_V40(1, receiver) // 获取 lRealHandle

// 然后才能控制
dev.PTZControlWithSpeed(PAN_RIGHT, 0, 4)
// 参数: lRealHandle, command, stop, speed
```

**限制：**
- ⚠️ 必须先启动视频预览
- ⚠️ 只能控制正在预览的通道
- ⚠️ 占用更多资源（视频流）

---

## 📖 **实际应用场景**

### 场景 1：只需要 PTZ 控制（不看视频）
```go
// 只需要 loginId
dev.Login()
dev.PTZControlWithSpeed_Other(1, PAN_RIGHT, 0, 4) // ✅
dev.Logout()
```

### 场景 2：预览视频 + PTZ 控制
```go
// 需要 loginId 和 lRealHandle
dev.Login()

receiver := &pkg.Receiver{}
receiver.Start()
dev.RealPlay_V40(1, receiver) // 获取 lRealHandle

// 两种控制方式都可以
dev.PTZControlWithSpeed(PAN_RIGHT, 0, 4)            // 方式1: 使用 lRealHandle
dev.PTZControlWithSpeed_Other(1, PAN_RIGHT, 0, 4)  // 方式2: 使用 loginId

dev.StopRealPlay()
dev.Logout()
```

### 场景 3：多通道预览 + PTZ 控制
```go
dev.Login() // loginId = 1

// 启动多个通道预览
receiver1 := &pkg.Receiver{}
receiver1.Start()
lRealHandle1, _ := dev.RealPlay_V40(1, receiver1) // 通道1: lRealHandle = 100

receiver2 := &pkg.Receiver{}
receiver2.Start()
lRealHandle2, _ := dev.RealPlay_V40(2, receiver2) // 通道2: lRealHandle = 101

// 使用 loginId 可以控制任意通道
dev.PTZControlWithSpeed_Other(1, PAN_RIGHT, 0, 4) // 控制通道1
dev.PTZControlWithSpeed_Other(2, PAN_LEFT, 0, 4)  // 控制通道2

// 使用 lRealHandle 只能控制对应通道
// 这里有问题：SDK 不知道你想用哪个 lRealHandle
// 所以推荐用 PTZControlWithSpeed_Other
```

---

## 🐛 **之前的 Bug 原因**

### 错误的实现（已修复）
```go
// ❌ 错误：使用 loginId 调用需要 lRealHandle 的函数
func (device *HKDevice) PTZControlWithSpeed(...) {
    C.NET_DVR_PTZControlWithSpeed(C.LONG(device.loginId), ...) // 错误！
}
```

### 正确的实现
```go
// ✅ 正确：使用 lRealHandle
func (device *HKDevice) PTZControlWithSpeed(...) {
    if device.lRealHandle == 0 {
        return false, errors.New("需要先启动预览")
    }
    C.NET_DVR_PTZControlWithSpeed(C.LONG(device.lRealHandle), ...) // 正确！
}
```

**这就是为什么之前 `PTZControlWithSpeed` 不工作的原因！**

---

## 💡 **记忆口诀**

```
loginId   - 设备的钥匙（一把钥匙开一个设备）
lRealHandle - 视频的遥控器（一个遥控器控一个视频流）

PTZ 控制：
- 用钥匙（loginId）+ 通道号 = 灵活，推荐 ✅
- 用遥控器（lRealHandle）= 需要先开视频 ⚠️
```

---

## 📚 **相关 API 速查**

### 需要 loginId 的函数
- `NET_DVR_Logout_V30(loginId)`
- `NET_DVR_GetDVRConfig(loginId, ...)`
- `NET_DVR_PTZControl_Other(loginId, channel, ...)`
- `NET_DVR_PTZControlWithSpeed_Other(loginId, channel, ...)`
- `NET_DVR_RealPlay_V40(loginId, ...)`

### 需要 lRealHandle 的函数
- `NET_DVR_StopRealPlay(lRealHandle)`
- `NET_DVR_PTZControl(lRealHandle, ...)`
- `NET_DVR_PTZControlWithSpeed(lRealHandle, ...)`
- `NET_DVR_SaveRealData(lRealHandle, ...)`
- `NET_DVR_OpenSound(lRealHandle)`

---

**总结：loginId 是设备级别的连接，lRealHandle 是视频流级别的连接。大多数情况下推荐使用 loginId + 通道号的方式（_Other 系列函数）！**

