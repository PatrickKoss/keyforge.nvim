package entities

import (
	"testing"
)

func TestNewTower(t *testing.T) {
	pos := Position{X: 5, Y: 5}
	tower := NewTower(1, TowerArrow, pos)

	if tower.ID != 1 {
		t.Errorf("Expected ID 1, got %d", tower.ID)
	}
	if tower.Type != TowerArrow {
		t.Errorf("Expected TowerArrow, got %v", tower.Type)
	}
	if tower.Level != 0 {
		t.Errorf("Expected level 0, got %d", tower.Level)
	}

	info := TowerTypes[TowerArrow]
	if tower.Damage != info.Damage {
		t.Errorf("Expected damage %d, got %d", info.Damage, tower.Damage)
	}
	if tower.Range != info.Range {
		t.Errorf("Expected range %f, got %f", info.Range, tower.Range)
	}
}

func TestTowerInRange(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	// Within range
	if !tower.InRange(Position{X: 6, Y: 5}) {
		t.Error("Position (6,5) should be in range")
	}
	if !tower.InRange(Position{X: 5, Y: 6}) {
		t.Error("Position (5,6) should be in range")
	}

	// Out of range
	if tower.InRange(Position{X: 15, Y: 5}) {
		t.Error("Position (15,5) should be out of range")
	}
}

func TestTowerFindTarget(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	enemies := []*Enemy{
		NewEnemy(1, EnemyBug, Position{X: 6, Y: 5}),  // In range
		NewEnemy(2, EnemyBug, Position{X: 20, Y: 5}), // Out of range
	}
	enemies[0].PathIndex = 2
	enemies[1].PathIndex = 1

	target := tower.FindTarget(enemies)
	if target == nil {
		t.Fatal("Expected to find a target")
	}
	if target.ID != 1 {
		t.Error("Should target enemy in range")
	}
}

func TestTowerFindTargetPriority(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	enemies := []*Enemy{
		NewEnemy(1, EnemyBug, Position{X: 6, Y: 5}),
		NewEnemy(2, EnemyBug, Position{X: 4, Y: 5}),
	}
	enemies[0].PathIndex = 2 // Further along
	enemies[1].PathIndex = 5 // Furthest along

	target := tower.FindTarget(enemies)
	if target == nil {
		t.Fatal("Expected to find a target")
	}
	if target.ID != 2 {
		t.Error("Should target enemy furthest along path")
	}
}

func TestTowerFindTargetIgnoresDead(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	enemies := []*Enemy{
		NewEnemy(1, EnemyBug, Position{X: 6, Y: 5}),
	}
	enemies[0].Dead = true

	target := tower.FindTarget(enemies)
	if target != nil {
		t.Error("Should not target dead enemies")
	}
}

func TestTowerCanUpgrade(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	if !tower.CanUpgrade() {
		t.Error("New tower should be upgradeable")
	}

	// Upgrade to max
	info := TowerTypes[TowerArrow]
	for range len(info.Upgrades) {
		tower.Upgrade()
	}

	if tower.CanUpgrade() {
		t.Error("Maxed tower should not be upgradeable")
	}
}

func TestTowerUpgrade(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})
	initialDamage := tower.Damage
	initialRange := tower.Range

	success := tower.Upgrade()
	if !success {
		t.Error("First upgrade should succeed")
	}
	if tower.Level != 1 {
		t.Errorf("Expected level 1, got %d", tower.Level)
	}
	if tower.Damage <= initialDamage {
		t.Error("Damage should increase after upgrade")
	}
	if tower.Range <= initialRange {
		t.Error("Range should increase after upgrade")
	}
}

func TestTowerUpgradeCost(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})
	info := TowerTypes[TowerArrow]

	cost := tower.UpgradeCost()
	if cost != info.Upgrades[0].Cost {
		t.Errorf("Expected cost %d, got %d", info.Upgrades[0].Cost, cost)
	}

	tower.Upgrade()
	cost = tower.UpgradeCost()
	if cost != info.Upgrades[1].Cost {
		t.Errorf("Expected cost %d, got %d", info.Upgrades[1].Cost, cost)
	}
}

func TestTowerUpdate(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 5, Y: 5})

	enemies := []*Enemy{
		NewEnemy(1, EnemyBug, Position{X: 6, Y: 5}),
	}

	// First update should fire (cooldown starts at 0)
	proj := tower.Update(0.1, enemies)
	if proj == nil {
		t.Fatal("Tower should fire on first update")
	}
	if proj.Damage != tower.Damage {
		t.Errorf("Projectile damage should be %d, got %d", tower.Damage, proj.Damage)
	}

	// Immediate update should not fire (on cooldown)
	proj = tower.Update(0.1, enemies)
	if proj != nil {
		t.Error("Tower should be on cooldown")
	}
}

func TestAllTowerTypes(t *testing.T) {
	types := []TowerType{TowerArrow, TowerLSP, TowerRefactor}

	for _, ttype := range types {
		info, ok := TowerTypes[ttype]
		if !ok {
			t.Errorf("Missing info for tower type %v", ttype)
			continue
		}
		if info.Name == "" {
			t.Errorf("Tower type %v has empty name", ttype)
		}
		if info.Cost <= 0 {
			t.Errorf("Tower type %v has invalid cost: %d", ttype, info.Cost)
		}
		if info.Damage <= 0 {
			t.Errorf("Tower type %v has invalid damage: %d", ttype, info.Damage)
		}
		if info.Range <= 0 {
			t.Errorf("Tower type %v has invalid range: %f", ttype, info.Range)
		}
	}
}

func TestProjectileUpdate(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 0, Y: 0})
	enemy := NewEnemy(1, EnemyBug, Position{X: 5, Y: 0})

	proj := NewProjectile(tower, enemy)

	// Update projectile
	reached := proj.Update(0.3)
	if reached {
		t.Error("Projectile should not have reached target yet")
	}
	if proj.Pos.X <= 0 {
		t.Error("Projectile should have moved")
	}
}

func TestProjectileReachesTarget(t *testing.T) {
	tower := NewTower(1, TowerArrow, Position{X: 0, Y: 0})
	enemy := NewEnemy(1, EnemyBug, Position{X: 1, Y: 0})

	proj := NewProjectile(tower, enemy)

	// Update until reached
	for i := 0; i < 100 && !proj.Done; i++ {
		proj.Update(0.1)
	}

	if !proj.Done {
		t.Error("Projectile should have reached target")
	}
}

func TestTowerStats(t *testing.T) {
	// Verify rebalanced tower stats per design document
	tests := []struct {
		towerType  TowerType
		name       string
		cost       int
		damage     int
		towerRange float64
		cooldown   float64
	}{
		{TowerArrow, "Arrow", 50, 8, 2.5, 0.8},
		{TowerLSP, "LSP", 100, 20, 5.0, 1.5},
		{TowerRefactor, "Refactor", 150, 12, 3.0, 1.0},
	}

	for _, tc := range tests {
		info := TowerTypes[tc.towerType]
		if info.Name != tc.name {
			t.Errorf("%s: expected name %s, got %s", tc.name, tc.name, info.Name)
		}
		if info.Cost != tc.cost {
			t.Errorf("%s: expected cost %d, got %d", tc.name, tc.cost, info.Cost)
		}
		if info.Damage != tc.damage {
			t.Errorf("%s: expected damage %d, got %d", tc.name, tc.damage, info.Damage)
		}
		if info.Range != tc.towerRange {
			t.Errorf("%s: expected range %v, got %v", tc.name, tc.towerRange, info.Range)
		}
		if info.Cooldown != tc.cooldown {
			t.Errorf("%s: expected cooldown %v, got %v", tc.name, tc.cooldown, info.Cooldown)
		}
	}
}

func TestTowerUpgradeScaling(t *testing.T) {
	// Verify upgrades provide consistent scaling:
	// +15% damage (approx), +0.3 range, -10% cooldown (0.9 multiplier)
	for ttype, info := range TowerTypes {
		if len(info.Upgrades) == 0 {
			continue
		}

		for i, upgrade := range info.Upgrades {
			// Range bonus should be 0.3
			if upgrade.RangeBonus != 0.3 {
				t.Errorf("%s upgrade %d: expected range bonus 0.3, got %v", info.Name, i+1, upgrade.RangeBonus)
			}

			// Cooldown multiplier should be 0.9 (10% faster)
			if upgrade.CooldownMult != 0.9 {
				t.Errorf("%s upgrade %d: expected cooldown mult 0.9, got %v", info.Name, i+1, upgrade.CooldownMult)
			}

			// Damage bonus should be positive
			if upgrade.DamageBonus <= 0 {
				t.Errorf("%s upgrade %d: expected positive damage bonus, got %d", info.Name, i+1, upgrade.DamageBonus)
			}

			// Cost should be approximately 60% of base tower cost
			expectedCostMin := int(float64(info.Cost) * 0.5)
			expectedCostMax := int(float64(info.Cost) * 1.5)
			if upgrade.Cost < expectedCostMin || upgrade.Cost > expectedCostMax {
				t.Errorf("%s upgrade %d: cost %d outside expected range [%d, %d]",
					info.Name, i+1, upgrade.Cost, expectedCostMin, expectedCostMax)
			}
		}

		// Verify tower type is one of the expected types
		if ttype != TowerArrow && ttype != TowerLSP && ttype != TowerRefactor {
			// Skip placeholder tower types
			continue
		}
	}
}

func TestTowerRoleDistinction(t *testing.T) {
	arrow := TowerTypes[TowerArrow]
	lsp := TowerTypes[TowerLSP]
	refactor := TowerTypes[TowerRefactor]

	// Arrow should be fastest (lowest cooldown)
	if arrow.Cooldown >= lsp.Cooldown || arrow.Cooldown >= refactor.Cooldown {
		t.Error("Arrow tower should have the fastest attack speed (lowest cooldown)")
	}

	// LSP should have longest range
	if lsp.Range <= arrow.Range || lsp.Range <= refactor.Range {
		t.Error("LSP tower should have the longest range")
	}

	// Arrow should be cheapest
	if arrow.Cost >= lsp.Cost || arrow.Cost >= refactor.Cost {
		t.Error("Arrow tower should be the cheapest")
	}

	// Refactor should be most expensive
	if refactor.Cost <= arrow.Cost || refactor.Cost <= lsp.Cost {
		t.Error("Refactor tower should be the most expensive")
	}
}
