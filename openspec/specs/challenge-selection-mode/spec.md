# challenge-selection-mode Specification

## Purpose
TBD - created by archiving change add-challenge-modes. Update Purpose after archive.
## Requirements
### Requirement: Challenge Selection Menu Entry
The system SHALL display a "Challenge Selection" option below the level selection list in the start screen menu.

#### Scenario: Menu visibility
- **WHEN** the player is on the start screen
- **THEN** a "Challenge Selection" option SHALL be visible below "Challenge Mode"
- **AND** it SHALL be selectable using standard navigation keys (j/k)

#### Scenario: Menu selection
- **WHEN** the player selects "Challenge Selection" and presses Enter
- **THEN** the game SHALL transition to challenge selection state
- **AND** the challenge list SHALL be displayed

### Requirement: Challenge List Display
The system SHALL display all available challenges in a browsable list organized by category.

#### Scenario: List organization
- **WHEN** the challenge selection screen is displayed
- **THEN** challenges SHALL be grouped by category (movement, text-objects, etc.)
- **AND** category headers SHALL be visible
- **AND** each challenge SHALL show name and difficulty indicator

#### Scenario: List navigation
- **WHEN** the player uses j/k keys in the challenge list
- **THEN** the selection cursor SHALL move between challenges
- **AND** category headers SHALL be skippable (not selectable)

#### Scenario: List scrolling
- **WHEN** the challenge list exceeds the visible area
- **THEN** the list SHALL scroll to keep the selected item visible
- **AND** scroll position SHALL be maintained during navigation

### Requirement: Challenge Preview
The system SHALL display a preview of the currently hovered challenge in a side panel.

#### Scenario: Preview content
- **WHEN** a challenge is hovered (selected but not started)
- **THEN** a preview panel SHALL display:
  - Challenge name
  - Category
  - Difficulty level (★☆☆, ★★☆, ★★★)
  - Full description
  - Truncated initial buffer preview

#### Scenario: Preview updates
- **WHEN** the player navigates to a different challenge
- **THEN** the preview panel SHALL update immediately
- **AND** the previous preview content SHALL be replaced

### Requirement: Challenge Start from Selection
The system SHALL allow the player to start a selected challenge from the list.

#### Scenario: Start challenge
- **WHEN** the player presses Enter on a selected challenge
- **THEN** the game SHALL transition to challenge practice state
- **AND** the selected challenge SHALL be sent to Neovim
- **AND** the player's position in the list SHALL be remembered

#### Scenario: Challenge index tracking
- **WHEN** a challenge is started from selection
- **THEN** the system SHALL track the selected challenge index
- **AND** this index SHALL be used for sequential progression

### Requirement: Sequential Challenge Progression
The system SHALL automatically progress to the next challenge in the list after completion.

#### Scenario: Progress to next challenge
- **WHEN** a challenge is completed (success or failure) from selection mode
- **THEN** a notification SHALL be displayed (success/failure)
- **AND** after the notification, the next challenge in the list SHALL load
- **AND** if at the end of the list, it SHALL wrap to the beginning

#### Scenario: Skip category headers
- **WHEN** progressing to the next challenge
- **THEN** category headers SHALL be skipped
- **AND** the next actual challenge SHALL be loaded

### Requirement: Selection Mode Notifications
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

### Requirement: Return to Selection List
The system SHALL provide a way to return from challenge practice to the selection list.

#### Scenario: Back button during practice
- **WHEN** the player is practicing a challenge from selection mode
- **THEN** a back option SHALL be available to return to the selection list
- **AND** pressing the back keymap SHALL return to the selection list
- **AND** the list position SHALL be restored to the last selected challenge

#### Scenario: Back from Neovim buffer
- **WHEN** the player is in a Neovim challenge buffer started from selection
- **THEN** canceling the challenge (Escape) SHALL return to the selection list
- **AND** the current challenge progress SHALL be discarded

### Requirement: Exit to Main Menu
The system SHALL provide a way to exit challenge selection and return to the main menu.

#### Scenario: Exit from selection list
- **WHEN** the player presses Escape while viewing the selection list
- **THEN** the game SHALL return to the start screen (level selection)
- **AND** the challenge selection state SHALL be cleared

#### Scenario: Exit from challenge practice
- **WHEN** the player uses the exit keymap during challenge practice
- **THEN** the game SHALL return to the start screen (level selection)
- **AND** both selection and practice states SHALL be cleared

### Requirement: Two-Column Selection Layout
The system SHALL render the selection screen with a two-column layout matching the level selection pattern.

#### Scenario: Layout structure
- **WHEN** challenge selection is displayed
- **THEN** a left column SHALL contain the challenge list
- **AND** a right column SHALL contain the challenge preview
- **AND** the layout SHALL match the level selection screen style

#### Scenario: Column sizing
- **WHEN** rendering the two-column layout
- **THEN** the list column SHALL be wide enough for challenge names and difficulty
- **AND** the preview column SHALL be wide enough for description and buffer preview
- **AND** consistent padding SHALL be applied

