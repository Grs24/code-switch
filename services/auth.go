package services

import (
	"codeswitch/config"
	"encoding/json"
	"fmt"
	"sync"
)

// UserInfo 用户信息
type UserInfo struct {
	Token    string                 `json:"token"`
	UserID   string                 `json:"userId"`
	Email    string                 `json:"email"`
	Mobile   string                 `json:"mobile"`
	UserInfo map[string]interface{} `json:"userInfo"` // 完整的后端用户数据
}

// AuthService 认证服务
type AuthService struct {
	mu       sync.RWMutex
	userInfo *UserInfo
	isLogin  bool
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		isLogin: false,
	}
}

// GetLoginURL 获取登录页 URL
func (s *AuthService) GetLoginURL() string {
	return config.GetLoginURL()
}

// GetAPIURL 获取 API 基础 URL
func (s *AuthService) GetAPIURL() string {
	return config.GetAPIURL()
}

// GetAICodeBaseURL 获取 AICoding Base URL
func (s *AuthService) GetAICodeBaseURL() string {
	return config.GetAICodeBaseURL()
}

// GetSubscriptionURL 获取订阅页面 URL
func (s *AuthService) GetSubscriptionURL() string {
	// 根据 LoginURL 推导订阅页面 URL
	loginURL := config.GetLoginURL()
	// 移除 /desktop-login 后缀，添加 /subscribe
	baseURL := loginURL
	if len(baseURL) > 14 && baseURL[len(baseURL)-14:] == "/desktop-login" {
		baseURL = baseURL[:len(baseURL)-14]
	}
	return baseURL + "/subscribe"
}

// SetUserInfo 设置用户信息（从前端调用）
func (s *AuthService) SetUserInfo(userInfoJSON string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userInfo UserInfo
	if err := json.Unmarshal([]byte(userInfoJSON), &userInfo); err != nil {
		return fmt.Errorf("解析用户信息失败: %w", err)
	}

	s.userInfo = &userInfo
	s.isLogin = true

	tokenPreview := userInfo.Token
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}
	fmt.Printf("[AuthService] 用户登录成功: %s (Token: %s)\n", userInfo.Email, tokenPreview)
	return nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo() (*UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isLogin || s.userInfo == nil {
		return nil, fmt.Errorf("用户未登录")
	}

	return s.userInfo, nil
}

// IsLogin 检查是否已登录
func (s *AuthService) IsLogin() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isLogin
}

// Logout 退出登录
func (s *AuthService) Logout() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.userInfo = nil
	s.isLogin = false
	fmt.Println("[AuthService] 用户已退出登录")
}

// GetToken 获取用户 token
func (s *AuthService) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.userInfo != nil {
		return s.userInfo.Token
	}
	return ""
}
