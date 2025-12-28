# neovim-integration Spec Delta

## ADDED Requirements

### Requirement: Keymap Discovery API
The system SHALL provide an API to discover and cache the user's Neovim keybindings.

#### Scenario: Query normal mode keymaps
- **WHEN** `discover_keymaps()` is called
- **THEN** all normal mode mappings SHALL be retrieved via `vim.api.nvim_get_keymap('n')`
- **AND** results SHALL be cached for performance

#### Scenario: Query visual mode keymaps
- **WHEN** `discover_keymaps()` is called
- **THEN** visual mode mappings SHALL also be retrieved
- **AND** they SHALL be stored separately in the cache

#### Scenario: Parse keymap structure
- **WHEN** a keymap is retrieved
- **THEN** the lhs (key sequence) SHALL be extracted
- **AND** the rhs or callback SHALL be identified
- **AND** the description (if present) SHALL be stored

#### Scenario: Cache performance
- **WHEN** keymaps are queried multiple times
- **THEN** cached results SHALL be returned
- **AND** cache age SHALL be tracked
- **AND** stale cache (>5 minutes) MAY trigger refresh

### Requirement: Plugin Detection
The system SHALL detect installed Neovim plugins to enable plugin-specific challenges.

#### Scenario: Detect Telescope
- **WHEN** plugin detection runs
- **THEN** `pcall(require, 'telescope')` SHALL be attempted
- **AND** success SHALL mark Telescope as available
- **AND** Telescope-related challenges SHALL be enabled

#### Scenario: Detect file explorer
- **WHEN** plugin detection runs
- **THEN** nvim-tree and neo-tree SHALL be checked
- **AND** at least one detection SHALL enable file explorer challenges

#### Scenario: Detect git plugins
- **WHEN** plugin detection runs
- **THEN** fugitive SHALL be checked via `vim.fn.exists(':Git')`
- **AND** gitsigns SHALL be checked via require
- **AND** detection SHALL enable git-related challenges

#### Scenario: Detect motion plugins
- **WHEN** plugin detection runs
- **THEN** flash.nvim, leap.nvim, hop.nvim SHALL be checked
- **AND** detection SHALL enable enhanced motion challenges

#### Scenario: Detect surround plugins
- **WHEN** plugin detection runs
- **THEN** nvim-surround and mini.surround SHALL be checked
- **AND** detection SHALL enable surround operation challenges

### Requirement: Challenge Trigger Keybindings
The system SHALL provide keybindings for user-controlled challenge flow.

#### Scenario: Next challenge keybinding
- **WHEN** the game is active
- **THEN** `<leader>kn` SHALL trigger the next challenge
- **AND** the keybinding SHALL be configurable via setup()

#### Scenario: Complete challenge keybinding
- **WHEN** a challenge is active
- **THEN** `<leader>kc` SHALL validate and complete the challenge
- **AND** the result SHALL be sent to the game engine

#### Scenario: Skip challenge keybinding
- **WHEN** a challenge is active
- **THEN** `<leader>ks` SHALL skip the challenge
- **AND** a small penalty MAY be applied
- **AND** the next challenge SHALL become available

### Requirement: Gold Notification
The system SHALL notify the user when gold is earned from challenges.

#### Scenario: Display gold earned
- **WHEN** a challenge is completed successfully
- **THEN** the gold earned SHALL be displayed
- **AND** the speed bonus (if any) SHALL be shown
- **AND** the notification SHALL fade after a few seconds

#### Scenario: Update HUD gold
- **WHEN** gold is earned from a challenge
- **THEN** the game HUD gold display SHALL update
- **AND** the update SHALL be visible even during the next challenge
