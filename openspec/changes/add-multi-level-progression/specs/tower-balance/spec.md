# tower-balance Specification Delta

## Related Specs
- `tower-defense` (modifies Tower Types, Tower Upgrades)

## MODIFIED Requirements

### Requirement: Tower Types (tower-defense)
The system SHALL provide tower types with cost-proportional stats, creating distinct tactical roles.

#### Scenario: Arrow Tower (fast attacker)
- **WHEN** the player builds an Arrow Tower
- **THEN** it SHALL deal lower damage (8 base) at short range (2.5)
- **AND** it SHALL have fast attack speed (0.8s cooldown)
- **AND** it SHALL trigger "movement" category challenges
- **AND** its cost SHALL be lowest tier (50 gold)

#### Scenario: LSP Tower (sniper)
- **WHEN** the player builds an LSP Tower
- **THEN** it SHALL deal high damage (20 base) at long range (5.0)
- **AND** it SHALL have slow attack speed (1.5s cooldown)
- **AND** it SHALL trigger "lsp-navigation" category challenges
- **AND** its cost SHALL be mid-tier (100 gold)

#### Scenario: Refactor Tower (balanced area)
- **WHEN** the player builds a Refactor Tower
- **THEN** it SHALL deal moderate damage (12 base) at medium range (3.0)
- **AND** it SHALL have medium attack speed (1.0s cooldown)
- **AND** it SHALL deal area damage
- **AND** it SHALL trigger "text-objects" category challenges
- **AND** its cost SHALL be highest tier (150 gold)

### Requirement: Tower Upgrades (tower-defense)
The system SHALL allow tower upgrades that incrementally improve damage, range, and attack speed.

#### Scenario: Upgrade stat bonuses
- **WHEN** a tower is upgraded
- **THEN** damage SHALL increase by approximately 15% (additive to base)
- **AND** range SHALL increase by 0.3 units
- **AND** cooldown SHALL decrease by 10% (multiplicative)

#### Scenario: Upgrade cost scaling
- **WHEN** an upgrade is purchased
- **THEN** the cost SHALL be approximately 60% of the tower's base cost per tier
- **AND** higher tiers MAY cost slightly more

#### Scenario: Upgrade availability
- **WHEN** a tower has not reached max level (2 upgrades)
- **AND** the player has sufficient gold
- **THEN** the upgrade option SHALL be available
- **AND** hovering over the tower SHALL show upgrade cost

#### Scenario: Visual upgrade indicator
- **WHEN** a tower is upgraded
- **THEN** a visual indicator MAY show the upgrade level
- **AND** the tower's effective stats SHALL be displayed on hover
