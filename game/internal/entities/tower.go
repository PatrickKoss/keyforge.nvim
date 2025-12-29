package entities

import "math"

// Tower represents a defensive tower placed on the grid.
type Tower struct {
	ID           int
	Type         TowerType
	Pos          Position
	Level        int
	Damage       int
	Range        float64
	Cooldown     float64
	CooldownLeft float64
	Target       *Enemy
}

// NewTower creates a new tower at the specified position.
func NewTower(id int, towerType TowerType, pos Position) *Tower {
	info := TowerTypes[towerType]
	return &Tower{
		ID:           id,
		Type:         towerType,
		Pos:          pos,
		Level:        0,
		Damage:       info.Damage,
		Range:        info.Range,
		Cooldown:     info.Cooldown,
		CooldownLeft: 0,
		Target:       nil,
	}
}

// Info returns the tower type configuration.
func (t *Tower) Info() TowerInfo {
	return TowerTypes[t.Type]
}

// CanUpgrade returns true if the tower can be upgraded.
func (t *Tower) CanUpgrade() bool {
	info := t.Info()
	return t.Level < len(info.Upgrades)
}

// UpgradeCost returns the cost of the next upgrade, or 0 if maxed.
func (t *Tower) UpgradeCost() int {
	info := t.Info()
	if t.Level >= len(info.Upgrades) {
		return 0
	}
	return info.Upgrades[t.Level].Cost
}

// Upgrade applies the next upgrade level.
func (t *Tower) Upgrade() bool {
	info := t.Info()
	if t.Level >= len(info.Upgrades) {
		return false
	}
	upgrade := info.Upgrades[t.Level]
	t.Damage += upgrade.DamageBonus
	t.Range += upgrade.RangeBonus
	t.Cooldown *= upgrade.CooldownMult
	t.Level++
	return true
}

// InRange checks if a position is within the tower's range.
func (t *Tower) InRange(pos Position) bool {
	dx := t.Pos.X - pos.X
	dy := t.Pos.Y - pos.Y
	distSq := dx*dx + dy*dy
	return distSq <= t.Range*t.Range
}

// FindTarget finds the best enemy target within range
// Uses "first" strategy (furthest along path).
func (t *Tower) FindTarget(enemies []*Enemy) *Enemy {
	var bestTarget *Enemy
	bestProgress := -1.0

	for _, enemy := range enemies {
		if enemy.Dead || !t.InRange(enemy.Pos) {
			continue
		}
		// Calculate total progress (waypoint index + fractional progress)
		progress := float64(enemy.PathIndex) + enemy.PathProg
		if progress > bestProgress {
			bestProgress = progress
			bestTarget = enemy
		}
	}

	return bestTarget
}

// Update handles tower cooldown and targeting
// Returns a projectile if the tower fires, nil otherwise.
func (t *Tower) Update(dt float64, enemies []*Enemy) *Projectile {
	// Update cooldown
	if t.CooldownLeft > 0 {
		t.CooldownLeft -= dt
	}

	// Find target
	t.Target = t.FindTarget(enemies)
	if t.Target == nil {
		return nil
	}

	// Fire if ready
	if t.CooldownLeft <= 0 {
		t.CooldownLeft = t.Cooldown
		return NewProjectile(t, t.Target)
	}

	return nil
}

// Projectile represents a projectile fired by a tower.
type Projectile struct {
	ID       int
	Pos      Position
	Target   Position
	Damage   int
	Speed    float64
	TargetID int
	Done     bool
}

var projectileIDCounter = 0

// NewProjectile creates a new projectile aimed at an enemy.
func NewProjectile(tower *Tower, target *Enemy) *Projectile {
	projectileIDCounter++
	return &Projectile{
		ID:       projectileIDCounter,
		Pos:      tower.Pos,
		Target:   target.Pos,
		Damage:   tower.Damage,
		Speed:    10.0, // cells per second
		TargetID: target.ID,
		Done:     false,
	}
}

// Update moves the projectile toward its target
// Returns true if the projectile has reached its destination.
func (p *Projectile) Update(dt float64) bool {
	if p.Done {
		return true
	}

	dx := p.Target.X - p.Pos.X
	dy := p.Target.Y - p.Pos.Y
	dist := math.Sqrt(dx*dx + dy*dy)

	if dist < 0.1 {
		p.Pos = p.Target
		p.Done = true
		return true
	}

	// Normalize and move
	moveAmount := p.Speed * dt
	if moveAmount >= dist {
		p.Pos = p.Target
		p.Done = true
		return true
	}

	p.Pos.X += (dx / dist) * moveAmount
	p.Pos.Y += (dy / dist) * moveAmount
	return false
}
