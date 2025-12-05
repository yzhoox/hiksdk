package examples

import (
	"testing"
	"time"

	"github.com/samsaralc/hiksdk/core/auth"
	"github.com/samsaralc/hiksdk/core/ptz"
)

func TestPTZControl2(t *testing.T) {
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

	// 创建预置点管理器
	preset := ptz.NewPresetManager(session.LoginID, channel)
	// ==================== 步骤7: 预置点操作 ====================
	t.Log("\n[步骤7] 预置点操作")

	// 转到预置点2
	t.Log("  • 转到预置点2...")
	if err := preset.GotoPreset(2); err != nil {
		t.Logf("    ✗ 失败: %v", err)
	} else {
		t.Log("    ✓ 正在移动到预置点2...")
		time.Sleep(3 * time.Second)
		t.Log("    ✓ 已到达预置点2")
	}

}
