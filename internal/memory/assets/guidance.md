You have persistent memory: a zk vault reached with the `zk-memory` MCP tools (`mcp__plugin_zk-memory_zk-memory__*`). It replaces built-in auto-memory. The user's own notes are a separate vault behind the `zk` tools: read-only unless asked.

Save when the user states a preference, corrects you, says "remember", or you hit a hard-won gotcha. Never save task state, progress, or anything derivable from code or git history.

- Scope: project-specific facts go in `1 Projects/<project>` (the project key is in the digest). Facts useful in any repo go in `3 Resources`; ongoing concerns in `2 Areas`.
- Tag each note one of `user`, `feedback`, `project`, `reference`.
- One atomic note per fact; the title states the claim. Self-contained, no reliance on this session.
- Link every note to another with `[[ID|Title]]` (`find_related_notes`, `semantic_search_notes`); else link an `index` note for the topic, creating it if needed.
- Search first. `update_note` or `delete_note` instead of adding a contradicting note. A recalled note that is wrong: fix or delete it now.
- After saving, run `lint_notes` and fix every finding.
