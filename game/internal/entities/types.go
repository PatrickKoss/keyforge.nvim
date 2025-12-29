package entities

// Position represents a 2D coordinate on the game grid.
type Position struct {
	X, Y float64
}

// IntPos returns the integer position for grid placement.
func (p Position) IntPos() (int, int) {
	return int(p.X), int(p.Y)
}

// Distance calculates the distance to another position.
func (p Position) Distance(other Position) float64 {
	dx := p.X - other.X
	dy := p.Y - other.Y
	return dx*dx + dy*dy // squared distance for efficiency
}

// EntityType identifies the kind of entity.
type EntityType int

const (
	EntityEnemy EntityType = iota
	EntityTower
	EntityProjectile
)

// TowerType identifies different tower categories.
type TowerType int

const (
	TowerArrow TowerType = iota
	TowerLSP
	TowerRefactor
	TowerTelescope
	TowerMacro
	TowerGit
)

// TowerInfo contains configuration for each tower type.
type TowerInfo struct {
	Name       string
	Cost       int
	Damage     int
	Range      float64
	Cooldown   float64  // seconds between attacks
	Category   string   // primary challenge category (for backwards compat)
	Categories []string // all challenge categories this tower can trigger
	Symbol     string   // display character
	Color      string   // hex color
	Upgrades   []TowerUpgrade
}

// TowerUpgrade defines an upgrade tier.
type TowerUpgrade struct {
	Cost         int
	DamageBonus  int
	RangeBonus   float64
	CooldownMult float64 // multiplier (0.8 = 20% faster)
}

// EnemyType identifies different enemy variants.
type EnemyType int

const (
	EnemyMite    EnemyType = iota // Very weak, fast - early game fodder
	EnemyBug                      // Baseline enemy
	EnemyGremlin                  // Fast, medium health
	EnemyCrawler                  // Slow tank, high health
	EnemySpecter                  // Very fast, fragile
	EnemyDaemon                   // Late-game tank
	EnemyBoss                     // Final challenge
)

// EnemyInfo contains configuration for each enemy type.
type EnemyInfo struct {
	Name      string
	Health    int
	Speed     float64 // cells per second
	Symbol    string
	Color     string
	GoldValue int
}

// TowerTypes contains all tower configurations.
// Rebalanced: cost correlates with power, each tower has distinct role.
// Categories mapping:
// - Arrow: movement, buffer-management, window-management, quickfix, folding
// - LSP: lsp-navigation, telescope, diagnostics, formatting, harpoon
// - Refactor: text-objects, search-replace, refactoring, surround, git-operations.
var TowerTypes = map[TowerType]TowerInfo{
	TowerArrow: {
		Name:       "Arrow",
		Cost:       50,
		Damage:     8,   // Fast attacker, lower damage
		Range:      2.5, // Short range
		Cooldown:   0.8, // Fast attack speed
		Category:   "movement",
		Categories: []string{"movement", "buffer-management", "window-management", "quickfix", "folding"},
		Symbol:     "🏹",
		Color:      "#22c55e",
		Upgrades: []TowerUpgrade{
			{Cost: 30, DamageBonus: 1, RangeBonus: 0.3, CooldownMult: 0.9}, // +15% dmg (1.2 -> rounds to 1), +0.3 range, -10% cooldown
			{Cost: 60, DamageBonus: 2, RangeBonus: 0.3, CooldownMult: 0.9}, // Cumulative
		},
	},
	TowerLSP: {
		Name:       "LSP",
		Cost:       100,
		Damage:     20,  // Sniper, high damage
		Range:      5.0, // Long range
		Cooldown:   1.5, // Slower attack
		Category:   "lsp-navigation",
		Categories: []string{"lsp-navigation", "telescope", "diagnostics", "formatting", "harpoon"},
		Symbol:     "🔮",
		Color:      "#8b5cf6",
		Upgrades: []TowerUpgrade{
			{Cost: 60, DamageBonus: 3, RangeBonus: 0.3, CooldownMult: 0.9},  // +15% dmg (3), +0.3 range, -10% cooldown
			{Cost: 120, DamageBonus: 3, RangeBonus: 0.3, CooldownMult: 0.9}, // Cumulative
		},
	},
	TowerRefactor: {
		Name:       "Refactor",
		Cost:       150,
		Damage:     12,  // Balanced area damage
		Range:      3.0, // Medium range
		Cooldown:   1.0, // Medium attack speed
		Category:   "text-objects",
		Categories: []string{"text-objects", "search-replace", "refactoring", "surround", "git-operations"},
		Symbol:     "⚡",
		Color:      "#f59e0b",
		Upgrades: []TowerUpgrade{
			{Cost: 90, DamageBonus: 2, RangeBonus: 0.3, CooldownMult: 0.9},  // +15% dmg (1.8 -> 2), +0.3 range, -10% cooldown
			{Cost: 180, DamageBonus: 2, RangeBonus: 0.3, CooldownMult: 0.9}, // Cumulative
		},
	},
}

// EnemyTypes contains all enemy configurations.
var EnemyTypes = map[EnemyType]EnemyInfo{
	EnemyMite: {
		Name:      "Mite",
		Health:    5,
		Speed:     2.0,
		Symbol:    "🦟",
		Color:     "#a3e635",
		GoldValue: 2,
	},
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
	EnemyCrawler: {
		Name:      "Crawler",
		Health:    40,
		Speed:     0.6,
		Symbol:    "🐌",
		Color:     "#78716c",
		GoldValue: 15,
	},
	EnemySpecter: {
		Name:      "Specter",
		Health:    15,
		Speed:     3.5,
		Symbol:    "👻",
		Color:     "#c4b5fd",
		GoldValue: 8,
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
