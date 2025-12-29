package engine

import (
	"testing"

	"github.com/keyforge/keyforge/internal/entities"
)

func TestClassicLevel(t *testing.T) {
	level := ClassicLevel()

	t.Run("has correct ID and name", func(t *testing.T) {
		if level.ID != "classic" {
			t.Errorf("Expected ID 'classic', got '%s'", level.ID)
		}
		if level.Name != "Classic" {
			t.Errorf("Expected Name 'Classic', got '%s'", level.Name)
		}
	})

	t.Run("has valid grid dimensions", func(t *testing.T) {
		if level.GridWidth != 20 {
			t.Errorf("Expected GridWidth 20, got %d", level.GridWidth)
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

	t.Run("has 10 waves", func(t *testing.T) {
		if level.TotalWaves != 10 {
			t.Errorf("Expected 10 waves, got %d", level.TotalWaves)
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
		// Classic should have Bug, Gremlin, Daemon, Boss
		if len(level.EnemyTypes) != 4 {
			t.Errorf("Expected 4 enemy types, got %d", len(level.EnemyTypes))
		}
	})

	t.Run("has beginner difficulty", func(t *testing.T) {
		if level.Difficulty != LevelDifficultyBeginner {
			t.Errorf("Expected beginner difficulty, got %s", level.Difficulty)
		}
	})
}

func TestLevelRegistry(t *testing.T) {
	registry := NewLevelRegistry()

	t.Run("has at least one level", func(t *testing.T) {
		if registry.Count() == 0 {
			t.Error("Expected at least one level in registry")
		}
	})

	t.Run("GetAll returns levels", func(t *testing.T) {
		levels := registry.GetAll()
		if len(levels) == 0 {
			t.Error("Expected GetAll to return levels")
		}
	})

	t.Run("GetByID finds classic level", func(t *testing.T) {
		level := registry.GetByID("classic")
		if level == nil {
			t.Error("Expected to find classic level")
		}
		if level != nil && level.Name != "Classic" {
			t.Errorf("Expected Classic level, got %s", level.Name)
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
