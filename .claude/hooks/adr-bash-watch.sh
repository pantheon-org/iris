#!/usr/bin/env bash
# After a Bash tool call, detect branch switches and remind Claude if the new branch has no ADR.

INPUT=$(cat)
CMD=$(printf '%s' "$INPUT" | jq -r '.tool_input.command // ""')

# Only act on git branch-switching commands
if ! printf '%s' "$CMD" | grep -qE '^\s*git\s+(checkout|switch)\s'; then
  exit 0
fi

# Give git a moment to settle, then read current branch
sleep 0.1
BRANCH=$(git branch --show-current 2>/dev/null) || exit 0
[[ -z "$BRANCH" ]] && exit 0

TYPE="${BRANCH%%/*}"
SLUG="${BRANCH#*/}"

ADR_DIR=".context/decisions/branches"
[[ -d "$ADR_DIR" ]] || exit 0

ADR_FILE="$ADR_DIR/$TYPE/$SLUG.md"
[[ -f "$ADR_FILE" ]] && exit 0

printf '{"hookSpecificOutput":{"hookEventName":"PostToolUse","additionalContext":"Switched to branch \"%s\" — no ADR found at .context/decisions/branches/%s/%s.md. Consider creating one."}}\n' \
  "$BRANCH" "$TYPE" "$SLUG"
