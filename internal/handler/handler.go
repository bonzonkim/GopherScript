package handler

import (
	"fmt"

	"github.com/bonzonkim/gopher-script/internal/generator"
	"github.com/bonzonkim/gopher-script/internal/llm"
	"github.com/bonzonkim/gopher-script/internal/parser"
	"go.uber.org/zap"
)

// Handler orchestrates the transpilation process
type Handler struct {
	Logger    *zap.Logger
	LLMClient llm.Clienter
	Parser    *parser.Parser
	Generator *generator.Generator
}

// NewHandler creates a new Handler with all dependencies
func NewHandler(logger *zap.Logger, provider llm.Provider, apiKey string) (*Handler, error) {
	client, err := llm.NewClient(provider, apiKey, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client: %w", err)
	}

	return &Handler{
		Logger:    logger,
		LLMClient: client,
		Parser:    parser.NewParser(),
		Generator: generator.NewGenerator(logger),
	}, nil
}

// TranspileOptions contains options for the transpile operation
type TranspileOptions struct {
	InputPath  string
	OutputPath string
	Build      bool
	BinaryPath string
}

// TranspileResult contains the result of transpilation
type TranspileResult struct {
	GoCode     string
	OutputPath string
	BinaryPath string
}

// RequestLLM sends the script code to LLM for transpilation
func (h *Handler) RequestLLM(scriptType parser.ScriptType, code string) (string, error) {
	h.Logger.Info("Requesting LLM for transpilation",
		zap.String("scriptType", string(scriptType)),
		zap.Int("codeLength", len(code)))

	// Build appropriate prompt based on script type
	var llmScriptType llm.ScriptType
	switch scriptType {
	case parser.ScriptTypePython:
		llmScriptType = llm.ScriptTypePython
	case parser.ScriptTypeShell:
		llmScriptType = llm.ScriptTypeShell
	default:
		return "", fmt.Errorf("unsupported script type: %s", scriptType)
	}

	prompt := llm.BuildTranspilePrompt(llmScriptType, code)

	goCode, err := h.LLMClient.Generate(prompt)
	if err != nil {
		return "", fmt.Errorf("LLM request failed: %w", err)
	}

	h.Logger.Info("LLM transpilation completed", zap.Int("resultLength", len(goCode)))
	return goCode, nil
}

// Transpile converts a script file to Go and optionally builds it
func (h *Handler) Transpile(opts TranspileOptions) (*TranspileResult, error) {
	h.Logger.Info("Starting transpilation", zap.String("input", opts.InputPath))

	// Step 1: Parse the input file
	parsed, err := h.Parser.Parse(opts.InputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input file: %w", err)
	}

	// Check if script type is supported
	if !h.Parser.IsSupportedType(parsed.ScriptType) {
		return nil, fmt.Errorf("unsupported script type: %s (file: %s)", parsed.ScriptType, parsed.FileName)
	}

	h.Logger.Info("Parsed script file",
		zap.String("file", parsed.FileName),
		zap.String("type", string(parsed.ScriptType)))

	// Step 2: Request LLM for transpilation
	goCode, err := h.RequestLLM(parsed.ScriptType, parsed.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to transpile: %w", err)
	}

	// Step 3: Determine output path
	outputPath := opts.OutputPath
	if outputPath == "" {
		outputPath = h.Generator.GetDefaultOutputPath(opts.InputPath)
	}

	// Step 4: Generate Go file
	genResult, err := h.Generator.Generate(goCode, outputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Go file: %w", err)
	}

	result := &TranspileResult{
		GoCode:     genResult.GoCode,
		OutputPath: genResult.OutputPath,
	}

	// Step 5: Build if requested
	if opts.Build {
		binaryPath := opts.BinaryPath
		if binaryPath == "" {
			binaryPath = h.Generator.GetDefaultBinaryPath(outputPath)
		}

		if err := h.Generator.Build(outputPath, binaryPath); err != nil {
			return nil, fmt.Errorf("failed to build binary: %w", err)
		}

		result.BinaryPath = binaryPath
	}

	h.Logger.Info("Transpilation completed successfully",
		zap.String("output", result.OutputPath),
		zap.String("binary", result.BinaryPath))

	return result, nil
}
