package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	owner  string
	repo   string
	ctx    context.Context
}

// NewClient creates a new GitHub client
func NewClient(token, owner, repo string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &Client{
		client: client,
		owner:  owner,
		repo:   repo,
		ctx:    ctx,
	}
}

// GetCommitRange fetches all commits between two refs
func (c *Client) GetCommitRange(from, to string) ([]CommitData, error) {
	// Use GitHub's compare API to get commits between refs
	comparison, _, err := c.client.Repositories.CompareCommits(
		c.ctx,
		c.owner,
		c.repo,
		from,
		to,
		&github.ListOptions{PerPage: 250},
	)
	if err != nil {
		return nil, fmt.Errorf("compare commits: %w", err)
	}

	var commits []CommitData
	for _, commit := range comparison.Commits {
		// Get full commit details including diffs
		fullCommit, err := c.GetCommitDetails(commit.GetSHA())
		if err != nil {
			return nil, fmt.Errorf("get commit details for %s: %w", commit.GetSHA(), err)
		}
		commits = append(commits, *fullCommit)
	}

	return commits, nil
}

// GetCommitDetails fetches full details for a single commit
func (c *Client) GetCommitDetails(sha string) (*CommitData, error) {
	commit, _, err := c.client.Repositories.GetCommit(
		c.ctx,
		c.owner,
		c.repo,
		sha,
		&github.ListOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("get commit: %w", err)
	}

	// Extract commit data
	commitData := &CommitData{
		SHA:     commit.GetSHA(),
		Message: commit.GetCommit().GetMessage(),
		Date:    commit.GetCommit().GetAuthor().GetDate().Time,
		Stats: CommitStats{
			Additions: commit.GetStats().GetAdditions(),
			Deletions: commit.GetStats().GetDeletions(),
			Total:     commit.GetStats().GetTotal(),
		},
	}

	// Get author info
	if commit.GetAuthor() != nil {
		commitData.Author = commit.GetAuthor().GetLogin()
	} else if commit.GetCommit().GetAuthor() != nil {
		commitData.Author = commit.GetCommit().GetAuthor().GetName()
	}

	// Extract file changes
	for _, file := range commit.Files {
		fileChange := FileChange{
			Filename:  file.GetFilename(),
			Status:    file.GetStatus(),
			Additions: file.GetAdditions(),
			Deletions: file.GetDeletions(),
			Patch:     file.GetPatch(),
		}
		commitData.FilesChanged = append(commitData.FilesChanged, fileChange)
	}

	return commitData, nil
}

// ValidateAccess checks if the client has access to the repository
func (c *Client) ValidateAccess() error {
	_, _, err := c.client.Repositories.Get(c.ctx, c.owner, c.repo)
	if err != nil {
		return fmt.Errorf("validate repository access: %w", err)
	}
	return nil
}
