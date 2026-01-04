package config

import "os"

type Config struct {
	Provider        string
	GeminiAPIKey    string
	OpenAIAPIKey    string
	AnthropicAPIKey string
	Env             string
}

func NewConfig() *Config {
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "gemini" // default provider
	}

	geminiKey := os.Getenv("GEMINI_API_KEY")
	// Backward compatibility: use legacy API_KEY for Gemini
	if geminiKey == "" {
		geminiKey = os.Getenv("API_KEY")
	}

	return &Config{
		Provider:        provider,
		GeminiAPIKey:    geminiKey,
		OpenAIAPIKey:    os.Getenv("OPENAI_API_KEY"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		Env:             os.Getenv("ENV"),
	}
}

// GetAPIKey returns the API key for the specified provider
func (c *Config) GetAPIKey(provider string) string {
	switch provider {
	case "gemini":
		return c.GeminiAPIKey
	case "openai":
		return c.OpenAIAPIKey
	case "claude":
		return c.AnthropicAPIKey
	default:
		return ""
	}
}
