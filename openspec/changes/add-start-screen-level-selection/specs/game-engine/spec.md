## ADDED Requirements

### Requirement: Level Definition System
The game engine SHALL support multiple levels, each defining a unique path, wave configuration, allowed towers, and enemy types.

#### Scenario: Load level definition
- **WHEN** a level is selected from the start screen
- **THEN** the game SHALL load the level's path, grid size, and wave configuration
- **AND** only the level's allowed towers SHALL be available for placement
- **AND** only the level's enemy types SHALL spawn

#### Scenario: Classic level as default
- **WHEN** the game starts
- **THEN** the "Classic" level SHALL be pre-selected
- **AND** it SHALL use the existing S-shaped path and 10-wave progression
- **AND** all current tower and enemy types SHALL be available

#### Scenario: Level registry
- **WHEN** the level browser is displayed
- **THEN** all registered levels SHALL be listed
- **AND** each level SHALL show its name, difficulty, and wave count

### Requirement: Game Settings Configuration
The game engine SHALL support runtime configuration of difficulty, game speed, starting gold, and starting health.

#### Scenario: Difficulty setting
- **WHEN** difficulty is set to "easy", "normal", or "hard"
- **THEN** the economy multipliers SHALL be applied per EconomyConfig
- **AND** the setting SHALL persist for the entire game session

#### Scenario: Game speed setting
- **WHEN** game speed is set (0.5x, 1x, 1.5x, 2x)
- **THEN** all time-based updates SHALL be scaled by the speed multiplier
- **AND** enemy movement, tower cooldowns, and wave timers SHALL all scale proportionally

#### Scenario: Starting resources
- **WHEN** starting gold is configured (100-500)
- **THEN** the game SHALL begin with that gold amount
- **WHEN** starting health is configured (50-200)
- **THEN** the game SHALL begin with that health amount
- **AND** max health SHALL equal starting health

#### Scenario: Settings from command-line
- **WHEN** the game binary receives --difficulty, --game-speed, --starting-gold, --starting-health flags
- **THEN** those values SHALL become the default settings in the settings menu
- **AND** the user MAY still modify them before starting

### Requirement: Start Screen State Machine
The game engine SHALL implement a start screen flow with level selection and settings configuration before gameplay.

#### Scenario: Initial state
- **WHEN** the game binary starts
- **THEN** it SHALL enter the start screen state
- **AND** the level browser SHALL be displayed

#### Scenario: Level selection flow
- **WHEN** a level is highlighted and Enter is pressed
- **THEN** the game SHALL transition to the settings menu
- **AND** the selected level SHALL be stored for game initialization

#### Scenario: Settings confirmation
- **WHEN** "Start Game" is selected in the settings menu
- **THEN** the game SHALL initialize with the selected level and settings
- **AND** it SHALL transition to the playing state

#### Scenario: Back navigation
- **WHEN** ESC is pressed in the settings menu
- **THEN** the game SHALL return to level selection
- **WHEN** ESC is pressed in level selection
- **THEN** the game SHALL quit (or show quit confirmation)

### Requirement: Level Preview
The game engine SHALL display a preview of the selected level showing path, enemies, and towers.

#### Scenario: Path preview
- **WHEN** a level is highlighted in the browser
- **THEN** a mini-grid SHALL display the level's path
- **AND** the path SHALL be visually distinguished from empty cells

#### Scenario: Enemy preview
- **WHEN** a level is highlighted
- **THEN** the enemy types that appear in that level SHALL be listed
- **AND** each enemy type SHALL show its icon and name

#### Scenario: Tower preview
- **WHEN** a level is highlighted
- **THEN** the towers available for that level SHALL be listed
- **AND** each tower type SHALL show its icon and name

## MODIFIED Requirements

### Requirement: Game State Machine
The game engine SHALL maintain distinct states for gameplay phases including start screen states, a challenge-waiting state for nvim mode.

#### Scenario: State transitions
- **WHEN** the game starts
- **THEN** it SHALL begin in the start screen state
- **AND** it SHALL transition to level select when ready
- **WHEN** the game is playing
- **THEN** it SHALL transition to paused when pause is triggered
- **AND** it SHALL transition to challenge-waiting when challenge starts in nvim mode
- **AND** it SHALL transition to game_over when health reaches zero
- **AND** it SHALL transition to victory when all waves are completed

#### Scenario: Game over
- **WHEN** health reaches zero
- **THEN** the game SHALL enter the game_over state
- **AND** all game updates SHALL stop
- **AND** the game SHALL send a game_over notification to Neovim (in nvim mode)

#### Scenario: Victory
- **WHEN** the final wave is completed and all enemies are defeated
- **THEN** the game SHALL enter the victory state
- **AND** a victory notification SHALL be sent to Neovim (in nvim mode)

#### Scenario: Return to start screen
- **WHEN** the game ends (victory or game over)
- **AND** the user presses a key to continue
- **THEN** the game SHALL return to the start screen
- **AND** the previous level and settings SHALL be pre-selected
