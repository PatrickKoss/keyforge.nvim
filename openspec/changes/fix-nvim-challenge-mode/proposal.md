# Change: Fix Neovim Challenge Mode Integration

## Why
The challenge system in nvim mode is fundamentally broken. When a user starts a challenge:
1. The `--nvim-mode` flag is not passed to the game binary, so challenges use the internal vim emulator instead of real Neovim
2. Even when nvim mode works, challenges open in a scratch floating window where LSP, commands (`:5`), and other Neovim features don't work
3. The game continues running during challenges, causing divided attention and game-over while editing
4. Game state (enemy positions, health) diverges from what the user sees after completing a challenge

Users cannot use their normal vim bindings, LSP (K for hover, gd for go-to-definition), or ex commands during challenges.

## What Changes
- **Lua plugin**: Pass `--nvim-mode` flag when launching game binary
- **Challenge buffer**: Create a REAL buffer in a new tab instead of scratch floating window
  - Buffer has proper filetype for LSP attachment
  - User can use all normal vim features
  - Challenge file is created in temp directory with correct extension
- **Game state management**: Game PAUSES during challenge, stores state snapshot
  - When challenge completes, game resumes from paused state
  - No enemy movement or health loss during challenge editing
- **Challenge completion**: Use `<CR>` (enter) to submit, `<Esc>` to cancel (configurable keymaps)
- **Buffer restoration**: After challenge, return user to game tab seamlessly
- **Game over screen**: **BREAKING** - Detect game over state and display message in Neovim when game ends

## Impact
- Affected specs: `challenge-system`, `neovim-integration`, `game-engine`
- Affected code:
  - `lua/keyforge/init.lua` - Pass `--nvim-mode` flag
  - `lua/keyforge/challenge_queue.lua` - Rewrite buffer creation to use real files
  - `lua/keyforge/rpc.lua` - Handle game state sync messages
  - `game/internal/ui/model.go` - Pause game during challenge
  - `game/internal/nvim/protocol.go` - Add game state messages
