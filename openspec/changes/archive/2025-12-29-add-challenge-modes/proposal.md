# Proposal: Add Challenge Mode and Challenge Selection

## Status
PROPOSED

## Summary
Add two new gameplay modes accessible from the main menu level selection screen: Challenge Mode for endless practice with continuous challenges, and Challenge Selection for browsing and practicing specific challenges with preview functionality.

## Motivation
Currently, players can only practice vim keybindings during tower defense gameplay. This limits practice opportunities and doesn't allow targeted skill improvement. Players should be able to:
1. Practice challenges endlessly without tower defense mechanics (Challenge Mode)
2. Browse all available challenges and pick specific ones to practice (Challenge Selection)

Both modes enhance the learning experience by providing focused, distraction-free practice environments.

## Scope

### In Scope
- Two new menu items below level selection: "Challenge Mode" and "Challenge Selection"
- Challenge Mode: Endless challenge loop with success/failure notifications
- Challenge Selection: Browsable challenge list with preview and sequential progression
- Navigation back to main menu from both modes
- Notifications for success/failure that don't overlay challenge descriptions

### Out of Scope
- Statistics tracking and persistence across sessions
- Leaderboards or scoring comparisons
- Custom challenge creation interface
- Difficulty filtering in selection (may be future enhancement)

## Design Considerations

### Menu Integration
Both modes appear below the level list in the start screen, using consistent styling with existing menu items. This keeps all game modes in one place and maintains visual hierarchy.

### State Management
Two new game states will be added:
- `StateChallengeMode` - Endless practice loop
- `StateChallengeSelection` - Browsing and selecting challenges

### Notification Placement
Success/failure notifications will appear in a dedicated area (e.g., top-right corner or status bar) separate from the challenge description area, ensuring players always see the challenge requirements.

### Challenge Selection Preview
When hovering over a challenge in the selection list, a preview shows:
- Challenge name and category
- Difficulty level
- Description
- Initial buffer preview (truncated)

## Capabilities
- [challenge-practice-mode](specs/challenge-practice-mode/spec.md) - Endless challenge practice with notifications
- [challenge-selection-mode](specs/challenge-selection-mode/spec.md) - Challenge browsing and selection with preview

## Related Specs
- challenge-system - Existing challenge validation and scoring
- neovim-integration - RPC communication for challenge handling
- game-engine - State machine and UI rendering

## Tasks
See [tasks.md](tasks.md) for implementation plan.
