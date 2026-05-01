#!/usr/bin/env bash
# Remind Claude to create an ADR if none exists for the current branch.

BRANCH=$(git branch --show-current 2>/dev/null) || exit 0
[[ -z "$BRANCH" ]] && exit 0

TYPE="${BRANCH%%/*}"
SLUG="${BRANCH#*/}"

ADR_DIR=".context/decisions/branches"
[[ -d "$ADR_DIR" ]] || exit 0   # ADR system not initialised — silent no-op

ADR_FILE="$ADR_DIR/$TYPE/$SLUG.md"
[[ -f "$ADR_FILE" ]] && exit 0  # ADR exists — nothing to do

printf '{"hookSpecificOutput":{"hookEventName":"SessionStart","additionalContext":"No ADR found for branch \"%s\" (.context/decisions/branches/%s/%s.md). Offer to create one using the ADR system (.context/decisions/ + .claude/adr-index.toml)."}}\n' \
  "$BRANCH" "$TYPE" "$SLUG"
