package pkg

// Device 海康威视设备接口，定义了设备操作的标准方法
// 所有海康设备的实现都应该满足这个接口
type Device interface {
	// Login 登录设备（使用 NET_DVR_Login_V30）
	// 返回值：
	//   - int: 登录ID，大于0表示成功，-1表示失败
	//   - error: 错误信息，成功时为nil
	Login() (int, error)

	// LoginV4 登录设备（使用 NET_DVR_Login_V40）
	// 相比 Login 方法，V40 版本支持更多功能和更好的性能
	// 返回值：
	//   - int: 登录ID，大于0表示成功，-1表示失败
	//   - error: 错误信息，成功时为nil
	LoginV4() (int, error)

	// GetDeviceInfo 获取设备详细信息（推荐使用）
	// 返回值：
	//   - *DeviceInfo: 设备信息结构体，包含设备名称、序列号、通道数等
	//   - error: 错误信息，成功时为nil
	GetDeviceInfo() (*DeviceInfo, error)

	// GetDeiceInfo 获取设备详细信息
	// Deprecated: 请使用 GetDeviceInfo（拼写错误）
	// 返回值：
	//   - *DeviceInfo: 设备信息结构体，包含设备名称、序列号、通道数等
	//   - error: 错误信息，成功时为nil
	GetDeiceInfo() (*DeviceInfo, error)

	// GetChannelName 获取所有通道的名称
	// 返回值：
	//   - map[int]string: 通道ID到通道名称的映射
	//   - error: 错误信息，成功时为nil
	GetChannelName() (map[int]string, error)

	// Logout 登出设备，断开与设备的连接
	// 返回值：
	//   - error: 错误信息，成功时为nil
	Logout() error

	// SetAlarmCallBack 设置报警回调函数
	// 必须在 StartListenAlarmMsg 之前调用
	// 返回值：
	//   - error: 错误信息，成功时为nil
	SetAlarmCallBack() error

	// StartListenAlarmMsg 启动报警监听
	// 开始接收设备的报警信息，如移动侦测、遮挡报警等
	// 返回值：
	//   - error: 错误信息，成功时为nil
	StartListenAlarmMsg() error

	// StopListenAlarmMsg 停止报警监听
	// 停止接收设备的报警信息
	// 返回值：
	//   - error: 错误信息，成功时为nil
	StopListenAlarmMsg() error

	// PTZControlWithSpeed PTZ云台控制（带速度参数）
	// 使用实时播放句柄控制，需要先调用 RealPlay_V40 启动预览
	// 参数：
	//   - dwPTZCommand: PTZ命令，如 PAN_RIGHT(24)、TILT_UP(21) 等
	//   - dwStop: 停止标志，0=开始动作，1=停止动作
	//   - dwSpeed: 速度值，范围0-7，数值越大速度越快
	// 返回值：
	//   - bool: 操作是否成功
	//   - error: 错误信息，成功时为nil
	PTZControlWithSpeed(dwPTZCommand, dwStop, dwSpeed int) (bool, error)

	// PTZControlWithSpeed_Other PTZ云台控制（指定通道，带速度参数）
	// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
	// 参数：
	//   - lChannel: 通道号，从1开始
	//   - dwPTZCommand: PTZ命令，如 PAN_RIGHT(24)、TILT_UP(21) 等
	//   - dwStop: 停止标志，0=开始动作，1=停止动作
	//   - dwSpeed: 速度值，范围0-7，数值越大速度越快
	// 返回值：
	//   - bool: 操作是否成功
	//   - error: 错误信息，成功时为nil
	PTZControlWithSpeed_Other(lChannel, dwPTZCommand, dwStop, dwSpeed int) (bool, error)

	// PTZControl PTZ云台控制（无速度参数）
	// 使用实时播放句柄控制，需要先调用 RealPlay_V40 启动预览
	// 参数：
	//   - dwPTZCommand: PTZ命令，如 ZOOM_IN(11)、FOCUS_NEAR(13) 等
	//   - dwStop: 停止标志，0=开始动作，1=停止动作
	// 返回值：
	//   - bool: 操作是否成功
	//   - error: 错误信息，成功时为nil
	PTZControl(dwPTZCommand, dwStop int) (bool, error)

	// PTZControl_Other PTZ云台控制（指定通道，无速度参数）
	// 使用登录ID和通道号控制，不需要预览即可使用（推荐使用）
	// 参数：
	//   - lChannel: 通道号，从1开始
	//   - dwPTZCommand: PTZ命令，如 ZOOM_IN(11)、FOCUS_NEAR(13) 等
	//   - dwStop: 停止标志，0=开始动作，1=停止动作
	// 返回值：
	//   - bool: 操作是否成功
	//   - error: 错误信息，成功时为nil
	PTZControl_Other(lChannel, dwPTZCommand, dwStop int) (bool, error)

	// GetChannelPTZ 获取指定通道的PTZ位置信息
	// 获取云台当前的水平角度、垂直角度和变焦倍数
	// 参数：
	//   - channel: 通道号，从1开始
	GetChannelPTZ(channel int)

	// RealPlay_V40 启动实时视频预览（使用V40版本接口）
	// 参数：
	//   - ChannelId: 通道号，从1开始
	//   - receiver: 数据接收器，用于接收视频流数据
	// 返回值：
	//   - int: 预览句柄，大于0表示成功
	//   - error: 错误信息，成功时为nil
	RealPlay_V40(ChannelId int, receiver *Receiver) (int, error)

	// StopRealPlay 停止实时视频预览
	// 停止 RealPlay_V40 启动的视频预览
	StopRealPlay()
}

// DeviceInfo 设备信息结构体
// 包含设备的连接信息和基本属性
type DeviceInfo struct {
	// IP 设备IP地址
	// 例如："192.168.1.64"
	IP string

	// Port 设备端口号
	// 默认端口：8000
	Port int

	// UserName 登录用户名
	// 默认用户名通常是 "admin"
	UserName string

	// Password 登录密码
	// 设备的登录密码
	Password string

	// DeviceID 设备序列号
	// 设备的唯一标识符，由设备自动生成
	DeviceID string

	// DeviceName 设备名称
	// 设备的自定义名称，如 "前门摄像头"、"1号NVR" 等
	DeviceName string

	// ByChanNum 通道数量
	// 设备支持的通道总数，单个摄像头通常为1，NVR/DVR可能有多个
	ByChanNum int
}
