## MODIFIED Requirements

### Requirement: Game Launching
The plugin SHALL spawn the Go game binary with nvim-mode enabled to allow bidirectional RPC communication for challenges and game state.

#### Scenario: Start game with keybind
- **WHEN** the user presses `<leader>kf` (default)
- **THEN** a terminal split SHALL open
- **AND** the keyforge binary SHALL be executed with `--nvim-mode` flag
- **AND** RPC communication SHALL be established via stdin/stdout

#### Scenario: Auto-build on first run
- **WHEN** the game binary does not exist
- **THEN** the plugin SHALL attempt to build it using `make build`
- **AND** a progress indicator SHALL be shown
- **AND** the game SHALL start after successful build

#### Scenario: Binary not found
- **WHEN** the build fails or Go is not installed
- **THEN** an error message SHALL be displayed
- **AND** instructions for manual build SHALL be provided

## MODIFIED Requirements

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

## ADDED Requirements

### Requirement: Game End Handling
The plugin SHALL detect game over and victory states and display appropriate UI in Neovim.

#### Scenario: Game over notification
- **WHEN** the game sends a game_over notification
- **THEN** the plugin SHALL display a game over message
- **AND** the message SHALL include final stats (wave reached, gold earned)
- **AND** the user SHALL be offered restart or quit options

#### Scenario: Victory notification
- **WHEN** the game sends a victory notification
- **THEN** the plugin SHALL display a victory message
- **AND** the message SHALL include final stats (health remaining, gold earned, towers built)
- **AND** the user SHALL be offered play again or quit options

#### Scenario: Restart game
- **WHEN** the user chooses to restart after game over or victory
- **THEN** the plugin SHALL send a restart command to the game
- **AND** the game SHALL reset to initial state

## ADDED Requirements

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
