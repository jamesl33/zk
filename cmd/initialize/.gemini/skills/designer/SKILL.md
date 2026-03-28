---
name: designer
description: Synthesis of one-pager design documents from Zettelkasten notes. Use when the user requests a design or architectural plan based on existing research, fleeting notes, or literature in the vault.
---

# Zettelkasten Design Document Synthesis

This skill provides a structured methodology for using `zk` MCP tools to discover, map, and synthesize information from a Zettelkasten vault into a cohesive one-pager design document.

## 1. Discovery Methodology

To produce a high-quality design document, you must traverse the note graph to gather problem context, theoretical foundations, and existing technical components.

### Step A: Semantic Exploration

Start by performing a semantic search to identify the "seed" notes for the design topic.

-   Use `mcp_zk_semantic_search_notes` with a broad description of the design problem.
-   Look for `permanent` notes (your ideas) and `fleeting` notes (captured problems).

### Step B: Literature Anchoring

Identify the theoretical or external basis for the design.

-   Use `mcp_zk_find_notes_linked_from` on the seed notes to find `literature` notes.
-   Follow links to `bibliographic` notes to confirm the original source material.

### Step C: Contextual Expansion

Find related concepts that might not be directly linked but share semantic relevance.

-   Use `mcp_zk_find_related_notes` to find overlapping architectural patterns or previous designs.
-   Use `mcp_zk_regex_search_notes` for specific technical keywords (e.g., `#architecture`, `#api`, `#trigram`).

## 2. Synthesis & Template

Once you have identified the core notes, use the `assets/template.md` to structure the final document.

### Synthesis Rules:

1.  **Atomic Integration**: Map each section of the design document to specific atomic notes from your research.
2.  **Back-Linking**: Every major claim or architectural decision should ideally have a WikiLink `[[$ID|$TITLE]]` to the source note in the Zettelkasten.
3.  **Progressive Refinement**: If a section is thin (e.g., "Implementation Plan"), perform another targeted search to find actionable sub-tasks.

## 3. Tool Usage Examples

### Finding relevant notes

```bash
# Semantic search for the design topic
mcp_zk_semantic_search_notes query="Trigram indexing for Zettelkasten search"

# Find notes linked TO a key technical note
mcp_zk_find_notes_linked_to path="2 Areas/Linux/20260206094342.md"
```

### Reading and Validating

-   Use `read_file` to extract the content of identified notes.
-   Use `mcp_zk_lint_notes` after creating the design doc to ensure it follows vault standards.

---

-   **Template**: `assets/template.md`
