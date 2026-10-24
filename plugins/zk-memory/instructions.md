# Claude memory vault

You have a persistent, cross-project memory: a zk vault at `@VAULT@`, reached with the `zk-memory` MCP tools. It is yours. The user's own notes are a separate vault behind the `zk` tools: read-only unless the user asks otherwise.

## Recall

Relevant notes are injected as `<memory>` before prompts. Treat them as background that may be stale: verify anything about code or state before relying on it. Search with `zk-memory` `semantic_search_notes` when you need more.

## Save

Save only cross-project knowledge, useful outside the repo it was learned in: user preferences and corrections, tool quirks, hard-won gotchas. Knowledge specific to one project belongs in the built-in project memory, not here. Never save task state, progress, or anything derivable from code or git history.

- One atomic note per fact; the title states the claim. Split anything holding two facts. Write for a stranger: self-contained, no reliance on this session.
- Every note links to at least one other note with `[[ID|Title]]`, found via `find_related_notes` or `semantic_search_notes`. With no related note, link to an `index` note for the topic, creating it if needed.
- After saving, run `lint_notes` and fix every finding before finishing.
- Search first. Update the existing note, or delete it when wrong, instead of adding a contradicting one.
- Directory: `3 Resources` for facts; `2 Areas` for ongoing concerns; `0 Inbox` for fleeting notes.
- When asked to "save what you learned" (end of session or before compaction), save what qualifies, or reply "nothing to save".
