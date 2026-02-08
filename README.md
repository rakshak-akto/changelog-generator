# Changelog Generator

A powerful tool that generates human-readable changelogs from Git commits using OpenAI's language models. It automatically categorizes changes, generates summaries, and produces well-formatted markdown output.

## Features

- 🤖 **AI-Powered Analysis**: Uses OpenAI's GPT models to intelligently categorize and summarize commits
- 📊 **Smart Categorization**: Automatically groups changes into Features, Bug Fixes, Improvements, Breaking Changes, Documentation, and Internal
- 📝 **Human-Readable Output**: Generates clean, GitHub-flavored markdown with links and emoji
- ⚡ **Fast & Efficient**: Built in Go for excellent performance
- 🔧 **Highly Configurable**: Supports environment variables, config files, and CLI flags
- 🎯 **GitHub Integration**: Seamlessly fetches commits and diffs from GitHub repositories

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/rakshaksatsangi/changelog-generator.git
cd changelog-generator

# Build the binary
make build

# Or install to GOPATH/bin
make install
```

### Quick Start

```bash
# Build the project
cd /Users/rakshaksatsangi/repos/changelog-generator
make build

# Run the generator
./bin/changelog-generator generate v1.0.0..v1.1.0
```

## Usage

### Prerequisites

You'll need:
1. A GitHub personal access token (create one at https://github.com/settings/tokens)
2. An OpenAI API key (get one at https://platform.openai.com/api-keys)

### Environment Variables

Set these required environment variables:

```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxx
export OPENAI_API_KEY=sk-xxxxxxxxxxxxx
```

### Basic Usage

```bash
# Generate changelog for a commit range
changelog-generator generate v1.0.0..v1.1.0 --owner=facebook --repo=react

# Use HEAD as the end ref
changelog-generator generate v1.0.0..HEAD --owner=myorg --repo=myrepo

# Output to a specific file
changelog-generator generate v1.0.0..v1.1.0 --owner=myorg --repo=myrepo --output=RELEASE_NOTES.md

# Verbose mode for debugging
changelog-generator generate v1.0.0..v1.1.0 --owner=myorg --repo=myrepo --verbose

# Output to stdout
changelog-generator generate v1.0.0..v1.1.0 --owner=myorg --repo=myrepo --output=-
```

### Configuration File

Create a `.changelog.yaml` file in your project root or home directory:

```yaml
# GitHub configuration
repo_owner: facebook
repo_name: react

# OpenAI configuration
openai_model: gpt-4o
max_tokens: 4000
temperature: 0.3

# Output configuration
output_path: CHANGELOG.md
include_authors: true
include_dates: false

# Behavior
verbose: false
```

### Command-Line Flags

```bash
Flags:
      --owner string           Repository owner (required)
      --repo string            Repository name (required)
      --output string          Output file path (default "CHANGELOG.md")
      --model string           OpenAI model to use (default "gpt-4o")
      --verbose                Verbose output
      --include-authors        Include commit authors (default true)
      --include-dates          Include commit dates (default false)
  -h, --help                   Help for generate
```

## Example Output

```markdown
# Changelog: v1.0.0 → v1.1.0

## Summary

This release introduces new authentication features, performance improvements,
and several bug fixes to enhance stability and user experience.

## Highlights

- ⭐ Added OAuth2 authentication support with Google and GitHub providers
- ⭐ 40% performance improvement in API response times
- ⭐ Fixed critical security vulnerability in session handling

## 🚀 Features

- **Add OAuth2 authentication** ([`abc123f`](https://github.com/user/repo/commit/abc123f)) by @johndoe
  Implements OAuth2 flow with support for Google and GitHub providers.
  Includes token refresh and secure storage.

- **Add user profile dashboard** ([`def456a`](https://github.com/user/repo/commit/def456a)) by @janedoe
  New dashboard showing user activity, preferences, and recent changes.

## 🐛 Bug Fixes

- **Fix race condition in cache** ([`ghi789b`](https://github.com/user/repo/commit/ghi789b)) by @bobsmith
  Resolved concurrent access issues causing intermittent failures under high load.

- **Fix memory leak in websocket handler** ([`jkl012c`](https://github.com/user/repo/commit/jkl012c)) by @alicejones
  Fixed goroutine leak that caused memory usage to grow over time.

## ⚡ Improvements

- **Optimize database queries** ([`mno345d`](https://github.com/user/repo/commit/mno345d)) by @charlielee
  Added indexes and query optimization resulting in 40% faster response times.

## 🔧 Internal

- **Update dependencies** ([`pqr678e`](https://github.com/user/repo/commit/pqr678e)) by @davidkim
  Bumped all dependencies to latest versions for security patches.
```

## Architecture

The tool follows a clean, modular architecture:

```
pkg/
├── github/      # GitHub API client for fetching commits
├── llm/         # OpenAI client for AI-powered analysis
├── generator/   # Orchestrates the changelog generation
└── config/      # Configuration management

cmd/
└── cli/         # Command-line interface
```

### How It Works

1. **Fetch Commits**: Uses the GitHub API to fetch all commits in the specified range
2. **Prepare Data**: Extracts commit messages, file changes, and diffs
3. **AI Analysis**: Sends commit data to OpenAI for intelligent categorization and summarization
4. **Format Output**: Generates clean, formatted markdown with links and emoji

## Development

### Project Structure

```
changelog-generator/
├── cmd/
│   └── cli/          # CLI entry point
├── pkg/
│   ├── github/       # GitHub API integration
│   ├── llm/          # OpenAI/LLM integration
│   ├── generator/    # Core changelog generation
│   └── config/       # Configuration management
├── Makefile          # Build automation
├── go.mod            # Go module definition
└── README.md         # This file
```

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linters
make lint

# Clean build artifacts
make clean
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Test specific package
go test -v ./pkg/generator/...
```

## Configuration Priority

Configuration is loaded in the following priority order (highest to lowest):

1. Command-line flags
2. Environment variables
3. Config file (`.changelog.yaml`)
4. Defaults

## Troubleshooting

### Authentication Errors

If you see authentication errors:
- Ensure `GITHUB_TOKEN` is set and valid
- Check that your token has `repo` scope for private repositories
- Verify `OPENAI_API_KEY` is set and has sufficient credits

### Invalid Commit Range

If you get "no commits found" errors:
- Verify the refs exist: `git tag` or `git branch -a`
- Ensure the range is in the correct format: `from..to`
- Try using full commit SHAs instead of tags

### Rate Limiting

- GitHub API has rate limits (5000 requests/hour for authenticated requests)
- OpenAI has rate limits based on your plan
- Use `--verbose` to see detailed progress

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details

## Acknowledgments

- Inspired by [triggerdotdev/ai-changelog](https://github.com/triggerdotdev/ai-changelog)
- Built with [go-github](https://github.com/google/go-github)
- Powered by [OpenAI](https://openai.com/)

## Roadmap

- [ ] Add support for Anthropic Claude
- [ ] HTTP REST API server
- [ ] Commit scoring/filtering for cost optimization
- [ ] Additional output formats (JSON, HTML)
- [ ] Custom changelog templates
- [ ] Caching layer for repeated queries
- [ ] Webhook integration for automatic releases
- [ ] Multi-repository changelog aggregation
