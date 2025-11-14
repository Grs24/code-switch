package services

import (
	"codeswitch/config"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AuthService 认证服务 - 提供给前端的接口
type AuthService struct {
	app             *application.App
	deepLinkService *DeepLinkService
	oauthServer     *OAuthServer
	securityManager *SecurityManager
	isLoginFlow     bool
}

// NewAuthService 创建认证服务
func NewAuthService(deepLinkService *DeepLinkService) *AuthService {
	return &AuthService{
		deepLinkService: deepLinkService,
		oauthServer:     NewOAuthServer(),
		securityManager: NewSecurityManager(),
	}
}

// SetApp 设置应用实例
func (s *AuthService) SetApp(app *application.App) {
	s.app = app
}

// OpenLoginPage 打开登录页面
func (s *AuthService) OpenLoginPage() error {
	// 防止重复登录流程
	if s.isLoginFlow {
		return fmt.Errorf("登录流程正在进行中")
	}

	s.isLoginFlow = true
	defer func() {
		s.isLoginFlow = false
	}()

	// 启动本地 OAuth 服务器
	callbackURL, err := s.oauthServer.Start()
	if err != nil {
		return fmt.Errorf("启动 OAuth 服务器失败: %w", err)
	}

	// 确保服务器会被关闭
	defer func() {
		if err := s.oauthServer.Stop(); err != nil {
			log.Printf("停止 OAuth 服务器失败: %v", err)
		}
	}()

	// 生成安全参数
	nonce, err := s.securityManager.GenerateNonce(s.oauthServer.GetWindowId())
	if err != nil {
		return fmt.Errorf("生成安全参数失败: %w", err)
	}

	// 生成 PKCE 参数
	pkce, err := s.securityManager.GeneratePKCE()
	if err != nil {
		return fmt.Errorf("生成 PKCE 参数失败: %w", err)
	}

	// 构造登录 URL
	params := url.Values{}
	params.Add("callback", callbackURL)
	params.Add("state", nonce)
	params.Add("window_id", s.oauthServer.GetWindowId())
	params.Add("code_challenge", pkce.CodeChallenge)
	params.Add("code_challenge_method", pkce.Method)

	loginURL := fmt.Sprintf("%s?%s", config.Current.LoginURL, params.Encode())
	log.Printf("打开登录页面: %s\n", loginURL)

	// 在系统默认浏览器中打开
	if err := s.openBrowser(loginURL); err != nil {
		return fmt.Errorf("打开浏览器失败: %w", err)
	}

	// 等待 OAuth 回调结果（5分钟超时）
	result, err := s.oauthServer.WaitForResult(5 * time.Minute)
	if err != nil {
		return fmt.Errorf("等待登录结果失败: %w", err)
	}

	// 处理登录结果
	if !result.Success {
		if result.Error == "access_denied" {
			return fmt.Errorf("用户取消了登录")
		}
		return fmt.Errorf("登录失败: %s", result.Error)
	}

	// 验证 state 参数
	if err := s.securityManager.ValidateNonce(result.State, s.oauthServer.GetWindowId()); err != nil {
		return fmt.Errorf("安全验证失败: %w", err)
	}

	// 标记 nonce 为已使用
	s.securityManager.MarkNonceUsed(result.State)

	log.Printf("登录成功，收到授权码: %s\n", result.Code)

	// 使用授权码交换 token
	tokenInfo, err := s.exchangeCodeForToken(result.Code)
	if err != nil {
		return fmt.Errorf("交换 token 失败: %w", err)
	}

	// 保存认证信息
	if err := s.deepLinkService.SaveAuthInfo(tokenInfo.Token, tokenInfo.UserID, tokenInfo.Email); err != nil {
		return fmt.Errorf("保存认证信息失败: %w", err)
	}

	// 通知前端登录成功
	if s.app != nil {
		s.app.Event.Emit("auth:login-success", map[string]string{
			"user_id": tokenInfo.UserID,
			"email":   tokenInfo.Email,
		})
	}

	log.Printf("用户登录成功: %s (%s)\n", tokenInfo.Email, tokenInfo.UserID)
	return nil
}

// GetAuthInfo 获取认证信息
func (s *AuthService) GetAuthInfo() (*AuthInfo, error) {
	return s.deepLinkService.LoadAuthInfo()
}

// Logout 登出
func (s *AuthService) Logout() error {
	if err := s.deepLinkService.ClearAuthInfo(); err != nil {
		return err
	}

	// 通知前端登出成功
	if s.app != nil {
		s.app.Event.Emit("auth:logout")
	}

	return nil
}

// IsLoggedIn 检查是否已登录
func (s *AuthService) IsLoggedIn() bool {
	authInfo, err := s.deepLinkService.LoadAuthInfo()
	if err != nil {
		return false
	}
	return authInfo != nil && authInfo.Token != ""
}

// GetUserEmail 获取用户邮箱
func (s *AuthService) GetUserEmail() string {
	authInfo, err := s.deepLinkService.LoadAuthInfo()
	if err != nil || authInfo == nil {
		return ""
	}
	return authInfo.Email
}

// GetUserID 获取用户ID
func (s *AuthService) GetUserID() string {
	authInfo, err := s.deepLinkService.LoadAuthInfo()
	if err != nil || authInfo == nil {
		return ""
	}
	return authInfo.UserID
}

// TokenInfo token 信息
type TokenInfo struct {
	Token  string
	UserID string
	Email  string
}

// openBrowser 打开系统默认浏览器
func (s *AuthService) openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}
	return err
}

// exchangeCodeForToken 使用授权码交换访问令牌
func (s *AuthService) exchangeCodeForToken(code string) (*TokenInfo, error) {
	// 构造 API URL
	tokenURL := fmt.Sprintf("%s/auth/exchange-code", config.Current.APIURL)

	// 这里应该使用 HTTP 客户端调用你的后端 API
	// 实际实现时需要：
	// 1. 发送 POST 请求到 tokenURL
	// 2. 包含 code, code_verifier, client_id 等参数
	// 3. 验证响应并解析 token 信息

	/*
	client := &http.Client{Timeout: 30 * time.Second}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://127.0.0.1:52779/callback")
	data.Set("client_id", "your-client-id")
	data.Set("code_verifier", pkce.CodeVerifier) // PKCE 验证

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 token API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token API 返回错误状态: %d", resp.StatusCode)
	}

	var result struct {
		Token  string `json:"access_token"`
		UserID string `json:"user_id"`
		Email  string `json:"email"`
		Error  string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 token 响应失败: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("token API 返回错误: %s", result.Error)
	}

	return &TokenInfo{
		Token:  result.Token,
		UserID: result.UserID,
		Email:  result.Email,
	}, nil
	*/

	// 临时模拟实现（需要替换为实际的 API 调用）
	log.Printf("模拟调用 token API: %s (code: %s)\n", tokenURL, code)

	// 验证 code 格式（基本安全检查）
	if len(code) < 8 {
		return nil, fmt.Errorf("授权码格式无效")
	}

	// 模拟返回 token 信息
	return &TokenInfo{
		Token:  "secure-jwt-token-" + s.securityManager.GenerateRandomString(32),
		UserID: "user-" + code[:min(8, len(code))],
		Email:  "user@example.com",
	}, nil
}

// Cleanup 清理资源
func (s *AuthService) Cleanup() {
	if s.oauthServer != nil {
		s.oauthServer.Stop()
	}
}
