#!/bin/bash
#
# Setup git hooks for the conreq project
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "🔧 Setting up git hooks..."

# Copy pre-push hook
if [ -f "$PROJECT_ROOT/scripts/pre-push" ]; then
    cp "$PROJECT_ROOT/scripts/pre-push" "$PROJECT_ROOT/.git/hooks/pre-push"
    chmod +x "$PROJECT_ROOT/.git/hooks/pre-push"
    echo "✅ pre-push hook installed"
else
    echo "❌ pre-push script not found in scripts directory"
    exit 1
fi

echo "✨ Git hooks setup complete!"
echo ""
echo "The following hooks are now active:"
echo "  - pre-push: Runs fmt, vet, lint, and test checks before pushing"