# Design: Multi-Level Progression System

## Overview

This change spans multiple systems: level definitions, enemy types, tower stats, and UI rendering. The design prioritizes minimal code changes while enabling significant gameplay expansion.

## Architecture

### Component Relationships

```
┌─────────────────────────────────────────────────────────────┐
│                      Level Registry                          │
│  - Holds all 10 level definitions                           │
│  - Each level references: path, enemy types, wave configs   │
└─────────────────────┬───────────────────────────────────────┘
                      │
          ┌───────────┼───────────┐
          ▼           ▼           ▼
┌─────────────┐ ┌───────────┐ ┌─────────────┐
│ Path Defs   │ │ Enemy     │ │ Wave        │
│ (10 unique) │ │ Types (7) │ │ Generators  │
└─────────────┘ └───────────┘ └─────────────┘
                      │
                      ▼
┌─────────────────────────────────────────────────────────────┐
│                    Tower System                              │
│  - Rebalanced stats (damage, range, cooldown)               │
│  - Range visualization during placement/hover               │
└─────────────────────────────────────────────────────────────┘
```

### File Changes

```
game/internal/
├── entities/
│   └── types.go           # Add 3 new enemy types, rebalance tower stats
├── engine/
│   ├── level.go           # Add 9 new level definitions
│   └── wave.go            # Add per-level wave generators
└── ui/
    └── view.go            # Add range visualization rendering
```

## Implementation Details

### 1. Enemy Type Expansion

**Current enemies (entities/types.go):**
- Bug, Gremlin, Daemon, Boss

**New enemies to add:**
```go
const (
    EnemyMite    EnemyType = iota // New: very weak, fast
    EnemyBug                       // Existing
    EnemyGremlin                   // Existing
    EnemyCrawler                   // New: slow tank
    EnemySpecter                   // New: very fast, fragile
    EnemyDaemon                    // Existing
    EnemyBoss                      // Existing
)
```

**Enemy stat profiles:**

| Type | Health | Speed | Gold | Design Intent |
|------|--------|-------|------|---------------|
| Mite | 5 | 2.0 | 2 | Early fodder, teaches basic targeting |
| Bug | 10 | 1.5 | 5 | Baseline enemy |
| Gremlin | 25 | 2.5 | 10 | Speed challenge |
| Crawler | 40 | 0.6 | 15 | Tank, tests sustained DPS |
| Specter | 15 | 3.5 | 8 | Requires quick tower response |
| Daemon | 100 | 0.8 | 25 | Late-game tank |
| Boss | 500 | 0.5 | 100 | Final challenge |

### 2. Tower Stat Rebalancing

**Design philosophy:**
- Cost correlates with overall power
- Trade-offs between range, damage, and speed
- Each tower has a distinct role

**Rebalanced stats:**

```go
TowerArrow: {
    Cost:     50,
    Damage:   8,     // Down from 10
    Range:    2.5,   // Down from 3.0
    Cooldown: 0.8,   // Down from 1.0 (faster)
    // Role: Fast attacker, short range
}

TowerLSP: {
    Cost:     100,
    Damage:   20,    // Down from 25
    Range:    5.0,   // Same
    Cooldown: 1.5,   // Down from 2.0
    // Role: Long-range sniper
}

TowerRefactor: {
    Cost:     150,
    Damage:   12,    // Down from 15
    Range:    3.0,   // Up from 2.5
    Cooldown: 1.0,   // Down from 1.5
    // Role: Balanced area damage
}
```

**Upgrade progression (per tier):**
```go
TowerUpgrade{
    Cost:         baseCost * 0.6,  // 60% of tower cost
    DamageBonus:  baseDamage * 0.15,  // +15%
    RangeBonus:   0.3,
    CooldownMult: 0.9,  // 10% faster
}
```

### 3. Level Definitions

Each level is defined as a `Level` struct with:
- Unique path waypoints
- Wave count (scales with difficulty)
- Allowed enemy types (subset based on level)
- Custom wave generator function

**Level progression:**

| Level | Grid | Waves | Path Desc | Enemy Pool |
|-------|------|-------|-----------|------------|
| 1 | 20x14 | 5 | Straight | Mite, Bug |
| 2 | 20x14 | 6 | L-turn | Mite, Bug |
| 3 | 20x14 | 7 | S-curve | Bug, Gremlin |
| 4 | 22x14 | 7 | Zigzag | Bug, Gremlin, Crawler |
| 5 | 22x14 | 8 | Loop | Bug, Gremlin, Specter |
| 6 | 24x14 | 8 | Spiral-in | Gremlin, Crawler, Daemon |
| 7 | 24x14 | 9 | Maze-lite | Gremlin, Specter, Daemon |
| 8 | 26x14 | 9 | Snake | Crawler, Specter, Daemon |
| 9 | 26x14 | 10 | Complex | Specter, Daemon |
| 10 | 28x14 | 10 | Ultimate | All (Boss on wave 10) |

### 4. Wave Generation Per Level

Each level has a custom `WaveFunc` that:
1. Selects from the level's allowed enemy pool
2. Scales enemy count: 3-7 enemies per wave
3. Increases difficulty within the level (later waves harder)

**Wave generation algorithm:**
```go
func levelNWave(waveNum int) Wave {
    pool := levelNEnemyPool
    count := baseCount + (waveNum * scaling)  // 3-7 range

    // Weight selection toward harder enemies in later waves
    weights := calculateWeights(waveNum, pool)
    spawns := generateSpawns(count, pool, weights)

    return Wave{
        Number:    waveNum,
        Spawns:    spawns,
        BonusGold: baseBonusGold + (waveNum * goldScale),
    }
}
```

### 5. Range Visualization

**During tower placement:**
- Draw circular indicator centered on cursor
- Use dotted/dashed style to distinguish from solid game elements
- Color indicates valid (green) or invalid (red) placement

**On tower hover:**
- When cursor is on an existing tower, show its range
- Same visual style as placement preview

**Implementation approach:**
```go
// In renderGridCursor, after cursor rendering:
if g.State == engine.StatePlaying {
    if placing {
        renderRangeCircle(grid, g.CursorX, g.CursorY, selectedTowerRange)
    } else if tower := g.GetTowerAt(g.CursorX, g.CursorY); tower != nil {
        renderRangeCircle(grid, tower.Pos, tower.CurrentRange())
    }
}
```

**Range circle rendering:**
- Use Bresenham's circle algorithm to determine cells within range
- Overlay a subtle background color on affected cells
- Don't overwrite entities, only empty/path cells

## Trade-offs

### Keeping Wave Count Low Per Level
- **Pro:** Faster level completion, more levels tried
- **Con:** Less time to build economy per level
- **Decision:** 5-10 waves per level, with higher gold from later enemies

### All Levels Unlocked
- **Pro:** Players can jump to preferred difficulty
- **Con:** No sense of achievement from unlocking
- **Decision:** Simplicity wins; vim skill is the real unlock mechanism

### No New Tower Types
- **Pro:** Keeps balance simpler, less testing needed
- **Con:** Limits strategic variety
- **Decision:** Focus on enemy variety instead; towers can be expanded later

## Testing Strategy

1. **Unit tests for new enemy types** - Verify stats, gold values
2. **Level path validation** - Ensure all paths are connected, no gaps
3. **Wave balance testing** - Playtest each level for completability
4. **Range visualization** - Visual verification of circle accuracy
5. **Performance testing** - Ensure 60fps with complex paths and more enemies

## Migration

No migration needed. The Classic level becomes Level 5 (similar difficulty), and all new levels are added to the registry. Existing saves/state are not affected since we don't persist level progress.
