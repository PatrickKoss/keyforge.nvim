# neovim-integration Specification

## Purpose
TBD - created by archiving change add-keyforge-core. Update Purpose after archive.
## Requirements
### Requirement: Plugin Loading
The plugin SHALL auto-load when Neovim starts and register the keyforge module without impacting startup time by more than 100ms.

#### Scenario: Lazy initialization
- **WHEN** Neovim starts with keyforge installed
- **THEN** only minimal setup code SHALL execute
- **AND** full initialization SHALL be deferred until first use

#### Scenario: Plugin availability
- **WHEN** the user calls `require("keyforge")`
- **THEN** the module SHALL be available
- **AND** the `setup()` function SHALL accept configuration options

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

### Requirement: RPC Communication
The plugin SHALL establish bidirectional JSON-RPC communication with the game process over stdin/stdout, handling challenge requests, game state updates, and end-game notifications.

#### Scenario: Receive challenge request
- **WHEN** the game sends a challenge request message
- **THEN** the plugin SHALL parse the JSON-RPC request
- **AND** the challenge handler SHALL create a temp file and open it in a new tab

#### Scenario: Send challenge result
- **WHEN** a challenge is completed
- **THEN** the plugin SHALL send a JSON-RPC response to the game
- **AND** the response SHALL include success status, keystroke count, time, and gold earned

#### Scenario: Receive game state update
- **WHEN** the game sends a state update notification
- **THEN** the plugin SHALL update internal state tracking
- **AND** the plugin SHALL handle pause/resume state changes

#### Scenario: Handle connection loss
- **WHEN** the game process terminates unexpectedly
- **THEN** the plugin SHALL detect the disconnection
- **AND** clean up any active challenge buffers
- **AND** display an error message to the user

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

### Requirement: Session Cleanup
The plugin SHALL properly clean up resources when the game ends or Neovim exits.

#### Scenario: Game exit
- **WHEN** the user quits the game
- **THEN** the game process SHALL be terminated
- **AND** the terminal buffer SHALL be closed or left for user inspection
- **AND** no orphan processes SHALL remain

#### Scenario: Neovim exit
- **WHEN** Neovim is closed while game is running
- **THEN** the game process SHALL be terminated gracefully
- **AND** no zombie processes SHALL remain

### Requirement: Tab Management
The plugin SHALL manage Neovim tabs to provide seamless transitions between game and challenge views.

#### Scenario: Track game tab
- **WHEN** the game is launched
- **THEN** the plugin SHALL store the game tab ID
- **AND** the plugin SHALL restore focus to game tab after challenge completion

#### Scenario: Challenge tab creation
- **WHEN** a challenge is started
- **THEN** a new tab SHALL be created for the challenge buffer
- **AND** the original game tab SHALL remain open but unfocused

#### Scenario: Return to game tab
- **WHEN** a challenge is completed or cancelled
- **THEN** the challenge tab SHALL be closed
- **AND** focus SHALL return to the game tab
- **AND** the terminal SHALL be re-rendered to show current game state

