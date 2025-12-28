# Keyforge.nvim

A tower defense game integrated into Neovim that gamifies learning vim keybindings. Players defend against waves of bugs by completing kata-style editing challenges.

## Project Overview

This is a **dual-stack application**:
- **Lua plugin** (`lua/keyforge/`): Integrates with Neovim, handles challenges, validates user input
- **Go game engine** (`game/`): Standalone TUI using Bubbletea, communicates via JSON-RPC over stdin/stdout

## Quick Commands

```bash
# Build the game binary
make build

# Run standalone (for testing without Neovim)
make run

# Run Go tests
make test

# Format code
make fmt

# Lint (requires golangci-lint)
make lint

# Build for all platforms
make release

# Symlink for development
make install-dev
```

## Project Structure

```
keyforge.nvim/
├── plugin/keyforge.lua           # Auto-loader, creates :Keyforge commands
├── lua/keyforge/
│   ├── init.lua                  # Main module - setup(), start(), stop()
│   ├── rpc.lua                   # JSON-RPC 2.0 bidirectional communication
│   ├── challenges.lua            # Challenge validation, keystroke tracking
│   └── ui.lua                    # Challenge buffer management
├── game/                         # Go game engine
│   ├── cmd/keyforge/main.go      # Entry point
│   ├── bin/keyforge              # Built binary
│   └── internal/
│       ├── engine/
│       │   ├── game.go           # Game state, update loop, core logic
│       │   ├── wave.go           # Wave generation, difficulty scaling
│       │   ├── challenges.go     # Challenge loading from YAML
│       │   └── assets/challenges.yaml  # 30+ challenge definitions
│       ├── entities/
│       │   ├── types.go          # Tower/Enemy type configs
│       │   ├── tower.go          # Tower targeting, upgrades
│       │   ├── enemy.go          # Enemy pathing, health
│       │   └── effects.go        # Visual effects
│       ├── ui/
│       │   ├── model.go          # Bubbletea model, 60fps game loop
│       │   ├── view.go           # Grid rendering, HUD
│       │   └── styles.go         # Lipgloss colors, emoji chars
│       └── nvim/
│           └── protocol.go       # RPC message types
├── tests/                        # Lua tests (plenary-based)
├── openspec/                     # Spec-driven development
└── Makefile
```

## Domain Knowledge

### Game Mechanics

**Resources:**
- Gold: Earned by killing enemies or completing challenges (start: 200)
- Health: Lost when enemies reach the end (start: 100, game over at 0)

**Tower Types:**
| Tower | Cost | Category | Range | Description |
|-------|------|----------|-------|-------------|
| Arrow (🏹) | 50g | movement | 3.0 | Basic, fast attacks |
| LSP (🔮) | 100g | lsp-navigation | 5.0 | Long range, slower |
| Refactor (⚡) | 150g | text-objects | 2.5 | Area damage |

**Enemy Types:**
| Enemy | Health | Speed | Damage | Gold |
|-------|--------|-------|--------|------|
| Bug (🐛) | 50 | 1.5 | 10 | 15 |
| Gremlin (👹) | 80 | 2.5 | 15 | 25 |
| Daemon (👿) | 150 | 1.0 | 25 | 40 |
| Boss (💀) | 500 | 0.5 | 50 | 100 |

**Wave Progression (10 waves):**
- Waves 1-3: Bugs only
- Waves 4-6: Bugs + Gremlins
- Waves 7-9: Bugs + Gremlins + Daemons
- Wave 10: Boss wave

### Challenge System

Challenges are triggered when towers fire. Categories match tower types:
- `movement`: Basic vim motions (`$`, `w`, `f`, `G`)
- `text-objects`: `ciw`, `da(`, `yi"`
- `lsp-navigation`: `gd`, `gr`, `K`
- `search-replace`: `/`, `:s`, `:g`
- `refactoring`: Extract function, inline variable

**Validation types** (in `challenges.lua`):
- `exact_match`: Buffer matches expected exactly
- `contains`: Buffer contains expected string
- `function_exists`: Function pattern found
- `pattern`: Lua pattern match
- `different`: Content changed from initial

**Scoring:** `efficiency = par_keystrokes / actual_keystrokes`
- Gold reward scales with efficiency

### Architecture

**Communication Flow:**
1. Game spawns enemy → triggers challenge request (RPC)
2. Lua creates challenge buffer in Neovim
3. User edits with vim commands
4. User completes (`<leader>kc`) → Lua validates
5. Lua sends result back (RPC) with efficiency score
6. Game awards gold, resumes

**Key Patterns:**
- Separation: entities/ (domain) → engine/ (logic) → ui/ (rendering)
- Bubbletea Model-Update-View pattern for TUI
- JSON-RPC 2.0 with newline-delimited messages
- Embedded YAML for challenge definitions

## Go Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework (60fps event loop)
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `gopkg.in/yaml.v3` - Challenge YAML parsing

## Lua Dependencies

- `plenary.nvim` - Testing framework
- Neovim 0.11+ APIs: `vim.fn.jobstart`, `vim.on_key`, `vim.api.nvim_buf_*`

## Testing

```bash
# Go tests (691 lines across 3 files)
cd game && go test -v ./...

# Lua tests
nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"
```

## Neovim Commands

| Command | Description |
|---------|-------------|
| `:Keyforge` | Start the game |
| `:KeyforgeStop` | Stop the game |
| `:KeyforgeBuild` | Rebuild binary |
| `:KeyforgeComplete` | Complete current challenge |

## Game Controls

| Key | Action |
|-----|--------|
| `h/j/k/l` | Move cursor |
| `1/2/3` | Select tower |
| `Space/Enter` | Place tower |
| `u` | Upgrade tower |
| `p` | Pause |
| `q` | Quit |

## Rendering Notes

- Grid: 20x14 cells
- Emojis are 2-char width in terminals - use `lipgloss.Width()` for padding
- 60 FPS via Bubbletea tick messages

---

<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->