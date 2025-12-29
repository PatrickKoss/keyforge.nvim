# Implementation Tasks

## Phase 1: State Machine Foundation

### Task 1.1: Add New Game States
- [x] Add `StateChallengeMode`, `StateChallengeSelection`, and `StateChallengeSelectionPractice` to the game state enum.
- File: `game/internal/engine/game.go`
- Verification: States compile and are accessible

### Task 1.2: Extend Model for New Modes
- [x] Add model fields for challenge mode state (streak, notification) and selection state (list, index, offset).
- File: `game/internal/ui/model.go`
- Verification: Model initializes correctly with new fields

### Task 1.3: Add Notification System
- [x] Implement notification struct and display logic with auto-dismiss timing.
- Files: `game/internal/ui/model.go`, `game/internal/ui/view.go`
- Verification: Notifications display and dismiss after 2 seconds

## Phase 2: Start Screen Integration

### Task 2.1: Add Mode Options to Start Screen
- [x] Extend start screen to show "Challenge Mode" and "Challenge Selection" options below levels.
- File: `game/internal/ui/start_screen.go`
- Verification: Both options visible, selectable via j/k navigation

### Task 2.2: Implement Cross-Section Navigation
- [x] Handle navigation between level list and mode options sections.
- File: `game/internal/ui/model.go` (handleLevelSelectKeys)
- Verification: Can navigate from levels to modes and back

### Task 2.3: Add Mode Selection Handlers
- [x] Implement Enter key handling for mode options to transition to appropriate states.
- File: `game/internal/ui/model.go`
- Verification: Selecting modes transitions to correct states

## Phase 3: Challenge Mode Implementation

### Task 3.1: Create Challenge Mode View
- [x] Implement UI rendering for challenge mode with header, streak, notification area.
- File: `game/internal/ui/start_screen.go` (RenderChallengeMode function)
- Verification: Challenge mode screen renders correctly

### Task 3.2: Implement Challenge Mode Loop
- [x] Add logic to request next challenge after completion in challenge mode.
- File: `game/internal/ui/model.go`
- Verification: Challenges load continuously after completion

### Task 3.3: Implement Challenge Mode Key Handlers
- [x] Add key handling for challenge mode (Escape to exit, challenge completion triggers).
- File: `game/internal/ui/model.go` (handleChallengeModeKeys)
- Verification: Escape returns to main menu, notifications show on completion

### Task 3.4: Update Neovim RPC for Challenge Mode
- [x] Add RPC handling to track challenge mode context and handle completion specially.
- Files: `game/internal/ui/model.go`, `game/internal/nvim/protocol.go`
- Verification: Challenge completion in Neovim triggers next challenge in mode

## Phase 4: Challenge Selection Implementation

### Task 4.1: Build Challenge List Data Structure
- [x] Implement challenge list loading grouped by category for selection UI.
- File: `game/internal/ui/model.go`, `game/internal/engine/challenges.go`
- Verification: All challenges loaded and grouped correctly

### Task 4.2: Create Challenge Selection View
- [x] Implement two-column layout with list on left, preview on right.
- File: `game/internal/ui/start_screen.go` (RenderChallengeSelection function)
- Verification: Selection screen matches design wireframe

### Task 4.3: Implement Selection Navigation
- [x] Add j/k navigation for challenge list with scrolling.
- File: `game/internal/ui/model.go` (handleChallengeSelectionKeys)
- Verification: Can navigate full list, preview updates

### Task 4.4: Implement Challenge Preview Panel
- [x] Render challenge preview with name, category, difficulty, description, buffer preview.
- File: `game/internal/ui/start_screen.go` (renderChallengePreview function)
- Verification: Preview shows all fields, updates on navigation

### Task 4.5: Implement Challenge Start from Selection
- [x] Handle Enter to start selected challenge, track index for progression.
- File: `game/internal/ui/model.go`
- Verification: Selected challenge starts, index tracked

### Task 4.6: Implement Sequential Progression
- [x] After completion, auto-load next challenge in list with wrap-around.
- File: `game/internal/ui/model.go`
- Verification: Completing challenge loads next, wraps at end

### Task 4.7: Implement Back to Selection
- [x] Add back navigation from practice to selection list, restoring position.
- File: `game/internal/ui/model.go`
- Verification: Escape returns to selection list at same position

### Task 4.8: Implement Exit to Main Menu
- [x] Handle Escape from selection list to return to start screen.
- File: `game/internal/ui/model.go`
- Verification: Escape returns to level selection

## Phase 5: Lua Integration

### Task 5.1: Update RPC Protocol
- [x] Add Mode field to ChallengeData and ChallengeRequest structs.
- Files: `game/internal/nvim/protocol.go`, `game/internal/nvim/client.go`, `game/internal/nvim/socket.go`
- Verification: Mode field included in RPC messages

### Task 5.2: Update Challenge UI for Mode Context
- [x] Pass mode context (challenge_mode, challenge_selection) in challenge requests.
- Files: `game/internal/ui/model.go`
- Verification: Lua knows which mode triggered challenge

### Task 5.3: Add Back-to-Selection Keymap
- [x] Escape from challenge practice returns to selection or main menu.
- File: `game/internal/ui/model.go`
- Verification: Keymap works correctly in all modes

## Phase 6: Testing

### Task 6.1: Unit Tests for State Transitions
- [x] Test new state transitions and edge cases.
- File: `game/internal/ui/model_test.go`
- Verification: Tests pass for all state transition scenarios

### Task 6.2: Unit Tests for Challenge Selection
- [x] Test list navigation, preview updates, index tracking.
- File: `game/internal/ui/model_test.go`
- Verification: Tests pass for navigation, bounds checking, scrolling, and challenge start

### Task 6.3: Integration Tests for Full Flows
- [x] Test complete flows: menu → mode → challenge → completion → next.
- File: `game/internal/ui/model_test.go`
- Verification: Full flow tests pass for both Challenge Mode and Challenge Selection

## Summary

All tasks (Phase 1-6) are complete. The feature is fully functional and tested:
- Challenge Mode and Challenge Selection appear in the start menu
- Navigation between levels and modes works correctly
- Both modes are fully playable with notifications
- RPC protocol updated to support mode context
- All tests pass (unit tests and integration tests)
- All linting passes (0 issues)
