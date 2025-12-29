package entities

// Enemy represents an enemy moving along the path.
type Enemy struct {
	ID        int
	Type      EnemyType
	Pos       Position
	Health    int
	MaxHealth int
	Speed     float64
	PathIndex int     // current waypoint index
	PathProg  float64 // progress to next waypoint (0-1)
	Dead      bool
}

// NewEnemy creates a new enemy at the start of the path.
func NewEnemy(id int, enemyType EnemyType, startPos Position) *Enemy {
	info := EnemyTypes[enemyType]
	return &Enemy{
		ID:        id,
		Type:      enemyType,
		Pos:       startPos,
		Health:    info.Health,
		MaxHealth: info.Health,
		Speed:     info.Speed,
		PathIndex: 0,
		PathProg:  0,
		Dead:      false,
	}
}

// Info returns the enemy type configuration.
func (e *Enemy) Info() EnemyInfo {
	return EnemyTypes[e.Type]
}

// TakeDamage applies damage to the enemy and returns true if killed.
func (e *Enemy) TakeDamage(damage int) bool {
	e.Health -= damage
	if e.Health <= 0 {
		e.Health = 0
		e.Dead = true
		return true
	}
	return false
}

// HealthPercent returns the enemy's health as a percentage.
func (e *Enemy) HealthPercent() float64 {
	return float64(e.Health) / float64(e.MaxHealth)
}

// Update moves the enemy along the path
// Returns true if the enemy has reached the end.
func (e *Enemy) Update(dt float64, path []Position) bool {
	if e.Dead || e.PathIndex >= len(path)-1 {
		return e.PathIndex >= len(path)-1
	}

	// Calculate movement
	currentWaypoint := path[e.PathIndex]
	nextWaypoint := path[e.PathIndex+1]

	// Direction to next waypoint
	dx := nextWaypoint.X - currentWaypoint.X
	dy := nextWaypoint.Y - currentWaypoint.Y

	// Distance between waypoints (should be 1 for adjacent cells)
	dist := 1.0
	if dx != 0 && dy != 0 {
		dist = 1.414 // diagonal
	}

	// Progress along path segment
	e.PathProg += (e.Speed * dt) / dist

	// Move to next waypoint if we've reached current target
	for e.PathProg >= 1.0 && e.PathIndex < len(path)-1 {
		e.PathProg -= 1.0
		e.PathIndex++
		if e.PathIndex >= len(path)-1 {
			e.Pos = path[len(path)-1]
			return true // reached the end
		}
	}

	// Interpolate position
	if e.PathIndex < len(path)-1 {
		curr := path[e.PathIndex]
		next := path[e.PathIndex+1]
		e.Pos.X = curr.X + (next.X-curr.X)*e.PathProg
		e.Pos.Y = curr.Y + (next.Y-curr.Y)*e.PathProg
	}

	return false
}
