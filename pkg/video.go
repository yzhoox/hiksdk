package pkg

// 视频处理相关功能
// 本文件提供视频数据的基本处理功能

// VideoCodec 视频编码格式
type VideoCodec int

const (
	CodecH264 VideoCodec = iota
	CodecH265
	CodecMJPEG
)

// VideoFrame 视频帧数据
type VideoFrame struct {
	Data      []byte     // 帧数据
	Timestamp uint32     // 时间戳
	Codec     VideoCodec // 编码格式
	KeyFrame  bool       // 是否为关键帧
}

// String 返回编码格式的字符串表示
func (c VideoCodec) String() string {
	switch c {
	case CodecH264:
		return "H.264"
	case CodecH265:
		return "H.265"
	case CodecMJPEG:
		return "MJPEG"
	default:
		return "Unknown"
	}
}
