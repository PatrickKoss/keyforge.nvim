package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
)

const (
	GridWidth  = 20
	GridHeight = 14
	TargetFPS  = 60
)

// TickMsg is sent on each frame update
type TickMsg time.Time

// Model is the bubbletea model for the game
type Model struct {
	Game       *engine.Game
	LastUpdate time.Time
	Width      int
	Height     int
	Quitting   bool
}

// NewModel creates a new game model
func NewModel() Model {
	return Model{
		Game:       engine.NewGame(GridWidth, GridHeight),
		LastUpdate: time.Now(),
		Width:      GridWidth,
		Height:     GridHeight,
		Quitting:   false,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tickCmd()
}

// tickCmd returns a command that sends tick messages at 60fps
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/TargetFPS, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case TickMsg:
		now := time.Time(msg)
		dt := now.Sub(m.LastUpdate).Seconds()
		m.LastUpdate = now
		m.Game.Update(dt)
		return m, tickCmd()

	case tea.WindowSizeMsg:
		// Could adapt to window size
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global keys
	switch msg.String() {
	case "ctrl+c", "q":
		if m.Game.State == engine.StatePlaying || m.Game.State == engine.StatePaused {
			m.Quitting = true
			return m, tea.Quit
		}
		// In game over or victory, q quits
		m.Quitting = true
		return m, tea.Quit

	case "r":
		// Restart game
		if m.Game.State == engine.StateGameOver || m.Game.State == engine.StateVictory {
			m.Game = engine.NewGame(GridWidth, GridHeight)
			m.LastUpdate = time.Now()
		}
		return m, nil
	}

	// Game state specific keys
	switch m.Game.State {
	case engine.StatePlaying:
		return m.handlePlayingKeys(msg)
	case engine.StatePaused:
		return m.handlePausedKeys(msg)
	}

	return m, nil
}

func (m Model) handlePlayingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	// Movement (vim keys)
	case "h", "left":
		m.Game.MoveCursor(-1, 0)
	case "j", "down":
		m.Game.MoveCursor(0, 1)
	case "k", "up":
		m.Game.MoveCursor(0, -1)
	case "l", "right":
		m.Game.MoveCursor(1, 0)

	// Tower selection
	case "1":
		m.Game.SelectTower(entities.TowerArrow)
	case "2":
		m.Game.SelectTower(entities.TowerLSP)
	case "3":
		m.Game.SelectTower(entities.TowerRefactor)

	// Actions
	case " ", "enter":
		m.Game.PlaceTower()
	case "u":
		m.Game.UpgradeTower()
	case "p":
		m.Game.TogglePause()
	}

	return m, nil
}

func (m Model) handlePausedKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "p", " ", "enter":
		m.Game.TogglePause()
	}
	return m, nil
}

// View renders the model
func (m Model) View() string {
	if m.Quitting {
		return "Thanks for playing Keyforge!\n"
	}

	switch m.Game.State {
	case engine.StateGameOver:
		return RenderGameOver(&m)
	case engine.StateVictory:
		return RenderVictory(&m)
	default:
		return RenderGame(&m)
	}
}
