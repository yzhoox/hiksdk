package examples

import (
	"testing"
	"time"

	"github.com/samsaralc/hiksdk/core/alarm"
	"github.com/samsaralc/hiksdk/core/auth"
)

// TestAlarmListen 报警监听示例
func TestAlarmListen(t *testing.T) {
	t.Log("========================================")
	t.Log("海康威视 SDK - 报警监听示例")
	t.Log("========================================")

	// 设备连接凭据
	cred := &auth.Credentials{
		IP:       "192.168.1.64",
		Port:     8000,
		Username: "admin",
		Password: "password",
	}

	// 登录设备
	session, err := auth.LoginV40(cred)
	if err != nil {
		t.Skipf("登录失败: %v", err)
		return
	}
	t.Logf("登录成功 (ID: %d)", session.LoginID)
	defer auth.Logout(session.LoginID)
	defer auth.Cleanup()

	// 创建报警监听器
	t.Log("\n创建报警监听器...")
	listener := alarm.NewAlarmListener(session.LoginID)

	// 启动报警监听
	t.Log("启动报警监听...")
	if err := listener.Start(); err != nil {
		t.Errorf("启动监听失败: %v", err)
		return
	}
	defer listener.Stop()

	t.Log("\n监听中... 等待10秒接收报警消息")
	t.Log("（移动侦测、遮挡报警等）")

	// 等待一段时间以接收报警
	time.Sleep(10 * time.Second)

	t.Log("\n停止监听...")
	t.Log("示例完成!")
}
