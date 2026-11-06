#!/usr/bin/env bash set -e # Exit immediately if a command exits with a non-zero status

export PATH="$HOME/.local/bin:$PATH"

echo "Starting post-create setup..."

if ! command -v uv >/dev/null 2>&1; then
  echo "Installing uv..."
  curl -LsSf https://astral.sh/uv/install.sh | sh
  if ! grep -q 'export PATH="$HOME/.local/bin:$PATH"' "$HOME/.bashrc" 2>/dev/null; then
    echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$HOME/.bashrc"
  fi
fi

# add codegraph installation
curl -fsSL https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh | sh
codegraph init
codegraph sync

echo "Setup finished successfully!"