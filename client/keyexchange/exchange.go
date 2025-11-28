package keyexchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sqlbots-client/config"
	"sqlbots-client/encryption"
	"sqlbots-client/session"
)

// KeyExchangeResponse 密钥交换响应结构体
type KeyExchangeResponse struct {
	StatusCode string `json:"status_code"`
	SessionKey string `json:"session_key"` // 加密的会话密钥
	ExpiresIn  int    `json:"expires_in"`  // 过期时间（秒）
	Username   string `json:"username"`    // 用户名
	Message    string `json:"message,omitempty"`
}

// ExchangeKey 执行密钥交换，返回用户名
func ExchangeKey(cfg *config.Config, sessionManager *session.Manager) (string, error) {
	// 使用初始 ENCRYPTION_KEY 进行密钥交换
	initialKey := cfg.EncryptionKey
	
	// 构建请求（不需要加密，因为这是初始连接）
	requestBody := map[string]string{
		"API_KEY": cfg.APIKey,
	}
	
	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}
	
	// 发送 HTTP POST 请求
	url := fmt.Sprintf("%s/key-exchange", cfg.ServerURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	// 读取响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}
	
	// 解析响应
	var keyExchangeResp KeyExchangeResponse
	if err := json.Unmarshal(responseBody, &keyExchangeResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	
	// 检查状态码
	if keyExchangeResp.StatusCode != "SUCCESS" {
		return "", fmt.Errorf("key exchange failed: %s", keyExchangeResp.Message)
	}
	
	// 解密会话密钥（使用初始密钥）
	decryptedSessionKey, err := encryption.Decrypt(keyExchangeResp.SessionKey, initialKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt session key: %w", err)
	}
	
	// 保存会话密钥
	sessionManager.SetSessionKey(decryptedSessionKey, keyExchangeResp.ExpiresIn)
	
	// 返回用户名（如果存在）
	if keyExchangeResp.Username != "" {
		return keyExchangeResp.Username, nil
	}
	
	return "", nil
}

