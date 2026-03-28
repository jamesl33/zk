---
name: exporter
description: Export a Zettelkasten note as a standalone, self-contained Markdown file. Use this when you need to share a note with someone who doesn't have access to the vault, or when you want to publish a note as a blog post or document.
---

# Exporter

## Overview

This skill guides the transformation of a Zettelkasten note into a "flat" Markdown file. It ensures all internal context, such as linked notes and bibliographic references, is embedded or resolved so the final document is fully understandable on its own.

## Export Workflow

### 1. Analyze the Source Note

- Identify the target note to export.
- Scan for all internal WikiLinks: `[[$NOTE_ID|$TITLE]]`.
- Scan for bibliographic references (usually linked to notes in `5 Bibliography`).

### 2. Resolve Bibliographic References

- For any link to a `bibliographic` note:
  - Fetch the content of the bibliographic note.
  - Replace the internal link with a full citation in the format: `"[Title]" ($URL)`.
  - If it's a `literature` note, ensure the source is clearly credited at the bottom or inline.

### 3. Embed Linked Concepts

- For each internal link to a `permanent` or `literature` note:
  - Decide if the linked concept is essential for understanding the current note.
  - **Essential**: Fetch the content of the linked note and integrate its core idea directly into the text (e.g., "As discussed in [Topic]... [brief summary]").
  - **Supporting**: Replace the WikiLink with a plain text title or a summarized footnote.
  - **Crucial Definition**: If the note relies on a term defined in another note, embed that definition.

### 4. Flatten the Metadata

- Remove Zettelkasten-specific YAML frontmatter like `type`, `ID`, and `date`.
- Retain the `title` as a primary `# H1` header.
- Convert `tags` into a simple comma-separated list or remove them if not appropriate for the output format.

### 5. Review and Refine

- Ensure the narrative flow is natural. Since multiple notes might be merged, transitions between sections may need to be added.
- Verify that all `$NOTE_ID` references are gone.

## Example Transformation

See [references/examples.md](references/examples.md) for a before-and-after comparison of a note export.
