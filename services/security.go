package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
	"time"
)

// SecurityManager 安全管理器
type SecurityManager struct {
	nonceStore map[string]*NonceInfo
	mu         sync.RWMutex
}

// NonceInfo nonce 信息
type NonceInfo struct {
	Value      string
	WindowId   string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	Used       bool
}

// PKCEInfo PKCE 信息
type PKCEInfo struct {
	CodeVerifier string
	CodeChallenge string
	Method       string
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager() *SecurityManager {
	sm := &SecurityManager{
		nonceStore: make(map[string]*NonceInfo),
	}

	// 启动清理协程
	go sm.cleanupExpiredNonces()

	return sm
}

// GenerateNonce 生成安全的 nonce
func (sm *SecurityManager) GenerateNonce(windowId string) (string, error) {
	// 生成 32 字节的随机数
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("生成随机数失败: %w", err)
	}

	nonce := hex.EncodeToString(bytes)

	// 存储 nonce 信息
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.nonceStore[nonce] = &NonceInfo{
		Value:     nonce,
		WindowId:  windowId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute), // 10分钟有效期
		Used:      false,
	}

	return nonce, nil
}

// ValidateNonce 验证 nonce
func (sm *SecurityManager) ValidateNonce(nonce, windowId string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	info, exists := sm.nonceStore[nonce]
	if !exists {
		return fmt.Errorf("nonce 不存在")
	}

	// 检查是否已使用
	if info.Used {
		return fmt.Errorf("nonce 已被使用")
	}

	// 检查是否过期
	if time.Now().After(info.ExpiresAt) {
		return fmt.Errorf("nonce 已过期")
	}

	// 检查 windowId
	if info.WindowId != windowId {
		return fmt.Errorf("windowId 不匹配")
	}

	return nil
}

// MarkNonceUsed 标记 nonce 为已使用
func (sm *SecurityManager) MarkNonceUsed(nonce string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if info, exists := sm.nonceStore[nonce]; exists {
		info.Used = true
	}
}

// GeneratePKCE 生成 PKCE 参数
func (sm *SecurityManager) GeneratePKCE() (*PKCEInfo, error) {
	// 生成 code_verifier (43-128 字符的随机字符串)
	verifierBytes := make([]byte, 32)
	if _, err := rand.Read(verifierBytes); err != nil {
		return nil, fmt.Errorf("生成 code_verifier 失败: %w", err)
	}

	codeVerifier := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(verifierBytes)

	// 生成 code_challenge
	hash := sha256.Sum256([]byte(codeVerifier))
	codeChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	return &PKCEInfo{
		CodeVerifier:  codeVerifier,
		CodeChallenge: codeChallenge,
		Method:        "S256",
	}, nil
}

// ValidatePKCE 验证 PKCE
func (sm *SecurityManager) ValidatePKCE(codeVerifier, codeChallenge, method string) bool {
	if method != "S256" {
		return false
	}

	// 重新计算 challenge
	hash := sha256.Sum256([]byte(codeVerifier))
	expectedChallenge := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(hash[:])

	return expectedChallenge == codeChallenge
}

// GenerateRandomString 生成随机字符串
func (sm *SecurityManager) GenerateRandomString(length int) string {
	bytes := make([]byte, length/2)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// cleanupExpiredNonces 清理过期的 nonce
func (sm *SecurityManager) cleanupExpiredNonces() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sm.mu.Lock()
		now := time.Now()

		for nonce, info := range sm.nonceStore {
			if now.After(info.ExpiresAt) {
				delete(sm.nonceStore, nonce)
			}
		}

		sm.mu.Unlock()
	}
}

// GetNonceInfo 获取 nonce 信息（用于调试）
func (sm *SecurityManager) GetNonceInfo(nonce string) (*NonceInfo, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	info, exists := sm.nonceStore[nonce]
	if !exists {
		return nil, fmt.Errorf("nonce 不存在")
	}

	return info, nil
}

// IsSafeURL 检查 URL 是否安全
func (sm *SecurityManager) IsSafeURL(urlStr string) bool {
	// 检查协议
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return false
	}

	// 检查是否包含危险字符
	dangerousChars := []string{"<", ">", "\"", "'", "&", "\n", "\r", "\t"}
	for _, char := range dangerousChars {
		if strings.Contains(urlStr, char) {
			return false
		}
	}

	return true
}

// SanitizeString 清理字符串
func (sm *SecurityManager) SanitizeString(input string) string {
	// 移除危险字符
	dangerousChars := []string{"<", ">", "\"", "'", "&", "\n", "\r", "\t"}
	result := input

	for _, char := range dangerousChars {
		result = strings.ReplaceAll(result, char, "")
	}

	// 限制长度
	if len(result) > 1000 {
		result = result[:1000]
	}

	return result
}