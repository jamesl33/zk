# zk

A composable CLI for a Markdown Zettelkasten, with semantic search, an LSP server and an MCP server. Commands print paths and metadata so they pipe into `fzf`, your editor and shell functions.

- Notes get unique names, so links survive renames and moves.
- Search by fixed string, glob, regex or meaning (local embeddings).
- An MCP server, skills and hooks let Claude Code read and write the vault, and give it persistent memory ([zk-memory](#zk-memory)).

# Install

```
go install github.com/jamesl33/zk@latest
zk init ~/notes
```

`zk init [directory]` creates `.zk` and the PARA layout, leaving anything that exists untouched. Run commands from inside a vault, or point at one with `--vault <directory>` (env `ZK_VAULT`).

```
.
├── .zk            # Index and cache (disposable, gitignore it)
├── 0 Inbox        # Flat, fleeting notes
├── 1 Projects     # One directory per project
├── 2 Areas        # One directory per area
├── 3 Resources    # One directory per topic
├── 4 Archives     # Finished or inactive
└── 5 Bibliography # Notes linking to books/articles
```

# Commands

| Command | Does |
|---|---|
| `zk note create <type>` | New `fleeting`, `permanent`, `literature`, `index` or `bibliographic` note; prints its path |
| `zk note update <path>` | Open in `$EDITOR` (`-` reads paths from stdin) |
| `zk note find <path>` | Notes semantically similar to a note |
| `zk note links <path> [--to\|--from]` | Notes linked to/from a note |
| `zk note summarize <path>` | Summarize with Gemini |
| `zk note delete <path>` | Delete a note |
| `zk notes list [dir] [--fixed\|--glob\|--regex]` | List notes by title; `list tagged --with/--without` filters by tag |
| `zk notes search [dir] [--fixed\|--glob\|--regex]` | Search titles and content |
| `zk notes find <query> [--limit] [--max-distance]` | Semantic search from plain text |
| `zk notes pick` | Pick a note from `zk` output with `fzf` |
| `zk tags list\|generate\|delete` | List tags, generate them with Gemini, remove one everywhere |
| `zk links rewrite` | Rewrite links to use note titles |
| `zk lint [dir] [--archives]` | Lint notes |
| `zk lsp` | LSP server |
| `zk mcp` | MCP server over stdio |

## Shell

Every command prints paths, so the pattern is list, pick, open (fish):

```fish
function zkfn -d "Finds notes related to a given note, then picks and opens one"
    set -l candidates (zk note find $argv)
    test -n "$candidates"; or return 1
    set -l note (printf '%s\n' $candidates | zk notes pick)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zkfo -d "Find notes based on a semantic query then picks and opens one"
    set -l candidates (zk notes find $argv)
    test -n "$candidates"; or return 1
    set -l note (printf '%s\n' $candidates | zk notes pick)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zkgt -d "Generates tags for the given directory/note"
    zk tags generate $argv
end

function zklo -d "Lists notes, then picks and opens one"
    set -l notes (zk notes list $argv)
    test -n "$notes"; or return 1
    set -l note (printf '%s\n' $notes | zk notes pick)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zknb -d "Creates a new bibliographic note"
    set -l note (zk note create bibliographic $argv)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zknf -d "Creates a new fleeting note"
    set -l note (zk note create fleeting $argv)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zkni -d "Creates a new index note, fzf-picked directory (ctrl-n: type new directory)"
    set -l result (fd -t d . "1 Projects" "2 Areas" "3 Resources" | fzf --print-query --expect=ctrl-n)
    set -l dir $result[3]
    test "$result[2]" = ctrl-n; and set dir $result[1]
    test -n "$dir"; or return 1
    set -l note (zk note create index $dir $argv)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zknl -d "Creates a new literature note, fzf-picked directory (ctrl-n: type new directory)"
    set -l result (fd -t d . "1 Projects" "2 Areas" "3 Resources" | fzf --print-query --expect=ctrl-n)
    set -l dir $result[3]
    test "$result[2]" = ctrl-n; and set dir $result[1]
    test -n "$dir"; or return 1
    set -l note (zk note create literature $dir $argv)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zknp -d "Creates a new permanent note, fzf-picked directory (ctrl-n: type new directory)"
    set -l result (fd -t d . "1 Projects" "2 Areas" "3 Resources" | fzf --print-query --expect=ctrl-n)
    set -l dir $result[3]
    test "$result[2]" = ctrl-n; and set dir $result[1]
    test -n "$dir"; or return 1
    set -l note (zk note create permanent $dir $argv)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zkoa -d "Consumes the 'zk' note listing output then opens all the notes in the default editor"
    set -l paths

    while read path
        set paths $paths (echo $path | awk -F '\x01' '{ print $NF }')
    end

	$EDITOR $paths
end

function zkp -d "List notes, picks one then prints the path"
    set -l notes (zk notes list $argv)
    test -n "$notes"; or return 1
    printf '%s\n' $notes | zk notes pick
end

function zkqc -d "Quick-capture text as a fleeting note, no editor"
    set -l path (zk note create fleeting --title "$argv[1]")
    test -n "$path"; or return 1
    truncate -s -1 $path
    printf "%s\n" "$argv[2]" >> $path
end

function zkso -d "Search notes, picks one then opens it"
    set -l candidates (zk notes search $argv)
    test -n "$candidates"; or return 1
    set -l note (printf '%s\n' $candidates | zk notes pick)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end

function zkt -d "Lists all tags, picks one, finds notes that have the tag, picks one and updates it"
    set -l tags (zk tags list $argv)
    test -n "$tags"; or return 1
    set -l tag (printf '%s\n' $tags | fzf)
    test -n "$tag"; or return 1
    set -l notes (zk notes list tagged --with $tag)
    test -n "$notes"; or return 1
    set -l note (printf '%s\n' $notes | zk notes pick)
    test -n "$note"; or return 1
    printf '%s\n' $note | zk note update -
end
```

## NeoVim

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

-- Pick a PARA directory via 'fzf', then create the given note type inside it.
local zk_note_create_in_dir = function(type)
	require 'fzf-lua'.fzf_exec("fd -t d . '1 Projects' '2 Areas' '3 Resources'", {
		actions = {
			['default'] = function(selected)
				zk_file_edit(vim.fn.system { 'zk', 'note', 'create', type, selected[1] })
			end,
		},
	})
end

-- Create a new index note.
vim.keymap.set('n', '<leader>zkni', function() zk_note_create_in_dir('index') end)

-- Create a new literature note.
vim.keymap.set('n', '<leader>zknl', function() zk_note_create_in_dir('literature') end)

-- Create a new permanent note.
vim.keymap.set('n', '<leader>zknp', function() zk_note_create_in_dir('permanent') end)

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

| Feature | Backend | Setup |
|---|---|---|
| Embeddings: `notes find`, `note find`, `semantic_search_notes`, `find_related_notes` | Local [Ollama](https://ollama.com) | `ollama serve` and `ollama pull embeddinggemma`; `OLLAMA_HOST` for a remote instance |
| Generation: `note summarize`, `tags generate` | Gemini | `GEMINI_API_KEY` |

Embeddings run locally because they are called far more often than generation. Prompts, responses and embeddings are cached in `.zk/zk.sqlite3`, keyed on a checksum of the input. The cache does not key on model, so set `ZK_GEMINI_MODEL` or `ZK_OLLAMA_EMBED_MODEL` before switching, or stale results survive.

## Claude Code

Install the plugin from within Claude Code:

```
/plugin marketplace add jamesl33/zk
/plugin install zk@zk
```

It ships the skills (archive a project, turn fleeting notes into permanent ones, and so on) and the `zk mcp` server.

- **MCP tools:** `list_notes`, `lint_notes`, `regex_search_notes`, `semantic_search_notes`, `find_related_notes`, `find_notes_linked_to`, `find_notes_linked_from`, `read_note`, `create_note`, `update_note`, `delete_note`.
- **Per vault:** `zk initialize claude` writes `ZK.md` (the Zettelkasten/PARA conventions) and puts `@ZK.md` first in `CLAUDE.md`. It never touches `.claude/`. Rerun it to refresh `ZK.md`.
- **Vault `.gitignore`:** commit `CLAUDE.md`, ignore the generated files.

  ```
  .zk/
  ZK.md
  .claude/settings.local.json
  ```

Install the plugin globally and `zk initialize claude` per vault. Never export `ZK_VAULT` globally.

# zk-memory

A second plugin that replaces Claude Code's built-in auto-memory with a zk vault. Memory becomes notes you can search, link, lint, edit and diff in git.

```
/plugin install zk-memory@zk
```

Needs `zk` on `PATH` and Ollama running. Set `"autoMemoryEnabled": false` in `~/.claude/settings.json`, because a plugin cannot.

## How it works

The vault is `~/.claude/zk-memory` (override with `ZK_MEMORY_VAULT`), a git repo created on first use.

| Hook | Does |
|---|---|
| `SessionStart` (also after compaction) | Injects save guidance and a digest of note titles: this project's notes, plus `user` and `feedback` notes from any project |
| `UserPromptSubmit` | Injects notes semantically close to the prompt, once per session each |
| `PreCompact`, `SessionEnd` | Start a detached headless Claude that saves what the session missed, then lints and commits the vault |

- **Scope:** project notes live in `1 Projects/<repo>` (the repo name, shared by every worktree). Notes useful anywhere go in `3 Resources` or `2 Areas`.
- **Recall:** notes within cosine distance `0.65` of the prompt, up to 15 notes and 64 KB. Prompts under three words and `/` commands are skipped.
- **Saving:** only the transcript since the last save is read, and sessions where the user said nothing are skipped. The model is told to save only material learnings (preferences, corrections, non-obvious gotchas, decisions with reasons) and to update rather than duplicate.
- **Failure:** hooks fail open and never block Claude. Errors go to `.zk/state/hooks.log`. A failed save keeps its place and retries next time.

| Env var | Default | Sets |
|---|---|---|
| `ZK_MEMORY_VAULT` | `~/.claude/zk-memory` | Vault location |
| `ZK_MEMORY_MODEL` | `haiku` | Model for the save run |
| `ZK_MEMORY_MAX_DISTANCE` | `0.65` | Recall cutoff |

# Development

```
make build    # build/zk
make test     # PACKAGE=cmd/note TESTS=TestFoo to narrow
```

# Why

I tried many note tools and none stood the test of time. `zk` favours small composable commands over a single app. It was not designed for performance and has been fast enough for my vault.

# TODO

- [ ] Add support for Trigram Indexing (https://swtch.com/~rsc/regexp/regexp4.html)
- [ ] Make it easier to define/share shell functions/aliases
- [ ] Write a NeoVim plugin to make it easier to define/share that setup
- [ ] Add expiration to the AI prompt/embedding cache
- [ ] Cache the note name index in the LSP server, as each diagnostics/completion request walks the whole vault (~35µs per note, so ~180ms at 5,000 notes)
- [ ] Ignore links in frontmatter and code (blocks/spans) when linting and publishing diagnostics
