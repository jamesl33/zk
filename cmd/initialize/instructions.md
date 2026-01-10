# A Zettelkasten (Slip Box) Note Vault in Markdown

## General Instructions

- You are well versed in the Zettelkasten, PARA and CODE techniques
- You should use the Zettelkasten to answer questions
- You should do this, by using `zk` to search notes
- Before running a `zk` command you must pass `--help` to determine if the usage is correct
- After finding notes, use tags/links and semantic search to find other related notes
- If you don't find anything using `zk`, you can use normal tools
- You can search/read files (e.g. PDFs)
- Don't fall-back to using the internet

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

## Searching, Interacting and Updating the Zettelkasten

### Listing Notes

```
$ zk notes list
```

1. Notes can be listed to find their directory, title, tags and path
2. The command supports multiple methods of filtering
3. The filtering is based on *only* the note title
4. When multiple filters are provided, it acts as an *AND* operation
5. The output will be in the format `$DIRECTORY\0x1$TITLE\0x1$TAGS0x1$PATH`

You can (and should) still use standard shell utilities to list notes, if required.

### Searching Notes

```
$ zk notes search
```

1. The content of notes can be searched (full text)
2. The command supports multiple methods of filtering
3. The filtering is based on both the note title *and* the content
4. When multiple filters are provided, it acts as an *AND* operation
5. The output will be in the format `$DIRECTORY\0x1$TITLE\0x1$TAGS0x1$PATH`

You can (and should) still use standard shell utilities to search notes, if required.

### Finding Related Notes

One of the core goals of the Zettelkasten is to improve the ability to find related notes, ideas or concepts.

1. You should `list` and/or `search` to find notes on a subject
2. You should use tags to find related notes
3. You can use links to find related notes
4. You should use `find` to find notes that are semantically related to each-other

#### Using Tags

```
$ zk tags list
```

1. You can find all the tags that exist

```
$ zk notes list tagged
```

1. You can list notes that are tagged with or without given tags
2. The output will be in the format `$DIRECTORY\0x1$TITLE\0x1$TAGS0x1$PATH`

#### Using Links

Given notes are atomic, they're often linked using Markdown/Wiki links.

```
$ zk note links
```

1. You can list notes that are linked to or from a given note
2. The output will be in the format `$DIRECTORY\0x1$TITLE\0x1$TAGS0x1$PATH`

#### Using Semantic Search

```
$ zk notes find
```

1. You can find notes that are sematically related to a given note
2. The output will be in the format `$DIRECTORY\0x1$TITLE\0x1$TAGS0x1$PATH`

#### Command Composition

```
$ zk tags list | fzf | xargs -r zk notes list tagged --with | zk notes pick | zk note update -
```

1. `zk` is made to be compostable
2. You can use it within a pipeline to perform more complex searches/operations

### Creating Notes

```
$ zk note create
```

1. The initial note should be created using `zk`
2. It will output the path to the new note
3. You will add the note content/body after
4. You will run `zk tags generate $PATH` once the note is created, to ensure the tags are up-to-date
