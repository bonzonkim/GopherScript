package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"
)

func TestGenerator_Generate(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	g := NewGenerator(logger)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.go")

	goCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}
`

	result, err := g.Generate(goCode, outputPath)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if result.OutputPath != outputPath {
		t.Errorf("OutputPath mismatch: got %s, expected %s", result.OutputPath, outputPath)
	}

	// Verify file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}

	// Verify content
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "fmt.Println") {
		t.Error("Output file content is incorrect")
	}
}

func TestGenerator_Generate_WithMarkdownCodeBlock(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	g := NewGenerator(logger)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "output.go")

	// Code wrapped in markdown code block
	goCode := "```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}\n```"

	result, err := g.Generate(goCode, outputPath)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify markdown was stripped
	if strings.Contains(result.GoCode, "```") {
		t.Error("Markdown code block markers should be stripped")
	}
}

func TestGenerator_cleanCodeBlock(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	g := NewGenerator(logger)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with go code block",
			input:    "```go\npackage main\n```",
			expected: "package main",
		},
		{
			name:     "with plain code block",
			input:    "```\npackage main\n```",
			expected: "package main",
		},
		{
			name:     "no code block",
			input:    "package main",
			expected: "package main",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := g.cleanCodeBlock(tt.input)
			if strings.TrimSpace(result) != strings.TrimSpace(tt.expected) {
				t.Errorf("cleanCodeBlock() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestGenerator_GetDefaultOutputPath(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	g := NewGenerator(logger)

	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/script.py", "/path/to/script.go"},
		{"/path/to/script.sh", "/path/to/script.go"},
		{"./test.py", "test.go"},
	}

	for _, tt := range tests {
		result := g.GetDefaultOutputPath(tt.input)
		if result != tt.expected {
			t.Errorf("GetDefaultOutputPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestGenerator_GetDefaultBinaryPath(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	g := NewGenerator(logger)

	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/script.go", "/path/to/script"},
		{"./test.go", "test"},
	}

	for _, tt := range tests {
		result := g.GetDefaultBinaryPath(tt.input)
		if result != tt.expected {
			t.Errorf("GetDefaultBinaryPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}
