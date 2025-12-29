package engine

import (
	"math"
	"testing"
)

func TestDefaultEconomyConfig(t *testing.T) {
	config := DefaultEconomyConfig()

	if config.MobGoldMultiplier != 0.25 {
		t.Errorf("expected MobGoldMultiplier 0.25, got %v", config.MobGoldMultiplier)
	}
	if config.WaveBonusMultiplier != 0.50 {
		t.Errorf("expected WaveBonusMultiplier 0.50, got %v", config.WaveBonusMultiplier)
	}
	if config.ChallengeBaseGold != 25 {
		t.Errorf("expected ChallengeBaseGold 25, got %v", config.ChallengeBaseGold)
	}
	if config.ChallengeSpeedMaxMult != 2.0 {
		t.Errorf("expected ChallengeSpeedMaxMult 2.0, got %v", config.ChallengeSpeedMaxMult)
	}
}

func TestEconomyConfigForDifficulty(t *testing.T) {
	tests := []struct {
		difficulty       string
		expectedMobMult  float64
		expectedWaveMult float64
	}{
		{DifficultyEasy, 0.50, 0.75},
		{DifficultyNormal, 0.25, 0.50},
		{DifficultyHard, 0.0, 0.25},
		{"unknown", 0.25, 0.50}, // defaults to normal
	}

	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			config := EconomyConfigForDifficulty(tt.difficulty)
			if config.MobGoldMultiplier != tt.expectedMobMult {
				t.Errorf("expected MobGoldMultiplier %v, got %v", tt.expectedMobMult, config.MobGoldMultiplier)
			}
			if config.WaveBonusMultiplier != tt.expectedWaveMult {
				t.Errorf("expected WaveBonusMultiplier %v, got %v", tt.expectedWaveMult, config.WaveBonusMultiplier)
			}
		})
	}
}

func TestCalculateMobGold(t *testing.T) {
	tests := []struct {
		name       string
		multiplier float64
		baseGold   int
		expected   int
	}{
		{"normal difficulty - bug", 0.25, 5, 1},
		{"normal difficulty - gremlin", 0.25, 10, 3}, // 10 * 0.25 = 2.5 -> 3 (rounded)
		{"normal difficulty - daemon", 0.25, 25, 6},  // 25 * 0.25 = 6.25 -> 6 (rounded)
		{"normal difficulty - boss", 0.25, 100, 25},
		{"easy difficulty - bug", 0.50, 5, 3}, // 5 * 0.50 = 2.5 -> 3 (rounded)
		{"easy difficulty - boss", 0.50, 100, 50},
		{"hard difficulty - all zero", 0.0, 100, 0},
		{"minimum 1 gold", 0.25, 1, 1}, // 1 * 0.25 = 0.25 -> min 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := EconomyConfig{MobGoldMultiplier: tt.multiplier}
			result := config.CalculateMobGold(tt.baseGold)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateWaveBonus(t *testing.T) {
	tests := []struct {
		name       string
		multiplier float64
		baseBonus  int
		expected   int
	}{
		{"normal difficulty - wave 1", 0.50, 25, 13}, // 25 * 0.50 = 12.5 -> 13 (rounded)
		{"normal difficulty - wave 10", 0.50, 200, 100},
		{"easy difficulty - wave 1", 0.75, 25, 19}, // 25 * 0.75 = 18.75 -> 19 (rounded)
		{"hard difficulty - wave 1", 0.25, 25, 6},  // 25 * 0.25 = 6.25 -> 6 (rounded)
		{"minimum 1 gold", 0.25, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := EconomyConfig{WaveBonusMultiplier: tt.multiplier}
			result := config.CalculateWaveBonus(tt.baseBonus)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateSpeedBonus(t *testing.T) {
	config := DefaultEconomyConfig()

	tests := []struct {
		name     string
		timeMs   int
		parMs    int
		expected float64
	}{
		{"at par time", 5000, 5000, 1.0},
		{"slower than par", 7000, 5000, 1.0},
		{"2x faster", 2500, 5000, 1.5},          // formula: (5000/2500 - 1) * 0.5 + 1
		{"4x faster - capped", 1250, 5000, 2.0}, // would be 2.5 but capped at 2.0
		{"zero time", 0, 5000, 1.0},
		{"zero par", 5000, 0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.CalculateSpeedBonus(tt.timeMs, tt.parMs)
			if math.Abs(result-tt.expected) > 0.001 {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCalculateChallengeGold(t *testing.T) {
	config := DefaultEconomyConfig()

	tests := []struct {
		name       string
		baseGold   int
		difficulty int
		efficiency float64
		speedBonus float64
		expected   int
	}{
		{
			name:       "perfect run - difficulty 1",
			baseGold:   50,
			difficulty: 1,
			efficiency: 1.0,
			speedBonus: 1.0,
			expected:   62, // 50 * 1.25 * 1.0 * 1.0 = 62.5 -> 62
		},
		{
			name:       "perfect run with speed bonus",
			baseGold:   50,
			difficulty: 1,
			efficiency: 1.0,
			speedBonus: 2.0,
			expected:   125, // base * diffMult * effBonus * speed
		},
		{
			name:       "50% efficiency",
			baseGold:   50,
			difficulty: 1,
			efficiency: 0.5,
			speedBonus: 1.0,
			expected:   46, // 46.875 rounded down
		},
		{
			name:       "difficulty 2",
			baseGold:   50,
			difficulty: 2,
			efficiency: 1.0,
			speedBonus: 1.0,
			expected:   75, // base * 1.5 diffMult
		},
		{
			name:       "difficulty 3",
			baseGold:   50,
			difficulty: 3,
			efficiency: 1.0,
			speedBonus: 1.0,
			expected:   87, // 87.5 rounded down
		},
		{
			name:       "minimum 1 gold",
			baseGold:   1,
			difficulty: 1,
			efficiency: 0.0,
			speedBonus: 1.0,
			expected:   1, // would be 0.625, but min 1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.CalculateChallengeGold(tt.baseGold, tt.difficulty, tt.efficiency, tt.speedBonus)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGameUsesEconomyConfig(t *testing.T) {
	// Test that NewGame uses default economy
	game := NewGame(20, 14)
	if game.Economy.MobGoldMultiplier != 0.25 {
		t.Errorf("expected default economy config")
	}

	// Test that NewGameWithEconomy uses provided config
	customEconomy := EconomyConfigForDifficulty(DifficultyHard)
	game2 := NewGameWithEconomy(20, 14, customEconomy)
	if game2.Economy.MobGoldMultiplier != 0.0 {
		t.Errorf("expected hard difficulty economy config")
	}
}
