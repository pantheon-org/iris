// code-review-graph-status — run `code-review-graph status` once at session start
// to surface the current graph state in the context window.
//
// Mirrors the SessionStart `code-review-graph status` hook in .claude/settings.json.
// Fires on session.created (the first event for a new session).

export const CodeReviewGraphStatusPlugin = async ({ $, directory }) => {
  return {
    event: async ({ event }) => {
      if (event.type !== "session.created") return

      try {
        await $`code-review-graph status`.cwd(directory).nothrow()
      } catch {
        // non-fatal — graph may not be initialised yet
      }
    },
  }
}
