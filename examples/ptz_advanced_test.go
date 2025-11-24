package examples

import (
	"testing"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

// TestPTZAdvanced PTZ高级控制示例
func TestPTZAdvanced(t *testing.T) {
	t.Log("========================================")
	t.Log("海康威视 SDK - PTZ高级控制示例")
	t.Log("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64",
		Port:     8000,
		Username: "admin",
		Password: "asdf234.",
	}

	// 登录设备
	session, err := auth.LoginV30(cred)
	if err != nil {
		t.Skipf("登录失败: %v", err)
		return
	}
	t.Logf("登录成功 (ID: %d)", session.LoginID)
	defer auth.Logout(session.LoginID)
	defer auth.Cleanup()

	channel := 1

	// ==================== 云台移动控制 ====================
	t.Log("\n【云台移动控制】")
	demonstrateMovementTest(t, session.LoginID, channel)

	// ==================== 相机控制 ====================
	t.Log("\n【相机控制】")
	demonstrateCameraTest(t, session.LoginID, channel)

	// ==================== 辅助设备控制 ====================
	t.Log("\n【辅助设备控制】")
	demonstrateAuxiliaryTest(t, session.LoginID, channel)

	t.Log("\n========================================")
	t.Log("示例完成!")
	t.Log("========================================")
}

// 云台移动控制演示
func demonstrateMovementTest(t *testing.T, loginID int, channel int) {
	// 创建PTZ控制器
	ctrl := ptz.NewController(loginID, channel)

	t.Log("\n[1] 基础方向移动（自动控制时长）:")

	// 向右移动2秒
	t.Log("  • 向右移动2秒...")
	if err := ctrl.Right(5, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	// 向上移动2秒
	t.Log("  • 向上移动2秒...")
	if err := ctrl.Up(5, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	t.Log("\n[2] 组合方向移动:")

	// 右上斜向移动
	t.Log("  • 右上斜向移动2秒...")
	if err := ctrl.UpRight(4, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	t.Log("\n[3] 手动控制（自己控制开始和停止）:")

	// 手动控制左转
	t.Log("  • 开始左转...")
	if err := ctrl.StartLeft(5); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		time.Sleep(2 * time.Second)
		t.Log("  • 停止左转...")
		if err := ctrl.StopLeft(); err != nil {
			t.Logf("    ✗ 失败: %v", err)
		}
	}

	t.Log("\n[4] 自动扫描:")

	// 启动自动扫描
	t.Log("  • 启动自动扫描...")
	if err := ctrl.AutoScan(3); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		time.Sleep(3 * time.Second) // 测试环境缩短时间
		t.Log("  • 停止自动扫描...")
		if err := ctrl.StopAutoScan(); err != nil {
			t.Logf("    ✗ 失败: %v", err)
		}
	}
}

// 相机控制演示
func demonstrateCameraTest(t *testing.T, loginID int, channel int) {
	// 创建PTZ控制器
	ctrl := ptz.NewController(loginID, channel)

	t.Log("\n[1] 焦距控制（自动控制时长）:")

	// 焦距放大
	t.Log("  • 焦距放大（拉近）1秒...")
	if err := ctrl.ZoomIn(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// 焦距缩小
	t.Log("  • 焦距缩小（拉远）1秒...")
	if err := ctrl.ZoomOut(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	t.Log("\n[2] 焦点控制（手动开始/停止）:")

	// 焦点前调 - 手动控制
	t.Log("  • 开始焦点前调...")
	if err := ctrl.StartFocusNear(); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		time.Sleep(1 * time.Second)
		t.Log("  • 停止焦点前调...")
		if err := ctrl.StopFocusNear(); err != nil {
			t.Logf("    ✗ 失败: %v", err)
		}
	}

	t.Log("\n[3] 光圈控制:")

	// 光圈扩大
	t.Log("  • 光圈扩大（变亮）1秒...")
	if err := ctrl.IrisOpen(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	// 光圈缩小
	t.Log("  • 光圈缩小（变暗）1秒...")
	if err := ctrl.IrisClose(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	}
}

// 辅助设备控制演示
func demonstrateAuxiliaryTest(t *testing.T, loginID int, channel int) {
	// 创建PTZ控制器
	ctrl := ptz.NewController(loginID, channel)

	t.Log("\n[1] 灯光控制:")

	// 开启灯光
	t.Log("  • 开启灯光...")
	if err := ctrl.LightOn(); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		time.Sleep(2 * time.Second)

		// 关闭灯光
		t.Log("  • 关闭灯光...")
		if err := ctrl.LightOff(); err != nil {
			t.Logf("    ✗ 失败: %v", err)
		}
	}

	t.Log("\n[2] 雨刷控制:")

	// 开启雨刷
	t.Log("  • 开启雨刷...")
	if err := ctrl.WiperOn(); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		time.Sleep(2 * time.Second)

		// 关闭雨刷
		t.Log("  • 关闭雨刷...")
		if err := ctrl.WiperOff(); err != nil {
			t.Logf("    ✗ 失败: %v", err)
		}
	}

	t.Log("\n💡 说明:")
	t.Log("  • 辅助设备功能需要硬件支持")
	t.Log("  • 如果设备不支持某些功能，会返回错误码23（不支持该操作）")
}
