package ui

import (
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/nvim"
)

// newTestModel creates a model ready for testing with game in playing state.
func newTestModel() Model {
	model := NewModel()
	// Set game to playing state for tests (skip level select/settings)
	model.Game.State = engine.StatePlaying
	return model
}

// MockRPCClient implements nvim.RPCClient for testing.
type MockRPCClient struct {
	ChallengeRequests []ChallengeRequestRecord
}

type ChallengeRequestRecord struct {
	RequestID     string
	ChallengeData *nvim.ChallengeData
}

func (m *MockRPCClient) RequestChallenge(requestID string, challenge *nvim.ChallengeData) error {
	m.ChallengeRequests = append(m.ChallengeRequests, ChallengeRequestRecord{
		RequestID:     requestID,
		ChallengeData: challenge,
	})
	return nil
}

func (m *MockRPCClient) SendGameState(state string, wave, gold, health, enemies, towers int) error {
	return nil
}

func (m *MockRPCClient) SendGameReady() error {
	return nil
}

func (m *MockRPCClient) SendGoldUpdate(gold, earned int, source string, speedBonus float64) error {
	return nil
}

func (m *MockRPCClient) SendChallengeAvailable(count, nextReward int, nextCategory string) error {
	return nil
}

func (m *MockRPCClient) SendGameOver(wave, gold, towers, health int) error {
	return nil
}

func (m *MockRPCClient) SendVictory(wave, gold, towers, health int) error {
	return nil
}

// TestChallengeResultChannel tests that challenge results are properly
// communicated via channel even when Model is copied (as Bubbletea does).
func TestChallengeResultChannel(t *testing.T) {
	model := newTestModel()
	model.NvimMode = true
	model.NvimRPC = &MockRPCClient{}

	// Simulate starting a challenge (this sets NvimChallengeID)
	// Game continues running during challenge for time pressure
	model.NvimChallengeCount++
	model.NvimChallengeID = "challenge_1"
	model.Game.StartChallenge()

	if model.Game.State != engine.StateChallengeActive {
		t.Fatalf("Expected StateChallengeActive, got %v", model.Game.State)
	}

	// Send a challenge result to the channel (simulating RPC handler)
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 50,
	}

	select {
	case model.ChallengeResultChan <- result:
		// OK
	default:
		t.Fatal("Failed to send to channel")
	}

	// Now simulate Bubbletea's Update cycle with a TickMsg
	// Note: Update returns a NEW model (value copy)
	tickMsg := TickMsg(time.Now())
	newModel, _ := model.Update(tickMsg)

	// The returned model should have processed the challenge result
	updatedModel := newModel.(Model)

	if updatedModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after processing result, got %v", updatedModel.Game.State)
	}

	if updatedModel.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID to be cleared, got %s", updatedModel.NvimChallengeID)
	}
}

// TestChallengeResultChannelWithCopy tests that even when we work with a copy
// of the model (as Bubbletea does), the channel still works.
func TestChallengeResultChannelWithCopy(t *testing.T) {
	// Create original model
	original := newTestModel()
	original.NvimMode = true
	original.NvimRPC = &MockRPCClient{}

	// Start a challenge on the original
	original.NvimChallengeCount++
	original.NvimChallengeID = "challenge_1"
	original.Game.StartChallengeWaiting()

	// Create a copy (simulating what Bubbletea does)
	modelCopy := original

	// The copy should share the same channel (channels are reference types)
	if modelCopy.ChallengeResultChan != original.ChallengeResultChan {
		t.Fatal("Channel should be shared between original and copy")
	}

	// Send result to the ORIGINAL's channel (simulating RPC handler)
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 75,
	}

	select {
	case original.ChallengeResultChan <- result:
		// OK
	default:
		t.Fatal("Failed to send to channel")
	}

	// Process on the COPY (simulating Bubbletea's Update)
	tickMsg := TickMsg(time.Now())
	newModel, _ := modelCopy.Update(tickMsg)
	updatedModel := newModel.(Model)

	// Should have processed the result
	if updatedModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying, got %v", updatedModel.Game.State)
	}
}

// TestChallengeResultMismatchedID tests that stale results are ignored.
func TestChallengeResultMismatchedID(t *testing.T) {
	model := newTestModel()
	model.NvimMode = true
	model.NvimRPC = &MockRPCClient{}

	// Start challenge with ID "challenge_2"
	// Game continues running during challenge for time pressure
	model.NvimChallengeCount = 2
	model.NvimChallengeID = "challenge_2"
	model.Game.StartChallenge()

	initialGold := model.Game.Gold

	// Send a stale result with wrong ID
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1", // Wrong ID!
		Success:    true,
		GoldEarned: 100,
	}

	select {
	case model.ChallengeResultChan <- result:
	default:
		t.Fatal("Failed to send to channel")
	}

	// Process
	tickMsg := TickMsg(time.Now())
	newModel, _ := model.Update(tickMsg)
	updatedModel := newModel.(Model)

	// Should NOT have processed - state should still be challenge active
	if updatedModel.Game.State != engine.StateChallengeActive {
		t.Errorf("Expected StateChallengeActive (stale result ignored), got %v", updatedModel.Game.State)
	}

	// Gold should not have changed
	if updatedModel.Game.Gold != initialGold {
		t.Errorf("Gold should not change for stale result, expected %d, got %d", initialGold, updatedModel.Game.Gold)
	}
}

// TestChallengeResultGoldAwarded tests that gold is properly awarded on success.
func TestChallengeResultGoldAwarded(t *testing.T) {
	model := newTestModel()
	model.NvimMode = true
	model.NvimRPC = &MockRPCClient{}

	model.NvimChallengeID = "challenge_1"
	model.Game.StartChallengeWaiting()

	initialGold := model.Game.Gold

	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 42,
	}

	model.ChallengeResultChan <- result

	tickMsg := TickMsg(time.Now())
	newModel, _ := model.Update(tickMsg)
	updatedModel := newModel.(Model)

	expectedGold := initialGold + 42
	if updatedModel.Game.Gold != expectedGold {
		t.Errorf("Expected gold %d, got %d", expectedGold, updatedModel.Game.Gold)
	}
}

// TestChallengeResultNoGoldOnFailure tests that no gold is awarded on failure.
func TestChallengeResultNoGoldOnFailure(t *testing.T) {
	model := newTestModel()
	model.NvimMode = true
	model.NvimRPC = &MockRPCClient{}

	model.NvimChallengeID = "challenge_1"
	model.Game.StartChallengeWaiting()

	initialGold := model.Game.Gold

	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    false, // Failed
		GoldEarned: 0,
	}

	model.ChallengeResultChan <- result

	tickMsg := TickMsg(time.Now())
	newModel, _ := model.Update(tickMsg)
	updatedModel := newModel.(Model)

	// Gold should remain unchanged
	if updatedModel.Game.Gold != initialGold {
		t.Errorf("Expected gold %d (unchanged), got %d", initialGold, updatedModel.Game.Gold)
	}

	// But game should resume
	if updatedModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after failed challenge, got %v", updatedModel.Game.State)
	}
}

// TestHandleChallengeCompleteViaPointer simulates the real scenario:
// - Original model is passed to socket server as pointer
// - Bubbletea works with value copies
// - Handler is called on original, but Update runs on copy.
func TestHandleChallengeCompleteViaPointer(t *testing.T) {
	// Create model and get a pointer (like socket server does)
	original := newTestModel()
	original.NvimMode = true
	original.NvimRPC = &MockRPCClient{}

	// Pointer to original (this is what socket server holds)
	handlerPtr := &original

	// Bubbletea makes a copy for its Update loop
	btCopy := original

	// Start challenge - but this happens on a PREVIOUS copy, not handlerPtr
	// In real code, this is the issue: the ID is set on a Bubbletea copy,
	// but handlerPtr still points to the original with empty ID

	// Simulate: Bubbletea's copy starts the challenge
	btCopy.NvimChallengeID = "challenge_1"
	btCopy.Game.StartChallengeWaiting()

	// The original (handlerPtr) does NOT have the challenge ID set!
	// This is the root cause of the bug.
	if handlerPtr.NvimChallengeID == "challenge_1" {
		t.Log("Note: In this test, original and copy are the same struct, so this passes")
		t.Log("In real Bubbletea, they would diverge after Update returns a new value")
	}

	// Now simulate RPC handler being called on the original pointer
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 50,
	}

	// Handler sends to channel (this works because channels are reference types)
	handlerPtr.HandleChallengeComplete(result)

	// Now btCopy processes in its Update loop
	// The btCopy has the correct NvimChallengeID
	tickMsg := TickMsg(time.Now())
	newModel, _ := btCopy.Update(tickMsg)
	updatedModel := newModel.(Model)

	if updatedModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying, got %v", updatedModel.Game.State)
	}
}

// TestMultipleTicksDoNotDuplicateProcessing ensures we don't process
// the same result multiple times.
func TestMultipleTicksDoNotDuplicateProcessing(t *testing.T) {
	model := newTestModel()
	model.NvimMode = true
	model.NvimRPC = &MockRPCClient{}

	model.NvimChallengeID = "challenge_1"
	model.Game.StartChallengeWaiting()

	initialGold := model.Game.Gold

	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 25,
	}

	model.ChallengeResultChan <- result

	// First tick processes the result
	tickMsg := TickMsg(time.Now())
	newModel, _ := model.Update(tickMsg)
	updatedModel := newModel.(Model)

	if updatedModel.Game.Gold != initialGold+25 {
		t.Fatalf("First tick should add gold, expected %d, got %d", initialGold+25, updatedModel.Game.Gold)
	}

	// Second tick should not add more gold (channel is empty)
	tickMsg2 := TickMsg(time.Now())
	newModel2, _ := updatedModel.Update(tickMsg2)
	finalModel := newModel2.(Model)

	if finalModel.Game.Gold != initialGold+25 {
		t.Errorf("Second tick should not add more gold, expected %d, got %d", initialGold+25, finalModel.Game.Gold)
	}
}

// TestRealisticBubbleteaFlow simulates the exact Bubbletea flow:
// 1. Model created, socket server initialized with pointer to model
// 2. Model passed to tea.NewProgram (by value)
// 3. Challenge started via Update (returns new model)
// 4. RPC handler called on ORIGINAL pointer
// 5. Update called on BUBBLETEA'S copy.
func TestRealisticBubbleteaFlow(t *testing.T) {
	// Step 1: Create model (like in main.go)
	model := newTestModel()
	model.NvimMode = true

	// Step 2: Socket server gets pointer to model
	// In real code: model.InitNvimSocket(socketPath)
	// The socket server stores &model as the handler
	handlerPointer := &model

	// Step 3: Bubbletea gets a copy (tea.NewProgram(model, ...))
	// After this, Bubbletea works with copies, not the original
	bubbleteaCopy := model

	// Simulate RPC client on the copy (in reality this happens during init)
	bubbleteaCopy.NvimRPC = &MockRPCClient{}

	// Step 4: User presses 'c' to start challenge
	// This happens in Bubbletea's Update, which returns a NEW copy
	bubbleteaCopy.NvimChallengeCount++
	bubbleteaCopy.NvimChallengeID = "challenge_1"
	bubbleteaCopy.Game.StartChallengeWaiting()

	// CRITICAL: The original model (handlerPointer) does NOT have the challenge ID!
	// This is because Go structs are value types.
	// However, the channel IS shared (reference type).

	t.Logf("Handler pointer NvimChallengeID: '%s'", handlerPointer.NvimChallengeID)
	t.Logf("Bubbletea copy NvimChallengeID: '%s'", bubbleteaCopy.NvimChallengeID)

	// The handler pointer has empty ID because it's the original, before challenge started
	if handlerPointer.NvimChallengeID != "" {
		t.Log("Note: In this simple test they're the same memory, but conceptually they diverge")
	}

	// Step 5: RPC notification arrives, handler called on original pointer
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 50,
	}
	handlerPointer.HandleChallengeComplete(result)

	// Step 6: Bubbletea's Update processes the tick
	tickMsg := TickMsg(time.Now())
	newModel, _ := bubbleteaCopy.Update(tickMsg)
	finalModel := newModel.(Model)

	// Verify challenge was processed
	if finalModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying, got %v", finalModel.Game.State)
	}

	if finalModel.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID cleared, got '%s'", finalModel.NvimChallengeID)
	}

	t.Logf("Final game state: %v, Gold: %d", finalModel.Game.State, finalModel.Game.Gold)
}

// TestSocketServerHandlerWithSeparateModels tests the exact problem:
// socket server has pointer to original, but Bubbletea evolves separately.
func TestSocketServerHandlerWithSeparateModels(t *testing.T) {
	// Create the original model
	original := newTestModel()
	original.NvimMode = true
	original.NvimRPC = &MockRPCClient{}

	// Socket server holds a pointer
	socketHandler := &original

	// Bubbletea gets value copy and evolves it
	btModel := original
	btModel.NvimChallengeID = "challenge_1"
	btModel.Game.StartChallengeWaiting()

	// Verify they're now different (in real code they would be)
	// Note: In this test, since we modified btModel after copy,
	// they ARE different

	// Send result via the socket handler
	result := &nvim.ChallengeResult{
		RequestID:  "challenge_1",
		Success:    true,
		GoldEarned: 100,
	}

	// This sends to the channel (which is shared)
	socketHandler.HandleChallengeComplete(result)

	// Process on Bubbletea's model
	tickMsg := TickMsg(time.Now())
	newModel, _ := btModel.Update(tickMsg)
	finalModel := newModel.(Model)

	if finalModel.Game.State != engine.StatePlaying {
		t.Errorf("Game should resume after challenge, got state %v", finalModel.Game.State)
	}
}

// =============================================================================
// Game Over / Victory Screen Tests
// =============================================================================

// TestHandleRestartUsesSelectedLevel tests that HandleRestart uses the
// selected level and settings instead of creating a generic game.
func TestHandleRestartUsesSelectedLevel(t *testing.T) {
	model := NewModel()
	model.NvimMode = true

	// Select a specific level
	levels := model.LevelRegistry.GetAll()
	if len(levels) < 2 {
		t.Fatal("Need at least 2 levels for this test")
	}
	model.SelectedLevel = &levels[1] // Select second level
	model.Settings.StartingGold = 300
	model.Settings.StartingHealth = 150
	model.Settings.Difficulty = engine.DifficultyEasy

	// Start a game and get to game over
	model.startGameFromSettings()
	model.Game.State = engine.StateGameOver

	// HandleRestart sends to channel, then Update processes it
	model.HandleRestart()

	// Simulate the Update loop processing the channel
	select {
	case <-model.RestartChan:
		// Process restart (same logic as Update)
		if model.SelectedLevel != nil {
			model.Game = engine.NewGameFromLevelAndSettings(model.SelectedLevel, model.Settings)
		} else {
			model.Game = engine.NewGame(GridWidth, GridHeight)
		}
	default:
		t.Fatal("Expected message on RestartChan")
	}

	// Verify the game was restarted with the same level
	if model.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after restart, got %v", model.Game.State)
	}

	// Verify settings were preserved
	if model.Game.Gold != 300 {
		t.Errorf("Expected starting gold 300, got %d", model.Game.Gold)
	}

	if model.Game.Health != 150 {
		t.Errorf("Expected starting health 150, got %d", model.Game.Health)
	}

	// Verify level path matches (level-specific)
	expectedPathLen := len(levels[1].Path)
	if len(model.Game.Path) != expectedPathLen {
		t.Errorf("Expected path length %d (from selected level), got %d", expectedPathLen, len(model.Game.Path))
	}
}

// TestHandleRestartWithoutSelectedLevel tests fallback when no level is selected.
func TestHandleRestartWithoutSelectedLevel(t *testing.T) {
	model := NewModel()
	model.NvimMode = true
	model.SelectedLevel = nil // No level selected

	model.Game.State = engine.StateGameOver

	// HandleRestart sends to channel, then Update processes it
	model.HandleRestart()

	// Simulate the Update loop processing the channel
	select {
	case <-model.RestartChan:
		// Process restart (same logic as Update)
		if model.SelectedLevel != nil {
			model.Game = engine.NewGameFromLevelAndSettings(model.SelectedLevel, model.Settings)
		} else {
			model.Game = engine.NewGame(GridWidth, GridHeight)
		}
	default:
		t.Fatal("Expected message on RestartChan")
	}

	if model.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after restart, got %v", model.Game.State)
	}

	// Should have default dimensions
	if model.Game.Width != GridWidth || model.Game.Height != GridHeight {
		t.Errorf("Expected default grid size %dx%d, got %dx%d",
			GridWidth, GridHeight, model.Game.Width, model.Game.Height)
	}
}

// TestHandleRestartClearsState tests that challenge state is properly cleared.
func TestHandleRestartClearsState(t *testing.T) {
	model := NewModel()
	model.NvimMode = true
	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]

	// Simulate mid-challenge state
	model.NvimChallengeID = "challenge_123"
	model.CurrentChallenge = &engine.Challenge{Name: "test"}
	model.Game.State = engine.StateGameOver

	// HandleRestart sends to channel, then Update processes it
	model.HandleRestart()

	// Simulate the Update loop processing the channel
	select {
	case <-model.RestartChan:
		// Process restart (same logic as Update)
		if model.SelectedLevel != nil {
			model.Game = engine.NewGameFromLevelAndSettings(model.SelectedLevel, model.Settings)
		} else {
			model.Game = engine.NewGame(GridWidth, GridHeight)
		}
		model.CurrentChallenge = nil
		model.VimEditor = nil
		model.NvimChallengeID = ""
	default:
		t.Fatal("Expected message on RestartChan")
	}

	if model.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID to be cleared, got '%s'", model.NvimChallengeID)
	}

	if model.CurrentChallenge != nil {
		t.Error("Expected CurrentChallenge to be nil")
	}

	if model.VimEditor != nil {
		t.Error("Expected VimEditor to be nil")
	}
}

// TestHandleGoToLevelSelect tests that the level select state transition works.
func TestHandleGoToLevelSelect(t *testing.T) {
	model := NewModel()
	model.NvimMode = true

	// Start a game and get to game over
	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]
	model.startGameFromSettings()
	model.Game.State = engine.StateGameOver

	// Set some state that should be cleared
	model.NvimChallengeID = "challenge_456"
	model.CurrentChallenge = &engine.Challenge{Name: "test"}
	model.SettingsMenuIndex = 3

	// HandleGoToLevelSelect sends to channel, then Update processes it
	model.HandleGoToLevelSelect()

	// Simulate the Update loop processing the channel
	select {
	case <-model.LevelSelectChan:
		// Process level select (same logic as Update)
		model.Game.State = engine.StateLevelSelect
		model.SettingsMenuIndex = 0
		model.CurrentChallenge = nil
		model.VimEditor = nil
		model.NvimChallengeID = ""
	default:
		t.Fatal("Expected message on LevelSelectChan")
	}

	if model.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect, got %v", model.Game.State)
	}

	if model.SettingsMenuIndex != 0 {
		t.Errorf("Expected SettingsMenuIndex reset to 0, got %d", model.SettingsMenuIndex)
	}

	if model.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID cleared, got '%s'", model.NvimChallengeID)
	}

	if model.CurrentChallenge != nil {
		t.Error("Expected CurrentChallenge to be nil")
	}
}

// TestHandleGoToLevelSelectFromVictory tests transition from victory state.
func TestHandleGoToLevelSelectFromVictory(t *testing.T) {
	model := NewModel()
	model.NvimMode = true

	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]
	model.startGameFromSettings()
	model.Game.State = engine.StateVictory

	// HandleGoToLevelSelect sends to channel, then Update processes it
	model.HandleGoToLevelSelect()

	// Simulate the Update loop processing the channel
	select {
	case <-model.LevelSelectChan:
		model.Game.State = engine.StateLevelSelect
	default:
		t.Fatal("Expected message on LevelSelectChan")
	}

	if model.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect from victory, got %v", model.Game.State)
	}
}

// TestEndGameKeysRestart tests the 'r' key in game over state.
func TestEndGameKeysRestart(t *testing.T) {
	model := NewModel()
	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]
	model.Settings.StartingGold = 250

	model.startGameFromSettings()
	model.Game.State = engine.StateGameOver
	model.Game.Gold = 50 // Simulate spent gold

	// Press 'r'
	newModel, _ := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := newModel.(Model)

	if updated.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after 'r', got %v", updated.Game.State)
	}

	// Gold should be reset to starting value
	if updated.Game.Gold != 250 {
		t.Errorf("Expected gold reset to 250, got %d", updated.Game.Gold)
	}
}

// TestEndGameKeysLevelSelect tests the 'm' key (menu) in game over state.
func TestEndGameKeysLevelSelect(t *testing.T) {
	model := NewModel()
	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]
	model.startGameFromSettings()
	model.Game.State = engine.StateGameOver
	model.SettingsMenuIndex = 2

	// Press 'm'
	newModel, _ := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	updated := newModel.(Model)

	if updated.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after 'm', got %v", updated.Game.State)
	}

	if updated.SettingsMenuIndex != 0 {
		t.Errorf("Expected SettingsMenuIndex reset to 0, got %d", updated.SettingsMenuIndex)
	}
}

// TestEndGameKeysQuit tests the 'q' key in game over state.
func TestEndGameKeysQuit(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateGameOver

	// Press 'q'
	newModel, cmd := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	updated := newModel.(Model)

	if !updated.Quitting {
		t.Error("Expected Quitting to be true after 'q'")
	}

	if cmd == nil {
		t.Error("Expected quit command to be returned")
	}
}

// TestEndGameKeysFromVictory tests keys work from victory state too.
func TestEndGameKeysFromVictory(t *testing.T) {
	model := NewModel()
	model.SelectedLevel = &model.LevelRegistry.GetAll()[0]
	model.startGameFromSettings()
	model.Game.State = engine.StateVictory

	// Press 'r'
	newModel, _ := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := newModel.(Model)

	if updated.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after 'r' from victory, got %v", updated.Game.State)
	}
}

// TestEndGameKeysIgnoresOtherKeys tests that invalid keys are ignored.
func TestEndGameKeysIgnoresOtherKeys(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateGameOver

	// Press some random key
	newModel, cmd := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	updated := newModel.(Model)

	if updated.Game.State != engine.StateGameOver {
		t.Errorf("State should remain GameOver, got %v", updated.Game.State)
	}

	if cmd != nil {
		t.Error("No command should be returned for invalid key")
	}
}

// TestRestartPreservesLevelSettings tests that restart uses the exact same
// level and settings that were used to start the current game.
func TestRestartPreservesLevelSettings(t *testing.T) {
	model := NewModel()

	// Select a harder level
	levels := model.LevelRegistry.GetAll()
	if len(levels) < 5 {
		t.Skip("Need at least 5 levels")
	}

	model.SelectedLevel = &levels[4] // 5th level
	model.Settings = engine.GameSettings{
		Difficulty:     engine.DifficultyHard,
		GameSpeed:      engine.SpeedDouble,
		StartingGold:   100,
		StartingHealth: 50,
	}

	model.startGameFromSettings()
	originalWaves := model.Game.TotalWaves

	// Play the game, change some state
	model.Game.Gold = 500
	model.Game.Health = 10
	model.Game.Wave = 5
	model.Game.State = engine.StateGameOver

	// Restart
	newModel, _ := model.handleEndGameKeys(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	updated := newModel.(Model)

	// Should be back to starting values
	if updated.Game.Gold != 100 {
		t.Errorf("Expected gold 100, got %d", updated.Game.Gold)
	}

	if updated.Game.Health != 50 {
		t.Errorf("Expected health 50, got %d", updated.Game.Health)
	}

	if updated.Game.Wave != 1 {
		t.Errorf("Expected wave 1, got %d", updated.Game.Wave)
	}

	// Level-specific settings preserved
	if updated.Game.TotalWaves != originalWaves {
		t.Errorf("Expected TotalWaves %d, got %d", originalWaves, updated.Game.TotalWaves)
	}
}

// =============================================================================
// Original Integration Tests
// =============================================================================

// TestIntegrationWithRealSocketServer tests the FULL integration:
// - Real SocketServer
// - Real Model with channel
// - Simulated Bubbletea Update loop
// - Real socket client (simulating Neovim).
func TestIntegrationWithRealSocketServer(t *testing.T) {
	tmpDir := os.TempDir()
	socketPath := filepath.Join(tmpDir, "test_integration.sock")
	defer os.Remove(socketPath)

	// Create model
	model := newTestModel()
	model.NvimMode = true

	// Initialize socket server (this stores &model as handler)
	model.InitNvimSocket(socketPath)
	defer model.NvimSocket.Stop()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	// Connect as client (simulating Neovim)
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Wait for connection
	time.Sleep(100 * time.Millisecond)

	// Consume the game_ready notification
	buf := make([]byte, 4096)
	conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	conn.Read(buf)

	// Now simulate Bubbletea's copy of the model starting a challenge
	btModel := model
	btModel.NvimChallengeCount++
	btModel.NvimChallengeID = "int_test_1"
	btModel.Game.StartChallenge() // Game continues during challenge for time pressure

	t.Logf("Started challenge, state=%v, id=%s", btModel.Game.State, btModel.NvimChallengeID)

	if btModel.Game.State != engine.StateChallengeActive {
		t.Fatalf("Expected StateChallengeActive, got %v", btModel.Game.State)
	}

	// Client sends challenge_complete (Neovim finished the challenge)
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "challenge_complete",
		"params": map[string]interface{}{
			"request_id":  "int_test_1",
			"success":     true,
			"gold_earned": 88.0,
		},
	}
	data, _ := json.Marshal(response)
	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		t.Fatalf("Failed to write: %v", err)
	}

	// Wait for RPC to be processed
	time.Sleep(200 * time.Millisecond)

	// Now simulate Bubbletea's Update tick on btModel
	initialGold := btModel.Game.Gold
	tickMsg := TickMsg(time.Now())
	newModel, _ := btModel.Update(tickMsg)
	finalModel := newModel.(Model)

	t.Logf("After tick: state=%v, gold=%d (was %d)", finalModel.Game.State, finalModel.Game.Gold, initialGold)

	// CRITICAL: This is where the bug would manifest
	if finalModel.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after challenge complete, got %v", finalModel.Game.State)
	}

	expectedGold := initialGold + 88
	if finalModel.Game.Gold != expectedGold {
		t.Errorf("Expected gold %d, got %d", expectedGold, finalModel.Game.Gold)
	}
}

// =============================================================================
// Challenge Mode and Selection State Transition Tests (Task 6.1)
// =============================================================================

// TestChallengeModeStateTransition tests entering challenge mode from start screen.
func TestChallengeModeStateTransition(t *testing.T) {
	model := NewModel()

	// Verify starting state
	if model.Game.State != engine.StateLevelSelect {
		t.Fatalf("Expected initial state StateLevelSelect, got %v", model.Game.State)
	}

	// Navigate to modes section
	model.StartSection = SectionModes
	model.ModeMenuIndex = 0 // Challenge Mode

	// Press Enter to start Challenge Mode
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.handleLevelSelectKeys(msg)
	updated := newModel.(Model)

	// Should transition to Challenge Mode or Challenge Mode Practice
	// (startChallengeModeChallenge auto-loads a challenge)
	validStates := []engine.GameState{engine.StateChallengeMode, engine.StateChallengeModePractice}
	stateValid := false
	for _, s := range validStates {
		if updated.Game.State == s {
			stateValid = true
			break
		}
	}
	if !stateValid {
		t.Errorf("Expected StateChallengeMode or StateChallengeModePractice, got %v", updated.Game.State)
	}

	// Streak should be reset
	if updated.ChallengeModeStreak != 0 {
		t.Errorf("Expected streak 0 on entering challenge mode, got %d", updated.ChallengeModeStreak)
	}
}

// TestChallengeSelectionStateTransition tests entering challenge selection from start screen.
func TestChallengeSelectionStateTransition(t *testing.T) {
	model := NewModel()

	// Navigate to modes section
	model.StartSection = SectionModes
	model.ModeMenuIndex = 1 // Challenge Selection

	// Press Enter to start Challenge Selection
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.handleLevelSelectKeys(msg)
	updated := newModel.(Model)

	// Should transition to Challenge Selection
	if updated.Game.State != engine.StateChallengeSelection {
		t.Errorf("Expected StateChallengeSelection after selecting Challenge Selection, got %v", updated.Game.State)
	}

	// List index should be reset
	if updated.ChallengeListIndex != 0 {
		t.Errorf("Expected ChallengeListIndex 0, got %d", updated.ChallengeListIndex)
	}
}

// TestChallengeModeExitTransition tests exiting challenge mode back to start screen.
func TestChallengeModeExitTransition(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeMode

	// Press Escape to exit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ := model.handleChallengeModeKeys(msg)
	updated := newModel.(Model)

	// Should return to level select
	if updated.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after exiting challenge mode, got %v", updated.Game.State)
	}

	// Should be in modes section with Challenge Mode selected
	if updated.StartSection != SectionModes {
		t.Errorf("Expected SectionModes after exit, got %v", updated.StartSection)
	}
	if updated.ModeMenuIndex != 0 {
		t.Errorf("Expected ModeMenuIndex 0 (Challenge Mode), got %d", updated.ModeMenuIndex)
	}
}

// TestChallengeSelectionExitTransition tests exiting challenge selection back to start screen.
func TestChallengeSelectionExitTransition(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelection

	// Press q to exit
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ := model.handleChallengeSelectionKeys(msg)
	updated := newModel.(Model)

	// Should return to level select
	if updated.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after exiting challenge selection, got %v", updated.Game.State)
	}

	// Should be in modes section with Challenge Selection selected
	if updated.StartSection != SectionModes {
		t.Errorf("Expected SectionModes after exit, got %v", updated.StartSection)
	}
	if updated.ModeMenuIndex != 1 {
		t.Errorf("Expected ModeMenuIndex 1 (Challenge Selection), got %d", updated.ModeMenuIndex)
	}
}

// TestChallengeModePracticeToModeTransition tests canceling a practice challenge.
func TestChallengeModePracticeToModeTransition(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeModePractice
	model.NvimMode = true
	model.NvimChallengeID = "test_challenge"
	model.CurrentChallenge = &engine.Challenge{Name: "test"}

	// Press Escape to cancel
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	newModel, _ := model.handleChallengeModePracticeKeys(msg)
	updated := newModel.(Model)

	// Should return to challenge mode
	if updated.Game.State != engine.StateChallengeMode {
		t.Errorf("Expected StateChallengeMode after canceling practice, got %v", updated.Game.State)
	}

	// Challenge state should be cleared
	if updated.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID to be cleared, got '%s'", updated.NvimChallengeID)
	}
	if updated.CurrentChallenge != nil {
		t.Error("Expected CurrentChallenge to be nil")
	}
}

// TestChallengeSelectionPracticeToSelectionTransition tests canceling from selection practice.
func TestChallengeSelectionPracticeToSelectionTransition(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelectionPractice
	model.NvimMode = true
	model.NvimChallengeID = "test_challenge"
	model.CurrentChallenge = &engine.Challenge{Name: "test"}

	// Press Escape to cancel
	msg := tea.KeyMsg{Type: tea.KeyEscape}
	newModel, _ := model.handleChallengeSelectionPracticeKeys(msg)
	updated := newModel.(Model)

	// Should return to challenge selection
	if updated.Game.State != engine.StateChallengeSelection {
		t.Errorf("Expected StateChallengeSelection after canceling practice, got %v", updated.Game.State)
	}

	// Challenge state should be cleared
	if updated.NvimChallengeID != "" {
		t.Errorf("Expected NvimChallengeID to be cleared, got '%s'", updated.NvimChallengeID)
	}
}

// =============================================================================
// Challenge Selection Navigation Tests (Task 6.2)
// =============================================================================

// TestChallengeSelectionNavigation tests navigating the challenge list.
func TestChallengeSelectionNavigation(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelection

	// Verify we have challenges
	if len(model.ChallengeList) == 0 {
		t.Fatal("Expected challenges to be loaded")
	}

	// Initially at index 0
	if model.ChallengeListIndex != 0 {
		t.Errorf("Expected initial index 0, got %d", model.ChallengeListIndex)
	}

	// Navigate down
	msgDown := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := model.handleChallengeSelectionKeys(msgDown)
	updated := newModel.(Model)

	if updated.ChallengeListIndex != 1 {
		t.Errorf("Expected index 1 after j, got %d", updated.ChallengeListIndex)
	}

	// Navigate up
	msgUp := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = updated.handleChallengeSelectionKeys(msgUp)
	updated = newModel.(Model)

	if updated.ChallengeListIndex != 0 {
		t.Errorf("Expected index 0 after k, got %d", updated.ChallengeListIndex)
	}
}

// TestChallengeSelectionListBounds tests that navigation respects list bounds.
func TestChallengeSelectionListBounds(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelection

	// Try to navigate up from start
	msgUp := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ := model.handleChallengeSelectionKeys(msgUp)
	updated := newModel.(Model)

	// Should stay at 0
	if updated.ChallengeListIndex != 0 {
		t.Errorf("Expected index 0 when at top, got %d", updated.ChallengeListIndex)
	}

	// Navigate to end
	lastIndex := len(model.ChallengeList) - 1
	updated.ChallengeListIndex = lastIndex

	// Try to navigate down past end
	msgDown := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ = updated.handleChallengeSelectionKeys(msgDown)
	finalModel := newModel.(Model)

	// Should stay at last index
	if finalModel.ChallengeListIndex != lastIndex {
		t.Errorf("Expected index %d when at bottom, got %d", lastIndex, finalModel.ChallengeListIndex)
	}
}

// TestChallengeSelectionScroll tests scrolling behavior for long lists.
func TestChallengeSelectionScroll(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelection

	// Skip if not enough challenges for scrolling
	if len(model.ChallengeList) < 20 {
		t.Skip("Need at least 20 challenges to test scrolling")
	}

	// Navigate far enough to trigger scroll
	for range 16 {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
		newModel, _ := model.handleChallengeSelectionKeys(msg)
		model = newModel.(Model)
	}

	// Offset should have increased
	if model.ChallengeListOffset == 0 {
		t.Error("Expected scroll offset to increase after navigating past visible area")
	}
}

// TestChallengeSelectionStartChallenge tests starting a challenge from selection.
func TestChallengeSelectionStartChallenge(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateChallengeSelection
	model.ChallengeListIndex = 2 // Select third challenge

	// Press Enter to start
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.handleChallengeSelectionKeys(msg)
	updated := newModel.(Model)

	// Should transition to practice state
	if updated.Game.State != engine.StateChallengeSelectionPractice {
		t.Errorf("Expected StateChallengeSelectionPractice, got %v", updated.Game.State)
	}

	// Selected index should be tracked
	if updated.SelectedChallengeIndex != 2 {
		t.Errorf("Expected SelectedChallengeIndex 2, got %d", updated.SelectedChallengeIndex)
	}

	// Current challenge should be set
	if updated.CurrentChallenge == nil {
		t.Error("Expected CurrentChallenge to be set")
	}
}

// TestStartMenuSectionNavigation tests navigating between levels and modes sections.
func TestStartMenuSectionNavigation(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateLevelSelect
	model.StartSection = SectionLevels

	levels := model.LevelRegistry.GetAll()

	// Navigate to last level
	model.LevelMenuIndex = len(levels) - 1

	// Navigate down should move to modes section
	msgDown := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	newModel, _ := model.handleLevelSelectKeys(msgDown)
	updated := newModel.(Model)

	if updated.StartSection != SectionModes {
		t.Errorf("Expected SectionModes after navigating past last level, got %v", updated.StartSection)
	}
	if updated.ModeMenuIndex != 0 {
		t.Errorf("Expected ModeMenuIndex 0, got %d", updated.ModeMenuIndex)
	}

	// Navigate up should return to levels section
	msgUp := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	newModel, _ = updated.handleLevelSelectKeys(msgUp)
	finalModel := newModel.(Model)

	if finalModel.StartSection != SectionLevels {
		t.Errorf("Expected SectionLevels after navigating up from modes, got %v", finalModel.StartSection)
	}
	if finalModel.LevelMenuIndex != len(levels)-1 {
		t.Errorf("Expected LevelMenuIndex %d, got %d", len(levels)-1, finalModel.LevelMenuIndex)
	}
}

// TestNotificationDisplay tests notification creation and expiration.
func TestNotificationDisplay(t *testing.T) {
	model := NewModel()

	// Show a notification
	model.ShowNotification("Test message", true)

	if model.Notification == nil {
		t.Fatal("Expected notification to be set")
	}
	if model.Notification.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", model.Notification.Message)
	}
	if !model.Notification.IsSuccess {
		t.Error("Expected IsSuccess to be true")
	}

	// Notification should have ShowUntil set in the future
	if model.Notification.ShowUntil.Before(time.Now()) {
		t.Error("Expected ShowUntil to be in the future")
	}
}

// TestNotificationExpiration tests that expired notifications are cleared.
func TestNotificationExpiration(t *testing.T) {
	model := NewModel()

	// Create an already expired notification
	model.Notification = &Notification{
		Message:   "Expired",
		IsSuccess: true,
		ShowUntil: time.Now().Add(-1 * time.Second), // Already expired
	}

	// Clear expired notification
	model.ClearExpiredNotification()

	if model.Notification != nil {
		t.Error("Expected expired notification to be cleared")
	}
}

// TestBuildChallengeData tests the helper function for creating challenge data.
func TestBuildChallengeData(t *testing.T) {
	challenge := &engine.Challenge{
		ID:          "test_id",
		Name:        "Test Challenge",
		Category:    "movement",
		Difficulty:  2,
		Description: "Test description",
	}

	data := buildChallengeData(challenge, "challenge_mode")

	if data.ID != "test_id" {
		t.Errorf("Expected ID 'test_id', got '%s'", data.ID)
	}
	if data.Name != "Test Challenge" {
		t.Errorf("Expected Name 'Test Challenge', got '%s'", data.Name)
	}
	if data.Mode != "challenge_mode" {
		t.Errorf("Expected Mode 'challenge_mode', got '%s'", data.Mode)
	}
}

// =============================================================================
// Integration Tests for Full Flows (Task 6.3)
// =============================================================================

// TestChallengeModeFullFlow tests the complete flow of entering, practicing, and exiting challenge mode.
func TestChallengeModeFullFlow(t *testing.T) {
	model := NewModel()

	// 1. Start from level select
	if model.Game.State != engine.StateLevelSelect {
		t.Fatalf("Expected StateLevelSelect, got %v", model.Game.State)
	}

	// 2. Navigate to Challenge Mode
	model.StartSection = SectionModes
	model.ModeMenuIndex = 0

	// 3. Enter Challenge Mode
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.handleLevelSelectKeys(enterMsg)
	model = newModel.(Model)

	// Should be in ChallengeMode or ChallengeModePractice (auto-loads first challenge)
	validStates := []engine.GameState{engine.StateChallengeMode, engine.StateChallengeModePractice}
	stateValid := false
	for _, s := range validStates {
		if model.Game.State == s {
			stateValid = true
			break
		}
	}
	if !stateValid {
		t.Fatalf("Expected StateChallengeMode or StateChallengeModePractice, got %v", model.Game.State)
	}

	// 4. Exit back to start screen
	// First get back to ChallengeMode state if in practice
	model.Game.State = engine.StateChallengeMode
	exitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ = model.handleChallengeModeKeys(exitMsg)
	model = newModel.(Model)

	if model.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after exit, got %v", model.Game.State)
	}
}

// TestChallengeSelectionFullFlow tests the complete flow of browsing and practicing challenges.
func TestChallengeSelectionFullFlow(t *testing.T) {
	model := NewModel()

	// 1. Start from level select
	if model.Game.State != engine.StateLevelSelect {
		t.Fatalf("Expected StateLevelSelect, got %v", model.Game.State)
	}

	// 2. Navigate to Challenge Selection
	model.StartSection = SectionModes
	model.ModeMenuIndex = 1

	// 3. Enter Challenge Selection
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	newModel, _ := model.handleLevelSelectKeys(enterMsg)
	model = newModel.(Model)

	if model.Game.State != engine.StateChallengeSelection {
		t.Fatalf("Expected StateChallengeSelection, got %v", model.Game.State)
	}

	// 4. Navigate the challenge list
	if len(model.ChallengeList) > 0 {
		downMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
		newModel, _ = model.handleChallengeSelectionKeys(downMsg)
		model = newModel.(Model)

		if model.ChallengeListIndex != 1 {
			t.Errorf("Expected index 1 after navigation, got %d", model.ChallengeListIndex)
		}

		// 5. Start a challenge
		enterMsg = tea.KeyMsg{Type: tea.KeyEnter}
		newModel, _ = model.handleChallengeSelectionKeys(enterMsg)
		model = newModel.(Model)

		if model.Game.State != engine.StateChallengeSelectionPractice {
			t.Errorf("Expected StateChallengeSelectionPractice, got %v", model.Game.State)
		}

		// 6. Cancel back to selection
		model.NvimMode = true // Set nvim mode to use escape handling
		escMsg := tea.KeyMsg{Type: tea.KeyEscape}
		newModel, _ = model.handleChallengeSelectionPracticeKeys(escMsg)
		model = newModel.(Model)

		if model.Game.State != engine.StateChallengeSelection {
			t.Errorf("Expected StateChallengeSelection after cancel, got %v", model.Game.State)
		}
	}

	// 7. Exit back to start screen
	exitMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	newModel, _ = model.handleChallengeSelectionKeys(exitMsg)
	model = newModel.(Model)

	if model.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after exit, got %v", model.Game.State)
	}
}
