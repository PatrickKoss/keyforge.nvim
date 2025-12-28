# Tasks: Refactor Gold Economy and Challenge System

## Phase 1: Economy Module (Go)

- [x] **1.1** Create `game/internal/engine/economy.go` with `EconomyConfig` struct and default values
- [x] **1.2** Add `MobGoldMultiplier` to enemy gold calculation in `game.go:updateProjectiles()`
- [x] **1.3** Add `WaveBonusMultiplier` to wave completion bonus in `game.go:updateWaveSpawning()`
- [x] **1.4** Add speed bonus calculation function `CalculateSpeedBonus(timeMs, parTimeMs int) float64`
- [x] **1.5** Write unit tests for economy calculations in `economy_test.go`

## Phase 2: Keymap Discovery (Lua)

- [x] **2.1** Create `lua/keyforge/keymap_hints.lua` module structure
- [x] **2.2** Implement `discover_keymaps()` using `vim.api.nvim_get_keymap()`
- [x] **2.3** Implement `detect_plugins()` with pcall-based plugin detection
- [x] **2.4** Build action-to-keymap mapping table for common vim actions
- [x] **2.5** Implement `get_hint_for_action(action, category)` function
- [x] **2.6** Add keymap cache with `refresh_cache()` function
- [x] **2.7** Write plenary tests for keymap discovery in `tests/keymap_hints_spec.lua`

## Phase 3: Challenge Queue System (Lua)

- [x] **3.1** Create `lua/keyforge/challenge_queue.lua` module
- [x] **3.2** Implement challenge queue data structure with available/current/completed tracking
- [x] **3.3** Add `request_next()` function for user-triggered challenge start
- [x] **3.4** Implement `complete_current(result)` with gold calculation
- [x] **3.5** Add `get_challenge_with_hints(challenge)` to enrich challenges with keymap hints
- [x] **3.6** Update `init.lua` to expose new keybindings: `<leader>kn` (next challenge)
- [x] **3.7** Write tests for challenge queue in `tests/challenge_queue_spec.lua`

## Phase 4: Game State Changes (Go)

- [x] **4.1** Modify `StateChallengeActive` to NOT pause game loop in `game.go:Update()`
- [x] **4.2** Update `ui/model.go` to render game grid during challenge active state (already works)
- [x] **4.3** Add challenge status indicator to HUD in `ui/view.go`
- [x] **4.4** Extend `ChallengeResult` struct in `nvim/protocol.go` with `SpeedBonus` and `GoldEarned`
- [x] **4.5** Add `MethodStartChallenge` RPC handler for user-triggered challenges
- [x] **4.6** Update game tests for concurrent challenge/gameplay

## Phase 5: RPC Protocol Extensions

- [x] **5.1** Add `start_challenge` method to protocol (Neovim -> Game)
- [x] **5.2** Add `gold_update` notification (Game -> Neovim)
- [x] **5.3** Update `rpc.lua` to handle new methods
- [x] **5.4** Implement gold notification handler in Lua UI

## Phase 6: Plugin-Aware Challenges (YAML + Lua)

- [x] **6.1** Add Telescope challenges to `assets/challenges.yaml` (find_files, live_grep, buffers)
- [x] **6.2** Add nvim-tree/neo-tree navigation challenges
- [x] **6.3** Add fugitive/git-related challenges
- [x] **6.4** Add window/buffer navigation challenges
- [x] **6.5** Add surround-plugin challenges (nvim-surround, mini.surround)
- [x] **6.6** Update challenge metadata with `required_plugin` field
- [x] **6.7** Filter available challenges based on detected plugins in `challenge_queue.lua`

## Phase 7: UI Updates

- [x] **7.1** Update challenge buffer UI to display keymap hints
- [x] **7.2** Add gold earned animation/notification after challenge completion
- [x] **7.3** Add "Next Challenge" prompt with estimated reward display
- [x] **7.4** Update status line to show challenge availability indicator

## Phase 8: Integration & Balance

- [x] **8.1** Integration test: Complete game using only challenge gold
- [x] **8.2** Balance pass: Adjust gold values so game is winnable
- [x] **8.3** Add difficulty presets (easy: 50% mob gold, normal: 25%, hard: 0%)
- [x] **8.4** Update README with new mechanics documentation

## Dependencies

```
Phase 1 ──┬──> Phase 4
          │
Phase 2 ──┼──> Phase 3 ──> Phase 6
          │
          └──> Phase 5

Phase 4 + Phase 5 ──> Phase 7

All phases ──> Phase 8
```

## Parallelizable Work

- Phase 1 (Go economy) and Phase 2 (Lua keymap) can run in parallel
- Phase 6 (YAML challenges) can start once Phase 2 is complete
- Phase 7 (UI) requires Phase 4 and Phase 5

## Validation Checkpoints

After Phase 1: `go test ./internal/engine/...` passes - DONE
After Phase 2: `nvim --headless -c "luafile tests/keymap_hints_spec.lua"` passes
After Phase 4: Game runs without pausing during challenges - DONE
After Phase 6: 10+ new plugin-aware challenges available - DONE (20+ added)
After Phase 8: Full game playthrough successful with challenge-focused economy
