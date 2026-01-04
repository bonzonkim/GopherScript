package llm

import (
	"fmt"

	"go.uber.org/zap"
)

// Provider represents supported LLM providers
type Provider string

const (
	ProviderGemini Provider = "gemini"
	ProviderOpenAI Provider = "openai"
	ProviderClaude Provider = "claude"
)

// ValidProviders returns a list of valid provider names
func ValidProviders() []Provider {
	return []Provider{ProviderGemini, ProviderOpenAI, ProviderClaude}
}

// IsValid checks if the provider is valid
func (p Provider) IsValid() bool {
	switch p {
	case ProviderGemini, ProviderOpenAI, ProviderClaude:
		return true
	default:
		return false
	}
}

// NewClient creates an LLM client based on the provider
func NewClient(provider Provider, apiKey string, logger *zap.Logger) (Clienter, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required for provider %s", provider)
	}

	switch provider {
	case ProviderGemini:
		return NewGeminiClient(apiKey, logger), nil
	case ProviderOpenAI:
		return NewOpenAIClient(apiKey, logger), nil
	case ProviderClaude:
		return NewClaudeClient(apiKey, logger), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
