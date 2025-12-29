# Change: Expand Challenge Library to 150 Challenges

## Why

The current challenge library has ~37 challenges across 6 categories. To truly gamify learning vim keybindings and provide long-term engagement, we need a comprehensive library that:
1. Covers the user's complete nvim configuration (60+ custom keymaps)
2. Includes varied difficulty and duration (quick 2-keystroke challenges to complex multi-step refactoring)
3. Draws inspiration from proven vim training materials (Vim-Katas)
4. Provides extensive test coverage for validation logic

## What Changes

### New Challenge Categories (expanding from 6 to 15)
- **movement**: Basic motions (`$`, `^`, `w`, `f`, `G`, `}`, `%`) - expanded
- **text-objects**: `ciw`, `da(`, `di[`, `cit`, text manipulation - expanded
- **lsp-navigation**: `gd`, `gr`, `K`, go to interface, find implementation - expanded
- **search-replace**: `:s`, `:%s`, `:g`, regex patterns - expanded
- **refactoring**: Extract function, inline variable, rename - expanded
- **git-operations**: Stage hunk, diff, blame, conflict resolution - expanded
- **window-management** (NEW): Splits, tabs, resize, navigate windows
- **buffer-management** (NEW): Buffer switching, close, close all but current
- **folding** (NEW): `zR`, `zM`, `zK`, `za`, `zo`, `zc` with nvim-ufo
- **quickfix** (NEW): Navigate quickfix/location lists, open/close
- **diagnostics** (NEW): Navigate diagnostics, show float, quickfix
- **telescope** (NEW): Find files, live grep, buffers, symbols, resume
- **surround** (NEW): Add/change/delete surrounds with mini.surround
- **harpoon** (NEW): Quick file marks and navigation
- **formatting** (NEW): Format buffer, auto-fix imports

### Challenge Distribution (150 total)
| Category | Count | Difficulty Mix |
|----------|-------|----------------|
| movement | 20 | 10 easy, 7 medium, 3 hard |
| text-objects | 18 | 6 easy, 8 medium, 4 hard |
| lsp-navigation | 15 | 5 easy, 7 medium, 3 hard |
| search-replace | 12 | 4 easy, 5 medium, 3 hard |
| refactoring | 10 | 2 easy, 5 medium, 3 hard |
| git-operations | 12 | 4 easy, 5 medium, 3 hard |
| window-management | 10 | 5 easy, 4 medium, 1 hard |
| buffer-management | 8 | 4 easy, 3 medium, 1 hard |
| folding | 8 | 3 easy, 4 medium, 1 hard |
| quickfix | 8 | 3 easy, 4 medium, 1 hard |
| diagnostics | 6 | 2 easy, 3 medium, 1 hard |
| telescope | 10 | 4 easy, 4 medium, 2 hard |
| surround | 8 | 3 easy, 4 medium, 1 hard |
| harpoon | 5 | 2 easy, 2 medium, 1 hard |
| formatting | 5 | 2 easy, 2 medium, 1 hard |

### Challenge Duration Variety
- **Quick** (1-5 keystrokes): 50 challenges - instant gratification
- **Standard** (6-15 keystrokes): 60 challenges - typical editing tasks
- **Complex** (16-40 keystrokes): 30 challenges - multi-step operations
- **Expert** (40+ keystrokes): 10 challenges - advanced refactoring

### Dynamic Hint System
- Challenges can include a `hint` field for non-standard keymaps
- Hints are read from user's nvim config at runtime when possible
- Standard vim commands show built-in hints

### Unit Test Coverage
- Validation tests for each validation type (exact_match, contains, cursor_position, etc.)
- Failure condition tests (wrong content, wrong position, timeout)
- Success condition tests for all 150 challenges
- Edge case tests (empty buffers, unicode, multiline)

## Impact

- Affected specs: `challenge-system`
- Affected code:
  - `game/internal/engine/assets/challenges.yaml` (expanded)
  - `lua/keyforge/challenges.lua` (new validation types if needed)
  - `tests/keyforge/challenges_spec.lua` (extensive new tests)
  - `game/internal/engine/challenges_test.go` (Go-side tests)
