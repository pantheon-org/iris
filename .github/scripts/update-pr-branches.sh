#!/usr/bin/env bash
# update-pr-branches.sh
#
# For every open PR targeting main, trigger a merge-update via the GitHub API
# (equivalent to clicking "Update branch" in the UI).
#
# Required env vars:
#   GH_TOKEN  — GitHub token with `contents:write` and `pull-requests:read` scopes
#   GITHUB_REPOSITORY — owner/repo (e.g. pantheon-org/iris), set automatically in Actions
#
# Exit codes:
#   0 — all PRs updated (or already up to date)
#   1 — one or more PRs could not be updated (conflict or API error)

set -euo pipefail

: "${GH_TOKEN:?GH_TOKEN must be set}"
: "${GITHUB_REPOSITORY:?GITHUB_REPOSITORY must be set}"

BASE_BRANCH="${BASE_BRANCH:-main}"
API="https://api.github.com/repos/${GITHUB_REPOSITORY}"
AUTH_HEADER="Authorization: Bearer ${GH_TOKEN}"

echo "Fetching open PRs targeting '${BASE_BRANCH}'..."

prs=$(gh api \
  --paginate \
  "repos/${GITHUB_REPOSITORY}/pulls?state=open&base=${BASE_BRANCH}" \
  --jq '.[] | {number: .number, title: .title, head: .head.ref}')

if [[ -z "$prs" ]]; then
  echo "No open PRs found targeting '${BASE_BRANCH}'. Nothing to do."
  exit 0
fi

failed=0

while IFS= read -r pr_json; do
  number=$(echo "$pr_json" | jq -r '.number')
  title=$(echo "$pr_json"  | jq -r '.title')
  branch=$(echo "$pr_json" | jq -r '.head')

  echo ""
  echo "PR #${number} — ${title} (branch: ${branch})"

  response=$(curl --silent --write-out "\n%{http_code}" \
    -X PUT \
    -H "${AUTH_HEADER}" \
    -H "Accept: application/vnd.github+json" \
    -H "X-GitHub-Api-Version: 2022-11-28" \
    "${API}/pulls/${number}/update-branch" \
    -d '{}')

  body=$(echo "$response" | head -n -1)
  http_code=$(echo "$response" | tail -n 1)

  case "$http_code" in
    202)
      echo "  -> Merge update queued successfully."
      ;;
    422)
      echo "  -> Skipped: branch has a merge conflict with '${BASE_BRANCH}'."
      echo "     ${body}"
      failed=1
      ;;
    *)
      echo "  -> Unexpected response (HTTP ${http_code}):"
      echo "     ${body}"
      failed=1
      ;;
  esac
done < <(echo "$prs" | jq -c '.')

echo ""
if [[ "$failed" -eq 0 ]]; then
  echo "All PRs processed successfully."
else
  echo "One or more PRs could not be updated — see above for details."
  exit 1
fi
