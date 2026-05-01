# Branch ADR: release-please setup

## Meta
- **Branch**: `fix/release-please-*` (multiple)
- **Type**: fix
- **Status**: merged
- **Created**: 2026-05-01
- **Merged**: 2026-05-01
- **Author**: Claude
- **PRs**: #4, #6, #10

## Problem Statement
### Context
After introducing CalVer via `feat/calver-releases` (PR #3), the release-please GitHub Actions workflow required several config corrections before it produced a stable release. This ADR captures the collective decisions made across PRs #4, #6, and #10.

### Goals
- Stable, automated release-please workflow
- No manual version management
- Config that survives re-runs without corruption

### Non-Goals
- Changing the CalVer scheme (covered in `calver-releases` ADR)

## Decision Record
### Corrections Made

**PR #4 — Remove `bootstrap_sha` from manifest**
- `bootstrap_sha` in `.release-please-manifest.json` caused release-please to ignore commits before that SHA, producing incorrect changelogs. Removed entirely.
- Decision: do not set `bootstrap_sha` unless explicitly needed for a brownfield migration.

**PR #6 — Move `versioning-strategy` to workflow input**
- `versioning-strategy: always-bump-patch` was set in `release-please-config.json` but release-please v17 expects it as a workflow `inputs` parameter for the action, not in the config file.
- Decision: pass `versioning-strategy` as a `with:` input to the `google-github-actions/release-please-action` step.

**PR #10 — Correct `versioning` key name**
- The config file used key `versioning` where release-please v17 expects `versioning-strategy`.
- Decision: align config key names with the release-please v17 schema; validate against schema before merging CI config changes.

### Chosen Solution
Final stable config:
- `.release-please-manifest.json`: CalVer seed value, no `$schema`, no `bootstrap_sha`
- `release-please-config.json`: correct key names per v17 schema
- `release-please.yml`: CalVer pre-step + `versioning-strategy` passed as workflow input

## Implementation
### Key Changes
- `.release-please-manifest.json` — stripped to minimal valid CalVer seed
- `release-please-config.json` — corrected key names
- `.github/workflows/release-please.yml` — `versioning-strategy` as `with:` input; `version.txt` added as release artifact trigger

## Challenges & Solutions
- Multiple config keys changed between release-please v15 and v17 — the migration guide was not consulted upfront.

## Impact Assessment
- **Performance**: None.
- **Security**: None.
- **Maintenance**: Caution — release-please config is sensitive to version-specific key names. Pin the action version and re-validate config when upgrading.

## Outcome & Lessons
- First successful release `2026.5.1` published after PR #10.
- **Lesson**: validate release-please config in a dry-run or against the v17 schema before merging. Three fix PRs could have been one.
