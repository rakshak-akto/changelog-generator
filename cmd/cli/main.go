package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
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
	Use:   "generate [from]..[to] OR use --from-date/--to-date flags",
	Short: "Generate a changelog for a commit range or timeline",
	Long: `Generate a changelog from a range of commits or a date range.

Examples:
  # Ref mode (original behavior)
  changelog-generator generate v1.0.0..v1.1.0
  changelog-generator generate --owner=facebook --repo=react v18.2.0..v18.3.0
  changelog-generator generate abc123..HEAD
  changelog-generator generate --output=RELEASE.md v1.0.0..HEAD
  changelog-generator generate --show-scores v1.0.0..v1.1.0
  changelog-generator generate --min-score=7.0 v1.0.0..v1.1.0

  # Timeline mode (new)
  changelog-generator generate --from-date=2024-01-01 --to-date=2024-12-31 --owner=facebook --repo=react
  changelog-generator generate --from-date=2024-01-01 --to-date=2024-12-31 --interactive`,
	Args: cobra.MaximumNArgs(1), // Allow 0 args for timeline mode, 1 for ref mode
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

	// Timeline mode flags
	generateCmd.Flags().String("from-date", "", "Start date for timeline mode (YYYY-MM-DD)")
	generateCmd.Flags().String("to-date", "", "End date for timeline mode (YYYY-MM-DD)")
	generateCmd.Flags().Bool("interactive", false, "Interactively select repository")
}

// promptForRepository prompts user to select a repository interactively
// and optionally saves it to .changelog.local.yaml
func promptForRepository(cfg *config.Config) (owner, repo string, err error) {
	fmt.Println("\n🔍 Repository Selection")
	fmt.Println()

	// Prompt for owner
	ownerPrompt := &survey.Input{
		Message: "Repository owner (e.g., facebook, vercel, golang):",
	}
	if err := survey.AskOne(ownerPrompt, &owner, survey.WithValidator(survey.Required)); err != nil {
		return "", "", err
	}

	// Prompt for repo name
	repoPrompt := &survey.Input{
		Message: "Repository name (e.g., react, next.js, go):",
	}
	if err := survey.AskOne(repoPrompt, &repo, survey.WithValidator(survey.Required)); err != nil {
		return "", "", err
	}

	// Ask if user wants to save for future use
	saveConfig := false
	savePrompt := &survey.Confirm{
		Message: fmt.Sprintf("Save %s/%s to .changelog.local.yaml for future use?", owner, repo),
		Default: true,
	}
	if err := survey.AskOne(savePrompt, &saveConfig); err != nil {
		return "", "", err
	}

	if saveConfig {
		// Save to local config file
		cfg.RepoOwner = owner
		cfg.RepoName = repo
		if err := cfg.SaveLocal(); err != nil {
			fmt.Printf("⚠️  Warning: Could not save to .changelog.local.yaml: %v\n", err)
		} else {
			fmt.Printf("✓ Saved to .changelog.local.yaml (git-ignored)\n")
		}
	}

	fmt.Println()
	return owner, repo, nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// 1. Check for interactive mode first
	interactive, _ := cmd.Flags().GetBool("interactive")
	if interactive {
		owner, repo, err := promptForRepository(cfg)
		if err != nil {
			return fmt.Errorf("repository selection: %w", err)
		}
		cfg.RepoOwner = owner
		cfg.RepoName = repo
	}

	// 2. Detect mode: timeline vs ref-based
	fromDateStr, _ := cmd.Flags().GetString("from-date")
	toDateStr, _ := cmd.Flags().GetString("to-date")
	hasDateFlags := fromDateStr != "" || toDateStr != ""
	hasRefArg := len(args) == 1

	// Validate mode selection
	if hasDateFlags && hasRefArg {
		return fmt.Errorf("cannot use both date flags (--from-date/--to-date) and ref argument ([from]..[to])")
	}
	if !hasDateFlags && !hasRefArg {
		return fmt.Errorf("must specify either date range (--from-date/--to-date) or ref range ([from]..[to])")
	}

	// 3. Route to appropriate mode
	if hasDateFlags {
		return runTimelineMode(cmd, fromDateStr, toDateStr)
	}
	return runRefMode(cmd, args[0])
}

// runRefMode handles the original ref-based generation (v1.0.0..v1.1.0)
func runRefMode(cmd *cobra.Command, commitRange string) error {
	// Parse commit range
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
	if err := cfg.ValidateRepository(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("Changelog Generator v%s (Ref Mode)\n", version)
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
	return writeOutput(changelog.Markdown, "")
}

// runTimelineMode handles timeline-based generation (date range)
func runTimelineMode(cmd *cobra.Command, fromDateStr, toDateStr string) error {
	// Parse dates
	if fromDateStr == "" || toDateStr == "" {
		return fmt.Errorf("both --from-date and --to-date are required for timeline mode")
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		return fmt.Errorf("invalid --from-date format (expected YYYY-MM-DD): %w", err)
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		return fmt.Errorf("invalid --to-date format (expected YYYY-MM-DD): %w", err)
	}

	// Set timeline mode in config
	cfg.TimelineMode = true
	cfg.FromDate = fromDate
	cfg.ToDate = toDate

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	if err := cfg.ValidateTimeline(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}
	if err := cfg.ValidateRepository(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.Verbose {
		fmt.Printf("Changelog Generator v%s (Timeline Mode)\n", version)
		fmt.Printf("Repository: %s/%s\n", cfg.RepoOwner, cfg.RepoName)
		fmt.Printf("Timeline: %s to %s\n", fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
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

	// Generate timeline changelog
	if cfg.Verbose {
		fmt.Printf("Discovering releases from %s to %s...\n",
			fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	}

	changelog, err := gen.GenerateTimeline(fromDate, toDate)
	if err != nil {
		return fmt.Errorf("generate timeline changelog: %w", err)
	}

	// Generate timestamped filename for timeline mode
	// Format: {repo-name}-{day}-{day}-{month}-{year}-changelog.md
	// Example: akto-5-9-feb-2026-changelog.md
	if cfg.OutputPath == "CHANGELOG.md" || cfg.OutputPath == "" {
		fromDay := fromDate.Day()
		toDay := toDate.Day()
		month := strings.ToLower(fromDate.Format("Jan"))
		year := fromDate.Year()

		// Use repo name (just the repo part, not owner)
		repoName := cfg.RepoName

		cfg.OutputPath = fmt.Sprintf("%s-%d-%d-%s-%d-changelog.md",
			repoName, fromDay, toDay, month, year)
	}

	// Write output
	releaseCount := fmt.Sprintf(" (%d releases)", len(changelog.Releases))
	return writeOutput(changelog.Markdown, releaseCount)
}

// writeOutput writes the changelog to file or stdout
func writeOutput(markdown, suffix string) error {
	if cfg.OutputPath == "-" || cfg.OutputPath == "" {
		fmt.Println(markdown)
	} else {
		if err := os.WriteFile(cfg.OutputPath, []byte(markdown), 0644); err != nil {
			return fmt.Errorf("write output file: %w", err)
		}
		if cfg.Verbose {
			fmt.Printf("\n✓ Changelog written to %s%s\n", cfg.OutputPath, suffix)
		} else {
			fmt.Printf("Changelog written to %s%s\n", cfg.OutputPath, suffix)
		}
	}
	return nil
}
