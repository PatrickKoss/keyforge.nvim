# Design: Expanded Challenge Library

## Context

The keyforge.nvim plugin gamifies vim learning through tower defense mechanics. Towers fire at enemies, triggering challenges that the user must complete in Neovim. The current library has ~37 challenges, but the user's nvim config has 60+ custom keymaps that should be practiced.

**Stakeholders**: Plugin users learning vim, maintainers adding challenges
**Constraints**:
- Challenges defined in YAML for easy authoring
- Validation happens in Lua (Neovim side)
- Game engine (Go) selects challenges by category/difficulty
- Must support dynamic hints for custom keymaps

## Goals / Non-Goals

**Goals:**
- 150 diverse challenges covering user's full nvim configuration
- Mix of quick (2-keystroke) and complex (40+ keystroke) challenges
- Comprehensive test coverage for all validation scenarios
- Categories mapped to tower types for gameplay variety

**Non-Goals:**
- Custom challenge editor UI (out of scope)
- Per-user challenge persistence (future work)
- Challenge difficulty auto-adjustment (future work)

## Decisions

### Decision 1: Category Taxonomy

Map 15 categories to 3 tower types for gameplay balance:

| Tower | Categories |
|-------|------------|
| Arrow (movement) | movement, buffer-management, window-management, quickfix, folding |
| LSP (lsp-navigation) | lsp-navigation, telescope, diagnostics, formatting, harpoon |
| Refactor (text-objects) | text-objects, search-replace, refactoring, surround, git-operations |

**Rationale**: Groups related skills while ensuring each tower has diverse challenge pools.

### Decision 2: Challenge Duration Tiers

Introduce `duration_tier` field to challenges:

| Tier | Keystrokes | Par Time | Use Case |
|------|------------|----------|----------|
| quick | 1-5 | 5s | Frequent tower fires, momentum |
| standard | 6-15 | 15s | Core gameplay loop |
| complex | 16-40 | 45s | Boss waves, upgraded towers |
| expert | 40+ | 90s | Special events, achievements |

**Rationale**: Allows wave/difficulty system to select appropriate challenge lengths.

### Decision 3: Runtime Keymap Resolution

Instead of hardcoding keybindings, challenges specify an `hint_action` that gets resolved at runtime:

**Challenge YAML structure:**
```yaml
- id: telescope_find_files
  hint_action: "find_files"           # Action identifier
  hint_fallback: "Use your file finder" # If resolution fails
  required_plugin: "telescope"
```

**Keymap resolution strategy (in Lua):**
```lua
-- Action to search patterns mapping
local action_patterns = {
  find_files = {
    rhs = { "telescope.builtin.find_files", "Telescope find_files", "fzf#" },
    desc = { "find file", "find files", "file finder" },
  },
  format_buffer = {
    rhs = { "conform", "format", "lsp.*format" },
    desc = { "format", "formatting" },
  },
  goto_definition = {
    rhs = { "vim.lsp.buf.definition", "lsp.*definition" },
    desc = { "definition", "go to def" },
  },
  -- ... more actions
}

function M.resolve_keymap(action)
  local patterns = action_patterns[action]
  if not patterns then return nil end

  -- Search all normal mode keymaps
  for _, keymap in ipairs(vim.api.nvim_get_keymap("n")) do
    -- Check rhs patterns
    if keymap.rhs then
      for _, pattern in ipairs(patterns.rhs) do
        if keymap.rhs:lower():find(pattern:lower()) then
          return keymap.lhs, keymap.desc
        end
      end
    end
    -- Check desc patterns
    if keymap.desc then
      for _, pattern in ipairs(patterns.desc) do
        if keymap.desc:lower():find(pattern:lower()) then
          return keymap.lhs, keymap.desc
        end
      end
    end
  end

  return nil -- Fallback to hint_fallback
end
```

**Rationale**: Makes the game portable across any Neovim configuration. Users see their own keybindings, not hardcoded defaults.

### Decision 4: Plugin Availability Detection

Filter challenges based on installed plugins at session start:

**Plugin detection (in Lua):**
```lua
local plugin_cache = {}

-- Plugin name variations to try
local plugin_aliases = {
  ["telescope"] = { "telescope", "telescope.nvim" },
  ["nvim-surround"] = { "nvim-surround", "mini.surround", "surround" },
  ["harpoon"] = { "harpoon" },
  ["gitsigns"] = { "gitsigns", "gitsigns.nvim" },
  ["nvim-tree"] = { "nvim-tree", "nvim-tree.lua", "neo-tree" },
}

function M.plugin_available(name)
  -- Check cache first
  if plugin_cache[name] ~= nil then
    return plugin_cache[name]
  end

  -- Get aliases to try
  local aliases = plugin_aliases[name] or { name, name:gsub("-", "_"), name:gsub("%.nvim$", "") }

  for _, alias in ipairs(aliases) do
    -- Try requiring the module
    local ok = pcall(require, alias)
    if ok then
      plugin_cache[name] = true
      return true
    end
  end

  -- Try checking if command exists (for vimscript plugins)
  local cmd_name = name:gsub("%-", ""):gsub("^%l", string.upper)
  if vim.fn.exists(":" .. cmd_name) == 2 then
    plugin_cache[name] = true
    return true
  end

  plugin_cache[name] = false
  return false
end

function M.filter_challenges_by_plugins(challenges)
  local available = {}
  for _, challenge in ipairs(challenges) do
    if not challenge.required_plugin or M.plugin_available(challenge.required_plugin) then
      table.insert(available, challenge)
    end
  end
  return available
end
```

**Rationale**: Users without certain plugins (e.g., no harpoon) won't get challenges they can't complete. Caching prevents repeated pcall overhead.

### Decision 5: Validation Test Strategy

For 150 challenges, test coverage follows this pattern:
- **Unit tests per validation type**: 7 types x 5 scenarios = 35 tests
- **Challenge-specific tests**: Sample 30 challenges (20%) with win/fail cases = 60 tests
- **Edge case tests**: Unicode, empty, multiline, timeout = 20 tests
- **Integration tests**: End-to-end challenge flow = 10 tests

Total: ~125 new tests

**Rationale**: Full coverage of every challenge would create maintenance burden. Sample testing with comprehensive validation type coverage provides confidence.

## Challenge Catalog (150 Challenges)

### Movement (20 challenges)

**Quick (10)**
1. `movement_end_of_line` - Jump to end with `$`
2. `movement_start_of_line` - Jump to first non-blank with `^`
3. `movement_line_start` - Jump to column 0 with `0`
4. `movement_word_forward` - Move forward one word with `w`
5. `movement_word_backward` - Move backward one word with `b`
6. `movement_word_end` - Move to end of word with `e`
7. `movement_find_char` - Find character with `f{char}`
8. `movement_till_char` - Till character with `t{char}`
9. `movement_find_char_back` - Find backward with `F{char}`
10. `movement_till_char_back` - Till backward with `T{char}`

**Standard (7)**
11. `movement_word_forward_5` - Move forward 5 words with `5w`
12. `movement_goto_line` - Go to line N with `NG` or `:N`
13. `movement_paragraph_forward` - Next paragraph with `}`
14. `movement_paragraph_backward` - Previous paragraph with `{`
15. `movement_matching_bracket` - Match bracket with `%`
16. `movement_search_forward` - Search with `/pattern`
17. `movement_search_backward` - Search backward with `?pattern`

**Complex (3)**
18. `movement_mark_jump` - Set mark and return with `ma` and `'a`
19. `movement_global_mark` - Global mark across files
20. `movement_jump_list` - Navigate jump list with `<C-o>` and `<C-i>`

### Text Objects (18 challenges)

**Quick (6)**
21. `text_delete_word` - Delete word with `dw`
22. `text_change_word` - Change word with `cw`
23. `text_delete_line` - Delete line with `dd`
24. `text_yank_line` - Yank line with `yy`
25. `text_delete_to_end` - Delete to end with `D`
26. `text_change_to_end` - Change to end with `C`

**Standard (8)**
27. `text_change_inner_word` - Change inner word with `ciw`
28. `text_delete_around_word` - Delete around word with `daw`
29. `text_change_inner_quotes` - Change inside quotes with `ci"`
30. `text_delete_inner_parens` - Delete inside parens with `di(`
31. `text_change_inner_brackets` - Change inside brackets with `ci[`
32. `text_delete_inner_braces` - Delete inside braces with `di{`
33. `text_yank_inner_paragraph` - Yank paragraph with `yip`
34. `text_change_inner_tag` - Change inside HTML tag with `cit`

**Complex (4)**
35. `text_change_around_quotes` - Change around quotes with `ca"`
36. `text_visual_inner_function` - Select function body with `vif`
37. `text_delete_around_parens_nested` - Delete nested parens with `da(`
38. `text_change_sentence` - Change sentence with `cis`

### LSP Navigation (15 challenges)

**Quick (5)**
39. `lsp_hover` - Show hover info with `K`
40. `lsp_signature_help` - Signature help with `<C-k>` (insert)
41. `lsp_goto_definition` - Go to definition with `gd`
42. `lsp_goto_declaration` - Go to declaration with `gD`
43. `lsp_goto_implementation` - Go to implementation with `gi`

**Standard (7)**
44. `lsp_find_references` - Find references with `gr`
45. `lsp_type_definition` - Go to type definition with `<leader>D`
46. `lsp_code_action` - Trigger code action with `<leader>ca`
47. `lsp_rename_symbol` - Rename with `<leader>ra`
48. `lsp_document_symbols` - Document symbols with `<leader>fs`
49. `lsp_workspace_symbols` - Workspace symbols with `<leader>fS`
50. `lsp_diagnostic_float` - Show diagnostic with `<leader>de`

**Complex (3)**
51. `lsp_go_to_interface` - Navigate to interface definition
52. `lsp_find_implementation_of_interface` - Find implementations
53. `lsp_organize_imports` - Auto-fix imports with code action

### Search and Replace (12 challenges)

**Quick (4)**
54. `search_current_word` - Search word under cursor with `*`
55. `search_current_word_back` - Search backward with `#`
56. `search_next` - Next match with `n`
57. `search_prev` - Previous match with `N`

**Standard (5)**
58. `search_replace_line` - Replace on line with `:s/old/new/`
59. `search_replace_global` - Replace all with `:%s/old/new/g`
60. `search_replace_confirm` - Confirm each with `:%s/old/new/gc`
61. `search_delete_lines` - Delete matching lines with `:g/pattern/d`
62. `search_visual_selection` - Search selected text with `//` (visual)

**Complex (3)**
63. `search_regex_groups` - Use capture groups in replace
64. `search_replace_visual_selection` - Replace in selection with `<leader>r`
65. `search_very_magic` - Use very magic regex with `\v`

### Refactoring (10 challenges)

**Quick (2)**
66. `refactor_join_lines` - Join lines with `J`
67. `refactor_indent_line` - Indent with `>>` or `<<`

**Standard (5)**
68. `refactor_extract_variable` - Extract expression to variable
69. `refactor_inline_variable` - Inline variable
70. `refactor_rename_function` - Rename function everywhere
71. `refactor_add_parameter` - Add function parameter
72. `refactor_change_signature` - Modify function signature

**Complex (3)**
73. `refactor_extract_function` - Extract code to function
74. `refactor_move_function` - Move function to different location
75. `refactor_convert_to_arrow` - Convert function to arrow function

### Git Operations (12 challenges)

**Quick (4)**
76. `git_next_hunk` - Next git hunk with `]h`
77. `git_prev_hunk` - Previous git hunk with `[h`
78. `git_stage_hunk` - Stage hunk with `<leader>hs`
79. `git_reset_hunk` - Reset hunk with `<leader>hr`

**Standard (5)**
80. `git_preview_hunk` - Preview hunk with `<leader>hp`
81. `git_blame_line` - Show blame with `<leader>hb`
82. `git_toggle_blame` - Toggle line blame with `<leader>tb`
83. `git_diff_this` - Diff current file with `<leader>hd`
84. `git_stage_buffer` - Stage entire buffer with `<leader>hS`

**Complex (3)**
85. `git_conflict_choose_ours` - Choose ours with `<leader>co`
86. `git_conflict_choose_theirs` - Choose theirs with `<leader>ct`
87. `git_open_diffview` - Open diffview with `<leader>gd`

### Window Management (10 challenges)

**Quick (5)**
88. `window_go_left` - Go to left window with `<C-h>`
89. `window_go_down` - Go to lower window with `<C-j>`
90. `window_go_up` - Go to upper window with `<C-k>`
91. `window_go_right` - Go to right window with `<C-l>`
92. `window_split_horizontal` - Split horizontal with `<leader>sh`

**Standard (4)**
93. `window_split_vertical` - Split vertical with `<leader>sv`
94. `window_resize_increase_height` - Increase height with `<C-Up>`
95. `window_resize_decrease_width` - Decrease width with `<C-Left>`
96. `window_close` - Close window with `<C-w>c`

**Complex (1)**
97. `window_rotate_layout` - Rotate windows with `<C-w>r`

### Buffer Management (8 challenges)

**Quick (4)**
98. `buffer_next` - Next buffer with `<S-l>` or `:bn`
99. `buffer_prev` - Previous buffer with `<S-h>` or `:bp`
100. `buffer_close` - Close buffer with `<leader>bd`
101. `buffer_close_alt` - Close buffer with `<leader>x`

**Standard (3)**
102. `buffer_close_all_but_current` - Close others with `<leader>ba`
103. `buffer_close_keep_window` - Close buffer keep window with `<leader>bx`
104. `buffer_save` - Save buffer with `<leader>w`

**Complex (1)**
105. `buffer_reopen_closed` - Reopen last closed buffer

### Folding (8 challenges)

**Quick (3)**
106. `fold_toggle` - Toggle fold with `za`
107. `fold_open` - Open fold with `zo`
108. `fold_close` - Close fold with `zc`

**Standard (4)**
109. `fold_open_all` - Open all folds with `zR`
110. `fold_close_all` - Close all folds with `zM`
111. `fold_peek` - Peek fold preview with `zK`
112. `fold_open_recursive` - Open folds recursively with `zO`

**Complex (1)**
113. `fold_close_recursive` - Close folds recursively with `zC`

### Quickfix (8 challenges)

**Quick (3)**
114. `quickfix_next` - Next item with `]q`
115. `quickfix_prev` - Previous item with `[q`
116. `quickfix_close` - Close quickfix with `<leader>qc`

**Standard (4)**
117. `quickfix_open` - Open quickfix with `<leader>qo`
118. `quickfix_first` - First item with `[Q`
119. `quickfix_last` - Last item with `]Q`
120. `location_next` - Next location with `]l`

**Complex (1)**
121. `quickfix_do` - Run command on all items with `:cdo`

### Diagnostics (6 challenges)

**Quick (2)**
122. `diagnostic_next` - Next diagnostic with `]d`
123. `diagnostic_prev` - Previous diagnostic with `[d`

**Standard (3)**
124. `diagnostic_show_float` - Show diagnostic float with `<leader>de`
125. `diagnostic_goto_error` - Go to next error
126. `diagnostic_set_loclist` - Set location list from diagnostics

**Complex (1)**
127. `diagnostic_disable_buffer` - Disable diagnostics for buffer

### Telescope (10 challenges)

**Quick (4)**
128. `telescope_find_files` - Find files with `<leader>ff`
129. `telescope_buffers` - Find buffers with `<leader>fb`
130. `telescope_oldfiles` - Recent files with `<leader>fo`
131. `telescope_resume` - Resume last picker with `<leader>fr`

**Standard (4)**
132. `telescope_live_grep` - Live grep with `<leader>fg`
133. `telescope_keymaps` - Find keymaps with `<leader>fk`
134. `telescope_current_buffer` - Fuzzy find in buffer with `<leader>/`
135. `telescope_help_tags` - Search help with `<leader>fh`

**Complex (2)**
136. `telescope_from_directory` - Find from directory with `<leader>fd`
137. `telescope_lsp_references` - LSP references via telescope

### Surround (8 challenges)

**Quick (3)**
138. `surround_delete_quotes` - Delete surrounding quotes with `sd"`
139. `surround_delete_parens` - Delete surrounding parens with `sd(`
140. `surround_change_quotes` - Change quotes with `sr"'`

**Standard (4)**
141. `surround_add_word_quotes` - Surround word with quotes `saiw"`
142. `surround_add_word_parens` - Surround word with parens `saiw(`
143. `surround_add_visual` - Surround selection in visual mode
144. `surround_change_tag` - Change HTML tag with `srt`

**Complex (1)**
145. `surround_add_function` - Surround with function call

### Harpoon (5 challenges)

**Quick (2)**
146. `harpoon_add_file` - Add file to harpoon with `<leader>a`
147. `harpoon_toggle_menu` - Toggle menu with `<C-e>`

**Standard (2)**
148. `harpoon_goto_1` - Go to file 1 with `<leader>1`
149. `harpoon_goto_2` - Go to file 2 with `<leader>2`

**Complex (1)**
150. `harpoon_reorder` - Reorder files in harpoon menu

### Formatting (5 challenges)

**Quick (2)**
151. `format_buffer` - Format buffer with `<leader>fm`
152. `format_indent_block` - Indent block with `=`

**Standard (2)**
153. `format_auto_fix_imports` - Organize imports via code action
154. `format_retab` - Convert tabs/spaces with `:retab`

**Complex (1)**
155. `format_custom_range` - Format visual selection

*Note: Final count is 155 to allow for adjustments. We'll trim to exactly 150.*

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Challenge YAML becomes unwieldy | Split into multiple files by category |
| Test maintenance overhead | Sample testing strategy, not 1:1 per challenge |
| Plugin dependencies (telescope, harpoon) | `required_plugin` field skips unavailable |
| Dynamic hints fail | Fallback to static hints or vim help |

## Migration Plan

1. Add new challenges to existing `challenges.yaml`
2. Add new categories to category index
3. Update tower-to-category mapping
4. Add unit tests for new validation scenarios
5. Integration test with sample challenges

**Rollback**: Revert to previous `challenges.yaml` version

## Open Questions

1. Should we split `challenges.yaml` into multiple files per category?
2. What's the exact tower-to-category mapping for new categories?
3. Should expert-tier challenges only appear in boss waves?
