package engine

import "github.com/keyforge/keyforge/internal/entities"

// LevelDifficulty indicates the skill level required for a level.
type LevelDifficulty string

const (
	LevelDifficultyBeginner     LevelDifficulty = "beginner"
	LevelDifficultyIntermediate LevelDifficulty = "intermediate"
	LevelDifficultyAdvanced     LevelDifficulty = "advanced"
)

// WaveFunc is a function that generates wave configuration for a given wave number.
type WaveFunc func(waveNum int) Wave

// Level defines a game level with its path, waves, and available entities.
type Level struct {
	ID            string
	Name          string
	Description   string
	GridWidth     int
	GridHeight    int
	Path          []entities.Position
	TotalWaves    int
	WaveFunc      WaveFunc
	AllowedTowers []entities.TowerType
	EnemyTypes    []entities.EnemyType
	Difficulty    LevelDifficulty
}

// LevelRegistry holds all available levels.
type LevelRegistry struct {
	levels []Level
}

// NewLevelRegistry creates a registry with all built-in levels.
func NewLevelRegistry() *LevelRegistry {
	return &LevelRegistry{
		levels: []Level{
			ClassicLevel(),
		},
	}
}

// GetAll returns all available levels.
func (r *LevelRegistry) GetAll() []Level {
	return r.levels
}

// GetByID returns a level by its ID, or nil if not found.
func (r *LevelRegistry) GetByID(id string) *Level {
	for i := range r.levels {
		if r.levels[i].ID == id {
			return &r.levels[i]
		}
	}
	return nil
}

// Count returns the number of available levels.
func (r *LevelRegistry) Count() int {
	return len(r.levels)
}

// ClassicLevel returns the default level using the original game configuration.
func ClassicLevel() Level {
	return Level{
		ID:          "classic",
		Name:        "Classic",
		Description: "The original Keyforge experience. Master vim motions to defend against 10 waves of bugs.",
		GridWidth:   20,
		GridHeight:  14,
		Path:        classicPath(),
		TotalWaves:  10,
		WaveFunc:    GetWave, // Uses existing wave generation
		AllowedTowers: []entities.TowerType{
			entities.TowerArrow,
			entities.TowerLSP,
			entities.TowerRefactor,
		},
		EnemyTypes: []entities.EnemyType{
			entities.EnemyBug,
			entities.EnemyGremlin,
			entities.EnemyDaemon,
			entities.EnemyBoss,
		},
		Difficulty: LevelDifficultyBeginner,
	}
}

// classicPath returns the S-shaped path from the original game.
func classicPath() []entities.Position {
	return []entities.Position{
		{X: 0, Y: 3},
		{X: 1, Y: 3},
		{X: 2, Y: 3},
		{X: 3, Y: 3},
		{X: 4, Y: 3},
		{X: 5, Y: 3},
		{X: 6, Y: 3},
		{X: 7, Y: 3},
		{X: 8, Y: 3},
		{X: 8, Y: 4},
		{X: 8, Y: 5},
		{X: 8, Y: 6},
		{X: 8, Y: 7},
		{X: 7, Y: 7},
		{X: 6, Y: 7},
		{X: 5, Y: 7},
		{X: 4, Y: 7},
		{X: 3, Y: 7},
		{X: 2, Y: 7},
		{X: 2, Y: 8},
		{X: 2, Y: 9},
		{X: 2, Y: 10},
		{X: 3, Y: 10},
		{X: 4, Y: 10},
		{X: 5, Y: 10},
		{X: 6, Y: 10},
		{X: 7, Y: 10},
		{X: 8, Y: 10},
		{X: 9, Y: 10},
		{X: 10, Y: 10},
		{X: 11, Y: 10},
		{X: 12, Y: 10},
		{X: 13, Y: 10},
		{X: 14, Y: 10},
		{X: 15, Y: 10},
	}
}
