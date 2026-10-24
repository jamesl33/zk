#!/bin/sh
# Every NUDGE_TURNS finished turns, makes Claude save what it learned before stopping.
# The headless save run (precompact.sh) must not trigger hooks.
[ -z "$ZK_MEMORY_SAVING" ] || exit 0

. "$(dirname "$0")/vault.sh" || exit 0

NUDGE_TURNS="${ZK_MEMORY_NUDGE_TURNS:-10}"

input=$(cat)

# This stop was itself caused by a nudge, don't loop.
[ "$(printf %s "$input" | jq -r '.stop_hook_active // false')" = "false" ] || exit 0

session=$(printf %s "$input" | jq -r '.session_id // empty')
[ -n "$session" ] || exit 0

state="$V/.zk/state/$session"
turns=$(($(cat "$state" 2>/dev/null || echo 0) + 1))

if [ "$turns" -lt "$NUDGE_TURNS" ]; then
	echo "$turns" >"$state"
	exit 0
fi

echo 0 >"$state"

jq -n '{decision: "block", reason: "Before finishing: save any durable, cross-project learnings from this session to zk-memory (create_note/update_note, per the memory instructions), or reply \"nothing to save\"."}'
