package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
)

const (
	saltPrefix = "Salted__"
	saltLen    = 8
)

// evpBytesToKey 实现 EVP_BytesToKey 密钥派生（兼容 OpenSSL 和 CryptoJS）
func evpBytesToKey(password, salt []byte) (key, iv []byte) {
	keySize := 32 // AES-256 (32 bytes)
	ivSize := 16  // AES block size (16 bytes)

	// 第一次迭代：MD5(password + salt)
	hash := md5.New()
	hash.Write(password)
	hash.Write(salt)
	derived := hash.Sum(nil)

	// 生成 key (需要 32 字节，每次 MD5 产生 16 字节，所以需要 2 次)
	key = make([]byte, keySize)
	copy(key, derived) // 第一次：16 字节
	if keySize > 16 {
		// 第二次迭代：MD5(derived + password + salt)
		hash = md5.New()
		hash.Write(derived)
		hash.Write(password)
		hash.Write(salt)
		derived = hash.Sum(nil)
		copy(key[16:], derived) // 第二次：16 字节
	}

	// 生成 IV (需要 16 字节)
	// 使用当前的 derived (来自 key 生成的最后一次迭代)
	hash = md5.New()
	hash.Write(key) // 使用完整的 key
	hash.Write(password)
	hash.Write(salt)
	derived = hash.Sum(nil)
	iv = make([]byte, ivSize)
	copy(iv, derived)

	return key[:keySize], iv[:ivSize]
}

// Encrypt 加密数据（OpenSSL 兼容格式）
func Encrypt(plaintext, password string) (string, error) {
	// 生成随机盐值
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// 派生密钥和 IV
	key, iv := evpBytesToKey([]byte(password), salt)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// PKCS7 填充
	plaintextBytes := []byte(plaintext)
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	plaintextBytes = append(plaintextBytes, padtext...)

	// 加密
	ciphertext := make([]byte, len(plaintextBytes))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintextBytes)

	// 组合：Salted__ + salt + ciphertext
	combined := append([]byte(saltPrefix), salt...)
	combined = append(combined, ciphertext...)

	// Base64 编码
	return base64.StdEncoding.EncodeToString(combined), nil
}

// Decrypt 解密数据（OpenSSL 兼容格式）
func Decrypt(ciphertext, password string) (string, error) {
	// Base64 解码
	encryptedData, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// 检查格式
	if len(encryptedData) < len(saltPrefix)+saltLen {
		return "", errors.New("invalid encrypted data format")
	}

	// 检查前缀
	if string(encryptedData[:len(saltPrefix)]) != saltPrefix {
		return "", errors.New("invalid encrypted data format: missing salt prefix")
	}

	// 提取 salt
	salt := encryptedData[len(saltPrefix) : len(saltPrefix)+saltLen]

	// 提取加密数据
	ciphertextOnly := encryptedData[len(saltPrefix)+saltLen:]

	// 派生密钥和 IV
	key, iv := evpBytesToKey([]byte(password), salt)

	// 创建 AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// 检查数据长度
	if len(ciphertextOnly)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of block size")
	}

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertextOnly))
	mode.CryptBlocks(plaintext, ciphertextOnly)

	// 移除 PKCS7 填充
	padding := int(plaintext[len(plaintext)-1])
	if padding > aes.BlockSize || padding == 0 {
		return "", errors.New("invalid padding")
	}
	plaintext = plaintext[:len(plaintext)-padding]

	return string(plaintext), nil
}
