package config

import (
	"fmt"
	"os"
)

// Environment 环境类型
type Environment string

const (
	// EnvironmentDevelopment 本地开发环境
	EnvironmentDevelopment Environment = "development"
	// EnvironmentQA QA 测试环境
	EnvironmentQA Environment = "qa"
	// EnvironmentProduction 生产环境
	EnvironmentProduction Environment = "production"
)

// Config 应用配置
type Config struct {
	// 环境
	Environment Environment
	// Web 登录页 URL
	LoginURL string
	// API 基础 URL
	APIURL string
}

var (
	// 当前环境配置
	Current *Config
)

// 初始化配置
func init() {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development" // 默认开发环境（方便本地调试）
	}

	Current = &Config{
		Environment: Environment(env),
	}

	// 根据环境设置 URL
	switch Current.Environment {
	case EnvironmentDevelopment:
		Current.LoginURL = "http://localhost:3000/desktop-login"
		Current.APIURL = "https://wildai-qa.bihu.it" // 本地开发使用 QA 后端
	case EnvironmentQA:
		Current.LoginURL = "https://0011ai-qa.bihu.it/desktop-login"
		Current.APIURL = "https://wildai-qa.bihu.it"
	case EnvironmentProduction:
		Current.LoginURL = "https://0011.ai/desktop-login"
		Current.APIURL = "https://api.wildai.net"
	default:
		// 默认开发环境
		Current.LoginURL = "http://localhost:3000/desktop-login"
		Current.APIURL = "https://wildai-qa.bihu.it" // 本地开发使用 QA 后端
	}

	// 打印配置信息（方便调试）
	fmt.Printf("[Config] Environment: %s\n", Current.Environment)
	fmt.Printf("[Config] LoginURL: %s\n", Current.LoginURL)
	fmt.Printf("[Config] APIURL: %s\n", Current.APIURL)
}

// GetLoginURL 获取登录页 URL
func GetLoginURL() string {
	return Current.LoginURL
}

// GetAPIURL 获取 API URL
func GetAPIURL() string {
	return Current.APIURL
}

// IsDevelopment 是否开发环境
func IsDevelopment() bool {
	return Current.Environment == EnvironmentDevelopment
}

// IsQA 是否 QA 环境
func IsQA() bool {
	return Current.Environment == EnvironmentQA
}

// IsProduction 是否生产环境
func IsProduction() bool {
	return Current.Environment == EnvironmentProduction
}
