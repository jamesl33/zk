---
name: librarian
description: Capture an external source (article, book, paper, video, conversation) into bibliographic and literature notes. Use when the user shares or references something they read or watched and wants it saved to the Zettelkasten.
---

# Librarian Skill

This skill captures external material into the Zettelkasten, keeping what the source says (`literature`) separate from what you think about it (`permanent`). See the vault's instructions file section 4.2 for the exact literature note format (blockquoted quotes, progressive summarization).

## Workflow

1.  **Record the Source**: Use `create_note` (type `bibliographic`, path `5 Bibliography`) with the source's title and a body containing the citation (author, URL/ISBN, date accessed), laid out as in "Note Layouts" below.
2.  **Write the Literature Note**: Use `create_note` (type `literature`, same directory as the bibliographic note it's linked to) with a title describing the note's focus, and a body containing verbatim quotes in blockquotes and your own summaries as plain text alongside them. End the body with a `---` rule, a blank line, then the source link as a bullet: `- [[$NOTE_ID|$TITLE]]` back to the bibliographic note. Never put the link at the top.
3.  **Keep it Raw**: A literature note only records the source plus light commentary. If something sparks an original idea worth developing, don't expand it here — that becomes a new `permanent` note, processed later.
4.  **Hand Off**: Leave the literature note for later processing. The `maintainer` skill treats `literature` notes the same as `fleeting` notes: raw material waiting to be distilled into atomic `permanent` notes.

## Note Layouts

Use these layouts every time, so notes stay consistent. Both follow the vault's `CLAUDE.md` (sections 3 and 4.2).

Bibliographic body: the title and a link only, as one standard Markdown link. No authors, publisher or dates; the frontmatter `date` is the access date.

```
[$TITLE]($URL)
```

For a book, `$URL` is the ISBN-based product or catalogue page.

Literature body, in this order:

1. Opening paragraph: one or two sentences of plain summary.
2. Body: verbatim quotes in blockquotes, progressively summarized per `CLAUDE.md` 4.2 (**bold** key sections, ==highlight== the crucial details). Put your own summaries as plain text beside the quote, never inside it. Group related points as `-` bullets with a short lead-in line ending in a colon.
3. Footer: `---`, a blank line, then `- [[$NOTE_ID|$TITLE]]` to the bibliographic note.

Formatting details:
- One line per paragraph or bullet; never hard-wrap.
- Tags: lowercase snake_case, three to five.

## Rules

1. Never put your own analysis inside a blockquote — only verbatim quotes belong there.
2. One literature note per source, unless the source is long enough that splitting by section aids navigation.
3. Links to sources always go in the footer after `---`, never in the opening lines. Before creating a note, read an existing literature note (e.g. `3 Resources/AI/20260801120100.md`) to match its layout.
