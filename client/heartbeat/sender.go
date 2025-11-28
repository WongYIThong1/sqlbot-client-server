package heartbeat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sqlbots-client/config"
	"sqlbots-client/encryption"
	"sqlbots-client/hardware"
	"sqlbots-client/session"
)

const (
	heartbeatInterval = 10 * time.Minute
)

// HeartbeatResponse 心跳响应结构体
type HeartbeatResponse struct {
	StatusCode  string `json:"status_code"`
	LicenseInfo struct {
		ExpiresAt string `json:"expires_at"`
		PlanType  string `json:"plan_type"`
	} `json:"license_info"`
	MachineInfo struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		RegisteredAt string `json:"registered_at"`
	} `json:"machine_info"`
	Message string `json:"message,omitempty"`
}

// SendHeartbeat 发送心跳
func SendHeartbeat(cfg *config.Config, machineInfo *hardware.MachineInfo, sessionManager *session.Manager) (*HeartbeatResponse, error) {
	// 构建请求数据
	requestData := map[string]interface{}{
		"machine_id":   machineInfo.MachineID,
		"machine_name": machineInfo.MachineName,
		"ram":          machineInfo.RAM,
		"cores":        machineInfo.Cores,
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %w", err)
	}

	// 确定使用哪个加密密钥：优先使用会话密钥
	encryptionKey := cfg.EncryptionKey
	useSessionKey := false

	if sessionManager != nil {
		if sessionKey, valid := sessionManager.GetSessionKey(); valid {
			encryptionKey = sessionKey
			useSessionKey = true
		}
	}

	// 加密数据
	encryptedData, err := encryption.Encrypt(string(jsonData), encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"API_KEY":         cfg.APIKey,
		"encrypted_data":  encryptedData,
		"use_session_key": useSessionKey,
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// 发送 HTTP POST 请求
	url := fmt.Sprintf("%s/heartbeat", cfg.ServerURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var response struct {
		EncryptedData string `json:"encrypted_data"`
	}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 解密响应数据（使用相同的密钥）
	decryptedText, err := encryption.Decrypt(response.EncryptedData, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}

	// 解析解密后的响应
	var heartbeatResp HeartbeatResponse
	if err := json.Unmarshal([]byte(decryptedText), &heartbeatResp); err != nil {
		return nil, fmt.Errorf("failed to parse decrypted response: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return &heartbeatResp, fmt.Errorf("server returned status %d: %s", resp.StatusCode, heartbeatResp.StatusCode)
	}

	return &heartbeatResp, nil
}

// HandleHeartbeatResponse 处理心跳响应，根据状态码决定是否终止程序
func HandleHeartbeatResponse(resp *HeartbeatResponse) error {
	switch resp.StatusCode {
	case "SUCCESS":
		return nil
	case "INVALID_API_KEY":
		return fmt.Errorf("INVALID_API_KEY: API Key is invalid")
	case "LICENSE_EXPIRED":
		return fmt.Errorf("LICENSE_EXPIRED: License has expired")
	case "MACHINE_LIMIT_EXCEEDED":
		return fmt.Errorf("MACHINE_LIMIT_EXCEEDED: Maximum machines limit exceeded")
	case "DECRYPTION_FAILED":
		return fmt.Errorf("DECRYPTION_FAILED: Failed to decrypt data")
	case "SERVER_ERROR":
		return fmt.Errorf("SERVER_ERROR: Server error occurred")
	default:
		return fmt.Errorf("unknown status code: %s", resp.StatusCode)
	}
}
