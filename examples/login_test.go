package examples

import (
	"testing"

	"github.com/samsaralc/hiksdk/core/auth"
)

// TestLoginMethods 两种登录方式示例
func TestLoginMethods(t *testing.T) {
	t.Log("========================================")
	t.Log("海康威视 SDK - 登录方式示例")
	t.Log("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64", // 替换为你的设备IP
		Port:     8000,           // 替换为你的端口
		Username: "admin",        // 替换为你的用户名
		Password: "password",     // 替换为你的密码
	}

	t.Logf("\n设备连接信息:")
	t.Logf("  - IP地址: %s", cred.IP)
	t.Logf("  - 端口: %d", cred.Port)
	t.Logf("  - 用户名: %s", cred.Username)

	// ==================== 方式1: LoginV40 (推荐) ====================
	t.Log("\n========================================")
	t.Log("方式1: 使用 LoginV40（推荐）")
	t.Log("========================================")

	t.Log("\n[1] 使用LoginV40登录...")
	session1, err := auth.LoginV40(cred)
	if err != nil {
		t.Logf("✗ 登录失败: %v", err)
		t.Log("\n可能的原因:")
		t.Log("  1. 设备不在线或网络不可达")
		t.Log("  2. 用户名或密码错误")
		t.Log("  3. 设备端口配置错误")
		t.Skip("跳过测试（设备未连接）")
	} else {
		t.Logf("✓ 登录成功")
		t.Logf("  登录ID: %d", session1.LoginID)
		t.Logf("  设备序列号: %s", session1.SerialNumber)
		t.Logf("  通道数: %d", session1.ChannelNum)

		// 登出
		t.Log("\n[2] 登出设备...")
		if err := auth.Logout(session1.LoginID); err != nil {
			t.Errorf("✗ 登出失败: %v", err)
		}
	}

	// ==================== 方式2: LoginV30 (兼容旧设备) ====================
	t.Log("\n========================================")
	t.Log("方式2: 使用 LoginV30（兼容旧设备）")
	t.Log("========================================")

	t.Log("\n[1] 使用LoginV30登录...")
	session2, err := auth.LoginV30(cred)
	if err != nil {
		t.Logf("✗ 登录失败: %v", err)
	} else {
		t.Logf("✓ 登录成功")
		t.Logf("  登录ID: %d", session2.LoginID)
		t.Logf("  设备序列号: %s", session2.SerialNumber)
		t.Logf("  通道数: %d", session2.ChannelNum)

		// 登出
		t.Log("\n[2] 登出设备...")
		if err := auth.Logout(session2.LoginID); err != nil {
			t.Errorf("✗ 登出失败: %v", err)
		}
	}

	// ==================== 对比说明 ====================
	t.Log("\n========================================")
	t.Log("两种登录方式对比")
	t.Log("========================================")

	t.Log("\nLoginV40():")
	t.Log("  ✓ 推荐使用")
	t.Log("  ✓ 支持更多功能")
	t.Log("  ✓ 性能更好")
	t.Log("  ✓ 设备信息更详细")

	t.Log("\nLoginV30():")
	t.Log("  ✓ 兼容旧设备")
	t.Log("  ✓ 简单直接")

	t.Log("\n💡 建议:")
	t.Log("  1. 优先使用 LoginV40()")
	t.Log("  2. 如果失败，可尝试 LoginV30()")
	t.Log("  3. 登录后务必调用 Logout() 释放资源")

	// 程序结束时清理SDK
	defer auth.Cleanup()
}
