# Claude memory vault

You have a persistent, cross-project memory: a zk vault at `@VAULT@`, reached with the `zk-memory` MCP tools. It is yours. The user's own notes are a separate vault behind the `zk` tools: read-only unless the user asks otherwise.

## Recall

Relevant notes are injected as `<memory>` before prompts. Treat them as background that may be stale: verify anything about code or state before relying on it. Search with `zk-memory` `semantic_search_notes` when you need more.

## Save

Save only durable knowledge useful in other projects or sessions: user preferences and corrections, decisions and why, tool quirks, hard-won gotchas. Never save task state, progress, or anything derivable from code or git history.

- One atomic note per fact; the title states the claim.
- Search first. Update the existing note, or delete it when wrong, instead of adding a contradicting one.
- Directory: `1 Projects` for knowledge specific to one piece of work (e.g. a new feature in a repo), one sub-directory per project; `3 Resources` for cross-project facts; `2 Areas` for ongoing concerns; `0 Inbox` for fleeting notes. Archive a finished project's directory to `4 Archives`.
- When asked to "save what you learned" (end of session or before compaction), save what qualifies, or reply "nothing to save".
