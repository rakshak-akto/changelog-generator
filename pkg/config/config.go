package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the changelog generator
type Config struct {
	// GitHub
	GitHubToken string
	RepoOwner   string
	RepoName    string

	// OpenAI
	OpenAIAPIKey string
	OpenAIModel  string
	MaxTokens    int
	Temperature  float64

	// Output
	OutputPath     string
	IncludeAuthors bool
	IncludeDates   bool
	ShowScores     bool
	MinScore       float64

	// Behavior
	Verbose bool

	// Timeline mode
	TimelineMode bool
	FromDate     time.Time
	ToDate       time.Time
}

// Load loads configuration from environment, config file, and defaults
func Load() (*Config, error) {
	// Look for .changelog.local.yaml first (git-ignored, user-specific)
	// Fall back to .changelog.yaml (committed example/defaults)
	viper.SetConfigName(".changelog.local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Try to read .changelog.local.yaml
	if err := viper.ReadInConfig(); err != nil {
		// If .changelog.local.yaml doesn't exist, try .changelog.yaml
		viper.SetConfigName(".changelog")
		_ = viper.ReadInConfig() // Ignore error, config file is optional
	}

	// Set up environment variable support
	viper.SetEnvPrefix("CHANGELOG")
	viper.AutomaticEnv()

	// Create config with defaults
	cfg := &Config{
		GitHubToken:    getEnvOrViper("GITHUB_TOKEN", ""),
		RepoOwner:      viper.GetString("repo_owner"),
		RepoName:       viper.GetString("repo_name"),
		OpenAIAPIKey:   getEnvOrViper("OPENAI_API_KEY", ""),
		OpenAIModel:    viper.GetString("openai_model"),
		MaxTokens:      viper.GetInt("max_tokens"),
		Temperature:    viper.GetFloat64("temperature"),
		OutputPath:     viper.GetString("output_path"),
		IncludeAuthors: viper.GetBool("include_authors"),
		IncludeDates:   viper.GetBool("include_dates"),
		ShowScores:     viper.GetBool("show_scores"),
		MinScore:       viper.GetFloat64("min_score"),
		Verbose:        viper.GetBool("verbose"),
	}

	// Set defaults if not configured
	if cfg.OpenAIModel == "" {
		cfg.OpenAIModel = "gpt-4o"
	}
	if cfg.MaxTokens == 0 {
		cfg.MaxTokens = 4000
	}
	if cfg.Temperature == 0 {
		cfg.Temperature = 0.3
	}
	if cfg.OutputPath == "" {
		cfg.OutputPath = "CHANGELOG.md"
	}
	if !viper.IsSet("include_authors") {
		cfg.IncludeAuthors = true
	}

	return cfg, nil
}

// Validate checks that all required configuration is present
func (c *Config) Validate() error {
	if c.GitHubToken == "" {
		return fmt.Errorf("GitHub token is required (set GITHUB_TOKEN environment variable)")
	}
	// RepoOwner and RepoName are validated later (after interactive prompt if needed)
	// This allows --interactive flag to work without requiring --owner/--repo upfront
	if c.OpenAIAPIKey == "" {
		return fmt.Errorf("OpenAI API key is required (set OPENAI_API_KEY environment variable)")
	}
	return nil
}

// ValidateTimeline validates timeline-specific configuration
func (c *Config) ValidateTimeline() error {
	if c.FromDate.IsZero() {
		return fmt.Errorf("from-date is required in timeline mode")
	}
	if c.ToDate.IsZero() {
		return fmt.Errorf("to-date is required in timeline mode")
	}
	if c.FromDate.After(c.ToDate) {
		return fmt.Errorf("from-date must be before to-date")
	}
	return nil
}

// ValidateRepository validates that repository information is present
func (c *Config) ValidateRepository() error {
	if c.RepoOwner == "" {
		return fmt.Errorf("repository owner is required")
	}
	if c.RepoName == "" {
		return fmt.Errorf("repository name is required")
	}
	return nil
}

// SaveLocal saves repository configuration to .changelog.local.yaml
func (c *Config) SaveLocal() error {
	viper.Set("repo_owner", c.RepoOwner)
	viper.Set("repo_name", c.RepoName)

	// Write to .changelog.local.yaml in current directory
	return viper.WriteConfigAs(".changelog.local.yaml")
}

// getEnvOrViper gets a value from environment variable first, then viper
func getEnvOrViper(envVar, viperKey string) string {
	if val := os.Getenv(envVar); val != "" {
		return val
	}
	if viperKey != "" {
		return viper.GetString(viperKey)
	}
	return ""
}
