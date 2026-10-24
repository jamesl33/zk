#!/bin/sh
# Before context is summarised, a headless Claude reads the transcript and saves durable learnings, then compaction proceeds. Fails open.
[ -z "$ZK_MEMORY_SAVING" ] || exit 0

. "$(dirname "$0")/vault.sh" || exit 0

transcript=$(jq -r '.transcript_path // empty')
[ -r "$transcript" ] || exit 0

prompt=$(mktemp) || exit 0
trap 'rm -f "$prompt"' EXIT

{
	sed "s|@VAULT@|$V|g" "$(dirname "$0")/../instructions.md"
	printf '\nThe session below is about to be compacted. Save its durable, cross-project learnings with the zk-memory tools, per the instructions above, checking for existing notes first. Save nothing if there are none. Reply "done" when finished.\n\n<transcript>\n'
	jq -r 'select(.type == "user" or .type == "assistant") | "\(.type): " + (if (.message.content | type) == "string" then .message.content else ([.message.content[]? | select(.type == "text") | .text] | join("\n")) end)' "$transcript" | tail -c 100000
	printf '</transcript>\n'
} >"$prompt"

mcp=$(jq -n --arg cmd "$(dirname "$0")/mcp.sh" '{mcpServers: {"zk-memory": {command: $cmd, args: []}}}')

ZK_MEMORY_SAVING=1 claude -p --no-session-persistence --strict-mcp-config --mcp-config "$mcp" --allowedTools 'mcp__zk-memory__*' <"$prompt" >/dev/null 2>&1

exit 0
