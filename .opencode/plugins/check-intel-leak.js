// check-intel-leak — block Write/Edit tool calls that embed internal ticket IDs
// in scoped files (skills/, docs/, .context/).
//
// Mirrors .claude/hooks/check-intel-leak.sh PreToolUse behaviour.
// Uses permission.ask to deny the tool call when leaks are detected.

const TICKET_RE = /\b[A-Z]{2,6}-\d{3,}\b/g
const IGNORE_RE = /^(SHA|RFC|ISO|TLS|SSL|AES|RSA|UTF|BOM|ACL|DNS|API)-\d/

const SCOPED_DIRS = ["skills/", "docs/", ".context/"]

function inScope(filePath) {
  return SCOPED_DIRS.some((dir) => filePath.startsWith(dir) || filePath.includes(`/${dir}`))
}

function findLeaks(content) {
  const matches = []
  const lines = content.split("\n")
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    // Skip comment lines
    if (/^\s*#/.test(line) || /:\s*#/.test(line)) continue
    const tickets = [...line.matchAll(TICKET_RE)].map((m) => m[0]).filter((t) => !IGNORE_RE.test(t))
    if (tickets.length > 0) {
      matches.push(`  line ${i + 1}: ${line.trim()} (found: ${tickets.join(", ")})`)
    }
  }
  return matches
}

export const CheckIntelLeakPlugin = async () => {
  return {
    "permission.ask": async (input, output) => {
      const type = (input.type ?? "").toLowerCase()
      if (type !== "write" && type !== "edit") return

      const filePath = String(input.metadata?.file_path ?? input.metadata?.path ?? "")
      if (!filePath || !inScope(filePath)) return

      // For Write: full content. For Edit: the new_string being inserted.
      const content = String(
        input.metadata?.content ?? input.metadata?.new_string ?? "",
      )
      if (!content) return

      const leaks = findLeaks(content)
      if (leaks.length === 0) return

      output.status = "deny"
      // Surface reason via title mutation (best-effort — metadata is read-only after deny)
      console.error(
        `[check-intel-leak] BLOCKED ${filePath} — internal ticket IDs detected:\n${leaks.join("\n")}\nReplace with generic slugs before writing.`,
      )
    },
  }
}
