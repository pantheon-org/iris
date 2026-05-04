// doc-review — prompt OpenCode to check if docs need updating when the session goes idle.
// Mirrors .claude/hooks/doc-review.sh behaviour: runs once per session on idle,
// injects the doc-review prompt targeting README.md and AGENTS.md specifically.

export const DocReviewPlugin = async ({ client, $, directory }) => {
  return {
    event: async ({ event }) => {
      if (event.type !== "session.idle") return

      const sessionId = event.properties?.id
      if (!sessionId) return

      // Use a project-local flag file to deduplicate within the session
      const flag = `${directory}/.opencode/.doc-review-${sessionId}`

      const flagExists = await $`test -f ${flag}`.quiet().nothrow()
      if (flagExists.exitCode === 0) return

      await $`touch ${flag}`.quiet()

      await client.session.prompt({
        path: { id: sessionId },
        body: {
          parts: [
            {
              type: "text",
              text: "Before finishing: review the changes made in this session. If any source code, configuration, or feature changes were made, check whether README.md and AGENTS.md need to be updated to reflect those changes accurately. Update them if needed, then complete your response.",
            },
          ],
        },
      })
    },
  }
}
