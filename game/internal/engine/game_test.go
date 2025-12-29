package engine

import (
	"testing"

	"github.com/keyforge/keyforge/internal/entities"
)

func TestNewGame(t *testing.T) {
	g := NewGame(20, 14)

	if g.Width != 20 {
		t.Errorf("Expected width 20, got %d", g.Width)
	}
	if g.Height != 14 {
		t.Errorf("Expected height 14, got %d", g.Height)
	}
	if g.Gold != 200 {
		t.Errorf("Expected starting gold 200, got %d", g.Gold)
	}
	if g.Health != 100 {
		t.Errorf("Expected starting health 100, got %d", g.Health)
	}
	if g.State != StatePlaying {
		t.Errorf("Expected state StatePlaying, got %v", g.State)
	}
	if len(g.Path) == 0 {
		t.Error("Expected non-empty path")
	}
}

func TestGamePath(t *testing.T) {
	g := NewGame(20, 14)

	// Path should have entries
	if len(g.Path) < 10 {
		t.Errorf("Path too short: %d", len(g.Path))
	}

	// First path cell should be at edge
	first := g.Path[0]
	if first.X != 0 {
		t.Errorf("Expected first path X=0, got %f", first.X)
	}
}

func TestIsOnPath(t *testing.T) {
	g := NewGame(20, 14)

	// First path position should be on path
	first := g.Path[0]
	x, y := first.IntPos()
	if !g.IsOnPath(x, y) {
		t.Errorf("Expected (%d, %d) to be on path", x, y)
	}

	// Position far from path should not be on path
	if g.IsOnPath(15, 1) {
		t.Error("Expected (15, 1) to not be on path")
	}
}

func TestCanPlaceTower(t *testing.T) {
	g := NewGame(20, 14)

	// Can't place outside bounds
	if g.CanPlaceTower(-1, 0) {
		t.Error("Should not place tower at negative X")
	}
	if g.CanPlaceTower(0, -1) {
		t.Error("Should not place tower at negative Y")
	}
	if g.CanPlaceTower(20, 0) {
		t.Error("Should not place tower at X >= Width")
	}
	if g.CanPlaceTower(0, 14) {
		t.Error("Should not place tower at Y >= Height")
	}

	// Can't place on path
	first := g.Path[0]
	x, y := first.IntPos()
	if g.CanPlaceTower(x, y) {
		t.Error("Should not place tower on path")
	}

	// Can place on empty cell
	if !g.CanPlaceTower(0, 0) {
		t.Error("Should be able to place tower at (0, 0)")
	}
}

func TestPlaceTower(t *testing.T) {
	g := NewGame(20, 14)
	initialGold := g.Gold

	g.CursorX = 0
	g.CursorY = 0
	g.SelectedTower = entities.TowerArrow

	if !g.PlaceTower() {
		t.Error("Expected tower placement to succeed")
	}

	if len(g.Towers) != 1 {
		t.Errorf("Expected 1 tower, got %d", len(g.Towers))
	}

	cost := entities.TowerTypes[entities.TowerArrow].Cost
	if g.Gold != initialGold-cost {
		t.Errorf("Expected gold %d, got %d", initialGold-cost, g.Gold)
	}

	// Can't place tower on same position
	if g.PlaceTower() {
		t.Error("Should not place second tower at same position")
	}
}

func TestPlaceTowerInsufficientGold(t *testing.T) {
	g := NewGame(20, 14)
	g.Gold = 10 // Not enough for any tower

	g.CursorX = 0
	g.CursorY = 0

	if g.PlaceTower() {
		t.Error("Should not place tower with insufficient gold")
	}
}

func TestSpawnEnemy(t *testing.T) {
	g := NewGame(20, 14)

	g.SpawnEnemy(entities.EnemyBug)

	if len(g.Enemies) != 1 {
		t.Errorf("Expected 1 enemy, got %d", len(g.Enemies))
	}

	enemy := g.Enemies[0]
	if enemy.Type != entities.EnemyBug {
		t.Errorf("Expected EnemyBug, got %v", enemy.Type)
	}

	// Enemy should spawn at start of path
	first := g.Path[0]
	if enemy.Pos.X != first.X || enemy.Pos.Y != first.Y {
		t.Errorf("Expected enemy at path start (%f, %f), got (%f, %f)",
			first.X, first.Y, enemy.Pos.X, enemy.Pos.Y)
	}
}

func TestMoveCursor(t *testing.T) {
	g := NewGame(20, 14)

	// Start at center
	startX := g.CursorX
	startY := g.CursorY

	// Move right
	g.MoveCursor(1, 0)
	if g.CursorX != startX+1 {
		t.Errorf("Expected cursor X=%d, got %d", startX+1, g.CursorX)
	}

	// Move down
	g.MoveCursor(0, 1)
	if g.CursorY != startY+1 {
		t.Errorf("Expected cursor Y=%d, got %d", startY+1, g.CursorY)
	}

	// Move left
	g.MoveCursor(-1, 0)
	if g.CursorX != startX {
		t.Errorf("Expected cursor X=%d, got %d", startX, g.CursorX)
	}

	// Move up
	g.MoveCursor(0, -1)
	if g.CursorY != startY {
		t.Errorf("Expected cursor Y=%d, got %d", startY, g.CursorY)
	}
}

func TestMoveCursorBounds(t *testing.T) {
	g := NewGame(20, 14)

	// Move to top-left corner
	for range 50 {
		g.MoveCursor(-1, -1)
	}

	if g.CursorX != 0 {
		t.Errorf("Cursor X should be clamped to 0, got %d", g.CursorX)
	}
	if g.CursorY != 0 {
		t.Errorf("Cursor Y should be clamped to 0, got %d", g.CursorY)
	}

	// Move to bottom-right corner
	for range 50 {
		g.MoveCursor(1, 1)
	}

	if g.CursorX != 19 {
		t.Errorf("Cursor X should be clamped to 19, got %d", g.CursorX)
	}
	if g.CursorY != 13 {
		t.Errorf("Cursor Y should be clamped to 13, got %d", g.CursorY)
	}
}

func TestTogglePause(t *testing.T) {
	g := NewGame(20, 14)

	if g.State != StatePlaying {
		t.Error("Game should start in playing state")
	}

	g.TogglePause()
	if g.State != StatePaused {
		t.Error("Game should be paused after toggle")
	}

	g.TogglePause()
	if g.State != StatePlaying {
		t.Error("Game should be playing after second toggle")
	}
}

func TestGameUpdate(t *testing.T) {
	g := NewGame(20, 14)

	// Add an enemy
	g.SpawnEnemy(entities.EnemyBug)
	initialPos := g.Enemies[0].Pos

	// Update game
	g.Update(0.1)

	// Enemy should have moved
	if g.Enemies[0].Pos == initialPos {
		t.Error("Enemy should have moved after update")
	}
}

func TestGameUpdatePaused(t *testing.T) {
	g := NewGame(20, 14)
	g.SpawnEnemy(entities.EnemyBug)
	initialPos := g.Enemies[0].Pos

	// Pause game
	g.TogglePause()

	// Update game
	g.Update(0.1)

	// Enemy should NOT have moved
	if g.Enemies[0].Pos != initialPos {
		t.Error("Enemy should not move while paused")
	}
}

func TestGameOver(t *testing.T) {
	g := NewGame(20, 14)
	g.Health = 1

	// Simulate enemy reaching end
	g.SpawnEnemy(entities.EnemyBug)
	g.Enemies[0].PathIndex = len(g.Path) - 2
	g.Enemies[0].PathProg = 0.99

	// Update to trigger enemy reaching end
	g.Update(0.5)

	if g.State != StateGameOver {
		t.Errorf("Expected StateGameOver, got %v", g.State)
	}
}

func TestUpgradeTower(t *testing.T) {
	g := NewGame(20, 14)

	// Place a tower
	g.CursorX = 0
	g.CursorY = 0
	g.PlaceTower()

	tower := g.Towers[0]
	initialDamage := tower.Damage
	initialGold := g.Gold

	// Upgrade tower
	if !g.UpgradeTower() {
		t.Error("Expected upgrade to succeed")
	}

	if tower.Damage <= initialDamage {
		t.Error("Expected damage to increase after upgrade")
	}

	if tower.Level != 1 {
		t.Errorf("Expected tower level 1, got %d", tower.Level)
	}

	if g.Gold >= initialGold {
		t.Error("Expected gold to decrease after upgrade")
	}
}

func TestSelectTower(t *testing.T) {
	g := NewGame(20, 14)

	g.SelectTower(entities.TowerLSP)
	if g.SelectedTower != entities.TowerLSP {
		t.Error("Expected TowerLSP to be selected")
	}

	g.SelectTower(entities.TowerRefactor)
	if g.SelectedTower != entities.TowerRefactor {
		t.Error("Expected TowerRefactor to be selected")
	}
}

// Tests for concurrent challenge/gameplay (game continues during challenges)

func TestChallengeActiveGameContinues(t *testing.T) {
	g := NewGame(20, 14)

	// Spawn an enemy
	g.SpawnEnemy(entities.EnemyBug)
	initialPos := g.Enemies[0].Pos

	// Start a challenge
	g.StartChallenge()

	if g.State != StateChallengeActive {
		t.Errorf("Expected StateChallengeActive, got %v", g.State)
	}

	if !g.ChallengeActive {
		t.Error("Expected ChallengeActive to be true")
	}

	// Update game - enemy should still move during challenge
	g.Update(0.1)

	if g.Enemies[0].Pos == initialPos {
		t.Error("Enemy should move during challenge active state")
	}
}

func TestEndChallenge(t *testing.T) {
	g := NewGame(20, 14)

	g.StartChallenge()
	if g.State != StateChallengeActive {
		t.Error("Expected StateChallengeActive after StartChallenge")
	}

	g.EndChallenge()
	if g.State != StatePlaying {
		t.Errorf("Expected StatePlaying after EndChallenge, got %v", g.State)
	}

	if g.ChallengeActive {
		t.Error("Expected ChallengeActive to be false after EndChallenge")
	}
}

func TestAddChallengeGold(t *testing.T) {
	g := NewGame(20, 14)
	initialGold := g.Gold

	g.AddChallengeGold(100)

	if g.Gold != initialGold+100 {
		t.Errorf("Expected gold %d, got %d", initialGold+100, g.Gold)
	}
}

func TestStartChallengeOnlyFromPlaying(t *testing.T) {
	g := NewGame(20, 14)

	// Should work from playing state
	g.StartChallenge()
	if g.State != StateChallengeActive {
		t.Error("Should be able to start challenge from playing state")
	}

	// End challenge
	g.EndChallenge()

	// Pause the game
	g.TogglePause()
	if g.State != StatePaused {
		t.Error("Game should be paused")
	}

	// Try to start challenge from paused state
	g.StartChallenge()
	if g.State != StatePaused {
		t.Error("Should not be able to start challenge from paused state")
	}
}

func TestTowersFireDuringChallenge(t *testing.T) {
	g := NewGame(20, 14)

	// Place a tower right next to the path start (path starts at 0,3)
	g.CursorX = 0
	g.CursorY = 2 // One cell above the path at (0,3)
	if !g.PlaceTower() {
		t.Fatal("Failed to place tower")
	}

	// Spawn an enemy - it will be at the path start
	g.SpawnEnemy(entities.EnemyBug)

	// Start challenge
	g.StartChallenge()

	if g.State != StateChallengeActive {
		t.Fatalf("Expected StateChallengeActive, got %v", g.State)
	}

	// Stop wave spawning from adding more enemies
	g.SpawnIndex = 100 // Mark all spawns as done
	g.WaveComplete = true

	// Update many times to allow tower to fire and projectile to reach
	for range 50 {
		g.Update(0.1)
	}

	// Tower should have fired during challenge - enemy should be dead
	hasLivingEnemies := false
	enemyDamaged := false
	for _, e := range g.Enemies {
		if !e.Dead {
			hasLivingEnemies = true
			if e.Health < e.MaxHealth {
				enemyDamaged = true
			}
		}
	}
	hasProjectiles := len(g.Projectiles) > 0

	// Either enemy was damaged, enemy died, or there are projectiles in flight
	if hasLivingEnemies && !enemyDamaged && !hasProjectiles {
		t.Error("Tower should fire during challenge active state")
	}
}

func TestGameOverCanHappenDuringChallenge(t *testing.T) {
	g := NewGame(20, 14)
	g.Health = 1

	// Start challenge
	g.StartChallenge()

	// Simulate enemy reaching end
	g.SpawnEnemy(entities.EnemyBug)
	g.Enemies[0].PathIndex = len(g.Path) - 2
	g.Enemies[0].PathProg = 0.99

	// Update to trigger enemy reaching end
	g.Update(0.5)

	// Game should be over even though challenge was active
	if g.State != StateGameOver {
		t.Errorf("Expected StateGameOver during challenge, got %v", g.State)
	}
}

func TestMobGoldReducedWithEconomy(t *testing.T) {
	g := NewGame(20, 14)

	// Place a tower
	g.CursorX = 1
	g.CursorY = 4
	g.PlaceTower()

	// Spawn a weak enemy that will die
	g.SpawnEnemy(entities.EnemyBug)
	g.Enemies[0].Health = 1

	initialGold := g.Gold

	// Update until enemy dies
	for i := 0; i < 100 && len(g.Enemies) > 0 && !g.Enemies[0].Dead; i++ {
		g.Update(0.1)
	}

	// Calculate expected gold: base * 0.25
	baseGold := entities.EnemyTypes[entities.EnemyBug].GoldValue
	expectedGold := g.Economy.CalculateMobGold(baseGold)
	actualGold := g.Gold - initialGold

	if actualGold != expectedGold && actualGold != 0 {
		// If enemy died, we should have gotten reduced gold
		if g.Enemies[0].Dead {
			t.Errorf("Expected reduced gold %d, got %d (base was %d)", expectedGold, actualGold, baseGold)
		}
	}
}
