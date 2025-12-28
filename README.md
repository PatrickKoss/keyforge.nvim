# keyforge.nvim

A tower defense game integrated into Neovim that gamifies learning vim keybindings and plugin workflows. Defend against waves of bugs by completing kata-style editing challenges!

## Features

- **Tower Defense Gameplay**: Place towers, defeat enemies, survive waves
- **Challenge-Based Economy**: Gold primarily comes from completing vim kata challenges
- **Speed Bonuses**: Complete challenges faster for up to 2x gold multiplier
- **Plugin-Aware Challenges**: Challenges adapt to your installed plugins (Telescope, nvim-surround, etc.)
- **Keymap Hints**: See your actual keybindings in challenge hints, not just defaults
- **Difficulty Presets**: Easy (50% mob gold), Normal (25%), Hard (0% - challenges only)
- **Beautiful TUI**: Smooth 60fps rendering with emoji support
- **30+ Challenges**: Covering movement, text objects, LSP, search/replace, and refactoring

## Requirements

- Neovim 0.11+
- Go 1.21+ (for building the game binary)
- A terminal with Unicode/emoji support (iTerm2, Alacritty, Kitty, Windows Terminal)

## Installation

### Using lazy.nvim

```lua
{
  "yourusername/keyforge.nvim",
  dependencies = { "nvim-lua/plenary.nvim" },
  build = "make build",
  config = function()
    require("keyforge").setup({
      keybind = "<leader>kf",
    })
  end
}
```

### Using packer.nvim

```lua
use {
  "yourusername/keyforge.nvim",
  requires = { "nvim-lua/plenary.nvim" },
  run = "make build",
  config = function()
    require("keyforge").setup()
  end
}
```

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/keyforge.nvim ~/.local/share/nvim/lazy/keyforge.nvim

# Build the game binary
cd ~/.local/share/nvim/lazy/keyforge.nvim
make build

# Add to your init.lua
# require("keyforge").setup()
```

## Usage

### Starting the Game

```vim
" Using the command
:Keyforge

" Using the default keybind
<leader>kf
```

### Controls

#### In-Game (Tower Defense)

| Key | Action |
|-----|--------|
| `h/j/k/l` or Arrow Keys | Move cursor |
| `1`, `2`, `3` | Select tower type |
| `Space` or `Enter` | Place tower |
| `u` | Upgrade tower (when cursor on tower) |
| `p` | Pause/Resume game |
| `q` | Quit game |
| `r` | Restart (on game over) |

#### Challenge Controls

| Key | Action |
|-----|--------|
| `<leader>kn` | Start next challenge |
| `<leader>kc` | Complete current challenge (validate) |
| `<leader>ks` | Skip current challenge |

### Tower Types

| Tower | Cost | Category | Description |
|-------|------|----------|-------------|
| Arrow | 50g | Movement | Basic tower, triggers movement challenges |
| LSP | 100g | LSP Navigation | Long range, triggers go-to-definition challenges |
| Refactor | 150g | Text Objects | Area damage, triggers text manipulation challenges |

### Challenge Categories

1. **Movement Mastery**: `$`, `^`, `w`, `f`, `G`, `%`, etc.
2. **Text Objects**: `ciw`, `da(`, `yi"`, etc.
3. **LSP Navigation**: `gd`, `gr`, `K`, rename
4. **Search & Replace**: `/`, `:s`, `:g`
5. **Refactoring**: Extract function, inline variable

### Plugin-Aware Challenges

Keyforge detects your installed plugins and shows challenges tailored to your setup:

| Plugin | Challenges Unlocked |
|--------|---------------------|
| Telescope | Fuzzy find files, live grep, buffer search |
| nvim-surround / mini.surround | Change quotes, add/remove surroundings |
| fugitive / gitsigns | Git status, stage hunks, blame |
| nvim-tree / neo-tree | File navigation, create/delete files |

Challenge hints show **your actual keybindings**, not just defaults!

## Economy System

Keyforge uses a challenge-based economy where completing vim kata challenges is the primary source of gold:

### Gold Sources

| Source | Normal Difficulty |
|--------|-------------------|
| Mob Kills | 25% of base gold value |
| Wave Completion | 75% of bonus |
| **Challenge Completion** | **100% (primary source)** |

### Speed Bonus

Complete challenges faster for bonus gold:

| Speed | Bonus Multiplier |
|-------|------------------|
| Slower than par | 1.0x |
| At par | 1.0x |
| 2x faster | 1.5x |
| 4x+ faster | 2.0x (max) |

### Difficulty Presets

| Difficulty | Mob Gold | Description |
|------------|----------|-------------|
| Easy | 50% | Good for learning, some buffer from mobs |
| **Normal** | **25%** | Balanced, challenges are main income |
| Hard | 0% | Pure challenge mode, no mob gold |

## Configuration

```lua
require("keyforge").setup({
  -- Keybind to launch game (set to "" to disable)
  keybind = "<leader>kf",

  -- Challenge keybinds
  keybind_next_challenge = "<leader>kn",
  keybind_complete = "<leader>kc",
  keybind_skip = "<leader>ks",

  -- Difficulty level: "easy", "normal", "hard"
  difficulty = "normal",

  -- Use Nerd Font icons (set false for ASCII fallback)
  use_nerd_fonts = true,

  -- Starting resources
  starting_gold = 200,
  starting_health = 100,

  -- Auto-build binary on first run
  auto_build = true,
})
```

## Commands

| Command | Description |
|---------|-------------|
| `:Keyforge` | Start the game |
| `:KeyforgeStop` | Stop the running game |
| `:KeyforgeBuild` | Rebuild the game binary |
| `:KeyforgeComplete` | Complete current challenge |

## Development

### Building

```bash
# Build the game binary
make build

# Run the standalone game (for testing)
make run

# Run tests
make test
```

### Project Structure

```
keyforge.nvim/
├── plugin/keyforge.lua       # Auto-loader
├── lua/keyforge/
│   ├── init.lua              # Main module, setup
│   ├── rpc.lua               # JSON-RPC communication
│   ├── challenges.lua        # Challenge validation
│   └── ui.lua                # Buffer management
├── game/                     # Go game engine
│   ├── cmd/keyforge/         # Entry point
│   └── internal/
│       ├── engine/           # Game logic
│       ├── entities/         # Towers, enemies, effects
│       ├── ui/               # Bubbletea TUI
│       └── nvim/             # RPC protocol
└── tests/                    # Lua tests (plenary)
```

### Running Tests

```bash
# Go tests
cd game && go test ./...

# Lua tests (requires plenary.nvim)
nvim --headless -c "PlenaryBustedDirectory tests/ {minimal_init = 'tests/minimal_init.lua'}"
```

## Custom Challenges

Create custom challenges in `~/.config/nvim/keyforge-challenges/`:

```yaml
# my-challenge.yaml
id: my_custom_challenge
name: "My Challenge"
category: movement
difficulty: 2
description: "Navigate to the target using vim motions"

initial_buffer: |
  function example() {
    return 42;
  }

validation_type: exact_match
expected_buffer: |
  function example() {
    return 0;
  }

par_keystrokes: 5
gold_base: 50
```

### Validation Types

- `exact_match`: Buffer must exactly match `expected_buffer`
- `contains`: Buffer must contain `expected_content`
- `function_exists`: A function named `function_name` must exist
- `pattern`: Buffer must match Lua pattern
- `different`: Buffer must be different from initial

## Contributing

Contributions are welcome! Please read the contributing guidelines first.

## License

MIT License

## Credits

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Inspired by typing games and vim training tools
