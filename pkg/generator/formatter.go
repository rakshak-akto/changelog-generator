package generator

import (
	"fmt"
	"strings"

	"github.com/rakshaksatsangi/changelog-generator/pkg/config"
	"github.com/rakshaksatsangi/changelog-generator/pkg/llm"
)

// CategoryEmojis maps category names to emoji prefixes
var CategoryEmojis = map[string]string{
	"Features":         "🚀",
	"Improvements":     "⚡",
	"Bug Fixes":        "🐛",
	"Breaking Changes": "💥",
	"Documentation":    "📚",
	"Internal":         "🔧",
}

// CategoryOrder defines the order in which categories appear
var CategoryOrder = []string{
	"Breaking Changes",
	"Features",
	"Improvements",
	"Bug Fixes",
	"Documentation",
	"Internal",
}

// FormatMarkdown generates GitHub-flavored markdown from the changelog response
func FormatMarkdown(response *llm.ChangelogResponse, from, to string, cfg *config.Config) string {
	var sb strings.Builder

	// Title
	sb.WriteString(fmt.Sprintf("# Changelog: %s → %s\n\n", from, to))

	// Summary
	if response.Summary != "" {
		sb.WriteString("## Summary\n\n")
		sb.WriteString(response.Summary)
		sb.WriteString("\n\n")
	}

	// Highlights
	if len(response.Highlights) > 0 {
		sb.WriteString("## Highlights\n\n")
		for _, highlight := range response.Highlights {
			sb.WriteString(fmt.Sprintf("- ⭐ %s\n", highlight))
		}
		sb.WriteString("\n")
	}

	// Categories in order
	for _, category := range CategoryOrder {
		entries, exists := response.Categories[category]
		if !exists || len(entries) == 0 {
			continue
		}

		emoji := CategoryEmojis[category]
		if emoji == "" {
			emoji = "•"
		}

		sb.WriteString(fmt.Sprintf("## %s %s\n\n", emoji, category))

		for _, entry := range entries {
			// Format: **Title** ([SHA](link))
			commitLink := fmt.Sprintf("https://github.com/%s/%s/commit/%s",
				cfg.RepoOwner, cfg.RepoName, entry.SHA)

			sb.WriteString(fmt.Sprintf("- **%s** ([`%s`](%s))",
				entry.Title,
				entry.SHA[:7],
				commitLink,
			))

			// Add author if configured
			if cfg.IncludeAuthors && entry.Author != "" {
				sb.WriteString(fmt.Sprintf(" by @%s", entry.Author))
			}

			sb.WriteString("\n")

			// Add description if present
			if entry.Description != "" {
				// Indent description
				lines := strings.Split(entry.Description, "\n")
				for _, line := range lines {
					if line != "" {
						sb.WriteString(fmt.Sprintf("  %s\n", line))
					}
				}
			}

			sb.WriteString("\n")
		}
	}

	// Add any categories that weren't in our predefined order
	for category, entries := range response.Categories {
		// Skip if already processed
		alreadyProcessed := false
		for _, knownCategory := range CategoryOrder {
			if category == knownCategory {
				alreadyProcessed = true
				break
			}
		}
		if alreadyProcessed || len(entries) == 0 {
			continue
		}

		// Use default emoji for unknown categories
		sb.WriteString(fmt.Sprintf("## • %s\n\n", category))

		for _, entry := range entries {
			commitLink := fmt.Sprintf("https://github.com/%s/%s/commit/%s",
				cfg.RepoOwner, cfg.RepoName, entry.SHA)

			sb.WriteString(fmt.Sprintf("- **%s** ([`%s`](%s))",
				entry.Title,
				entry.SHA[:7],
				commitLink,
			))

			if cfg.IncludeAuthors && entry.Author != "" {
				sb.WriteString(fmt.Sprintf(" by @%s", entry.Author))
			}

			sb.WriteString("\n")

			if entry.Description != "" {
				lines := strings.Split(entry.Description, "\n")
				for _, line := range lines {
					if line != "" {
						sb.WriteString(fmt.Sprintf("  %s\n", line))
					}
				}
			}

			sb.WriteString("\n")
		}
	}

	return sb.String()
}
