package entities

// Position represents a 2D coordinate on the game grid
type Position struct {
	X, Y float64
}

// IntPos returns the integer position for grid placement
func (p Position) IntPos() (int, int) {
	return int(p.X), int(p.Y)
}

// Distance calculates the distance to another position
func (p Position) Distance(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return dx*dx + dy*dy // squared distance for efficiency
}

// EntityType identifies the kind of entity
type EntityType int

const (
	EntityEnemy EntityType = iota
	EntityTower
	EntityProjectile
)

// TowerType identifies different tower categories
type TowerType int

const (
	TowerArrow TowerType = iota
	TowerLSP
	TowerRefactor
	TowerTelescope
	TowerMacro
	TowerGit
)

// TowerInfo contains configuration for each tower type
type TowerInfo struct {
	Name       string
	Cost       int
	Damage     int
	Range      float64
	Cooldown   float64 // seconds between attacks
	Category   string  // challenge category
	Symbol     string  // display character
	Color      string  // hex color
	Upgrades   []TowerUpgrade
}

// TowerUpgrade defines an upgrade tier
type TowerUpgrade struct {
	Cost         int
	DamageBonus  int
	RangeBonus   float64
	CooldownMult float64 // multiplier (0.8 = 20% faster)
}

// EnemyType identifies different enemy variants
type EnemyType int

const (
	EnemyBug EnemyType = iota
	EnemyGremlin
	EnemyDaemon
	EnemyBoss
)

// EnemyInfo contains configuration for each enemy type
type EnemyInfo struct {
	Name      string
	Health    int
	Speed     float64 // cells per second
	Symbol    string
	Color     string
	GoldValue int
}

// TowerTypes contains all tower configurations
var TowerTypes = map[TowerType]TowerInfo{
	TowerArrow: {
		Name:     "Arrow",
		Cost:     50,
		Damage:   10,
		Range:    3.0,
		Cooldown: 1.0,
		Category: "movement",
		Symbol:   "🏹",
		Color:    "#22c55e",
		Upgrades: []TowerUpgrade{
			{Cost: 30, DamageBonus: 5, RangeBonus: 0.5, CooldownMult: 0.9},
			{Cost: 60, DamageBonus: 10, RangeBonus: 1.0, CooldownMult: 0.8},
		},
	},
	TowerLSP: {
		Name:     "LSP",
		Cost:     100,
		Damage:   25,
		Range:    5.0,
		Cooldown: 2.0,
		Category: "lsp-navigation",
		Symbol:   "🔮",
		Color:    "#8b5cf6",
		Upgrades: []TowerUpgrade{
			{Cost: 60, DamageBonus: 15, RangeBonus: 1.0, CooldownMult: 0.85},
			{Cost: 120, DamageBonus: 30, RangeBonus: 2.0, CooldownMult: 0.7},
		},
	},
	TowerRefactor: {
		Name:     "Refactor",
		Cost:     150,
		Damage:   15,
		Range:    2.5,
		Cooldown: 1.5,
		Category: "text-objects",
		Symbol:   "⚡",
		Color:    "#f59e0b",
		Upgrades: []TowerUpgrade{
			{Cost: 80, DamageBonus: 10, RangeBonus: 0.5, CooldownMult: 0.9},
			{Cost: 150, DamageBonus: 20, RangeBonus: 1.0, CooldownMult: 0.75},
		},
	},
}

// EnemyTypes contains all enemy configurations
var EnemyTypes = map[EnemyType]EnemyInfo{
	EnemyBug: {
		Name:      "Bug",
		Health:    10,
		Speed:     1.5,
		Symbol:    "🐛",
		Color:     "#ef4444",
		GoldValue: 5,
	},
	EnemyGremlin: {
		Name:      "Gremlin",
		Health:    25,
		Speed:     2.5,
		Symbol:    "👹",
		Color:     "#f97316",
		GoldValue: 10,
	},
	EnemyDaemon: {
		Name:      "Daemon",
		Health:    100,
		Speed:     0.8,
		Symbol:    "👿",
		Color:     "#dc2626",
		GoldValue: 25,
	},
	EnemyBoss: {
		Name:      "Boss",
		Health:    500,
		Speed:     0.5,
		Symbol:    "💀",
		Color:     "#7c2d12",
		GoldValue: 100,
	},
}
