# challenge-system Spec Delta

## MODIFIED Requirements

### Requirement: User-Triggered Challenges
The system SHALL allow users to manually trigger challenges via a keybinding instead of automatic tower-based triggering.

#### Scenario: Start challenge on demand
- **WHEN** the user presses the "next challenge" keybinding (`<leader>kn`)
- **THEN** the next available challenge SHALL be loaded
- **AND** the challenge buffer SHALL be created
- **AND** keystroke tracking SHALL begin
- **AND** the game SHALL continue running (not pause)

#### Scenario: No challenge while one is active
- **WHEN** a challenge is already active
- **AND** the user presses the "next challenge" keybinding
- **THEN** the request SHALL be ignored
- **AND** a notification MAY be shown

#### Scenario: Challenge queue empty
- **WHEN** the user requests a challenge
- **AND** no challenges are available
- **THEN** a notification SHALL inform the user
- **AND** challenges SHALL refresh after a cooldown

### Requirement: Concurrent Gameplay During Challenges
The system SHALL continue running the game loop while a challenge is active, creating time pressure.

#### Scenario: Game continues during challenge
- **WHEN** a challenge is active
- **THEN** enemies SHALL continue moving along the path
- **AND** towers SHALL continue firing at enemies
- **AND** the player MAY lose health if enemies reach the exit

#### Scenario: Challenge timeout awareness
- **WHEN** the player is solving a challenge
- **THEN** elapsed time SHALL be visible
- **AND** par time SHALL be displayed for speed bonus reference

### Requirement: Speed Bonus Scoring
The system SHALL award a speed bonus multiplier for completing challenges faster than the par time.

#### Scenario: Calculate speed bonus
- **WHEN** a challenge is completed in less than par time
- **THEN** a speed bonus multiplier SHALL be calculated
- **AND** the multiplier SHALL be `min(2.0, 1.0 + (par_time / actual_time - 1) * 0.5)`
- **AND** the bonus SHALL cap at 2.0x

#### Scenario: No penalty for slow completion
- **WHEN** a challenge is completed slower than par time
- **THEN** the speed bonus SHALL be 1.0 (no bonus, no penalty)
- **AND** the base gold reward SHALL still be awarded

#### Scenario: Gold calculation with speed bonus
- **WHEN** the efficiency and speed bonus are calculated
- **THEN** gold reward SHALL be `base_gold * difficulty_mult * efficiency_mult * speed_bonus`
- **AND** minimum reward SHALL remain 1 gold

## ADDED Requirements

### Requirement: Keybinding Hints
The system SHALL display context-aware keybinding hints based on the user's actual Neovim configuration.

#### Scenario: Discover user keymaps
- **WHEN** the game starts
- **THEN** the system SHALL query `vim.api.nvim_get_keymap()` for all modes
- **AND** the keymap cache SHALL be populated
- **AND** detected plugins SHALL be recorded

#### Scenario: Display relevant hints
- **WHEN** a challenge is displayed
- **THEN** the challenge description SHALL include keybinding hints
- **AND** hints SHALL reflect the user's actual mappings
- **AND** if no custom mapping exists, default vim bindings SHALL be shown

#### Scenario: Plugin-specific hints
- **WHEN** a challenge requires a plugin (e.g., Telescope)
- **AND** the plugin is detected
- **THEN** hints SHALL show the user's bindings for that plugin
- **AND** the plugin name SHALL be mentioned in the hint

#### Scenario: Refresh keymaps
- **WHEN** the user triggers a keymap refresh command
- **THEN** the keymap cache SHALL be updated
- **AND** subsequent hints SHALL reflect any changes

### Requirement: Plugin-Aware Challenges
The system SHALL support challenges that require specific plugins and filter availability based on detected plugins.

#### Scenario: Plugin detection
- **WHEN** the game initializes
- **THEN** common plugins SHALL be detected via `pcall(require, 'plugin_name')`
- **AND** detected plugins SHALL be stored for challenge filtering

#### Scenario: Filter challenges by plugin
- **WHEN** challenges are loaded
- **AND** a challenge has `required_plugin` metadata
- **THEN** the challenge SHALL only be available if the plugin is detected
- **AND** challenges without plugin requirements SHALL always be available

#### Scenario: Telescope challenges
- **WHEN** Telescope is detected
- **THEN** challenges for find_files, live_grep, buffers SHALL be available
- **AND** hints SHALL show the user's Telescope keybindings

#### Scenario: File explorer challenges
- **WHEN** nvim-tree or neo-tree is detected
- **THEN** file navigation challenges SHALL be available
- **AND** hints SHALL show relevant tree navigation bindings

### Requirement: Challenge Availability Indicator
The system SHALL indicate when challenges are available and show estimated rewards.

#### Scenario: Show available challenge count
- **WHEN** the game HUD is rendered
- **THEN** the number of available challenges SHALL be visible
- **AND** a prompt to start the next challenge MAY be shown

#### Scenario: Show estimated reward
- **WHEN** a challenge is about to start
- **THEN** the base gold reward SHALL be displayed
- **AND** the potential speed bonus range SHALL be indicated
