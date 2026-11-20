#!/usr/bin/env bash
# Enforces Conventional Commits and the no-AI-attribution rule (AGENTS.md §6).
set -euo pipefail

msg_file="$1"
subject="$(head -n1 "$msg_file")"
pattern='^(feat|fix|docs|chore|refactor|test|perf|build|ci|revert)(\([a-z0-9-]+\))?!?: .{1,72}$'

if [[ "$subject" =~ ^(Merge|Revert) ]]; then
  exit 0
fi

if ! [[ "$subject" =~ $pattern ]]; then
  echo "Commit subject must follow Conventional Commits, max 72 chars after the type:" >&2
  echo "  feat(api): add entry search" >&2
  echo "Got: $subject" >&2
  exit 1
fi

if grep -qiE '^Co-Authored-By: .*(claude|anthropic|copilot|openai)' "$msg_file"; then
  echo "Remove the AI co-author trailer (AGENTS.md §6)." >&2
  exit 1
fi
