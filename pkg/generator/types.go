package generator

import "github.com/rakshaksatsangi/changelog-generator/pkg/llm"

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
