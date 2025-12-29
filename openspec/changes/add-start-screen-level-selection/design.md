## Context
The game currently starts directly into gameplay with hardcoded settings. This change introduces a start screen flow that allows level selection, game configuration, and level preview before starting. The design must support future level expansion while keeping the current experience as the default.

## Goals / Non-Goals
- Goals:
  - Create an intuitive start screen with keyboard navigation
  - Support level preview with path, enemies, and tower visualization
  - Allow game settings configuration (difficulty, speed, resources)
  - Pass nvim config defaults to the game as initial settings
  - Design level system for easy future level additions
- Non-Goals:
  - Level editor (levels are code-defined)
  - Online leaderboards or level sharing
  - Dynamic level generation

## Decisions

### Decision: Start Screen Flow
The start screen will have three sequential phases:
1. **Level Browser** - Grid/list of available levels with preview panel
2. **Settings Menu** - Difficulty, speed, starting resources configuration
3. **Game Start** - Transition to playing state

Rationale: Sequential flow keeps the UI simple and prevents overwhelming new players. Settings apply to the selected level.

### Decision: Level Definition Structure
Levels will be defined in Go code as structs containing:
```go
type Level struct {
    ID          string
    Name        string
    Description string
    GridSize    struct{ Width, Height int }
    Path        []Position
    TotalWaves  int
    WaveFunc    func(waveNum int) Wave  // Custom wave generator
    AllowedTowers []TowerType
    EnemyTypes    []EnemyType           // Enemies that appear in this level
    Difficulty    string                // "beginner", "intermediate", "advanced"
}
```
Rationale: Code-defined levels allow complex wave logic and validation at compile time. YAML levels can be added later if needed.

### Decision: Level Preview Rendering
The preview will show:
- Mini-grid with path highlighted
- List of enemy types that appear (with icons)
- List of available towers (with icons)
- Wave count and difficulty indicator

Rationale: Players need to understand what they're selecting. Visual preview is more intuitive than text descriptions alone.

### Decision: Game Settings
Configurable settings:
| Setting | Options | Default | Description |
|---------|---------|---------|-------------|
| Difficulty | easy/normal/hard | normal | Economy multipliers (from existing EconomyConfig) |
| Game Speed | 0.5x/1x/1.5x/2x | 1x | Time multiplier for game updates |
| Starting Gold | 100-500 (slider) | 200 | Initial gold |
| Starting Health | 50-200 (slider) | 100 | Initial health |

Rationale: These settings exist in the codebase but aren't exposed. Speed multiplier is new and helps advanced players or those short on time.

### Decision: Nvim Config Integration
The Lua config will be passed via command-line flags:
```bash
keyforge --nvim-mode --rpc-socket /tmp/... --difficulty normal --starting-gold 200 --starting-health 100 --game-speed 1.0
```
These become the default values in the settings menu (user can still change them).

Rationale: Command-line flags are simple and don't require protocol changes. The game already uses flags for `--nvim-mode` and `--rpc-socket`.

### Decision: State Machine Update
New states for start screen flow:
```
StateStartScreen -> StateLevelSelect -> StateSettings -> StatePlaying
```
The existing `StateMenu` can be repurposed as `StateStartScreen`.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Single level feels incomplete | Design "Classic" level as the polished default; future levels are bonus content |
| Settings complexity | Keep settings to one screen; use sensible defaults from nvim config |
| Preview rendering complexity | Start with simple ASCII preview; enhance later if needed |

## Migration Plan
1. Existing behavior (immediate game start) becomes pressing Enter on default level with default settings
2. No breaking changes to existing config - new settings are additive
3. Current wave generation moves into "Classic" level definition

## Open Questions
- Should levels be locked/unlocked based on completion? (Recommendation: No, keep all unlocked for MVP)
- Should settings persist between sessions? (Recommendation: Use nvim config as source of truth)
