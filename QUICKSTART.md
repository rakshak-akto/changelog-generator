# Quick Start Guide

Get up and running with Changelog Generator in 5 minutes!

## Prerequisites

- Go 1.23+ installed
- Git installed
- A GitHub account
- An OpenAI account

## Step 1: Get Your API Keys (2 minutes)

### GitHub Token

1. Go to https://github.com/settings/tokens
2. Click "Generate new token" → "Generate new token (classic)"
3. Name it "Changelog Generator"
4. Select scopes:
   - `public_repo` (for public repositories)
   - `repo` (for private repositories)
5. Click "Generate token" and copy it

### OpenAI API Key

1. Go to https://platform.openai.com/api-keys
2. Click "Create new secret key"
3. Name it "Changelog Generator"
4. Copy the key

## Step 2: Set Environment Variables (30 seconds)

```bash
# Add to your ~/.bashrc or ~/.zshrc
export GITHUB_TOKEN="ghp_your_github_token_here"
export OPENAI_API_KEY="sk-your_openai_key_here"

# Apply immediately
source ~/.bashrc  # or source ~/.zshrc
```

## Step 3: Build the Tool (1 minute)

```bash
# Navigate to the project
cd /Users/rakshaksatsangi/repos/changelog-generator

# Run the automated setup
./scripts/setup.sh

# Or build manually
make build
```

## Step 4: Generate Your First Changelog (1 minute)

### Try with a public repository first:

```bash
./bin/changelog-generator generate v18.2.0..v18.3.0 \
  --owner=facebook \
  --repo=react \
  --verbose
```

This will:
1. ✅ Fetch commits from React's GitHub repository
2. ✅ Send them to OpenAI for analysis
3. ✅ Generate a beautiful changelog
4. ✅ Save it to `CHANGELOG.md`

### Or use your own repository:

```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=your-username \
  --repo=your-repo \
  --verbose
```

## Step 5: Review the Output (30 seconds)

Open `CHANGELOG.md` to see your generated changelog:

```bash
cat CHANGELOG.md
# or
code CHANGELOG.md
# or
open CHANGELOG.md
```

You should see something like:

```markdown
# Changelog: v1.0.0 → v1.1.0

## Summary
This release introduces new features and bug fixes...

## Highlights
- ⭐ New OAuth2 authentication
- ⭐ Performance improvements
- ⭐ Critical bug fixes

## 🚀 Features
- **Add OAuth2 support** ([`abc123`](link))
  Implements OAuth2 with Google and GitHub providers.

## 🐛 Bug Fixes
- **Fix race condition** ([`def456`](link))
  Resolved concurrent access issues.
```

## Common Issues & Solutions

### "GitHub token is required"

```bash
# Check if token is set
echo $GITHUB_TOKEN

# If empty, set it
export GITHUB_TOKEN="ghp_your_token_here"
```

### "OpenAI API key is required"

```bash
# Check if key is set
echo $OPENAI_API_KEY

# If empty, set it
export OPENAI_API_KEY="sk-your_key_here"
```

### "Build failed"

```bash
# Check Go version
go version  # Should be 1.23+

# Install dependencies
go mod download

# Try building again
make build
```

### "No commits found"

Make sure:
1. The repository exists
2. The tags/refs are correct
3. Your GitHub token has access

```bash
# Verify refs exist
git ls-remote https://github.com/owner/repo

# Or check locally
git tag
git branch -a
```

## Next Steps

### Customize Your Configuration

Create `.changelog.yaml` in your project:

```yaml
repo_owner: your-username
repo_name: your-repo
openai_model: gpt-4o
verbose: true
include_authors: true
output_path: CHANGELOG.md
```

Then simply run:
```bash
./bin/changelog-generator generate v1.0.0..v1.1.0
```

### Install Globally

```bash
make install
# Now you can use it anywhere:
changelog-generator generate v1.0.0..v1.1.0 --owner=org --repo=repo
```

### Integrate with CI/CD

Add to `.github/workflows/release.yml`:

```yaml
name: Generate Changelog
on:
  release:
    types: [created]

jobs:
  changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.23'

      - name: Build changelog generator
        run: |
          git clone https://github.com/rakshaksatsangi/changelog-generator.git
          cd changelog-generator
          make build

      - name: Generate changelog
        run: |
          ./changelog-generator/bin/changelog-generator generate \
            ${{ github.event.release.previous_tag }}..${{ github.event.release.tag_name }} \
            --owner=${{ github.repository_owner }} \
            --repo=${{ github.event.repository.name }}
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
```

## Pro Tips

### 1. Use Verbose Mode for Debugging

```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=org --repo=repo --verbose
```

### 2. Start with Small Ranges

Test with recent commits first:
```bash
./bin/changelog-generator generate HEAD~5..HEAD \
  --owner=org --repo=repo
```

### 3. Output to Stdout

Pipe to other tools:
```bash
./bin/changelog-generator generate v1.0.0..v1.1.0 \
  --owner=org --repo=repo --output=- | pbcopy
```

### 4. Use Different Models

```bash
# Higher quality (more expensive)
--model=gpt-4

# Default (balanced)
--model=gpt-4o

# Faster/cheaper (lower quality)
--model=gpt-3.5-turbo
```

### 5. Batch Generate Multiple Versions

```bash
for tag in v1.0.0 v1.1.0 v1.2.0; do
  if [ -n "$prev" ]; then
    ./bin/changelog-generator generate $prev..$tag \
      --owner=org --repo=repo \
      --output=CHANGELOG_$tag.md
  fi
  prev=$tag
done
```

## Getting Help

- **Documentation**: See [README.md](README.md) and [USAGE.md](USAGE.md)
- **Contributing**: See [CONTRIBUTING.md](CONTRIBUTING.md)
- **Issues**: Report at https://github.com/rakshaksatsangi/changelog-generator/issues

## Success! 🎉

You're now ready to generate beautiful changelogs with AI!

**What's Next?**
1. ⭐ Star the repository
2. 📖 Read the full [USAGE.md](USAGE.md) for advanced features
3. 🤝 Contribute improvements (see [CONTRIBUTING.md](CONTRIBUTING.md))
4. 📣 Share with your team

---

**Estimated Cost**: $0.01 - $0.10 per changelog
**Time Saved**: Hours of manual changelog writing
**Quality**: Professional, consistent, human-readable
