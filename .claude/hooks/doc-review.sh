#!/usr/bin/env bash
# doc-review — prompt Claude to check if docs need updating at session end.

raw=$(cat)
session_id=$(jq -r '.session_id // "unknown"' <<<"$raw")
flag="/tmp/civitas-doc-review-${session_id}"

# Already ran this session — let Claude stop normally
[[ -f "$flag" ]] && exit 0

touch "$flag"

jq -n '{
  decision: "block",
  reason: "Before finishing: review the changes made in this session. If any source code, configuration, or feature changes were made, check whether README.md or any docs/*.md files need to be updated to reflect those changes accurately. Update them if needed, then complete your response."
}'
