# Change: Add Keyforge Core Tower Defense Game

## Why

Neovim users need an engaging, gamified way to learn vim keybindings and plugin workflows. Current resources (tutorials, cheat sheets, practice files) lack engagement and don't provide measurable progress. A tower defense game integrated into Neovim combines skill practice with entertainment, making keybinding mastery feel rewarding rather than tedious.

## What Changes

This proposal introduces the complete keyforge.nvim system:

- **Go Game Engine**: Bubbletea-based TUI rendering a 60fps tower defense game with smooth animations
- **Neovim Plugin**: Lua integration for launching the game, managing challenges, and communicating via JSON-RPC
- **Challenge System**: Kata-style text editing challenges that validate user proficiency and award resources
- **Tower Defense Mechanics**: Multiple tower types mapped to vim skill categories, enemy waves, and progression systems

### Components

1. **Game Engine** (Go + bubbletea)
   - Main game loop with state machine
   - Physics system for movement and collision
   - Entity management (towers, enemies, projectiles)
   - Terminal UI rendering with lipgloss styling

2. **Neovim Integration** (Lua)
   - Plugin loader and keybind registration
   - JSON-RPC communication over stdin/stdout
   - Challenge buffer management
   - Keystroke tracking and validation

3. **Challenge System**
   - YAML-defined kata challenges
   - Multi-category support (movement, LSP, text objects, search, refactoring, git)
   - Efficiency scoring based on keystroke count and time
   - Progressive difficulty scaling

4. **Tower Defense Mechanics**
   - 6 tower types aligned with vim skill categories
   - Wave-based enemy spawning with increasing difficulty
   - Resource economy (gold from challenges, health system)
   - Tower placement, upgrades, and abilities

## Impact

- Affected specs: `game-engine`, `neovim-integration`, `challenge-system`, `tower-defense` (all new)
- Affected code: Creates new Go module under `game/` and Lua plugin under `lua/keyforge/`
- Dependencies: Go 1.21+, Neovim 0.11+, bubbletea, lipgloss, plenary.nvim
