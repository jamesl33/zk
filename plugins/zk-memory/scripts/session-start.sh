#!/bin/sh
# Ensures the vault exists and gives Claude its memory instructions (also after /clear and compaction).
# The headless save run (extract.sh) must not trigger hooks.
[ -z "$ZK_MEMORY_SAVING" ] || exit 0

. "$(dirname "$0")/vault.sh" || exit 0

# Drop stale per-session state.
find "$V/.zk/state" -type f -mtime +7 -delete 2>/dev/null

sed "s|@VAULT@|$V|g" "$(dirname "$0")/../instructions.md" |
	jq -Rs '{hookSpecificOutput: {hookEventName: "SessionStart", additionalContext: .}}'
