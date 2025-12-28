# Design: Refactored Gold Economy and Challenge System

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         Neovim Plugin (Lua)                      │
├─────────────────────────────────────────────────────────────────┤
│  keymap_hints.lua     │  challenge_queue.lua  │  ui.lua         │
│  - Discover keymaps   │  - Track available    │  - Split view   │
│  - Cache mappings     │  - User trigger       │  - Challenge    │
│  - Plugin detection   │  - Gold calculation   │    buffer       │
└────────────┬──────────┴──────────┬────────────┴────────┬────────┘
             │                     │                     │
             │              JSON-RPC Protocol            │
             │                     │                     │
┌────────────▼─────────────────────▼─────────────────────▼────────┐
│                       Go Game Engine                             │
├─────────────────────────────────────────────────────────────────┤
│  engine/game.go           │  engine/economy.go (NEW)            │
│  - Game loop continues    │  - Reduced mob gold                 │
│  - State: challenge_active│  - Challenge reward calc            │
│  - No pause on challenge  │  - Speed bonus formula              │
└─────────────────────────────────────────────────────────────────┘
```

## Component Design

### 1. Keymap Discovery System (Lua)

```lua
-- lua/keyforge/keymap_hints.lua

-- Structure for cached keymaps
local keymap_cache = {
  normal = {},      -- mode -> action -> keybinding
  visual = {},
  insert = {},
  plugins = {},     -- detected plugins
  last_refresh = 0,
}

-- Action categories mapped to common operations
local action_mappings = {
  ["telescope.find_files"] = { action = "find_file", category = "lsp-navigation" },
  ["telescope.live_grep"] = { action = "grep_project", category = "search-replace" },
  -- ... more mappings
}
```

**Key functions:**
- `discover_keymaps()`: Query nvim_get_keymap for all modes
- `detect_plugins()`: Check for common plugins via pcall(require)
- `get_hint_for_action(action)`: Return user's keybinding for an action
- `refresh_cache()`: Update cache (called on game start, config change)

### 2. Challenge Queue System (Lua)

```lua
-- lua/keyforge/challenge_queue.lua

local queue = {
  available = {},     -- challenges ready to be started
  current = nil,      -- active challenge (if any)
  completed = {},     -- session history
  pending_gold = 0,   -- gold earned but not yet sent to game
}
```

**Key functions:**
- `request_next()`: User-triggered, starts next challenge
- `complete_current(result)`: Validate and calculate gold
- `get_challenge_with_hints(challenge)`: Enrich challenge with keybinding hints
- `calculate_speed_bonus(time_ms, par_time_ms)`: Speed multiplier formula

### 3. Economy Module (Go)

```go
// game/internal/engine/economy.go

type EconomyConfig struct {
    MobGoldMultiplier     float64 // 0.25 = 25% of original
    WaveBonusMultiplier   float64 // 0.50 = 50% of original
    ChallengeBaseGold     int     // Base gold for difficulty 1
    ChallengeSpeedMaxMult float64 // Max speed bonus multiplier (2.0)
}

func DefaultEconomyConfig() EconomyConfig {
    return EconomyConfig{
        MobGoldMultiplier:     0.25,
        WaveBonusMultiplier:   0.50,
        ChallengeBaseGold:     25,
        ChallengeSpeedMaxMult: 2.0,
    }
}
```

### 4. Game State Changes (Go)

Current states: `Menu, Playing, Paused, ChallengeActive, WaveComplete, GameOver, Victory`

**Change**: `ChallengeActive` no longer pauses the game loop. Instead:
- Game continues updating enemies, towers, projectiles
- UI renders both game grid and indicates challenge is active
- Lua side manages challenge buffer separately

### 5. RPC Protocol Extensions

**New methods:**

```go
// Neovim -> Game
MethodStartChallenge    = "start_challenge"     // User pressed trigger
MethodChallengeComplete = "challenge_complete"  // Challenge finished (existing, extended)

// Game -> Neovim
MethodChallengeAvailable = "challenge_available" // Notify available challenges
MethodGoldUpdate         = "gold_update"         // Gold changed (challenge reward)
```

**Extended ChallengeResult:**

```go
type ChallengeResult struct {
    RequestID      string  `json:"request_id"`
    Success        bool    `json:"success"`
    Skipped        bool    `json:"skipped,omitempty"`
    KeystrokeCount int     `json:"keystroke_count,omitempty"`
    TimeMs         int     `json:"time_ms,omitempty"`
    Efficiency     float64 `json:"efficiency,omitempty"`
    SpeedBonus     float64 `json:"speed_bonus,omitempty"`     // NEW
    GoldEarned     int     `json:"gold_earned,omitempty"`     // NEW
}
```

## Gold Calculation Formula

```
base_gold = challenge.gold_base
difficulty_mult = 1.0 + (difficulty * 0.25)  // 1.25 for d1, 1.5 for d2, 1.75 for d3
efficiency_mult = 0.5 + (efficiency * 0.5)   // 50-100% based on keystrokes

speed_bonus = 1.0
if time_ms < par_time_ms:
    speed_ratio = par_time_ms / time_ms
    speed_bonus = min(2.0, 1.0 + (speed_ratio - 1.0) * 0.5)

total_gold = floor(base_gold * difficulty_mult * efficiency_mult * speed_bonus)
```

Example:
- Base: 50g, Difficulty 2 (1.5x), 90% efficiency (0.95x), 1.5x speed bonus
- Total: 50 * 1.5 * 0.95 * 1.5 = 106g

## Keybinding Hint Display

Challenge description format with hints:

```
╭─ Challenge: Find Configuration File ─────────────────────╮
│                                                          │
│ Use your fuzzy finder to locate and open 'init.lua'      │
│                                                          │
│ ┌─ Your Keybindings ─────────────────────────────────┐   │
│ │ Find files: <leader>ff  (telescope.find_files)    │   │
│ │ Live grep:  <leader>fg  (telescope.live_grep)     │   │
│ │ Buffers:    <leader>fb  (telescope.buffers)       │   │
│ └───────────────────────────────────────────────────┘   │
│                                                          │
│ Par: 3 keystrokes │ Base reward: 40g │ Difficulty: ★★☆  │
╰──────────────────────────────────────────────────────────╯
```

## Plugin Detection Strategy

```lua
local plugins_to_detect = {
  { name = "telescope", require_path = "telescope" },
  { name = "nvim-tree", require_path = "nvim-tree" },
  { name = "neo-tree", require_path = "neo-tree" },
  { name = "fugitive", check = function() return vim.fn.exists(":Git") > 0 end },
  { name = "gitsigns", require_path = "gitsigns" },
  { name = "nvim-surround", require_path = "nvim-surround" },
  { name = "mini.surround", require_path = "mini.surround" },
  { name = "flash.nvim", require_path = "flash" },
  { name = "leap.nvim", require_path = "leap" },
}
```

For each detected plugin, relevant challenges become available and hints show the user's actual bindings for that plugin's actions.

## View Layout During Challenge

```
┌─────────────────────────────────────────────────────────────────┐
│ Wave 3/10 │ Gold: 450 │ Health: ████████░░ 80% │ Towers: 5     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Tower Defense Grid (continues updating)                       │
│   Enemies moving, towers firing                                 │
│                                                                 │
│   [Game shrinks or shows mini-map when challenge active]        │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ CHALLENGE ACTIVE │ Time: 00:04 │ Keystrokes: 3                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   Challenge buffer (Neovim buffer with syntax highlighting)     │
│                                                                 │
│   const message = "old value";                                  │
│   │                                                             │
│                                                                 │
│   Goal: Change "old value" to "new value"                       │
│   Hint: ci" (change inside quotes)                              │
│                                                                 │
├─────────────────────────────────────────────────────────────────┤
│ <leader>kc Complete │ <leader>ks Skip │ <leader>kn Next        │
└─────────────────────────────────────────────────────────────────┘
```

## Trade-offs Considered

### Option A: Pause game during challenge (rejected)
- Pros: Less stressful, focus on learning
- Cons: Removes time pressure, less engaging, doesn't reward speed

### Option B: Separate challenge/tower phases (rejected)
- Pros: Clear separation of concerns
- Cons: Breaks real-time feel, reduces strategic depth

### Option C: Concurrent play (chosen)
- Pros: Time pressure rewards mastery, strategic choices matter
- Cons: More complex, may feel overwhelming initially

**Mitigation**: Difficulty settings can adjust mob speed and challenge par times.

## Migration Path

1. Phase 1: Add economy module with configurable multipliers
2. Phase 2: Implement keymap discovery (can be tested standalone)
3. Phase 3: Add user-triggered challenge flow
4. Phase 4: Create plugin-aware challenges
5. Phase 5: Tune balance through playtesting
