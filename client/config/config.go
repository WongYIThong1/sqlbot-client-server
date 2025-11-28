package config

import (
	"flag"
	"fmt"
	"os"
)

// Config 配置结构体
type Config struct {
	APIKey        string
	ServerURL     string
	EncryptionKey string
}

// LoadConfig 加载配置（从命令行参数和环境变量）
func LoadConfig() (*Config, error) {
	config := &Config{}

	// 先读取环境变量作为默认值
	defaultAPIKey := os.Getenv("API_KEY")
	defaultServerURL := os.Getenv("SERVER_URL")
	if defaultServerURL == "" {
		defaultServerURL = "https://api.sqlbots.online"
	}
	defaultEncryptionKey := os.Getenv("ENCRYPTION_KEY")

	// 命令行参数
	flag.StringVar(&config.APIKey, "api-key", defaultAPIKey, "API Key (required, can also use API_KEY env var)")
	flag.StringVar(&config.ServerURL, "server-url", defaultServerURL, "Server URL (default: https://api.sqlbots.online)")
	flag.StringVar(&config.EncryptionKey, "encryption-key", defaultEncryptionKey, "Encryption key (required, can also use ENCRYPTION_KEY env var)")
	flag.Parse()

	// 如果命令行参数为空，再次尝试从环境变量读取
	if config.APIKey == "" {
		config.APIKey = os.Getenv("API_KEY")
	}
	if config.ServerURL == "" {
		if envURL := os.Getenv("SERVER_URL"); envURL != "" {
			config.ServerURL = envURL
		} else {
			config.ServerURL = "https://api.sqlbots.online"
		}
	}
	if config.EncryptionKey == "" {
		config.EncryptionKey = os.Getenv("ENCRYPTION_KEY")
	}

	// 验证必填字段
	if config.APIKey == "" {
		return nil, fmt.Errorf("API_KEY is required (use --api-key or API_KEY environment variable)")
	}
	if config.EncryptionKey == "" {
		return nil, fmt.Errorf("ENCRYPTION_KEY is required (use --encryption-key or ENCRYPTION_KEY environment variable)")
	}

	return config, nil
}
