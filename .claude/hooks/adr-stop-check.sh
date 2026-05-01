#!/usr/bin/env bash
# Before Claude stops, check for unfilled ADR sections on the current branch.

BRANCH=$(git branch --show-current 2>/dev/null) || exit 0
[[ -z "$BRANCH" ]] && exit 0

TYPE="${BRANCH%%/*}"
SLUG="${BRANCH#*/}"

ADR_DIR=".context/decisions/branches"
ADR_FILE="$ADR_DIR/$TYPE/$SLUG.md"
[[ -f "$ADR_FILE" ]] || exit 0  # No ADR — nothing to check

# Collect section headers whose content is still only a comment placeholder
EMPTY_SECTIONS=$(awk '
  /^## / { section = substr($0, 4); next }
  /^### / { section = substr($0, 5); next }
  /^\[.*\]$/ || /^\[.*\]/ { placeholder=1; next }
  /^\s*$/ { next }
  { placeholder=0 }
  END { }
' "$ADR_FILE")

# Simpler: grep for lines that are only HTML/markdown comments (the fill-in markers)
UNFILLED=$(grep -n '^\s*<!-- .*-->$\|^\s*\[.*\]\s*$' "$ADR_FILE" | grep -v '^[0-9]*:\s*$' | wc -l | tr -d ' ')

[[ "$UNFILLED" -eq 0 ]] && exit 0

printf '{"hookSpecificOutput":{"hookEventName":"Stop","additionalContext":"ADR for branch \"%s\" (.context/decisions/branches/%s/%s.md) has %s unfilled placeholder(s). Consider updating it before ending the session."}}\n' \
  "$BRANCH" "$TYPE" "$SLUG" "$UNFILLED"
