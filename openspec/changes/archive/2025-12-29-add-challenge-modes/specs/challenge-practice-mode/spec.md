# challenge-practice-mode Specification

## Purpose
Provides an endless challenge practice mode where players can practice vim keybindings continuously without tower defense mechanics, receiving immediate feedback on success or failure.

## ADDED Requirements

### Requirement: Challenge Mode Menu Entry
The system SHALL display a "Challenge Mode" option below the level selection list in the start screen menu.

#### Scenario: Menu visibility
- **WHEN** the player is on the start screen
- **THEN** a "Challenge Mode" option SHALL be visible below the level list
- **AND** it SHALL be selectable using standard navigation keys (j/k)

#### Scenario: Menu selection
- **WHEN** the player selects "Challenge Mode" and presses Enter
- **THEN** the game SHALL transition to challenge mode state
- **AND** the first challenge SHALL be presented immediately

### Requirement: Endless Challenge Loop
The system SHALL present challenges continuously in an endless loop, automatically loading the next challenge after completion.

#### Scenario: Challenge completion triggers next
- **WHEN** a challenge is completed (success or failure)
- **THEN** a brief notification SHALL be shown (1-2 seconds)
- **AND** the next random challenge SHALL be loaded automatically

#### Scenario: Challenge variety
- **WHEN** selecting the next challenge
- **THEN** the system SHALL use the existing ChallengeSelector for variety
- **AND** recently completed challenges SHALL be deprioritized

#### Scenario: Session continuity
- **WHEN** multiple challenges are completed in sequence
- **THEN** a streak counter SHALL be maintained for successful completions
- **AND** the streak SHALL reset on failure

### Requirement: Challenge Mode Notifications
The system SHALL display success/failure notifications that do not overlay the challenge description.

#### Scenario: Success notification
- **WHEN** a challenge is completed successfully
- **THEN** a success notification SHALL appear (e.g., "✓ Success!")
- **AND** it SHALL be positioned in the title bar or top-right corner
- **AND** it SHALL NOT overlay the challenge description or buffer

#### Scenario: Failure notification
- **WHEN** a challenge submission fails validation
- **THEN** a failure notification SHALL appear (e.g., "✗ Try again")
- **AND** it SHALL be positioned in the title bar or top-right corner
- **AND** it SHALL NOT overlay the challenge description or buffer

#### Scenario: Notification auto-dismiss
- **WHEN** a notification is displayed
- **THEN** it SHALL auto-dismiss after 2 seconds
- **OR** it SHALL dismiss when the next challenge loads

### Requirement: Challenge Mode Exit
The system SHALL provide a way to exit challenge mode and return to the main menu.

#### Scenario: Exit via keyboard
- **WHEN** the player presses Escape in the game UI (not during Neovim editing)
- **THEN** the game SHALL return to the start screen (level selection)
- **AND** the challenge mode state SHALL be cleared

#### Scenario: Exit from Neovim buffer
- **WHEN** the player is in a Neovim challenge buffer
- **THEN** a keymap SHALL be available to return to the main menu
- **AND** the current challenge progress SHALL be discarded

### Requirement: Challenge Mode Display
The system SHALL render a dedicated challenge mode screen with mode indicator, streak counter, and challenge content.

#### Scenario: Mode header
- **WHEN** challenge mode is active
- **THEN** a header SHALL display "CHALLENGE MODE"
- **AND** the current streak count SHALL be visible
- **AND** the notification area SHALL be visible in the header

#### Scenario: Challenge content area
- **WHEN** a challenge is displayed
- **THEN** the challenge name and difficulty SHALL be shown
- **AND** the description SHALL be clearly visible
- **AND** the buffer preview/content area SHALL be displayed
- **AND** help text for controls SHALL be shown at the bottom
