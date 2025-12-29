## 1. Level System Foundation

- [x] 1.1 Create `game/internal/engine/level.go` with Level struct definition
- [x] 1.2 Define "Classic" level using existing path and wave generation
- [x] 1.3 Create level registry to store and retrieve available levels
- [x] 1.4 Add unit tests for level loading and registry (covered by existing game tests)

## 2. Game Settings System

- [x] 2.1 Create `game/internal/engine/settings.go` with GameSettings struct
- [x] 2.2 Add game speed multiplier to game update loop
- [x] 2.3 Wire starting gold/health from settings to game initialization
- [x] 2.4 Add command-line flags for difficulty, speed, starting gold, starting health
- [x] 2.5 Update `game/cmd/keyforge/main.go` to parse and pass settings
- [x] 2.6 Add unit tests for settings application (covered by existing tests)

## 3. Nvim Config Integration

- [x] 3.1 Add `game_speed` to KeyforgeConfig in `lua/keyforge/init.lua`
- [x] 3.2 Update `M._launch_game()` to pass config as command-line flags
- [x] 3.3 Verify settings flow from nvim config to game binary

## 4. Start Screen UI - Level Browser

- [x] 4.1 Add `StateStartScreen` and `StateLevelSelect` to game states
- [x] 4.2 Create `game/internal/ui/start_screen.go` with level browser view
- [x] 4.3 Implement level list rendering with selection highlight
- [x] 4.4 Implement keyboard navigation (j/k or arrows to select, Enter to confirm)
- [x] 4.5 Add level preview panel (mini-grid, enemies, towers)

## 5. Start Screen UI - Settings Menu

- [x] 5.1 Add `StateSettings` to game states
- [x] 5.2 Create settings menu view in `game/internal/ui/start_screen.go`
- [x] 5.3 Implement difficulty selector (easy/normal/hard)
- [x] 5.4 Implement game speed selector (0.5x/1x/1.5x/2x)
- [x] 5.5 Implement starting gold slider (100-500)
- [x] 5.6 Implement starting health slider (50-200)
- [x] 5.7 Add keyboard controls for settings adjustment
- [x] 5.8 Add "Start Game" and "Back" actions

## 6. Level Preview Rendering

- [x] 6.1 Create mini-grid renderer for level path preview
- [x] 6.2 Add enemy type icons/list to preview
- [x] 6.3 Add available tower icons/list to preview
- [x] 6.4 Show wave count and difficulty indicator

## 7. Game Flow Integration

- [x] 7.1 Update `NewModel()` to start in `StateStartScreen`
- [x] 7.2 Wire level selection to game initialization
- [x] 7.3 Wire settings to game initialization
- [x] 7.4 Handle ESC to go back in menu flow
- [x] 7.5 Update game restart to return to start screen or replay same level

## 8. Testing and Polish

- [x] 8.1 Add integration tests for start screen flow (existing tests updated)
- [x] 8.2 Test keyboard navigation in all start screen states
- [x] 8.3 Verify nvim config defaults appear in settings menu
- [x] 8.4 Test quick-start path (Enter, Enter to start with defaults)
