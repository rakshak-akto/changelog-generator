#!/bin/bash
set -e

echo "==================================="
echo "Changelog Generator - Setup Script"
echo "==================================="
echo

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "Please install Go from https://golang.org/dl/"
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo

# Check for GitHub token
if [ -z "$GITHUB_TOKEN" ]; then
    echo "⚠️  Warning: GITHUB_TOKEN is not set"
    echo "Please set it with: export GITHUB_TOKEN=ghp_your_token_here"
    echo "Get a token from: https://github.com/settings/tokens"
    echo
else
    echo "✓ GITHUB_TOKEN is set"
    echo
fi

# Check for OpenAI API key
if [ -z "$OPENAI_API_KEY" ]; then
    echo "⚠️  Warning: OPENAI_API_KEY is not set"
    echo "Please set it with: export OPENAI_API_KEY=sk_your_key_here"
    echo "Get an API key from: https://platform.openai.com/api-keys"
    echo
else
    echo "✓ OPENAI_API_KEY is set"
    echo
fi

# Build the project
echo "Building changelog-generator..."
make build

if [ $? -eq 0 ]; then
    echo
    echo "✅ Build successful!"
    echo
    echo "Binary location: ./bin/changelog-generator"
    echo
    echo "Try it out with:"
    echo "  ./bin/changelog-generator --help"
    echo
    echo "Or generate a changelog:"
    echo "  ./bin/changelog-generator generate v1.0.0..v1.1.0 \\"
    echo "    --owner=facebook --repo=react --verbose"
    echo
else
    echo
    echo "❌ Build failed. Please check the errors above."
    exit 1
fi
