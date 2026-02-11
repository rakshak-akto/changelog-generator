package generator

import (
	"fmt"
	"strings"
	"time"

	"github.com/rakshaksatsangi/changelog-generator/pkg/config"
	"github.com/rakshaksatsangi/changelog-generator/pkg/github"
	"github.com/rakshaksatsangi/changelog-generator/pkg/llm"
)

// Generator orchestrates the changelog generation workflow
type Generator struct {
	githubClient *github.Client
	llmClient    *llm.OpenAIClient
	config       *config.Config
}

// NewGenerator creates a new changelog generator
func NewGenerator(githubClient *github.Client, llmClient *llm.OpenAIClient, cfg *config.Config) *Generator {
	return &Generator{
		githubClient: githubClient,
		llmClient:    llmClient,
		config:       cfg,
	}
}

// Generate creates a changelog for the specified commit range
func (g *Generator) Generate(from, to string) (*Changelog, error) {
	if g.config.Verbose {
		fmt.Printf("Fetching commits from %s to %s...\n", from, to)
	}

	// 1. Fetch commits from GitHub
	commits, err := g.githubClient.GetCommitRange(from, to)
	if err != nil {
		return nil, fmt.Errorf("fetch commits: %w", err)
	}

	if len(commits) == 0 {
		return nil, fmt.Errorf("no commits found in range %s..%s", from, to)
	}

	if g.config.Verbose {
		fmt.Printf("Found %d commits\n", len(commits))
		fmt.Println("Preparing commits for LLM analysis...")
	}

	// 2. Prepare commits for LLM (with diffs summarized to fit token limits)
	commitInfos := g.prepareCommitsForLLM(commits)

	if g.config.Verbose {
		fmt.Println("Sending to OpenAI for changelog generation...")
	}

	// 3. Send to OpenAI for changelog generation
	response, err := g.llmClient.GenerateChangelog(llm.ChangelogRequest{
		Commits:  commitInfos,
		RepoName: fmt.Sprintf("%s/%s", g.config.RepoOwner, g.config.RepoName),
		FromRef:  from,
		ToRef:    to,
	})
	if err != nil {
		return nil, fmt.Errorf("generate changelog: %w", err)
	}

	if g.config.Verbose {
		fmt.Println("Formatting changelog as markdown...")
	}

	// 4. Format as markdown
	markdown := g.formatAsMarkdown(response, from, to)

	return &Changelog{
		Summary:    response.Summary,
		Highlights: response.Highlights,
		Categories: response.Categories,
		Markdown:   markdown,
		FromRef:    from,
		ToRef:      to,
		RepoName:   fmt.Sprintf("%s/%s", g.config.RepoOwner, g.config.RepoName),
	}, nil
}

// prepareCommitsForLLM converts GitHub commits to LLM-friendly format
func (g *Generator) prepareCommitsForLLM(commits []github.CommitData) []llm.CommitInfo {
	commitInfos := make([]llm.CommitInfo, 0, len(commits))

	for _, commit := range commits {
		// Extract file names
		fileNames := make([]string, 0, len(commit.FilesChanged))
		for _, file := range commit.FilesChanged {
			fileNames = append(fileNames, file.Filename)
		}

		// Limit files shown to first 20 to avoid token overflow
		if len(fileNames) > 20 {
			fileNames = append(fileNames[:20], fmt.Sprintf("... and %d more files", len(fileNames)-20))
		}

		// Create a summary of the diffs
		diffSummary := ""
		if len(commit.FilesChanged) > 0 {
			// For token efficiency, only include diff summary for files with significant changes
			significantChanges := []string{}
			for _, file := range commit.FilesChanged {
				if file.Additions+file.Deletions > 10 { // Only show files with >10 line changes
					if file.Patch != "" {
						summary := llm.SummarizeDiff(file.Patch)
						if summary != "" {
							significantChanges = append(significantChanges, fmt.Sprintf("%s: %s", file.Filename, summary))
						}
					}
				}
			}
			if len(significantChanges) > 0 {
				// Limit to 3 most significant files
				if len(significantChanges) > 3 {
					significantChanges = significantChanges[:3]
				}
				diffSummary = strings.Join(significantChanges, "\n")
			}
		}

		commitInfo := llm.CommitInfo{
			SHA:          commit.SHA,
			Message:      commit.Message,
			Author:       commit.Author,
			Date:         commit.Date,
			FilesChanged: fileNames,
			DiffSummary:  diffSummary,
			Stats:        fmt.Sprintf("+%d/-%d", commit.Stats.Additions, commit.Stats.Deletions),
		}

		commitInfos = append(commitInfos, commitInfo)
	}

	return commitInfos
}

// formatAsMarkdown formats the LLM response as markdown
func (g *Generator) formatAsMarkdown(response *llm.ChangelogResponse, from, to string) string {
	return FormatMarkdown(response, from, to, g.config)
}

// GenerateTimeline generates a changelog for multiple releases in a date range
func (g *Generator) GenerateTimeline(from, to time.Time) (*TimelineChangelog, error) {
	// 1. Discover releases within timeline
	timelineReleases, err := g.githubClient.GetTimelineReleases(from, to)
	if err != nil {
		return nil, fmt.Errorf("discover releases: %w", err)
	}

	if g.config.Verbose {
		fmt.Printf("Found %d releases in timeline\n\n", len(timelineReleases))
	}

	// 2. Process each release
	var releaseChangelogs []ReleaseChangelog
	for i, release := range timelineReleases {
		if g.config.Verbose {
			fmt.Printf("[%d/%d] Processing %s → %s (%d commits)...\n",
				i+1, len(timelineReleases), release.FromRef, release.ToRef, release.CommitCount)
		}

		// Prepare commits for LLM
		commitInfos := g.prepareCommitsForLLM(release.Commits)

		// Generate changelog for this release
		response, err := g.llmClient.GenerateChangelog(llm.ChangelogRequest{
			Commits:  commitInfos,
			RepoName: fmt.Sprintf("%s/%s", g.config.RepoOwner, g.config.RepoName),
			FromRef:  release.FromRef,
			ToRef:    release.ToRef,
		})
		if err != nil {
			return nil, fmt.Errorf("generate changelog for %s: %w", release.ToRef, err)
		}

		releaseChangelogs = append(releaseChangelogs, ReleaseChangelog{
			FromRef:    release.FromRef,
			ToRef:      release.ToRef,
			FromDate:   release.FromDate,
			ToDate:     release.ToDate,
			Summary:    response.Summary,
			Highlights: response.Highlights,
			Categories: response.Categories,
			Commits:    release.Commits,
		})
	}

	if g.config.Verbose {
		fmt.Println()
	}

	// 3. Build timeline changelog
	timeline := &TimelineChangelog{
		FromDate: from,
		ToDate:   to,
		RepoName: fmt.Sprintf("%s/%s", g.config.RepoOwner, g.config.RepoName),
		Releases: releaseChangelogs,
	}

	// 4. Format as markdown
	timeline.Markdown = g.formatTimelineAsMarkdown(timeline)

	return timeline, nil
}
