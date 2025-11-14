package services

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
)

// DeepLinkService 处理 Deep Link 回调
type DeepLinkService struct {
	onAuthCallback func(token string, userID string, email string)
}

// NewDeepLinkService 创建 Deep Link 服务
func NewDeepLinkService() *DeepLinkService {
	return &DeepLinkService{}
}

// SetAuthCallback 设置登录回调
func (s *DeepLinkService) SetAuthCallback(callback func(token string, userID string, email string)) {
	s.onAuthCallback = callback
}

// HandleDeepLink 处理 Deep Link URL - 支持 OAuth 和 Legacy 两种格式
func (s *DeepLinkService) HandleDeepLink(deepLinkURL string) error {
	fmt.Printf("收到 Deep Link: %s\n", deepLinkURL)

	// 解析 URL
	parsedURL, err := url.Parse(deepLinkURL)
	if err != nil {
		return fmt.Errorf("解析 Deep Link 失败: %w", err)
	}

	// 检查 scheme
	if parsedURL.Scheme != "codeswitch" {
		return fmt.Errorf("不支持的协议: %s", parsedURL.Scheme)
	}

	// 检查 host (路径)
	if parsedURL.Host != "auth" {
		return fmt.Errorf("不支持的路径: %s", parsedURL.Host)
	}

	// 解析参数
	query := parsedURL.Query()

	// 检查是否为新 OAuth 格式（包含 code 参数）
	if code := query.Get("code"); code != "" {
		return s.handleOAuthCallback(deepLinkURL, parsedURL, query)
	}

	// 兼容旧的直接传 token 格式
	return s.handleLegacyCallback(deepLinkURL, parsedURL, query)
}

// handleOAuthCallback 处理新的 OAuth 回调格式
func (s *DeepLinkService) handleOAuthCallback(deepLinkURL string, parsedURL *url.URL, query url.Values) error {
	nonce := query.Get("nonce")
	windowId := query.Get("windowId")
	code := query.Get("code")
	errorParam := query.Get("error")

	fmt.Printf("OAuth 回调 - Nonce: %s, WindowId: %s, Code: %s, Error: %s\n",
		nonce, windowId, code, errorParam)

	// 检查错误
	if errorParam != "" {
		if errorParam == "access_denied" {
			fmt.Println("用户取消了登录")
		} else {
			return fmt.Errorf("登录失败: %s", errorParam)
		}
		return nil
	}

	// 验证必要参数
	if nonce == "" {
		return fmt.Errorf("OAuth 回调缺少 nonce 参数")
	}

	if code == "" {
		return fmt.Errorf("OAuth 回调缺少 code 参数")
	}

	fmt.Printf("OAuth 回调验证成功，授权码: %s\n", code)

	// 注意：这里的 code 已经在 AuthService.OpenLoginPage 中处理了
	// Deep Link 回调只是确认流程完成，不需要重复处理

	return nil
}

// handleLegacyCallback 处理旧的兼容格式
func (s *DeepLinkService) handleLegacyCallback(deepLinkURL string, parsedURL *url.URL, query url.Values) error {
	token := query.Get("token")
	userID := query.Get("user_id")
	email := query.Get("email")

	if token == "" {
		return fmt.Errorf("缺少 token 参数")
	}

	if userID == "" {
		return fmt.Errorf("缺少 user_id 参数")
	}

	tokenPreview := token
	if len(token) > 20 {
		tokenPreview = token[:20] + "..."
	}
	fmt.Printf("Legacy 格式解析成功 - Token: %s, UserID: %s, Email: %s\n",
		tokenPreview, userID, email)

	// 保存认证信息
	if err := s.SaveAuthInfo(token, userID, email); err != nil {
		return fmt.Errorf("保存认证信息失败: %w", err)
	}

	// 触发回调
	if s.onAuthCallback != nil {
		s.onAuthCallback(token, userID, email)
	}

	return nil
}

// AuthInfo 认证信息结构
type AuthInfo struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// SaveAuthInfo 保存认证信息到本地文件
func (s *DeepLinkService) SaveAuthInfo(token, userID, email string) error {
	// 获取配置目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}

	configDir := filepath.Join(homeDir, ".code-switch")

	// 确保目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 构建认证信息
	authInfo := AuthInfo{
		Token:  token,
		UserID: userID,
		Email:  email,
	}

	// 序列化为 JSON
	data, err := json.MarshalIndent(authInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化认证信息失败: %w", err)
	}

	// 保存到文件
	authFile := filepath.Join(configDir, "auth.json")
	if err := os.WriteFile(authFile, data, 0600); err != nil {
		return fmt.Errorf("写入认证文件失败: %w", err)
	}

	fmt.Printf("认证信息已保存到: %s\n", authFile)
	return nil
}

// LoadAuthInfo 加载认证信息
func (s *DeepLinkService) LoadAuthInfo() (*AuthInfo, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("获取用户目录失败: %w", err)
	}

	authFile := filepath.Join(homeDir, ".code-switch", "auth.json")

	// 检查文件是否存在
	if _, err := os.Stat(authFile); os.IsNotExist(err) {
		return nil, nil // 文件不存在，返回 nil 不报错
	}

	// 读取文件
	data, err := os.ReadFile(authFile)
	if err != nil {
		return nil, fmt.Errorf("读取认证文件失败: %w", err)
	}

	// 解析 JSON
	var authInfo AuthInfo
	if err := json.Unmarshal(data, &authInfo); err != nil {
		return nil, fmt.Errorf("解析认证信息失败: %w", err)
	}

	return &authInfo, nil
}

// ClearAuthInfo 清除认证信息
func (s *DeepLinkService) ClearAuthInfo() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("获取用户目录失败: %w", err)
	}

	authFile := filepath.Join(homeDir, ".code-switch", "auth.json")

	if err := os.Remove(authFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除认证文件失败: %w", err)
	}

	fmt.Println("认证信息已清除")
	return nil
}
