---
name: researcher
description: Provides a methodology for effectively searching a Zettelkasten, using a combination of regular expression, semantic, and link-based search tools.
---

# Zettelkasten Researcher

This skill provides a methodology for searching your Zettelkasten.

## Searching Methodology

To search your Zettelkasten effectively, use the available tools based on what you're looking for.

### For specific text, titles, or tags:

Use `regex_search_notes`. This is best when you know a specific keyword, phrase, tag, or title.

- **To find a tag:** `regex_search_notes(expression='thru_hiking')`
- **To find a title:** `regex_search_notes(expression='(?i)progressive summarization')`
- **To find content:** `regex_search_notes(expression='a specific phrase in a note')`

### For general ideas or concepts:

Use `semantic_search_notes`. This is best for finding notes related to a topic, even if they don't contain the exact keywords.

- **To find a concept:** `semantic_search_notes(query='what have I written about personal productivity?')`

### For exploring connections:

Once you have a specific note, you can explore its connections:

- `find_notes_linked_from`: See which notes are directly referenced by the current note.
- `find_notes_linked_to`: See which other notes reference the current note.
- `find_related_notes`: Find other notes that are conceptually similar to the current one.

### To browse a directory:

If you know the general area your note is in (e.g., a specific project), use `list_notes` to see all the notes in that directory.

- **To list notes in a project:** `list_notes(path='1 Projects/Bike (2026)/')`
