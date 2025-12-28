# Design: Keyforge Core Architecture

## Context

Keyforge is a tower defense game integrated into Neovim that gamifies learning vim keybindings. The architecture must handle real-time game rendering while seamlessly integrating with Neovim's editing environment for challenge validation.

**Stakeholders**: Neovim users wanting to improve keybinding proficiency
**Constraints**: Must run in terminal, maintain 60fps, communicate with Neovim without blocking

## Goals / Non-Goals

### Goals
- Smooth 60fps game rendering in terminal
- Seamless Neovim integration with minimal plugin load time (<500ms)
- Accurate keystroke tracking and validation for challenges
- Extensible challenge and tower system via YAML configuration
- Cross-platform support (Linux, macOS, WSL)

### Non-Goals
- Graphical UI outside terminal (no Electron, no browser)
- Multiplayer support (future enhancement)
- Real-time Neovim buffer synchronization during gameplay
- Sound/audio effects (terminal bell only)

## Decisions

### Decision 1: Go + Bubbletea for Game Engine

**What**: Use Go with bubbletea TUI framework for the game engine, separate from Neovim.

**Why**:
- Bubbletea provides Elm-architecture for predictable state management
- Lipgloss enables rich terminal styling with true color support
- Go compiles to single binary, simplifying distribution
- Separation from Neovim allows independent game loop without editor blocking

**Alternatives considered**:
- Pure Lua in Neovim: Limited animation capabilities, would block editor
- Python + Rich: Slower startup, dependency management complexity
- Rust + Ratatui: Steeper learning curve, similar benefits to Go

### Decision 2: JSON-RPC over stdin/stdout

**What**: Use JSON-RPC for bidirectional communication between Go game and Neovim plugin.

**Why**:
- Simple protocol, easy to debug (human-readable)
- No external dependencies (sockets, files)
- Low latency for local process communication
- Neovim already supports jobstart with stdin/stdout pipes

**Protocol format**:
```json
{"jsonrpc": "2.0", "method": "request_challenge", "params": {...}, "id": 1}
{"jsonrpc": "2.0", "result": {...}, "id": 1}
```

**Alternatives considered**:
- Unix sockets: Platform compatibility issues on Windows
- Named pipes: More complex setup
- HTTP: Overhead, requires port management

### Decision 3: Challenge Validation in Lua

**What**: Perform all challenge validation in Neovim Lua plugin, not in Go.

**Why**:
- Direct access to buffer state and vim APIs
- Can track keystrokes via `vim.on_key()` callback
- Validation logic stays close to the editing context
- Game engine remains stateless about editor internals

**Trade-off**: Validation logic in Lua may be slower than Go, but latency is acceptable for human-speed interactions.

### Decision 4: YAML-Based Content Definition

**What**: Define challenges, towers, waves, and enemy types in YAML files.

**Why**:
- Human-readable and easy to modify
- Users can create custom challenges without coding
- Supports hot-reloading for development
- Separates content from code

**File locations**:
- `game/assets/challenges.yaml` - Kata definitions
- `game/assets/towers.yaml` - Tower configurations
- `game/assets/waves.yaml` - Wave spawn patterns

### Decision 5: Entity-Component Architecture (Simplified)

**What**: Use simple struct-based entities with behavior methods, not full ECS.

**Why**:
- Tower defense has limited entity types (enemies, towers, projectiles)
- Full ECS adds complexity without benefit at this scale
- Go's struct embedding provides sufficient composition
- Easier to understand and maintain

**Structure**:
```go
type Entity interface {
    Update(dt float64)
    Render() string
    Position() (x, y int)
}

type Enemy struct {
    BaseEntity
    Health, Speed, PathIndex int
}
```

## Risks / Trade-offs

### Risk: Terminal compatibility
Different terminals have varying support for Unicode, true color, and refresh rates.
- **Mitigation**: Test on popular terminals (iTerm2, Alacritty, Windows Terminal, Kitty). Provide fallback ASCII mode.

### Risk: RPC latency affecting gameplay
If challenge requests/responses are slow, game flow could feel interrupted.
- **Mitigation**: Keep RPC messages small. Pre-load challenge content. Use async processing where possible.

### Risk: Keystroke tracking accuracy
`vim.on_key()` may miss some inputs or count internal mappings.
- **Mitigation**: Test thoroughly with various keybind configurations. Document known limitations.

### Trade-off: Separate binary vs embedded
Running a separate Go process adds complexity (spawning, cleanup, path management).
- **Accepted**: Benefits of smooth rendering and separation of concerns outweigh complexity.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         NEOVIM                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  lua/keyforge/                                           │   │
│  │  ├── init.lua      (setup, keybinds)                    │   │
│  │  ├── rpc.lua       (JSON-RPC handler)                   │   │
│  │  ├── challenges.lua (validation, scoring)               │   │
│  │  └── ui.lua        (buffer management)                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              │ vim.fn.jobstart()                │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  Terminal Split (game renders here)                      │   │
│  └─────────────────────────────────────────────────────────┘   │
└──────────────────────────────│───────────────────────────────────┘
                               │
                    stdin/stdout (JSON-RPC)
                               │
┌──────────────────────────────▼───────────────────────────────────┐
│                      GO GAME BINARY                              │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  cmd/keyforge/main.go                                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  internal/                                               │   │
│  │  ├── engine/    (game loop, physics, waves)             │   │
│  │  ├── entities/  (tower, enemy, projectile)              │   │
│  │  ├── ui/        (bubbletea model, view, styles)         │   │
│  │  └── nvim/      (RPC client, protocol)                  │   │
│  └─────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │  assets/        (challenges.yaml, towers.yaml, etc.)    │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow: Challenge Cycle

1. Game decides challenge is needed (timer, event, or player request)
2. Game → Neovim: `{"method": "request_challenge", "params": {"category": "movement", "difficulty": 2}}`
3. Neovim selects appropriate challenge, opens buffer with initial content
4. User solves challenge in Neovim buffer
5. User signals completion (keybind or auto-detect)
6. Neovim validates: compares buffer state, counts keystrokes, measures time
7. Neovim → Game: `{"method": "challenge_complete", "params": {"success": true, "keystrokes": 12, "time_ms": 3500}}`
8. Game awards resources, updates UI, continues wave

## Open Questions

1. **Auto-validation vs manual completion**: Should challenges auto-detect completion or require explicit command?
   - Leaning toward: Hybrid - auto-detect for simple challenges, manual for complex ones

2. **Challenge interruption**: What happens if a wave reaches the end while user is mid-challenge?
   - Leaning toward: Pause game during active challenges, or provide "emergency cancel"

3. **Replay system**: Should we record and replay challenge solutions?
   - Leaning toward: Defer to future enhancement
