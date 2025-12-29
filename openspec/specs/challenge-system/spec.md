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
The system SHALL support 15 challenge categories mapped to vim skill areas: movement, text-objects, lsp-navigation, search-replace, refactoring, git-operations, window-management, buffer-management, folding, quickfix, diagnostics, telescope, surround, harpoon, and formatting.

#### Scenario: Category filtering
- **WHEN** the game requests a challenge of category "movement"
- **THEN** only movement-related challenges SHALL be considered
- **AND** difficulty SHALL further filter the selection

#### Scenario: Tower-category mapping
- **WHEN** a player places an Arrow Tower
- **THEN** subsequent challenges for that tower SHALL be from movement, buffer-management, window-management, quickfix, or folding categories
- **WHEN** a player places an LSP Tower
- **THEN** subsequent challenges for that tower SHALL be from lsp-navigation, telescope, diagnostics, formatting, or harpoon categories
- **WHEN** a player places a Refactor Tower
- **THEN** subsequent challenges for that tower SHALL be from text-objects, search-replace, refactoring, surround, or git-operations categories

#### Scenario: New category availability
- **WHEN** the game requests a challenge of category "telescope"
- **THEN** telescope-specific challenges SHALL be returned
- **AND** required plugins SHALL be checked before selection

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
- **WHEN** a challenge uses "cursor_position" validation
- **THEN** the cursor SHALL be at the expected row and column
- **WHEN** a challenge uses "contains" validation
- **THEN** the file content SHALL contain the expected string
- **WHEN** a challenge uses "pattern" validation
- **THEN** the file content SHALL match the Lua pattern

#### Scenario: Challenge timeout
- **WHEN** a challenge exceeds the configured timeout (tier-based: quick=5s, standard=15s, complex=45s, expert=90s)
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

### Requirement: Extended Challenge Categories
The system SHALL support 15 challenge categories mapped to vim skill areas and user keybindings: movement, text-objects, lsp-navigation, search-replace, refactoring, git-operations, window-management, buffer-management, folding, quickfix, diagnostics, telescope, surround, harpoon, and formatting.

#### Scenario: New category registration
- **WHEN** challenges are loaded from YAML
- **THEN** all 15 categories SHALL be indexed
- **AND** each category SHALL have at least 5 challenges

#### Scenario: Category-to-tower mapping
- **WHEN** a tower fires and requests a challenge
- **THEN** the category SHALL be selected from the tower's category pool
- **AND** Arrow towers SHALL use movement, buffer-management, window-management, quickfix, folding
- **AND** LSP towers SHALL use lsp-navigation, telescope, diagnostics, formatting, harpoon
- **AND** Refactor towers SHALL use text-objects, search-replace, refactoring, surround, git-operations

### Requirement: Challenge Duration Tiers
The system SHALL classify challenges by duration tier based on expected keystroke count: quick (1-5), standard (6-15), complex (16-40), and expert (40+).

#### Scenario: Quick tier selection
- **WHEN** a challenge with 1-5 par_keystrokes is loaded
- **THEN** its duration_tier SHALL be "quick"
- **AND** it SHALL have a par_time of 5 seconds

#### Scenario: Standard tier selection
- **WHEN** a challenge with 6-15 par_keystrokes is loaded
- **THEN** its duration_tier SHALL be "standard"
- **AND** it SHALL have a par_time of 15 seconds

#### Scenario: Complex tier selection
- **WHEN** a challenge with 16-40 par_keystrokes is loaded
- **THEN** its duration_tier SHALL be "complex"
- **AND** it SHALL have a par_time of 45 seconds

#### Scenario: Expert tier selection
- **WHEN** a challenge with 40+ par_keystrokes is loaded
- **THEN** its duration_tier SHALL be "expert"
- **AND** it SHALL have a par_time of 90 seconds

### Requirement: Runtime Keymap Resolution
The system SHALL dynamically resolve keybindings from the user's Neovim configuration at runtime using `vim.api.nvim_get_keymap()` and keymap descriptions.

#### Scenario: Resolve keymap by action pattern
- **WHEN** a challenge specifies `hint_action: "find_files"`
- **THEN** the system SHALL search all keymaps for one that triggers telescope.builtin.find_files or has a matching description
- **AND** the user's actual keymap (e.g., `<leader>ff`) SHALL be shown in the hint

#### Scenario: Resolve keymap by rhs pattern
- **WHEN** searching for a keymap
- **THEN** the system SHALL check the `rhs` field for matching Lua function calls or command patterns
- **AND** the system SHALL check the `desc` field for matching descriptions (case-insensitive)

#### Scenario: Keymap resolution fallback
- **WHEN** no matching keymap is found for an action
- **THEN** the system SHALL display a generic hint (e.g., "Use your find files keymap")
- **AND** the challenge SHALL still be playable

#### Scenario: Standard vim commands
- **WHEN** a challenge uses a standard vim command (e.g., `$`, `ciw`, `dd`)
- **THEN** the hint SHALL show the standard command directly
- **AND** no keymap resolution SHALL be needed

### Requirement: Plugin Availability Detection
The system SHALL detect which plugins are installed in the user's Neovim configuration and filter challenges accordingly.

#### Scenario: Detect Lua plugin availability
- **WHEN** a challenge has `required_plugin: "telescope"`
- **THEN** the system SHALL check if `require("telescope")` succeeds
- **AND** if the plugin is not available, the challenge SHALL be excluded from selection

#### Scenario: Detect plugin by command
- **WHEN** a challenge requires a plugin that provides vim commands
- **THEN** the system SHALL check `vim.fn.exists(":CommandName")` as a fallback
- **AND** a return value of 2 indicates the command exists

#### Scenario: Plugin name variations
- **WHEN** checking plugin availability
- **THEN** the system SHALL try common name variations (e.g., "nvim-surround" -> "nvim_surround", "mini.surround")
- **AND** the system SHALL try requiring the plugin's main module

#### Scenario: Cache plugin detection results
- **WHEN** plugin availability is checked
- **THEN** the result SHALL be cached for the session
- **AND** subsequent checks for the same plugin SHALL use the cached result

### Requirement: Dynamic Hint System
The system SHALL support dynamic hint resolution for challenges using action-based lookups instead of hardcoded keybindings.

#### Scenario: Action-based hint resolution
- **WHEN** a challenge has a `hint_action` field (e.g., "format_buffer", "goto_definition")
- **THEN** the system SHALL resolve the user's keymap for that action
- **AND** the hint SHALL display the user's actual keybinding

#### Scenario: Static hint fallback
- **WHEN** a challenge has a `hint_fallback` field
- **THEN** this text SHALL be used if keymap resolution fails

#### Scenario: No hint needed
- **WHEN** a challenge uses only standard vim commands
- **THEN** the description SHALL serve as the only guidance
- **AND** no dynamic resolution SHALL be performed

### Requirement: Challenge Library Scale
The system SHALL maintain a library of at least 150 challenges across all categories with varied difficulty distribution.

#### Scenario: Minimum challenge count
- **WHEN** challenges are loaded
- **THEN** at least 150 challenges SHALL be available
- **AND** each difficulty level (1, 2, 3) SHALL have at least 30 challenges

#### Scenario: Category balance
- **WHEN** counting challenges per category
- **THEN** no category SHALL have fewer than 5 challenges
- **AND** no category SHALL have more than 25 challenges

### Requirement: Challenge Validation Test Coverage
The system SHALL have comprehensive test coverage for challenge validation logic covering all validation types and edge cases.

#### Scenario: Validation type coverage
- **WHEN** running validation tests
- **THEN** tests SHALL exist for exact_match, contains, cursor_position, function_exists, pattern, different, and cursor_on_char
- **AND** each type SHALL have at least 3 test cases (success, failure, edge case)

#### Scenario: Challenge sample testing
- **WHEN** running challenge tests
- **THEN** at least 20% of challenges SHALL have explicit win/fail test cases
- **AND** all new categories SHALL have at least 2 sample tests each

#### Scenario: Edge case coverage
- **WHEN** testing validation logic
- **THEN** tests SHALL cover empty buffers, unicode content, multiline content, and timeout scenarios

