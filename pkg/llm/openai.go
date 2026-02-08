package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

// OpenAIClient wraps the OpenAI API client
type OpenAIClient struct {
	client      *openai.Client
	model       string
	maxTokens   int
	temperature float64
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey, model string, maxTokens int, temperature float64) *OpenAIClient {
	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	return &OpenAIClient{
		client:      &client,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}
}

// GenerateChangelog generates a changelog using OpenAI
func (c *OpenAIClient) GenerateChangelog(req ChangelogRequest) (*ChangelogResponse, error) {
	// Build the prompt
	prompt := BuildChangelogPrompt(req)

	// Create chat completion request
	ctx := context.Background()
	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		},
		Model:       openai.ChatModel(c.model),
		MaxTokens:   param.NewOpt(int64(c.maxTokens)),
		Temperature: param.NewOpt(c.temperature),
	}

	chatCompletion, err := c.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("create chat completion: %w", err)
	}

	// Extract the response
	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := chatCompletion.Choices[0].Message.Content

	// Parse the JSON response
	response, err := ParseChangelogResponse(content)
	if err != nil {
		return nil, fmt.Errorf("parse changelog response: %w", err)
	}

	return response, nil
}

// TruncateDiff truncates a diff to a reasonable size for token limits
func TruncateDiff(diff string, maxLines int) string {
	lines := strings.Split(diff, "\n")
	if len(lines) <= maxLines {
		return diff
	}

	truncated := strings.Join(lines[:maxLines], "\n")
	return truncated + fmt.Sprintf("\n... (%d more lines truncated)", len(lines)-maxLines)
}

// SummarizeDiff creates a brief summary of changes from a diff
func SummarizeDiff(diff string) string {
	if diff == "" {
		return ""
	}

	lines := strings.Split(diff, "\n")
	additions := 0
	deletions := 0

	for _, line := range lines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			additions++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletions++
		}
	}

	// Get a sample of the changes
	sample := TruncateDiff(diff, 10)

	return fmt.Sprintf("+%d/-%d lines. Sample:\n%s", additions, deletions, sample)
}
