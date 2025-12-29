## MODIFIED Requirements

### Requirement: Configuration
The plugin SHALL accept user configuration for keybinds, difficulty, visual settings, gameplay options, and game settings defaults.

#### Scenario: Custom keybind
- **WHEN** the user configures `keybind = "<leader>g"`
- **THEN** that keybind SHALL launch the game instead of the default

#### Scenario: Difficulty setting
- **WHEN** the user sets `difficulty = "hard"`
- **THEN** the setting SHALL be passed to the game process as a command-line flag
- **AND** it SHALL appear as the default in the game's settings menu

#### Scenario: Game speed setting
- **WHEN** the user sets `game_speed = 1.5`
- **THEN** the setting SHALL be passed to the game process as a command-line flag
- **AND** it SHALL appear as the default in the game's settings menu

#### Scenario: Starting resources setting
- **WHEN** the user sets `starting_gold = 300` or `starting_health = 150`
- **THEN** those values SHALL be passed to the game process as command-line flags
- **AND** they SHALL appear as defaults in the game's settings menu

#### Scenario: Default configuration
- **WHEN** no configuration is provided
- **THEN** sensible defaults SHALL be applied (difficulty=normal, game_speed=1.0, starting_gold=200, starting_health=100)
- **AND** the game SHALL be playable without any setup

### Requirement: Game Launching
The plugin SHALL spawn the Go game binary in a terminal split when the user invokes the start command or keybind, passing configuration as command-line flags.

#### Scenario: Start game with keybind
- **WHEN** the user presses `<leader>kf` (default)
- **THEN** a terminal split SHALL open
- **AND** the keyforge binary SHALL be executed with config flags
- **AND** RPC communication SHALL be established

#### Scenario: Pass config to binary
- **WHEN** the game is launched
- **THEN** the binary SHALL receive --difficulty, --game-speed, --starting-gold, --starting-health flags
- **AND** the values SHALL match the user's nvim configuration

#### Scenario: Auto-build on first run
- **WHEN** the game binary does not exist
- **THEN** the plugin SHALL attempt to build it using `make build`
- **AND** a progress indicator SHALL be shown
- **AND** the game SHALL start after successful build

#### Scenario: Binary not found
- **WHEN** the build fails or Go is not installed
- **THEN** an error message SHALL be displayed
- **AND** instructions for manual build SHALL be provided
