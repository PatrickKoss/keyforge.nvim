package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/keyforge/keyforge/internal/engine"
)

// TestRenderStartScreen tests the level selection screen rendering.
func TestRenderStartScreen(t *testing.T) {
	model := NewModel()

	t.Run("renders logo", func(t *testing.T) {
		output := RenderStartScreen(&model)
		// Logo uses ASCII art with box-drawing characters
		if !strings.Contains(output, "██") {
			t.Error("Expected logo to contain ASCII art blocks")
		}
	})

	t.Run("renders level list", func(t *testing.T) {
		output := RenderStartScreen(&model)
		if !strings.Contains(output, "Select Level") {
			t.Error("Expected 'Select Level' title")
		}
		if !strings.Contains(output, "Classic") {
			t.Error("Expected 'Classic' level in list")
		}
	})

	t.Run("renders level preview", func(t *testing.T) {
		output := RenderStartScreen(&model)
		// Preview should show enemies and towers
		if !strings.Contains(output, "Enemies:") {
			t.Error("Expected 'Enemies:' in preview")
		}
		if !strings.Contains(output, "Towers:") {
			t.Error("Expected 'Towers:' in preview")
		}
	})

	t.Run("renders help text", func(t *testing.T) {
		output := RenderStartScreen(&model)
		if !strings.Contains(output, "[j/k]") {
			t.Error("Expected help text with j/k navigation")
		}
		if !strings.Contains(output, "[Enter]") {
			t.Error("Expected help text with Enter to configure")
		}
	})
}

// TestRenderSettingsScreen tests the settings menu rendering.
func TestRenderSettingsScreen(t *testing.T) {
	model := NewModel()
	model.SelectedLevel = model.LevelRegistry.GetByID("level-5") // Classic level

	t.Run("renders title", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Game Settings") {
			t.Error("Expected 'Game Settings' title")
		}
	})

	t.Run("renders level info", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Level: Classic") {
			t.Error("Expected selected level info")
		}
	})

	t.Run("renders difficulty options", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Difficulty") {
			t.Error("Expected Difficulty setting")
		}
		if !strings.Contains(output, "Easy") {
			t.Error("Expected Easy option")
		}
		if !strings.Contains(output, "Normal") {
			t.Error("Expected Normal option")
		}
		if !strings.Contains(output, "Hard") {
			t.Error("Expected Hard option")
		}
	})

	t.Run("renders game speed options", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Game Speed") {
			t.Error("Expected Game Speed setting")
		}
		if !strings.Contains(output, "0.5x") {
			t.Error("Expected 0.5x speed option")
		}
		if !strings.Contains(output, "2x") {
			t.Error("Expected 2x speed option")
		}
	})

	t.Run("renders gold slider", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Starting Gold") {
			t.Error("Expected Starting Gold setting")
		}
	})

	t.Run("renders health slider", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Starting Health") {
			t.Error("Expected Starting Health setting")
		}
	})

	t.Run("renders start game button", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "Start Game") {
			t.Error("Expected Start Game button")
		}
	})

	t.Run("renders help text", func(t *testing.T) {
		output := RenderSettingsScreen(&model)
		if !strings.Contains(output, "[h/l]") {
			t.Error("Expected help text with h/l navigation")
		}
		if !strings.Contains(output, "[Esc]") {
			t.Error("Expected help text with Esc to go back")
		}
	})
}

// TestLevelSelectNavigation tests keyboard navigation in level selection.
func TestLevelSelectNavigation(t *testing.T) {
	t.Run("j moves selection down", func(t *testing.T) {
		model := NewModel()
		model.LevelMenuIndex = 0

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m := newModel.(Model)

		// If there's more than one level, index should increase
		// Currently only one level, so it should stay at 0
		levels := m.LevelRegistry.GetAll()
		expectedIndex := 0
		if len(levels) > 1 {
			expectedIndex = 1
		}
		if m.LevelMenuIndex != expectedIndex {
			t.Errorf("Expected LevelMenuIndex %d, got %d", expectedIndex, m.LevelMenuIndex)
		}
	})

	t.Run("k moves selection up", func(t *testing.T) {
		model := NewModel()
		model.LevelMenuIndex = 1 // Start at second item (if exists)

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		m := newModel.(Model)

		if m.LevelMenuIndex != 0 {
			t.Errorf("Expected LevelMenuIndex 0, got %d", m.LevelMenuIndex)
		}
	})

	t.Run("down arrow moves selection down", func(t *testing.T) {
		model := NewModel()
		model.LevelMenuIndex = 0

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyDown})
		m := newModel.(Model)

		// With only one level, should stay at 0
		levels := m.LevelRegistry.GetAll()
		if len(levels) == 1 && m.LevelMenuIndex != 0 {
			t.Errorf("Expected LevelMenuIndex 0 with one level, got %d", m.LevelMenuIndex)
		}
	})

	t.Run("up arrow moves selection up", func(t *testing.T) {
		model := NewModel()
		model.LevelMenuIndex = 1

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyUp})
		m := newModel.(Model)

		if m.LevelMenuIndex != 0 {
			t.Errorf("Expected LevelMenuIndex 0, got %d", m.LevelMenuIndex)
		}
	})

	t.Run("enter goes to settings screen", func(t *testing.T) {
		model := NewModel()

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m := newModel.(Model)

		if m.Game.State != engine.StateSettings {
			t.Errorf("Expected StateSettings after Enter, got %v", m.Game.State)
		}
		if m.SelectedLevel == nil {
			t.Error("Expected SelectedLevel to be set")
		}
	})

	t.Run("q quits application", func(t *testing.T) {
		model := NewModel()

		_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

		// Should return tea.Quit command
		if cmd == nil {
			t.Error("Expected quit command")
		}
	})
}

// TestSettingsNavigation tests keyboard navigation in settings screen.
func TestSettingsNavigation(t *testing.T) {
	t.Run("j moves to next setting", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m := newModel.(Model)

		if m.SettingsMenuIndex != 1 {
			t.Errorf("Expected SettingsMenuIndex 1, got %d", m.SettingsMenuIndex)
		}
	})

	t.Run("k moves to previous setting", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 2

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		m := newModel.(Model)

		if m.SettingsMenuIndex != 1 {
			t.Errorf("Expected SettingsMenuIndex 1, got %d", m.SettingsMenuIndex)
		}
	})

	t.Run("cannot go below 0", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		m := newModel.(Model)

		if m.SettingsMenuIndex != 0 {
			t.Errorf("Expected SettingsMenuIndex 0, got %d", m.SettingsMenuIndex)
		}
	})

	t.Run("cannot go above max index", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 4 // Start Game button is last

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		m := newModel.(Model)

		if m.SettingsMenuIndex != 4 {
			t.Errorf("Expected SettingsMenuIndex 4, got %d", m.SettingsMenuIndex)
		}
	})

	t.Run("escape goes back to level select", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEsc})
		m := newModel.(Model)

		if m.Game.State != engine.StateLevelSelect {
			t.Errorf("Expected StateLevelSelect after Esc, got %v", m.Game.State)
		}
	})
}

// TestDifficultyAdjustment tests changing difficulty setting.
func TestDifficultyAdjustment(t *testing.T) {
	t.Run("l increases difficulty", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0 // Difficulty
		model.Settings.Difficulty = engine.DifficultyNormal

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.Difficulty != engine.DifficultyHard {
			t.Errorf("Expected difficulty Hard, got %s", m.Settings.Difficulty)
		}
	})

	t.Run("h decreases difficulty", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0 // Difficulty
		model.Settings.Difficulty = engine.DifficultyNormal

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.Difficulty != engine.DifficultyEasy {
			t.Errorf("Expected difficulty Easy, got %s", m.Settings.Difficulty)
		}
	})

	t.Run("cannot go below easy", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0
		model.Settings.Difficulty = engine.DifficultyEasy

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.Difficulty != engine.DifficultyEasy {
			t.Errorf("Expected difficulty Easy, got %s", m.Settings.Difficulty)
		}
	})

	t.Run("cannot go above hard", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 0
		model.Settings.Difficulty = engine.DifficultyHard

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.Difficulty != engine.DifficultyHard {
			t.Errorf("Expected difficulty Hard, got %s", m.Settings.Difficulty)
		}
	})
}

// TestGameSpeedAdjustment tests changing game speed setting.
func TestGameSpeedAdjustment(t *testing.T) {
	t.Run("l increases speed", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 1 // Game Speed
		model.Settings.GameSpeed = engine.SpeedNormal

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.GameSpeed != engine.SpeedFast {
			t.Errorf("Expected speed 1.5x, got %v", m.Settings.GameSpeed)
		}
	})

	t.Run("h decreases speed", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 1
		model.Settings.GameSpeed = engine.SpeedNormal

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.GameSpeed != engine.SpeedHalf {
			t.Errorf("Expected speed 0.5x, got %v", m.Settings.GameSpeed)
		}
	})

	t.Run("cannot go below 0.5x", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 1
		model.Settings.GameSpeed = engine.SpeedHalf

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.GameSpeed != engine.SpeedHalf {
			t.Errorf("Expected speed 0.5x, got %v", m.Settings.GameSpeed)
		}
	})

	t.Run("cannot go above 2x", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 1
		model.Settings.GameSpeed = engine.SpeedDouble

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.GameSpeed != engine.SpeedDouble {
			t.Errorf("Expected speed 2x, got %v", m.Settings.GameSpeed)
		}
	})
}

// TestGoldSliderAdjustment tests changing starting gold.
func TestGoldSliderAdjustment(t *testing.T) {
	t.Run("l increases gold by 25", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 2 // Starting Gold
		model.Settings.StartingGold = 200

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.StartingGold != 225 {
			t.Errorf("Expected gold 225, got %d", m.Settings.StartingGold)
		}
	})

	t.Run("h decreases gold by 25", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 2
		model.Settings.StartingGold = 200

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.StartingGold != 175 {
			t.Errorf("Expected gold 175, got %d", m.Settings.StartingGold)
		}
	})

	t.Run("cannot go below 100", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 2
		model.Settings.StartingGold = 100

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.StartingGold != 100 {
			t.Errorf("Expected gold 100, got %d", m.Settings.StartingGold)
		}
	})

	t.Run("cannot go above 500", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 2
		model.Settings.StartingGold = 500

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.StartingGold != 500 {
			t.Errorf("Expected gold 500, got %d", m.Settings.StartingGold)
		}
	})
}

// TestHealthSliderAdjustment tests changing starting health.
func TestHealthSliderAdjustment(t *testing.T) {
	t.Run("l increases health by 10", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 3 // Starting Health
		model.Settings.StartingHealth = 100

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.StartingHealth != 110 {
			t.Errorf("Expected health 110, got %d", m.Settings.StartingHealth)
		}
	})

	t.Run("h decreases health by 10", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 3
		model.Settings.StartingHealth = 100

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.StartingHealth != 90 {
			t.Errorf("Expected health 90, got %d", m.Settings.StartingHealth)
		}
	})

	t.Run("cannot go below 50", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 3
		model.Settings.StartingHealth = 50

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
		m := newModel.(Model)

		if m.Settings.StartingHealth != 50 {
			t.Errorf("Expected health 50, got %d", m.Settings.StartingHealth)
		}
	})

	t.Run("cannot go above 200", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 3
		model.Settings.StartingHealth = 200

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		m := newModel.(Model)

		if m.Settings.StartingHealth != 200 {
			t.Errorf("Expected health 200, got %d", m.Settings.StartingHealth)
		}
	})
}

// TestStartGameFromSettings tests starting the game from settings screen.
func TestStartGameFromSettings(t *testing.T) {
	t.Run("enter on Start Game button starts game", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 4 // Start Game button
		model.Settings = engine.GameSettings{
			Difficulty:     engine.DifficultyEasy,
			GameSpeed:      engine.SpeedDouble,
			StartingGold:   300,
			StartingHealth: 150,
		}

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m := newModel.(Model)

		if m.Game.State != engine.StatePlaying {
			t.Errorf("Expected StatePlaying, got %v", m.Game.State)
		}
	})

	t.Run("game uses selected settings", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 4
		model.Settings = engine.GameSettings{
			Difficulty:     engine.DifficultyEasy,
			GameSpeed:      engine.SpeedDouble,
			StartingGold:   300,
			StartingHealth: 150,
		}

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m := newModel.(Model)

		if m.Game.Gold != 300 {
			t.Errorf("Expected starting gold 300, got %d", m.Game.Gold)
		}
		if m.Game.Health != 150 {
			t.Errorf("Expected starting health 150, got %d", m.Game.Health)
		}
		if m.Game.GameSpeed != engine.SpeedDouble {
			t.Errorf("Expected game speed 2x, got %v", m.Game.GameSpeed)
		}
	})

	t.Run("game uses selected level", func(t *testing.T) {
		model := NewModel()
		model.Game.State = engine.StateSettings
		model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
		model.SettingsMenuIndex = 4
		model.Settings = engine.DefaultGameSettings()

		newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m := newModel.(Model)

		// Level 5 (Classic) has 8 waves
		if m.Game.TotalWaves != 8 {
			t.Errorf("Expected 8 total waves from Classic level, got %d", m.Game.TotalWaves)
		}
	})
}

// TestQuickStartPath tests Enter, Enter to start with defaults.
func TestQuickStartPath(t *testing.T) {
	model := NewModel()

	// First Enter: select level, go to settings
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := newModel.(Model)

	if m.Game.State != engine.StateSettings {
		t.Fatalf("Expected StateSettings after first Enter, got %v", m.Game.State)
	}

	// Move to Start Game button (index 4)
	m.SettingsMenuIndex = 4

	// Second Enter: start game
	newModel2, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := newModel2.(Model)

	if m2.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after second Enter, got %v", m2.Game.State)
	}

	// Verify defaults are used
	defaults := engine.DefaultGameSettings()
	if m2.Game.Gold != defaults.StartingGold {
		t.Errorf("Expected default gold %d, got %d", defaults.StartingGold, m2.Game.Gold)
	}
	if m2.Game.Health != defaults.StartingHealth {
		t.Errorf("Expected default health %d, got %d", defaults.StartingHealth, m2.Game.Health)
	}
}

// TestMiniGridPreview tests the mini grid rendering.
func TestMiniGridPreview(t *testing.T) {
	level := engine.ClassicLevel()

	output := renderMiniGrid(&level)

	t.Run("has expected height", func(t *testing.T) {
		lines := strings.Split(output, "\n")
		if len(lines) != 7 {
			t.Errorf("Expected 7 lines, got %d", len(lines))
		}
	})

	t.Run("contains path characters", func(t *testing.T) {
		if !strings.Contains(output, "░") {
			t.Error("Expected path characters in mini grid")
		}
	})
}

// TestSliderRendering tests slider bar rendering.
func TestSliderRendering(t *testing.T) {
	t.Run("renders at minimum", func(t *testing.T) {
		output := renderSlider(100, 100, 500)
		if !strings.HasPrefix(output, "[") || !strings.HasSuffix(output, "]") {
			t.Error("Expected slider to be enclosed in brackets")
		}
	})

	t.Run("renders at maximum", func(t *testing.T) {
		output := renderSlider(500, 100, 500)
		if !strings.Contains(output, "█") {
			t.Error("Expected filled characters at maximum")
		}
	})

	t.Run("renders at middle", func(t *testing.T) {
		output := renderSlider(300, 100, 500)
		// Should have both filled and unfilled characters
		if !strings.Contains(output, "█") || !strings.Contains(output, "░") {
			t.Error("Expected both filled and unfilled characters at middle value")
		}
	})
}

// TestDifficultyIndex tests mapping difficulty string to index.
func TestDifficultyIndex(t *testing.T) {
	tests := []struct {
		difficulty string
		expected   int
	}{
		{engine.DifficultyEasy, 0},
		{engine.DifficultyNormal, 1},
		{engine.DifficultyHard, 2},
		{"invalid", 1}, // defaults to normal
	}

	for _, tc := range tests {
		t.Run(tc.difficulty, func(t *testing.T) {
			result := difficultyIndex(tc.difficulty)
			if result != tc.expected {
				t.Errorf("Expected index %d for %s, got %d", tc.expected, tc.difficulty, result)
			}
		})
	}
}

// TestSpeedIndex tests mapping GameSpeed to index.
func TestSpeedIndex(t *testing.T) {
	tests := []struct {
		speed    engine.GameSpeed
		expected int
	}{
		{engine.SpeedHalf, 0},
		{engine.SpeedNormal, 1},
		{engine.SpeedFast, 2},
		{engine.SpeedDouble, 3},
		{engine.GameSpeed(0.75), 1}, // invalid defaults to normal
	}

	for _, tc := range tests {
		t.Run(tc.speed.String(), func(t *testing.T) {
			result := speedIndex(tc.speed)
			if result != tc.expected {
				t.Errorf("Expected index %d, got %d", tc.expected, result)
			}
		})
	}
}

// TestDifficultyIcon tests the star icons for level difficulty.
func TestDifficultyIcon(t *testing.T) {
	tests := []struct {
		difficulty engine.LevelDifficulty
		expected   string
	}{
		{engine.LevelDifficultyBeginner, "★☆☆"},
		{engine.LevelDifficultyIntermediate, "★★☆"},
		{engine.LevelDifficultyAdvanced, "★★★"},
	}

	for _, tc := range tests {
		t.Run(string(tc.difficulty), func(t *testing.T) {
			result := difficultyIcon(tc.difficulty)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

// TestGameOverMenuKey tests 'm' during game over goes back to menu.
func TestGameOverMenuKey(t *testing.T) {
	model := newTestModel()
	model.Game.State = engine.StateGameOver

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m := newModel.(Model)

	if m.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after 'm' from game over, got %v", m.Game.State)
	}
}

// TestGameVictoryMenuKey tests 'm' during victory goes back to menu.
func TestGameVictoryMenuKey(t *testing.T) {
	model := newTestModel()
	model.Game.State = engine.StateVictory

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	m := newModel.(Model)

	if m.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after 'm' from victory, got %v", m.Game.State)
	}
}

// TestQuitDuringPlayingReturnsToStartScreen tests 'q' during playing goes back to start screen.
func TestQuitDuringPlayingReturnsToStartScreen(t *testing.T) {
	model := newTestModel()
	model.Game.State = engine.StatePlaying

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m := newModel.(Model)

	if m.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after 'q' from playing, got %v", m.Game.State)
	}
}

// TestQuitDuringPausedReturnsToStartScreen tests 'q' while paused goes back to start screen.
func TestQuitDuringPausedReturnsToStartScreen(t *testing.T) {
	model := newTestModel()
	model.Game.State = engine.StatePaused

	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m := newModel.(Model)

	if m.Game.State != engine.StateLevelSelect {
		t.Errorf("Expected StateLevelSelect after 'q' from paused, got %v", m.Game.State)
	}
}

// TestGameRestartKeepsLevelAndSettings tests 'r' restarts with same settings.
func TestGameRestartKeepsLevelAndSettings(t *testing.T) {
	model := NewModel()
	model.Game.State = engine.StateSettings
	model.SelectedLevel = model.LevelRegistry.GetByID("level-5")
	model.Settings = engine.GameSettings{
		Difficulty:     engine.DifficultyEasy,
		GameSpeed:      engine.SpeedDouble,
		StartingGold:   300,
		StartingHealth: 150,
	}

	// Start the game
	model.SettingsMenuIndex = 4
	newModel, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m := newModel.(Model)

	// Simulate game over
	m.Game.State = engine.StateGameOver

	// Restart with 'r'
	newModel2, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	m2 := newModel2.(Model)

	// Should be playing again with same settings
	if m2.Game.State != engine.StatePlaying {
		t.Errorf("Expected StatePlaying after restart, got %v", m2.Game.State)
	}
	if m2.Game.Gold != 300 {
		t.Errorf("Expected starting gold 300 after restart, got %d", m2.Game.Gold)
	}
	if m2.Game.Health != 150 {
		t.Errorf("Expected starting health 150 after restart, got %d", m2.Game.Health)
	}
	if m2.Game.GameSpeed != engine.SpeedDouble {
		t.Errorf("Expected game speed 2x after restart, got %v", m2.Game.GameSpeed)
	}
}
