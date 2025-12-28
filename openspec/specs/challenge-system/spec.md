# challenge-system Specification

## Purpose
TBD - created by archiving change add-keyforge-core. Update Purpose after archive.
## Requirements
### Requirement: Challenge Definition
The system SHALL support YAML-defined challenges (katas) with initial buffer content, validation rules, and reward configuration.

#### Scenario: Load challenge from YAML
- **WHEN** a challenge is requested by category and difficulty
- **THEN** a matching challenge SHALL be loaded from the assets
- **AND** the challenge SHALL include description, initial buffer, and validation criteria

#### Scenario: Custom challenge directory
- **WHEN** the user configures a custom challenges path
- **THEN** challenges from that directory SHALL be loaded
- **AND** they SHALL be available alongside built-in challenges

### Requirement: Challenge Categories
The system SHALL support multiple challenge categories mapped to vim skill areas: movement, text-objects, lsp-navigation, search-replace, refactoring, and git-operations.

#### Scenario: Category filtering
- **WHEN** the game requests a challenge of category "movement"
- **THEN** only movement-related challenges SHALL be considered
- **AND** difficulty SHALL further filter the selection

#### Scenario: Tower-category mapping
- **WHEN** a player places an LSP Tower
- **THEN** subsequent challenges for that tower SHALL be from the "lsp-navigation" category

### Requirement: Challenge Buffer
The system SHALL create a dedicated Neovim buffer for each challenge with the initial content and appropriate filetype.

#### Scenario: Buffer creation
- **WHEN** a challenge is started
- **THEN** a new buffer SHALL be created with the challenge content
- **AND** the buffer SHALL have the appropriate filetype for syntax highlighting
- **AND** the cursor SHALL be positioned at the starting location

#### Scenario: Buffer isolation
- **WHEN** the user is solving a challenge
- **THEN** changes SHALL be isolated to the challenge buffer
- **AND** user's other buffers SHALL not be affected

### Requirement: Keystroke Tracking
The system SHALL track all keystrokes during an active challenge for efficiency scoring.

#### Scenario: Count keystrokes
- **WHEN** a challenge is active
- **THEN** every keystroke SHALL be counted
- **AND** the count SHALL be accurate regardless of mappings

#### Scenario: Exclude non-editing keys
- **WHEN** the user presses non-editing keys (e.g., Escape to exit insert mode)
- **THEN** they SHALL still be counted as part of the solution
- **AND** the efficiency comparison SHALL use optimal keystroke count from the challenge definition

### Requirement: Challenge Validation
The system SHALL validate challenge completion by comparing the final file content against expected outcomes, with support for multiple validation strategies.

#### Scenario: Successful validation
- **WHEN** the file content matches the expected outcome
- **THEN** the challenge SHALL be marked as successful
- **AND** the completion result SHALL be sent to the game

#### Scenario: Failed validation
- **WHEN** the file content does not match expected outcome
- **THEN** the challenge SHALL be marked as failed
- **AND** the user MAY retry or skip (with penalty)

#### Scenario: Validation types
- **WHEN** a challenge uses "exact_match" validation
- **THEN** the file content SHALL exactly match the expected content
- **WHEN** a challenge uses "function_exists" validation
- **THEN** the specified function SHALL exist in the file
- **WHEN** a challenge uses "different" validation
- **THEN** the file content SHALL differ from the initial content

#### Scenario: Challenge timeout
- **WHEN** a challenge exceeds the configured timeout (default 5 minutes)
- **THEN** the challenge SHALL be automatically cancelled
- **AND** it SHALL be recorded as a timeout (no gold penalty beyond skip)

### Requirement: Efficiency Scoring
The system SHALL calculate an efficiency score based on keystroke count and time compared to optimal solutions.

#### Scenario: Calculate efficiency
- **WHEN** a challenge is completed successfully
- **THEN** efficiency SHALL be calculated as (optimal_keystrokes / actual_keystrokes)
- **AND** the score SHALL be capped at 1.0 (100%)

#### Scenario: Time bonus
- **WHEN** a challenge is completed under the par time
- **THEN** a time bonus multiplier SHALL be applied to the gold reward

#### Scenario: Gold calculation
- **WHEN** the efficiency score is calculated
- **THEN** gold reward SHALL be base_gold * efficiency * difficulty_multiplier
- **AND** minimum reward SHALL be 1 gold for any successful completion

### Requirement: Challenge Timer
The system SHALL track time spent on each challenge and display it to the user.

#### Scenario: Timer display
- **WHEN** a challenge is active
- **THEN** elapsed time SHALL be displayed in the game UI
- **AND** par time (if defined) SHALL also be shown

#### Scenario: Time tracking
- **WHEN** the challenge starts
- **THEN** a timer SHALL begin
- **WHEN** the challenge completes
- **THEN** total elapsed time SHALL be recorded in milliseconds

### Requirement: Challenge Completion Controls
The system SHALL provide configurable keymaps for submitting or cancelling challenges within the challenge buffer.

#### Scenario: Submit challenge
- **WHEN** the user presses the submit keymap (default `<CR>` in normal mode)
- **THEN** the current buffer content SHALL be validated
- **AND** the result SHALL be sent to the game
- **AND** the challenge buffer SHALL be cleaned up

#### Scenario: Cancel challenge
- **WHEN** the user presses the cancel keymap (default `<Esc>` in normal mode)
- **THEN** the challenge SHALL be marked as skipped
- **AND** no gold SHALL be awarded
- **AND** the challenge buffer SHALL be cleaned up

#### Scenario: Close tab without explicit action
- **WHEN** the user closes the challenge tab (`:q`, `:bd`, etc.) without pressing submit
- **THEN** the challenge SHALL be treated as cancelled/skipped
- **AND** the challenge buffer SHALL be cleaned up

