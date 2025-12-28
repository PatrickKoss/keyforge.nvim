## ADDED Requirements

### Requirement: Challenge Pause State
The game engine SHALL pause all game updates during challenge completion in nvim mode to allow focused editing.

#### Scenario: Pause on challenge start
- **WHEN** a challenge is started in nvim mode
- **THEN** the game SHALL enter a paused state
- **AND** enemy movement SHALL stop
- **AND** tower firing SHALL stop
- **AND** wave timers SHALL stop

#### Scenario: Display challenge overlay
- **WHEN** the game is paused for a challenge
- **THEN** the game view SHALL display "CHALLENGE IN PROGRESS" indicator
- **AND** the current game state (health, gold, wave) SHALL remain visible
- **AND** the game grid SHALL show frozen enemy positions

#### Scenario: Resume on challenge complete
- **WHEN** a challenge result is received from Neovim
- **THEN** the game SHALL resume from the paused state
- **AND** enemies SHALL continue from their frozen positions
- **AND** wave timers SHALL continue from where they stopped

#### Scenario: Handle timeout during pause
- **WHEN** a challenge times out (no result received within limit)
- **THEN** the game SHALL resume automatically
- **AND** the challenge SHALL be recorded as skipped

## MODIFIED Requirements

### Requirement: Game State Machine
The game engine SHALL maintain distinct states for gameplay phases including a challenge-waiting state for nvim mode.

#### Scenario: State transitions
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

## ADDED Requirements

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
