package entities

// Effect represents a visual effect on the grid.
type Effect struct {
	ID       int
	Type     EffectType
	Pos      Position
	Duration float64 // total duration in seconds
	Elapsed  float64 // time elapsed
	Frame    int     // current animation frame
	Done     bool
}

// EffectType identifies different visual effects.
type EffectType int

const (
	EffectExplosion EffectType = iota
	EffectHit
	EffectSpawn
	EffectTowerFire
	EffectLevelUp
	EffectGoldGain
)

// EffectInfo contains configuration for each effect type.
type EffectInfo struct {
	Frames   []string // animation frames
	Duration float64  // total duration
	Color    string   // hex color
}

// EffectTypes contains all effect configurations.
var EffectTypes = map[EffectType]EffectInfo{
	EffectExplosion: {
		Frames:   []string{"💥", "✨", "·", " "},
		Duration: 0.4,
		Color:    "#f97316",
	},
	EffectHit: {
		Frames:   []string{"✦", "✧", "·"},
		Duration: 0.2,
		Color:    "#fbbf24",
	},
	EffectSpawn: {
		Frames:   []string{"○", "◎", "●"},
		Duration: 0.3,
		Color:    "#ef4444",
	},
	EffectTowerFire: {
		Frames:   []string{"⚡", "✶", "✦"},
		Duration: 0.15,
		Color:    "#22c55e",
	},
	EffectLevelUp: {
		Frames:   []string{"⬆", "⇧", "↑", "✨"},
		Duration: 0.5,
		Color:    "#8b5cf6",
	},
	EffectGoldGain: {
		Frames:   []string{"+$", "💰", "✨"},
		Duration: 0.6,
		Color:    "#fbbf24",
	},
}

var effectIDCounter = 0

// NewEffect creates a new visual effect.
func NewEffect(effectType EffectType, pos Position) *Effect {
	effectIDCounter++
	info := EffectTypes[effectType]
	return &Effect{
		ID:       effectIDCounter,
		Type:     effectType,
		Pos:      pos,
		Duration: info.Duration,
		Elapsed:  0,
		Frame:    0,
		Done:     false,
	}
}

// Info returns the effect type configuration.
func (e *Effect) Info() EffectInfo {
	return EffectTypes[e.Type]
}

// Update advances the effect animation
// Returns true if the effect has finished.
func (e *Effect) Update(dt float64) bool {
	if e.Done {
		return true
	}

	e.Elapsed += dt
	info := e.Info()

	// Calculate current frame
	frameCount := len(info.Frames)
	if frameCount > 0 {
		progress := e.Elapsed / e.Duration
		e.Frame = int(progress * float64(frameCount))
		if e.Frame >= frameCount {
			e.Frame = frameCount - 1
		}
	}

	// Check if effect is complete
	if e.Elapsed >= e.Duration {
		e.Done = true
		return true
	}

	return false
}

// CurrentFrame returns the current animation frame character.
func (e *Effect) CurrentFrame() string {
	info := e.Info()
	if e.Frame >= 0 && e.Frame < len(info.Frames) {
		return info.Frames[e.Frame]
	}
	return ""
}

// EffectManager manages all active effects.
type EffectManager struct {
	Effects []*Effect
}

// NewEffectManager creates a new effect manager.
func NewEffectManager() *EffectManager {
	return &EffectManager{
		Effects: make([]*Effect, 0),
	}
}

// Add creates and adds a new effect.
func (em *EffectManager) Add(effectType EffectType, pos Position) *Effect {
	effect := NewEffect(effectType, pos)
	em.Effects = append(em.Effects, effect)
	return effect
}

// Update updates all effects and removes completed ones.
func (em *EffectManager) Update(dt float64) {
	active := make([]*Effect, 0, len(em.Effects))
	for _, effect := range em.Effects {
		effect.Update(dt)
		if !effect.Done {
			active = append(active, effect)
		}
	}
	em.Effects = active
}

// GetEffectAt returns the effect at a position (if any).
func (em *EffectManager) GetEffectAt(x, y int) *Effect {
	// Return the most recent effect at position
	for i := len(em.Effects) - 1; i >= 0; i-- {
		effect := em.Effects[i]
		ex, ey := effect.Pos.IntPos()
		if ex == x && ey == y && !effect.Done {
			return effect
		}
	}
	return nil
}

// Clear removes all effects.
func (em *EffectManager) Clear() {
	em.Effects = make([]*Effect, 0)
}
