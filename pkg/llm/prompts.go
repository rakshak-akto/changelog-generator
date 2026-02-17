package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// BuildChangelogPrompt creates the prompt for changelog generation
func BuildChangelogPrompt(req ChangelogRequest) string {
	var sb strings.Builder

	sb.WriteString("You are a technical writer creating a changelog for a software release.\n\n")
	sb.WriteString(fmt.Sprintf("Repository: %s\n", req.RepoName))
	sb.WriteString(fmt.Sprintf("Range: %s → %s\n\n", req.FromRef, req.ToRef))
	sb.WriteString(fmt.Sprintf("Total commits: %d\n\n", len(req.Commits)))

	sb.WriteString("Commits (most recent first):\n")
	sb.WriteString("---\n\n")

	for i, commit := range req.Commits {
		sb.WriteString(fmt.Sprintf("%d. Commit: %s\n", i+1, commit.SHA[:8]))
		sb.WriteString(fmt.Sprintf("   Author: %s\n", commit.Author))
		sb.WriteString(fmt.Sprintf("   Date: %s\n", commit.Date.Format("2006-01-02")))
		sb.WriteString(fmt.Sprintf("   Message: %s\n", commit.Message))

		if len(commit.FilesChanged) > 0 {
			sb.WriteString(fmt.Sprintf("   Files: %s\n", strings.Join(commit.FilesChanged, ", ")))
		}

		if commit.Stats != "" {
			sb.WriteString(fmt.Sprintf("   Stats: %s\n", commit.Stats))
		}

		if commit.DiffSummary != "" {
			sb.WriteString(fmt.Sprintf("   Changes: %s\n", commit.DiffSummary))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("Generate a structured changelog with:\n\n")
	sb.WriteString("1. **Categories**: Organize commits into these categories:\n")
	sb.WriteString("   - Features: New functionality or capabilities\n")
	sb.WriteString("   - Improvements: Enhancements to existing features\n")
	sb.WriteString("   - Bug Fixes: Bug fixes and error corrections\n")
	sb.WriteString("   - Breaking Changes: Changes that break backward compatibility\n")
	sb.WriteString("   - Documentation: Documentation updates\n")
	sb.WriteString("   - Internal: Internal changes, refactoring, or dependencies\n\n")

	sb.WriteString("2. **For each commit**:\n")
	sb.WriteString("   - title: Concise, user-facing title (max 80 chars)\n")
	sb.WriteString("   - description: Brief explanation of the impact (1-2 sentences)\n")
	sb.WriteString("   - importance_score: Rate 0-10 (10=critical/major impact, 5=moderate, 1=minor)\n")
	sb.WriteString("   - Include the SHA and author\n\n")

	sb.WriteString("3. **Top highlights**: Select 3-5 most important changes across all categories\n\n")

	sb.WriteString("4. **Release summary**: Write 2-3 sentences summarizing this release\n\n")

	sb.WriteString("Output ONLY valid JSON with this structure:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"summary\": \"2-3 sentence release summary\",\n")
	sb.WriteString("  \"highlights\": [\"highlight 1\", \"highlight 2\", \"highlight 3\"],\n")
	sb.WriteString("  \"categories\": {\n")
	sb.WriteString("    \"Features\": [\n")
	sb.WriteString("      {\"sha\": \"abc123\", \"title\": \"...\", \"description\": \"...\", \"author\": \"...\", \"importance_score\": 8.5}\n")
	sb.WriteString("    ],\n")
	sb.WriteString("    \"Bug Fixes\": [...],\n")
	sb.WriteString("    ...\n")
	sb.WriteString("  }\n")
	sb.WriteString("}\n\n")
	sb.WriteString("Importance Score Guidelines:\n")
	sb.WriteString("- 9-10: Critical/Breaking changes, major new features, security fixes\n")
	sb.WriteString("- 7-8: Significant features, important bug fixes, notable improvements\n")
	sb.WriteString("- 5-6: Moderate features/fixes, useful enhancements\n")
	sb.WriteString("- 3-4: Minor features/fixes, small improvements\n")
	sb.WriteString("- 1-2: Trivial changes, documentation, internal refactoring\n\n")
	sb.WriteString("Important:\n")
	sb.WriteString("- Only include categories that have commits\n")
	sb.WriteString("- Write from the user's perspective (what changed for them)\n")
	sb.WriteString("- Be concise and clear\n")
	sb.WriteString("- Use the exact category names listed above\n")
	sb.WriteString("- Include importance_score for EVERY commit\n")
	sb.WriteString("- Output ONLY the JSON, no additional text\n")

	return sb.String()
}

// BuildPRChangelogPrompt creates the prompt for PR-based release notes
func BuildPRChangelogPrompt(req PRChangelogRequest) string {
	var sb strings.Builder

	sb.WriteString("You are a technical writer creating release notes for a software release.\n\n")
	sb.WriteString(fmt.Sprintf("Repository: %s\n", req.RepoName))
	sb.WriteString(fmt.Sprintf("Release: %s\n\n", req.ToRef))
	sb.WriteString(fmt.Sprintf("This release contains %d pull requests.\n\n", len(req.PRs)))

	sb.WriteString("Pull Requests:\n")
	sb.WriteString("---\n\n")

	for i, pr := range req.PRs {
		sb.WriteString(fmt.Sprintf("%d. PR #%d: %s\n", i+1, pr.Number, pr.Title))
		sb.WriteString(fmt.Sprintf("   Author: %s\n", pr.Author))
		if pr.Body != "" {
			// Truncate long PR bodies
			body := pr.Body
			if len(body) > 500 {
				body = body[:500] + "..."
			}
			sb.WriteString(fmt.Sprintf("   Description: %s\n", body))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("For each pull request, write a single concise sentence summarizing its user-facing impact.\n")
	sb.WriteString("Focus on WHAT changed from the user's perspective, not implementation details.\n\n")
	sb.WriteString("Output ONLY valid JSON with this structure:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"entries\": [\n")
	sb.WriteString("    {\"number\": 4208, \"summary\": \"One sentence describing what this PR does for users.\"},\n")
	sb.WriteString("    ...\n")
	sb.WriteString("  ]\n")
	sb.WriteString("}\n\n")
	sb.WriteString("Important:\n")
	sb.WriteString("- Include an entry for EVERY pull request\n")
	sb.WriteString("- Each summary must be a single concise sentence\n")
	sb.WriteString("- Write from the user's perspective\n")
	sb.WriteString("- Output ONLY the JSON, no additional text\n")

	return sb.String()
}

// ParsePRChangelogResponse parses the JSON response for PR-based release notes
func ParsePRChangelogResponse(jsonStr string) (*PRChangelogResponse, error) {
	jsonStr = strings.TrimSpace(jsonStr)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	var response PRChangelogResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return nil, fmt.Errorf("parse PR changelog JSON response: %w", err)
	}

	return &response, nil
}

// ParseChangelogResponse parses the JSON response from the LLM
func ParseChangelogResponse(jsonStr string) (*ChangelogResponse, error) {
	// Clean up the response - remove markdown code blocks if present
	jsonStr = strings.TrimSpace(jsonStr)
	jsonStr = strings.TrimPrefix(jsonStr, "```json")
	jsonStr = strings.TrimPrefix(jsonStr, "```")
	jsonStr = strings.TrimSuffix(jsonStr, "```")
	jsonStr = strings.TrimSpace(jsonStr)

	var response ChangelogResponse
	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		return nil, fmt.Errorf("parse JSON response: %w", err)
	}

	return &response, nil
}
