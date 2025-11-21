package pkg

// 音频处理相关功能
// 本文件提供音频数据的基本处理功能

// AudioCodec 音频编码格式
type AudioCodec int

const (
	CodecAAC AudioCodec = iota
	CodecG711A
	CodecG711U
	CodecPCM
)

// AudioFrame 音频帧数据
type AudioFrame struct {
	Data      []byte     // 帧数据
	Timestamp uint32     // 时间戳
	Codec     AudioCodec // 编码格式
}

// String 返回编码格式的字符串表示
func (c AudioCodec) String() string {
	switch c {
	case CodecAAC:
		return "AAC"
	case CodecG711A:
		return "G.711 A-law"
	case CodecG711U:
		return "G.711 μ-law"
	case CodecPCM:
		return "PCM"
	default:
		return "Unknown"
	}
}
