// code-review-graph-update — run `code-review-graph update --skip-flows` after
// any file edit or write, keeping the knowledge graph incrementally up to date.
//
// Mirrors the PostToolUse Edit|Write|Bash hook in .claude/settings.json.
// Uses file.edited events rather than tool.execute.after to stay lightweight:
// OpenCode emits file.edited after actual filesystem changes.

const DEBOUNCE_MS = 5000 // coalesce rapid saves into one graph update

let debounceTimer = null

export const CodeReviewGraphUpdatePlugin = async ({ $, directory }) => {
  return {
    event: async ({ event }) => {
      if (event.type !== "file.edited") return

      // Debounce: wait for a quiet period before running the update
      if (debounceTimer) clearTimeout(debounceTimer)

      debounceTimer = setTimeout(async () => {
        debounceTimer = null
        try {
          await $`code-review-graph update --skip-flows`.cwd(directory).quiet().nothrow()
        } catch {
          // non-fatal — graph may not be initialised yet
        }
      }, DEBOUNCE_MS)
    },
  }
}
