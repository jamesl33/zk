#!/bin/sh
# Injects the notes closest to the prompt. Fails open: any problem injects nothing.
# The headless save run (extract.sh) must not trigger hooks.
[ -z "$ZK_MEMORY_SAVING" ] || exit 0

. "$(dirname "$0")/vault.sh" || exit 0

# Cosine distance beyond which a note is not relevant, so trivial prompts ("yes", "continue") inject nothing.
# Calibrated (embeddinggemma, 4-note vault) as D = mean - k*stddev of the top-1 distance of unrelated prompts:
# mean 0.83, stddev 0.07 (n=7), k=2.5 gives 0.65. Related prompts matched at 0.45-0.60. Tune via the environment.
MAX_DISTANCE="${ZK_MEMORY_MAX_DISTANCE:-0.65}"

input=$(cat)
prompt=$(printf %s "$input" | jq -r '.prompt // empty') || exit 0
session=$(printf %s "$input" | jq -r '.session_id // empty')

# Skip what cannot match, saving the embedding round trip.
case "$prompt" in
'' | /*) exit 0 ;;
esac

[ "$(printf %s "$prompt" | wc -w)" -ge 3 ] || exit 0

query=$(printf %s "$prompt" | head -c 2000)

hits=$(timeout 3 zk --vault "$V" notes find --json --limit 5 --max-distance "$MAX_DISTANCE" -- "$query" 2>/dev/null) || exit 0

[ -n "$hits" ] || exit 0

# Each hit is a note as JSON (path and frontmatter); the body is the file past its frontmatter.
# Notes already injected in this session are skipped.
[ -n "$session" ] || exit 0
recalled="$V/.zk/state/$session.recalled"
touch "$recalled"

context=$(printf %s "$hits" | jq -r '[.path, .frontmatter.date // "unknown", .frontmatter.title] | @tsv' | while IFS="$(printf '\t')" read -r path date title; do
	grep -Fxq -- "$path" "$recalled" && continue
	printf '%s\n' "$path" >>"$recalled"
	printf '### %s (%s, saved %s)\n' "$title" "$path" "$date"
	awk 'body {print; next} /^---$/ && ++markers == 2 {body = 1}' "$V/$path"
done | head -c 2000)

[ -n "$context" ] || exit 0

printf '<memory>\nNotes from Claude'"'"'s memory vault that may be relevant. Older notes may be stale: verify anything about code or state before relying on it. If a note is wrong or outdated, say so in the form "Memory note <path> is stale: <what is wrong>" so it is removed.\n\n%s\n</memory>' "$context" |
	jq -Rs '{hookSpecificOutput: {hookEventName: "UserPromptSubmit", additionalContext: .}}'
