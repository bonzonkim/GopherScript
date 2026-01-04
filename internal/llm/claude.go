package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

const (
	claudeAPIURL = "https://api.anthropic.com/v1/messages"
	claudeModel  = "claude-sonnet-4-20250514"
)

// ClaudeClient implements Clienter for Anthropic Claude API
type ClaudeClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewClaudeClient creates a new Anthropic Claude API client
func NewClaudeClient(apiKey string, logger *zap.Logger) *ClaudeClient {
	return &ClaudeClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logger,
	}
}

// ClaudeRequest represents the request body for Claude API
type ClaudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []ClaudeMessage `json:"messages"`
}

// ClaudeMessage represents a message in the request
type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ClaudeResponse represents the response from Claude API
type ClaudeResponse struct {
	ID      string               `json:"id"`
	Type    string               `json:"type"`
	Content []ClaudeContentBlock `json:"content"`
	Error   *ClaudeError         `json:"error,omitempty"`
}

// ClaudeContentBlock represents a content block in the response
type ClaudeContentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ClaudeError represents an error from Claude API
type ClaudeError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// Generate sends a prompt to Claude API and returns the response
func (c *ClaudeClient) Generate(prompt string) (string, error) {
	c.logger.Debug("Sending request to Claude API")

	reqBody := ClaudeRequest{
		Model:     claudeModel,
		MaxTokens: 8192,
		Messages: []ClaudeMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, claudeAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if claudeResp.Error != nil {
		return "", fmt.Errorf("API error [%s]: %s", claudeResp.Error.Type, claudeResp.Error.Message)
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Claude API")
	}

	// Find the first text block in the response
	var result string
	for _, block := range claudeResp.Content {
		if block.Type == "text" {
			result = block.Text
			break
		}
	}

	if result == "" {
		return "", fmt.Errorf("no text content in Claude API response")
	}

	c.logger.Debug("Received response from Claude API", zap.Int("length", len(result)))

	return result, nil
}
