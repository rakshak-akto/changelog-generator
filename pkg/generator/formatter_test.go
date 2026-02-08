package generator

import (
	"strings"
	"testing"

	"github.com/rakshaksatsangi/changelog-generator/pkg/config"
	"github.com/rakshaksatsangi/changelog-generator/pkg/llm"
)

func TestFormatMarkdown(t *testing.T) {
	response := &llm.ChangelogResponse{
		Summary: "This is a test release with new features and bug fixes.",
		Highlights: []string{
			"Added OAuth2 authentication",
			"Fixed critical bug in cache",
		},
		Categories: map[string][]llm.ChangelogEntry{
			"Features": {
				{
					SHA:             "abc123def456",
					Title:           "Add OAuth2 authentication",
					Description:     "Implements OAuth2 flow with Google and GitHub providers.",
					Author:          "johndoe",
					ImportanceScore: 8.5,
				},
			},
			"Bug Fixes": {
				{
					SHA:             "def456ghi789",
					Title:           "Fix race condition in cache",
					Description:     "Resolved concurrent access issues.",
					Author:          "janedoe",
					ImportanceScore: 7.0,
				},
			},
		},
	}

	cfg := &config.Config{
		RepoOwner:      "testorg",
		RepoName:       "testrepo",
		IncludeAuthors: true,
	}

	markdown := FormatMarkdown(response, "v1.0.0", "v1.1.0", cfg)

	// Verify markdown structure
	if markdown == "" {
		t.Error("Expected non-empty markdown")
	}

	// Check for required elements
	requiredStrings := []string{
		"# Changelog: v1.0.0 → v1.1.0",
		"## Summary",
		"This is a test release",
		"## Highlights",
		"Added OAuth2 authentication",
		"Fixed critical bug",
		"## 🚀 Features",
		"Add OAuth2 authentication",
		"abc123d", // First 7 chars of SHA
		"@johndoe",
		"## 🐛 Bug Fixes",
		"Fix race condition",
		"@janedoe",
		"https://github.com/testorg/testrepo/commit/abc123def456",
	}

	for _, str := range requiredStrings {
		if !strings.Contains(markdown, str) {
			t.Errorf("Expected markdown to contain %q\nGot:\n%s", str, markdown)
		}
	}
}

func TestFormatMarkdownWithoutAuthors(t *testing.T) {
	response := &llm.ChangelogResponse{
		Summary:    "Test release",
		Highlights: []string{"Feature 1"},
		Categories: map[string][]llm.ChangelogEntry{
			"Features": {
				{
					SHA:             "abc123",
					Title:           "Test feature",
					Description:     "Test description",
					Author:          "john",
					ImportanceScore: 6.0,
				},
			},
		},
	}

	cfg := &config.Config{
		RepoOwner:      "org",
		RepoName:       "repo",
		IncludeAuthors: false, // Disabled
	}

	markdown := FormatMarkdown(response, "v1.0.0", "v1.1.0", cfg)

	// Should not contain author
	if strings.Contains(markdown, "@john") {
		t.Error("Expected markdown to not contain author when IncludeAuthors is false")
	}

	// But should contain other elements
	if !strings.Contains(markdown, "Test feature") {
		t.Error("Expected markdown to contain feature title")
	}
}

func TestCategoryEmojis(t *testing.T) {
	expectedEmojis := map[string]string{
		"Features":         "🚀",
		"Improvements":     "⚡",
		"Bug Fixes":        "🐛",
		"Breaking Changes": "💥",
		"Documentation":    "📚",
		"Internal":         "🔧",
	}

	for category, expectedEmoji := range expectedEmojis {
		if CategoryEmojis[category] != expectedEmoji {
			t.Errorf("Expected emoji for %s to be %s, got %s",
				category, expectedEmoji, CategoryEmojis[category])
		}
	}
}

func TestCategoryOrder(t *testing.T) {
	expectedOrder := []string{
		"Breaking Changes",
		"Features",
		"Improvements",
		"Bug Fixes",
		"Documentation",
		"Internal",
	}

	if len(CategoryOrder) != len(expectedOrder) {
		t.Errorf("Expected %d categories in order, got %d",
			len(expectedOrder), len(CategoryOrder))
	}

	for i, expected := range expectedOrder {
		if i >= len(CategoryOrder) || CategoryOrder[i] != expected {
			t.Errorf("Expected category at position %d to be %s, got %s",
				i, expected, CategoryOrder[i])
		}
	}
}
