package ptz

/*
#include <stdio.h>
#include <stdlib.h>
#include "../hiksdk_wrapper.h"
*/
import "C"
import (
	"fmt"
	"log"
	"time"

	"github.com/samsaralc/hiksdk/core"
)

// ==================== 云台移动命令常量（来自官方文档表 5.10）====================
const (
	// ========== 基本云台移动 ==========
	// TILT_UP 云台上仰
	TILT_UP = 21
	// TILT_DOWN 云台下俯
	TILT_DOWN = 22
	// PAN_LEFT 云台左转
	PAN_LEFT = 23
	// PAN_RIGHT 云台右转
	PAN_RIGHT = 24

	// ========== 组合移动 ==========
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
)

// ==================== 相机控制命令常量（来自官方文档表 5.10）====================
const (
	// ========== 焦距控制 ==========
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
)

// ==================== 辅助设备命令常量（来自官方文档表 5.10）====================
const (
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
)

// 控制参数常量（来自官方文档）
const (
	// PTZ_START 开始动作
	PTZ_START = 0
	// PTZ_STOP 停止动作
	PTZ_STOP = 1

	// 速度范围：1-7（根据官方文档 5.6.3 和 5.6.4）
	MinSpeed     = 1
	MaxSpeed     = 7
	DefaultSpeed = 4
)

// ==================== Controller 统一的PTZ控制器 ====================

// Controller PTZ统一控制器
// 封装云台移动、相机控制、辅助设备控制的所有操作
type Controller struct {
	userID  int // 登录句柄
	channel int // 通道号
}

// NewController 创建PTZ控制器
// 参数：
//   - userID: 登录句柄
//   - channel: 通道号
func NewController(userID int, channel int) *Controller {
	return &Controller{
		userID:  userID,
		channel: channel,
	}
}

// ==================== 云台移动控制（带持续时间，自动停止）====================

// Up 云台上仰（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) Up(speed int, duration time.Duration) error {
	return c.move(TILT_UP, speed, duration)
}

// Down 云台下俯（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) Down(speed int, duration time.Duration) error {
	return c.move(TILT_DOWN, speed, duration)
}

// Left 云台左转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) Left(speed int, duration time.Duration) error {
	return c.move(PAN_LEFT, speed, duration)
}

// Right 云台右转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) Right(speed int, duration time.Duration) error {
	return c.move(PAN_RIGHT, speed, duration)
}

// UpLeft 云台上仰并左转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) UpLeft(speed int, duration time.Duration) error {
	return c.move(UP_LEFT, speed, duration)
}

// UpRight 云台上仰并右转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) UpRight(speed int, duration time.Duration) error {
	return c.move(UP_RIGHT, speed, duration)
}

// DownLeft 云台下俯并左转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) DownLeft(speed int, duration time.Duration) error {
	return c.move(DOWN_LEFT, speed, duration)
}

// DownRight 云台下俯并右转（自动控制时长后停止）
// 参数：
//   - speed: 速度（1-7）
//   - duration: 持续时间
func (c *Controller) DownRight(speed int, duration time.Duration) error {
	return c.move(DOWN_RIGHT, speed, duration)
}

// AutoScan 云台左右自动扫描（持续扫描，需要手动停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) AutoScan(speed int) error {
	if err := c.validateSpeed(speed); err != nil {
		return err
	}

	// 开始自动扫描
	if err := c.controlWithSpeed(PAN_AUTO, PTZ_START, speed); err != nil {
		return fmt.Errorf("启动自动扫描失败: %w", err)
	}

	log.Printf("✓ 启动自动扫描（通道%d，速度%d）", c.channel, speed)
	return nil
}

// StopAutoScan 停止自动扫描
func (c *Controller) StopAutoScan() error {
	if err := c.controlWithSpeed(PAN_AUTO, PTZ_STOP, DefaultSpeed); err != nil {
		return fmt.Errorf("停止自动扫描失败: %w", err)
	}

	log.Printf("✓ 停止自动扫描（通道%d）", c.channel)
	return nil
}

// ==================== 云台移动控制（手动开始/停止）====================

// StartUp 开始云台上仰（需手动调用StopUp停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartUp(speed int) error {
	return c.startMove(TILT_UP, speed, "上仰")
}

// StopUp 停止云台上仰
func (c *Controller) StopUp() error {
	return c.stopMove(TILT_UP, "上仰")
}

// StartDown 开始云台下俯（需手动调用StopDown停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartDown(speed int) error {
	return c.startMove(TILT_DOWN, speed, "下俯")
}

// StopDown 停止云台下俯
func (c *Controller) StopDown() error {
	return c.stopMove(TILT_DOWN, "下俯")
}

// StartLeft 开始云台左转（需手动调用StopLeft停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartLeft(speed int) error {
	return c.startMove(PAN_LEFT, speed, "左转")
}

// StopLeft 停止云台左转
func (c *Controller) StopLeft() error {
	return c.stopMove(PAN_LEFT, "左转")
}

// StartRight 开始云台右转（需手动调用StopRight停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartRight(speed int) error {
	return c.startMove(PAN_RIGHT, speed, "右转")
}

// StopRight 停止云台右转
func (c *Controller) StopRight() error {
	return c.stopMove(PAN_RIGHT, "右转")
}

// StartUpLeft 开始云台上仰并左转（需手动调用StopUpLeft停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartUpLeft(speed int) error {
	return c.startMove(UP_LEFT, speed, "上仰左转")
}

// StopUpLeft 停止云台上仰并左转
func (c *Controller) StopUpLeft() error {
	return c.stopMove(UP_LEFT, "上仰左转")
}

// StartUpRight 开始云台上仰并右转（需手动调用StopUpRight停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartUpRight(speed int) error {
	return c.startMove(UP_RIGHT, speed, "上仰右转")
}

// StopUpRight 停止云台上仰并右转
func (c *Controller) StopUpRight() error {
	return c.stopMove(UP_RIGHT, "上仰右转")
}

// StartDownLeft 开始云台下俯并左转（需手动调用StopDownLeft停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartDownLeft(speed int) error {
	return c.startMove(DOWN_LEFT, speed, "下俯左转")
}

// StopDownLeft 停止云台下俯并左转
func (c *Controller) StopDownLeft() error {
	return c.stopMove(DOWN_LEFT, "下俯左转")
}

// StartDownRight 开始云台下俯并右转（需手动调用StopDownRight停止）
// 参数：
//   - speed: 速度（1-7）
func (c *Controller) StartDownRight(speed int) error {
	return c.startMove(DOWN_RIGHT, speed, "下俯右转")
}

// StopDownRight 停止云台下俯并右转
func (c *Controller) StopDownRight() error {
	return c.stopMove(DOWN_RIGHT, "下俯右转")
}

// ==================== 相机控制（带持续时间，自动停止）====================

// ZoomIn 焦距放大（拉近，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) ZoomIn(duration time.Duration) error {
	return c.adjustCamera(ZOOM_IN, duration, "焦距放大")
}

// ZoomOut 焦距缩小（拉远，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) ZoomOut(duration time.Duration) error {
	return c.adjustCamera(ZOOM_OUT, duration, "焦距缩小")
}

// FocusNear 焦点前调（聚焦近处，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) FocusNear(duration time.Duration) error {
	return c.adjustCamera(FOCUS_NEAR, duration, "焦点前调")
}

// FocusFar 焦点后调（聚焦远处，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) FocusFar(duration time.Duration) error {
	return c.adjustCamera(FOCUS_FAR, duration, "焦点后调")
}

// IrisOpen 光圈扩大（变亮，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) IrisOpen(duration time.Duration) error {
	return c.adjustCamera(IRIS_OPEN, duration, "光圈扩大")
}

// IrisClose 光圈缩小（变暗，自动控制时长后停止）
// 参数：
//   - duration: 持续时间
func (c *Controller) IrisClose(duration time.Duration) error {
	return c.adjustCamera(IRIS_CLOSE, duration, "光圈缩小")
}

// ==================== 相机控制（手动开始/停止）====================

// StartZoomIn 开始焦距放大（需手动调用StopZoomIn停止）
func (c *Controller) StartZoomIn() error {
	return c.startCamera(ZOOM_IN, "焦距放大")
}

// StopZoomIn 停止焦距放大
func (c *Controller) StopZoomIn() error {
	return c.stopCamera(ZOOM_IN, "焦距放大")
}

// StartZoomOut 开始焦距缩小（需手动调用StopZoomOut停止）
func (c *Controller) StartZoomOut() error {
	return c.startCamera(ZOOM_OUT, "焦距缩小")
}

// StopZoomOut 停止焦距缩小
func (c *Controller) StopZoomOut() error {
	return c.stopCamera(ZOOM_OUT, "焦距缩小")
}

// StartFocusNear 开始焦点前调（需手动调用StopFocusNear停止）
func (c *Controller) StartFocusNear() error {
	return c.startCamera(FOCUS_NEAR, "焦点前调")
}

// StopFocusNear 停止焦点前调
func (c *Controller) StopFocusNear() error {
	return c.stopCamera(FOCUS_NEAR, "焦点前调")
}

// StartFocusFar 开始焦点后调（需手动调用StopFocusFar停止）
func (c *Controller) StartFocusFar() error {
	return c.startCamera(FOCUS_FAR, "焦点后调")
}

// StopFocusFar 停止焦点后调
func (c *Controller) StopFocusFar() error {
	return c.stopCamera(FOCUS_FAR, "焦点后调")
}

// StartIrisOpen 开始光圈扩大（需手动调用StopIrisOpen停止）
func (c *Controller) StartIrisOpen() error {
	return c.startCamera(IRIS_OPEN, "光圈扩大")
}

// StopIrisOpen 停止光圈扩大
func (c *Controller) StopIrisOpen() error {
	return c.stopCamera(IRIS_OPEN, "光圈扩大")
}

// StartIrisClose 开始光圈缩小（需手动调用StopIrisClose停止）
func (c *Controller) StartIrisClose() error {
	return c.startCamera(IRIS_CLOSE, "光圈缩小")
}

// StopIrisClose 停止光圈缩小
func (c *Controller) StopIrisClose() error {
	return c.stopCamera(IRIS_CLOSE, "光圈缩小")
}

// ==================== 辅助设备控制 ====================

// LightOn 接通灯光电源
func (c *Controller) LightOn() error {
	return c.switchDevice(LIGHT_PWRON, true, "灯光")
}

// LightOff 关闭灯光电源
func (c *Controller) LightOff() error {
	return c.switchDevice(LIGHT_PWRON, false, "灯光")
}

// WiperOn 接通雨刷
func (c *Controller) WiperOn() error {
	return c.switchDevice(WIPER_PWRON, true, "雨刷")
}

// WiperOff 关闭雨刷
func (c *Controller) WiperOff() error {
	return c.switchDevice(WIPER_PWRON, false, "雨刷")
}

// FanOn 接通风扇
func (c *Controller) FanOn() error {
	return c.switchDevice(FAN_PWRON, true, "风扇")
}

// FanOff 关闭风扇
func (c *Controller) FanOff() error {
	return c.switchDevice(FAN_PWRON, false, "风扇")
}

// HeaterOn 接通加热器
func (c *Controller) HeaterOn() error {
	return c.switchDevice(HEATER_PWRON, true, "加热器")
}

// HeaterOff 关闭加热器
func (c *Controller) HeaterOff() error {
	return c.switchDevice(HEATER_PWRON, false, "加热器")
}

// AuxDevice1On 接通辅助设备1
func (c *Controller) AuxDevice1On() error {
	return c.switchDevice(AUX_PWRON1, true, "辅助设备1")
}

// AuxDevice1Off 关闭辅助设备1
func (c *Controller) AuxDevice1Off() error {
	return c.switchDevice(AUX_PWRON1, false, "辅助设备1")
}

// AuxDevice2On 接通辅助设备2
func (c *Controller) AuxDevice2On() error {
	return c.switchDevice(AUX_PWRON2, true, "辅助设备2")
}

// AuxDevice2Off 关闭辅助设备2
func (c *Controller) AuxDevice2Off() error {
	return c.switchDevice(AUX_PWRON2, false, "辅助设备2")
}

// ==================== 内部实现函数 ====================

// move 云台移动（带速度和时长）
func (c *Controller) move(cmd, speed int, duration time.Duration) error {
	// 验证速度
	if err := c.validateSpeed(speed); err != nil {
		return err
	}

	// 开始移动
	if err := c.controlWithSpeed(cmd, PTZ_START, speed); err != nil {
		return err
	}

	// 等待指定时间
	time.Sleep(duration)

	// 停止移动
	if err := c.controlWithSpeed(cmd, PTZ_STOP, speed); err != nil {
		return err
	}

	return nil
}

// startMove 开始云台移动（手动控制）
func (c *Controller) startMove(cmd, speed int, actionName string) error {
	if err := c.validateSpeed(speed); err != nil {
		return err
	}

	if err := c.controlWithSpeed(cmd, PTZ_START, speed); err != nil {
		return fmt.Errorf("开始%s失败: %w", actionName, err)
	}

	log.Printf("✓ 开始%s（通道%d，速度%d）", actionName, c.channel, speed)
	return nil
}

// stopMove 停止云台移动（手动控制）
func (c *Controller) stopMove(cmd int, actionName string) error {
	if err := c.controlWithSpeed(cmd, PTZ_STOP, DefaultSpeed); err != nil {
		return fmt.Errorf("停止%s失败: %w", actionName, err)
	}

	log.Printf("✓ 停止%s（通道%d）", actionName, c.channel)
	return nil
}

// adjustCamera 相机调整（带时长）
func (c *Controller) adjustCamera(cmd int, duration time.Duration, actionName string) error {
	// 开始调整
	if err := c.controlWithSpeed(cmd, PTZ_START, DefaultSpeed); err != nil {
		return fmt.Errorf("%s失败: %w", actionName, err)
	}

	// 等待指定时间
	time.Sleep(duration)

	// 停止调整
	if err := c.controlWithSpeed(cmd, PTZ_STOP, DefaultSpeed); err != nil {
		return fmt.Errorf("停止%s失败: %w", actionName, err)
	}

	log.Printf("✓ %s（通道%d，持续%v）", actionName, c.channel, duration)
	return nil
}

// startCamera 开始相机调整（手动控制）
func (c *Controller) startCamera(cmd int, actionName string) error {
	if err := c.controlWithSpeed(cmd, PTZ_START, DefaultSpeed); err != nil {
		return fmt.Errorf("开始%s失败: %w", actionName, err)
	}

	log.Printf("✓ 开始%s（通道%d）", actionName, c.channel)
	return nil
}

// stopCamera 停止相机调整（手动控制）
func (c *Controller) stopCamera(cmd int, actionName string) error {
	if err := c.controlWithSpeed(cmd, PTZ_STOP, DefaultSpeed); err != nil {
		return fmt.Errorf("停止%s失败: %w", actionName, err)
	}

	log.Printf("✓ 停止%s（通道%d）", actionName, c.channel)
	return nil
}

// switchDevice 辅助设备开关
func (c *Controller) switchDevice(cmd int, turnOn bool, deviceName string) error {
	action := PTZ_START // 0=开启
	actionName := "开启"
	if !turnOn {
		action = PTZ_STOP // 1=关闭
		actionName = "关闭"
	}

	if err := c.controlWithSpeed(cmd, action, DefaultSpeed); err != nil {
		return fmt.Errorf("%s%s失败: %w", actionName, deviceName, err)
	}

	log.Printf("✓ %s%s（通道%d）", actionName, deviceName, c.channel)
	return nil
}

// controlWithSpeed 带速度的云台控制（底层调用）
func (c *Controller) controlWithSpeed(cmd, stop, speed int) error {
	if c.userID < 0 {
		return fmt.Errorf("无效的登录ID：%d", c.userID)
	}

	ret := C.NET_DVR_PTZControlWithSpeed_Other(
		C.LONG(c.userID),
		C.LONG(c.channel),
		C.DWORD(cmd),
		C.DWORD(stop),
		C.DWORD(speed),
	)

	if ret != C.TRUE {
		return core.NewHKError(fmt.Sprintf("PTZ控制[通道:%d 命令:%d]", c.channel, cmd))
	}

	return nil
}

// validateSpeed 验证速度范围
func (c *Controller) validateSpeed(speed int) error {
	if speed < MinSpeed || speed > MaxSpeed {
		return fmt.Errorf("速度超出范围：%d（有效范围：%d-%d）", speed, MinSpeed, MaxSpeed)
	}
	return nil
}
