#!/bin/bash
# Install Git hooks for the Sentire project

set -e

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "Error: This script must be run from the root of the git repository"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
echo "Installing pre-commit hook..."
cp scripts/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

echo "✅ Pre-commit hook installed successfully!"
echo ""
echo "The hook will now automatically format Go files before each commit."
echo "To disable the hook temporarily, use: git commit --no-verify"
echo ""
echo "You can also format all Go files manually with: make fmt"