package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ScriptType represents the type of script file
type ScriptType string

const (
	ScriptTypePython  ScriptType = "python"
	ScriptTypeShell   ScriptType = "shell"
	ScriptTypeUnknown ScriptType = "unknown"
)

// ParseResult contains the parsed script information
type ParseResult struct {
	FilePath   string
	FileName   string
	ScriptType ScriptType
	Content    string
}

// Parser handles script file parsing
type Parser struct{}

// NewParser creates a new Parser instance
func NewParser() *Parser {
	return &Parser{}
}

// Parse reads and analyzes a script file
func (p *Parser) Parse(filePath string) (*ParseResult, error) {
	// Check if file exists
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s: %w", filePath, err)
		}
		return nil, fmt.Errorf("failed to stat file: %s: %w", filePath, err)
	}

	// Check if it's a regular file
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a file: %s", filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %s: %w", filePath, err)
	}

	// Detect script type
	scriptType := p.detectScriptType(filePath, string(content))

	return &ParseResult{
		FilePath:   filePath,
		FileName:   filepath.Base(filePath),
		ScriptType: scriptType,
		Content:    string(content),
	}, nil
}

// detectScriptType determines the script type based on extension and content
func (p *Parser) detectScriptType(filePath string, content string) ScriptType {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check by extension first
	switch ext {
	case ".py":
		return ScriptTypePython
	case ".sh", ".bash", ".zsh":
		return ScriptTypeShell
	}

	// Check by shebang
	lines := strings.Split(content, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if strings.HasPrefix(firstLine, "#!") {
			if strings.Contains(firstLine, "python") {
				return ScriptTypePython
			}
			if strings.Contains(firstLine, "bash") || strings.Contains(firstLine, "sh") || strings.Contains(firstLine, "zsh") {
				return ScriptTypeShell
			}
		}
	}

	return ScriptTypeUnknown
}

// IsSupportedType checks if the script type is supported for transpilation
func (p *Parser) IsSupportedType(scriptType ScriptType) bool {
	return scriptType == ScriptTypePython || scriptType == ScriptTypeShell
}
