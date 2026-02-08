package config

import (
	"fmt"
	"os"

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
}

// Load loads configuration from environment, config file, and defaults
func Load() (*Config, error) {
	// Set up viper
	viper.SetConfigName(".changelog")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

	// Read config file if it exists (optional)
	_ = viper.ReadInConfig()

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
	if c.RepoOwner == "" {
		return fmt.Errorf("repository owner is required (use --owner flag)")
	}
	if c.RepoName == "" {
		return fmt.Errorf("repository name is required (use --repo flag)")
	}
	if c.OpenAIAPIKey == "" {
		return fmt.Errorf("OpenAI API key is required (set OPENAI_API_KEY environment variable)")
	}
	return nil
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
