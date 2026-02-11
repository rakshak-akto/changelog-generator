package github

import (
	"context"
	"fmt"
	"sort"
	"time"

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

// ListAllTags fetches all tags from the repository with pagination
func (c *Client) ListAllTags() ([]TagInfo, error) {
	var allTags []TagInfo
	opts := &github.ListOptions{PerPage: 100}

	for {
		tags, resp, err := c.client.Repositories.ListTags(
			c.ctx,
			c.owner,
			c.repo,
			opts,
		)
		if err != nil {
			return nil, fmt.Errorf("list tags: %w", err)
		}

		for _, tag := range tags {
			// Get commit details to extract date
			commit, _, err := c.client.Repositories.GetCommit(
				c.ctx,
				c.owner,
				c.repo,
				tag.GetCommit().GetSHA(),
				&github.ListOptions{},
			)
			if err != nil {
				return nil, fmt.Errorf("get commit for tag %s: %w", tag.GetName(), err)
			}

			tagInfo := TagInfo{
				Name:       tag.GetName(),
				SHA:        tag.GetCommit().GetSHA(),
				CommitSHA:  tag.GetCommit().GetSHA(),
				CommitDate: commit.GetCommit().GetCommitter().GetDate().Time,
				Message:    "", // Tag message not available from ListTags
			}
			allTags = append(allTags, tagInfo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allTags, nil
}

// ListAllReleases fetches all GitHub releases with pagination
func (c *Client) ListAllReleases() ([]ReleaseInfo, error) {
	var allReleases []ReleaseInfo
	opts := &github.ListOptions{PerPage: 100}

	for {
		releases, resp, err := c.client.Repositories.ListReleases(
			c.ctx,
			c.owner,
			c.repo,
			opts,
		)
		if err != nil {
			return nil, fmt.Errorf("list releases: %w", err)
		}

		for _, release := range releases {
			releaseInfo := ReleaseInfo{
				TagName:     release.GetTagName(),
				Name:        release.GetName(),
				PublishedAt: release.GetPublishedAt().Time,
				CreatedAt:   release.GetCreatedAt().Time,
				Body:        release.GetBody(),
				Author:      release.GetAuthor().GetLogin(),
				Draft:       release.GetDraft(),
				Prerelease:  release.GetPrerelease(),
			}
			allReleases = append(allReleases, releaseInfo)
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return allReleases, nil
}

// GetReleaseRefsInTimeline discovers all tags and releases within a date range
// Returns deduplicated, sorted list of release references
func (c *Client) GetReleaseRefsInTimeline(from, to time.Time) ([]ReleaseRef, error) {
	// Fetch tags
	tags, err := c.ListAllTags()
	if err != nil {
		return nil, fmt.Errorf("fetch tags: %w", err)
	}

	// Fetch releases
	releases, err := c.ListAllReleases()
	if err != nil {
		return nil, fmt.Errorf("fetch releases: %w", err)
	}

	// Convert to unified ReleaseRef format
	refMap := make(map[string]ReleaseRef) // Map to deduplicate by name

	// Add tags
	for _, tag := range tags {
		if (tag.CommitDate.Equal(from) || tag.CommitDate.After(from)) &&
			(tag.CommitDate.Equal(to) || tag.CommitDate.Before(to)) {
			refMap[tag.Name] = ReleaseRef{
				Name:         tag.Name,
				Date:         tag.CommitDate,
				Type:         "tag",
				IsPrerelease: false,
			}
		}
	}

	// Add releases (may override tags with same name, prioritizing release data)
	for _, release := range releases {
		if release.Draft {
			continue // Skip draft releases
		}

		pubDate := release.PublishedAt
		if (pubDate.Equal(from) || pubDate.After(from)) &&
			(pubDate.Equal(to) || pubDate.Before(to)) {
			// If tag exists with same name, update with release info
			refMap[release.TagName] = ReleaseRef{
				Name:         release.TagName,
				Date:         pubDate,
				Type:         "release",
				IsPrerelease: release.Prerelease,
			}
		}
	}

	// Convert map to slice
	var refs []ReleaseRef
	for _, ref := range refMap {
		refs = append(refs, ref)
	}

	// Sort by date ascending
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Date.Before(refs[j].Date)
	})

	return refs, nil
}

// GetTimelineReleases builds TimelineRelease objects for consecutive ref pairs
func (c *Client) GetTimelineReleases(from, to time.Time) ([]TimelineRelease, error) {
	// Get all release refs in timeline
	refs, err := c.GetReleaseRefsInTimeline(from, to)
	if err != nil {
		return nil, err
	}

	if len(refs) == 0 {
		return nil, fmt.Errorf("no tags or releases found between %s and %s",
			from.Format("2006-01-02"), to.Format("2006-01-02"))
	}

	// Build timeline releases from consecutive pairs
	var timelineReleases []TimelineRelease
	for i := 0; i < len(refs)-1; i++ {
		fromRef := refs[i]
		toRef := refs[i+1]

		// Fetch commits between these refs
		commits, err := c.GetCommitRange(fromRef.Name, toRef.Name)
		if err != nil {
			return nil, fmt.Errorf("get commits %s..%s: %w", fromRef.Name, toRef.Name, err)
		}

		timelineReleases = append(timelineReleases, TimelineRelease{
			FromRef:     fromRef.Name,
			ToRef:       toRef.Name,
			FromDate:    fromRef.Date,
			ToDate:      toRef.Date,
			CommitCount: len(commits),
			Commits:     commits,
		})
	}

	return timelineReleases, nil
}
