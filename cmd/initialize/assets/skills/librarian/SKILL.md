---
name: librarian
description: Capture an external source (article, book, paper, video, conversation) into bibliographic and literature notes. Use when the user shares or references something they read or watched and wants it saved to the Zettelkasten.
---

# Librarian Skill

This skill captures external material into the Zettelkasten, keeping what the source says (`literature`) separate from what you think about it (`permanent`). See the vault's instructions file section 4.2 for the exact literature note format (blockquoted quotes, progressive summarization).

## Workflow

1.  **Record the Source**: Use `create_note` (type `bibliographic`, path `5 Bibliography`) with the source's title and a body containing the citation (author, URL/ISBN, date accessed).

2.  **Write the Literature Note**: Use `create_note` (type `literature`, same directory as the bibliographic note it's linked to) with a title describing the note's focus, and a body containing verbatim quotes in blockquotes and your own summaries as plain text alongside them. End the body with a `---` rule, a blank line, then the source link as a bullet: `- [[$NOTE_ID|$TITLE]]` back to the bibliographic note. Never put the link at the top.

3.  **Keep it Raw**: A literature note only records the source plus light commentary. If something sparks an original idea worth developing, don't expand it here — that becomes a new `permanent` note, processed later.

4.  **Hand Off**: Leave the literature note for later processing. The `maintainer` skill treats `literature` notes the same as `fleeting` notes: raw material waiting to be distilled into atomic `permanent` notes.

## Rules

1. Never put your own analysis inside a blockquote — only verbatim quotes belong there.
2. One literature note per source, unless the source is long enough that splitting by section aids navigation.
3. Links to sources always go in the footer after `---`, never in the opening lines. Before creating a note, read an existing literature note (e.g. `3 Resources/AI/20260801120100.md`) to match its layout.
