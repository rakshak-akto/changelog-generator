package github

import "time"

// CommitData represents a commit with all its details
type CommitData struct {
	SHA          string
	Message      string
	Author       string
	Date         time.Time
	FilesChanged []FileChange
	Stats        CommitStats
}

// FileChange represents a file modification in a commit
type FileChange struct {
	Filename  string
	Status    string // "added", "modified", "deleted", "renamed"
	Additions int
	Deletions int
	Patch     string // The diff content
}

// CommitStats provides aggregate statistics for a commit
type CommitStats struct {
	Additions int
	Deletions int
	Total     int
}

// TagInfo represents a Git tag with metadata
type TagInfo struct {
	Name       string    // Tag name (e.g., "v1.0.0")
	SHA        string    // Tag object SHA
	CommitSHA  string    // Commit SHA the tag points to
	CommitDate time.Time // Date of the commit
	Message    string    // Tag message (for annotated tags)
}

// ReleaseInfo represents a GitHub release
type ReleaseInfo struct {
	TagName     string    // Associated tag name
	Name        string    // Release name/title
	PublishedAt time.Time // When the release was published
	CreatedAt   time.Time // When the release was created
	Body        string    // Release notes
	Author      string    // Release author
	Draft       bool      // Is draft?
	Prerelease  bool      // Is prerelease?
}

// ReleaseRef represents a unified tag or release reference
type ReleaseRef struct {
	Name         string    // Tag/release name (e.g., "v1.0.0")
	Date         time.Time // Date of tag commit or release publication
	Type         string    // "tag" or "release"
	IsPrerelease bool      // For releases
}

// PullRequestData represents a pull request with its details
type PullRequestData struct {
	Number int
	Title  string
	Author string
	URL    string
	Body   string   // PR description (for LLM context)
	Labels []string
}

// TimelineRelease represents a release period with its commits and PRs
type TimelineRelease struct {
	FromRef      string            // Starting tag/release name
	ToRef        string            // Ending tag/release name
	FromDate     time.Time         // Date of from ref
	ToDate       time.Time         // Date of to ref
	CommitCount  int               // Number of commits
	Commits      []CommitData      // Actual commits
	PullRequests []PullRequestData // PRs in this release
}
