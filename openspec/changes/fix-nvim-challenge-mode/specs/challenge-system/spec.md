## MODIFIED Requirements

### Requirement: Challenge Buffer
The system SHALL create a dedicated Neovim buffer for each challenge using a temporary file with the appropriate file extension to enable full Neovim features including LSP.

#### Scenario: Buffer creation
- **WHEN** a challenge is started
- **THEN** a temporary file SHALL be created with the challenge content
- **AND** the file SHALL have the appropriate extension for syntax highlighting and LSP (e.g., `.lua`, `.go`, `.ts`)
- **AND** the file SHALL be opened in a new tab
- **AND** the cursor SHALL be positioned at the starting location

#### Scenario: Buffer isolation
- **WHEN** the user is solving a challenge
- **THEN** changes SHALL be isolated to the challenge file
- **AND** user's other buffers SHALL not be affected

#### Scenario: LSP support
- **WHEN** the challenge file is opened
- **THEN** LSP servers SHALL attach based on filetype
- **AND** user SHALL have access to hover (K), go-to-definition (gd), and other LSP features
- **AND** ex commands (e.g., `:5`, `:s/old/new/`) SHALL work normally

#### Scenario: Challenge cleanup
- **WHEN** the challenge is completed or cancelled
- **THEN** the temporary file SHALL be deleted
- **AND** the challenge tab SHALL be closed
- **AND** focus SHALL return to the game tab

## MODIFIED Requirements

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

## ADDED Requirements

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
