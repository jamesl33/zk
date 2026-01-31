---
name: maintainer
description: "A skill to guide the process of converting fleeting notes into permanent notes, including a linting and fixing step."
---

# Maintainer Skill

This skill outlines the process for converting `fleeting` notes into `permanent` notes within the Zettelkasten.

## Workflow

1.  **Review the Inbox**: Regularly go through your `fleeting` notes in the `0 Inbox` directory. For each note, consider if the idea is still interesting or relevant. If it's no longer valuable, it should be deleted.

2.  **Synthesize and Refine**: If the idea is valuable, the next step is to process it.
    *   **Rewrite**: Rephrase the note in your own words. This is crucial for ensuring you've understood the concept. The new note should be self-contained and understandable without any external context.
    *   **Atomize**: Ensure the note is "atomic"—meaning it focuses on a single idea. If a fleeting note contains multiple distinct ideas, break it down into several new `permanent` notes.

3.  **Connect to the Network**: Think about how this new, atomic idea fits within your existing knowledge.
    *   Search your vault for related notes.
    *   Add links from your new note to existing ones.
    *   Crucially, open the existing notes and add links back to your new note. This bidirectional linking is what builds a web of knowledge.

4.  **File and Format**:
    *   Create a new note with a unique timestamp ID (`YYYYMMDDHHMMSS`).
    *   Set the frontmatter `type` to `permanent`.
    *   Give it a clear, descriptive `title`.
    *   Add relevant `tags` to make it discoverable.
    *   Move the newly created `permanent` note to the appropriate location within `1 Projects`, `2 Areas`, or `3 Resources`.

5.  **Archive the Original**: Once the `fleeting` note has been fully processed into one or more `permanent` notes, delete the original file from the `0 Inbox` to keep it clean.

6.  **Lint and Fix**: Run `lint_notes` to check for any issues, such as broken links, and fix any errors that are found.
