package core

import (
	"fmt"
	"io"
	"log"
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
// 支持PS流的接收、缓冲和同步处理
type Receiver struct {
	PSMouth  chan []byte // PS 流数据通道（对外接口）
	psBuffer []byte      // PS 数据缓冲区（内部使用）
	running  bool        // 运行状态标志
}

// Start 启动接收器
// 初始化PS流数据通道，准备接收数据
// 返回值：
//   - error: 错误信息，成功时为nil
func (r *Receiver) Start() error {
	if r.running {
		return fmt.Errorf("接收器已经在运行")
	}

	// 创建 PS 数据通道，缓冲区大小500
	// 如果数据产生速度过快，可能需要增大缓冲区
	r.PSMouth = make(chan []byte, 500)
	r.psBuffer = make([]byte, 0, 1024*1024) // 预分配1MB缓冲区
	r.running = true

	log.Println("✓ 数据接收器已启动")
	return nil
}

// Stop 停止接收器
// 关闭数据通道，清理资源
func (r *Receiver) Stop() {
	if !r.running {
		return
	}

	if r.PSMouth != nil {
		close(r.PSMouth)
		r.PSMouth = nil
	}

	// 清空缓冲区
	r.psBuffer = nil
	r.running = false

	log.Println("✓ 数据接收器已停止")
}

// ReadPSData 读取 PS 流数据
// 从设备接收原始数据，提取PS包并发送到处理通道
// 参数：
//   - data: 从设备接收的原始数据
//
// 返回值：
//   - error: 处理错误信息
func (r *Receiver) ReadPSData(data []byte) error {
	if !r.running {
		return fmt.Errorf("接收器未运行")
	}

	if len(data) == 0 {
		return nil // 空数据不是错误
	}

	// 将数据添加到缓冲区
	r.psBuffer = append(r.psBuffer, data...)

	// 防止缓冲区无限增长
	const maxBufferSize = 10 * 1024 * 1024 // 10MB
	if len(r.psBuffer) > maxBufferSize {
		log.Printf("警告: PS缓冲区过大(%d bytes)，清空缓冲区", len(r.psBuffer))
		r.psBuffer = r.psBuffer[:0] // 清空但保留容量
		return fmt.Errorf("PS缓冲区溢出")
	}

	// 处理缓冲区中的完整 PS 包
	packetsProcessed := 0
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
			packetsProcessed++
		default:
			// 通道满了，跳过这个数据包
			log.Printf("警告: PS通道满（%d/%d），跳过数据包", len(r.PSMouth), cap(r.PSMouth))
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
