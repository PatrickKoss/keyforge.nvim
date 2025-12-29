package engine

import (
	"testing"
)

func TestDefaultGameSettings(t *testing.T) {
	settings := DefaultGameSettings()

	t.Run("has normal difficulty", func(t *testing.T) {
		if settings.Difficulty != DifficultyNormal {
			t.Errorf("Expected normal difficulty, got %s", settings.Difficulty)
		}
	})

	t.Run("has normal speed", func(t *testing.T) {
		if settings.GameSpeed != SpeedNormal {
			t.Errorf("Expected 1x speed, got %v", settings.GameSpeed)
		}
	})

	t.Run("has 200 starting gold", func(t *testing.T) {
		if settings.StartingGold != 200 {
			t.Errorf("Expected 200 starting gold, got %d", settings.StartingGold)
		}
	})

	t.Run("has 100 starting health", func(t *testing.T) {
		if settings.StartingHealth != 100 {
			t.Errorf("Expected 100 starting health, got %d", settings.StartingHealth)
		}
	})
}

func TestGameSettingsValidate(t *testing.T) {
	t.Run("clamps gold below minimum", func(t *testing.T) {
		settings := GameSettings{StartingGold: 50}
		settings.Validate()
		if settings.StartingGold != 100 {
			t.Errorf("Expected gold clamped to 100, got %d", settings.StartingGold)
		}
	})

	t.Run("clamps gold above maximum", func(t *testing.T) {
		settings := GameSettings{StartingGold: 1000}
		settings.Validate()
		if settings.StartingGold != 500 {
			t.Errorf("Expected gold clamped to 500, got %d", settings.StartingGold)
		}
	})

	t.Run("accepts valid gold", func(t *testing.T) {
		settings := GameSettings{StartingGold: 300}
		settings.Validate()
		if settings.StartingGold != 300 {
			t.Errorf("Expected gold unchanged at 300, got %d", settings.StartingGold)
		}
	})

	t.Run("clamps health below minimum", func(t *testing.T) {
		settings := GameSettings{StartingHealth: 25}
		settings.Validate()
		if settings.StartingHealth != 50 {
			t.Errorf("Expected health clamped to 50, got %d", settings.StartingHealth)
		}
	})

	t.Run("clamps health above maximum", func(t *testing.T) {
		settings := GameSettings{StartingHealth: 500}
		settings.Validate()
		if settings.StartingHealth != 200 {
			t.Errorf("Expected health clamped to 200, got %d", settings.StartingHealth)
		}
	})

	t.Run("accepts valid health", func(t *testing.T) {
		settings := GameSettings{StartingHealth: 150}
		settings.Validate()
		if settings.StartingHealth != 150 {
			t.Errorf("Expected health unchanged at 150, got %d", settings.StartingHealth)
		}
	})

	t.Run("corrects invalid difficulty", func(t *testing.T) {
		settings := GameSettings{Difficulty: "invalid"}
		settings.Validate()
		if settings.Difficulty != DifficultyNormal {
			t.Errorf("Expected difficulty corrected to normal, got %s", settings.Difficulty)
		}
	})

	t.Run("accepts easy difficulty", func(t *testing.T) {
		settings := GameSettings{Difficulty: DifficultyEasy}
		settings.Validate()
		if settings.Difficulty != DifficultyEasy {
			t.Errorf("Expected easy difficulty preserved, got %s", settings.Difficulty)
		}
	})

	t.Run("accepts hard difficulty", func(t *testing.T) {
		settings := GameSettings{Difficulty: DifficultyHard}
		settings.Validate()
		if settings.Difficulty != DifficultyHard {
			t.Errorf("Expected hard difficulty preserved, got %s", settings.Difficulty)
		}
	})

	t.Run("corrects invalid game speed", func(t *testing.T) {
		settings := GameSettings{GameSpeed: 3.0}
		settings.Validate()
		if settings.GameSpeed != SpeedNormal {
			t.Errorf("Expected speed corrected to 1x, got %v", settings.GameSpeed)
		}
	})

	t.Run("accepts valid game speeds", func(t *testing.T) {
		speeds := GameSpeedOptions()
		for _, speed := range speeds {
			settings := GameSettings{GameSpeed: speed}
			settings.Validate()
			if settings.GameSpeed != speed {
				t.Errorf("Expected speed %v preserved, got %v", speed, settings.GameSpeed)
			}
		}
	})
}

func TestGameSpeed(t *testing.T) {
	t.Run("SpeedHalf is 0.5", func(t *testing.T) {
		if SpeedHalf != 0.5 {
			t.Errorf("Expected SpeedHalf = 0.5, got %v", SpeedHalf)
		}
	})

	t.Run("SpeedNormal is 1.0", func(t *testing.T) {
		if SpeedNormal != 1.0 {
			t.Errorf("Expected SpeedNormal = 1.0, got %v", SpeedNormal)
		}
	})

	t.Run("SpeedFast is 1.5", func(t *testing.T) {
		if SpeedFast != 1.5 {
			t.Errorf("Expected SpeedFast = 1.5, got %v", SpeedFast)
		}
	})

	t.Run("SpeedDouble is 2.0", func(t *testing.T) {
		if SpeedDouble != 2.0 {
			t.Errorf("Expected SpeedDouble = 2.0, got %v", SpeedDouble)
		}
	})
}

func TestGameSpeedString(t *testing.T) {
	tests := []struct {
		speed    GameSpeed
		expected string
	}{
		{SpeedHalf, "0.5x"},
		{SpeedNormal, "1x"},
		{SpeedFast, "1.5x"},
		{SpeedDouble, "2x"},
		{GameSpeed(0.75), "1x"}, // Invalid speed defaults to 1x
	}

	for _, tc := range tests {
		t.Run(tc.expected, func(t *testing.T) {
			result := tc.speed.String()
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestGameSpeedOptions(t *testing.T) {
	options := GameSpeedOptions()

	t.Run("returns 4 options", func(t *testing.T) {
		if len(options) != 4 {
			t.Errorf("Expected 4 speed options, got %d", len(options))
		}
	})

	t.Run("options are in order", func(t *testing.T) {
		expected := []GameSpeed{SpeedHalf, SpeedNormal, SpeedFast, SpeedDouble}
		for i, speed := range options {
			if speed != expected[i] {
				t.Errorf("Option %d: expected %v, got %v", i, expected[i], speed)
			}
		}
	})
}

func TestGetEconomyConfig(t *testing.T) {
	t.Run("easy difficulty returns easy economy", func(t *testing.T) {
		settings := GameSettings{Difficulty: DifficultyEasy}
		economy := settings.GetEconomyConfig()
		expected := EconomyConfigForDifficulty(DifficultyEasy)
		if economy.MobGoldMultiplier != expected.MobGoldMultiplier {
			t.Errorf("Expected MobGoldMultiplier %v, got %v",
				expected.MobGoldMultiplier, economy.MobGoldMultiplier)
		}
	})

	t.Run("normal difficulty returns normal economy", func(t *testing.T) {
		settings := GameSettings{Difficulty: DifficultyNormal}
		economy := settings.GetEconomyConfig()
		expected := EconomyConfigForDifficulty(DifficultyNormal)
		if economy.MobGoldMultiplier != expected.MobGoldMultiplier {
			t.Errorf("Expected MobGoldMultiplier %v, got %v",
				expected.MobGoldMultiplier, economy.MobGoldMultiplier)
		}
	})

	t.Run("hard difficulty returns hard economy", func(t *testing.T) {
		settings := GameSettings{Difficulty: DifficultyHard}
		economy := settings.GetEconomyConfig()
		expected := EconomyConfigForDifficulty(DifficultyHard)
		if economy.MobGoldMultiplier != expected.MobGoldMultiplier {
			t.Errorf("Expected MobGoldMultiplier %v, got %v",
				expected.MobGoldMultiplier, economy.MobGoldMultiplier)
		}
	})
}

func TestNewGameFromLevelAndSettings(t *testing.T) {
	level := ClassicLevel()
	settings := GameSettings{
		Difficulty:     DifficultyEasy,
		GameSpeed:      SpeedDouble,
		StartingGold:   300,
		StartingHealth: 150,
	}

	game := NewGameFromLevelAndSettings(&level, settings)

	t.Run("uses level grid dimensions", func(t *testing.T) {
		if game.Width != level.GridWidth {
			t.Errorf("Expected width %d, got %d", level.GridWidth, game.Width)
		}
		if game.Height != level.GridHeight {
			t.Errorf("Expected height %d, got %d", level.GridHeight, game.Height)
		}
	})

	t.Run("uses settings starting gold", func(t *testing.T) {
		if game.Gold != settings.StartingGold {
			t.Errorf("Expected gold %d, got %d", settings.StartingGold, game.Gold)
		}
	})

	t.Run("uses settings starting health", func(t *testing.T) {
		if game.Health != settings.StartingHealth {
			t.Errorf("Expected health %d, got %d", settings.StartingHealth, game.Health)
		}
		if game.MaxHealth != settings.StartingHealth {
			t.Errorf("Expected max health %d, got %d", settings.StartingHealth, game.MaxHealth)
		}
	})

	t.Run("uses settings game speed", func(t *testing.T) {
		if game.GameSpeed != settings.GameSpeed {
			t.Errorf("Expected speed %v, got %v", settings.GameSpeed, game.GameSpeed)
		}
	})

	t.Run("uses level path", func(t *testing.T) {
		if len(game.Path) != len(level.Path) {
			t.Errorf("Expected path length %d, got %d", len(level.Path), len(game.Path))
		}
	})

	t.Run("uses level total waves", func(t *testing.T) {
		if game.TotalWaves != level.TotalWaves {
			t.Errorf("Expected total waves %d, got %d", level.TotalWaves, game.TotalWaves)
		}
	})

	t.Run("uses level wave function", func(t *testing.T) {
		if game.WaveFunc == nil {
			t.Error("Expected WaveFunc to be set")
		}
	})

	t.Run("uses level allowed towers", func(t *testing.T) {
		if len(game.AllowedTowers) != len(level.AllowedTowers) {
			t.Errorf("Expected %d allowed towers, got %d",
				len(level.AllowedTowers), len(game.AllowedTowers))
		}
	})

	t.Run("selects first allowed tower", func(t *testing.T) {
		if game.SelectedTower != level.AllowedTowers[0] {
			t.Errorf("Expected selected tower %v, got %v",
				level.AllowedTowers[0], game.SelectedTower)
		}
	})

	t.Run("starts in playing state", func(t *testing.T) {
		if game.State != StatePlaying {
			t.Errorf("Expected StatePlaying, got %v", game.State)
		}
	})

	t.Run("applies economy from difficulty", func(t *testing.T) {
		expectedEconomy := EconomyConfigForDifficulty(settings.Difficulty)
		if game.Economy.MobGoldMultiplier != expectedEconomy.MobGoldMultiplier {
			t.Errorf("Expected economy MobGoldMultiplier %v, got %v",
				expectedEconomy.MobGoldMultiplier, game.Economy.MobGoldMultiplier)
		}
	})
}

func TestGameSpeedAffectsUpdate(t *testing.T) {
	level := ClassicLevel()

	t.Run("double speed processes faster", func(t *testing.T) {
		settingsNormal := GameSettings{
			Difficulty: DifficultyNormal,
			GameSpeed:  SpeedNormal,
		}
		settingsDouble := GameSettings{
			Difficulty: DifficultyNormal,
			GameSpeed:  SpeedDouble,
		}

		gameNormal := NewGameFromLevelAndSettings(&level, settingsNormal)
		gameDouble := NewGameFromLevelAndSettings(&level, settingsDouble)

		// Test wave countdown - it's predictable for speed testing
		gameNormal.WaveComplete = true
		gameNormal.WaveCountdown = 3.0
		gameDouble.WaveComplete = true
		gameDouble.WaveCountdown = 3.0

		// Update both for 1 second
		gameNormal.Update(1.0)
		gameDouble.Update(1.0)

		// Double speed should have reduced countdown by 2 seconds
		if gameNormal.WaveCountdown != 2.0 {
			t.Errorf("Normal speed: expected countdown 2.0, got %v", gameNormal.WaveCountdown)
		}
		if gameDouble.WaveCountdown != 1.0 {
			t.Errorf("Double speed: expected countdown 1.0, got %v", gameDouble.WaveCountdown)
		}
	})
}
