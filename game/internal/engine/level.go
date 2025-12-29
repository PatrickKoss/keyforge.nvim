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
			Level1(),
			Level2(),
			Level3(),
			Level4(),
			Level5(), // Classic level
			Level6(),
			Level7(),
			Level8(),
			Level9(),
			Level10(),
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

// AllTowers returns all available tower types.
var AllTowers = []entities.TowerType{
	entities.TowerArrow,
	entities.TowerLSP,
	entities.TowerRefactor,
}

// Level1 - Straight path, easiest level for beginners.
func Level1() Level {
	return Level{
		ID:            "level-1",
		Name:          "The Straight Path",
		Description:   "A simple straight path. Perfect for learning the basics.",
		GridWidth:     20,
		GridHeight:    14,
		Path:          level1Path(),
		TotalWaves:    5,
		WaveFunc:      level1Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyMite, entities.EnemyBug},
		Difficulty:    LevelDifficultyBeginner,
	}
}

func level1Path() []entities.Position {
	// Straight horizontal path (~15 cells)
	path := make([]entities.Position, 0, 15)
	for x := range 15 {
		path = append(path, entities.Position{X: float64(x), Y: 7})
	}
	return path
}

// Level2 - L-turn path, introduces turning.
func Level2() Level {
	return Level{
		ID:            "level-2",
		Name:          "The Corner",
		Description:   "An L-shaped path with one turn. Watch the corner!",
		GridWidth:     20,
		GridHeight:    14,
		Path:          level2Path(),
		TotalWaves:    6,
		WaveFunc:      level2Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyMite, entities.EnemyBug},
		Difficulty:    LevelDifficultyBeginner,
	}
}

func level2Path() []entities.Position {
	// L-shaped path (~20 cells)
	path := make([]entities.Position, 0, 20)
	// Horizontal segment
	for x := range 10 {
		path = append(path, entities.Position{X: float64(x), Y: 3})
	}
	// Vertical segment
	for y := 4; y < 14; y++ {
		path = append(path, entities.Position{X: 9, Y: float64(y)})
	}
	return path
}

// Level3 - S-curve path, more complex pathing.
func Level3() Level {
	return Level{
		ID:            "level-3",
		Name:          "The Serpent",
		Description:   "A winding S-curve. Enemies take their time.",
		GridWidth:     20,
		GridHeight:    14,
		Path:          level3Path(),
		TotalWaves:    7,
		WaveFunc:      level3Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin},
		Difficulty:    LevelDifficultyBeginner,
	}
}

func level3Path() []entities.Position {
	// S-curve path (~25 cells)
	path := make([]entities.Position, 0, 25)
	// Top horizontal
	for x := range 8 {
		path = append(path, entities.Position{X: float64(x), Y: 3})
	}
	// Down
	for y := 4; y < 8; y++ {
		path = append(path, entities.Position{X: 7, Y: float64(y)})
	}
	// Left
	for x := 6; x >= 2; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 7})
	}
	// Down
	for y := 8; y < 12; y++ {
		path = append(path, entities.Position{X: 2, Y: float64(y)})
	}
	// Right to exit
	for x := 3; x < 10; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 11})
	}
	return path
}

// Level4 - Zigzag path, introduces tank enemies.
func Level4() Level {
	return Level{
		ID:            "level-4",
		Name:          "Zigzag Valley",
		Description:   "A zigzag path through the valley. Tank enemies appear!",
		GridWidth:     22,
		GridHeight:    14,
		Path:          level4Path(),
		TotalWaves:    7,
		WaveFunc:      level4Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin, entities.EnemyCrawler},
		Difficulty:    LevelDifficultyIntermediate,
	}
}

func level4Path() []entities.Position {
	// Zigzag path (~30 cells)
	path := make([]entities.Position, 0, 30)
	// First leg right
	for x := range 6 {
		path = append(path, entities.Position{X: float64(x), Y: 2})
	}
	// Down
	for y := 3; y < 6; y++ {
		path = append(path, entities.Position{X: 5, Y: float64(y)})
	}
	// Left
	for x := 4; x >= 2; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 5})
	}
	// Down
	for y := 6; y < 9; y++ {
		path = append(path, entities.Position{X: 2, Y: float64(y)})
	}
	// Right
	for x := 3; x < 8; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 8})
	}
	// Down
	for y := 9; y < 12; y++ {
		path = append(path, entities.Position{X: 7, Y: float64(y)})
	}
	// Right to exit
	for x := 8; x < 15; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 11})
	}
	return path
}

// Level5 - Classic level (the original S-path).
func Level5() Level {
	return Level{
		ID:            "level-5",
		Name:          "Classic",
		Description:   "The original Keyforge experience. Fast enemies join the fray!",
		GridWidth:     22,
		GridHeight:    14,
		Path:          classicPath(),
		TotalWaves:    8,
		WaveFunc:      level5Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyBug, entities.EnemyGremlin, entities.EnemySpecter},
		Difficulty:    LevelDifficultyIntermediate,
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

// Level6 - Spiral-in path, daemons appear.
func Level6() Level {
	return Level{
		ID:            "level-6",
		Name:          "The Spiral",
		Description:   "A spiral path toward the center. Daemons emerge!",
		GridWidth:     24,
		GridHeight:    14,
		Path:          level6Path(),
		TotalWaves:    8,
		WaveFunc:      level6Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyGremlin, entities.EnemyCrawler, entities.EnemyDaemon},
		Difficulty:    LevelDifficultyIntermediate,
	}
}

func level6Path() []entities.Position {
	// Spiral-in path (~40 cells)
	path := make([]entities.Position, 0, 40)
	// Outer ring - top
	for x := range 12 {
		path = append(path, entities.Position{X: float64(x), Y: 1})
	}
	// Outer ring - right
	for y := 2; y < 12; y++ {
		path = append(path, entities.Position{X: 11, Y: float64(y)})
	}
	// Outer ring - bottom (going left)
	for x := 10; x >= 3; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 11})
	}
	// Inner approach - up
	for y := 10; y >= 5; y-- {
		path = append(path, entities.Position{X: 3, Y: float64(y)})
	}
	// Inner approach - right to center
	for x := 4; x < 8; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 5})
	}
	return path
}

// Level7 - Maze-lite path, speed and power mix.
func Level7() Level {
	return Level{
		ID:            "level-7",
		Name:          "The Labyrinth",
		Description:   "A maze-like path. Speed meets power!",
		GridWidth:     24,
		GridHeight:    14,
		Path:          level7Path(),
		TotalWaves:    9,
		WaveFunc:      level7Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyGremlin, entities.EnemySpecter, entities.EnemyDaemon},
		Difficulty:    LevelDifficultyAdvanced,
	}
}

func level7Path() []entities.Position {
	// Maze-lite path (~45 cells)
	path := make([]entities.Position, 0, 45)
	// Entry horizontal
	for x := range 5 {
		path = append(path, entities.Position{X: float64(x), Y: 2})
	}
	// Down
	for y := 3; y < 7; y++ {
		path = append(path, entities.Position{X: 4, Y: float64(y)})
	}
	// Right
	for x := 5; x < 10; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 6})
	}
	// Up
	for y := 5; y >= 2; y-- {
		path = append(path, entities.Position{X: 9, Y: float64(y)})
	}
	// Right
	for x := 10; x < 15; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 2})
	}
	// Down
	for y := 3; y < 10; y++ {
		path = append(path, entities.Position{X: 14, Y: float64(y)})
	}
	// Left
	for x := 13; x >= 8; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 9})
	}
	// Down to exit
	for y := 10; y < 14; y++ {
		path = append(path, entities.Position{X: 8, Y: float64(y)})
	}
	return path
}

// Level8 - Snake path, late game enemy mix.
func Level8() Level {
	return Level{
		ID:            "level-8",
		Name:          "The Serpent's Lair",
		Description:   "A long snaking path. Only the strong survive.",
		GridWidth:     26,
		GridHeight:    14,
		Path:          level8Path(),
		TotalWaves:    9,
		WaveFunc:      level8Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemyCrawler, entities.EnemySpecter, entities.EnemyDaemon},
		Difficulty:    LevelDifficultyAdvanced,
	}
}

func level8Path() []entities.Position {
	// Snake path (~50 cells)
	path := make([]entities.Position, 0, 50)
	// Row 1 - right
	for x := range 12 {
		path = append(path, entities.Position{X: float64(x), Y: 2})
	}
	// Down
	for y := 3; y < 5; y++ {
		path = append(path, entities.Position{X: 11, Y: float64(y)})
	}
	// Row 2 - left
	for x := 10; x >= 2; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 4})
	}
	// Down
	for y := 5; y < 7; y++ {
		path = append(path, entities.Position{X: 2, Y: float64(y)})
	}
	// Row 3 - right
	for x := 3; x < 14; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 6})
	}
	// Down
	for y := 7; y < 9; y++ {
		path = append(path, entities.Position{X: 13, Y: float64(y)})
	}
	// Row 4 - left
	for x := 12; x >= 4; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 8})
	}
	// Down
	for y := 9; y < 11; y++ {
		path = append(path, entities.Position{X: 4, Y: float64(y)})
	}
	// Final row - right to exit
	for x := 5; x < 16; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 10})
	}
	return path
}

// Level9 - Complex winding path, pre-boss.
func Level9() Level {
	return Level{
		ID:            "level-9",
		Name:          "The Gauntlet",
		Description:   "A complex winding path. Prepare for the final challenge!",
		GridWidth:     26,
		GridHeight:    14,
		Path:          level9Path(),
		TotalWaves:    10,
		WaveFunc:      level9Wave,
		AllowedTowers: AllTowers,
		EnemyTypes:    []entities.EnemyType{entities.EnemySpecter, entities.EnemyDaemon},
		Difficulty:    LevelDifficultyAdvanced,
	}
}

func level9Path() []entities.Position {
	// Complex winding path (~55 cells)
	path := make([]entities.Position, 0, 55)
	// Start - horizontal
	for x := range 6 {
		path = append(path, entities.Position{X: float64(x), Y: 1})
	}
	// Down
	for y := 2; y < 6; y++ {
		path = append(path, entities.Position{X: 5, Y: float64(y)})
	}
	// Right
	for x := 6; x < 12; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 5})
	}
	// Up
	for y := 4; y >= 1; y-- {
		path = append(path, entities.Position{X: 11, Y: float64(y)})
	}
	// Right
	for x := 12; x < 18; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 1})
	}
	// Down
	for y := 2; y < 8; y++ {
		path = append(path, entities.Position{X: 17, Y: float64(y)})
	}
	// Left
	for x := 16; x >= 8; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 7})
	}
	// Down
	for y := 8; y < 11; y++ {
		path = append(path, entities.Position{X: 8, Y: float64(y)})
	}
	// Right
	for x := 9; x < 15; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 10})
	}
	// Down to exit
	for y := 11; y < 14; y++ {
		path = append(path, entities.Position{X: 14, Y: float64(y)})
	}
	return path
}

// Level10 - Ultimate challenge, all enemies including boss.
func Level10() Level {
	return Level{
		ID:            "level-10",
		Name:          "The Ultimate Challenge",
		Description:   "The final test. Face all enemies and the Boss!",
		GridWidth:     28,
		GridHeight:    14,
		Path:          level10Path(),
		TotalWaves:    10,
		WaveFunc:      level10Wave,
		AllowedTowers: AllTowers,
		EnemyTypes: []entities.EnemyType{
			entities.EnemyBug,
			entities.EnemyGremlin,
			entities.EnemyCrawler,
			entities.EnemySpecter,
			entities.EnemyDaemon,
			entities.EnemyBoss,
		},
		Difficulty: LevelDifficultyAdvanced,
	}
}

func level10Path() []entities.Position {
	// Ultimate path (~60+ cells)
	path := make([]entities.Position, 0, 65)
	// Outer perimeter start - top
	for x := range 14 {
		path = append(path, entities.Position{X: float64(x), Y: 1})
	}
	// Down right side
	for y := 2; y < 6; y++ {
		path = append(path, entities.Position{X: 13, Y: float64(y)})
	}
	// Left
	for x := 12; x >= 4; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 5})
	}
	// Down
	for y := 6; y < 9; y++ {
		path = append(path, entities.Position{X: 4, Y: float64(y)})
	}
	// Right
	for x := 5; x < 16; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 8})
	}
	// Up
	for y := 7; y >= 3; y-- {
		path = append(path, entities.Position{X: 15, Y: float64(y)})
	}
	// Right
	for x := 16; x < 22; x++ {
		path = append(path, entities.Position{X: float64(x), Y: 3})
	}
	// Down
	for y := 4; y < 12; y++ {
		path = append(path, entities.Position{X: 21, Y: float64(y)})
	}
	// Left to exit
	for x := 20; x >= 10; x-- {
		path = append(path, entities.Position{X: float64(x), Y: 11})
	}
	return path
}

// ClassicLevel returns Level5 for backward compatibility.
func ClassicLevel() Level {
	return Level5()
}
