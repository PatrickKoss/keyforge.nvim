package engine

import (
	"testing"

	"github.com/keyforge/keyforge/internal/entities"
)

func TestClassicLevel(t *testing.T) {
	level := ClassicLevel()

	t.Run("has correct ID and name", func(t *testing.T) {
		// ClassicLevel now returns Level5
		if level.ID != "level-5" {
			t.Errorf("Expected ID 'level-5', got '%s'", level.ID)
		}
		if level.Name != "Classic" {
			t.Errorf("Expected Name 'Classic', got '%s'", level.Name)
		}
	})

	t.Run("has valid grid dimensions", func(t *testing.T) {
		if level.GridWidth < 20 {
			t.Errorf("Expected GridWidth >= 20, got %d", level.GridWidth)
		}
		if level.GridHeight != 14 {
			t.Errorf("Expected GridHeight 14, got %d", level.GridHeight)
		}
	})

	t.Run("has non-empty path", func(t *testing.T) {
		if len(level.Path) == 0 {
			t.Error("Expected non-empty path")
		}
		// Classic path should have 35 positions
		if len(level.Path) != 35 {
			t.Errorf("Expected 35 path positions, got %d", len(level.Path))
		}
	})

	t.Run("path starts at left edge", func(t *testing.T) {
		if len(level.Path) == 0 {
			t.Skip("Path is empty")
		}
		start := level.Path[0]
		if start.X != 0 {
			t.Errorf("Expected path to start at X=0, got X=%.0f", start.X)
		}
	})

	t.Run("has waves", func(t *testing.T) {
		if level.TotalWaves < 5 {
			t.Errorf("Expected at least 5 waves, got %d", level.TotalWaves)
		}
	})

	t.Run("has wave function", func(t *testing.T) {
		if level.WaveFunc == nil {
			t.Error("Expected WaveFunc to be set")
		}
	})

	t.Run("wave function returns valid waves", func(t *testing.T) {
		if level.WaveFunc == nil {
			t.Skip("WaveFunc is nil")
		}
		for waveNum := 1; waveNum <= level.TotalWaves; waveNum++ {
			wave := level.WaveFunc(waveNum)
			if len(wave.Spawns) == 0 {
				t.Errorf("Wave %d has no spawns", waveNum)
			}
			if wave.BonusGold <= 0 {
				t.Errorf("Wave %d has no bonus gold", waveNum)
			}
		}
	})

	t.Run("has allowed towers", func(t *testing.T) {
		if len(level.AllowedTowers) == 0 {
			t.Error("Expected at least one allowed tower")
		}
		// Classic should have Arrow, LSP, Refactor
		expectedTowers := []entities.TowerType{
			entities.TowerArrow,
			entities.TowerLSP,
			entities.TowerRefactor,
		}
		if len(level.AllowedTowers) != len(expectedTowers) {
			t.Errorf("Expected %d allowed towers, got %d", len(expectedTowers), len(level.AllowedTowers))
		}
	})

	t.Run("has enemy types", func(t *testing.T) {
		if len(level.EnemyTypes) == 0 {
			t.Error("Expected at least one enemy type")
		}
	})
}

func TestLevelRegistry(t *testing.T) {
	registry := NewLevelRegistry()

	t.Run("has exactly 10 levels", func(t *testing.T) {
		if registry.Count() != 10 {
			t.Errorf("Expected exactly 10 levels, got %d", registry.Count())
		}
	})

	t.Run("GetAll returns all 10 levels", func(t *testing.T) {
		levels := registry.GetAll()
		if len(levels) != 10 {
			t.Errorf("Expected GetAll to return 10 levels, got %d", len(levels))
		}
	})

	t.Run("GetByID finds all levels by ID", func(t *testing.T) {
		expectedIDs := []string{
			"level-1", "level-2", "level-3", "level-4", "level-5",
			"level-6", "level-7", "level-8", "level-9", "level-10",
		}
		for _, id := range expectedIDs {
			level := registry.GetByID(id)
			if level == nil {
				t.Errorf("Expected to find level with ID %s", id)
			}
		}
	})

	t.Run("GetByID returns nil for unknown level", func(t *testing.T) {
		level := registry.GetByID("nonexistent")
		if level != nil {
			t.Error("Expected nil for unknown level ID")
		}
	})

	t.Run("Count matches GetAll length", func(t *testing.T) {
		levels := registry.GetAll()
		if registry.Count() != len(levels) {
			t.Errorf("Count() = %d, but GetAll() returned %d levels", registry.Count(), len(levels))
		}
	})

	t.Run("levels have unique IDs", func(t *testing.T) {
		levels := registry.GetAll()
		seen := make(map[string]bool)
		for _, level := range levels {
			if seen[level.ID] {
				t.Errorf("Duplicate level ID: %s", level.ID)
			}
			seen[level.ID] = true
		}
	})
}

func TestClassicPath(t *testing.T) {
	path := classicPath()

	t.Run("path is continuous", func(t *testing.T) {
		for i := 1; i < len(path); i++ {
			prev := path[i-1]
			curr := path[i]
			dx := curr.X - prev.X
			dy := curr.Y - prev.Y
			// Each step should be exactly 1 unit in one direction
			if (dx != 0 && dy != 0) || (dx == 0 && dy == 0) {
				t.Errorf("Path is not continuous at index %d: prev=(%.0f,%.0f), curr=(%.0f,%.0f)",
					i, prev.X, prev.Y, curr.X, curr.Y)
			}
			if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
				t.Errorf("Path has a gap at index %d", i)
			}
		}
	})

	t.Run("path stays within grid bounds", func(t *testing.T) {
		for i, pos := range path {
			if pos.X < 0 || pos.X >= 20 {
				t.Errorf("Path position %d has X=%.0f out of bounds [0,20)", i, pos.X)
			}
			if pos.Y < 0 || pos.Y >= 14 {
				t.Errorf("Path position %d has Y=%.0f out of bounds [0,14)", i, pos.Y)
			}
		}
	})
}

func TestLevelDifficulty(t *testing.T) {
	tests := []struct {
		difficulty LevelDifficulty
		expected   string
	}{
		{LevelDifficultyBeginner, "beginner"},
		{LevelDifficultyIntermediate, "intermediate"},
		{LevelDifficultyAdvanced, "advanced"},
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			if string(tc.difficulty) != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, tc.difficulty)
			}
		})
	}
}

func TestAllLevelPathsConnected(t *testing.T) {
	registry := NewLevelRegistry()
	levels := registry.GetAll()

	for _, level := range levels {
		t.Run(level.Name, func(t *testing.T) {
			path := level.Path
			if len(path) == 0 {
				t.Error("Path is empty")
				return
			}

			// Check each step is connected (adjacent cell)
			for i := 1; i < len(path); i++ {
				prev := path[i-1]
				curr := path[i]
				dx := curr.X - prev.X
				dy := curr.Y - prev.Y

				// Each step should be exactly 1 unit in one direction (no diagonal)
				if (dx != 0 && dy != 0) || (dx == 0 && dy == 0) {
					t.Errorf("Path not continuous at index %d: prev=(%.0f,%.0f), curr=(%.0f,%.0f)",
						i, prev.X, prev.Y, curr.X, curr.Y)
				}
				if dx < -1 || dx > 1 || dy < -1 || dy > 1 {
					t.Errorf("Path has gap at index %d", i)
				}
			}
		})
	}
}

func TestAllLevelPathsInBounds(t *testing.T) {
	registry := NewLevelRegistry()
	levels := registry.GetAll()

	for _, level := range levels {
		t.Run(level.Name, func(t *testing.T) {
			for i, pos := range level.Path {
				if pos.X < 0 || pos.X >= float64(level.GridWidth) {
					t.Errorf("Path position %d has X=%.0f out of bounds [0,%d)", i, pos.X, level.GridWidth)
				}
				if pos.Y < 0 || pos.Y >= float64(level.GridHeight) {
					t.Errorf("Path position %d has Y=%.0f out of bounds [0,%d)", i, pos.Y, level.GridHeight)
				}
			}
		})
	}
}

func TestLevelPathLengthProgression(t *testing.T) {
	registry := NewLevelRegistry()
	levels := registry.GetAll()

	// Expected minimum path lengths (approximately increasing with level)
	minLengths := map[string]int{
		"level-1":  10, // ~15 cells
		"level-2":  15, // ~20 cells
		"level-3":  20, // ~25 cells
		"level-4":  25, // ~30 cells
		"level-5":  30, // ~35 cells (classic)
		"level-6":  35, // ~40 cells
		"level-7":  40, // ~45 cells
		"level-8":  45, // ~50 cells
		"level-9":  50, // ~55 cells
		"level-10": 55, // ~60+ cells
	}

	for _, level := range levels {
		t.Run(level.Name, func(t *testing.T) {
			minLen := minLengths[level.ID]
			if len(level.Path) < minLen {
				t.Errorf("Level %s path too short: expected >= %d, got %d",
					level.ID, minLen, len(level.Path))
			}
		})
	}

	// Verify general trend: later levels have longer paths
	t.Run("path length increases with level", func(t *testing.T) {
		level1 := registry.GetByID("level-1")
		level5 := registry.GetByID("level-5")
		level10 := registry.GetByID("level-10")

		if len(level5.Path) <= len(level1.Path) {
			t.Errorf("Level 5 path (%d) should be longer than level 1 (%d)",
				len(level5.Path), len(level1.Path))
		}
		if len(level10.Path) <= len(level5.Path) {
			t.Errorf("Level 10 path (%d) should be longer than level 5 (%d)",
				len(level10.Path), len(level5.Path))
		}
	})
}

func TestLevelEnemyPools(t *testing.T) {
	// Test that each level has the correct enemy pool as designed
	tests := []struct {
		levelID       string
		expectedTypes []entities.EnemyType
	}{
		{"level-1", []entities.EnemyType{entities.EnemyMite, entities.EnemyBug}},
		{"level-2", []entities.EnemyType{entities.EnemyMite, entities.EnemyBug}},
		{"level-3", []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin}},
		{"level-4", []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin, entities.EnemyCrawler}},
		{"level-5", []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin, entities.EnemySpecter}},
		{"level-6", []entities.EnemyType{entities.EnemyGremlin, entities.EnemyCrawler, entities.EnemyDaemon}},
		{"level-7", []entities.EnemyType{entities.EnemyGremlin, entities.EnemySpecter, entities.EnemyDaemon}},
		{"level-8", []entities.EnemyType{entities.EnemyCrawler, entities.EnemySpecter, entities.EnemyDaemon}},
		{"level-9", []entities.EnemyType{entities.EnemySpecter, entities.EnemyDaemon}},
		{"level-10", []entities.EnemyType{
			entities.EnemyBug, entities.EnemyGremlin, entities.EnemyCrawler,
			entities.EnemySpecter, entities.EnemyDaemon, entities.EnemyBoss,
		}},
	}

	registry := NewLevelRegistry()

	for _, tc := range tests {
		t.Run(tc.levelID, func(t *testing.T) {
			level := registry.GetByID(tc.levelID)
			if level == nil {
				t.Fatalf("Level %s not found", tc.levelID)
			}

			if len(level.EnemyTypes) != len(tc.expectedTypes) {
				t.Errorf("Expected %d enemy types, got %d",
					len(tc.expectedTypes), len(level.EnemyTypes))
			}

			// Check each expected type is present
			for _, expected := range tc.expectedTypes {
				found := false
				for _, actual := range level.EnemyTypes {
					if actual == expected {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected enemy type %v not found in level", expected)
				}
			}
		})
	}
}

func TestLevelDifficultyProgression(t *testing.T) {
	registry := NewLevelRegistry()

	// Verify difficulty progression matches design
	expectedDifficulties := map[string]LevelDifficulty{
		"level-1":  LevelDifficultyBeginner,
		"level-2":  LevelDifficultyBeginner,
		"level-3":  LevelDifficultyBeginner,
		"level-4":  LevelDifficultyIntermediate,
		"level-5":  LevelDifficultyIntermediate,
		"level-6":  LevelDifficultyIntermediate,
		"level-7":  LevelDifficultyAdvanced,
		"level-8":  LevelDifficultyAdvanced,
		"level-9":  LevelDifficultyAdvanced,
		"level-10": LevelDifficultyAdvanced,
	}

	for levelID, expectedDiff := range expectedDifficulties {
		t.Run(levelID, func(t *testing.T) {
			level := registry.GetByID(levelID)
			if level == nil {
				t.Fatalf("Level %s not found", levelID)
			}
			if level.Difficulty != expectedDiff {
				t.Errorf("Expected difficulty %s, got %s", expectedDiff, level.Difficulty)
			}
		})
	}
}

func TestLevelWaveFunctions(t *testing.T) {
	registry := NewLevelRegistry()
	levels := registry.GetAll()

	for _, level := range levels {
		t.Run(level.Name, func(t *testing.T) {
			if level.WaveFunc == nil {
				t.Error("WaveFunc is nil")
				return
			}

			// Test all waves generate valid data
			for waveNum := 1; waveNum <= level.TotalWaves; waveNum++ {
				wave := level.WaveFunc(waveNum)

				if wave.Number != waveNum {
					t.Errorf("Wave %d has wrong number: %d", waveNum, wave.Number)
				}
				if len(wave.Spawns) == 0 {
					t.Errorf("Wave %d has no spawns", waveNum)
				}
				if len(wave.Spawns) > 10 {
					t.Errorf("Wave %d has too many spawns: %d (max 10)", waveNum, len(wave.Spawns))
				}
				if wave.BonusGold <= 0 {
					t.Errorf("Wave %d has no bonus gold", waveNum)
				}
			}
		})
	}
}

func TestLevelWaveEnemiesMatchPool(t *testing.T) {
	registry := NewLevelRegistry()
	levels := registry.GetAll()

	for _, level := range levels {
		t.Run(level.Name, func(t *testing.T) {
			if level.WaveFunc == nil {
				t.Skip("WaveFunc is nil")
			}

			// Build a set of allowed enemy types
			allowed := make(map[entities.EnemyType]bool)
			for _, et := range level.EnemyTypes {
				allowed[et] = true
			}

			// Check all waves only spawn allowed enemies
			for waveNum := 1; waveNum <= level.TotalWaves; waveNum++ {
				wave := level.WaveFunc(waveNum)
				for _, spawn := range wave.Spawns {
					if !allowed[spawn.Type] {
						t.Errorf("Wave %d spawns %v which is not in level's enemy pool",
							waveNum, spawn.Type)
					}
				}
			}
		})
	}
}
