# game-engine Specification

## Purpose
TBD - created by archiving change add-keyforge-core. Update Purpose after archive.
## Requirements
### Requirement: Game Loop
The game engine SHALL implement a bubbletea-based main loop that processes input, updates game state, and renders the UI at a target of 60 frames per second.

#### Scenario: Smooth rendering
- **WHEN** the game is running
- **THEN** the frame rate SHALL remain stable at approximately 60fps
- **AND** input lag SHALL be imperceptible (<16ms response)

#### Scenario: Game pause
- **WHEN** the user presses the pause key
- **THEN** the game state SHALL freeze
- **AND** a pause overlay SHALL be displayed
- **AND** pressing pause again SHALL resume the game

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

### Requirement: Grid Rendering
The game engine SHALL render a playable grid using Unicode box-drawing characters and optional Nerd Font icons for entities.

#### Scenario: Grid display
- **WHEN** the game view is rendered
- **THEN** a bordered game grid SHALL be displayed
- **AND** towers, enemies, and projectiles SHALL be visible at their positions
- **AND** the path SHALL be clearly marked

#### Scenario: Fallback rendering
- **WHEN** the terminal does not support Nerd Fonts
- **THEN** ASCII fallback characters SHALL be used for all entities

### Requirement: Physics System
The game engine SHALL implement movement and collision detection for all moving entities.

#### Scenario: Enemy movement
- **WHEN** an enemy is active on the grid
- **THEN** it SHALL move along the predefined path at its configured speed
- **AND** position updates SHALL be smooth (interpolated between cells)

#### Scenario: Projectile collision
- **WHEN** a projectile occupies the same cell as an enemy
- **THEN** the collision SHALL be detected
- **AND** damage SHALL be applied to the enemy
- **AND** the projectile SHALL be removed

### Requirement: HUD Display
The game engine SHALL display a heads-up display showing current wave, gold, health, and available tower options.

#### Scenario: Resource display
- **WHEN** the game is active
- **THEN** current gold amount SHALL be visible
- **AND** current health SHALL be visible with a visual bar
- **AND** current wave number and total waves SHALL be displayed

#### Scenario: Tower shop
- **WHEN** the player is not in a challenge
- **THEN** available tower types SHALL be displayed with their costs
- **AND** towers the player cannot afford SHALL be visually dimmed

### Requirement: Game State Notifications
The game engine SHALL send state change notifications to Neovim when running in nvim mode.

#### Scenario: Game over notification
- **WHEN** the game enters game_over state in nvim mode
- **THEN** a game_over RPC notification SHALL be sent
- **AND** the notification SHALL include wave_reached, final_gold, towers_built

#### Scenario: Victory notification
- **WHEN** the game enters victory state in nvim mode
- **THEN** a victory RPC notification SHALL be sent
- **AND** the notification SHALL include final_health, final_gold, towers_built

#### Scenario: Restart command
- **WHEN** the game receives a restart RPC command
- **THEN** the game SHALL reset to initial state
- **AND** the game SHALL enter playing state
- **AND** the game SHALL notify Neovim of the state change

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

