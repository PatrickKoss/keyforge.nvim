# tower-defense Specification

## Purpose
TBD - created by archiving change add-keyforge-core. Update Purpose after archive.
## Requirements
### Requirement: Tower Types
The system SHALL provide multiple tower types, each mapped to a vim skill category and with unique attributes (damage, range, attack speed, cost).

#### Scenario: Arrow Tower
- **WHEN** the player builds an Arrow Tower
- **THEN** it SHALL deal moderate damage at medium range
- **AND** it SHALL trigger "movement" category challenges
- **AND** its cost SHALL be the lowest tier (50 gold base)

#### Scenario: LSP Tower
- **WHEN** the player builds an LSP Tower
- **THEN** it SHALL deal high damage at long range
- **AND** it SHALL trigger "lsp-navigation" category challenges
- **AND** its cost SHALL be mid-tier (100 gold base)

#### Scenario: Refactor Tower
- **WHEN** the player builds a Refactor Tower
- **THEN** it SHALL deal area damage
- **AND** it SHALL trigger "text-objects" category challenges
- **AND** its cost SHALL be higher tier (150 gold base)

### Requirement: Tower Placement
The system SHALL allow players to place towers on valid grid positions using cursor movement and a place action.

#### Scenario: Valid placement
- **WHEN** the player selects an empty, non-path cell
- **AND** has sufficient gold
- **THEN** a tower MAY be placed at that position
- **AND** gold SHALL be deducted

#### Scenario: Invalid placement
- **WHEN** the player attempts to place a tower on the path or an occupied cell
- **THEN** placement SHALL be rejected
- **AND** visual feedback SHALL indicate the error

#### Scenario: Insufficient funds
- **WHEN** the player attempts to place a tower without sufficient gold
- **THEN** placement SHALL be rejected
- **AND** a "not enough gold" message SHALL be displayed

### Requirement: Tower Targeting
Towers SHALL automatically target enemies within range, prioritizing based on configurable strategies.

#### Scenario: Auto-targeting
- **WHEN** an enemy enters a tower's range
- **THEN** the tower SHALL begin attacking
- **AND** projectiles SHALL be fired at the target

#### Scenario: Target priority
- **WHEN** multiple enemies are in range
- **THEN** the tower SHALL target based on strategy (first, strongest, weakest, closest)
- **AND** default strategy SHALL be "first" (furthest along path)

### Requirement: Tower Upgrades
The system SHALL allow players to upgrade placed towers to increase their effectiveness.

#### Scenario: Upgrade available
- **WHEN** a tower has not reached max level
- **AND** the player has sufficient gold
- **THEN** the upgrade option SHALL be available

#### Scenario: Apply upgrade
- **WHEN** the player upgrades a tower
- **THEN** tower stats (damage, range, or speed) SHALL increase
- **AND** gold SHALL be deducted
- **AND** visual appearance MAY change to indicate level

### Requirement: Enemy Types
The system SHALL provide multiple enemy types with varying health, speed, and special abilities.

#### Scenario: Bug enemy (basic)
- **WHEN** a Bug spawns
- **THEN** it SHALL have low health (10 base)
- **AND** medium speed
- **AND** no special abilities

#### Scenario: Gremlin enemy (fast)
- **WHEN** a Gremlin spawns
- **THEN** it SHALL have medium health (25 base)
- **AND** high speed
- **AND** no special abilities

#### Scenario: Daemon enemy (tank)
- **WHEN** a Daemon spawns
- **THEN** it SHALL have high health (100 base)
- **AND** slow speed
- **AND** no special abilities

#### Scenario: Boss enemy
- **WHEN** a Boss spawns
- **THEN** it SHALL have very high health (500 base)
- **AND** slow speed
- **AND** may have special abilities (regeneration, armor)

### Requirement: Enemy Pathfinding
Enemies SHALL follow a predefined path from spawn point to exit point.

#### Scenario: Path following
- **WHEN** an enemy is spawned
- **THEN** it SHALL follow the path waypoints in order
- **AND** movement SHALL be smooth between waypoints

#### Scenario: Reach exit
- **WHEN** an enemy reaches the path exit
- **THEN** the player SHALL lose health equal to enemy's remaining health percentage
- **AND** the enemy SHALL be removed from the grid

### Requirement: Wave System
The system SHALL spawn enemies in waves with increasing difficulty.

#### Scenario: Wave start
- **WHEN** a wave begins
- **THEN** enemies SHALL spawn according to the wave definition
- **AND** spawn timing SHALL follow configured intervals

#### Scenario: Wave completion
- **WHEN** all enemies in a wave are defeated
- **THEN** the wave SHALL be marked complete
- **AND** the player SHALL have time before the next wave
- **AND** bonus gold MAY be awarded

#### Scenario: Final wave
- **WHEN** the final wave is completed
- **THEN** the game SHALL transition to victory state
- **AND** final statistics SHALL be displayed

### Requirement: Resource Economy
The system SHALL manage gold currency earned from challenges and spent on towers.

#### Scenario: Earn gold
- **WHEN** a challenge is completed successfully
- **THEN** gold SHALL be awarded based on efficiency
- **AND** the HUD SHALL update to show new total

#### Scenario: Spend gold
- **WHEN** gold is spent on a tower or upgrade
- **THEN** the gold total SHALL decrease
- **AND** the transaction SHALL only succeed if sufficient gold exists

#### Scenario: Starting resources
- **WHEN** a new game begins
- **THEN** the player SHALL start with configured starting gold (default 200)
- **AND** starting health SHALL be configured (default 100)

### Requirement: Health System
The system SHALL track player health that decreases when enemies reach the exit.

#### Scenario: Take damage
- **WHEN** an enemy reaches the exit
- **THEN** player health SHALL decrease
- **AND** the health bar SHALL update visually

#### Scenario: Game over
- **WHEN** player health reaches zero
- **THEN** the game SHALL end
- **AND** game-over screen SHALL display statistics

### Requirement: Projectile System
The system SHALL manage projectiles fired by towers with visual rendering and collision detection.

#### Scenario: Fire projectile
- **WHEN** a tower attacks an enemy
- **THEN** a projectile SHALL be created
- **AND** it SHALL travel toward the target position

#### Scenario: Projectile hit
- **WHEN** a projectile collides with an enemy
- **THEN** damage SHALL be applied
- **AND** the projectile SHALL be removed
- **AND** a hit effect MAY be displayed

#### Scenario: Projectile miss
- **WHEN** a projectile reaches its target position but the enemy has moved
- **THEN** the projectile SHALL be removed without damage
- **OR** homing projectiles SHALL continue to track (tower-dependent)

