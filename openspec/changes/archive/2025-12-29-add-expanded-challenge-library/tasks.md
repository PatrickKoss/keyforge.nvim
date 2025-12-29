# Tasks: Expanded Challenge Library

## 1. Challenge YAML Expansion

### 1.1 Movement Challenges (20 total)
- [x] 1.1.1 Add 10 quick movement challenges (end_of_line, start_of_line, 0, w, b, e, f, t, F, T)
- [x] 1.1.2 Add 7 standard movement challenges (5w, NG, }, {, %, /pattern, ?pattern)
- [x] 1.1.3 Add 3 complex movement challenges (marks, global marks, jump list)

### 1.2 Text Object Challenges (18 total)
- [x] 1.2.1 Add 6 quick text object challenges (dw, cw, dd, yy, D, C)
- [x] 1.2.2 Add 8 standard text object challenges (ciw, daw, ci", di(, ci[, di{, yip, cit)
- [x] 1.2.3 Add 4 complex text object challenges (ca", vif, nested da(, cis)

### 1.3 LSP Navigation Challenges (15 total)
- [x] 1.3.1 Add 5 quick LSP challenges (K, C-k, gd, gD, gi)
- [x] 1.3.2 Add 7 standard LSP challenges (gr, leader-D, leader-ca, leader-ra, leader-fs, leader-fS, leader-de)
- [x] 1.3.3 Add 3 complex LSP challenges (go to interface, find implementations, organize imports)

### 1.4 Search and Replace Challenges (12 total)
- [x] 1.4.1 Add 4 quick search challenges (*, #, n, N)
- [x] 1.4.2 Add 5 standard search challenges (:s, :%s/g, :%s/gc, :g/d, // visual)
- [x] 1.4.3 Add 3 complex search challenges (regex groups, visual replace, very magic)

### 1.5 Refactoring Challenges (10 total)
- [x] 1.5.1 Add 2 quick refactoring challenges (J, >>)
- [x] 1.5.2 Add 5 standard refactoring challenges (extract variable, inline, rename, add param, change signature)
- [x] 1.5.3 Add 3 complex refactoring challenges (extract function, move function, convert to arrow)

### 1.6 Git Operations Challenges (12 total)
- [x] 1.6.1 Add 4 quick git challenges (]h, [h, leader-hs, leader-hr)
- [x] 1.6.2 Add 5 standard git challenges (leader-hp, leader-hb, leader-tb, leader-hd, leader-hS)
- [x] 1.6.3 Add 3 complex git challenges (conflict ours, conflict theirs, diffview)

### 1.7 Window Management Challenges (10 total)
- [x] 1.7.1 Add 5 quick window challenges (C-h, C-j, C-k, C-l, leader-sh)
- [x] 1.7.2 Add 4 standard window challenges (leader-sv, C-Up, C-Left, C-w c)
- [x] 1.7.3 Add 1 complex window challenge (rotate layout)

### 1.8 Buffer Management Challenges (8 total)
- [x] 1.8.1 Add 4 quick buffer challenges (S-l, S-h, leader-bd, leader-x)
- [x] 1.8.2 Add 3 standard buffer challenges (leader-ba, leader-bx, leader-w)
- [x] 1.8.3 Add 1 complex buffer challenge (reopen closed)

### 1.9 Folding Challenges (8 total)
- [x] 1.9.1 Add 3 quick folding challenges (za, zo, zc)
- [x] 1.9.2 Add 4 standard folding challenges (zR, zM, zK, zO)
- [x] 1.9.3 Add 1 complex folding challenge (zC recursive)

### 1.10 Quickfix Challenges (8 total)
- [x] 1.10.1 Add 3 quick quickfix challenges (]q, [q, leader-qc)
- [x] 1.10.2 Add 4 standard quickfix challenges (leader-qo, [Q, ]Q, ]l)
- [x] 1.10.3 Add 1 complex quickfix challenge (:cdo)

### 1.11 Diagnostics Challenges (6 total)
- [x] 1.11.1 Add 2 quick diagnostic challenges (]d, [d)
- [x] 1.11.2 Add 3 standard diagnostic challenges (leader-de, goto error, set loclist)
- [x] 1.11.3 Add 1 complex diagnostic challenge (disable buffer)

### 1.12 Telescope Challenges (10 total)
- [x] 1.12.1 Add 4 quick telescope challenges (leader-ff, leader-fb, leader-fo, leader-fr)
- [x] 1.12.2 Add 4 standard telescope challenges (leader-fg, leader-fk, leader-/, leader-fh)
- [x] 1.12.3 Add 2 complex telescope challenges (leader-fd, LSP references)

### 1.13 Surround Challenges (8 total)
- [x] 1.13.1 Add 3 quick surround challenges (sd", sd(, sr"')
- [x] 1.13.2 Add 4 standard surround challenges (saiw", saiw(, visual surround, srt)
- [x] 1.13.3 Add 1 complex surround challenge (function wrap)

### 1.14 Harpoon Challenges (5 total)
- [x] 1.14.1 Add 2 quick harpoon challenges (leader-a, C-e)
- [x] 1.14.2 Add 2 standard harpoon challenges (leader-1, leader-2)
- [x] 1.14.3 Add 1 complex harpoon challenge (reorder)

### 1.15 Formatting Challenges (5 total)
- [x] 1.15.1 Add 2 quick formatting challenges (leader-fm, =)
- [x] 1.15.2 Add 2 standard formatting challenges (organize imports, :retab)
- [x] 1.15.3 Add 1 complex formatting challenge (visual range format)

## 2. Game Engine Updates

### 2.1 Category System
- [x] 2.1.1 Add new categories to ChallengeManager index
- [x] 2.1.2 Update tower-to-category mapping in tower types
- [x] 2.1.3 Add duration_tier derivation from par_keystrokes
- [x] 2.1.4 Add par_time calculation based on duration_tier

### 2.2 Challenge Selection
- [x] 2.2.1 Update GetRandomChallenge to support new categories
- [x] 2.2.2 Add duration_tier filtering option
- [x] 2.2.3 Verify plugin requirements before selection

## 3. Lua Plugin Updates

### 3.1 Runtime Keymap Resolution
- [x] 3.1.1 Create action_patterns table mapping action IDs to rhs/desc search patterns
- [x] 3.1.2 Implement resolve_keymap(action) function using vim.api.nvim_get_keymap()
- [x] 3.1.3 Add fallback to hint_fallback when keymap resolution fails
- [x] 3.1.4 Handle standard vim commands (no resolution needed for $, ciw, dd, etc.)
- [x] 3.1.5 Format resolved keymaps for display (e.g., "<leader>ff" -> "Space f f")

### 3.2 Plugin Availability Detection
- [x] 3.2.1 Create plugin_aliases table with common name variations
- [x] 3.2.2 Implement plugin_available(name) function using pcall(require, ...)
- [x] 3.2.3 Add command existence check fallback (vim.fn.exists)
- [x] 3.2.4 Implement session-level caching for plugin detection results
- [x] 3.2.5 Create filter_challenges_by_plugins(challenges) function
- [x] 3.2.6 Integrate plugin filtering into challenge selection flow

### 3.3 Challenge Hint Display
- [x] 3.3.1 Update challenge UI to call resolve_keymap for hint_action fields
- [x] 3.3.2 Display resolved keymap in challenge description
- [x] 3.3.3 Show fallback text when resolution fails

### 3.4 Validation
- [x] 3.4.1 Verify all validation types work with new challenges
- [x] 3.4.2 Add any missing validation edge case handling

## 4. Unit Tests

### 4.1 Go Tests (game/internal/engine/)
- [x] 4.1.1 Test new category indexing
- [x] 4.1.2 Test duration_tier calculation
- [x] 4.1.3 Test tower-to-category mapping
- [x] 4.1.4 Test challenge count validation (>= 150)
- [x] 4.1.5 Test category balance (5-25 per category)

### 4.2 Lua Tests (tests/keyforge/)
- [x] 4.2.1 Add validation type tests (exact_match success/fail/edge)
- [x] 4.2.2 Add validation type tests (contains success/fail/edge)
- [x] 4.2.3 Add validation type tests (cursor_position success/fail/edge)
- [x] 4.2.4 Add validation type tests (function_exists success/fail/edge)
- [x] 4.2.5 Add validation type tests (pattern success/fail/edge)
- [x] 4.2.6 Add validation type tests (different success/fail/edge)
- [x] 4.2.7 Add validation type tests (cursor_on_char success/fail/edge)

### 4.3 Sample Challenge Tests
- [x] 4.3.1 Add win/fail tests for 5 movement challenges
- [x] 4.3.2 Add win/fail tests for 5 text-object challenges
- [x] 4.3.3 Add win/fail tests for 5 LSP challenges
- [x] 4.3.4 Add win/fail tests for 3 search-replace challenges
- [x] 4.3.5 Add win/fail tests for 3 refactoring challenges
- [x] 4.3.6 Add win/fail tests for 3 git-operations challenges
- [x] 4.3.7 Add win/fail tests for 2 window-management challenges
- [x] 4.3.8 Add win/fail tests for 2 buffer-management challenges
- [x] 4.3.9 Add win/fail tests for 2 folding challenges
- [x] 4.3.10 Add win/fail tests for 2 quickfix challenges
- [x] 4.3.11 Add win/fail tests for 2 telescope challenges
- [x] 4.3.12 Add win/fail tests for 2 surround challenges

### 4.4 Keymap Resolution Tests
- [x] 4.4.1 Test resolve_keymap returns user's actual keymap for telescope actions
- [x] 4.4.2 Test resolve_keymap returns nil when keymap not found
- [x] 4.4.3 Test resolve_keymap matches by rhs pattern
- [x] 4.4.4 Test resolve_keymap matches by desc pattern (case-insensitive)
- [x] 4.4.5 Test standard vim commands skip resolution

### 4.5 Plugin Detection Tests
- [x] 4.5.1 Test plugin_available returns true for installed plugins
- [x] 4.5.2 Test plugin_available returns false for missing plugins
- [x] 4.5.3 Test plugin_available tries name variations (nvim-surround -> mini.surround)
- [x] 4.5.4 Test plugin_available caches results
- [x] 4.5.5 Test filter_challenges_by_plugins excludes unavailable plugin challenges
- [x] 4.5.6 Test filter_challenges_by_plugins keeps challenges with no required_plugin

### 4.6 Edge Case Tests
- [x] 4.6.1 Test validation with empty buffer
- [x] 4.6.2 Test validation with unicode content
- [x] 4.6.3 Test validation with multiline (>100 lines) content
- [x] 4.6.4 Test challenge timeout scenarios
- [x] 4.6.5 Test plugin unavailable handling

## 5. Integration Testing
- [x] 5.1 Manual playtest with new challenge pool
- [x] 5.2 Verify tower-category distribution in gameplay
- [x] 5.3 Verify difficulty progression feels balanced
- [x] 5.4 Test keymap resolution with different nvim configs
- [x] 5.5 Verify challenges filtered correctly based on installed plugins
