package pkg

// PTZ 命令常量定义
// 这些常量用于 PTZ 控制函数的 dwPTZCommand 参数

const (
	// ========== 辅助设备控制 ==========
	// LIGHT_PWRON 接通灯光电源
	LIGHT_PWRON = 2
	// WIPER_PWRON 接通雨刷开关
	WIPER_PWRON = 3
	// FAN_PWRON 接通风扇开关
	FAN_PWRON = 4
	// HEATER_PWRON 接通加热器开关
	HEATER_PWRON = 5
	// AUX_PWRON1 接通辅助设备开关1
	AUX_PWRON1 = 6
	// AUX_PWRON2 接通辅助设备开关2
	AUX_PWRON2 = 7

	// ========== 预置点操作 ==========
	// SET_PRESET 设置预置点
	SET_PRESET = 8
	// CLE_PRESET 清除预置点
	CLE_PRESET = 9
	// GOTO_PRESET 快球转到预置点
	GOTO_PRESET = 39

	// ========== 焦距控制（需要带速度的控制方法） ==========
	// ZOOM_IN 焦距变大(倍率变大)
	ZOOM_IN = 11
	// ZOOM_OUT 焦距变小(倍率变小)
	ZOOM_OUT = 12

	// ========== 焦点控制 ==========
	// FOCUS_NEAR 焦点前调
	FOCUS_NEAR = 13
	// FOCUS_FAR 焦点后调
	FOCUS_FAR = 14

	// ========== 光圈控制 ==========
	// IRIS_OPEN 光圈扩大
	IRIS_OPEN = 15
	// IRIS_CLOSE 光圈缩小
	IRIS_CLOSE = 16

	// ========== 基本云台移动（需要带速度的控制方法） ==========
	// TILT_UP 云台上仰
	TILT_UP = 21
	// TILT_DOWN 云台下俯
	TILT_DOWN = 22
	// PAN_LEFT 云台左转
	PAN_LEFT = 23
	// PAN_RIGHT 云台右转
	PAN_RIGHT = 24

	// ========== 组合移动（需要带速度的控制方法） ==========
	// UP_LEFT 云台上仰和左转
	UP_LEFT = 25
	// UP_RIGHT 云台上仰和右转
	UP_RIGHT = 26
	// DOWN_LEFT 云台下俯和左转
	DOWN_LEFT = 27
	// DOWN_RIGHT 云台下俯和右转
	DOWN_RIGHT = 28

	// ========== 自动扫描 ==========
	// PAN_AUTO 云台左右自动扫描
	PAN_AUTO = 29
	// PAN_CIRCLE 云台自动圆周扫描
	PAN_CIRCLE = 50
	// DRAG_PTZ 拖动PTZ
	DRAG_PTZ = 51

	// ========== 巡航控制 ==========
	// FILL_PRE_SEQ 将预置点加入巡航序列
	FILL_PRE_SEQ = 30
	// SET_SEQ_DWELL 设置巡航点停顿时间
	SET_SEQ_DWELL = 31
	// SET_SEQ_SPEED 设置巡航速度
	SET_SEQ_SPEED = 32
	// CLE_PRE_SEQ 将预置点从巡航序列中删除
	CLE_PRE_SEQ = 33
	// RUN_SEQ 开始巡航
	RUN_SEQ = 37
	// STOP_SEQ 停止巡航
	STOP_SEQ = 38
	// DEL_SEQ 删除巡航路径
	DEL_SEQ = 43

	// ========== 轨迹控制 ==========
	// STA_MEM_CRUISE 开始记录轨迹
	STA_MEM_CRUISE = 34
	// STO_MEM_CRUISE 停止记录轨迹
	STO_MEM_CRUISE = 35
	// RUN_CRUISE 开始轨迹
	RUN_CRUISE = 36
	// STOP_CRUISE 停止轨迹
	STOP_CRUISE = 44
	// DELETE_CRUISE 删除单条轨迹
	DELETE_CRUISE = 45
	// DELETE_ALL_CRUISE 删除所有轨迹
	DELETE_ALL_CRUISE = 46

	// ========== 其他控制 ==========
	// LINEAR_SCAN 区域扫描
	LINEAR_SCAN = 52
	// CLE_ALL_PRESET 预置点全部清除
	CLE_ALL_PRESET = 53
	// CLE_ALL_SEQ 巡航全部清除
	CLE_ALL_SEQ = 54
	// CLE_ALL_CRUISE 轨迹全部清除
	CLE_ALL_CRUISE = 55
	// POPUP_MENU 显示操作菜单
	POPUP_MENU = 56

	// ========== 组合控制（移动+变焦） ==========
	// TILT_DOWN_ZOOM_IN 云台下俯+焦距变大
	TILT_DOWN_ZOOM_IN = 58
	// TILT_DOWN_ZOOM_OUT 云台下俯+焦距变小
	TILT_DOWN_ZOOM_OUT = 59
	// PAN_LEFT_ZOOM_IN 云台左转+焦距变大
	PAN_LEFT_ZOOM_IN = 60
	// PAN_LEFT_ZOOM_OUT 云台左转+焦距变小
	PAN_LEFT_ZOOM_OUT = 61
	// PAN_RIGHT_ZOOM_IN 云台右转+焦距变大
	PAN_RIGHT_ZOOM_IN = 62
	// PAN_RIGHT_ZOOM_OUT 云台右转+焦距变小
	PAN_RIGHT_ZOOM_OUT = 63
	// UP_LEFT_ZOOM_IN 云台上仰左转+焦距变大
	UP_LEFT_ZOOM_IN = 64
	// UP_LEFT_ZOOM_OUT 云台上仰左转+焦距变小
	UP_LEFT_ZOOM_OUT = 65
	// UP_RIGHT_ZOOM_IN 云台上仰右转+焦距变大
	UP_RIGHT_ZOOM_IN = 66
	// UP_RIGHT_ZOOM_OUT 云台上仰右转+焦距变小
	UP_RIGHT_ZOOM_OUT = 67
	// DOWN_LEFT_ZOOM_IN 云台下俯左转+焦距变大
	DOWN_LEFT_ZOOM_IN = 68
	// DOWN_LEFT_ZOOM_OUT 云台下俯左转+焦距变小
	DOWN_LEFT_ZOOM_OUT = 69
	// DOWN_RIGHT_ZOOM_IN 云台下俯右转+焦距变大
	DOWN_RIGHT_ZOOM_IN = 70
	// DOWN_RIGHT_ZOOM_OUT 云台下俯右转+焦距变小
	DOWN_RIGHT_ZOOM_OUT = 71
	// TILT_UP_ZOOM_IN 云台上仰+焦距变大
	TILT_UP_ZOOM_IN = 72
	// TILT_UP_ZOOM_OUT 云台上仰+焦距变小
	TILT_UP_ZOOM_OUT = 73
)

// PTZ 控制参数说明
const (
	// PTZ_STOP 停止动作（dwStop 参数）
	PTZ_STOP = 1
	// PTZ_START 开始动作（dwStop 参数）
	PTZ_START = 0

	// PTZ_SPEED_MIN 最小速度
	PTZ_SPEED_MIN = 0
	// PTZ_SPEED_MAX 最大速度
	PTZ_SPEED_MAX = 7
	// PTZ_SPEED_DEFAULT 默认速度
	PTZ_SPEED_DEFAULT = 4
)

// GetPTZCommandName 获取 PTZ 命令的名称（用于调试）
func GetPTZCommandName(command int) string {
	names := map[int]string{
		LIGHT_PWRON:         "接通灯光电源",
		WIPER_PWRON:         "接通雨刷开关",
		FAN_PWRON:           "接通风扇开关",
		HEATER_PWRON:        "接通加热器开关",
		AUX_PWRON1:          "接通辅助设备1",
		AUX_PWRON2:          "接通辅助设备2",
		SET_PRESET:          "设置预置点",
		CLE_PRESET:          "清除预置点",
		GOTO_PRESET:         "转到预置点",
		ZOOM_IN:             "焦距放大",
		ZOOM_OUT:            "焦距缩小",
		FOCUS_NEAR:          "焦点前调",
		FOCUS_FAR:           "焦点后调",
		IRIS_OPEN:           "光圈扩大",
		IRIS_CLOSE:          "光圈缩小",
		TILT_UP:             "云台上仰",
		TILT_DOWN:           "云台下俯",
		PAN_LEFT:            "云台左转",
		PAN_RIGHT:           "云台右转",
		UP_LEFT:             "上仰左转",
		UP_RIGHT:            "上仰右转",
		DOWN_LEFT:           "下俯左转",
		DOWN_RIGHT:          "下俯右转",
		PAN_AUTO:            "左右自动扫描",
		PAN_CIRCLE:          "圆周扫描",
		DRAG_PTZ:            "拖动PTZ",
		FILL_PRE_SEQ:        "加入巡航序列",
		SET_SEQ_DWELL:       "设置巡航停顿时间",
		SET_SEQ_SPEED:       "设置巡航速度",
		CLE_PRE_SEQ:         "从巡航序列删除",
		RUN_SEQ:             "开始巡航",
		STOP_SEQ:            "停止巡航",
		DEL_SEQ:             "删除巡航路径",
		STA_MEM_CRUISE:      "开始记录轨迹",
		STO_MEM_CRUISE:      "停止记录轨迹",
		RUN_CRUISE:          "开始轨迹",
		STOP_CRUISE:         "停止轨迹",
		DELETE_CRUISE:       "删除单条轨迹",
		DELETE_ALL_CRUISE:   "删除所有轨迹",
		LINEAR_SCAN:         "区域扫描",
		CLE_ALL_PRESET:      "清除所有预置点",
		CLE_ALL_SEQ:         "清除所有巡航",
		CLE_ALL_CRUISE:      "清除所有轨迹",
		POPUP_MENU:          "显示操作菜单",
		TILT_DOWN_ZOOM_IN:   "下俯+放大",
		TILT_DOWN_ZOOM_OUT:  "下俯+缩小",
		PAN_LEFT_ZOOM_IN:    "左转+放大",
		PAN_LEFT_ZOOM_OUT:   "左转+缩小",
		PAN_RIGHT_ZOOM_IN:   "右转+放大",
		PAN_RIGHT_ZOOM_OUT:  "右转+缩小",
		UP_LEFT_ZOOM_IN:     "上仰左转+放大",
		UP_LEFT_ZOOM_OUT:    "上仰左转+缩小",
		UP_RIGHT_ZOOM_IN:    "上仰右转+放大",
		UP_RIGHT_ZOOM_OUT:   "上仰右转+缩小",
		DOWN_LEFT_ZOOM_IN:   "下俯左转+放大",
		DOWN_LEFT_ZOOM_OUT:  "下俯左转+缩小",
		DOWN_RIGHT_ZOOM_IN:  "下俯右转+放大",
		DOWN_RIGHT_ZOOM_OUT: "下俯右转+缩小",
		TILT_UP_ZOOM_IN:     "上仰+放大",
		TILT_UP_ZOOM_OUT:    "上仰+缩小",
	}

	if name, ok := names[command]; ok {
		return name
	}
	return "未知命令"
}
