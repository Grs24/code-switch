package services

// Provider0011Config 0011 供应商配置
type Provider0011Config struct {
	ClaudeID int
	CodexID  int
	Name     string
	APIURL   string
	Site     string
	Icon     string
	Tint     string
	Accent   string
}

// DefaultProvider0011Config 默认的 0011 供应商配置
var DefaultProvider0011Config = Provider0011Config{
	ClaudeID: 100,
	CodexID:  200,
	Name:     "0011",
	APIURL:   "https://aicoding.2233.ai",
	Site:     "https://0011.ai",
	Icon:     "aicoding",
	Tint:     "rgba(10, 132, 255, 0.14)",
	Accent:   "#0aff5cff",
}

// CreateProvider0011 创建 0011 供应商
func CreateProvider0011(providerType string, apiKey string) Provider {
	config := DefaultProvider0011Config

	providerID := config.ClaudeID
	if providerType == "codex" {
		providerID = config.CodexID
	}

	enabled := apiKey != ""

	return Provider{
		ID:      providerID,
		Name:    config.Name,
		APIURL:  config.APIURL,
		APIKey:  apiKey,
		Site:    config.Site,
		Icon:    config.Icon,
		Tint:    config.Tint,
		Accent:  config.Accent,
		Enabled: enabled,
	}
}

// Is0011Provider 检查是否为 0011 供应商
func Is0011Provider(id int) bool {
	config := DefaultProvider0011Config
	return id == config.ClaudeID || id == config.CodexID
}

// Get0011ProviderID 获取 0011 供应商的 ID
func Get0011ProviderID(providerType string) int {
	config := DefaultProvider0011Config
	if providerType == "codex" {
		return config.CodexID
	}
	return config.ClaudeID
}
