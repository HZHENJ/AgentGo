package utils

import (
    "os"
    "testing"
    "time"

    "agentgo/pkg/conf"
)

// TestSendEmail_Smoke
// 集成测试（可选执行）：当设置 ENABLE_EMAIL_TEST=1 且提供 TEST_EMAIL_TO 时，
// 读取本地 config（config/config.yaml 或 APP_ENV=prod 时的 config.prod.yaml），
// 调用 SendEmail 真实发送一封测试邮件以验证配置有效性。
func TestSendEmail_Smoke(t *testing.T) {
    if os.Getenv("ENABLE_EMAIL_TEST") != "1" {
        t.Skip("skip: set ENABLE_EMAIL_TEST=1 to run this integration test")
    }

    to := os.Getenv("TEST_EMAIL_TO")
    if to == "" {
        t.Skip("skip: set TEST_EMAIL_TO to a valid recipient email address")
    }

    err := os.Chdir("../..")
	if err != nil {
		t.Fatalf("Failed to change working directory: %v", err)
	}

    // 可选：APP_ENV=prod 时会读取 config/config.prod.yaml
    // 否则读取 config/config.yaml
    conf.Init()

    code := "Test-" + time.Now().Format(time.RFC3339)
    msg := "AgentGo Email Smoke Test"

    if err := SendEmail(to, code, msg); err != nil {
        t.Fatalf("SendEmail failed: %v", err)
    }
}
