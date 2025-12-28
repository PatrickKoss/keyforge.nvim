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
	types := []EnemyType{EnemyBug, EnemyGremlin, EnemyDaemon, EnemyBoss}

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
	}
}
