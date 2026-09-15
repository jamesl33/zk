---
name: archiver
description: Archiving completed or inactive projects in the Zettelkasten. Use when a project in '1 Projects' is done or on hold and needs to be moved to '4 Archives' while extracting useful atomic insights into permanent notes.
---

# Archiver

This skill guides the process of archiving completed or inactive project directories in the Zettelkasten. It ensures that valuable knowledge, lessons learned, and insights are not lost in "cold storage," but are instead extracted, formalized into permanent notes, integrated into active areas or resources, and linked.

## Workflow

To archive a project directory, execute the following steps in sequence:

### 1. Identify and Review the Project
*   Locate the inactive or completed project directory within `1 Projects/` (e.g., `1 Projects/MyProject/`).
*   Review all notes, documents, and files in this directory to identify key ideas, technical patterns, solutions, or insights that have long-term value beyond the project itself.

### 2. Extract Atomic Permanent Notes
*   For every valuable, reusable concept or insight identified, create a new `permanent` note using `zk note create permanent` or the `create_note` tool (this generates the ID and frontmatter, but leaves the body empty — write it with `update_note`, see step 3).
*   **Atomize**: Ensure each new note represents exactly one self-contained, atomic idea.
*   **Metadata**:
    - Set the frontmatter `type` to `permanent`.
    - Provide a descriptive and clear `title`.
    - Add relevant, lowercase, snake_case `tags`.
*   **Placement**: Save these permanent notes in the appropriate active directories:
    - `2 Areas/{Area}/` if they relate to an ongoing area of responsibility.
    - `3 Resources/{Resource}/` if they relate to a topic of ongoing interest or research.

### 3. Build Links
*   Search your Zettelkasten using `regex_search_notes` or `semantic_search_notes` for existing related notes. If the topic has an `index` note, link to it as well.
*   Use `update_note` to write each new permanent note's body, including WikiLinks `[[$NOTE_ID|$TITLE]]` to relevant existing notes — a note with none is rejected by `lint_notes` (`orphan-note`).
*   Link to/from any literature, bibliographic, or project notes that remain relevant.

### 4. Archive the Project Directory
*   Move the entire project directory from `1 Projects/{ProjectName}` into `4 Archives/{ProjectName}`.
*   Ensure all project-specific fleeting notes, meeting notes, and attachments stay within the moved directory.
*   *Command:* Use your shell tool to move it (e.g., `mv "1 Projects/MyProject" "4 Archives/"`). There is no MCP tool for moving files.

### 5. Lint and Validate
*   Run the Zettelkasten linter (`lint_notes` or `zk lint`) to check for any broken links or frontmatter issues.
*   Fix any errors to keep the knowledge network clean.

## Rules

1.  **Extract Before Moving**: Always perform the extraction of atomic permanent notes *before* moving the project folder to `4 Archives`. Once inside `4 Archives`, the content should be considered archived and read-only.
2.  **No Dead Ends**: Every newly created permanent note must link to at least one existing active note.
3.  **Preserve Directory Structure**: When moving a project directory to `4 Archives`, preserve its contents and structure exactly. Do not flatten or delete project files.
4.  **Do Not Remove `.gitkeep` Files**: When moving folders, ensure any `.gitkeep` files in empty directories are preserved so Git tracks the structure properly.
