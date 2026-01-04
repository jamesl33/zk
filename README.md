# zk

`zk` allows interacting, searching and updating a Markdown Zettelkasten; it's designed to allow command composition.

# Usage

I use `zk` to mange my notes, which are organized to follow [PARA](https://fortelabs.com/blog/para).

```
.
├── .zk
├── 0 Inbox        # Flat, fleeting notes
├── 1 Projects     # Directories containing a flat hierarchy of "project" notes
├── 2 Areas        # Directories containing a flat hierarchy "area" notes
├── 3 Resources    # Directories containing a flat hierarchy "resource" notes
├── 4 Archives     # Directories/notes that are not longer required
└── 5 Bibliography # Notes which link to books/articals
```

## CLI

The `zk` command exposes useful commands to enable finding notes that are required; a non-exhaustive list of the
commands/functionality is as follows.

```sh
# Create a new note
zk note create [type] <path>

# List notes that link to another
zk note links --to <path>

# List notes linked from a given note
zk note links --from <path>

# List all notes
zk notes list <path>

# List notes matching a pattern
zk notes list --fixed <string> --glob <pattern> --regex '<pattern>'

# Search the contents of notes
zk notes search --fixed <string> --glob <pattern> --regex '<pattern>'

# Find related notes
zk notes find <path>

# List all tags
zk tags list

# Generate new tags for a note
zk tags generate <path>

# Remove a tag globally
zk tags delete <tag>
```

### Aliases/Functions

The tool is designed to be compostable, allow aliases/functions which improve overall functionality.

```fish
# Create a new note, then open it in an editor
function zknp
    zk note create permanent $argv | zk note update -
end
```

```fish
# List all the notes, pick one in 'fzf', then open it in an editor
function zklo
    zk notes list $argv | zk notes pick | zk note update -
end
```

```fish
# List all the tags, select one, find notes tagged with it, then edit it
function zkt
    zk tags list $argv | fzf | xargs -r zk notes tagged --with | zk notes pick | zk note update -
end
```

## NeoVim

In a similar way that functions/aliases can be used to improve `zk`, the same can be done in NeoVim.

```lua
-- Edit the given file
local zk_file_edit = function(output)
	vim.cmd.edit { output:match("^%s*(.-)%s*$") }
end

-- Extracts the path from the selected line, then opens the file.
local zk_fzf_file_edit = function(selected, opts)
	require "fzf-lua.actions".file_edit({ string.match(selected[1], "[^\x01]+$") }, opts)
end

-- Common 'fzf' options for 'zk.'
local zk_fzf_opts = { ["--ansi"] = true, ["--with-nth"] = "{1} {2} [{3}]", ["--delimiter"] = "\x01" }

-- Common 'fzf' actions for 'zk.'
local zk_fzf_actions = { ["enter"] = zk_fzf_file_edit }

-- Create a new permanent note.
vim.keymap.set(
	'n',
	'<leader>zknp',
	function()
		zk_file_edit(vim.fn.system { 'zk', 'note', 'create', 'permanent', vim.fn.input('Path: ', '', 'dir') })
	end
)

-- List notes.
vim.keymap.set(
	'n',
	'<leader>zklo',
	function()
		require 'fzf-lua'.fzf_exec("zk notes list",
			{ fzf_opts = zk_fzf_opts, actions = zk_fzf_actions })
	end
)

-- List links to the current note.
vim.keymap.set(
	'n',
	'<leader>zklt',
	function()
		require 'fzf-lua'.fzf_exec(string.format("zk note links --to '%s'", vim.fn.expand("%")),
			{ fzf_opts = zk_fzf_opts, actions = zk_fzf_actions })
	end
)

-- List links from the current note.
vim.keymap.set(
	'n',
	'<leader>zklf',
	function()
		require 'fzf-lua'.fzf_exec(string.format("zk note links --from '%s'", vim.fn.expand("%")),
			{ fzf_opts = zk_fzf_opts, actions = zk_fzf_actions })
	end
)

-- Generate tags for the current note.
vim.keymap.set(
	'n',
	'<leader>zkgt',
	function()
		-- Write out the file, generation should have the latest changes
		vim.cmd.write { vim.fn.expand("%") }

		-- Issue generation
		vim.fn.system { 'zk', 'tags', 'generate', vim.fn.expand("%") }

		-- Reload the file, changes occur outside the editor
		vim.cmd.edit {}
	end
)
```

# Performance

For the amount of notes I have, the performance hasn't been a problem although admittedly it wasn't designed with
performance as a key consideration; there's room for optimization in the future, if required.

# Design

- Composable commands, each simple, in conjunction allowing more complex behavior
- Improve "discoverability" of notes, via fixed string, glob, regex and similarly search
- Unique naming to allow links to survive renames/moves (e.g. between directories)

# Why?

I've used multiple tools for this process in the past, and struggled to find anything that fit my desired usage in a way
that has stood the test of time.
