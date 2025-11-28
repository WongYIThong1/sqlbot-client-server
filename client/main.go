package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"sqlbots-client/config"
	"sqlbots-client/hardware"
	"sqlbots-client/heartbeat"
	"sqlbots-client/keyexchange"
	"sqlbots-client/session"
	"sqlbots-client/ui"
)

const version = "v1.0"

func main() {
	// 显示欢迎界面
	ui.ClearScreen()
	ui.ShowBanner(version)

	// 提示输入 API Key
	ui.ShowLoginPrompt()
	apiKey, err := ui.HideInput()
	if err != nil {
		fmt.Printf("\n❌ Failed to read input: %v\n", err)
		os.Exit(1)
	}

	if apiKey == "" {
		ui.ShowError("API Key cannot be empty")
		os.Exit(1)
	}

	// 创建配置（使用输入的 API Key）
	cfg := &config.Config{
		APIKey:        apiKey,
		ServerURL:     getEnvOrDefault("SERVER_URL", "https://api.sqlbots.online"),
		EncryptionKey: getEnvOrDefault("ENCRYPTION_KEY", ""),
	}

	// 验证 ENCRYPTION_KEY
	if cfg.EncryptionKey == "" {
		ui.ShowError("ENCRYPTION_KEY is required (use ENCRYPTION_KEY environment variable)")
		os.Exit(1)
	}

	// 获取机器信息
	machineInfo, err := hardware.GetMachineInfo()
	if err != nil {
		ui.ShowError(fmt.Sprintf("Failed to get machine info: %v", err))
		os.Exit(1)
	}

	// 初始化会话密钥管理器
	sessionManager := session.NewManager()

	// 进行密钥交换（获取用户名）
	username, err := keyexchange.ExchangeKey(cfg, sessionManager)
	if err != nil {
		ui.ShowError(fmt.Sprintf("Authentication failed: %v", err))
		os.Exit(1)
	}

	// 如果服务器没有返回用户名，使用默认值
	if username == "" {
		username = "User"
	}

	// 显示登录成功界面
	ui.ShowLoggedIn(username, version)

	// 显示机器信息（可选，可以注释掉）
	// fmt.Printf("Machine ID: %s\n", machineInfo.MachineID)
	// fmt.Printf("Machine Name: %s\n", machineInfo.MachineName)
	// fmt.Printf("RAM: %d GB | CPU Cores: %d\n\n", machineInfo.RAM, machineInfo.Cores)

	// 设置优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 启动时立即发送首次心跳
	if err := sendHeartbeatAndHandle(cfg, machineInfo, sessionManager); err != nil {
		ui.ShowError(fmt.Sprintf("Initial heartbeat failed: %v", err))
		os.Exit(1)
	}

	// 设置定时器每 10 分钟发送心跳
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	// 心跳循环
	go func() {
		for {
			select {
			case <-ticker.C:
				// 检查会话密钥是否即将过期（提前5分钟刷新）
				if !sessionManager.HasValidSession() {
					_, _ = keyexchange.ExchangeKey(cfg, sessionManager)
				}

				// 静默发送心跳（不显示日志，避免干扰界面）
				if err := sendHeartbeatAndHandle(cfg, machineInfo, sessionManager); err != nil {
					// 检查是否是致命错误
					if isFatalError(err) {
						ui.ShowError(fmt.Sprintf("Fatal error: %v", err))
						os.Exit(1)
					}
					// 非致命错误继续运行（静默处理）
				}
			case <-sigChan:
				ui.ClearScreen()
				fmt.Println("Shutting down...")
				return
			}
		}
	}()

	// 等待退出信号
	<-sigChan
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// sendHeartbeatAndHandle 发送心跳并处理响应
func sendHeartbeatAndHandle(cfg *config.Config, machineInfo *hardware.MachineInfo, sessionManager *session.Manager) error {
	resp, err := heartbeat.SendHeartbeat(cfg, machineInfo, sessionManager)
	if err != nil {
		return err
	}

	// 处理响应，检查是否需要终止程序
	if err := heartbeat.HandleHeartbeatResponse(resp); err != nil {
		return err
	}

	// 静默处理成功响应（不显示日志）
	_ = resp
	return nil
}

// isFatalError 判断是否是致命错误（需要终止程序）
func isFatalError(err error) bool {
	if err == nil {
		return false
	}

	errStr := err.Error()
	fatalErrors := []string{
		"INVALID_API_KEY",
		"LICENSE_EXPIRED",
		"MACHINE_LIMIT_EXCEEDED",
	}

	for _, fatalErr := range fatalErrors {
		if strings.Contains(errStr, fatalErr) {
			return true
		}
	}

	return false
}
