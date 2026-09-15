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

The `zk` command exposes useful commands to enable finding notes that are required; see [#Aliases/Functions](#aliasesfunctions) for some usage examples.

### Aliases/Functions

The tool is designed to be compostable, allow aliases/functions which improve overall functionality.

```fish
function zkfn -d "Finds notes related to a given note, then picks and opens one"
    zk note find $argv | zk notes pick | zk note update -
end

function zkfo -d "Find notes based on a semantic query then picks and opens one"
    zk notes find $argv | zk notes pick | zk note update -
end

function zkgt -d "Generates tags for the given directory/note"
    zk tags generate $argv
end

function zklo -d "Lists notes, then picks and opens one"
    zk notes list $argv | zk notes pick | zk note update -
end

function zknb -d "Creates a new bibliographic note"
    zk note create bibliographic $argv | zk note update -
end

function zknf -d "Creates a new fleeting note"
    zk note create fleeting $argv | zk note update -
end

function zkni -d "Creates a new index note"
    zk note create index $argv | zk note update -
end

function zknl -d "Creates a new literature note"
    zk note create literature $argv | zk note update -
end

function zknp -d "Creates a new permanent note"
    zk note create permanent $argv | zk note update -
end

function zkp -d "List notes, picks one then prints the path"
    zk notes list $argv | zk notes pick
end

function zkso -d "Search notes, picks one then opens it"
    zk notes search $argv | zk notes pick | zk note update -
end

function zkt -d "Lists all tags, picks one, finds notes that have the tag, picks one and updates it"
    zk tags list $argv | fzf | xargs -r zk notes list tagged --with | zk notes pick | zk note update -
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
local zk_fzf_opts = {
	["--ansi"] = true,
	["--with-nth"] = "{1} {2} {3}",
	["--delimiter"] = "\x01",
	["--preview"] = "bat --color=always --style=numbers {4}",
}

-- Common 'fzf' actions for 'zk.'
local zk_fzf_actions = { ["enter"] = zk_fzf_file_edit }

-- Create a new bibliographic note.
vim.keymap.set(
	'n',
	'<leader>zknb',
	function()
		zk_file_edit(vim.fn.system { 'zk', 'note', 'create', 'bibliographic' })
	end
)

-- Create a new fleeting note.
vim.keymap.set(
	'n',
	'<leader>zknf',
	function()
		zk_file_edit(vim.fn.system { 'zk', 'note', 'create', 'fleeting' })
	end
)

-- Create a new index note.
vim.keymap.set(
	'n',
	'<leader>zkni',
	function()
		zk_file_edit(vim.fn.system { 'zk', 'note', 'create', 'index', vim.fn.input('Path: ', '', 'dir') })
	end
)

-- Create a new literature note.
vim.keymap.set(
	'n',
	'<leader>zknl',
	function()
		zk_file_edit(vim.fn.system { 'zk', 'note', 'create', 'literature', vim.fn.input('Path: ', '', 'dir') })
	end
)

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

-- Live search notes.
vim.keymap.set(
	'n',
	'<leader>zkso',
	function()
		require 'fzf-lua'.live_grep({ previewer = false, cmd = "zk notes search --regex", hidden = false, fzf_opts = zk_fzf_opts, actions = zk_fzf_actions, git_icons = false, file_icons = false, exec_empty_query = true })
	end
)

-- List similar notes.
vim.keymap.set(
	'n',
	'<leader>zkfn',
	function()
		require 'fzf-lua'.fzf_exec(string.format("zk note find '%s'", vim.fn.expand("%")),
			{ fzf_opts = zk_fzf_opts, actions = zk_fzf_actions })
	end
)

-- List links to/from the current note.
vim.keymap.set(
	'n',
	'<leader>zkla',
	function()
		require 'fzf-lua'.fzf_exec(string.format("zk note links '%s'", vim.fn.expand("%")),
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

# AI

`zk` has some AI-assisted features, backed by the Gemini API (requires a `GEMINI_API_KEY`):

- `zk note summarize <path>` and `zk tags generate [directory | path]` use Gemini to summarize a note or generate its tags. Both cache prompts/responses/embeddings in `.zk/zk.sqlite3`, keyed on a checksum of the input; if you need a different model, set `ZK_GEMINI_MODEL` and/or `ZK_GEMINI_EMBED_MODEL` before switching, since the cache doesn't key on model and won't notice a change on its own.
- `zk mcp` runs `zk` as a standard [MCP](https://modelcontextprotocol.io) server over stdio, so any MCP-speaking client (not just Gemini CLI) can search, read and write your vault. It exposes: `list_notes`, `lint_notes`, `regex_search_notes`, `semantic_search_notes`, `find_related_notes`, `find_notes_linked_to`, `find_notes_linked_from`, `read_note`, `create_note`, `update_note`.
- `zk initialize` bootstraps [Gemini CLI](https://github.com/google-gemini/gemini-cli) for a vault: it writes `GEMINI.md` (the Zettelkasten/PARA conventions the assistant should follow), `.gemini/settings.json` (registers `zk mcp` as an MCP server), `.gemini/policies/` and `.gemini/skills/` (workflows like archiving a project or converting fleeting notes into permanent ones, built on the MCP tools above).

# Performance

For the amount of notes I have, the performance hasn't been a problem although admittedly it wasn't designed with performance as a key consideration; there's room for optimization in the future, if required.

# Design

- Composable commands, each simple, in conjunction allowing more complex behavior
- Improve "discoverability" of notes, via fixed string, glob, regex and similarly search
- Unique naming to allow links to survive renames/moves (e.g. between directories)

# Why?

I've used multiple tools for this process in the past, and struggled to find anything that fit my desired usage in a way that has stood the test of time.

# TODO

- [ ] Add support for Trigram Indexing (https://swtch.com/~rsc/regexp/regexp4.html)
- [ ] Make it easier to define/share shell functions/aliases
- [ ] Write a NeoVim plugin to make it easier to define/share that setup
