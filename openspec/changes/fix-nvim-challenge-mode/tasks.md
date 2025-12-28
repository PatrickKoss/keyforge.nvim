## 1. Unix Socket RPC Architecture

- [x] 1.1 Create `game/internal/nvim/socket.go` with SocketServer implementation
- [x] 1.2 Add `RPCClient` interface to `protocol.go` for polymorphism
- [x] 1.3 Add `--rpc-socket` CLI flag to `main.go`
- [x] 1.4 Add `InitNvimSocket()` method to `model.go`
- [x] 1.5 Update RPC calls in `model.go` to use `NvimRPC` interface

## 2. Lua Socket Client

- [x] 2.1 Refactor `lua/keyforge/rpc.lua` to use `vim.loop.new_pipe()` for socket
- [x] 2.2 Implement `connect(socket_path, on_connect, on_error)` function
- [x] 2.3 Implement `disconnect()` with proper cleanup
- [x] 2.4 Keep handler registration unchanged (`M.on()`, `M.register_handlers()`)

## 3. Neovim Plugin Integration

- [x] 3.1 Generate unique socket path (`/tmp/keyforge-{pid}-{time}.sock`)
- [x] 3.2 Pass `--nvim-mode --rpc-socket <path>` to game binary
- [x] 3.3 Implement connection retry with exponential backoff (200ms start, 1s max)
- [x] 3.4 Clean up socket file in `M.stop()` and on game exit
- [x] 3.5 Clean up stale sockets on plugin load

## 4. Game Engine Updates (from prior work)

- [x] 4.1 `StateChallengeWaiting` state for pausing during challenges
- [x] 4.2 `HandleChallengeComplete` processes results from Neovim
- [x] 4.3 Game over/victory RPC notifications implemented
- [x] 4.4 "CHALLENGE IN PROGRESS" status in game view

## 5. Challenge Buffer System (from prior work)

- [x] 5.1 `challenge_buffer.lua` creates real temp files for challenges
- [x] 5.2 Opens file in new tab with proper filetype for LSP
- [x] 5.3 Keystroke tracking via `challenges.lua`
- [x] 5.4 Submit/cancel keymaps with validation
- [x] 5.5 Info floating window with challenge details

## 6. Game Over/Victory UI (from prior work)

- [x] 6.1 `game_over.lua` displays ASCII art for end states
- [x] 6.2 Stats shown: wave, gold, health, towers
- [x] 6.3 Restart/quit prompts visible

## 7. Build Verification

- [x] 7.1 Go tests pass (`make test`)
- [x] 7.2 Go build passes (`make build`)

## 8. Manual Testing (requires user)

- [ ] 8.1 Launch game with `:Keyforge`, verify movement (h/j/k/l)
- [ ] 8.2 Press 'c' to start challenge, verify buffer opens in new tab
- [ ] 8.3 Verify LSP works in challenge buffer (K for hover, gd for definition)
- [ ] 8.4 Verify `:5` and other ex commands work
- [ ] 8.5 Submit challenge and verify return to game with gold award
- [ ] 8.6 Test game over notification displays properly
