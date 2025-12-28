package engine

import (
	"github.com/keyforge/keyforge/internal/entities"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateChallengeActive  // Internal challenge (standalone mode) - game continues
	StateChallengeWaiting // Nvim challenge - game paused while user edits
	StateWaveComplete
	StateGameOver
	StateVictory
)

// Game holds all game state and logic
type Game struct {
	State      GameState
	Width      int
	Height     int
	Gold       int
	Health     int
	MaxHealth  int
	Wave       int
	TotalWaves int

	Path        []entities.Position
	Towers      []*entities.Tower
	Enemies     []*entities.Enemy
	Projectiles []*entities.Projectile
	Effects     *entities.EffectManager

	// Cursor for tower placement
	CursorX int
	CursorY int

	// Selected tower type for placement
	SelectedTower entities.TowerType

	// Wave management
	WaveTimer     float64 // time until next spawn
	SpawnIndex    int     // current spawn in wave
	WaveComplete  bool
	WaveCountdown float64 // countdown between waves

	// Economy configuration
	Economy EconomyConfig

	// Challenge state
	ChallengeActive bool // indicates a challenge is being solved

	// ID counters
	nextEnemyID int
	nextTowerID int
}

// NewGame creates a new game with default settings
func NewGame(width, height int) *Game {
	return NewGameWithEconomy(width, height, DefaultEconomyConfig())
}

// NewGameWithEconomy creates a new game with a specific economy configuration
func NewGameWithEconomy(width, height int, economy EconomyConfig) *Game {
	g := &Game{
		State:           StatePlaying,
		Width:           width,
		Height:          height,
		Gold:            200,
		Health:          100,
		MaxHealth:       100,
		Wave:            1,
		TotalWaves:      10,
		Towers:          make([]*entities.Tower, 0),
		Enemies:         make([]*entities.Enemy, 0),
		Projectiles:     make([]*entities.Projectile, 0),
		Effects:         entities.NewEffectManager(),
		CursorX:         width / 2,
		CursorY:         height / 2,
		SelectedTower:   entities.TowerArrow,
		WaveTimer:       0,
		SpawnIndex:      0,
		WaveComplete:    false,
		WaveCountdown:   3.0,
		Economy:         economy,
		ChallengeActive: false,
		nextEnemyID:     0,
		nextTowerID:     0,
	}
	g.Path = g.createDefaultPath()
	return g
}

// createDefaultPath creates a winding path across the map
func (g *Game) createDefaultPath() []entities.Position {
	// Create an S-shaped path
	path := []entities.Position{
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
	return path
}

// IsOnPath checks if a position is part of the path
func (g *Game) IsOnPath(x, y int) bool {
	for _, p := range g.Path {
		px, py := p.IntPos()
		if px == x && py == y {
			return true
		}
	}
	return false
}

// HasTower checks if there's a tower at the position
func (g *Game) HasTower(x, y int) bool {
	for _, t := range g.Towers {
		tx, ty := t.Pos.IntPos()
		if tx == x && ty == y {
			return true
		}
	}
	return false
}

// GetTowerAt returns the tower at the position, or nil
func (g *Game) GetTowerAt(x, y int) *entities.Tower {
	for _, t := range g.Towers {
		tx, ty := t.Pos.IntPos()
		if tx == x && ty == y {
			return t
		}
	}
	return nil
}

// CanPlaceTower checks if a tower can be placed at the position
func (g *Game) CanPlaceTower(x, y int) bool {
	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return false
	}
	if g.IsOnPath(x, y) {
		return false
	}
	if g.HasTower(x, y) {
		return false
	}
	return true
}

// PlaceTower attempts to place a tower at the cursor position
func (g *Game) PlaceTower() bool {
	info := entities.TowerTypes[g.SelectedTower]
	if g.Gold < info.Cost {
		return false
	}
	if !g.CanPlaceTower(g.CursorX, g.CursorY) {
		return false
	}

	g.nextTowerID++
	tower := entities.NewTower(g.nextTowerID, g.SelectedTower, entities.Position{
		X: float64(g.CursorX),
		Y: float64(g.CursorY),
	})
	g.Towers = append(g.Towers, tower)
	g.Gold -= info.Cost
	return true
}

// UpgradeTower attempts to upgrade the tower at cursor position
func (g *Game) UpgradeTower() bool {
	tower := g.GetTowerAt(g.CursorX, g.CursorY)
	if tower == nil || !tower.CanUpgrade() {
		return false
	}
	cost := tower.UpgradeCost()
	if g.Gold < cost {
		return false
	}
	tower.Upgrade()
	g.Gold -= cost
	return true
}

// SpawnEnemy spawns an enemy at the start of the path
func (g *Game) SpawnEnemy(enemyType entities.EnemyType) {
	if len(g.Path) == 0 {
		return
	}
	g.nextEnemyID++
	enemy := entities.NewEnemy(g.nextEnemyID, enemyType, g.Path[0])
	g.Enemies = append(g.Enemies, enemy)
}

// Update advances the game state by dt seconds
func (g *Game) Update(dt float64) {
	// Always update effects (even when paused for visual continuity)
	g.Effects.Update(dt)

	// Game continues during challenges (ChallengeActive state) but not when paused
	if g.State != StatePlaying && g.State != StateChallengeActive {
		return
	}

	// Update wave spawning
	g.updateWaveSpawning(dt)

	// Update enemies
	g.updateEnemies(dt)

	// Update towers and create projectiles
	g.updateTowers(dt)

	// Update projectiles and handle collisions
	g.updateProjectiles(dt)

	// Check win/lose conditions
	g.checkGameEnd()
}

// StartChallenge marks a challenge as active (game continues running)
// Used for standalone mode where internal vim editor handles the challenge
func (g *Game) StartChallenge() {
	if g.State == StatePlaying {
		g.State = StateChallengeActive
		g.ChallengeActive = true
	}
}

// StartChallengeWaiting marks a challenge as waiting (game paused)
// Used for nvim mode where user edits in a real Neovim buffer
func (g *Game) StartChallengeWaiting() {
	if g.State == StatePlaying {
		g.State = StateChallengeWaiting
		g.ChallengeActive = true
	}
}

// EndChallenge returns to normal playing state
func (g *Game) EndChallenge() {
	if g.State == StateChallengeActive || g.State == StateChallengeWaiting {
		g.State = StatePlaying
		g.ChallengeActive = false
	}
}

// AddChallengeGold adds gold from completing a challenge
func (g *Game) AddChallengeGold(gold int) {
	g.Gold += gold
}

func (g *Game) updateWaveSpawning(dt float64) {
	if g.WaveComplete {
		g.WaveCountdown -= dt
		if g.WaveCountdown <= 0 {
			g.WaveComplete = false
			g.Wave++
			g.SpawnIndex = 0
			g.WaveCountdown = 3.0
		}
		return
	}

	wave := GetWave(g.Wave)
	if g.SpawnIndex >= len(wave.Spawns) {
		// Check if wave is complete (all enemies dead)
		if len(g.Enemies) == 0 {
			g.WaveComplete = true
			// Apply economy multiplier to wave bonus
			g.Gold += g.Economy.CalculateWaveBonus(wave.BonusGold)
		}
		return
	}

	g.WaveTimer -= dt
	if g.WaveTimer <= 0 {
		spawn := wave.Spawns[g.SpawnIndex]
		g.SpawnEnemy(spawn.Type)
		g.SpawnIndex++
		if g.SpawnIndex < len(wave.Spawns) {
			g.WaveTimer = wave.Spawns[g.SpawnIndex].Delay
		}
	}
}

func (g *Game) updateEnemies(dt float64) {
	aliveEnemies := make([]*entities.Enemy, 0, len(g.Enemies))
	for _, enemy := range g.Enemies {
		if enemy.Dead {
			continue
		}
		reachedEnd := enemy.Update(dt, g.Path)
		if reachedEnd {
			// Damage player based on remaining health
			damage := int(float64(enemy.MaxHealth) * enemy.HealthPercent() * 0.1)
			if damage < 1 {
				damage = 1
			}
			g.Health -= damage
			enemy.Dead = true
			// Add escape effect at end of path
			if len(g.Path) > 0 {
				g.Effects.Add(entities.EffectHit, g.Path[len(g.Path)-1])
			}
		}
		if !enemy.Dead {
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}
	g.Enemies = aliveEnemies
}

func (g *Game) updateTowers(dt float64) {
	for _, tower := range g.Towers {
		projectile := tower.Update(dt, g.Enemies)
		if projectile != nil {
			g.Projectiles = append(g.Projectiles, projectile)
			// Add tower fire effect
			g.Effects.Add(entities.EffectTowerFire, tower.Pos)
		}
	}
}

func (g *Game) updateProjectiles(dt float64) {
	activeProjectiles := make([]*entities.Projectile, 0, len(g.Projectiles))
	for _, proj := range g.Projectiles {
		reached := proj.Update(dt)
		if reached {
			// Find enemy at target and deal damage
			for _, enemy := range g.Enemies {
				if enemy.ID == proj.TargetID && !enemy.Dead {
					killed := enemy.TakeDamage(proj.Damage)
					// Add hit effect
					g.Effects.Add(entities.EffectHit, enemy.Pos)
					if killed {
						// Apply economy multiplier to mob gold
						baseGold := enemy.Info().GoldValue
						g.Gold += g.Economy.CalculateMobGold(baseGold)
						// Add explosion effect for kill
						g.Effects.Add(entities.EffectExplosion, enemy.Pos)
					}
					break
				}
			}
		}
		if !proj.Done {
			activeProjectiles = append(activeProjectiles, proj)
		}
	}
	g.Projectiles = activeProjectiles
}

func (g *Game) checkGameEnd() {
	if g.Health <= 0 {
		g.Health = 0
		g.State = StateGameOver
	}
	if g.Wave > g.TotalWaves && len(g.Enemies) == 0 {
		g.State = StateVictory
	}
}

// MoveCursor moves the placement cursor
func (g *Game) MoveCursor(dx, dy int) {
	g.CursorX += dx
	g.CursorY += dy
	if g.CursorX < 0 {
		g.CursorX = 0
	}
	if g.CursorX >= g.Width {
		g.CursorX = g.Width - 1
	}
	if g.CursorY < 0 {
		g.CursorY = 0
	}
	if g.CursorY >= g.Height {
		g.CursorY = g.Height - 1
	}
}

// SelectTower selects a tower type for placement
func (g *Game) SelectTower(t entities.TowerType) {
	g.SelectedTower = t
}

// TogglePause toggles the pause state
func (g *Game) TogglePause() {
	if g.State == StatePlaying {
		g.State = StatePaused
	} else if g.State == StatePaused {
		g.State = StatePlaying
	}
}
