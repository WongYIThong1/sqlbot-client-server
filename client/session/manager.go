package session

import (
	"sync"
	"time"
)

// SessionKey 会话密钥结构
type SessionKey struct {
	Key       string
	ExpiresAt time.Time
}

// Manager 会话密钥管理器
type Manager struct {
	mu        sync.RWMutex
	sessionKey *SessionKey
}

// NewManager 创建新的会话密钥管理器
func NewManager() *Manager {
	return &Manager{}
}

// SetSessionKey 设置会话密钥
func (m *Manager) SetSessionKey(key string, expiresIn int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.sessionKey = &SessionKey{
		Key:       key,
		ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second),
	}
}

// GetSessionKey 获取会话密钥（如果未过期）
func (m *Manager) GetSessionKey() (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if m.sessionKey == nil {
		return "", false
	}
	
	if time.Now().After(m.sessionKey.ExpiresAt) {
		return "", false
	}
	
	return m.sessionKey.Key, true
}

// ClearSessionKey 清除会话密钥
func (m *Manager) ClearSessionKey() {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.sessionKey = nil
}

// HasValidSession 检查是否有有效的会话密钥
func (m *Manager) HasValidSession() bool {
	_, valid := m.GetSessionKey()
	return valid
}

