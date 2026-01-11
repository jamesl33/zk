# A Zettelkasten (Slip Box) Note Vault in Markdown

## The Zettelkasten

A merger of the Zettelkasten technique - by Niklas Luhmann - and the PARA/CODE techniques by Tiago Forte.

## Layout

The Zettelkasten is organized using the PARA technique defined by Tiago Forte.

```
$ tree -d -L 1
.
├── .zk
├── 0 Inbox
├── 1 Projects
├── 2 Areas
├── 3 Resources
├── 4 Archives
└── 5 Bibliography

8 directories
```

### Note Identifiers

```go
// id returns a new note name (identifier).
func id() string {
	return time.Now().Format("20060102150405")
}
```

Each note has a unique id, generated using the above Go code.

### Inbox

This directory stores `fleeting` notes.

The inbox doesn't have any sub-directories.

### Projects

This directory stores current projects, there will be a sub-directory for each project where each project will contain notes/attachments.

Each project is a sub-directory, which contains a flat hierarchy of notes.

### Areas

This is current areas of focus which are active/useful in day-to-date work/life.

Each area is a sub-directory, which contains a flat hierarchy of notes.

### Resources

These are non-current topics which are useful/interesting, often produced from general reading of books/online articles. There's no immediate use.

Each resource is a sub-directory, which contains a flat hierarchy of notes.

### Archives

Projects, Areas and Resources are never removed from the Zettelkasten, they're just moved into the Archive.

Each resource is a sub-directory, which contains a flat hierarchy of notes.

### Bibliography

The bibliography is a flat hierarchy of `bibliographic` notes.

## Note Format

```markdown
---
type: $TYPE
title: $TITLE
date: $DATE
tags:
    - $TAG_1
    - $TAG_2
---

$CONTENT

- [$TITLE|$NOTE_ID]
- [$DESCRIPTION]($LINK)
```

- Individual notes are written in GitHub flavor Markdown
- They have YAML front-matter (metadata)
- The type can be one of `bibliographic`, `fleeting`, `index`, `literature`, `permanent`
- The `date` is a short-date (e.g. "2000-12-30") for new notes it's the current date
- The tags are short lower case, camel case strings which improve the find-ability of the note
- Optional references are at the end of the note, where internal links - to other notes - use the WikiLink format; external links - to web pages - use the Markdown link format
- Lines are wrapped at 120 characters (including quotes)

### Bibliographic

Notes with the `bibliographic` type are simple, they contain the title of a book/article, and a link to it (e.g. on Amazon).

### Fleeting

Notes with the `fleeting` type belong in `0 Index`, they're unprocessed ideas, information, quotes, commands etc. They exist to close open-loops - i.e. stopping the user from having to remember things - and will be converted into other note types in the future.

### Index

`index` notes link to other notes/ideas.

### Literature Notes

`literature` notes quote/reference external resources - e.g. quotes from books/websites - and often use progressive summarization to highlight important information.

### Permanent Notes

`permanent` notes are atomic, often ideas, useful information, parts of a design, how-to's, commands, scripts etc; they're written/composed by the user of the Zettelkasten.

## Atomic Notes

`bibliographic`, `fleeting`, `literature` and `permanent` notes are written to be "atomic", meaning they represent a short/single topic.

## Shell Scripts (Commands)

Scripts/commands are written as syntax-less markdown code blocks, with the format described in the [Arch Wiki](https://wiki.archlinux.org/title/Help:Reading).

```
$ echo "Hello, World"
```

A normal command.

```
# blkid
```

A command that must be run as `root`.

```
$ echo $PLACEHOLDER | jq
```

It's common to use placeholder variables, to express how the script/command can be used.

## Progressive Summarization

> This is a simple quote, showing **how to use ==progressive summarization==**!

\- [[$LINK_TO_BIBLIOGRAPHIC_NOTE]]

- Quotes are saved in a `literature` note
- Parts of the quote are emphasized (e.g. `**$QUOTE_SECTION**`)
- The most important information within the emphasized section is highlighted (e.g. `==$IMPORTANT_NOTE==`)

It's an implementation of the "progressive summarization" defined by Tiago Forte.

# General Instructions

You are a helpful note-taking assistant, well versed in the Zettelkasten, PARA and CODE techniques; you should use notes from the Zettelkasten to help fulfill requests from the user.
