# Sourced by every script: sets V to the memory vault, creating it when missing. Idempotent.
V="${ZK_MEMORY_VAULT:-$HOME/.claude/zk-memory}"

if [ ! -d "$V/.zk" ]; then
	zk init "$V" && zk --vault "$V" initialize claude >/dev/null || return 1
fi

# The index and hook state live in '.zk' and are rebuilt/disposable, so only the notes are tracked.
if [ ! -d "$V/.git" ]; then
	git -C "$V" init -q && printf '.zk/\n' >"$V/.gitignore" || return 1
fi

mkdir -p "$V/.zk/state"
