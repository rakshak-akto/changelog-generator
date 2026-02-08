package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rakshaksatsangi/changelog-generator/pkg/config"
	"github.com/rakshaksatsangi/changelog-generator/pkg/generator"
	"github.com/rakshaksatsangi/changelog-generator/pkg/github"
	"github.com/rakshaksatsangi/changelog-generator/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	cfg     *config.Config
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "changelog-generator",
	Short: "Generate human-readable changelogs from Git commits using AI",
	Long: `Changelog Generator analyzes Git commits and generates structured,
human-readable changelogs using OpenAI's language models.`,
	Version: version,
}

var generateCmd = &cobra.Command{
	Use:   "generate [from]..[to]",
	Short: "Generate a changelog for a commit range",
	Long: `Generate a changelog from a range of commits.

Examples:
  changelog-generator generate v1.0.0..v1.1.0
  changelog-generator generate --owner=facebook --repo=react v18.2.0..v18.3.0
  changelog-generator generate abc123..HEAD
  changelog-generator generate --output=RELEASE.md v1.0.0..HEAD
  changelog-generator generate --show-scores v1.0.0..v1.1.0
  changelog-generator generate --min-score=7.0 v1.0.0..v1.1.0`,
	Args: cobra.ExactArgs(1),
	RunE: runGenerate,
}

func init() {
	// Load config
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Add generate command
	rootCmd.AddCommand(generateCmd)

	// Flags for generate command
	generateCmd.Flags().StringVar(&cfg.RepoOwner, "owner", cfg.RepoOwner, "Repository owner (required)")
	generateCmd.Flags().StringVar(&cfg.RepoName, "repo", cfg.RepoName, "Repository name (required)")
	generateCmd.Flags().StringVar(&cfg.OutputPath, "output", cfg.OutputPath, "Output file path")
	generateCmd.Flags().StringVar(&cfg.OpenAIModel, "model", cfg.OpenAIModel, "OpenAI model to use")
	generateCmd.Flags().BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "Verbose output")
	generateCmd.Flags().BoolVar(&cfg.IncludeAuthors, "include-authors", cfg.IncludeAuthors, "Include commit authors")
	generateCmd.Flags().BoolVar(&cfg.IncludeDates, "include-dates", cfg.IncludeDates, "Include commit dates")
	generateCmd.Flags().BoolVar(&cfg.ShowScores, "show-scores", cfg.ShowScores, "Show importance scores for each commit")
	generateCmd.Flags().Float64Var(&cfg.MinScore, "min-score", cfg.MinScore, "Minimum importance score to include (0-10)")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Parse commit range
	commitRange := args[0]
	parts := strings.Split(commitRange, "..")
	if len(parts) != 2 {
		return fmt.Errorf("invalid commit range format, expected 'from..to', got '%s'", commitRange)
	}
	from, to := parts[0], parts[1]

	if from == "" || to == "" {
		return fmt.Errorf("both 'from' and 'to' refs must be specified")
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("Changelog Generator v%s\n", version)
		fmt.Printf("Repository: %s/%s\n", cfg.RepoOwner, cfg.RepoName)
		fmt.Printf("Range: %s..%s\n", from, to)
		fmt.Printf("Model: %s\n", cfg.OpenAIModel)
		fmt.Println()
	}

	// Create clients
	githubClient := github.NewClient(cfg.GitHubToken, cfg.RepoOwner, cfg.RepoName)
	llmClient := llm.NewOpenAIClient(cfg.OpenAIAPIKey, cfg.OpenAIModel, cfg.MaxTokens, cfg.Temperature)

	// Validate GitHub access
	if cfg.Verbose {
		fmt.Println("Validating GitHub access...")
	}
	if err := githubClient.ValidateAccess(); err != nil {
		return fmt.Errorf("GitHub access validation failed: %w", err)
	}

	// Create generator
	gen := generator.NewGenerator(githubClient, llmClient, cfg)

	// Generate changelog
	changelog, err := gen.Generate(from, to)
	if err != nil {
		return fmt.Errorf("generate changelog: %w", err)
	}

	// Write output
	if cfg.OutputPath == "-" || cfg.OutputPath == "" {
		// Write to stdout
		fmt.Println(changelog.Markdown)
	} else {
		// Write to file
		if err := os.WriteFile(cfg.OutputPath, []byte(changelog.Markdown), 0644); err != nil {
			return fmt.Errorf("write output file: %w", err)
		}
		if cfg.Verbose {
			fmt.Printf("\n✓ Changelog written to %s\n", cfg.OutputPath)
		} else {
			fmt.Printf("Changelog written to %s\n", cfg.OutputPath)
		}
	}

	return nil
}
