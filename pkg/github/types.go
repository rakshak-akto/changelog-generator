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
