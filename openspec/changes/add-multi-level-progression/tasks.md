# Tasks: Add Multi-Level Progression System

## Phase 1: Enemy Type Expansion

- [x] **1.1** Add new enemy type constants (EnemyMite, EnemyCrawler, EnemySpecter) to `entities/types.go`
- [x] **1.2** Add EnemyInfo configurations for new enemies with health, speed, gold, symbol, color
- [x] **1.3** Add enemy rendering in `ui/styles.go` (colors) and `ui/view.go` (render cases)
- [x] **1.4** Write unit tests for new enemy types verifying stats and gold values
- [x] **1.5** Manually test enemy spawning and movement in standalone mode

## Phase 2: Tower Stat Rebalancing

- [x] **2.1** Update TowerInfo stats in `entities/types.go` (damage, range, cooldown per design)
- [x] **2.2** Update TowerUpgrade configs for balanced scaling (+15% damage, +0.3 range, -10% cooldown)
- [x] **2.3** Write unit tests verifying tower DPS calculations at each upgrade tier
- [x] **2.4** Playtest tower balance with existing Classic level

## Phase 3: Range Visualization

- [x] **3.1** Add `ShowingRange` field to Game struct to track range display state
- [x] **3.2** Implement `renderRangeCircle()` function using Bresenham circle algorithm
- [x] **3.3** Call range rendering in `renderGridCursor()` during placement mode
- [x] **3.4** Show range when cursor hovers over existing tower
- [x] **3.5** Add distinct visual style (semi-transparent, dotted) for range overlay
- [x] **3.6** Write unit tests for range circle calculation (cells within range)

## Phase 4: Level Definitions

- [x] **4.1** Create path definitions for levels 1-10 in `engine/level.go`
  - Level 1: Straight line (~15 cells)
  - Level 2: Single L-turn (~20 cells)
  - Level 3: S-curve (~25 cells)
  - Level 4: Zigzag (~30 cells)
  - Level 5: Classic path (existing, ~35 cells)
  - Level 6: Spiral-in (~40 cells)
  - Level 7: Maze-lite (~45 cells)
  - Level 8: Snake (~50 cells)
  - Level 9: Complex winding (~55 cells)
  - Level 10: Ultimate challenge (~60+ cells)
- [x] **4.2** Create Level structs for each level with ID, Name, Description, Difficulty
- [x] **4.3** Register all 10 levels in `NewLevelRegistry()`
- [x] **4.4** Write unit tests verifying path connectivity for all levels
- [x] **4.5** Write unit tests verifying level registry contains exactly 10 levels

## Phase 5: Per-Level Wave Generation

- [x] **5.1** Create wave generator functions for each level in `engine/wave.go`
- [x] **5.2** Implement enemy pool selection based on level (use design table)
- [x] **5.3** Implement weighted enemy selection favoring harder enemies in later waves
- [x] **5.4** Ensure 3-7 enemies per wave across all levels
- [x] **5.5** Configure bonus gold scaling per wave and level
- [x] **5.6** Write unit tests for wave generation verifying enemy counts and types
- [x] **5.7** Write unit tests verifying enemy pool restrictions per level

## Phase 6: Integration & Polish

- [x] **6.1** Update level browser UI to show all 10 levels with difficulty indicators
- [x] **6.2** Update level preview to show enemy icons for each level
- [x] **6.3** Verify start screen navigation works with 10 levels
- [x] **6.4** Run full test suite (`make test`)
- [x] **6.5** Run linter (`make lint`) and fix any issues
- [ ] **6.6** Playtest complete progression from level 1 through 10
- [ ] **6.7** Performance test with complex paths to verify 60fps target

## Dependencies

- Phase 1 (enemies) can run in parallel with Phase 2 (towers)
- Phase 3 (range viz) can run in parallel with Phase 1 and 2
- Phase 4 (levels) depends on Phase 1 (new enemy types must exist)
- Phase 5 (waves) depends on Phase 1 and Phase 4
- Phase 6 (integration) depends on all prior phases

## Acceptance Criteria

### Unit Tests (Required)

All unit tests in `game/internal/` using Go's `testing` package:

| Test File | Test Cases | Description |
|-----------|------------|-------------|
| `entities/types_test.go` | `TestEnemyTypes` | Verify all 7 enemy types have correct health, speed, gold values |
| `entities/types_test.go` | `TestEnemyTypeCount` | Assert exactly 7 enemy types exist |
| `entities/types_test.go` | `TestTowerStats` | Verify rebalanced tower damage, range, cooldown |
| `entities/types_test.go` | `TestTowerUpgradeScaling` | Verify +15% damage, +0.3 range, -10% cooldown per tier |
| `engine/level_test.go` | `TestLevelRegistryCount` | Assert exactly 10 levels in registry |
| `engine/level_test.go` | `TestLevelPathConnectivity` | For each level, verify all waypoints are adjacent |
| `engine/level_test.go` | `TestLevelPathLengths` | Verify path lengths increase with level number |
| `engine/level_test.go` | `TestLevelEnemyPools` | Verify each level's enemy pool matches design table |
| `engine/level_test.go` | `TestLevelDifficulty` | Verify difficulty progression (beginner/intermediate/advanced) |
| `engine/wave_test.go` | `TestWaveEnemyCount` | Verify 3-7 enemies per wave for all levels |
| `engine/wave_test.go` | `TestWaveEnemyPoolRestriction` | Verify only allowed enemies spawn per level |
| `engine/wave_test.go` | `TestWaveGoldScaling` | Verify bonus gold increases with wave number |
| `ui/range_test.go` | `TestRangeCircleCalculation` | Verify cells within range are correctly identified |
| `ui/range_test.go` | `TestRangeCircleBounds` | Verify range doesn't exceed grid boundaries |

### Integration Tests (Required)

Integration tests verifying component interactions:

| Test File | Test Cases | Description |
|-----------|------------|-------------|
| `engine/game_test.go` | `TestGameWithLevel` | Create game from each level, verify initialization |
| `engine/game_test.go` | `TestLevelWaveProgression` | Run through all waves of levels 1, 5, 10 |
| `engine/game_test.go` | `TestEnemySpawnFromLevel` | Spawn enemies and verify only allowed types appear |
| `engine/game_test.go` | `TestTowerDamageOnEnemies` | Place tower, spawn enemy, verify damage application |
| `engine/game_test.go` | `TestGoldEconomyBalance` | Complete waves, verify gold distribution (60-70% from challenges) |

### Validation Criteria (Manual)

1. All 10 levels appear in start screen level selector
2. Each level has a unique, visually distinct path
3. Enemy variety per level matches design table (3-7 types)
4. Tower range shows during placement and on hover
5. All existing tests pass (`make test` returns 0)
6. Linter passes (`make lint` returns 0)
7. Game runs at 60fps on all levels
8. Each level is completable (manual playtest)
