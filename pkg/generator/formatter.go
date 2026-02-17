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

// formatTimelineAsMarkdown formats a timeline changelog as PR-based release notes
func (g *Generator) formatTimelineAsMarkdown(timeline *TimelineChangelog) string {
	var b strings.Builder

	// Title and metadata
	b.WriteString(fmt.Sprintf("# Release Notes: %s\n\n", timeline.RepoName))
	b.WriteString(fmt.Sprintf("**Timeline:** %s to %s\n\n",
		timeline.FromDate.Format("January 2, 2006"),
		timeline.ToDate.Format("January 2, 2006")))

	// Each release section
	for i, release := range timeline.Releases {
		b.WriteString(fmt.Sprintf("## [Release %s]\n\n", release.ToRef))

		if len(release.PullRequests) > 0 {
			for _, pr := range release.PullRequests {
				// Format: - PR title by @author in PR_URL
				b.WriteString(fmt.Sprintf("- %s by @%s in %s\n", pr.Title, pr.Author, pr.URL))

				// Add LLM summary indented
				if summary, ok := release.PRSummaries[pr.Number]; ok && summary != "" {
					b.WriteString(fmt.Sprintf("    - %s\n", summary))
				}
			}
		} else {
			b.WriteString("_No pull requests in this release._\n")
		}

		b.WriteString("\n")

		// Separator between releases
		if i < len(timeline.Releases)-1 {
			b.WriteString("---\n\n")
		}
	}

	return b.String()
}
