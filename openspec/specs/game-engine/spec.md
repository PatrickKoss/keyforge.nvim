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
The game engine SHALL manage game state through distinct phases: menu, playing, paused, challenge-active, wave-complete, game-over, and victory.

#### Scenario: State transitions
- **WHEN** the game is in "playing" state and a challenge is requested
- **THEN** the state SHALL transition to "challenge-active"
- **AND** game updates SHALL pause until challenge completes

#### Scenario: Game over
- **WHEN** player health reaches zero
- **THEN** the state SHALL transition to "game-over"
- **AND** final statistics SHALL be displayed

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

