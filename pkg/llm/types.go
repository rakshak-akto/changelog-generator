package llm

import "time"

// ChangelogRequest represents a request to generate a changelog
type ChangelogRequest struct {
	Commits  []CommitInfo
	RepoName string
	FromRef  string
	ToRef    string
}

// CommitInfo contains the information about a commit for LLM processing
type CommitInfo struct {
	SHA          string
	Message      string
	Author       string
	Date         time.Time
	FilesChanged []string
	DiffSummary  string
	Stats        string
}

// ChangelogResponse represents the structured response from the LLM
type ChangelogResponse struct {
	Summary    string                       `json:"summary"`
	Highlights []string                     `json:"highlights"`
	Categories map[string][]ChangelogEntry  `json:"categories"`
}

// ChangelogEntry represents a single entry in the changelog
type ChangelogEntry struct {
	SHA             string  `json:"sha"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Author          string  `json:"author"`
	ImportanceScore float64 `json:"importance_score"` // 0-10 scale, 10 being most important
}

// PRInfo contains pull request information for LLM processing
type PRInfo struct {
	Number int
	Title  string
	Author string
	Body   string
}

// PRChangelogRequest represents a request to generate PR-based release notes
type PRChangelogRequest struct {
	PRs      []PRInfo
	RepoName string
	FromRef  string
	ToRef    string
}

// PRChangelogResponse represents the LLM response for PR-based release notes
type PRChangelogResponse struct {
	Entries []PRSummaryEntry `json:"entries"`
}

// PRSummaryEntry represents a single PR summary from the LLM
type PRSummaryEntry struct {
	Number  int    `json:"number"`
	Summary string `json:"summary"`
}
