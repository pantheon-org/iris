Manage Architecture Decision Records (ADRs) for this project following the Claude ADR System Guide.

The ADR system lives entirely under `.claude/`:

- `.claude/adr-index.toml` — master index (metadata only, no content)
- `.context/decisions/branches/<type>/<slug>.md` — active branch ADRs
- `.context/decisions/merged/YYYY-MM/<slug>.md` — archived post-merge ADRs

The argument passed to this command determines the action: `$ARGUMENTS`

---

## init

Create the ADR directory structure and a skeleton `adr-index.toml` if they don't already exist.

1. Create: `.context/decisions/branches/feat/`, `.context/decisions/branches/chore/`, `.context/decisions/branches/docs/`, `.context/decisions/branches/fix/`, `.context/decisions/merged/`
2. Create `.claude/adr-index.toml` with this skeleton (skip if file already exists):

```toml
[meta]
project = "<repo name from git remote or directory name>"
created = "<today YYYY-MM-DD>"
version = "1.0"

[active.feat]
[active.chore]
[active.docs]
[active.fix]

[merged]

[categories]

[cross_references]
```

1. Confirm the paths created.

---

## create

Create an ADR for the current branch.

1. Run `git branch --show-current` to get the branch name.
2. Parse type prefix (`feat`, `chore`, `docs`, `fix`). Default to `chore` if unrecognised.
3. Slug = branch name with the `<type>/` prefix stripped.
4. Create `.context/decisions/branches/<type>/<slug>.md` using this template (abort if file already exists and confirm with user first):

```markdown
# Branch ADR: <branch-name>

## Meta
- **Branch**: `<branch-name>`
- **Type**: <type>
- **Status**: active
- **Created**: <YYYY-MM-DD>
- **Author**: Claude
- **PR**: <!-- Add PR number when opened -->

## Problem Statement
### Context
<!-- Describe the situation requiring this branch -->

### Goals
<!-- What this branch must achieve -->

### Non-Goals
<!-- Explicitly out of scope -->

## Decision Record
### Options Considered
<!-- List alternatives evaluated -->

### Chosen Solution
<!-- What was decided and implemented -->

### Rationale
<!-- Why this option over alternatives -->

## Implementation
### Key Changes
<!-- Files/modules changed -->

### Testing Strategy
<!-- How correctness is verified -->

## Challenges & Solutions
<!-- Technical or process obstacles encountered -->

## Impact Assessment
- **Performance**: <!-- impact or "none" -->
- **Security**: <!-- impact or "none" -->
- **Maintenance**: <!-- impact or "none" -->

## Outcome & Lessons
<!-- Fill in post-merge: results, metrics, lessons learned -->
```

1. Add an entry to `[active.<type>]` in `.claude/adr-index.toml`. Preserve existing entries:

```toml
"<slug>" = { file = "<type>/<slug>.md", created = "<today>", author = "Claude", tags = [], description = "" }
```

1. Report the file path created.

---

## update

Update the current branch's ADR.

1. Detect branch → find `.context/decisions/branches/<type>/<slug>.md`. If missing, offer to run `adr create` first.
2. If the user provided text after `update` in `$ARGUMENTS`, treat it as notes and place them in the most appropriate section.
3. Otherwise, show the current ADR content and ask which section to update.
4. Write the changes. If the user says the branch is merging, set `Status: merged` in the Meta section.

---

## list

List all ADRs from `.claude/adr-index.toml`.

Parse the TOML and print a table with columns: `TYPE`, `SLUG`, `CREATED`, `STATUS`, `DESCRIPTION`.
Show active ADRs first, then merged ones (labelled).
If the index doesn't exist, say so and suggest running `/adr init`.

---

## archive

Archive the current branch's ADR after it merges.

1. Detect branch → slug + type. Accept an explicit branch name as argument override.
2. Determine merge date (today's date).
3. Move `.context/decisions/branches/<type>/<slug>.md` → `.context/decisions/merged/YYYY-MM/<slug>.md` (create the month dir if needed).
4. Update `.claude/adr-index.toml`:
   - Remove from `[active.<type>]`
   - Add to `[merged]`: `"<slug>" = { file = "merged/YYYY-MM/<slug>.md", merged = "<YYYY-MM-DD>", pr = <number if known> }`
5. Confirm the move.

---

## status

Show the ADR for the current branch with completeness indicators.

1. Find `.context/decisions/branches/<type>/<slug>.md`. If missing, say so.
2. Display all sections.
3. For each section whose body consists only of a `<!-- ... -->` comment (unfilled placeholder), prefix it with `⚠ TODO`.
4. End with a summary: `N of M sections complete`.
