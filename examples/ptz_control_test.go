package examples

import (
	"testing"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

// TestPTZControl PTZ云台控制示例
func TestPTZControl(t *testing.T) {
	t.Log("========================================")
	t.Log("海康威视 SDK - PTZ控制示例")
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

	// 选择通道
	channel := 1

	// 创建PTZ控制器（统一控制云台、相机、辅助设备）
	ctrl := ptz.NewController(session.LoginID, channel)

	// 创建预置点管理器
	preset := ptz.NewPresetManager(session.LoginID, channel)

	// ==================== 步骤0: 设置原点预置点 ====================
	t.Log("\n[步骤0] 设置当前位置为原点（预置点1）")
	if err := preset.SetPreset(1); err != nil {
		t.Logf("  ✗ 设置原点失败: %v", err)
		t.Log("  ⚠️  继续测试，但无法回到原点")
	} else {
		t.Log("  ✓ 原点已设置（预置点1）")
	}

	// ==================== 步骤1: 基础方向控制（自动计时）====================
	t.Log("\n[步骤1] 基础方向控制（自动计时）")

	// 向右转动2秒
	t.Log("  • 向右转动2秒（速度7）...")
	if err := ctrl.Right(7, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// 向上转动2秒
	t.Log("  • 向上转动2秒（速度7）...")
	if err := ctrl.Up(7, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// 向左转动2秒
	t.Log("  • 向左转动2秒（速度6）...")
	if err := ctrl.Left(6, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// 向下转动2秒
	t.Log("  • 向下转动2秒（速度6）...")
	if err := ctrl.Down(6, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// ==================== 步骤2: 组合方向控制 ====================
	t.Log("\n[步骤2] 组合方向控制")

	// 右上斜向移动
	t.Log("  • 右上斜向移动2秒（速度5）...")
	if err := ctrl.UpRight(5, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// 左下斜向移动
	t.Log("  • 左下斜向移动2秒（速度5）...")
	if err := ctrl.DownLeft(5, 2*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// ==================== 步骤3: 手动开始/停止控制 ====================
	t.Log("\n[步骤3] 手动开始/停止控制（更灵活）")

	// 开始右转
	t.Log("  • 开始右转（速度4）...")
	if err := ctrl.StartRight(4); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 已开始，等待3秒...")
		time.Sleep(3 * time.Second)

		// 停止右转
		t.Log("  • 停止右转...")
		if err := ctrl.StopRight(); err != nil {
			t.Logf("    ✗ 停止失败: %v", err)
		} else {
			t.Log("    ✓ 已停止")
		}
	}

	// 开始上仰
	t.Log("  • 开始上仰（速度5）...")
	if err := ctrl.StartUp(5); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 已开始，等待2秒...")
		time.Sleep(2 * time.Second)

		// 停止上仰
		t.Log("  • 停止上仰...")
		if err := ctrl.StopUp(); err != nil {
			t.Logf("    ✗ 停止失败: %v", err)
		} else {
			t.Log("    ✓ 已停止")
		}
	}

	// ==================== 步骤4: 相机焦距控制 ====================
	t.Log("\n[步骤4] 相机焦距控制")

	// 焦距放大
	t.Log("  • 焦距放大（拉近）1秒...")
	if err := ctrl.ZoomIn(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	time.Sleep(500 * time.Millisecond)

	// 焦距缩小
	t.Log("  • 焦距缩小（拉远）1秒...")
	if err := ctrl.ZoomOut(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// ==================== 步骤5: 相机焦点控制（手动）====================
	t.Log("\n[步骤5] 相机焦点控制（手动开始/停止）")

	// 开始焦点前调
	t.Log("  • 开始焦点前调（聚焦近处）...")
	if err := ctrl.StartFocusNear(); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 已开始，等待800毫秒...")
		time.Sleep(800 * time.Millisecond)

		// 停止焦点前调
		t.Log("  • 停止焦点前调...")
		if err := ctrl.StopFocusNear(); err != nil {
			t.Logf("    ✗ 停止失败: %v", err)
		} else {
			t.Log("    ✓ 已停止")
		}
	}

	time.Sleep(500 * time.Millisecond)

	// 焦点后调
	t.Log("  • 焦点后调（聚焦远处）800毫秒...")
	if err := ctrl.FocusFar(800 * time.Millisecond); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// ==================== 步骤6: 相机光圈控制 ====================
	t.Log("\n[步骤6] 相机光圈控制")

	// 光圈扩大
	t.Log("  • 光圈扩大（变亮）1秒...")
	if err := ctrl.IrisOpen(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	time.Sleep(500 * time.Millisecond)

	// 光圈缩小
	t.Log("  • 光圈缩小（变暗）1秒...")
	if err := ctrl.IrisClose(1 * time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	// ==================== 步骤7: 预置点操作 ====================
	t.Log("\n[步骤7] 预置点操作")

	// 设置预置点2
	t.Log("  • 设置预置点2（当前位置）...")
	if err := preset.SetPreset(2); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 预置点2已设置")
	}

	// 移动到其他位置
	t.Log("  • 移动到其他位置（左转4秒）...")
	if err := ctrl.Left(4, 4*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	time.Sleep(1 * time.Second)

	// 转到预置点2
	t.Log("  • 转到预置点2...")
	if err := preset.GotoPreset(2); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 正在移动到预置点2...")
		time.Sleep(3 * time.Second)
		t.Log("    ✓ 已到达预置点2")
	}

	// 设置预置点3
	t.Log("  • 设置预置点3（当前位置）...")
	if err := preset.SetPreset(3); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 预置点3已设置")
	}

	// 再次移动到其他位置
	t.Log("  • 移动到其他位置（右转5秒）...")
	if err := ctrl.Right(5, 5*time.Second); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 完成")
	}

	time.Sleep(1 * time.Second)

	// 转到预置点3
	t.Log("  • 转到预置点3...")
	if err := preset.GotoPreset(3); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 正在移动到预置点3...")
		time.Sleep(3 * time.Second)
		t.Log("    ✓ 已到达预置点3")
	}

	// ==================== 步骤8: 自动扫描 ====================
	t.Log("\n[步骤8] 自动扫描测试")

	t.Log("  • 启动自动扫描（速度3）...")
	if err := ctrl.AutoScan(3); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 自动扫描已启动，运行5秒...")
		time.Sleep(5 * time.Second)

		// 停止自动扫描
		t.Log("  • 停止自动扫描...")
		if err := ctrl.StopAutoScan(); err != nil {
			t.Logf("    ✗ 停止失败: %v", err)
		} else {
			t.Log("    ✓ 自动扫描已停止")
		}
	}

	time.Sleep(1 * time.Second)

	// ==================== 步骤9: 复杂场景模拟 ====================
	t.Log("\n[步骤9] 复杂场景模拟 - 巡查监控区域")

	t.Log("  • 场景：依次查看监控区域的不同位置")

	// 移动到左上角
	t.Log("    1) 移动到左上角...")
	if err := ctrl.UpLeft(5, 3*time.Second); err != nil {
		t.Logf("       ✗ 失败: %v", err)
	} else {
		t.Log("       ✓ 到达左上角，放大查看...")
		ctrl.ZoomIn(1 * time.Second)
		time.Sleep(2 * time.Second)
		ctrl.ZoomOut(1 * time.Second)
	}

	// 移动到右上角
	t.Log("    2) 移动到右上角...")
	if err := ctrl.UpRight(5, 3*time.Second); err != nil {
		t.Logf("       ✗ 失败: %v", err)
	} else {
		t.Log("       ✓ 到达右上角，放大查看...")
		ctrl.ZoomIn(1 * time.Second)
		time.Sleep(2 * time.Second)
		ctrl.ZoomOut(1 * time.Second)
	}

	// 移动到右下角
	t.Log("    3) 移动到右下角...")
	if err := ctrl.DownRight(5, 3*time.Second); err != nil {
		t.Logf("       ✗ 失败: %v", err)
	} else {
		t.Log("       ✓ 到达右下角")
		time.Sleep(1 * time.Second)
	}

	// 移动到左下角
	t.Log("    4) 移动到左下角...")
	if err := ctrl.DownLeft(5, 3*time.Second); err != nil {
		t.Logf("       ✗ 失败: %v", err)
	} else {
		t.Log("       ✓ 到达左下角")
		time.Sleep(1 * time.Second)
	}

	t.Log("  ✓ 区域巡查完成")

	// ==================== 步骤10: 手动精细控制示例 ====================
	t.Log("\n[步骤10] 手动精细控制示例")

	t.Log("  • 精细调整云台位置（多次微调）")

	// 微调1：右转
	t.Log("    - 微调右转...")
	ctrl.StartRight(2)
	time.Sleep(500 * time.Millisecond)
	ctrl.StopRight()
	t.Log("      ✓ 完成")

	time.Sleep(300 * time.Millisecond)

	// 微调2：上仰
	t.Log("    - 微调上仰...")
	ctrl.StartUp(2)
	time.Sleep(500 * time.Millisecond)
	ctrl.StopUp()
	t.Log("      ✓ 完成")

	time.Sleep(300 * time.Millisecond)

	// 微调3：焦距
	t.Log("    - 微调焦距...")
	ctrl.StartZoomIn()
	time.Sleep(300 * time.Millisecond)
	ctrl.StopZoomIn()
	t.Log("      ✓ 完成")

	// ==================== 步骤11: 回到原点 ====================
	t.Log("\n[步骤11] 回到原点")

	t.Log("  • 正在返回原点（预置点0）...")
	if err := preset.GotoPreset(0); err != nil {
		t.Logf("  ✗ 返回原点失败: %v", err)
		t.Log("  💡 提示：如果原点未设置，此步骤会失败")
	} else {
		t.Log("  ✓ 正在移动到原点...")
		time.Sleep(5 * time.Second) // 给足够时间回到原点
		t.Log("  ✓ 已回到原点")
	}

	// ==================== 测试总结 ====================
	t.Log("\n========================================")
	t.Log("测试总结")
	t.Log("========================================")
	t.Log("✓ 基础方向控制：4个方向")
	t.Log("✓ 组合方向控制：斜向移动")
	t.Log("✓ 自动计时控制：简单易用")
	t.Log("✓ 手动开始/停止：灵活精确")
	t.Log("✓ 相机焦距控制：放大/缩小")
	t.Log("✓ 相机焦点控制：近/远聚焦")
	t.Log("✓ 相机光圈控制：变亮/变暗")
	t.Log("✓ 预置点管理：设置和调用")
	t.Log("✓ 自动扫描：左右扫描")
	t.Log("✓ 复杂场景：区域巡查")
	t.Log("✓ 回到原点：完成循环")
	t.Log("\n💡 使用建议:")
	t.Log("  1. 优先使用自动计时方法（如 Right(speed, duration)）")
	t.Log("  2. 需要精细控制时使用手动方法（如 StartRight/StopRight）")
	t.Log("  3. 建议设置预置点0为原点，方便随时回归")
	t.Log("  4. 速度范围：1-7，建议常用值：3-5")
	t.Log("\n示例完成!")
}
