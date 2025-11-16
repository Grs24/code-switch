package services

import (
	appConfig "codeswitch/config"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// AutoConfigService 自动配置服务
type AutoConfigService struct {
	authService           *AuthService
	claudeSettingsService *ClaudeSettingsService
	codexSettingsService  *CodexSettingsService
	relayAddr             string
}

// NewAutoConfigService 创建自动配置服务
func NewAutoConfigService(authService *AuthService, claudeService *ClaudeSettingsService, codexService *CodexSettingsService, relayAddr string) *AutoConfigService {
	return &AutoConfigService{
		authService:           authService,
		claudeSettingsService: claudeService,
		codexSettingsService:  codexService,
		relayAddr:             relayAddr,
	}
}

// AutoConfigResult 自动配置结果
type AutoConfigResult struct {
	Success              bool                `json:"success"`
	ClaudeStatus         string              `json:"claudeStatus"`
	CodexStatus          string              `json:"codexStatus"`
	Errors               []string            `json:"errors"`
	Warnings             []string            `json:"warnings"`
	OS                   string              `json:"os"`
	ConfigPaths          map[string][]string `json:"configPaths"`
	VerificationCommands map[string][]string `json:"verificationCommands"`
}

// DependencyCheckResult 依赖检查结果
type DependencyCheckResult struct {
	Passed   bool       `json:"passed"`
	NodeJS   NodeJSInfo `json:"nodejs"`
	Errors   []string   `json:"errors"`
	Warnings []string   `json:"warnings"`
}

// NodeJSInfo Node.js 信息
type NodeJSInfo struct {
	Installed        bool   `json:"installed"`
	Version          string `json:"version"`
	MeetsRequirement bool   `json:"meetsRequirement"`
	MinVersion       string `json:"minVersion"`
}

// ConfigStep 配置步骤
type ConfigStep struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"` // pending, running, success, error
	Error       string   `json:"error,omitempty"`
	Details     []string `json:"details,omitempty"`
}

// ConfigureAll 自动配置 Claude Code 和 Codex
func (s *AutoConfigService) ConfigureAll(apiKey string, baseURL string) (*AutoConfigResult, error) {
	// 获取 base URL，如果没有提供则使用默认值
	var apiBaseURL string
	if baseURL != "" {
		apiBaseURL = baseURL
	} else {
		apiBaseURL = appConfig.GetAICodeBaseURL()
	}
	result := &AutoConfigResult{
		Success:              true,
		Errors:               []string{},
		Warnings:             []string{},
		OS:                   runtime.GOOS,
		ConfigPaths:          make(map[string][]string),
		VerificationCommands: make(map[string][]string),
	}

	// 检查用户是否已登录
	if !s.authService.IsLogin() {
		result.Success = false
		result.Errors = append(result.Errors, "用户未登录")
		return result, fmt.Errorf("用户未登录")
	}

	// 验证 API Key
	if apiKey == "" {
		result.Success = false
		result.Errors = append(result.Errors, "API Key 不能为空")
		return result, fmt.Errorf("API Key 不能为空")
	}

	home, _ := os.UserHomeDir()

	// 配置 Claude Code
	if err := s.configureClaudeCode(apiKey, apiBaseURL); err != nil {
		result.ClaudeStatus = fmt.Sprintf("配置失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Claude Code: %v", err))
		result.Success = false
	} else {
		result.ClaudeStatus = "✅ 配置成功"
		result.ConfigPaths["claude"] = []string{
			filepath.Join(home, ".claude", "settings.json"),
		}
		result.VerificationCommands["claude"] = []string{
			"claude --version",
		}
	}

	// 配置 Codex
	if err := s.configureCodex(apiKey, apiBaseURL); err != nil {
		result.CodexStatus = fmt.Sprintf("配置失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Codex: %v", err))
		result.Success = false
	} else {
		result.CodexStatus = "✅ 配置成功"
		result.ConfigPaths["codex"] = []string{
			filepath.Join(home, ".codex", "config.toml"),
			filepath.Join(home, ".codex", "auth.json"),
		}
		result.VerificationCommands["codex"] = []string{
			"codex -V",
		}
	}

	return result, nil
}

// ConfigureClaudeCodeOnly 仅配置 Claude Code
func (s *AutoConfigService) ConfigureClaudeCodeOnly(apiKey string, baseURL string) (*AutoConfigResult, error) {
	// 获取 base URL，如果没有提供则使用默认值
	var apiBaseURL string
	if baseURL != "" {
		apiBaseURL = baseURL
	} else {
		apiBaseURL = appConfig.GetAICodeBaseURL()
	}
	result := &AutoConfigResult{
		Success:              true,
		Errors:               []string{},
		Warnings:             []string{},
		OS:                   runtime.GOOS,
		ConfigPaths:          make(map[string][]string),
		VerificationCommands: make(map[string][]string),
	}

	// 检查用户是否已登录
	if !s.authService.IsLogin() {
		result.Success = false
		result.Errors = append(result.Errors, "用户未登录")
		return result, fmt.Errorf("用户未登录")
	}

	// 验证 API Key
	if apiKey == "" {
		result.Success = false
		result.Errors = append(result.Errors, "API Key 不能为空")
		return result, fmt.Errorf("API Key 不能为空")
	}

	home, _ := os.UserHomeDir()

	// 配置 Claude Code
	if err := s.configureClaudeCode(apiKey, apiBaseURL); err != nil {
		result.ClaudeStatus = fmt.Sprintf("配置失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Claude Code: %v", err))
		result.Success = false
	} else {
		result.ClaudeStatus = "✅ 配置成功"
		result.ConfigPaths["claude"] = []string{
			filepath.Join(home, ".claude", "settings.json"),
		}
		result.VerificationCommands["claude"] = []string{
			"claude --version",
		}
	}

	result.CodexStatus = "未配置"

	return result, nil
}

// ConfigureCodexOnly 仅配置 Codex
func (s *AutoConfigService) ConfigureCodexOnly(apiKey string, baseURL string) (*AutoConfigResult, error) {
	// 获取 base URL，如果没有提供则使用默认值
	var apiBaseURL string
	if baseURL != "" {
		apiBaseURL = baseURL
	} else {
		apiBaseURL = appConfig.GetAICodeBaseURL()
	}
	result := &AutoConfigResult{
		Success:              true,
		Errors:               []string{},
		Warnings:             []string{},
		OS:                   runtime.GOOS,
		ConfigPaths:          make(map[string][]string),
		VerificationCommands: make(map[string][]string),
	}

	// 检查用户是否已登录
	if !s.authService.IsLogin() {
		result.Success = false
		result.Errors = append(result.Errors, "用户未登录")
		return result, fmt.Errorf("用户未登录")
	}

	// 验证 API Key
	if apiKey == "" {
		result.Success = false
		result.Errors = append(result.Errors, "API Key 不能为空")
		return result, fmt.Errorf("API Key 不能为空")
	}

	result.ClaudeStatus = "未配置"

	home, _ := os.UserHomeDir()

	// 配置 Codex
	if err := s.configureCodex(apiKey, apiBaseURL); err != nil {
		result.CodexStatus = fmt.Sprintf("配置失败: %v", err)
		result.Errors = append(result.Errors, fmt.Sprintf("Codex: %v", err))
		result.Success = false
	} else {
		result.CodexStatus = "✅ 配置成功"
		result.ConfigPaths["codex"] = []string{
			filepath.Join(home, ".codex", "config.toml"),
			filepath.Join(home, ".codex", "auth.json"),
		}
		result.VerificationCommands["codex"] = []string{
			"codex -V",
		}
	}

	return result, nil
}

// VerifyConfiguration 验证配置是否正确
func (s *AutoConfigService) VerifyConfiguration() (*AutoConfigResult, error) {
	result := &AutoConfigResult{
		Success: true,
		Errors:  []string{},
		OS:      runtime.GOOS,
	}

	home, err := os.UserHomeDir()
	if err != nil {
		result.Success = false
		result.Errors = append(result.Errors, fmt.Sprintf("无法获取用户目录: %v", err))
		return result, err
	}

	// 验证 Claude Code 配置
	claudeSettingsPath := filepath.Join(home, ".claude", "settings.json")
	if _, err := os.Stat(claudeSettingsPath); os.IsNotExist(err) {
		result.ClaudeStatus = "配置文件不存在"
		result.Success = false
	} else {
		// 读取并验证配置内容
		content, err := os.ReadFile(claudeSettingsPath)
		if err != nil {
			result.ClaudeStatus = fmt.Sprintf("无法读取配置: %v", err)
			result.Success = false
		} else {
			var config map[string]interface{}
			if err := json.Unmarshal(content, &config); err != nil {
				result.ClaudeStatus = "配置格式错误"
				result.Success = false
			} else {
				env, ok := config["env"].(map[string]interface{})
				if !ok || env["ANTHROPIC_BASE_URL"] == nil || env["ANTHROPIC_AUTH_TOKEN"] == nil {
					result.ClaudeStatus = "配置内容不完整"
					result.Success = false
				} else {
					baseURL, _ := env["ANTHROPIC_BASE_URL"].(string)
					expectedBaseURL := appConfig.GetAICodeBaseURL()
					if baseURL != expectedBaseURL {
						result.ClaudeStatus = fmt.Sprintf("⚠️ base_url 不正确: %s (期望: %s)", baseURL, expectedBaseURL)
						result.Success = false
					} else {
						result.ClaudeStatus = "✅ 配置正确"
					}
				}
			}
		}
	}

	// 验证 Codex 配置
	codexConfigPath := filepath.Join(home, ".codex", "config.toml")
	codexAuthPath := filepath.Join(home, ".codex", "auth.json")

	configExists := true
	authExists := true

	if _, err := os.Stat(codexConfigPath); os.IsNotExist(err) {
		configExists = false
	}

	if _, err := os.Stat(codexAuthPath); os.IsNotExist(err) {
		authExists = false
	}

	if !configExists && !authExists {
		result.CodexStatus = "配置文件不存在"
		result.Success = false
	} else if !configExists {
		result.CodexStatus = "config.toml 不存在"
		result.Success = false
	} else if !authExists {
		result.CodexStatus = "auth.json 不存在"
		result.Success = false
	} else {
		// 验证 auth.json 内容
		authContent, err := os.ReadFile(codexAuthPath)
		if err != nil {
			result.CodexStatus = fmt.Sprintf("无法读取 auth.json: %v", err)
			result.Success = false
		} else {
			var authConfig map[string]string
			if err := json.Unmarshal(authContent, &authConfig); err != nil {
				result.CodexStatus = "auth.json 格式错误"
				result.Success = false
			} else if authConfig["OPENAI_API_KEY"] == "" {
				result.CodexStatus = "API Key 未设置"
				result.Success = false
			} else {
				result.CodexStatus = "✅ 配置正确"
			}
		}
	}

	return result, nil
}

// configureClaudeCode 配置 Claude Code
func (s *AutoConfigService) configureClaudeCode(apiKey string, baseURL string) error {
	// 如果没有提供 baseURL，使用默认值
	if baseURL == "" {
		baseURL = appConfig.GetAICodeBaseURL()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	claudeDir := filepath.Join(home, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")

	// 创建目录
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 备份现有配置
	if _, err := os.Stat(settingsPath); err == nil {
		backupPath := filepath.Join(claudeDir, "settings.backup.json")
		content, _ := os.ReadFile(settingsPath)
		_ = os.WriteFile(backupPath, content, 0o600)
	}

	// 生成配置 - 使用用户指定的 API 服务器地址
	settings := map[string]interface{}{
		"env": map[string]string{
			"ANTHROPIC_BASE_URL":   baseURL,
			"ANTHROPIC_AUTH_TOKEN": apiKey,
		},
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("生成配置失败: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0o600); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}

	// 配置环境变量到 shell 配置文件
	if err := s.configureShellEnv(apiKey, baseURL); err != nil {
		fmt.Printf("[AutoConfig] 警告: 环境变量配置失败: %v\n", err)
		// 不返回错误，因为 settings.json 已经配置成功
	}

	fmt.Printf("[AutoConfig] Claude Code 配置成功: %s\n", settingsPath)
	return nil
}

// configureShellEnv 配置 shell 环境变量
func (s *AutoConfigService) configureShellEnv(apiKey string, baseURL string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	// 确定 shell 配置文件
	var shellRC string
	if runtime.GOOS == "windows" {
		// Windows 不需要配置 shell 环境变量
		return nil
	}

	// 检测使用的 shell
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		shellRC = filepath.Join(home, ".zshrc")
	} else {
		shellRC = filepath.Join(home, ".bashrc")
	}

	// 读取现有配置
	content, err := os.ReadFile(shellRC)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("读取 shell 配置失败: %w", err)
	}

	existingContent := string(content)

	// 检查是否已经配置过
	baseURLLine := fmt.Sprintf(`export ANTHROPIC_BASE_URL="%s"`, baseURL)
	tokenLine := fmt.Sprintf(`export ANTHROPIC_AUTH_TOKEN="%s"`, apiKey)

	// 移除旧的配置（如果存在）
	lines := strings.Split(existingContent, "\n")
	var newLines []string
	for _, line := range lines {
		// 跳过旧的 ANTHROPIC 配置
		if strings.Contains(line, "ANTHROPIC_BASE_URL") || strings.Contains(line, "ANTHROPIC_AUTH_TOKEN") {
			continue
		}
		newLines = append(newLines, line)
	}

	// 添加新配置
	newLines = append(newLines, "")
	newLines = append(newLines, "# Claude Code 配置 (由 Code Switch 自动生成)")
	newLines = append(newLines, baseURLLine)
	newLines = append(newLines, tokenLine)

	// 写回文件
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(shellRC, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("写入 shell 配置失败: %w", err)
	}

	fmt.Printf("[AutoConfig] Shell 环境变量配置成功: %s\n", shellRC)
	return nil
}

// configureCodex 配置 Codex
func (s *AutoConfigService) configureCodex(apiKey string, baseURL string) error {
	// 如果没有提供 baseURL，使用默认值
	if baseURL == "" {
		baseURL = appConfig.GetAICodeBaseURL()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户目录: %w", err)
	}

	codexDir := filepath.Join(home, ".codex")
	configPath := filepath.Join(codexDir, "config.toml")
	authPath := filepath.Join(codexDir, "auth.json")

	// 创建目录
	if err := os.MkdirAll(codexDir, 0o755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	// 备份现有配置
	if _, err := os.Stat(configPath); err == nil {
		backupPath := filepath.Join(codexDir, "config.backup.toml")
		content, _ := os.ReadFile(configPath)
		_ = os.WriteFile(backupPath, content, 0o600)
	}

	// 生成 config.toml - 使用用户指定的 API 服务器地址
	config := map[string]interface{}{
		"disable_response_storage": true,
		"preferred_auth_method":    "apikey",
		"model":                    "gpt-5-codex",
		"model_provider":           "aicoding",
		"model_providers": map[string]interface{}{
			"aicoding": map[string]interface{}{
				"name":                 "aicoding",
				"base_url":             baseURL,
				"env_key":              "OPENAI_API_KEY",
				"wire_api":             "responses",
				"requires_openai_auth": false,
			},
		},
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("生成配置失败: %w", err)
	}

	// 移除 [model_providers] 标题行
	cleaned := removeModelProvidersHeader(data)
	if err := os.WriteFile(configPath, cleaned, 0o600); err != nil {
		return fmt.Errorf("写入配置失败: %w", err)
	}

	// 备份现有 auth.json
	if _, err := os.Stat(authPath); err == nil {
		backupPath := filepath.Join(codexDir, "auth.backup.json")
		content, _ := os.ReadFile(authPath)
		_ = os.WriteFile(backupPath, content, 0o600)
	}

	// 生成 auth.json
	authData := map[string]string{
		"OPENAI_API_KEY": apiKey,
	}

	authJSON, err := json.MarshalIndent(authData, "", "  ")
	if err != nil {
		return fmt.Errorf("生成 auth.json 失败: %w", err)
	}

	if err := os.WriteFile(authPath, authJSON, 0o600); err != nil {
		return fmt.Errorf("写入 auth.json 失败: %w", err)
	}

	fmt.Printf("[AutoConfig] Codex 配置成功: %s\n", configPath)
	return nil
}

// GetSystemInfo 获取系统信息
func (s *AutoConfigService) GetSystemInfo() map[string]string {
	info := map[string]string{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}

	// 检测 shell
	if runtime.GOOS != "windows" {
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/bash"
		}
		info["shell"] = shell
	} else {
		info["shell"] = "powershell"
	}

	return info
}

// CheckCodexInstalled 检查 Codex 是否已安装
func (s *AutoConfigService) CheckCodexInstalled() bool {
	cmd := exec.Command("codex", "-V")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// GetToolVersions 获取工具版本信息
func (s *AutoConfigService) GetToolVersions() map[string]string {
	versions := make(map[string]string)

	// 获取 Claude Code 版本
	claudeCmd := exec.Command("claude", "--version")
	if output, err := claudeCmd.Output(); err == nil {
		versions["claude"] = strings.TrimSpace(string(output))
	} else {
		versions["claude"] = "未安装或未找到"
	}

	// 获取 Codex 版本
	codexCmd := exec.Command("codex", "-V")
	if output, err := codexCmd.Output(); err == nil {
		versions["codex"] = strings.TrimSpace(string(output))
	} else {
		versions["codex"] = "未安装或未找到"
	}

	return versions
}

// CheckConfigurationStatus 检查配置状态
func (s *AutoConfigService) CheckConfigurationStatus() (*AutoConfigResult, error) {
	result := &AutoConfigResult{
		Success:              false,
		Errors:               []string{},
		Warnings:             []string{},
		OS:                   runtime.GOOS,
		ConfigPaths:          make(map[string][]string),
		VerificationCommands: make(map[string][]string),
	}

	home, err := os.UserHomeDir()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("无法获取用户目录: %v", err))
		return result, err
	}

	claudeConfigured := false
	codexConfigured := false

	// 检查 Claude Code 配置
	claudeSettingsPath := filepath.Join(home, ".claude", "settings.json")
	if _, err := os.Stat(claudeSettingsPath); err == nil {
		// 读取配置文件
		content, err := os.ReadFile(claudeSettingsPath)
		if err == nil {
			var config map[string]interface{}
			if err := json.Unmarshal(content, &config); err == nil {
				env, ok := config["env"].(map[string]interface{})
				if ok && env["ANTHROPIC_BASE_URL"] != nil && env["ANTHROPIC_AUTH_TOKEN"] != nil {
					claudeConfigured = true
					result.ClaudeStatus = "✅ 已配置"
					result.ConfigPaths["claude"] = []string{claudeSettingsPath}
					result.VerificationCommands["claude"] = []string{"claude --version"}
				}
			}
		}
	}

	if !claudeConfigured {
		result.ClaudeStatus = "未配置"
	}

	// 检查 Codex 配置
	codexConfigPath := filepath.Join(home, ".codex", "config.toml")
	codexAuthPath := filepath.Join(home, ".codex", "auth.json")

	configExists := false
	authExists := false

	if _, err := os.Stat(codexConfigPath); err == nil {
		configExists = true
	}

	if _, err := os.Stat(codexAuthPath); err == nil {
		authExists = true
	}

	if configExists && authExists {
		// 验证 auth.json 内容
		authContent, err := os.ReadFile(codexAuthPath)
		if err == nil {
			var authConfig map[string]string
			if err := json.Unmarshal(authContent, &authConfig); err == nil {
				if authConfig["OPENAI_API_KEY"] != "" {
					codexConfigured = true
					result.CodexStatus = "✅ 已配置"
					result.ConfigPaths["codex"] = []string{codexConfigPath, codexAuthPath}
					result.VerificationCommands["codex"] = []string{"codex -V"}
				}
			}
		}
	}

	if !codexConfigured {
		result.CodexStatus = "未配置"
	}

	// 如果两个都配置了，标记为成功
	if claudeConfigured && codexConfigured {
		result.Success = true
	} else if claudeConfigured || codexConfigured {
		result.Success = true
		result.Warnings = append(result.Warnings, "部分工具已配置")
	}

	return result, nil
}

// CheckDependencies 检查依赖
func (s *AutoConfigService) CheckDependencies() (*DependencyCheckResult, error) {
	result := &DependencyCheckResult{
		Passed:   true,
		Errors:   []string{},
		Warnings: []string{},
		NodeJS: NodeJSInfo{
			MinVersion: "v16.0.0",
		},
	}

	// 检查 Node.js
	nodeInfo, err := s.checkNodeJS()
	if err != nil {
		result.Passed = false
		result.Errors = append(result.Errors, fmt.Sprintf("检查 Node.js 失败: %v", err))
	} else {
		result.NodeJS = *nodeInfo
		if !nodeInfo.Installed {
			result.Passed = false
			result.Errors = append(result.Errors, "Node.js 未安装，请先安装 Node.js v16.0.0 或更高版本")
		} else if !nodeInfo.MeetsRequirement {
			result.Passed = false
			result.Errors = append(result.Errors, fmt.Sprintf("Node.js 版本过低 (当前: %s, 需要: %s+)", nodeInfo.Version, nodeInfo.MinVersion))
		}
	}

	return result, nil
}

// checkNodeJS 检查 Node.js 安装和版本
func (s *AutoConfigService) checkNodeJS() (*NodeJSInfo, error) {
	info := &NodeJSInfo{
		MinVersion: "v16.0.0",
	}

	// 执行 node --version
	cmd := exec.Command("node", "--version")
	output, err := cmd.Output()
	if err != nil {
		info.Installed = false
		return info, nil
	}

	version := strings.TrimSpace(string(output))
	info.Installed = true
	info.Version = version

	// 比较版本
	info.MeetsRequirement = s.compareVersion(version, info.MinVersion) >= 0

	return info, nil
}

// compareVersion 比较版本号
func (s *AutoConfigService) compareVersion(current, required string) int {
	// 移除 'v' 前缀
	current = strings.TrimPrefix(current, "v")
	required = strings.TrimPrefix(required, "v")

	// 分割版本号
	currentParts := strings.Split(current, ".")
	requiredParts := strings.Split(required, ".")

	// 比较每个部分
	maxLen := len(currentParts)
	if len(requiredParts) > maxLen {
		maxLen = len(requiredParts)
	}

	for i := 0; i < maxLen; i++ {
		var currentPart, requiredPart int
		if i < len(currentParts) {
			fmt.Sscanf(currentParts[i], "%d", &currentPart)
		}
		if i < len(requiredParts) {
			fmt.Sscanf(requiredParts[i], "%d", &requiredPart)
		}

		if currentPart > requiredPart {
			return 1
		} else if currentPart < requiredPart {
			return -1
		}
	}

	return 0
}

// getBaseURL 获取基础 URL
func (s *AutoConfigService) getBaseURL() string {
	addr := strings.TrimSpace(s.relayAddr)
	if addr == "" {
		addr = ":18100"
	}
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		return addr
	}
	host := addr
	if strings.HasPrefix(host, ":") {
		host = "127.0.0.1" + host
	}
	if !strings.Contains(host, "://") {
		host = "http://" + host
	}
	return host
}

// removeModelProvidersHeader 移除 [model_providers] 标题行
func removeModelProvidersHeader(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	result := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "[model_providers]" {
			continue
		}
		result = append(result, line)
	}
	return []byte(strings.Join(result, "\n"))
}

// GetDefaultBaseURL 获取默认的 Base URL
func (s *AutoConfigService) GetDefaultBaseURL() string {
	return appConfig.GetAICodeBaseURL()
}
