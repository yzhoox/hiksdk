package consts

// MAX_SERIALNO_LEN 序列号最大长度（与C头文件中定义一致）
const MAX_SERIALNO_LEN = 48

// 报警类型常量定义
const (
	COMM_ALARM_RULE = 0x1101 // 行为分析报警
	COMM_ALARM_V30  = 0x4000 // V30报警信息（移动侦测、视频丢失、遮挡、IO等）
	COMM_ALARM_V40  = 0x4007 // V40报警信息（扩展报警类型）
)
