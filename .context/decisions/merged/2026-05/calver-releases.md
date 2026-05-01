# Branch ADR: feat/calver-releases

## Meta
- **Branch**: `feat/calver-releases`
- **Type**: feat
- **Status**: merged
- **Created**: 2026-05-01
- **Merged**: 2026-05-01
- **Author**: Claude
- **PR**: #3

## Problem Statement
### Context
The project needed a versioning scheme for releases. The default for release-please is SemVer. Given that Iris is a tool (not a library), communicating year/month in the version is more useful to users than MAJOR.MINOR.PATCH semantics. Additionally, breaking changes are unlikely to be signalled via MAJOR bumps in a CLI tool context.

### Goals
- Version strings that encode release date (year + month)
- Automated release management via release-please
- Patch bumps for fixes within a calendar month

### Non-Goals
- SemVer MAJOR/MINOR semantics
- Manual version management

## Decision Record
### Options Considered

**SemVer (MAJOR.MINOR.PATCH)**
- release-please default; well-understood.
- Cons: MAJOR bumps carry library-breaking-change semantics that don't apply to a CLI tool; version numbers carry no temporal information.

**CalVer (YYYY.M.PATCH)**
- Encodes release year and month; patch auto-increments within the month.
- Cons: release-please has no built-in CalVer support — requires a custom pre-step.

### Chosen Solution
**CalVer `YYYY.M.PATCH`** (e.g. `2026.5.1`)

Implementation: a pre-step in the release-please GitHub Actions workflow rewrites `.release-please-manifest.json` to `YYYY.M.0` when year/month rolls over; release-please then applies `always-bump-patch` strategy from there.

### Rationale
Iris is a CLI tool, not a library. CalVer is more informative for end-users ("this release is from May 2026") and sidesteps the semantic weight of SemVer MAJOR bumps. The custom pre-step is small and self-contained.

## Implementation
### Key Changes
- `.github/workflows/release-please.yml` — added checkout + CalVer bump pre-step
- `release-please-config.json` — `versioning-strategy: always-bump-patch`
- `.release-please-manifest.json` — seeded with `2026.5.0`; removed `$schema` (release-please v17 parsed it as a version)

### Testing Strategy
- Validated by observing the release-please PR that produced `2026.5.1`

## Challenges & Solutions
- **`$schema` in manifest**: release-please v17 treated the `$schema` key as a version string, causing parse failures. Removed in this branch.
- **No native CalVer support**: worked around via a rewrite pre-step in CI rather than forking release-please.
- **`bootstrap_sha` not needed**: removed in follow-up PR #4 after it caused unnecessary CI failures.
- **`versioning-strategy` as workflow input**: moved to workflow input in PR #6 to allow override without config file changes.
- **`versioning` key naming**: corrected in PR #10 (config key was `versioning` not `versioning-strategy` in the config file).

## Impact Assessment
- **Performance**: None.
- **Security**: None.
- **Maintenance**: The CalVer pre-step is ~10 lines of shell in CI — low maintenance burden.

## Outcome & Lessons
- First release `2026.5.1` successfully published.
- Multiple follow-up fixes (#4, #6, #10) were needed to stabilise release-please config — suggests the release-please setup should be validated end-to-end in a throwaway branch before merging to main.
