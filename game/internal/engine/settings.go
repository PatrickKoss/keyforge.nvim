package engine

// GameSpeed represents the game speed multiplier.
type GameSpeed float64

const (
	SpeedHalf   GameSpeed = 0.5
	SpeedNormal GameSpeed = 1.0
	SpeedFast   GameSpeed = 1.5
	SpeedDouble GameSpeed = 2.0
)

// GameSpeedOptions returns all available game speed options.
func GameSpeedOptions() []GameSpeed {
	return []GameSpeed{SpeedHalf, SpeedNormal, SpeedFast, SpeedDouble}
}

// String returns a human-readable label for the game speed.
func (s GameSpeed) String() string {
	switch s {
	case SpeedHalf:
		return "0.5x"
	case SpeedNormal:
		return "1x"
	case SpeedFast:
		return "1.5x"
	case SpeedDouble:
		return "2x"
	default:
		return "1x"
	}
}

// GameSettings holds all configurable game settings.
type GameSettings struct {
	Difficulty     string    // "easy", "normal", "hard"
	GameSpeed      GameSpeed // Time multiplier
	StartingGold   int       // Initial gold (100-500)
	StartingHealth int       // Initial health (50-200)
}

// DefaultGameSettings returns settings with sensible defaults.
func DefaultGameSettings() GameSettings {
	return GameSettings{
		Difficulty:     DifficultyNormal,
		GameSpeed:      SpeedNormal,
		StartingGold:   200,
		StartingHealth: 100,
	}
}

// Validate ensures settings are within valid ranges.
func (s *GameSettings) Validate() {
	// Validate difficulty
	switch s.Difficulty {
	case DifficultyEasy, DifficultyNormal, DifficultyHard:
		// Valid
	default:
		s.Difficulty = DifficultyNormal
	}

	// Validate game speed
	validSpeed := false
	for _, speed := range GameSpeedOptions() {
		if s.GameSpeed == speed {
			validSpeed = true
			break
		}
	}
	if !validSpeed {
		s.GameSpeed = SpeedNormal
	}

	// Clamp starting gold (100-500)
	if s.StartingGold < 100 {
		s.StartingGold = 100
	}
	if s.StartingGold > 500 {
		s.StartingGold = 500
	}

	// Clamp starting health (50-200)
	if s.StartingHealth < 50 {
		s.StartingHealth = 50
	}
	if s.StartingHealth > 200 {
		s.StartingHealth = 200
	}
}

// GetEconomyConfig returns the economy configuration based on difficulty.
func (s *GameSettings) GetEconomyConfig() EconomyConfig {
	return EconomyConfigForDifficulty(s.Difficulty)
}
