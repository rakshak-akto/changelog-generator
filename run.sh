#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BINARY="$SCRIPT_DIR/bin/changelog-generator"

# ── Load .env ────────────────────────────────────────────────────────
if [ -f "$SCRIPT_DIR/.env" ]; then
    set -a
    source "$SCRIPT_DIR/.env"
    set +a
else
    echo "Error: .env file not found. Create a .env with GITHUB_TOKEN, OPENAI_API_KEY, REPO_OWNER, REPO_NAME."
    exit 1
fi

# ── Preflight checks ────────────────────────────────────────────────
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Install it from https://golang.org/dl/"
    exit 1
fi

if [ -z "$GITHUB_TOKEN" ] || [ "$GITHUB_TOKEN" = "ghp_your_token_here" ]; then
    echo "Error: Set a valid GITHUB_TOKEN in .env"
    exit 1
fi

if [ -z "$OPENAI_API_KEY" ] || [ "$OPENAI_API_KEY" = "sk-your_key_here" ]; then
    echo "Error: Set a valid OPENAI_API_KEY in .env"
    exit 1
fi

# ── Build if binary is missing or source is newer ────────────────────
needs_build=false

if [ ! -f "$BINARY" ]; then
    needs_build=true
else
    newer=$(find "$SCRIPT_DIR" -name '*.go' -newer "$BINARY" -print -quit 2>/dev/null)
    if [ -n "$newer" ]; then
        needs_build=true
    fi
fi

if [ "$needs_build" = true ]; then
    echo "Building changelog-generator..."
    (cd "$SCRIPT_DIR" && go build -o "$BINARY" ./cmd/cli)
    echo "Build complete."
    echo
fi

# ── Prompt for dates ─────────────────────────────────────────────────
echo "=== Changelog Generator ($REPO_OWNER/$REPO_NAME) ==="
echo

read -rp "Start date (YYYY-MM-DD): " from_date
read -rp "End date   (YYYY-MM-DD): " to_date
echo

# ── Run ──────────────────────────────────────────────────────────────
echo "Generating changelog for $from_date to $to_date ..."
echo
exec "$BINARY" generate \
    --from-date="$from_date" \
    --to-date="$to_date" \
    --owner="$REPO_OWNER" \
    --repo="$REPO_NAME" \
    --verbose
