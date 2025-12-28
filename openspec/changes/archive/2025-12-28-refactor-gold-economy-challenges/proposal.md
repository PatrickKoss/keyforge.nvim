# Refactor Gold Economy and Challenge System

## Summary

Shift the primary gold economy from killing mobs to completing vim kata challenges. Players now control challenge timing via a button, allowing strategic breaks for tower placement. Challenges display context-aware hints based on the user's actual Neovim keymaps and installed plugins.

## Motivation

The current system awards gold primarily for killing mobs, which doesn't reinforce the core learning loop—practicing vim keybindings. By making challenges the primary gold source:

1. **Learning incentive**: Players must actively practice vim skills to progress
2. **Speed matters**: Faster challenge completion = more gold = more towers
3. **Strategic depth**: Players choose when to build vs when to continue challenges
4. **Personalized learning**: Hints reflect the user's actual keymap configuration

## Scope

### In Scope

- Modify gold economy: mob kills give reduced gold (25%), challenges are primary source
- User-controlled challenge triggering via keybind (not automatic)
- Game continues running while user solves challenges (real-time pressure)
- Dynamic keybinding hints based on `vim.api.nvim_get_keymap()` runtime inspection
- New challenges for common plugins (Telescope, nvim-tree, fugitive) within existing categories
- Challenge queue/selection UI showing available challenges and their rewards

### Out of Scope

- New tower types
- Wave system changes (beyond gold reduction)
- Multiplayer features
- Challenge editor/creator

## Design Decisions

### Gold Economy Balance

| Source | Current | Proposed |
|--------|---------|----------|
| Bug kill | 15g | 4g (25%) |
| Gremlin kill | 25g | 6g (25%) |
| Daemon kill | 40g | 10g (25%) |
| Boss kill | 100g | 25g (25%) |
| Wave bonus | 20-50g | 10-25g (50%) |
| Challenge (base) | N/A | 20-100g (based on difficulty) |
| Challenge (speed bonus) | N/A | up to 2x multiplier |

### Challenge Flow

```
[Tower View]                    [Challenge View]
     |                               |
     |-- User presses <leader>kn --> |
     |                               | Challenge starts
     |   Game continues running      | Timer visible
     |   Mobs still move/attack      | Keymap hints shown
     |                               |
     |<-- Challenge complete --------|
     |                               |
     |   Gold awarded                |
     |   User can place towers       |
     |                               |
     |-- User presses <leader>kn --> | (next challenge)
```

### Keybinding Discovery

The Lua plugin will:
1. Query `vim.api.nvim_get_keymap('n')` for normal mode mappings
2. Detect installed plugins via `pcall(require, 'telescope')` etc.
3. Build a hint database mapping actions to user's actual keybindings
4. Display relevant hints in challenge description based on challenge category

Example hint for a Telescope challenge:
```
Challenge: Find a file by name
Your keybinding: <leader>ff (Telescope find_files)
Hint: Use your fuzzy finder to locate 'config.lua'
```

### New Challenge Types (within existing categories)

**lsp-navigation** (extended):
- Telescope: find files, grep, buffers, help tags
- Quick fix list navigation
- Location list navigation

**movement** (extended):
- Window navigation (C-w commands)
- Tab navigation
- Buffer switching

**text-objects** (extended):
- Tree-sitter based text objects (if installed)
- Custom surround operations

**search-replace** (extended):
- Spectre-style project-wide search/replace
- Quickfix-based replacements

## Risks and Mitigations

| Risk | Mitigation |
|------|------------|
| Game too hard without mob gold | Tune challenge rewards to compensate; playtest balance |
| Keybinding detection slow | Cache keymap at game start, refresh on config change |
| Users with vanilla vim config | Provide default hint fallbacks for standard vim bindings |
| Challenge queue empty | Auto-refresh challenges; minimum pool of 30+ challenges |

## Success Criteria

1. Players can complete a full 10-wave game using only challenge gold
2. Keybinding hints correctly reflect user's actual mappings 90%+ of time
3. Game feels challenging but fair at normal difficulty
4. New plugin-aware challenges cover: Telescope, nvim-tree, fugitive basics

## Related Specs

- `challenge-system`: Modified to support user-triggered challenges and keymap hints
- `tower-defense`: Modified to reduce mob gold rewards
- `neovim-integration`: Extended for keymap discovery API
