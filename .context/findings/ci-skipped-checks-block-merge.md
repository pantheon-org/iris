# Finding: Skipped CI checks block merge on release-please PRs

**Date:** 2026-05-01  
**Affected PRs:** #42 (release 2026.5.4), any docs/chore PR with no `.go` changes

## Problem

Release-please PRs (and docs-only PRs) touch only `CHANGELOG.md` and version files — no `.go`, `go.mod`, or `go.sum` changes. The `test` and `build` jobs in `ci.yml` were gated behind:

```yaml
if: needs.changes.outputs.code == 'true'
```

When no code files change, these jobs report a `skipped` conclusion. GitHub branch protection treats a **skipped required check as unsatisfied**, permanently blocking merge. The PR shows "Waiting for status to be reported" with no way to unblock it without admin intervention.

## Root cause

`test` and `build` were listed as required checks in branch protection. GitHub does not distinguish between `skipped` and `not run` — both block merge.

## Fix (merged in #45)

Added a `ci-ok` aggregator job that always runs (`if: always()`) and passes when `test`/`build` are either `success` or `skipped`:

```yaml
ci-ok:
  runs-on: ubuntu-latest
  needs: [lint, test, build]
  if: always()
  steps:
    - name: check required jobs
      run: |
        if [[ "${{ needs.lint.result }}" != "success" ]]; then exit 1; fi
        if [[ "${{ needs.test.result }}" != "success" && "${{ needs.test.result }}" != "skipped" ]]; then exit 1; fi
        if [[ "${{ needs.build.result }}" != "success" && "${{ needs.build.result }}" != "skipped" ]]; then exit 1; fi
```

Branch protection required checks updated from `lint + test + build` → `lint + ci-ok`.

## Unblocking the stuck PR (#42)

The fix landed after PR #42 was already open. The release-please branch only gets a new CI run on a `push` event — a comment or review does not trigger it. Since `enforce_admins` was disabled in branch protection, PR #42 was merged with `gh pr merge --admin`.

## Prevention

- Never list jobs with `if:` skip conditions as required checks directly.
- Always use an aggregator job (`ci-ok` pattern) as the single required check gate.
- The `ci.yml` push trigger on `release-please--branches--*` ensures future release PRs get a CI run, but only after a new push to that branch.
