# level-system Specification Delta

## Related Specs
- `game-engine` (modifies Level Definition System, Level Registry, Level Preview)

## MODIFIED Requirements

### Requirement: Level Definition System (game-engine)
The game engine SHALL support 10 levels, each defining a unique path, wave configuration, allowed enemies, and difficulty tier.

#### Scenario: Level count
- **WHEN** the level registry is initialized
- **THEN** it SHALL contain exactly 10 levels
- **AND** levels SHALL be numbered 1 through 10

#### Scenario: Level difficulty progression
- **WHEN** levels are listed in the browser
- **THEN** they SHALL be ordered by difficulty (1 = easiest, 10 = hardest)
- **AND** difficulty indicators SHALL show: beginner (1-3), intermediate (4-6), advanced (7-10)

#### Scenario: Level path uniqueness
- **WHEN** a level is loaded
- **THEN** its path SHALL be unique from all other levels
- **AND** path length SHALL increase with level number (approx 15-60+ cells)

#### Scenario: Level enemy pool
- **WHEN** a level is loaded
- **THEN** only enemies from that level's pool SHALL spawn
- **AND** earlier levels SHALL have easier enemies (Mite, Bug)
- **AND** later levels SHALL include harder enemies (Daemon, Boss)

#### Scenario: Load level definition
- **WHEN** a level is selected from the start screen
- **THEN** the game SHALL load the level's path, grid size, and wave configuration
- **AND** only the level's allowed enemies SHALL spawn
- **AND** all tower types SHALL be available for placement

### Requirement: Level Registry (game-engine)
The level registry SHALL provide access to all 10 levels with metadata for the browser.

#### Scenario: Level registry contents
- **WHEN** the level browser is displayed
- **THEN** all 10 levels SHALL be listed
- **AND** each level SHALL show: number, name, difficulty, wave count, path preview

#### Scenario: Level metadata
- **WHEN** a level is retrieved from the registry
- **THEN** it SHALL include: ID, Name, Description, GridWidth, GridHeight, Path, TotalWaves, EnemyTypes, Difficulty

### Requirement: Level Preview (game-engine)
The level preview SHALL show path layout and enemy roster for selection.

#### Scenario: Path preview
- **WHEN** a level is highlighted in the browser
- **THEN** a mini-grid SHALL display the level's path
- **AND** longer paths SHALL be visually distinguishable from shorter ones

#### Scenario: Enemy preview
- **WHEN** a level is highlighted
- **THEN** the enemy types for that level SHALL be displayed
- **AND** each enemy type SHALL show its icon, name, and health indicator

## ADDED Requirements

### Requirement: Level Path Definitions
The system SHALL define 10 unique paths with progressive complexity.

#### Scenario: Level 1 path (Straight)
- **WHEN** level 1 is loaded
- **THEN** the path SHALL be approximately straight with minimal turns
- **AND** path length SHALL be approximately 15 cells

#### Scenario: Level 5 path (Classic)
- **WHEN** level 5 is loaded
- **THEN** the path SHALL match the original S-shaped Classic path
- **AND** path length SHALL be approximately 35 cells

#### Scenario: Level 10 path (Ultimate)
- **WHEN** level 10 is loaded
- **THEN** the path SHALL be the most complex with many turns
- **AND** path length SHALL be approximately 60+ cells
- **AND** grid size MAY be larger (28x14) to accommodate complexity

#### Scenario: Path connectivity
- **WHEN** any level path is loaded
- **THEN** all waypoints SHALL be connected (no gaps)
- **AND** each waypoint SHALL be adjacent to the next (Manhattan distance = 1)
