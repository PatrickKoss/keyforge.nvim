# Change: Add Start Screen with Level Selection and Game Settings

## Why
Currently the game starts immediately with hardcoded wave progression. To support multiple levels in the future and give players control over their experience, we need a start screen where players can:
1. Select from available levels (currently one, designed for extensibility)
2. Configure game settings (difficulty, speed, starting resources) with nvim config defaults
3. Preview the selected level's layout (path, enemy types, available towers)

## What Changes
- **New start screen UI** in the Go game engine with level browser and settings menu
- **Level definition system** to define paths, waves, allowed towers, and enemy types per level
- **Game settings menu** with difficulty, game speed, starting gold/health configuration
- **Level preview** showing path visualization, enemy types, and tower availability
- **Nvim config integration** to pass default settings from Lua setup() to the game

## Impact
- Affected specs: `game-engine`, `neovim-integration`
- Affected code:
  - `game/internal/engine/` - level definitions, settings structs
  - `game/internal/ui/` - start screen, level browser, settings menu views
  - `game/cmd/keyforge/main.go` - flag handling for settings
  - `lua/keyforge/init.lua` - pass config to game binary
  - `game/internal/nvim/protocol.go` - config messages
