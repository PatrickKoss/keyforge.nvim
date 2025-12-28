# tower-defense Spec Delta

## MODIFIED Requirements

### Requirement: Resource Economy
The system SHALL manage gold currency with challenges as the primary source and reduced rewards from mob kills.

#### Scenario: Reduced mob gold
- **WHEN** an enemy is killed by tower damage
- **THEN** gold SHALL be awarded at 25% of the base value
- **AND** Bug kill SHALL award 4g (was 15g)
- **AND** Gremlin kill SHALL award 6g (was 25g)
- **AND** Daemon kill SHALL award 10g (was 40g)
- **AND** Boss kill SHALL award 25g (was 100g)

#### Scenario: Reduced wave bonus
- **WHEN** a wave is completed
- **THEN** the wave bonus gold SHALL be 50% of the original value
- **AND** the bonus SHALL still scale with wave number

#### Scenario: Challenge gold as primary income
- **WHEN** a challenge is completed successfully
- **THEN** gold SHALL be awarded based on base_gold, difficulty, efficiency, and speed bonus
- **AND** challenge rewards SHALL be the primary method of earning gold

#### Scenario: Starting resources (unchanged)
- **WHEN** a new game begins
- **THEN** the player SHALL start with configured starting gold (default 200)
- **AND** starting health SHALL be configured (default 100)

#### Scenario: Economy balance
- **GIVEN** a player who completes challenges at average efficiency
- **THEN** the total gold earned SHALL be sufficient to build towers for wave progression
- **AND** faster challenge completion SHALL provide meaningful economic advantage

## ADDED Requirements

### Requirement: Economy Configuration
The system SHALL support configurable economy multipliers for difficulty tuning.

#### Scenario: Easy difficulty
- **WHEN** the game is set to easy difficulty
- **THEN** mob gold multiplier SHALL be 0.50 (50%)
- **AND** challenge rewards SHALL remain at 100%

#### Scenario: Normal difficulty
- **WHEN** the game is set to normal difficulty
- **THEN** mob gold multiplier SHALL be 0.25 (25%)
- **AND** challenge rewards SHALL remain at 100%

#### Scenario: Hard difficulty
- **WHEN** the game is set to hard difficulty
- **THEN** mob gold multiplier SHALL be 0.00 (0%)
- **AND** all gold MUST come from challenges

#### Scenario: Custom economy config
- **WHEN** the user provides custom economy configuration
- **THEN** the configured multipliers SHALL override defaults
- **AND** invalid values SHALL fall back to normal difficulty
