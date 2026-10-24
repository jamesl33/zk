#!/bin/sh
# On PreCompact and SessionEnd, a detached headless Claude reads the part of the transcript not yet seen and saves durable learnings, then the vault is committed. Fails open.
[ -z "$ZK_MEMORY_SAVING" ] || exit 0

# Detach, so neither compaction nor exit waits on (or kills) the run.
if [ -z "$ZK_MEMORY_BG" ]; then
	input=$(cat)
	printf %s "$input" | ZK_MEMORY_BG=1 setsid "$0" >/dev/null 2>&1 &
	exit 0
fi

. "$(dirname "$0")/vault.sh" || exit 0

MODEL="${ZK_MEMORY_MODEL:-haiku}"
MAX_CHUNKS=5

input=$(cat)
transcript=$(printf %s "$input" | jq -r '.transcript_path // empty')
session=$(printf %s "$input" | jq -r '.session_id // empty')
[ -r "$transcript" ] && [ -n "$session" ] || exit 0

# Transcript lines already extracted, so PreCompact and SessionEnd never repeat work.
state="$V/.zk/state/$session.extracted"
seen=$(cat "$state" 2>/dev/null || echo 0)
total=$(wc -l <"$transcript")
[ "$total" -gt "$seen" ] || exit 0

work=$(mktemp -d) || exit 0
trap 'rm -rf "$work"' EXIT

# Tool calls and results hold the gotchas, so keep them, clipped.
tail -n +$((seen + 1)) "$transcript" | jq -r '
	def clip: if length > 500 then .[:500] + "…" else . end;
	def flat: if type == "string" then . else ([.[]? | .text? // empty] | join("\n")) end;
	def render: if type == "string" then . else ([.[]? |
		if .type == "text" then .text
		elif .type == "tool_use" then "[tool \(.name)] \(.input | tojson | clip)"
		elif .type == "tool_result" then "[result] \(.content | flat | clip)"
		else empty end] | join("\n")) end;
	select(.type == "user" or .type == "assistant") | "\(.type): " + (.message.content | render)
' | split -C 100000 - "$work/chunk."

mcp=$(jq -n --arg cmd "$(dirname "$0")/mcp.sh" '{mcpServers: {"zk-memory": {command: $cmd, args: []}}}')

n=0
for chunk in "$work"/chunk.*; do
	[ -e "$chunk" ] || break
	# ponytail: transcripts past MAX_CHUNKS chunks lose their tail, raise it if that bites
	n=$((n + 1))
	[ "$n" -le "$MAX_CHUNKS" ] || break

	{
		sed "s|@VAULT@|$V|g" "$(dirname "$0")/../instructions.md"
		printf '\n# Vault note format (ZK.md)\n\n'
		cat "$V/ZK.md"
		printf '\nThe session transcript below (tool calls and results clipped) is being closed or compacted. Save its durable, cross-project learnings with the zk-memory tools, per the instructions above, checking for existing notes first. Save nothing if there are none. Reply "done" when finished.\n\n<transcript>\n'
		cat "$chunk"
		printf '</transcript>\n'
	} >"$work/prompt"

	ZK_MEMORY_SAVING=1 claude -p --model "$MODEL" --no-session-persistence --strict-mcp-config --mcp-config "$mcp" --allowedTools 'mcp__zk-memory__*' <"$work/prompt" >/dev/null 2>&1
done

# Notes must pass lint (e.g. no orphans); give the model one pass to fix what it left.
lint=$(zk --vault "$V" lint 2>&1)
if [ -n "$lint" ]; then
	{
		sed "s|@VAULT@|$V|g" "$(dirname "$0")/../instructions.md"
		printf '\n# Vault note format (ZK.md)\n\n'
		cat "$V/ZK.md"
		printf '\n`zk lint` reports problems in the vault. Fix every one with the zk-memory tools, per the instructions above, then reply "done".\n\n%%s\n' "$lint"
	} >"$work/prompt"
	ZK_MEMORY_SAVING=1 claude -p --model "$MODEL" --no-session-persistence --strict-mcp-config --mcp-config "$mcp" --allowedTools 'mcp__zk-memory__*' <"$work/prompt" >/dev/null 2>&1
fi

echo "$total" >"$state"

git -C "$V" add -A && git -C "$V" -c commit.gpgsign=false -c user.name=Claude -c user.email=claude@localhost commit -qm "Save memories" 2>/dev/null
exit 0
