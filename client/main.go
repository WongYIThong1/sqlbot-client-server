package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Configuration
const (
	DefaultServerURL = "https://api.sqlbots.online"
	HeartbeatInterval = 10 * time.Minute
	DefaultEncryptionKey = "385918897f57e8fedf49c4230f316661bdc6b6c2157094d9b43715ba87100851" // Must match server ENCRYPTION_KEY
)

// HeartbeatRequest represents the heartbeat data structure
type HeartbeatRequest struct {
	APIKey     string `json:"apiKey"`
	MachineName string `json:"machineName"`
	MachineID  string `json:"machineId"`
	RAM        int    `json:"ram"`
	Cores      int    `json:"cores"`
}

// HeartbeatResponse represents the server response
type HeartbeatResponse struct {
	Data string `json:"data"`
}

// DecryptedResponse represents the decrypted response data
type DecryptedResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	Error            string `json:"error"`
	Code             string `json:"code"`
	LicenseValidUntil string `json:"licenseValidUntil"`
	IsNewMachine     bool   `json:"isNewMachine"`
	HardwareChanged  bool   `json:"hardwareChanged"`
}

// Client represents the client application
type Client struct {
	apiKey      string
	machineID   string
	machineName string
	serverURL   string
	encKey      []byte
	stopChan    chan os.Signal
}

// NewClient creates a new client instance
func NewClient(apiKey string) (*Client, error) {
	// Get machine ID
	machineID, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("failed to get machine ID: %v", err)
	}

	// Get machine name (hostname)
	machineName, err := os.Hostname()
	if err != nil {
		machineName = "Unknown"
	}

	// Get server URL from environment or use default
	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = DefaultServerURL
	}

	// Get encryption key from environment or use default
	encKeyStr := os.Getenv("ENCRYPTION_KEY")
	if encKeyStr == "" {
		encKeyStr = DefaultEncryptionKey
	}

	// Prepare encryption key (pad to 32 bytes for AES-256)
	encKey := []byte(encKeyStr)
	if len(encKey) < 32 {
		padded := make([]byte, 32)
		copy(padded, encKey)
		encKey = padded
	} else if len(encKey) > 32 {
		encKey = encKey[:32]
	}

	return &Client{
		apiKey:      apiKey,
		machineID:   machineID,
		machineName: machineName,
		serverURL:   serverURL,
		encKey:      encKey,
		stopChan:    make(chan os.Signal, 1),
	}, nil
}

// GetSystemInfo collects RAM and CPU core information
func (c *Client) GetSystemInfo() (ramGB int, cores int, err error) {
	// Get memory info
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get memory info: %v", err)
	}
	ramGB = int(memInfo.Total / (1024 * 1024 * 1024)) // Convert bytes to GB

	// Get CPU info
	cpuCount, err := cpu.Counts(true)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get CPU info: %v", err)
	}
	cores = cpuCount

	return ramGB, cores, nil
}

// evpBytesToKey implements OpenSSL's EVP_BytesToKey for key derivation
// This matches CryptoJS behavior when using a string password
func evpBytesToKey(password []byte, salt []byte, keyLen int, ivLen int) (key []byte, iv []byte) {
	m := []byte{}
	i := 0
	key = make([]byte, keyLen)
	iv = make([]byte, ivLen)

	for len(m) < keyLen+ivLen {
		hash := md5.New()
		if i > 0 {
			hash.Write(m[len(m)-16:])
		}
		hash.Write(password)
		hash.Write(salt)
		m = append(m, hash.Sum(nil)...)
		i++
	}

	copy(key, m[:keyLen])
	copy(iv, m[keyLen:keyLen+ivLen])
	return
}

// Encrypt encrypts data using AES-CBC with OpenSSL-compatible format
// Format: Salt (8 bytes) + IV (16 bytes) + Ciphertext
// This matches CryptoJS.AES.encrypt() behavior
func (c *Client) Encrypt(plaintext []byte) (string, error) {
	// Generate random salt (8 bytes)
	salt := make([]byte, 8)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// Derive key and IV using EVP_BytesToKey (OpenSSL compatible)
	key, iv := evpBytesToKey(c.encKey, salt, 32, 16) // AES-256 key (32 bytes), IV (16 bytes)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Pad plaintext to block size
	plaintext = pkcs7Pad(plaintext, aes.BlockSize)

	// Encrypt
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	// Combine: Salt + IV + Ciphertext
	result := append(salt, iv...)
	result = append(result, ciphertext...)

	// Return base64 encoded
	return base64.StdEncoding.EncodeToString(result), nil
}

// Decrypt decrypts data using AES-CBC with OpenSSL-compatible format
// Format: Salt (8 bytes) + IV (16 bytes) + Ciphertext
func (c *Client) Decrypt(ciphertext string) ([]byte, error) {
	// Decode base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	// Minimum size: 8 (salt) + 16 (IV) + 16 (min ciphertext block)
	if len(data) < 40 {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract salt, IV, and ciphertext
	salt := data[:8]
	iv := data[8:24]
	ciphertextBytes := data[24:]

	if len(ciphertextBytes)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// Derive key and IV using EVP_BytesToKey
	key, derivedIV := evpBytesToKey(c.encKey, salt, 32, 16)

	// Use derived IV (should match, but we use the one from data)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextBytes))
	mode.CryptBlocks(plaintext, ciphertextBytes)

	// Remove padding
	plaintext = pkcs7Unpad(plaintext)

	return plaintext, nil
}

// pkcs7Pad adds PKCS7 padding
func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpad removes PKCS7 padding
func pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	return data[:(length - unpadding)]
}

// SendHeartbeat sends a heartbeat to the server
func (c *Client) SendHeartbeat() (*DecryptedResponse, error) {
	// Get system info
	ramGB, cores, err := c.GetSystemInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get system info: %v", err)
	}

	// Create heartbeat request
	heartbeat := HeartbeatRequest{
		APIKey:      c.apiKey,
		MachineName: c.machineName,
		MachineID:   c.machineID,
		RAM:         ramGB,
		Cores:       cores,
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(heartbeat)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal heartbeat: %v", err)
	}

	// Encrypt data
	encryptedData, err := c.Encrypt(jsonData)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %v", err)
	}

	// Create HTTP request
	requestBody := map[string]string{
		"data": encryptedData,
	}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Send POST request
	resp, err := http.Post(c.serverURL+"/heartbeat", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send heartbeat: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	var response HeartbeatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Decrypt response
	decryptedBytes, err := c.Decrypt(response.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt response: %v", err)
	}

	// Parse decrypted response
	var decryptedResp DecryptedResponse
	if err := json.Unmarshal(decryptedBytes, &decryptedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal decrypted response: %v", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		return &decryptedResp, fmt.Errorf("server returned status %d: %s", resp.StatusCode, decryptedResp.Error)
	}

	return &decryptedResp, nil
}

// Run starts the client heartbeat loop
func (c *Client) Run() error {
	log.Println("Starting SQLBots Client...")
	log.Printf("Machine ID: %s", c.machineID)
	log.Printf("Machine Name: %s", c.machineName)

	// Setup signal handling for graceful shutdown
	signal.Notify(c.stopChan, os.Interrupt, syscall.SIGTERM)

	// Send initial heartbeat immediately
	log.Println("Sending initial heartbeat...")
	resp, err := c.SendHeartbeat()
	if err != nil {
		return fmt.Errorf("initial heartbeat failed: %v", err)
	}

	// Check response
	if !resp.Success {
		if resp.Code == "LICENSE_EXPIRED" {
			log.Fatal("License has expired. Exiting...")
		}
		if resp.Code == "MACHINE_LIMIT_EXCEEDED" {
			log.Fatal("Maximum machine limit reached. Unable to configure new machine.")
		}
		if resp.Code == "INVALID_API_KEY" {
			log.Fatal("Invalid API key. Exiting...")
		}
		return fmt.Errorf("heartbeat failed: %s", resp.Error)
	}

	log.Printf("Heartbeat successful: %s", resp.Message)
	if resp.IsNewMachine {
		log.Println("Machine registered successfully")
	}
	if resp.HardwareChanged {
		log.Println("Machine hardware info updated")
	}

	// Start periodic heartbeat
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-c.stopChan:
			log.Println("Shutting down...")
			return nil

		case <-ticker.C:
			log.Println("Sending heartbeat...")
			resp, err := c.SendHeartbeat()
			if err != nil {
				log.Printf("Heartbeat error: %v", err)
				continue
			}

			// Check response
			if !resp.Success {
				if resp.Code == "LICENSE_EXPIRED" {
					log.Fatal("License has expired. Exiting...")
				}
				if resp.Code == "MACHINE_LIMIT_EXCEEDED" {
					log.Fatal("Maximum machine limit reached. Unable to configure new machine.")
				}
				log.Printf("Heartbeat failed: %s", resp.Error)
				continue
			}

			log.Printf("Heartbeat successful: %s", resp.Message)
			if resp.LicenseValidUntil != "" {
				log.Printf("License valid until: %s", resp.LicenseValidUntil)
			}
		}
	}
}

func main() {
	// Get API key from environment variable or command line argument
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		if len(os.Args) > 1 {
			apiKey = os.Args[1]
		} else {
			log.Fatal("API_KEY environment variable or command line argument is required")
		}
	}

	// Create client
	client, err := NewClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Run client
	if err := client.Run(); err != nil {
		log.Fatalf("Client error: %v", err)
	}
}

