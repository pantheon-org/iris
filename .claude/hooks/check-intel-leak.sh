#!/usr/bin/env bash
# check-intel-leak — detect internal ticket IDs before files are written.
#
# Two modes:
#   Claude hook (no args): reads JSON from stdin, called by PreToolUse hook.
#   File scan (args):      called with file paths for manual/CI runs.
#
# Exit 0 = clean/allow.  Exit 2 = block (hard) with JSON output on stdout.

TICKET_RE='[A-Z]{2,6}-[0-9]{3,}'
IGNORE_RE='^(SHA|RFC|ISO|TLS|SSL|AES|RSA|UTF|BOM|ACL|DNS|API)-[0-9]'

in_scope() {
  case "$1" in
    skills/*|docs/*|.context/*) return 0 ;;
    *) return 1 ;;
  esac
}

find_leaks() {
  local content="$1"
  grep -nE "$TICKET_RE" <<<"$content" \
    | grep -v '^\s*#' \
    | grep -Ev "$IGNORE_RE" \
    | grep -oE "[0-9]+:.*$TICKET_RE.*" \
    | grep -Ev ":[[:space:]]*#"
}

# ── Claude hook mode (no arguments) ──────────────────────────────────────────

if [[ $# -eq 0 ]]; then
  raw=$(cat)

  tool_name=$(jq -r '.tool_name // ""' <<<"$raw")
  file_path=$(jq -r '.tool_input.file_path // ""' <<<"$raw")

  in_scope "$file_path" || exit 0

  case "$tool_name" in
    Write) content=$(jq -r '.tool_input.content // ""' <<<"$raw") ;;
    Edit)  content=$(jq -r '.tool_input.new_string // ""' <<<"$raw") ;;
    *)     exit 0 ;;
  esac

  leaks=$(find_leaks "$content")
  [[ -z "$leaks" ]] && exit 0

  reason="Intel leak blocked in ${file_path}.\nFound ticket-style IDs (e.g. PROJ-123):\n${leaks}\n\nReplace with generic slugs before writing."

  jq -n --arg reason "$reason" '{
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "block",
      permissionDecisionReason: $reason
    }
  }'
  exit 2
fi

# ── File scan mode (arguments = file paths) ───────────────────────────────────

errors=0

for file in "$@"; do
  in_scope "$file" || continue
  [[ -f "$file" ]] || continue

  leaks=$(find_leaks "$(cat "$file")")
  if [[ -n "$leaks" ]]; then
    echo -e "\nINTEL LEAK in: $file\n$leaks" >&2
    ((errors++))
  fi
done

if ((errors > 0)); then
  echo -e "\n[BLOCKED] Remove internal ticket IDs before committing." >&2
  echo "Replace with generic slugs (e.g. PROJ-123 → cleanup-plan)." >&2
  exit 1
fi

exit 0
