package llm

import (
	"testing"
)

func TestBuildTranspilePrompt_Python(t *testing.T) {
	code := `print("Hello, World!")`
	prompt := BuildTranspilePrompt(ScriptTypePython, code)

	if prompt == "" {
		t.Error("Prompt should not be empty")
	}

	// Check that prompt contains key elements
	if !contains(prompt, "python") {
		t.Error("Prompt should mention python")
	}

	if !contains(prompt, code) {
		t.Error("Prompt should contain the original code")
	}

	if !contains(prompt, "Go") {
		t.Error("Prompt should mention Go")
	}
}

func TestBuildTranspilePrompt_Shell(t *testing.T) {
	code := `echo "Hello, World!"`
	prompt := BuildTranspilePrompt(ScriptTypeShell, code)

	if prompt == "" {
		t.Error("Prompt should not be empty")
	}

	if !contains(prompt, "shell") {
		t.Error("Prompt should mention shell")
	}

	if !contains(prompt, code) {
		t.Error("Prompt should contain the original code")
	}
}

func TestBuildRefinePrompt(t *testing.T) {
	goCode := `package main

func main() {
	fmt.Println("test")
}
`
	errorMsg := "undefined: fmt"

	prompt := BuildRefinePrompt(goCode, errorMsg)

	if prompt == "" {
		t.Error("Prompt should not be empty")
	}

	if !contains(prompt, errorMsg) {
		t.Error("Prompt should contain the error message")
	}

	if !contains(prompt, goCode) {
		t.Error("Prompt should contain the Go code")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
