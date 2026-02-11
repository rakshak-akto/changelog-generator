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
			// Skip entries below minimum score threshold
			if cfg.MinScore > 0 && entry.ImportanceScore < cfg.MinScore {
				continue
			}

			// Format: **Title** ([SHA](link))
			commitLink := fmt.Sprintf("https://github.com/%s/%s/commit/%s",
				cfg.RepoOwner, cfg.RepoName, entry.SHA)

			// Get short SHA (first 7 chars or full if shorter)
			shortSHA := entry.SHA
			if len(shortSHA) > 7 {
				shortSHA = shortSHA[:7]
			}

			sb.WriteString(fmt.Sprintf("- **%s** ([`%s`](%s))",
				entry.Title,
				shortSHA,
				commitLink,
			))

			// Add score if configured
			if cfg.ShowScores {
				scoreIndicator := getScoreIndicator(entry.ImportanceScore)
				sb.WriteString(fmt.Sprintf(" %s **[%.1f]**", scoreIndicator, entry.ImportanceScore))
			}

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
			// Skip entries below minimum score threshold
			if cfg.MinScore > 0 && entry.ImportanceScore < cfg.MinScore {
				continue
			}

			commitLink := fmt.Sprintf("https://github.com/%s/%s/commit/%s",
				cfg.RepoOwner, cfg.RepoName, entry.SHA)

			// Get short SHA (first 7 chars or full if shorter)
			shortSHA := entry.SHA
			if len(shortSHA) > 7 {
				shortSHA = shortSHA[:7]
			}

			sb.WriteString(fmt.Sprintf("- **%s** ([`%s`](%s))",
				entry.Title,
				shortSHA,
				commitLink,
			))

			// Add score if configured
			if cfg.ShowScores {
				scoreIndicator := getScoreIndicator(entry.ImportanceScore)
				sb.WriteString(fmt.Sprintf(" %s **[%.1f]**", scoreIndicator, entry.ImportanceScore))
			}

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

// getScoreIndicator returns a visual indicator based on the importance score
func getScoreIndicator(score float64) string {
	switch {
	case score >= 9.0:
		return "🔴" // Critical
	case score >= 7.0:
		return "🟠" // High
	case score >= 5.0:
		return "🟡" // Medium
	case score >= 3.0:
		return "🟢" // Low
	default:
		return "⚪" // Trivial
	}
}

// formatTimelineAsMarkdown formats a timeline changelog as markdown
func (g *Generator) formatTimelineAsMarkdown(timeline *TimelineChangelog) string {
	var b strings.Builder

	// Title and metadata
	b.WriteString(fmt.Sprintf("# Changelog: %s\n\n", timeline.RepoName))
	b.WriteString(fmt.Sprintf("**Timeline:** %s to %s\n\n",
		timeline.FromDate.Format("January 2, 2006"),
		timeline.ToDate.Format("January 2, 2006")))
	b.WriteString(fmt.Sprintf("**Total Releases:** %d\n\n", len(timeline.Releases)))

	// Table of contents
	b.WriteString("## 📋 Releases\n\n")
	for _, release := range timeline.Releases {
		anchor := sanitizeAnchor(release.ToRef)
		b.WriteString(fmt.Sprintf("- [%s](#%s) — %s\n",
			release.ToRef,
			anchor,
			release.ToDate.Format("Jan 2, 2006")))
	}
	b.WriteString("\n---\n\n")

	// Each release section
	for i, release := range timeline.Releases {
		b.WriteString(fmt.Sprintf("## %s\n\n", release.ToRef))
		b.WriteString(fmt.Sprintf("**Released:** %s\n\n", release.ToDate.Format("January 2, 2006")))
		b.WriteString(fmt.Sprintf("**Range:** `%s` → `%s`\n\n", release.FromRef, release.ToRef))

		// Summary
		if release.Summary != "" {
			b.WriteString("### Summary\n\n")
			b.WriteString(release.Summary)
			b.WriteString("\n\n")
		}

		// Highlights
		if len(release.Highlights) > 0 {
			b.WriteString("### ✨ Highlights\n\n")
			for _, highlight := range release.Highlights {
				b.WriteString(fmt.Sprintf("- %s\n", highlight))
			}
			b.WriteString("\n")
		}

		// Individual commits
		if len(release.Commits) > 0 {
			b.WriteString("### 📝 Commits\n\n")
			for _, commit := range release.Commits {
				// Get short SHA (first 7 chars)
				shortSHA := commit.SHA
				if len(shortSHA) > 7 {
					shortSHA = shortSHA[:7]
				}

				// Format commit link
				commitLink := fmt.Sprintf("https://github.com/%s/%s/commit/%s",
					g.config.RepoOwner, g.config.RepoName, commit.SHA)

				// Extract first line of commit message
				message := commit.Message
				if idx := strings.Index(message, "\n"); idx != -1 {
					message = message[:idx]
				}

				// Format: - [SHA] Message by @author
				b.WriteString(fmt.Sprintf("- [`%s`](%s) %s", shortSHA, commitLink, message))

				if commit.Author != "" {
					b.WriteString(fmt.Sprintf(" by @%s", commit.Author))
				}

				b.WriteString("\n")
			}
			b.WriteString("\n")
		}

		// Categories (reuse existing formatting logic)
		b.WriteString(formatCategoriesForRelease(release.Categories, g.config))

		// Separator between releases
		if i < len(timeline.Releases)-1 {
			b.WriteString("\n---\n\n")
		}
	}

	return b.String()
}

// sanitizeAnchor converts a string to a valid markdown anchor
func sanitizeAnchor(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// formatCategoriesForRelease formats categories with minimal headers (for timeline mode)
func formatCategoriesForRelease(categories map[string][]llm.ChangelogEntry, cfg *config.Config) string {
	var b strings.Builder

	// Category order
	categoryOrder := []string{
		"Breaking Changes",
		"New Features",
		"Enhancements",
		"Bug Fixes",
		"Performance",
		"Documentation",
		"Dependencies",
		"Refactoring",
		"Testing",
		"CI/CD",
		"Other",
	}

	for _, categoryName := range categoryOrder {
		entries, exists := categories[categoryName]
		if !exists || len(entries) == 0 {
			continue
		}

		b.WriteString(fmt.Sprintf("### %s\n\n", categoryName))

		for _, entry := range entries {
			// Apply min-score filter if configured
			if cfg.MinScore > 0 && entry.ImportanceScore < cfg.MinScore {
				continue
			}

			// Format entry
			b.WriteString(fmt.Sprintf("- %s", entry.Description))

			// Add commit link
			if entry.SHA != "" {
				shortSHA := entry.SHA
				if len(shortSHA) > 7 {
					shortSHA = shortSHA[:7]
				}
				b.WriteString(fmt.Sprintf(" ([`%s`](https://github.com/%s/%s/commit/%s))",
					shortSHA, cfg.RepoOwner, cfg.RepoName, entry.SHA))
			}

			// Add importance score if enabled
			if cfg.ShowScores {
				b.WriteString(fmt.Sprintf(" `[%.1f]`", entry.ImportanceScore))
			}

			b.WriteString("\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}
