---
name: maintainer
description: "A skill to guide the process of converting fleeting and literature notes into permanent notes, and reviewing Projects/Areas/Resources for items that should be archived, including a linting and fixing step."
---

# Maintainer Skill

This skill outlines the process for converting `fleeting` notes into `permanent` notes within the Zettelkasten. `literature` notes are produced by the `librarian` skill and are references, not raw material — leave them as-is unless an idea inside one is worth distilling into a `permanent` note.

## Workflow

1.  **Review the Queue**: Regularly go through your `fleeting` notes in the `0 Inbox` directory and any unprocessed `literature` notes. For each, consider if the idea is still interesting or relevant. If not, it should be deleted (a `fleeting` note) or left as-is (a `literature` note stays as the record of the source).

2.  **Synthesize and Refine**: If the idea is valuable, the next step is to process it.
    *   **Rewrite**: Rephrase the note in your own words. This is crucial for ensuring you've understood the concept. The new note should be self-contained and understandable without any external context.
    *   **Atomize**: Ensure the note is "atomic"—meaning it focuses on a single idea. If a fleeting note contains multiple distinct ideas, break it down into several new `permanent` notes.

3.  **Connect to the Network**: Think about how this new, atomic idea fits within your existing knowledge. A `permanent` note with no links is rejected by `lint_notes` (`orphan-note`), so this step is mandatory, not optional.
    *   Search your vault for related notes using `regex_search_notes` or `semantic_search_notes`.
    *   If the topic already has an `index` note, link to it — that's its purpose. If the topic is new and substantial enough to gather multiple notes over time, consider creating one (`create_note`, type `index`).
    *   Add links from your new note to existing ones (or to the topic's `index` note, if nothing more specific fits yet). Links don't need to be bidirectional — a link from the new note is enough.

4.  **File and Format**:
    *   Use `create_note` (type `permanent`, with a `title`, `tags`, and the rewritten, atomic body including its links) to write the new note into the appropriate location within `1 Projects`, `2 Areas`, or `3 Resources` — it generates the timestamp ID automatically.
    *   If a fleeting or literature note splits into an existing note instead of a new one, use `update_note` to append to it directly.

5.  **Archive the Original**: Once a `fleeting` note has been fully processed into one or more `permanent` notes, delete the original file from the `0 Inbox` using your shell tool to keep it clean. There is no MCP tool for deletion. Processed `literature` notes can stay where they are — they remain the citable record of the source.

6.  **Lint and Fix**: Run `lint_notes` to check for any issues, such as broken links, and fix any errors that are found.

7.  **Review Projects, Areas, and Resources**: Check `1 Projects`, `2 Areas`, and `3 Resources` for
    directories that no longer belong there:
    *   A project whose goal is met or abandoned.
    *   An area of responsibility you no longer maintain.
    *   A resource you're no longer interested in.

    Confirm with the user before handing any of these off to the `archiver` skill — don't archive
    on your own judgment alone.

## Rules

1. Don't touch the `4 Archives` directory, consider it to be read-only
2. Don't remove any `.gitkeep` files
