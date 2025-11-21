package pkg

import (
	"fmt"
	"io"
)

// ChanReader 将 channel 转换为 io.Reader
type ChanReader chan []byte

func (r ChanReader) Read(buf []byte) (n int, err error) {
	b, ok := <-r
	if !ok {
		return 0, io.EOF
	}
	copy(buf, b)
	return len(b), nil
}

// Receiver 数据接收器，用于接收设备发送的视频/音频数据
type Receiver struct {
	PSMouth  chan []byte // PS 流数据通道
	psBuffer []byte      // PS 数据缓冲区
}

// Start 启动接收器
func (r *Receiver) Start() error {
	r.PSMouth = make(chan []byte, 500) // 创建 PS 数据通道
	return nil
}

// Stop 停止接收器
func (r *Receiver) Stop() {
	if r.PSMouth != nil {
		close(r.PSMouth)
	}
}

// ReadPSData 读取 PS 流数据
func (r *Receiver) ReadPSData(data []byte) error {
	// 将数据添加到缓冲区
	r.psBuffer = append(r.psBuffer, data...)

	// 处理缓冲区中的完整 PS 包
	for {
		syncedData, remaining := r.extractSynchronizedPSData(r.psBuffer)
		if syncedData == nil {
			// 没有找到完整的PS包，保留剩余数据等待更多数据
			r.psBuffer = remaining
			break
		}

		// 发送同步后的PS数据到处理通道
		select {
		case r.PSMouth <- syncedData:
			// 成功发送数据到PS处理通道
		default:
			// 通道满了，跳过这个数据包
			fmt.Printf("PS通道满了，跳过数据包，当前缓冲区大小: %d/%d\n", len(r.PSMouth), cap(r.PSMouth))
		}

		// 更新缓冲区为剩余数据
		r.psBuffer = remaining
	}

	return nil
}

// PS 流起始码
const (
	StartCodePS       = 0x000001BA // PS 包起始码
	StartCodeSYS      = 0x000001BB // 系统头起始码
	StartCodeMAP      = 0x000001BC // 节目流映射
	StartCodeVideo    = 0x000001E0 // 视频流
	StartCodeVideo1   = 0x000001E1 // 视频流1
	StartCodeVideo2   = 0x000001E2 // 视频流2
	StartCodeAudio    = 0x000001C0 // 音频流
	PrivateStreamCode = 0x000001BD // 私有流
)

// extractSynchronizedPSData 从缓冲区中提取同步的 PS 数据包
func (r *Receiver) extractSynchronizedPSData(buffer []byte) ([]byte, []byte) {
	if len(buffer) < 4 {
		return nil, buffer // 数据不足，返回所有数据等待更多
	}

	// 寻找 PS 起始码
	startIndex := -1
	for i := 0; i <= len(buffer)-4; i++ {
		if buffer[i] == 0x00 && buffer[i+1] == 0x00 && buffer[i+2] == 0x01 {
			// 检查第四个字节是否为有效的 PS 起始码
			startCode := uint32(buffer[i])<<24 | uint32(buffer[i+1])<<16 |
				uint32(buffer[i+2])<<8 | uint32(buffer[i+3])

			switch startCode {
			case StartCodePS, StartCodeVideo, StartCodeVideo1,
				StartCodeVideo2, StartCodeAudio, StartCodeMAP,
				StartCodeSYS, PrivateStreamCode:
				startIndex = i
				goto found
			}
		}
	}

found:
	if startIndex == -1 {
		// 没有找到有效起始码
		if len(buffer) > 3 {
			// 保留最后 3 个字节，丢弃其余数据
			return nil, buffer[len(buffer)-3:]
		}
		return nil, buffer
	}

	// 寻找下一个起始码来确定当前包的结束位置
	nextStartIndex := -1
	for i := startIndex + 4; i <= len(buffer)-4; i++ {
		if buffer[i] == 0x00 && buffer[i+1] == 0x00 && buffer[i+2] == 0x01 {
			startCode := uint32(buffer[i])<<24 | uint32(buffer[i+1])<<16 |
				uint32(buffer[i+2])<<8 | uint32(buffer[i+3])

			switch startCode {
			case StartCodePS, StartCodeVideo, StartCodeVideo1,
				StartCodeVideo2, StartCodeAudio, StartCodeMAP,
				StartCodeSYS, PrivateStreamCode:
				nextStartIndex = i
				goto nextFound
			}
		}
	}

nextFound:
	if nextStartIndex == -1 {
		// 没有找到下一个起始码，返回从当前起始码到缓冲区末尾的所有数据
		return buffer[startIndex:], nil
	}

	// 返回从当前起始码到下一个起始码之间的数据
	return buffer[startIndex:nextStartIndex], buffer[nextStartIndex:]
}
