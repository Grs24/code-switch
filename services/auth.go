package services

import (
	"codeswitch/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	// API 请求相关常量
	apiKeyEndpoint    = "/aicoding0011/xapi/user/token"
	headerAccessToken = "accesstoken"
	headerContentType = "Content-Type"
	headerClientType  = "clienttype"
	headerLanguage    = "language"

	// 默认值
	defaultClientType = "1"
	defaultLanguage   = "zh"
)

var (
	// 错误定义
	ErrNotLoggedIn     = fmt.Errorf("用户未登录")
	ErrProviderNotInit = fmt.Errorf("ProviderService 未初始化")
	ErrNoAPIKey        = fmt.Errorf("未找到 API Key")
)

// UserInfo 用户信息
type UserInfo struct {
	Token    string         `json:"token"`
	UserID   string         `json:"userId"`
	Email    string         `json:"email"`
	Mobile   string         `json:"mobile"`
	UserInfo map[string]any `json:"userInfo"` // 完整的后端用户数据
}

// AuthService 认证服务
type AuthService struct {
	mu              sync.RWMutex
	userInfo        *UserInfo
	isLogin         bool
	providerService *ProviderService
	claudeSettings  *ClaudeSettingsService
	codexSettings   *CodexSettingsService
	httpClient      *http.Client
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{
		isLogin: false,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetProviderService 设置 ProviderService（用于依赖注入）
func (s *AuthService) SetProviderService(ps *ProviderService) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providerService = ps
}

// SetClaudeSettings 设置 ClaudeSettingsService（用于依赖注入）
func (s *AuthService) SetClaudeSettings(cs *ClaudeSettingsService) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.claudeSettings = cs
}

// SetCodexSettings 设置 CodexSettingsService（用于依赖注入）
func (s *AuthService) SetCodexSettings(cs *CodexSettingsService) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.codexSettings = cs
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
	// 同时输出到控制台和日志文件
	logMsg := fmt.Sprintf("[AuthService] ========== 开始处理用户登录 ==========\n")
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	logMsg = fmt.Sprintf("[AuthService] 接收到的用户信息 JSON: %s\n", userInfoJSON)
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	s.mu.Lock()
	defer s.mu.Unlock()

	var userInfo UserInfo
	if err := json.Unmarshal([]byte(userInfoJSON), &userInfo); err != nil {
		fmt.Printf("[AuthService] ❌ 解析用户信息失败: %v\n", err)
		return fmt.Errorf("解析用户信息失败: %w", err)
	}

	s.userInfo = &userInfo
	s.isLogin = true

	tokenPreview := userInfo.Token
	if len(tokenPreview) > 10 {
		tokenPreview = tokenPreview[:10] + "..."
	}
	logMsg = fmt.Sprintf("[AuthService] ✅ 用户登录成功: %s (Token: %s)\n", userInfo.Email, tokenPreview)
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	// 登录成功后，自动创建 0011 供应商卡片
	logMsg = "[AuthService] 开始创建默认供应商...\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if err := s.ensureDefaultProvider(); err != nil {
		logMsg = fmt.Sprintf("[AuthService] 创建默认供应商失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		// 不阻断登录流程，只记录错误
	}

	// 登录成功后，自动启用代理
	logMsg = "[AuthService] 开始启用代理...\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if err := s.enableProxies(); err != nil {
		logMsg = fmt.Sprintf("[AuthService] 启用代理失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		// 不阻断登录流程，只记录错误
	}

	logMsg = "[AuthService] ========== 登录处理完成 ==========\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	return nil
}

// ApiResponse 通用 API 响应结构
type ApiResponse struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message,omitempty"`
	Msg     string          `json:"msg,omitempty"`
}

// ApiKeyResponse API Key 响应结构
type ApiKeyResponse struct {
	List []ApiKeyData `json:"list"`
}

// ApiKeyData API Key 数据
type ApiKeyData struct {
	ID     int    `json:"id"`
	Token  string `json:"token"`
	Status int    `json:"status"`
}

// fetchUserApiKey 获取用户的 API Key
func (s *AuthService) fetchUserApiKey() (string, error) {
	if s.userInfo == nil {
		return "", ErrNotLoggedIn
	}

	apiURL := config.GetAPIURL()
	url := apiURL + apiKeyEndpoint

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置认证 token（使用与前端相同的请求头）
	req.Header.Set(headerAccessToken, s.userInfo.Token)
	req.Header.Set(headerContentType, "application/json")
	req.Header.Set(headerClientType, defaultClientType)
	req.Header.Set(headerLanguage, defaultLanguage)

	logMsg := fmt.Sprintf("[AuthService] 正在请求 API Key: %s\n", url)
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logMsg = fmt.Sprintf("[AuthService] HTTP 请求失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logMsg = fmt.Sprintf("[AuthService] 读取响应失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	logMsg = fmt.Sprintf("[AuthService] API 响应状态: %d\n", resp.StatusCode)
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	logMsg = fmt.Sprintf("[AuthService] API 响应内容: %s\n", string(body))
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败: status=%d, body=%s", resp.StatusCode, string(body))
	}

	// 先解析外层响应
	var apiResp ApiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		logMsg = fmt.Sprintf("[AuthService] JSON 解析失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", fmt.Errorf("解析响应失败: %w, body=%s", err, string(body))
	}

	logMsg = fmt.Sprintf("[AuthService] API 响应 code: %d\n", apiResp.Code)
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	// 检查业务状态码
	if apiResp.Code != 0 && apiResp.Code != 200 {
		errMsg := apiResp.Message
		if errMsg == "" {
			errMsg = apiResp.Msg
		}
		logMsg = fmt.Sprintf("[AuthService] API 业务错误: code=%d, msg=%s\n", apiResp.Code, errMsg)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", fmt.Errorf("API 错误: %s", errMsg)
	}

	// 解析 data 字段
	var apiKeyResp ApiKeyResponse
	if err := json.Unmarshal(apiResp.Data, &apiKeyResp); err != nil {
		logMsg = fmt.Sprintf("[AuthService] 解析 data 字段失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", fmt.Errorf("解析 data 失败: %w", err)
	}

	logMsg = fmt.Sprintf("[AuthService] 解析到的 API Key 列表长度: %d\n", len(apiKeyResp.List))
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if len(apiKeyResp.List) == 0 {
		logMsg = "[AuthService] ⚠️ API Key 列表为空\n"
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return "", ErrNoAPIKey
	}

	logMsg = fmt.Sprintf("[AuthService] 获取到 API Key: %s...\n", apiKeyResp.List[0].Token[:10])
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	return apiKeyResp.List[0].Token, nil
}

// ensureDefaultProvider 确保 0011 供应商存在（Claude Code 和 Codex）
func (s *AuthService) ensureDefaultProvider() error {
	logMsg := "[AuthService] 检查 ProviderService...\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if s.providerService == nil {
		logMsg = "[AuthService] ❌ ProviderService 未初始化\n"
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return ErrProviderNotInit
	}

	// 尝试获取用户的 API Key（两种类型共用）
	logMsg = "[AuthService] 开始获取 API Key...\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	apiKey, err := s.fetchUserApiKey()
	if err != nil {
		logMsg = fmt.Sprintf("[AuthService] 获取 API Key 失败: %v，将使用空 API Key\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		apiKey = ""
	} else {
		keyPreview := apiKey
		if len(keyPreview) > 10 {
			keyPreview = keyPreview[:10] + "..."
		}
		logMsg = fmt.Sprintf("[AuthService] 成功获取 API Key: %s\n", keyPreview)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
	}

	// 为 Claude Code 和 Codex 创建 0011 供应商
	providerTypes := []string{"claude", "codex"}
	for _, providerType := range providerTypes {
		providerID := Get0011ProviderID(providerType)
		logMsg = fmt.Sprintf("[AuthService] 为 %s 创建供应商 (ID=%d)...\n", providerType, providerID)
		fmt.Print(logMsg)
		s.writeLog(logMsg)

		if err := s.ensureProviderForType(providerType, providerID, apiKey); err != nil {
			logMsg = fmt.Sprintf("[AuthService] 创建 %s 0011 供应商失败: %v\n", providerType, err)
			fmt.Print(logMsg)
			s.writeLog(logMsg)
		}
	}

	logMsg = "[AuthService] 默认供应商创建完成\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	return nil
}

// ensureProviderForType 为指定类型创建或更新 0011 供应商
func (s *AuthService) ensureProviderForType(providerType string, providerID int, apiKey string) error {
	providers, err := s.providerService.LoadProviders(providerType)
	if err != nil {
		return fmt.Errorf("加载 %s 供应商列表失败: %w", providerType, err)
	}

	// 查找是否已存在 0011 供应商
	var existing0011 *Provider
	var existing0011Index int = -1
	for i, p := range providers {
		if p.ID == providerID {
			existing0011 = &providers[i]
			existing0011Index = i
			break
		}
	}

	needsSave := false

	if existing0011 == nil {
		// 不存在，创建新的 0011 供应商
		defaultProvider := CreateProvider0011(providerType, apiKey)
		providers = append([]Provider{defaultProvider}, providers...)
		needsSave = true

		if apiKey != "" {
			fmt.Printf("[AuthService] 已自动创建并启用 %s 0011 供应商卡片 (ID=%d)\n", providerType, providerID)
		} else {
			fmt.Printf("[AuthService] 已自动创建 %s 0011 供应商卡片 (ID=%d，未启用)\n", providerType, providerID)
		}
	} else if apiKey != "" && existing0011.APIKey == "" {
		// 已存在但 API Key 为空，更新 API Key 并启用
		existing0011.APIKey = apiKey
		existing0011.Enabled = true
		providers[existing0011Index] = *existing0011
		needsSave = true

		fmt.Printf("[AuthService] 已更新 %s 0011 供应商的 API Key 并启用 (ID=%d)\n", providerType, providerID)
	} else if apiKey != "" && existing0011.APIKey != "" {
		fmt.Printf("[AuthService] %s 0011 供应商已存在且有 API Key，跳过更新 (ID=%d)\n", providerType, providerID)
	} else {
		fmt.Printf("[AuthService] %s 0011 供应商已存在 (ID=%d)\n", providerType, providerID)
	}

	// 保存更改
	if needsSave {
		if err := s.providerService.SaveProviders(providerType, providers); err != nil {
			return fmt.Errorf("保存 %s 供应商失败: %w", providerType, err)
		}
	}

	return nil
}

// enableProxies 启用 Claude Code 和 Codex 代理
func (s *AuthService) enableProxies() error {
	// 启用 Claude Code 代理
	if s.claudeSettings != nil {
		if err := s.enableClaudeProxy(); err != nil {
			fmt.Printf("[AuthService] Claude Code 代理配置失败: %v\n", err)
		}
	}

	// 启用 Codex 代理
	if s.codexSettings != nil {
		if err := s.enableCodexProxy(); err != nil {
			fmt.Printf("[AuthService] Codex 代理配置失败: %v\n", err)
		}
	}

	return nil
}

// enableClaudeProxy 启用 Claude Code 代理
func (s *AuthService) enableClaudeProxy() error {
	status, err := s.claudeSettings.ProxyStatus()
	if err != nil {
		return fmt.Errorf("检查代理状态失败: %w", err)
	}

	if status.Enabled {
		fmt.Println("[AuthService] Claude Code 代理已启用，跳过")
		return nil
	}

	if err := s.claudeSettings.EnableProxy(); err != nil {
		return fmt.Errorf("启用代理失败: %w", err)
	}

	fmt.Println("[AuthService] 已自动启用 Claude Code 代理")
	return nil
}

// enableCodexProxy 启用 Codex 代理
func (s *AuthService) enableCodexProxy() error {
	status, err := s.codexSettings.ProxyStatus()
	if err != nil {
		return fmt.Errorf("检查代理状态失败: %w", err)
	}

	if status.Enabled {
		fmt.Println("[AuthService] Codex 代理已启用，跳过")
		return nil
	}

	if err := s.codexSettings.EnableProxy(); err != nil {
		return fmt.Errorf("启用代理失败: %w", err)
	}

	fmt.Println("[AuthService] 已自动启用 Codex 代理")
	return nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo() (*UserInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.isLogin || s.userInfo == nil {
		return nil, ErrNotLoggedIn
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

// RefreshProviderApiKey 刷新 0011 供应商的 API Key
// 用于用户在界面上更新 API Key 后同步到供应商卡片
func (s *AuthService) RefreshProviderApiKey() error {
	logMsg := "[AuthService] ========== 开始刷新 0011 供应商 API Key ==========\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	if !s.isLogin {
		return ErrNotLoggedIn
	}

	if s.providerService == nil {
		return ErrProviderNotInit
	}

	// 获取最新的 API Key
	logMsg = "[AuthService] 获取最新的 API Key...\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	apiKey, err := s.fetchUserApiKey()
	if err != nil {
		logMsg = fmt.Sprintf("[AuthService] 获取 API Key 失败: %v\n", err)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
		return fmt.Errorf("获取 API Key 失败: %w", err)
	}

	// 更新 Claude Code 和 Codex 的 0011 供应商
	providerTypes := []string{"claude", "codex"}
	for _, providerType := range providerTypes {
		providerID := Get0011ProviderID(providerType)

		logMsg = fmt.Sprintf("[AuthService] 更新 %s 0011 供应商 API Key (ID=%d)...\n", providerType, providerID)
		fmt.Print(logMsg)
		s.writeLog(logMsg)

		providers, err := s.providerService.LoadProviders(providerType)
		if err != nil {
			logMsg = fmt.Sprintf("[AuthService] 加载 %s 供应商列表失败: %v\n", providerType, err)
			fmt.Print(logMsg)
			s.writeLog(logMsg)
			continue
		}

		// 查找 0011 供应商
		found := false
		for i, p := range providers {
			if p.ID == providerID {
				providers[i].APIKey = apiKey
				providers[i].Enabled = true
				found = true
				break
			}
		}

		if !found {
			logMsg = fmt.Sprintf("[AuthService] 未找到 %s 0011 供应商，跳过\n", providerType)
			fmt.Print(logMsg)
			s.writeLog(logMsg)
			continue
		}

		// 保存更新
		if err := s.providerService.SaveProviders(providerType, providers); err != nil {
			logMsg = fmt.Sprintf("[AuthService] 保存 %s 供应商失败: %v\n", providerType, err)
			fmt.Print(logMsg)
			s.writeLog(logMsg)
			continue
		}

		logMsg = fmt.Sprintf("[AuthService] ✅ %s 0011 供应商 API Key 已更新\n", providerType)
		fmt.Print(logMsg)
		s.writeLog(logMsg)
	}

	logMsg = "[AuthService] ========== API Key 刷新完成 ==========\n"
	fmt.Print(logMsg)
	s.writeLog(logMsg)

	return nil
}

// writeLog 写入日志到文件
func (s *AuthService) writeLog(msg string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	logFile := filepath.Join(home, ".code-switch", "auth.log")
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(time.Now().Format("2006-01-02 15:04:05") + " " + msg)
}
