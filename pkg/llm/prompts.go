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
	sb.WriteString("   - Include the SHA and author\n\n")

	sb.WriteString("3. **Top highlights**: Select 3-5 most important changes across all categories\n\n")

	sb.WriteString("4. **Release summary**: Write 2-3 sentences summarizing this release\n\n")

	sb.WriteString("Output ONLY valid JSON with this structure:\n")
	sb.WriteString("{\n")
	sb.WriteString("  \"summary\": \"2-3 sentence release summary\",\n")
	sb.WriteString("  \"highlights\": [\"highlight 1\", \"highlight 2\", \"highlight 3\"],\n")
	sb.WriteString("  \"categories\": {\n")
	sb.WriteString("    \"Features\": [{\"sha\": \"abc123\", \"title\": \"...\", \"description\": \"...\", \"author\": \"...\"}],\n")
	sb.WriteString("    \"Bug Fixes\": [...],\n")
	sb.WriteString("    ...\n")
	sb.WriteString("  }\n")
	sb.WriteString("}\n\n")
	sb.WriteString("Important:\n")
	sb.WriteString("- Only include categories that have commits\n")
	sb.WriteString("- Write from the user's perspective (what changed for them)\n")
	sb.WriteString("- Be concise and clear\n")
	sb.WriteString("- Use the exact category names listed above\n")
	sb.WriteString("- Output ONLY the JSON, no additional text\n")

	return sb.String()
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
