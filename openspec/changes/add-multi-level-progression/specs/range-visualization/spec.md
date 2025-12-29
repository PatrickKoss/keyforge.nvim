# range-visualization Specification Delta

## Related Specs
- `tower-defense` (modifies Tower Placement)
- `game-engine` (modifies Grid Rendering)

## ADDED Requirements

### Requirement: Tower Range Visualization
The system SHALL display tower attack range to help players make informed placement decisions.

#### Scenario: Range display during placement
- **WHEN** the player is in tower placement mode
- **AND** a tower type is selected
- **THEN** the selected tower's range SHALL be displayed as a circular overlay
- **AND** the overlay SHALL be centered on the cursor position
- **AND** cells within range SHALL be visually highlighted

#### Scenario: Range display on hover
- **WHEN** the cursor is positioned over an existing tower
- **AND** the game is in playing state (not challenge)
- **THEN** that tower's current range (including upgrades) SHALL be displayed
- **AND** the overlay SHALL be centered on the tower position

#### Scenario: Range overlay styling
- **WHEN** a range overlay is displayed
- **THEN** it SHALL use a semi-transparent or dotted visual style
- **AND** it SHALL NOT obscure enemies, projectiles, or the path
- **AND** it SHALL be distinguishable from the cursor and other UI elements

#### Scenario: Range calculation
- **WHEN** range is calculated for display
- **THEN** it SHALL use the tower's effective range (base + upgrade bonuses)
- **AND** cells SHALL be considered "in range" if their center is within the range radius

#### Scenario: Range display absence
- **WHEN** the cursor is not on a tower
- **AND** the player is not in placement mode
- **THEN** no range overlay SHALL be displayed

### Requirement: Range-Based Tower Selection
The system SHALL help players understand tower range differences when selecting towers.

#### Scenario: Shop range indicator
- **WHEN** the tower shop is displayed
- **THEN** each tower's range value MAY be shown alongside cost
- **AND** players SHALL be able to compare ranges before selection

#### Scenario: Placement validity with range context
- **WHEN** placing a tower
- **THEN** the range overlay SHALL help visualize path coverage
- **AND** players can assess whether the position covers desired path segments
