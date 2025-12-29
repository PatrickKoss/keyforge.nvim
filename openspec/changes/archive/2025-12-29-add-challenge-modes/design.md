# Design: Challenge Mode and Challenge Selection

## Overview
This document describes the architecture for adding two standalone practice modes to Keyforge: Challenge Mode (endless practice) and Challenge Selection (targeted practice with browsing).

## Architecture

### State Machine Extensions

Current game states:
```
StateLevelSelect → StateSettings → StatePlaying → StateGameOver/Victory
```

New states:
```
StateLevelSelect
    ├── StateSettings → StatePlaying
    ├── StateChallengeMode → (loops challenges until exit)
    └── StateChallengeSelection → StateChallengeSelectionPractice → (returns to selection)
```

### Go Game Engine Changes

#### New States (`game/internal/engine/game.go`)
```go
const (
    // ... existing states ...
    StateChallengeMode           // Endless challenge practice
    StateChallengeSelection      // Challenge list browsing
    StateChallengeSelectionPractice // Doing a selected challenge
)
```

#### New Model Fields (`game/internal/ui/model.go`)
```go
type Model struct {
    // ... existing fields ...

    // Challenge mode state
    ChallengeModeActive    bool
    ChallengeModeStreak    int      // Successful challenges in a row

    // Challenge selection state
    ChallengeList          []engine.Challenge  // All loaded challenges
    ChallengeListIndex     int                 // Currently hovered challenge
    ChallengeListOffset    int                 // Scroll offset for long lists
    SelectedChallengeIndex int                 // Which challenge is being practiced

    // Notification state
    LastResultSuccess      bool
    LastResultTime         time.Time
    ShowNotification       bool
}
```

### UI Layout

#### Start Screen with New Options
```
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║   ██╗  ██╗███████╗██╗   ██╗███████╗ ██████╗ ██████╗  ...    ║
║                                                              ║
║   Select Level                     ┌─────────────────────┐   ║
║                                    │ Level 1: The Path   │   ║
║   ★☆☆ Level 1: The Path       →   │ Waves: 5            │   ║
║   ★☆☆ Level 2: Forked Roads       │ Difficulty: Beginner│   ║
║   ★★☆ Level 3: The Maze           │                     │   ║
║   ...                              │ ░░░░░░░░░░░░░       │   ║
║                                    └─────────────────────┘   ║
║   ─────────────────                                          ║
║   ⚔️  Challenge Mode                                         ║
║   📋 Challenge Selection                                     ║
║                                                              ║
║   [j/k] Select  [Enter] Start  [q] Quit                      ║
╚══════════════════════════════════════════════════════════════╝
```

#### Challenge Mode Screen
```
╔══════════════════════════════════════════════════════════════╗
║  CHALLENGE MODE              Streak: 5  ✓ Success!           ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  Movement Challenge (★☆☆)                                    ║
║  ─────────────────────────                                   ║
║  Move cursor to end of line using $                          ║
║                                                              ║
║  ┌──────────────────────────────────────────────────────┐   ║
║  │ The quick brown fox jumps over the lazy dog█         │   ║
║  └──────────────────────────────────────────────────────┘   ║
║                                                              ║
║  [Enter] Submit  [Esc] Back to Menu                          ║
╚══════════════════════════════════════════════════════════════╝
```

#### Challenge Selection Screen
```
╔══════════════════════════════════════════════════════════════╗
║  CHALLENGE SELECTION                                         ║
╠═══════════════════════════════╦══════════════════════════════╣
║  movement (12)                ║  Preview: Jump to End        ║
║  ► Jump to End            ★☆☆║  ─────────────────────────── ║
║    Word Hop               ★☆☆║  Category: movement          ║
║    Find the X             ★☆☆║  Difficulty: ★☆☆             ║
║  text-objects (8)             ║                              ║
║    Change Inside Quotes   ★★☆║  Move cursor to end of       ║
║    Delete Inside Parens   ★★☆║  line using $                ║
║  search-replace (5)           ║                              ║
║    Simple Replace         ★☆☆║  Buffer:                     ║
║    Global Replace         ★★☆║  ┌────────────────────────┐ ║
║  ...                          ║  │ The quick brown fox... │ ║
║                               ║  └────────────────────────┘ ║
╠═══════════════════════════════╩══════════════════════════════╣
║  [j/k] Navigate  [Enter] Start  [Esc] Back to Menu           ║
╚══════════════════════════════════════════════════════════════╝
```

### Notification System

Notifications appear in a fixed position that doesn't overlay the challenge description:
- Position: Top-right corner or inline with title
- Duration: 2 seconds auto-dismiss
- Styling: Green for success (✓), Red for failure (✗)

```go
type Notification struct {
    Message   string
    IsSuccess bool
    ShowUntil time.Time
}
```

### RPC Communication

New RPC methods for Neovim integration:

**Game → Neovim:**
- `challenge_mode_start` - Entering challenge mode
- `challenge_selection_start` - Entering selection mode
- `challenge_selection_list` - Send challenge list for Lua-side rendering

**Neovim → Game:**
- `exit_challenge_mode` - Return to main menu from challenge mode
- `select_challenge` - Select a specific challenge from selection list
- `exit_challenge_selection` - Return to main menu from selection

### Challenge Flow

#### Challenge Mode Flow
1. User selects "Challenge Mode" from menu
2. Game enters `StateChallengeMode`
3. Random challenge selected and sent to Neovim
4. User completes challenge in Neovim buffer
5. Result received → Show notification (doesn't overlay description)
6. After 1 second delay, next random challenge loads
7. `Esc` in game UI returns to main menu

#### Challenge Selection Flow
1. User selects "Challenge Selection" from menu
2. Game enters `StateChallengeSelection`
3. Challenge list rendered with categories and previews
4. User navigates with j/k, sees preview on right
5. User presses Enter → Game enters `StateChallengeSelectionPractice`
6. Challenge sent to Neovim
7. Result received → Show notification
8. After notification, next challenge in list loads (or wrap)
9. "Back" button in challenge UI returns to selection list
10. `Esc` from selection list returns to main menu

### Menu Navigation Extension

The start screen needs to track a combined index across levels and mode options:

```go
type StartMenuSection int

const (
    SectionLevels StartMenuSection = iota
    SectionModes
)

// In Model:
StartSection      StartMenuSection
ModeMenuIndex     int  // 0 = Challenge Mode, 1 = Challenge Selection
```

Navigation logic:
- `j` at bottom of levels → move to modes section
- `k` at top of modes → move back to levels section
- `Enter` on level → settings screen
- `Enter` on mode → corresponding mode screen

## Trade-offs

### Decision: Separate States vs Mode Flags
**Chosen:** Separate states (`StateChallengeMode`, `StateChallengeSelection`)
**Reason:** Cleaner state machine, explicit transitions, easier testing
**Alternative:** Single `StatePractice` with mode flag - would complicate state logic

### Decision: Notification Placement
**Chosen:** Top-right inline with title bar
**Reason:** Never overlays content, visible but not intrusive
**Alternative:** Modal popup - would interrupt flow and overlay content

### Decision: Challenge Selection with Preview
**Chosen:** Two-column layout (list + preview)
**Reason:** Matches existing level selection pattern, provides context before starting
**Alternative:** Full-screen preview on hover - too disruptive

## Testing Strategy

1. **Unit tests:** State transitions, notification timing
2. **Integration tests:** Challenge flow completion, menu navigation
3. **Manual tests:** Visual layout, notification visibility, preview accuracy
