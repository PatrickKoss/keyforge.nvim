# Proposal: Add Multi-Level Progression System

## Summary

Expand the game from a single "Classic" level to 10 progressively difficult levels. Each level features unique enemy paths, scaled mob spawns (3-7 enemy types per level), new enemy variants with distinct speed/health profiles, rebalanced towers with cost-based stats and range visualization, and improved upgrade mechanics.

## Motivation

The current game has only the "Classic" level available in the start screen level selector. This limits replayability and doesn't provide a clear progression path for players to improve their vim skills. A 10-level system with escalating difficulty creates:

1. **Clear progression** - Players start with simple paths and basic enemies, graduating to complex layouts and challenging mob compositions
2. **Extended gameplay** - More content to explore beyond the single level
3. **Skill-appropriate challenges** - Easy levels for beginners, hard levels for vim veterans
4. **Balanced variety** - 3-7 enemy types per level keeps waves interesting without overwhelming

## Scope

### In Scope

1. **10 Levels with unique paths** - Each level has a different path layout; higher levels have longer, more complex paths
2. **New enemy types** - Add slow/tanky and fast/fragile variants to the existing 4 types
3. **Enemy progression by level** - Easy enemies (Bug, Mite) in levels 1-3; medium enemies in levels 4-6; hard enemies (Daemon, Boss) in levels 7-10
4. **Tower stat rebalancing** - Damage, attack speed, and range tied to cost tier
5. **Range visualization** - Show tower range during placement; show on hover after placement
6. **Upgrade improvements** - Upgrades increase damage, range, and attack speed by small amounts
7. **Gold economy tuning** - High-health enemies yield more gold; low-health yield less; challenges remain primary gold source

### Out of Scope

- New tower types (keep existing 3)
- New challenge categories
- Visual path decorations
- Level unlock system (all levels available immediately)
- Persistent progression/save system

## Design Decisions

### Level Path Complexity Scaling

| Level | Path Segments | Approx Length | Description |
|-------|--------------|---------------|-------------|
| 1 | 2 | 15 cells | Straight line |
| 2 | 3 | 20 cells | Single turn |
| 3 | 4 | 25 cells | L-shape |
| 4-5 | 5-6 | 30-35 cells | S-curve |
| 6-7 | 7-8 | 40-45 cells | Zigzag |
| 8-9 | 9-10 | 50-55 cells | Spiral approach |
| 10 | 12+ | 60+ cells | Complex maze |

### New Enemy Types

| Enemy | Health | Speed | Gold | Levels |
|-------|--------|-------|------|--------|
| Mite (new) | 5 | 2.0 | 2 | 1-4 |
| Bug | 10 | 1.5 | 5 | 1-6 |
| Gremlin | 25 | 2.5 | 10 | 3-8 |
| Crawler (new) | 40 | 0.6 | 15 | 4-8 |
| Daemon | 100 | 0.8 | 25 | 6-10 |
| Specter (new) | 15 | 3.5 | 8 | 5-9 |
| Boss | 500 | 0.5 | 100 | 10 |

Key attributes:
- **Mite**: Very low health, fast, early-game fodder
- **Crawler**: High health, very slow, mid-game tank
- **Specter**: Low health, very fast, requires quick targeting

### Enemy Distribution Per Level

| Level | Enemy Types | Count Range | Notes |
|-------|------------|-------------|-------|
| 1 | Mite, Bug | 3-4 | Introduction |
| 2 | Mite, Bug | 3-4 | Basic variety |
| 3 | Bug, Gremlin | 4-5 | Speed introduced |
| 4 | Bug, Gremlin, Crawler | 4-5 | Tank introduced |
| 5 | Bug, Gremlin, Specter | 5-6 | Fast enemy |
| 6 | Gremlin, Crawler, Daemon | 5-6 | Heavy mix |
| 7 | Gremlin, Specter, Daemon | 5-6 | Speed & power |
| 8 | Crawler, Specter, Daemon | 5-7 | Late game mix |
| 9 | Specter, Daemon | 5-7 | Pre-boss |
| 10 | All except Mite + Boss | 6-7 | Boss wave |

### Tower Stat Rebalancing

| Tower | Cost | Damage | Range | Cooldown | Notes |
|-------|------|--------|-------|----------|-------|
| Arrow | 50g | 8 | 2.5 | 0.8s | Fast, short range |
| LSP | 100g | 20 | 5.0 | 1.5s | Sniper, slow |
| Refactor | 150g | 12 | 3.0 | 1.0s | Area damage |

### Tower Upgrade Scaling

Each upgrade tier provides:
- +15% damage (additive)
- +0.3 range
- -10% cooldown (multiplicative)

### Range Visualization

- **During placement**: Show circular range indicator around cursor
- **After placement**: Show range only when cursor hovers over tower
- Use semi-transparent overlay or dotted circle

### Gold Economy

Challenge gold remains primary income (60-70% of total). Mob gold scales with health:
- Under 20 HP: 2-5 gold
- 20-50 HP: 8-15 gold
- 50-150 HP: 20-30 gold
- 150+ HP: 50-100 gold

## Related Specs

- `tower-defense` - Enemy Types, Tower Types, Wave System, Resource Economy
- `game-engine` - Level Definition System, Level Registry, Level Preview

## Risks

1. **Balance tuning** - May need iteration to ensure all 10 levels are completable
2. **Path rendering** - Complex paths need clear visual distinction
3. **Performance** - More enemy types and longer paths may impact frame rate

## Alternatives Considered

1. **Procedural level generation** - Rejected: Would make difficulty unpredictable
2. **Level unlock system** - Rejected: Adds complexity, better to let players self-select difficulty
3. **Fewer levels (5)** - Rejected: 10 provides better granularity for progression
