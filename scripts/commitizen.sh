#!/bin/bash

set -e

COMMIT_MSG_FILE=$1

if ! [ -f "$COMMIT_MSG_FILE" ]; then
    echo "Cannot read commit message file: $COMMIT_MSG_FILE"
    exit 1
fi

# Read the first line of the commit message
COMMIT_MSG=$(head -n1 "$COMMIT_MSG_FILE" 2>/dev/null || echo "")

# Skip if empty
if [ -z "$COMMIT_MSG" ]; then
    echo "Empty commit message, skipping validation"
    exit 0
fi

# Skip merge commits
if echo "$COMMIT_MSG" | grep -qE '^Merge (branch|pull request)'; then
    echo "Merge commit detected, skipping validation"
    exit 0
fi

# Validate conventional commits format
commit_regex='^(feat|fix|docs|style|refactor|perf|test|chore|build|ci|revert)(\(.+\))?(!)?: .{1,50}'

if ! echo "$COMMIT_MSG" | grep -qE "$commit_regex"; then
    echo ""
    echo "Invalid commit message format!"
    echo ""
    echo "Commit messages must follow the Conventional Commits specification:"
    echo "  <type>[optional scope]: <description>"
    echo ""
    echo "Types: feat, fix, docs, style, refactor, perf, test, chore, build, ci, revert"
    echo ""
    echo "Examples:"
    echo "  feat: add new TUI navigation component"
    echo "  fix(installer): resolve disk partition detection"
    echo "  docs: update installation instructions"
    echo "  feat!: introduce breaking API changes"
    echo ""
    echo "Your commit message:"
    echo "  '$COMMIT_MSG'"
    echo ""
    exit 1
fi

# Check commit message length
MSG_LEN=${#COMMIT_MSG}
if [ "$MSG_LEN" -gt 72 ]; then
    echo "Commit message is too long (> 72 characters)"
    echo "Current length: $MSG_LEN characters"
    echo "Your commit message:"
    echo "  '$COMMIT_MSG'"
    exit 1
fi
