package llm

import "fmt"

// ScriptType represents the type of script being converted
type ScriptType string

const (
	ScriptTypePython ScriptType = "python"
	ScriptTypeShell  ScriptType = "shell"
)

// BuildTranspilePrompt creates a prompt for transpiling script code to Go
func BuildTranspilePrompt(scriptType ScriptType, code string) string {
	baseInstruction := `You are an expert Go programmer. Convert the following %s script to idiomatic Go code.

Requirements:
1. Use proper error handling with wrapped errors
2. Follow Go naming conventions (camelCase for unexported, PascalCase for exported)
3. Add necessary imports
4. Include a main function that can be compiled into a standalone binary
5. Add brief comments explaining the logic
6. Use the standard library when possible
7. Return ONLY the Go code without any explanation or markdown formatting

%s script to convert:
` + "```\n%s\n```"

	return fmt.Sprintf(baseInstruction, scriptType, scriptType, code)
}

// BuildRefinePrompt creates a prompt for refining generated Go code
func BuildRefinePrompt(goCode string, errorMessage string) string {
	return fmt.Sprintf(`The following Go code has a compilation error. Please fix it and return only the corrected Go code without any explanation or markdown formatting.

Error message:
%s

Go code to fix:
`+"```go\n%s\n```", errorMessage, goCode)
}
