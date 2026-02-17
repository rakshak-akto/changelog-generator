package generator

import (
	"time"

	"github.com/rakshaksatsangi/changelog-generator/pkg/github"
	"github.com/rakshaksatsangi/changelog-generator/pkg/llm"
)

// Changelog represents the complete generated changelog
type Changelog struct {
	Summary    string
	Highlights []string
	Categories map[string][]llm.ChangelogEntry
	Markdown   string
	FromRef    string
	ToRef      string
	RepoName   string
}

// TimelineChangelog represents a changelog covering multiple releases
type TimelineChangelog struct {
	FromDate time.Time
	ToDate   time.Time
	RepoName string
	Releases []ReleaseChangelog
	Markdown string
}

// ReleaseChangelog represents a single release within a timeline
type ReleaseChangelog struct {
	FromRef      string
	ToRef        string
	FromDate     time.Time
	ToDate       time.Time
	Summary      string
	Highlights   []string
	Categories   map[string][]llm.ChangelogEntry
	Commits      []github.CommitData      // Individual commits in this release
	PullRequests []github.PullRequestData  // PRs in this release
	PRSummaries  map[int]string            // PR number → LLM summary
}
