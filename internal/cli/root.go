package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/bonzonkim/gopher-script/config"
	"github.com/bonzonkim/gopher-script/internal/handler"
	"github.com/bonzonkim/gopher-script/internal/llm"
	"github.com/bonzonkim/gopher-script/internal/logger"
	"github.com/spf13/cobra"
)

var (
	outputPath string
	binaryPath string
	build      bool
	verbose    bool
	provider   string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gopherscript [file]",
		Short: "Converts a script file to Go.",
		Long: `GopherScript is a tool to convert Python or Shell scripts into idiomatic Go code.

It uses LLM to intelligently transpile your scripts into
standalone Go binaries.

Supported LLM providers:
  - gemini  (Google Gemini, default)
  - openai  (OpenAI GPT-4o)
  - claude  (Anthropic Claude)

Examples:
  gopherscript script.py                       # Convert using default provider (Gemini)
  gopherscript script.py --provider openai     # Convert using OpenAI GPT
  gopherscript script.py --provider claude     # Convert using Anthropic Claude
  gopherscript script.sh -o output.go          # Convert Shell script with custom output
  gopherscript script.py --build               # Convert and build binary
  gopherscript script.py -o main.go -b bin     # Convert with custom output and binary path`,
		Args: cobra.MinimumNArgs(1),
		RunE: runTranspile,
	}

	// Add flags
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for the generated Go file")
	cmd.Flags().StringVarP(&binaryPath, "binary", "b", "", "Output path for the compiled binary (requires --build)")
	cmd.Flags().BoolVar(&build, "build", false, "Build the generated Go code into a binary")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	cmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider to use (gemini, openai, claude)")

	return cmd
}

func runTranspile(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Load configuration
	cfg := config.NewConfig()

	// Determine provider (CLI flag overrides env var)
	selectedProvider := cfg.Provider
	if provider != "" {
		selectedProvider = provider
	}

	// Validate provider
	llmProvider := llm.Provider(selectedProvider)
	if !llmProvider.IsValid() {
		validProviders := make([]string, len(llm.ValidProviders()))
		for i, p := range llm.ValidProviders() {
			validProviders[i] = string(p)
		}
		return fmt.Errorf("invalid provider '%s'. Valid providers: %s", selectedProvider, strings.Join(validProviders, ", "))
	}

	// Get API key for selected provider
	apiKey := cfg.GetAPIKey(selectedProvider)
	if apiKey == "" {
		var envVar string
		switch llmProvider {
		case llm.ProviderGemini:
			envVar = "GEMINI_API_KEY (or API_KEY)"
		case llm.ProviderOpenAI:
			envVar = "OPENAI_API_KEY"
		case llm.ProviderClaude:
			envVar = "ANTHROPIC_API_KEY"
		}
		return fmt.Errorf("%s environment variable is not set for provider '%s'", envVar, selectedProvider)
	}

	// Determine environment for logger
	env := cfg.Env
	if env == "" {
		env = "dev"
	}
	if verbose {
		env = "dev"
	}

	// Initialize logger
	log := logger.NewLogger(env)
	defer log.Logger.Sync()

	// Create handler
	h, err := handler.NewHandler(log.Logger, llmProvider, apiKey)
	if err != nil {
		return fmt.Errorf("failed to initialize handler: %w", err)
	}

	// Run transpilation
	opts := handler.TranspileOptions{
		InputPath:  inputPath,
		OutputPath: outputPath,
		Build:      build,
		BinaryPath: binaryPath,
	}

	result, err := h.Transpile(opts)
	if err != nil {
		return fmt.Errorf("transpilation failed: %w", err)
	}

	// Print success message
	fmt.Fprintf(os.Stdout, "✅ Successfully transpiled: %s (using %s)\n", inputPath, selectedProvider)
	fmt.Fprintf(os.Stdout, "   Go file: %s\n", result.OutputPath)

	if result.BinaryPath != "" {
		fmt.Fprintf(os.Stdout, "   Binary:  %s\n", result.BinaryPath)
	}

	return nil
}
