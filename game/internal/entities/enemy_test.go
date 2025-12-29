package entities

import (
	"testing"
)

func TestNewEnemy(t *testing.T) {
	pos := Position{X: 0, Y: 0}
	enemy := NewEnemy(1, EnemyBug, pos)

	if enemy.ID != 1 {
		t.Errorf("Expected ID 1, got %d", enemy.ID)
	}
	if enemy.Type != EnemyBug {
		t.Errorf("Expected EnemyBug, got %v", enemy.Type)
	}
	if enemy.Dead {
		t.Error("New enemy should not be dead")
	}

	info := EnemyTypes[EnemyBug]
	if enemy.Health != info.Health {
		t.Errorf("Expected health %d, got %d", info.Health, enemy.Health)
	}
	if enemy.Speed != info.Speed {
		t.Errorf("Expected speed %f, got %f", info.Speed, enemy.Speed)
	}
}

func TestEnemyTakeDamage(t *testing.T) {
	enemy := NewEnemy(1, EnemyBug, Position{X: 0, Y: 0})
	initialHealth := enemy.Health

	killed := enemy.TakeDamage(5)
	if killed {
		t.Error("Enemy should not be killed by 5 damage")
	}
	if enemy.Health != initialHealth-5 {
		t.Errorf("Expected health %d, got %d", initialHealth-5, enemy.Health)
	}

	// Kill the enemy
	killed = enemy.TakeDamage(100)
	if !killed {
		t.Error("Enemy should be killed by 100 damage")
	}
	if !enemy.Dead {
		t.Error("Enemy should be marked as dead")
	}
	if enemy.Health != 0 {
		t.Errorf("Dead enemy health should be 0, got %d", enemy.Health)
	}
}

func TestEnemyHealthPercent(t *testing.T) {
	enemy := NewEnemy(1, EnemyBug, Position{X: 0, Y: 0})

	// Full health
	percent := enemy.HealthPercent()
	if percent != 1.0 {
		t.Errorf("Expected 100%% health, got %f", percent)
	}

	// Half health
	enemy.Health = enemy.MaxHealth / 2
	percent = enemy.HealthPercent()
	if percent != 0.5 {
		t.Errorf("Expected 50%% health, got %f", percent)
	}
}

func TestEnemyUpdate(t *testing.T) {
	path := []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
		{X: 2, Y: 0},
		{X: 3, Y: 0},
	}

	enemy := NewEnemy(1, EnemyBug, path[0])

	// Update and check movement
	reachedEnd := enemy.Update(0.5, path)
	if reachedEnd {
		t.Error("Enemy should not have reached end yet")
	}
	if enemy.Pos.X <= 0 {
		t.Error("Enemy should have moved forward")
	}
}

func TestEnemyReachesEnd(t *testing.T) {
	path := []Position{
		{X: 0, Y: 0},
		{X: 1, Y: 0},
	}

	enemy := NewEnemy(1, EnemyBug, path[0])
	enemy.PathIndex = 0
	enemy.PathProg = 0.99

	// Update to trigger reaching end
	reachedEnd := enemy.Update(0.5, path)
	if !reachedEnd {
		t.Error("Enemy should have reached the end")
	}
}

func TestEnemyInfo(t *testing.T) {
	enemy := NewEnemy(1, EnemyGremlin, Position{})

	info := enemy.Info()
	if info.Name != "Gremlin" {
		t.Errorf("Expected 'Gremlin', got %s", info.Name)
	}
}

func TestAllEnemyTypes(t *testing.T) {
	// All 7 enemy types
	types := []EnemyType{EnemyMite, EnemyBug, EnemyGremlin, EnemyCrawler, EnemySpecter, EnemyDaemon, EnemyBoss}

	for _, etype := range types {
		info, ok := EnemyTypes[etype]
		if !ok {
			t.Errorf("Missing info for enemy type %v", etype)
			continue
		}
		if info.Name == "" {
			t.Errorf("Enemy type %v has empty name", etype)
		}
		if info.Health <= 0 {
			t.Errorf("Enemy type %v has invalid health: %d", etype, info.Health)
		}
		if info.Speed <= 0 {
			t.Errorf("Enemy type %v has invalid speed: %f", etype, info.Speed)
		}
		if info.GoldValue <= 0 {
			t.Errorf("Enemy type %v has invalid gold value: %d", etype, info.GoldValue)
		}
	}
}

func TestEnemyTypeCount(t *testing.T) {
	// Verify exactly 7 enemy types exist
	expectedCount := 7
	actualCount := len(EnemyTypes)
	if actualCount != expectedCount {
		t.Errorf("Expected %d enemy types, got %d", expectedCount, actualCount)
	}
}

func TestEnemyTypeStats(t *testing.T) {
	// Verify specific stats per design document
	tests := []struct {
		enemyType EnemyType
		name      string
		health    int
		speed     float64
		gold      int
	}{
		{EnemyMite, "Mite", 5, 2.0, 2},
		{EnemyBug, "Bug", 10, 1.5, 5},
		{EnemyGremlin, "Gremlin", 25, 2.5, 10},
		{EnemyCrawler, "Crawler", 40, 0.6, 15},
		{EnemySpecter, "Specter", 15, 3.5, 8},
		{EnemyDaemon, "Daemon", 100, 0.8, 25},
		{EnemyBoss, "Boss", 500, 0.5, 100},
	}

	for _, tc := range tests {
		info := EnemyTypes[tc.enemyType]
		if info.Name != tc.name {
			t.Errorf("%s: expected name %s, got %s", tc.name, tc.name, info.Name)
		}
		if info.Health != tc.health {
			t.Errorf("%s: expected health %d, got %d", tc.name, tc.health, info.Health)
		}
		if info.Speed != tc.speed {
			t.Errorf("%s: expected speed %v, got %v", tc.name, tc.speed, info.Speed)
		}
		if info.GoldValue != tc.gold {
			t.Errorf("%s: expected gold %d, got %d", tc.name, tc.gold, info.GoldValue)
		}
	}
}

func TestEnemyGoldScalesWithHealth(t *testing.T) {
	// Verify higher health enemies give more gold (generally)
	// Gold should scale: low health (2-5), medium (8-15), high (25-100)
	lowHealthEnemies := []EnemyType{EnemyMite, EnemyBug}
	medHealthEnemies := []EnemyType{EnemyGremlin, EnemyCrawler, EnemySpecter}
	highHealthEnemies := []EnemyType{EnemyDaemon, EnemyBoss}

	var lowGoldMax, medGoldMin, medGoldMax, highGoldMin int

	for _, et := range lowHealthEnemies {
		if EnemyTypes[et].GoldValue > lowGoldMax {
			lowGoldMax = EnemyTypes[et].GoldValue
		}
	}
	for _, et := range medHealthEnemies {
		g := EnemyTypes[et].GoldValue
		if medGoldMin == 0 || g < medGoldMin {
			medGoldMin = g
		}
		if g > medGoldMax {
			medGoldMax = g
		}
	}
	for _, et := range highHealthEnemies {
		g := EnemyTypes[et].GoldValue
		if highGoldMin == 0 || g < highGoldMin {
			highGoldMin = g
		}
	}

	if lowGoldMax >= medGoldMin {
		t.Errorf("Low health enemies should give less gold than medium health: max low=%d, min med=%d", lowGoldMax, medGoldMin)
	}
	if medGoldMax >= highGoldMin {
		t.Errorf("Medium health enemies should give less gold than high health: max med=%d, min high=%d", medGoldMax, highGoldMin)
	}
}
