## Context

Keyforge is a tower defense game that runs as a Bubbletea TUI inside a Neovim terminal. Challenges are vim kata-style editing tasks that reward gold. The current implementation has two modes:
1. **Standalone mode**: Game handles challenges internally with a basic vim emulator
2. **Nvim mode**: Game delegates challenges to Neovim via JSON-RPC

The Lua plugin currently launches the game without `--nvim-mode`, so all challenges use the limited internal emulator. Even when nvim mode is enabled, the challenge buffer is a scratch floating window that doesn't support LSP, ex commands, or proper vim semantics.

**Stakeholders**: Users who want to practice real vim skills with their custom keymaps and LSP configurations.

## Goals / Non-Goals

**Goals:**
- Challenges use REAL Neovim buffers with full feature support (LSP, commands, user keymaps)
- Game pauses during challenge so user can focus on editing
- Challenge workflow is seamless (start challenge -> edit -> submit -> back to game)
- Game state is consistent after challenge completion
- Game over/victory screens are visible and actionable in Neovim

**Non-Goals:**
- Real-time game updates during challenge (this creates the problem we're solving)
- Supporting multiple simultaneous challenges
- Multiplayer or networked game state
- Persisting game state across Neovim sessions

## Decisions

### Decision 1: Pause game during challenge
**What**: When a challenge starts, game enters `StatePaused` and sends state snapshot to Neovim.
**Why**:
- User can focus on editing without divided attention
- Game state stays consistent (no ghost enemies)
- Simpler to reason about than real-time sync
**Alternatives considered**:
- Keep game running: Adds complexity, user frustration, state sync issues
- Slow-motion game: Still divides attention, still needs state sync

### Decision 2: Real file buffer instead of scratch buffer
**What**: Challenge content is written to a temp file (e.g., `/tmp/keyforge_challenge_abc123.lua`) and opened as a real buffer.
**Why**:
- LSP servers attach based on filetype, which requires a real file
- Ex commands like `:5` work normally
- User's vim config applies (statusline, colorscheme, etc.)
**Alternatives considered**:
- Scratch buffer with manual LSP attach: Complex, unreliable across LSP servers
- Virtual buffer: No persistence, limited tooling support

### Decision 3: Challenge in new tab with callback on close
**What**: Challenge opens in a new tab. When user closes the tab (`:q`, `:wq`, etc.) or presses submit keymap, challenge is validated.
**Why**:
- Clean separation between game and challenge
- User can use `:wq` workflow naturally
- Tab close is a natural "I'm done" signal
**Alternatives considered**:
- Split window: Cluttered, hard to focus
- Replace game buffer: Loses game visual state, confusing

### Decision 4: State snapshot and restore
**What**: When challenge starts, Go sends full game state. When challenge ends, Lua sends result back and Go restores from paused state.
**Why**:
- Game can display "CHALLENGE IN PROGRESS" during editing
- State is deterministic - no drift from expected positions
- Handles edge cases like game-over-during-challenge gracefully

### Decision 5: Configurable completion keymaps in challenge buffer
**What**: User can configure keymaps for submit (`<CR>` default) and cancel (`<Esc>` default) in the challenge buffer.
**Why**:
- `<CR>` in normal mode is intuitive for "submit"
- `<Esc>` is standard for "cancel/abort"
- Users may want different keymaps based on workflow
**Alternatives considered**:
- Only `:w` to submit: Doesn't work for cursor-position challenges
- Only close tab: Too easy to accidentally quit

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Temp files left behind on crash | Use `vim.loop.fs_unlink` on buffer wipe, also clean on plugin load |
| LSP slow to attach | Pre-create file with content, delay validation until LSP ready or timeout |
| User closes tab without submitting | Treat as skip, send skip result to game |
| Game process dies during challenge | Detect disconnection, clean up challenge, show error |
| Terminal rendering glitches after tab switch | Force redraw game terminal on focus return |

## Migration Plan

1. **Phase 1**: Add `--nvim-mode` flag to game launch (minimal fix, enables RPC)
2. **Phase 2**: Rewrite challenge buffer to use temp files in new tab
3. **Phase 3**: Add game pause/resume during challenge
4. **Phase 4**: Add game over/victory detection and display in Neovim

**Rollback**: Can revert to standalone mode by removing `--nvim-mode` flag.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CHALLENGE FLOW                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────┐                      ┌─────────────────────────────┐   │
│  │  GAME (Go TUI)  │                      │      NEOVIM (Lua)           │   │
│  │                 │                      │                             │   │
│  │  User presses   │                      │                             │   │
│  │  'c' on tower   │                      │                             │   │
│  │       │         │                      │                             │   │
│  │       ▼         │                      │                             │   │
│  │  game.Pause()   │                      │                             │   │
│  │  snapshot state │                      │                             │   │
│  │       │         │   RPC: request_      │                             │   │
│  │       │         │   challenge          │                             │   │
│  │       └─────────┼──────────────────────▶  1. Create temp file        │   │
│  │                 │   {id, category,     │     /tmp/kf_xxx.lua         │   │
│  │                 │    difficulty,       │                             │   │
│  │                 │    state_snapshot}   │  2. Write initial content   │   │
│  │                 │                      │                             │   │
│  │  Display:       │                      │  3. Open in new tab         │   │
│  │  "CHALLENGE     │                      │     with proper filetype    │   │
│  │   IN PROGRESS"  │                      │                             │   │
│  │                 │                      │  4. Wait for LSP attach     │   │
│  │       ...       │                      │                             │   │
│  │                 │                      │  5. User edits with full    │   │
│  │  (game paused,  │                      │     vim features            │   │
│  │   no updates)   │                      │                             │   │
│  │                 │                      │  6. User presses <CR>       │   │
│  │                 │                      │     or closes tab           │   │
│  │                 │   RPC: challenge_    │                             │   │
│  │       ┌─────────│◀──────────────────────  7. Validate buffer         │   │
│  │       │         │   complete           │     Send result             │   │
│  │       ▼         │   {id, success,      │                             │   │
│  │  game.Resume()  │    gold, ...}        │  8. Clean up temp file      │   │
│  │  award gold     │                      │     Return to game tab      │   │
│  │                 │                      │                             │   │
│  └─────────────────┘                      └─────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

## Open Questions

1. Should we support "practice mode" where game doesn't pause? (Deferred - can add later as config option)
2. Should challenge timer continue while LSP is attaching? (Proposed: Start timer after first keystroke)
3. How to handle LSP servers that never attach? (Proposed: 5s timeout, allow completion without LSP)
