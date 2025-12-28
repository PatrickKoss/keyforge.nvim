# Tasks: Add Keyforge Core

## Phase 1: Proof of Concept

### 1.1 Go Project Setup
- [x] 1.1.1 Initialize Go module under `game/` with go.mod
- [x] 1.1.2 Add bubbletea and lipgloss dependencies
- [x] 1.1.3 Create cmd/keyforge/main.go entry point
- [x] 1.1.4 Set up Makefile with build targets

### 1.2 Basic Game Loop
- [x] 1.2.1 Implement bubbletea Model interface in `internal/ui/model.go`
- [x] 1.2.2 Create game grid rendering with box drawing characters
- [x] 1.2.3 Implement basic input handling (quit, pause)
- [x] 1.2.4 Add FPS-limited update loop (60fps target)

### 1.3 Core Entities
- [x] 1.3.1 Create Enemy struct with position, health, speed
- [x] 1.3.2 Implement hardcoded path following for enemies
- [x] 1.3.3 Create Tower struct with position, range, damage, cooldown
- [x] 1.3.4 Implement tower targeting and shooting logic

### 1.4 Basic Combat
- [x] 1.4.1 Create Projectile struct with position and velocity
- [x] 1.4.2 Implement projectile-enemy collision detection
- [x] 1.4.3 Add damage application and enemy death
- [x] 1.4.4 Implement player health reduction when enemies reach end

## Phase 2: Neovim Integration

### 2.1 Lua Plugin Structure
- [x] 2.1.1 Create `plugin/keyforge.lua` auto-load file
- [x] 2.1.2 Create `lua/keyforge/init.lua` with setup function
- [x] 2.1.3 Implement keybind registration (`<leader>kf`)
- [x] 2.1.4 Add configuration options (difficulty, colors, etc.)

### 2.2 Process Management
- [x] 2.2.1 Implement Go binary spawning in terminal split
- [x] 2.2.2 Add binary build check and auto-build on first run
- [x] 2.2.3 Create cleanup on buffer close/Neovim exit

### 2.3 RPC Communication
- [x] 2.3.1 Define JSON-RPC message types in Go (`internal/nvim/protocol.go`)
- [x] 2.3.2 Implement RPC client in Go (`internal/nvim/client.go`)
- [x] 2.3.3 Create Lua RPC handler (`lua/keyforge/rpc.lua`)
- [x] 2.3.4 Test bidirectional message passing

### 2.4 Challenge Buffer
- [x] 2.4.1 Create `lua/keyforge/ui.lua` for buffer management
- [x] 2.4.2 Implement challenge buffer creation with kata content
- [x] 2.4.3 Add manual completion command (`:KeyforgeComplete`)

## Phase 3: Challenge System

### 3.1 Challenge Definitions
- [x] 3.1.1 Create `game/internal/engine/assets/challenges.yaml` schema
- [x] 3.1.2 Define 30+ challenges across different categories
- [x] 3.1.3 Implement YAML loading in Go with embed

### 3.2 Validation Logic
- [x] 3.2.1 Create `lua/keyforge/challenges.lua` validation module
- [x] 3.2.2 Implement buffer state comparison (before/after)
- [x] 3.2.3 Add keystroke tracking via `vim.on_key()`
- [x] 3.2.4 Implement time tracking for challenge completion

### 3.3 Scoring System
- [x] 3.3.1 Calculate efficiency score (optimal vs actual keystrokes)
- [x] 3.3.2 Implement gold reward calculation
- [x] 3.3.3 Add difficulty multipliers
- [x] 3.3.4 Send completion results to game via RPC

## Phase 4: Core Gameplay Loop

### 4.1 Tower Types
- [x] 4.1.1 Implement Arrow Tower (basic motions)
- [x] 4.1.2 Implement LSP Tower (navigation challenges)
- [x] 4.1.3 Implement Refactor Tower (text object challenges)
- [x] 4.1.4 Add tower selection UI and placement

### 4.2 Enemy Types
- [x] 4.2.1 Create Bug enemy (basic, low health)
- [x] 4.2.2 Create Gremlin enemy (medium health, faster)
- [x] 4.2.3 Create Daemon enemy (high health, slow)
- [x] 4.2.4 Create Boss enemy type with special mechanics

### 4.3 Wave System
- [x] 4.3.1 Define wave configurations in Go
- [x] 4.3.2 Implement wave spawner with timing
- [x] 4.3.3 Add wave progression (difficulty scaling)
- [x] 4.3.4 Implement win/lose conditions

### 4.4 Economy
- [x] 4.4.1 Implement gold tracking and display
- [x] 4.4.2 Add tower purchase with gold cost
- [x] 4.4.3 Implement tower upgrade system
- [x] 4.4.4 Balance resource economy

## Phase 5: Content & Polish

### 5.1 Challenge Library
- [x] 5.1.1 Create 10+ movement mastery challenges
- [x] 5.1.2 Create 10+ text object challenges
- [x] 5.1.3 Create 5+ LSP navigation challenges
- [x] 5.1.4 Create 5+ search/replace challenges

### 5.2 Visual Polish
- [x] 5.2.1 Add projectile animations
- [x] 5.2.2 Implement explosion/death effects
- [x] 5.2.3 Add tower attack animations
- [x] 5.2.4 Polish UI with lipgloss styles

### 5.3 User Experience
- [x] 5.3.1 Add tutorial/first-time experience
- [x] 5.3.2 Implement pause menu
- [x] 5.3.3 Add game over screen with stats
- [x] 5.3.4 Create help/controls display

## Phase 6: Testing & Documentation

### 6.1 Go Tests
- [x] 6.1.1 Unit tests for engine logic
- [x] 6.1.2 Unit tests for entity behaviors
- [x] 6.1.3 Integration tests for RPC

### 6.2 Lua Tests
- [x] 6.2.1 Challenge validation tests (plenary.nvim)
- [x] 6.2.2 RPC handler tests
- [x] 6.2.3 UI helper tests

### 6.3 Documentation
- [x] 6.3.1 Write README with installation instructions
- [x] 6.3.2 Document configuration options
- [x] 6.3.3 Create custom challenge guide
