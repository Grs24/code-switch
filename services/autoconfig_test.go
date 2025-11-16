package services

import (
	appConfig "codeswitch/config"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSystemInfo(t *testing.T) {
	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	info := autoConfig.GetSystemInfo()

	if info["os"] == "" {
		t.Error("OS should not be empty")
	}

	if info["arch"] == "" {
		t.Error("Arch should not be empty")
	}

	if info["shell"] == "" {
		t.Error("Shell should not be empty")
	}

	t.Logf("System Info: %+v", info)
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name      string
		relayAddr string
		expected  string
	}{
		{
			name:      "default port",
			relayAddr: ":18100",
			expected:  "http://127.0.0.1:18100",
		},
		{
			name:      "with http",
			relayAddr: "http://localhost:8080",
			expected:  "http://localhost:8080",
		},
		{
			name:      "with https",
			relayAddr: "https://api.example.com",
			expected:  "https://api.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authService := NewAuthService()
			claudeService := NewClaudeSettingsService(tt.relayAddr)
			codexService := NewCodexSettingsService(tt.relayAddr)
			autoConfig := NewAutoConfigService(authService, claudeService, codexService, tt.relayAddr)

			result := autoConfig.getBaseURL()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestConfigureAllWithoutLogin(t *testing.T) {
	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	result, err := autoConfig.ConfigureAll("test-api-key", "")

	if err == nil {
		t.Error("Should return error when user is not logged in")
	}

	if result.Success {
		t.Error("Success should be false when user is not logged in")
	}

	if len(result.Errors) == 0 {
		t.Error("Should have errors when user is not logged in")
	}
}

func TestConfigureClaudeCode(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	// 模拟登录
	authService.SetUserInfo(`{"token":"test-api-key","email":"test@example.com"}`)

	err := autoConfig.configureClaudeCode("test-api-key", "")
	if err != nil {
		t.Errorf("Configure Claude Code failed: %v", err)
	}

	// 验证配置文件是否创建
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Error("Settings file should be created")
	}

	// 读取并验证配置内容
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		t.Errorf("Failed to read settings file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Settings file should not be empty")
	}

	// 验证 base_url 是否正确
	var config map[string]interface{}
	if err := json.Unmarshal(content, &config); err != nil {
		t.Errorf("Failed to parse settings: %v", err)
	}

	env, ok := config["env"].(map[string]interface{})
	if !ok {
		t.Error("env field should exist")
	}

	baseURL, ok := env["ANTHROPIC_BASE_URL"].(string)
	if !ok {
		t.Error("ANTHROPIC_BASE_URL should exist")
	}

	expectedBaseURL := appConfig.GetAICodeBaseURL()
	if baseURL != expectedBaseURL {
		t.Errorf("ANTHROPIC_BASE_URL should be %s, got: %s", expectedBaseURL, baseURL)
	}

	t.Logf("Claude Code settings: %s", string(content))
}

func TestConfigureCodex(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	// 模拟登录
	authService.SetUserInfo(`{"token":"test-api-key","email":"test@example.com"}`)

	err := autoConfig.configureCodex("test-api-key", "")
	if err != nil {
		t.Errorf("Configure Codex failed: %v", err)
	}

	// 验证配置文件是否创建
	configPath := filepath.Join(tmpDir, ".codex", "config.toml")
	authPath := filepath.Join(tmpDir, ".codex", "auth.json")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}

	if _, err := os.Stat(authPath); os.IsNotExist(err) {
		t.Error("Auth file should be created")
	}

	// 读取并验证配置内容
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Errorf("Failed to read config file: %v", err)
	}

	if len(configContent) == 0 {
		t.Error("Config file should not be empty")
	}

	// 验证 base_url 是否正确
	configStr := string(configContent)
	expectedBaseURL := appConfig.GetAICodeBaseURL()
	if !strings.Contains(configStr, expectedBaseURL) {
		t.Errorf("Config should contain %s, got: %s", expectedBaseURL, configStr)
	}

	authContent, err := os.ReadFile(authPath)
	if err != nil {
		t.Errorf("Failed to read auth file: %v", err)
	}

	if len(authContent) == 0 {
		t.Error("Auth file should not be empty")
	}

	t.Logf("Codex config: %s", string(configContent))
	t.Logf("Codex auth: %s", string(authContent))
}

func TestConfigureClaudeCodeOnly(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	// 模拟登录
	authService.SetUserInfo(`{"token":"user-login-token","email":"test@example.com"}`)

	result, err := autoConfig.ConfigureClaudeCodeOnly("sk-test-api-key", "")
	if err != nil {
		t.Errorf("Configure Claude Code Only failed: %v", err)
	}

	if !result.Success {
		t.Error("Result should be success")
	}

	if result.ClaudeStatus != "✅ 配置成功" {
		t.Errorf("Claude status should be success, got: %s", result.ClaudeStatus)
	}

	if result.CodexStatus != "未配置" {
		t.Errorf("Codex status should be unconfigured, got: %s", result.CodexStatus)
	}

	// 验证配置文件是否创建
	settingsPath := filepath.Join(tmpDir, ".claude", "settings.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		t.Error("Claude settings file should be created")
	}

	// 验证 Codex 配置文件未创建
	codexConfigPath := filepath.Join(tmpDir, ".codex", "config.toml")
	if _, err := os.Stat(codexConfigPath); err == nil {
		t.Error("Codex config file should not be created")
	}
}

func TestConfigureCodexOnly(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	// 模拟登录
	authService.SetUserInfo(`{"token":"user-login-token","email":"test@example.com"}`)

	result, err := autoConfig.ConfigureCodexOnly("sk-test-api-key", "")
	if err != nil {
		t.Errorf("Configure Codex Only failed: %v", err)
	}

	if !result.Success {
		t.Error("Result should be success")
	}

	if result.CodexStatus != "✅ 配置成功" {
		t.Errorf("Codex status should be success, got: %s", result.CodexStatus)
	}

	if result.ClaudeStatus != "未配置" {
		t.Errorf("Claude status should be unconfigured, got: %s", result.ClaudeStatus)
	}

	// 验证配置文件是否创建
	codexConfigPath := filepath.Join(tmpDir, ".codex", "config.toml")
	codexAuthPath := filepath.Join(tmpDir, ".codex", "auth.json")

	if _, err := os.Stat(codexConfigPath); os.IsNotExist(err) {
		t.Error("Codex config file should be created")
	}

	if _, err := os.Stat(codexAuthPath); os.IsNotExist(err) {
		t.Error("Codex auth file should be created")
	}

	// 验证 Claude 配置文件未创建
	claudeSettingsPath := filepath.Join(tmpDir, ".claude", "settings.json")
	if _, err := os.Stat(claudeSettingsPath); err == nil {
		t.Error("Claude settings file should not be created")
	}
}

func TestVerifyConfiguration(t *testing.T) {
	// 创建临时目录用于测试
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	authService := NewAuthService()
	claudeService := NewClaudeSettingsService(":18100")
	codexService := NewCodexSettingsService(":18100")
	autoConfig := NewAutoConfigService(authService, claudeService, codexService, ":18100")

	// 模拟登录
	authService.SetUserInfo(`{"token":"user-login-token","email":"test@example.com"}`)

	// 先配置
	_, err := autoConfig.ConfigureAll("sk-test-api-key", "")
	if err != nil {
		t.Errorf("Configure failed: %v", err)
	}

	// 验证配置
	result, err := autoConfig.VerifyConfiguration()
	if err != nil {
		t.Errorf("Verify configuration failed: %v", err)
	}

	if !result.Success {
		t.Errorf("Verification should succeed, errors: %v", result.Errors)
	}

	if result.ClaudeStatus != "✅ 配置正确" {
		t.Errorf("Claude status should be correct, got: %s", result.ClaudeStatus)
	}

	if result.CodexStatus != "✅ 配置正确" {
		t.Errorf("Codex status should be correct, got: %s", result.CodexStatus)
	}

	t.Logf("Verification result: Claude=%s, Codex=%s", result.ClaudeStatus, result.CodexStatus)
}
