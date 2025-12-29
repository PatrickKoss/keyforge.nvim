# enemy-variety Specification Delta

## Related Specs
- `tower-defense` (modifies Enemy Types, Resource Economy)

## MODIFIED Requirements

### Requirement: Enemy Types (tower-defense)
The system SHALL provide 7 enemy types with distinct health, speed, and gold profiles to enable varied combat encounters.

#### Scenario: Mite enemy (fodder)
- **WHEN** a Mite spawns
- **THEN** it SHALL have very low health (5 base)
- **AND** fast speed (2.0 cells/sec)
- **AND** minimal gold value (2)
- **AND** it SHALL appear in levels 1-4

#### Scenario: Bug enemy (basic)
- **WHEN** a Bug spawns
- **THEN** it SHALL have low health (10 base)
- **AND** medium speed (1.5 cells/sec)
- **AND** low gold value (5)
- **AND** it SHALL appear in levels 1-6

#### Scenario: Gremlin enemy (fast)
- **WHEN** a Gremlin spawns
- **THEN** it SHALL have medium health (25 base)
- **AND** high speed (2.5 cells/sec)
- **AND** medium gold value (10)
- **AND** it SHALL appear in levels 3-8

#### Scenario: Crawler enemy (tank)
- **WHEN** a Crawler spawns
- **THEN** it SHALL have high health (40 base)
- **AND** very slow speed (0.6 cells/sec)
- **AND** medium-high gold value (15)
- **AND** it SHALL appear in levels 4-8

#### Scenario: Specter enemy (glass cannon)
- **WHEN** a Specter spawns
- **THEN** it SHALL have low health (15 base)
- **AND** very high speed (3.5 cells/sec)
- **AND** medium gold value (8)
- **AND** it SHALL appear in levels 5-9

#### Scenario: Daemon enemy (heavy tank)
- **WHEN** a Daemon spawns
- **THEN** it SHALL have very high health (100 base)
- **AND** slow speed (0.8 cells/sec)
- **AND** high gold value (25)
- **AND** it SHALL appear in levels 6-10

#### Scenario: Boss enemy
- **WHEN** a Boss spawns
- **THEN** it SHALL have extreme health (500 base)
- **AND** very slow speed (0.5 cells/sec)
- **AND** very high gold value (100)
- **AND** it SHALL only appear in level 10 final wave

### Requirement: Resource Economy (tower-defense)
The system SHALL award gold from enemy kills scaled to enemy difficulty, with challenges remaining the primary gold source.

#### Scenario: Enemy gold scaling
- **WHEN** an enemy is killed
- **THEN** gold SHALL be awarded based on enemy type
- **AND** high-health enemies SHALL give more gold (15-100)
- **AND** low-health enemies SHALL give less gold (2-10)

#### Scenario: Challenge vs kill gold ratio
- **WHEN** a player completes a level
- **THEN** approximately 60-70% of total gold earned SHALL come from challenges
- **AND** approximately 30-40% SHALL come from enemy kills and wave bonuses

## ADDED Requirements

### Requirement: Enemy Level Distribution
The system SHALL restrict which enemy types appear in each level based on difficulty.

#### Scenario: Early level enemies (1-3)
- **WHEN** a wave spawns in levels 1-3
- **THEN** only Mite and Bug enemies SHALL appear
- **AND** no Gremlin, Crawler, Specter, Daemon, or Boss SHALL spawn

#### Scenario: Mid level enemies (4-6)
- **WHEN** a wave spawns in levels 4-6
- **THEN** Bug, Gremlin, Crawler, and Specter enemies MAY appear
- **AND** the specific pool SHALL vary per level per the level definition

#### Scenario: Late level enemies (7-10)
- **WHEN** a wave spawns in levels 7-10
- **THEN** Gremlin, Specter, Daemon enemies MAY appear
- **AND** Boss SHALL only appear in level 10

#### Scenario: Wave enemy count
- **WHEN** a wave is generated for any level
- **THEN** it SHALL contain between 3 and 7 distinct enemy types total
- **AND** enemy count SHALL scale with wave number within the level
