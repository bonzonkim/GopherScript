package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	// Create a temporary Python file for testing
	tmpDir := t.TempDir()
	pythonFile := filepath.Join(tmpDir, "test.py")
	pythonContent := `#!/usr/bin/env python3
print("Hello, World!")
`
	if err := os.WriteFile(pythonFile, []byte(pythonContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	p := NewParser()

	result, err := p.Parse(pythonFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.ScriptType != ScriptTypePython {
		t.Errorf("Expected ScriptTypePython, got %s", result.ScriptType)
	}

	if result.Content != pythonContent {
		t.Errorf("Content mismatch")
	}

	if result.FileName != "test.py" {
		t.Errorf("Expected filename 'test.py', got %s", result.FileName)
	}
}

func TestParser_Parse_ShellScript(t *testing.T) {
	tmpDir := t.TempDir()
	shellFile := filepath.Join(tmpDir, "test.sh")
	shellContent := `#!/bin/bash
echo "Hello, World!"
`
	if err := os.WriteFile(shellFile, []byte(shellContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	p := NewParser()

	result, err := p.Parse(shellFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.ScriptType != ScriptTypeShell {
		t.Errorf("Expected ScriptTypeShell, got %s", result.ScriptType)
	}
}

func TestParser_Parse_FileNotExists(t *testing.T) {
	p := NewParser()

	_, err := p.Parse("/nonexistent/file.py")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestParser_Parse_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	p := NewParser()

	_, err := p.Parse(tmpDir)
	if err == nil {
		t.Error("Expected error for directory, got nil")
	}
}

func TestParser_detectScriptType_ByShebang(t *testing.T) {
	tmpDir := t.TempDir()

	// File with no extension but Python shebang
	noExtFile := filepath.Join(tmpDir, "script")
	content := `#!/usr/bin/python3
print("test")
`
	if err := os.WriteFile(noExtFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	p := NewParser()
	result, err := p.Parse(noExtFile)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if result.ScriptType != ScriptTypePython {
		t.Errorf("Expected ScriptTypePython from shebang, got %s", result.ScriptType)
	}
}

func TestParser_IsSupportedType(t *testing.T) {
	p := NewParser()

	tests := []struct {
		scriptType ScriptType
		expected   bool
	}{
		{ScriptTypePython, true},
		{ScriptTypeShell, true},
		{ScriptTypeUnknown, false},
	}

	for _, tt := range tests {
		result := p.IsSupportedType(tt.scriptType)
		if result != tt.expected {
			t.Errorf("IsSupportedType(%s) = %v, expected %v", tt.scriptType, result, tt.expected)
		}
	}
}
