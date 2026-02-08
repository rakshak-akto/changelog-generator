# Changelog Generator - Usage Guide

## Quick Start

### 1. Set up authentication

```bash
# Set your GitHub token (required)
export GITHUB_TOKEN=ghp_your_github_token_here

# Set your OpenAI API key (required)
export OPENAI_API_KEY=sk-your_openai_key_here
```

### 2. Build the tool

```bash
cd /Users/rakshaksatsangi/repos/changelog-generator
make build
```

### 3. Generate your first changelog

```bash
# Basic usage - generates CHANGELOG.md
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=facebook \
  --repo=react

# With custom output file
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=myorg \
  --repo=myrepo \
  --output=RELEASE_NOTES.md

# Verbose mode to see progress
./bin/changelog-generator generate v1.0.0..HEAD \
  --owner=myorg \
  --repo=myrepo \
  --verbose
```

## Installation Options

### Option 1: Build from source (recommended for development)

```bash
git clone https://github.com/rakshaksatsangi/changelog-generator.git
cd changelog-generator
make build
# Binary will be at ./bin/changelog-generator
```

### Option 2: Install to GOPATH

```bash
cd changelog-generator
make install
# Binary will be at $GOPATH/bin/changelog-generator
```

### Option 3: Manual build

```bash
go build -o changelog-generator ./cmd/cli
```

## Configuration

The tool supports three levels of configuration (in priority order):

### 1. Command-line flags (highest priority)

```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=myorg \
  --repo=myrepo \
  --model=gpt-4o \
  --output=CHANGELOG.md \
  --verbose \
  --include-authors
```

### 2. Environment variables

```bash
export GITHUB_TOKEN=ghp_xxxxx
export OPENAI_API_KEY=sk-xxxxx
export CHANGELOG_REPO_OWNER=myorg
export CHANGELOG_REPO_NAME=myrepo
export CHANGELOG_OPENAI_MODEL=gpt-4o
export CHANGELOG_OUTPUT_PATH=CHANGELOG.md
export CHANGELOG_VERBOSE=true
```

### 3. Configuration file (lowest priority)

Create `.changelog.yaml` in your project root or home directory:

```yaml
# GitHub configuration
repo_owner: myorg
repo_name: myrepo

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

## Authentication Setup

### GitHub Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Give it a descriptive name (e.g., "Changelog Generator")
4. Select scopes:
   - For public repos: `public_repo`
   - For private repos: `repo` (full control)
5. Click "Generate token"
6. Copy the token and set it as an environment variable:

```bash
export GITHUB_TOKEN=ghp_your_token_here
```

**Security tip**: Add this to your `~/.bashrc` or `~/.zshrc` but NEVER commit it to git!

### OpenAI API Key

1. Go to https://platform.openai.com/api-keys
2. Click "Create new secret key"
3. Give it a name (e.g., "Changelog Generator")
4. Copy the key and set it as an environment variable:

```bash
export OPENAI_API_KEY=sk-your_key_here
```

## Usage Examples

### Example 1: Basic usage with a public repository

```bash
./bin/changelog-generator generate v18.2.0..v18.3.0 \
  --owner=facebook \
  --repo=react
```

**Output**: `CHANGELOG.md` with all changes between React v18.2.0 and v18.3.0

### Example 2: Generate for your own repository

```bash
./bin/changelog-generator generate v1.0.0..HEAD \
  --owner=myusername \
  --repo=my-project \
  --verbose
```

**Output**: Verbose progress updates, then `CHANGELOG.md`

### Example 3: Custom output file and model

```bash
./bin/changelog-generator generate abc123..def456 \
  --owner=myorg \
  --repo=myrepo \
  --output=RELEASE_v2.0.0.md \
  --model=gpt-4o
```

**Output**: `RELEASE_v2.0.0.md` using GPT-4o model

### Example 4: Output to stdout (for piping)

```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=myorg \
  --repo=myrepo \
  --output=- | cat
```

### Example 5: Using commit SHAs instead of tags

```bash
./bin/changelog-generator generate \
  abc123def456..789ghi012jkl \
  --owner=myorg \
  --repo=myrepo
```

### Example 6: With configuration file

Create `.changelog.yaml`:
```yaml
repo_owner: myorg
repo_name: myrepo
openai_model: gpt-4o
verbose: true
include_authors: true
```

Then simply run:
```bash
./bin/changelog-generator generate v1.0.0..v1.1.0
```

## Command Reference

### generate

Generate a changelog for a commit range.

**Usage:**
```bash
changelog-generator generate [from]..[to] [flags]
```

**Arguments:**
- `[from]..[to]`: Git commit range (tags, branches, or SHAs)
  - Examples: `v1.0.0..v1.1.0`, `main..develop`, `abc123..def456`

**Flags:**
- `--owner string`: Repository owner (required)
- `--repo string`: Repository name (required)
- `--output string`: Output file path (default: "CHANGELOG.md")
  - Use `-` for stdout
- `--model string`: OpenAI model (default: "gpt-4o")
  - Options: `gpt-4o`, `gpt-4`, `gpt-4-turbo`, `gpt-3.5-turbo`
- `--verbose`: Enable verbose output
- `--include-authors`: Include commit authors (default: true)
- `--include-dates`: Include commit dates (default: false)
- `-h, --help`: Help for generate command

## Understanding the Output

The generated changelog has this structure:

```markdown
# Changelog: v1.0.0 → v1.1.0

## Summary
A 2-3 sentence overview of the release.

## Highlights
- ⭐ Most important change 1
- ⭐ Most important change 2
- ⭐ Most important change 3

## 💥 Breaking Changes
Critical changes that require user action.

## 🚀 Features
New functionality added.

## ⚡ Improvements
Enhancements to existing features.

## 🐛 Bug Fixes
Bugs that were fixed.

## 📚 Documentation
Documentation updates.

## 🔧 Internal
Internal changes, refactoring, dependencies.
```

Each entry includes:
- **Title**: Human-readable description
- **Commit link**: Clickable SHA linking to GitHub
- **Author**: GitHub username (if enabled)
- **Description**: 1-2 sentence explanation of impact

## Tips & Best Practices

### 1. Start with smaller ranges

Test with a small commit range first to verify everything works:
```bash
./bin/changelog-generator generate HEAD~10..HEAD --owner=myorg --repo=myrepo --verbose
```

### 2. Use verbose mode for debugging

Always use `--verbose` when troubleshooting:
```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 --owner=myorg --repo=myrepo --verbose
```

### 3. Choose the right model

- **gpt-4o** (default): Best balance of quality and cost
- **gpt-4**: Higher quality, more expensive
- **gpt-3.5-turbo**: Faster and cheaper, lower quality

### 4. Review and edit the output

The AI-generated changelog is a great starting point, but always review and edit for:
- Accuracy of categorization
- Completeness of descriptions
- Proper highlighting of breaking changes

### 5. Use tags for releases

Create git tags for releases to make changelog generation easier:
```bash
git tag -a v1.1.0 -m "Release v1.1.0"
git push origin v1.1.0
```

### 6. Automate in CI/CD

Add changelog generation to your release workflow:

```yaml
# .github/workflows/release.yml
- name: Generate Changelog
  run: |
    changelog-generator generate ${{ github.event.release.previous_tag }}..${{ github.event.release.tag_name }} \
      --owner=${{ github.repository_owner }} \
      --repo=${{ github.event.repository.name }} \
      --output=CHANGELOG.md
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
    OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

## Troubleshooting

### Error: "GitHub token is required"

**Solution**: Set the `GITHUB_TOKEN` environment variable:
```bash
export GITHUB_TOKEN=ghp_your_token_here
```

### Error: "OpenAI API key is required"

**Solution**: Set the `OPENAI_API_KEY` environment variable:
```bash
export OPENAI_API_KEY=sk-your_key_here
```

### Error: "repository owner is required"

**Solution**: Add the `--owner` flag:
```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 --owner=myorg --repo=myrepo
```

### Error: "validate repository access"

**Possible causes:**
1. Repository doesn't exist
2. Token doesn't have access (private repo needs `repo` scope)
3. Owner/repo names are incorrect

**Solution**: Verify the repository exists and your token has access:
```bash
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/repos/owner/repo
```

### Error: "no commits found in range"

**Possible causes:**
1. Invalid refs (tags/branches don't exist)
2. Range is backwards (to..from instead of from..to)

**Solution**: Verify the refs exist:
```bash
git tag  # List all tags
git branch -a  # List all branches
```

### Error: "rate limit exceeded"

**Causes:**
- GitHub API rate limit (5000/hour authenticated)
- OpenAI API rate limit

**Solutions:**
- Wait for rate limit to reset
- For OpenAI: Upgrade your API plan
- For GitHub: Use a different token or wait

### Poor changelog quality

**Possible causes:**
1. Too many commits (>100) overwhelming the model
2. Poor commit messages
3. Wrong model selected

**Solutions:**
1. Break into smaller ranges
2. Use a better model (gpt-4o or gpt-4)
3. Improve commit message quality in your repository

## Advanced Usage

### Custom commit range queries

```bash
# Last 10 commits
./bin/changelog-generator generate HEAD~10..HEAD --owner=org --repo=repo

# All commits on a branch since main
./bin/changelog-generator generate main..feature-branch --owner=org --repo=repo

# Between two commits
./bin/changelog-generator generate abc123..def456 --owner=org --repo=repo
```

### Combining with other tools

```bash
# Generate and commit the changelog
./bin/changelog-generator generate v1.0.0..v1.1.0 --owner=org --repo=repo
git add CHANGELOG.md
git commit -m "Update changelog for v1.1.0"

# Generate and open in editor
./bin/changelog-generator generate v1.0.0..v1.1.0 --owner=org --repo=repo && code CHANGELOG.md

# Generate multiple versions
for version in v1.0.0 v1.1.0 v1.2.0; do
  ./bin/changelog-generator generate $prev..$version --owner=org --repo=repo --output=CHANGELOG_$version.md
  prev=$version
done
```

## Cost Estimation

### GitHub API
- **Free tier**: 5,000 requests/hour
- **Cost**: Free
- **Usage**: ~1 request per commit + 1 for range comparison

### OpenAI API
- **gpt-4o**: ~$2.50 per 1M input tokens, ~$10 per 1M output tokens
- **Estimated cost**: $0.01 - $0.10 per changelog (50-100 commits)
- **Factors**: Number of commits, size of diffs, model choice

Example for a 50-commit release:
- Input: ~15K tokens (commit messages + diffs)
- Output: ~2K tokens (changelog)
- Cost: ~$0.05 with gpt-4o

## Next Steps

1. **Test the tool** with a small public repository
2. **Review the output** and verify quality
3. **Customize configuration** for your needs
4. **Integrate into your workflow** (CI/CD, release process)
5. **Provide feedback** and report issues

## Getting Help

- **Documentation**: See [README.md](README.md)
- **Issues**: Report bugs at https://github.com/rakshaksatsangi/changelog-generator/issues
- **Examples**: Check the examples in this guide
