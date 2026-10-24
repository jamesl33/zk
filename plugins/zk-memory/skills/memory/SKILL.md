---
name: memory
description: Save, update or forget a memory in Claude's persistent zk vault. Use when the user says "remember", states a preference, corrects you, or a recalled memory is wrong.
---

Follow the guidance in the `<memory-digest>` given at session start. In short:

- Search first (`semantic_search_notes`), then `update_note` or `delete_note` rather than adding a contradicting note.
- Project facts go in `1 Projects/<project key>`, cross-project facts in `3 Resources`, ongoing concerns in `2 Areas`. Tag one of `user`, `feedback`, `project`, `reference`.
- One atomic note per fact, titled as the claim, linked to another note or an `index` note.
- Run `lint_notes` afterwards and fix every finding.
