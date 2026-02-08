# Contributing to Changelog Generator

Thank you for your interest in contributing to Changelog Generator! This document provides guidelines and instructions for contributing.

## Getting Started

### Prerequisites

- Go 1.23 or higher
- Git
- GitHub account
- (For testing) GitHub personal access token
- (For testing) OpenAI API key

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/changelog-generator.git
   cd changelog-generator
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build the Project**
   ```bash
   make build
   ```

4. **Run Tests**
   ```bash
   make test
   ```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

### 2. Make Your Changes

- Write clean, readable code
- Follow existing code style
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make vet
```

### 4. Commit Your Changes

We use conventional commits. Format your commit messages as:

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
git commit -m "feat(llm): add support for Anthropic Claude"
git commit -m "fix(formatter): handle short SHA values correctly"
git commit -m "docs(readme): add installation instructions"
```

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Code Style Guidelines

### Go Code Style

1. **Follow Standard Go Conventions**
   - Use `gofmt` for formatting
   - Follow effective Go guidelines
   - Use meaningful variable names

2. **Package Structure**
   ```
   pkg/
   ├── package/
   │   ├── package.go       # Main implementation
   │   ├── package_test.go  # Tests
   │   └── types.go         # Type definitions
   ```

3. **Error Handling**
   ```go
   // Good: Wrap errors with context
   if err != nil {
       return fmt.Errorf("fetch commits: %w", err)
   }

   // Bad: Return raw errors
   if err != nil {
       return err
   }
   ```

4. **Comments**
   ```go
   // Good: Document public functions
   // GenerateChangelog creates a changelog for the specified commit range.
   // It fetches commits from GitHub and uses OpenAI to generate structured output.
   func GenerateChangelog(from, to string) (*Changelog, error) {
       // ...
   }
   ```

5. **Testing**
   ```go
   func TestFeature(t *testing.T) {
       tests := []struct {
           name    string
           input   string
           want    string
           wantErr bool
       }{
           // Test cases...
       }

       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test implementation...
           })
       }
   }
   ```

## Areas for Contribution

### High Priority

1. **Add Anthropic Claude Support**
   - Create `pkg/llm/anthropic.go`
   - Implement provider interface
   - Add configuration options
   - Update documentation

2. **HTTP REST API**
   - Create `cmd/server/` package
   - Implement REST endpoints
   - Add request validation
   - Create OpenAPI spec

3. **Improve Test Coverage**
   - Add integration tests
   - Mock GitHub/OpenAI clients
   - Add edge case tests
   - Test error paths

### Medium Priority

4. **Commit Filtering**
   - Add importance scoring
   - Filter trivial commits
   - Configurable thresholds
   - Cost optimization

5. **Additional Output Formats**
   - JSON output
   - HTML output
   - Custom templates
   - Multiple formats in one run

6. **Caching Layer**
   - Cache commit data
   - Cache LLM responses
   - Configurable TTL
   - Storage options

### Low Priority

7. **Enhanced Configuration**
   - Interactive setup wizard
   - Configuration validation
   - More granular controls
   - Profile management

8. **UI Improvements**
   - Progress bars
   - Better error messages
   - Color output
   - Interactive prompts

9. **Documentation**
   - More examples
   - Video tutorials
   - Blog posts
   - Architecture diagrams

## Testing Guidelines

### Unit Tests

- Test public functions
- Test edge cases
- Test error conditions
- Use table-driven tests

Example:
```go
func TestFormatMarkdown(t *testing.T) {
    tests := []struct {
        name     string
        response *ChangelogResponse
        want     string
    }{
        {
            name: "basic formatting",
            response: &ChangelogResponse{
                Summary: "Test",
                Categories: map[string][]Entry{
                    "Features": {{Title: "Feature 1"}},
                },
            },
            want: "# Changelog",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := FormatMarkdown(tt.response)
            if !strings.Contains(got, tt.want) {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Integration Tests

For features that interact with external services:

```go
// +build integration

func TestGitHubIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    token := os.Getenv("GITHUB_TOKEN")
    if token == "" {
        t.Skip("GITHUB_TOKEN not set")
    }

    // Test implementation...
}
```

Run with: `go test -tags=integration ./...`

## Documentation Guidelines

### Code Documentation

1. **Package Comments**
   ```go
   // Package generator orchestrates changelog generation.
   // It coordinates between GitHub API and LLM services to produce
   // human-readable changelogs from commit histories.
   package generator
   ```

2. **Function Comments**
   ```go
   // Generate creates a changelog for the specified commit range.
   // It fetches commits from GitHub, analyzes them with an LLM,
   // and formats the results as markdown.
   //
   // The from and to parameters can be tags, branches, or commit SHAs.
   // Returns an error if the range is invalid or API calls fail.
   func Generate(from, to string) (*Changelog, error) {
       // ...
   }
   ```

### README Updates

When adding features, update:
- Feature list
- Usage examples
- Configuration options
- Troubleshooting section

### USAGE.md Updates

For new functionality, add:
- Command examples
- Configuration examples
- Common use cases
- Troubleshooting tips

## Pull Request Process

1. **Before Submitting**
   - [ ] Tests pass (`make test`)
   - [ ] Code is formatted (`make fmt`)
   - [ ] No lint errors (`make vet`)
   - [ ] Documentation updated
   - [ ] Commit messages follow conventions

2. **PR Description Should Include**
   - What: Brief description of changes
   - Why: Motivation for changes
   - How: Technical approach
   - Testing: How you tested the changes
   - Screenshots: For UI changes

3. **PR Template**
   ```markdown
   ## Description
   Brief description of what this PR does.

   ## Motivation
   Why is this change needed?

   ## Changes
   - Change 1
   - Change 2

   ## Testing
   How was this tested?

   ## Screenshots
   (if applicable)

   ## Checklist
   - [ ] Tests added/updated
   - [ ] Documentation updated
   - [ ] Code formatted
   - [ ] No lint errors
   ```

4. **Review Process**
   - Maintainers will review within 1-2 days
   - Address review comments
   - Keep PR focused and small
   - Be responsive to feedback

## Issue Guidelines

### Reporting Bugs

Use this template:

```markdown
## Description
Clear description of the bug.

## Steps to Reproduce
1. Step 1
2. Step 2
3. Step 3

## Expected Behavior
What should happen?

## Actual Behavior
What actually happens?

## Environment
- OS: macOS/Linux/Windows
- Go Version: 1.23
- Project Version: v0.1.0

## Additional Context
Any other relevant information.
```

### Requesting Features

Use this template:

```markdown
## Feature Description
Clear description of the feature.

## Use Case
Why is this feature needed?

## Proposed Solution
How should this work?

## Alternatives Considered
Other approaches you've thought about.

## Additional Context
Any other relevant information.
```

## Code Review Guidelines

### For Reviewers

- Be constructive and respectful
- Explain reasoning behind suggestions
- Focus on code quality and maintainability
- Approve when ready, request changes if needed

### For Contributors

- Respond to all review comments
- Ask for clarification if needed
- Don't take criticism personally
- Update PR based on feedback

## Release Process

(For maintainers)

1. Update version in `cmd/cli/main.go`
2. Update CHANGELOG.md
3. Create and push tag: `git tag -a v0.2.0 -m "Release v0.2.0"`
4. Push tag: `git push origin v0.2.0`
5. Create GitHub release with notes

## Questions?

- **General Questions**: Open a discussion on GitHub
- **Bug Reports**: Open an issue with bug template
- **Feature Requests**: Open an issue with feature template
- **Security Issues**: Email maintainers directly

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Acknowledgments

Thank you for contributing to Changelog Generator! Every contribution helps make this tool better for everyone.
