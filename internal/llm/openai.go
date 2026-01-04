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
	openAIAPIURL = "https://api.openai.com/v1/chat/completions"
	openAIModel  = "gpt-4o"
)

// OpenAIClient implements Clienter for OpenAI API
type OpenAIClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewOpenAIClient creates a new OpenAI API client
func NewOpenAIClient(apiKey string, logger *zap.Logger) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
		logger: logger,
	}
}

// OpenAIChatRequest represents the request body for OpenAI Chat API
type OpenAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []OpenAIChatMessage `json:"messages"`
}

// OpenAIChatMessage represents a message in the chat
type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatResponse represents the response from OpenAI Chat API
type OpenAIChatResponse struct {
	ID      string              `json:"id"`
	Choices []OpenAIChatChoice  `json:"choices"`
	Error   *OpenAIErrorWrapper `json:"error,omitempty"`
}

// OpenAIChatChoice represents a choice in the response
type OpenAIChatChoice struct {
	Message      OpenAIChatMessage `json:"message"`
	FinishReason string            `json:"finish_reason"`
}

// OpenAIErrorWrapper represents an error from OpenAI API
type OpenAIErrorWrapper struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// Generate sends a prompt to OpenAI API and returns the response
func (c *OpenAIClient) Generate(prompt string) (string, error) {
	c.logger.Debug("Sending request to OpenAI API")

	reqBody := OpenAIChatRequest{
		Model: openAIModel,
		Messages: []OpenAIChatMessage{
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

	req, err := http.NewRequest(http.MethodPost, openAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResp OpenAIChatResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if openAIResp.Error != nil {
		return "", fmt.Errorf("API error [%s]: %s", openAIResp.Error.Type, openAIResp.Error.Message)
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from OpenAI API")
	}

	result := openAIResp.Choices[0].Message.Content
	c.logger.Debug("Received response from OpenAI API", zap.Int("length", len(result)))

	return result, nil
}
