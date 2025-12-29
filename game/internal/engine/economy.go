package engine

import "math"

// EconomyConfig holds configuration for the game's gold economy.
type EconomyConfig struct {
	MobGoldMultiplier     float64 // 0.25 = 25% of original mob gold
	WaveBonusMultiplier   float64 // 0.50 = 50% of original wave bonus
	ChallengeBaseGold     int     // Base gold for difficulty 1 challenges
	ChallengeSpeedMaxMult float64 // Max speed bonus multiplier (2.0)
}

// Difficulty presets.
const (
	DifficultyEasy   = "easy"
	DifficultyNormal = "normal"
	DifficultyHard   = "hard"
)

// DefaultEconomyConfig returns the normal difficulty economy config.
func DefaultEconomyConfig() EconomyConfig {
	return EconomyConfig{
		MobGoldMultiplier:     0.25,
		WaveBonusMultiplier:   0.50,
		ChallengeBaseGold:     25,
		ChallengeSpeedMaxMult: 2.0,
	}
}

// EconomyConfigForDifficulty returns economy config for a given difficulty preset.
func EconomyConfigForDifficulty(difficulty string) EconomyConfig {
	switch difficulty {
	case DifficultyEasy:
		return EconomyConfig{
			MobGoldMultiplier:     0.50, // 50% mob gold
			WaveBonusMultiplier:   0.75, // 75% wave bonus
			ChallengeBaseGold:     25,
			ChallengeSpeedMaxMult: 2.0,
		}
	case DifficultyHard:
		return EconomyConfig{
			MobGoldMultiplier:     0.0,  // No mob gold
			WaveBonusMultiplier:   0.25, // 25% wave bonus
			ChallengeBaseGold:     25,
			ChallengeSpeedMaxMult: 2.0,
		}
	default: // Normal
		return DefaultEconomyConfig()
	}
}

// CalculateMobGold calculates the gold reward for killing an enemy.
func (e EconomyConfig) CalculateMobGold(baseGold int) int {
	gold := int(math.Round(float64(baseGold) * e.MobGoldMultiplier))
	if gold < 1 && e.MobGoldMultiplier > 0 {
		gold = 1 // Minimum 1 gold if multiplier is non-zero
	}
	return gold
}

// CalculateWaveBonus calculates the gold bonus for completing a wave.
func (e EconomyConfig) CalculateWaveBonus(baseBonus int) int {
	bonus := int(math.Round(float64(baseBonus) * e.WaveBonusMultiplier))
	if bonus < 1 && e.WaveBonusMultiplier > 0 {
		bonus = 1 // Minimum 1 gold if multiplier is non-zero
	}
	return bonus
}

// CalculateSpeedBonus calculates the speed bonus multiplier for challenge completion
// Returns a multiplier between 1.0 and ChallengeSpeedMaxMult.
func (e EconomyConfig) CalculateSpeedBonus(timeMs, parTimeMs int) float64 {
	if timeMs <= 0 || parTimeMs <= 0 {
		return 1.0
	}

	// If completed slower than par, no bonus
	if timeMs >= parTimeMs {
		return 1.0
	}

	// Calculate speed ratio (how much faster than par)
	speedRatio := float64(parTimeMs) / float64(timeMs)

	// Apply formula: 1.0 + (speedRatio - 1.0) * 0.5
	// This gives a gentler curve: 2x faster = 1.5x bonus, not 2x
	bonus := 1.0 + (speedRatio-1.0)*0.5

	// Cap at max multiplier
	if bonus > e.ChallengeSpeedMaxMult {
		bonus = e.ChallengeSpeedMaxMult
	}

	return bonus
}

// CalculateChallengeGold calculates the total gold reward for a challenge.
func (e EconomyConfig) CalculateChallengeGold(baseGold, difficulty int, efficiency, speedBonus float64) int {
	// Difficulty multiplier: 1.0 + (difficulty * 0.25)
	// d1 = 1.25, d2 = 1.5, d3 = 1.75
	difficultyMult := 1.0 + float64(difficulty)*0.25

	// Efficiency multiplier: 0.5 + (efficiency * 0.5)
	// 50% base + up to 50% for efficiency
	efficiencyMult := 0.5 + efficiency*0.5

	// Total gold
	gold := float64(baseGold) * difficultyMult * efficiencyMult * speedBonus

	// Minimum 1 gold for any successful completion
	result := int(math.Floor(gold))
	if result < 1 {
		result = 1
	}

	return result
}
