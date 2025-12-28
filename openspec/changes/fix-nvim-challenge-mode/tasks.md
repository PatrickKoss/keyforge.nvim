## 1. Enable nvim-mode in plugin launch

- [x] 1.1 Update `lua/keyforge/init.lua` to pass `--nvim-mode` flag to game binary
- [x] 1.2 Update `lua/keyforge/init.lua` to use `termopen` with stderr handler for RPC
- [x] 1.3 Register RPC handlers for `request_challenge`, `game_state`, and `game_over`
- [x] 1.4 Build passes, game launches with RPC communication

## 2. Implement game pause during challenge

- [x] 2.1 Add `StateChallengeWaiting` state to game engine (distinct from `StateChallengeActive`)
- [x] 2.2 Update `startChallenge()` in `model.go` to use `StartChallengeWaiting()` in nvim mode
- [x] 2.3 Add `handleChallengeWaitingKeys()` for nvim mode (only allows Escape to cancel)
- [x] 2.4 Game loop already skips updates when not in StatePlaying or StateChallengeActive
- [x] 2.5 Add "CHALLENGE IN PROGRESS - Game paused" status to game view

## 3. Rewrite challenge buffer to use real files

- [x] 3.1 Create `lua/keyforge/challenge_buffer.lua` module for file-based challenge buffers
- [x] 3.2 Implement `create_temp_file()` - writes content to temp file with correct extension
- [x] 3.3 Implement `start_challenge()` - opens file in new tab with keymaps
- [x] 3.4 Set up `BufWipeout` autocmd for challenge completion/cancel
- [x] 3.5 Implement `submit_challenge()` - reads file, validates, sends result
- [x] 3.6 Implement `_cleanup()` - deletes temp file, returns to game tab
- [x] 3.7 Add challenge info floating window (top-right corner)

## 4. Update validation for file-based buffers

- [x] 4.1 `challenge_buffer.lua` reads from file path and validates via `challenges.lua`
- [x] 4.2 Keystroke tracking via existing `challenges.start_tracking()` / `stop_tracking()`
- [x] 4.3 Keystroke tracking uses `vim.on_key()` which works across tabs
- [x] 4.4 Add timeout handling (configurable via `challenge_timeout`, default 300s)

## 5. Handle game state restoration

- [x] 5.1 `HandleChallengeComplete` already exists and calls `Game.EndChallenge()`
- [x] 5.2 Game state is frozen during challenge (no enemy movement)
- [x] 5.3 Edge case handled - game can't end during StateChallengeWaiting
- [x] 5.4 `focus_game_tab()` returns to game and calls `startinsert`

## 6. Add game over/victory detection

- [x] 6.1 Add `game_over` and `victory` RPC notifications in `protocol.go` and `client.go`
- [x] 6.2 Create `lua/keyforge/game_over.lua` module for end-game UI
- [x] 6.3 Display game over/victory screen with ASCII art and stats
- [x] 6.4 Handle restart command (`HandleRestart()` in Go, `restart_game` RPC)

## 7. Configuration and keymaps

- [x] 7.1 Add config options `keybind_submit` (default `<CR>`) and `keybind_cancel` (default `<Esc>`)
- [x] 7.2 Add config option `challenge_timeout` (default 300 seconds)
- [x] 7.3 Game pause behavior is default in nvim mode (no config needed)
- [x] 7.4 Documentation deferred to separate task

## 8. Testing and cleanup

- [x] 8.1 Go tests pass (`make test`)
- [x] 8.2 Go build passes (`make build`)
- [ ] 8.3 Manual testing with actual Neovim integration (requires user testing)
- [ ] 8.4 Test with multiple LSP servers (requires user testing)
- [ ] 8.5 Test game over during challenge scenario (requires user testing)
