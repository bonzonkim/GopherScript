package generator

import (
	"fmt"
	"go/format"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"go.uber.org/zap"
)

// Generator handles Go code generation and building
type Generator struct {
	logger *zap.Logger
}

// NewGenerator creates a new Generator instance
func NewGenerator(logger *zap.Logger) *Generator {
	return &Generator{
		logger: logger,
	}
}

// GenerateResult contains the result of code generation
type GenerateResult struct {
	GoCode     string
	OutputPath string
	BinaryPath string
}

// Generate formats and writes Go code to a file
func (g *Generator) Generate(goCode string, outputPath string) (*GenerateResult, error) {
	// Clean up markdown code blocks if present
	cleanCode := g.cleanCodeBlock(goCode)

	// Format the Go code
	formatted, err := format.Source([]byte(cleanCode))
	if err != nil {
		g.logger.Warn("Failed to format Go code, using raw code", zap.Error(err))
		formatted = []byte(cleanCode)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Write the Go file
	if err := os.WriteFile(outputPath, formatted, 0644); err != nil {
		return nil, fmt.Errorf("failed to write Go file: %w", err)
	}

	g.logger.Info("Generated Go file", zap.String("path", outputPath))

	return &GenerateResult{
		GoCode:     string(formatted),
		OutputPath: outputPath,
	}, nil
}

// Build compiles the Go code into a static binary
func (g *Generator) Build(goFilePath string, binaryPath string) error {
	g.logger.Info("Building binary", zap.String("source", goFilePath), zap.String("output", binaryPath))

	// Build with static linking flags
	cmd := exec.Command("go", "build",
		"-ldflags", "-s -w",
		"-o", binaryPath,
		goFilePath,
	)

	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %s: %w", string(output), err)
	}

	g.logger.Info("Successfully built binary", zap.String("path", binaryPath))
	return nil
}

// cleanCodeBlock removes markdown code block markers from the code
func (g *Generator) cleanCodeBlock(code string) string {
	// Remove ```go or ``` at the beginning
	code = strings.TrimSpace(code)

	// Pattern to match code blocks
	codeBlockPattern := regexp.MustCompile("(?s)^```(?:go)?\\s*\\n(.*)\\n```$")
	if matches := codeBlockPattern.FindStringSubmatch(code); len(matches) > 1 {
		return matches[1]
	}

	// Also handle case where there's no trailing newline before ```
	codeBlockPattern2 := regexp.MustCompile("(?s)^```(?:go)?\\s*\\n(.*)```$")
	if matches := codeBlockPattern2.FindStringSubmatch(code); len(matches) > 1 {
		return matches[1]
	}

	return code
}

// GetDefaultOutputPath generates a default output path based on input file
func (g *Generator) GetDefaultOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)

	return filepath.Join(dir, name+".go")
}

// GetDefaultBinaryPath generates a default binary path based on Go file
func (g *Generator) GetDefaultBinaryPath(goFilePath string) string {
	dir := filepath.Dir(goFilePath)
	base := filepath.Base(goFilePath)
	name := strings.TrimSuffix(base, ".go")

	return filepath.Join(dir, name)
}
