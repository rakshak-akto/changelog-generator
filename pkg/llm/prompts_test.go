package llm

import (
	"testing"
	"time"
)

func TestBuildChangelogPrompt(t *testing.T) {
	req := ChangelogRequest{
		Commits: []CommitInfo{
			{
				SHA:          "abc123def456",
				Message:      "Add new feature",
				Author:       "johndoe",
				Date:         time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				FilesChanged: []string{"main.go", "README.md"},
				DiffSummary:  "Added 50 lines, removed 10 lines",
				Stats:        "+50/-10",
			},
		},
		RepoName: "test/repo",
		FromRef:  "v1.0.0",
		ToRef:    "v1.1.0",
	}

	prompt := BuildChangelogPrompt(req)

	// Verify prompt contains key elements
	if prompt == "" {
		t.Error("Expected non-empty prompt")
	}

	// Check for required sections
	requiredStrings := []string{
		"test/repo",
		"v1.0.0",
		"v1.1.0",
		"abc123de", // First 8 chars of SHA
		"Add new feature",
		"johndoe",
		"Features",
		"Bug Fixes",
	}

	for _, str := range requiredStrings {
		if !contains(prompt, str) {
			t.Errorf("Expected prompt to contain %q", str)
		}
	}
}

func TestParseChangelogResponse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid JSON",
			input: `{
				"summary": "Test release",
				"highlights": ["Feature 1", "Fix 2"],
				"categories": {
					"Features": [
						{
							"sha": "abc123",
							"title": "Add feature",
							"description": "Added a new feature",
							"author": "john"
						}
					]
				}
			}`,
			wantErr: false,
		},
		{
			name: "JSON with markdown code blocks",
			input: "```json\n" + `{
				"summary": "Test release",
				"highlights": [],
				"categories": {}
			}` + "\n```",
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   "not valid json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := ParseChangelogResponse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChangelogResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && resp == nil {
				t.Error("Expected non-nil response for valid input")
			}
		})
	}
}

func TestTruncateDiff(t *testing.T) {
	diff := ""
	for i := 0; i < 100; i++ {
		diff += "line " + string(rune(i)) + "\n"
	}

	truncated := TruncateDiff(diff, 10)

	if !contains(truncated, "truncated") {
		t.Error("Expected truncated message")
	}
}

func TestSummarizeDiff(t *testing.T) {
	diff := `--- a/file.go
+++ b/file.go
@@ -1,3 +1,5 @@
+added line 1
+added line 2
 existing line
-removed line
 existing line 2
`

	summary := SummarizeDiff(diff)

	if summary == "" {
		t.Error("Expected non-empty summary")
	}

	if !contains(summary, "+") || !contains(summary, "-") {
		t.Error("Expected summary to contain addition/deletion counts")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[0:1] == substr[0:1] && contains(s[1:], substr[1:])) || contains(s[1:], substr)))
}
