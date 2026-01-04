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
	geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"
)

// Clienter defines the interface for LLM clients
type Clienter interface {
	Generate(prompt string) (string, error)
}

// GeminiClient implements Clienter for Google Gemini API
type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
	logger     *zap.Logger
}

// NewGeminiClient creates a new Gemini API client
func NewGeminiClient(apiKey string, logger *zap.Logger) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		logger: logger,
	}
}

// GeminiRequest represents the request body for Gemini API
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

// Content represents a content block in the request
type Content struct {
	Parts []Part `json:"parts"`
}

// Part represents a part of the content
type Part struct {
	Text string `json:"text"`
}

// GeminiResponse represents the response from Gemini API
type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
	Error      *APIError   `json:"error,omitempty"`
}

// Candidate represents a candidate response
type Candidate struct {
	Content Content `json:"content"`
}

// APIError represents an error from the API
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Generate sends a prompt to Gemini API and returns the response
func (c *GeminiClient) Generate(prompt string) (string, error) {
	c.logger.Debug("Sending request to Gemini API")

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: prompt},
				},
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s?key=%s", geminiAPIURL, c.apiKey)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var geminiResp GeminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if geminiResp.Error != nil {
		return "", fmt.Errorf("API error [%d]: %s", geminiResp.Error.Code, geminiResp.Error.Message)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	result := geminiResp.Candidates[0].Content.Parts[0].Text
	c.logger.Debug("Received response from Gemini API", zap.Int("length", len(result)))

	return result, nil
}
