package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
	"github.com/keyforge/keyforge/internal/nvim"
	"github.com/keyforge/keyforge/internal/vim"
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
	Game             *engine.Game
	LastUpdate       time.Time
	Width            int
	Height           int
	Quitting         bool
	ChallengeManager *engine.ChallengeManager
	CurrentChallenge *engine.Challenge
	VimEditor        *vim.Editor

	// Neovim integration
	NvimMode           bool
	NvimClient         *nvim.Client
	NvimChallengeID    string       // Current challenge request ID
	NvimChallengeCount int          // Counter for generating unique IDs
	PrevGameState      engine.GameState // Track state changes for notifications
}

// NewModel creates a new game model
func NewModel() Model {
	cm, _ := engine.NewChallengeManager()
	return Model{
		Game:             engine.NewGame(GridWidth, GridHeight),
		LastUpdate:       time.Now(),
		Width:            GridWidth,
		Height:           GridHeight,
		Quitting:         false,
		ChallengeManager: cm,
		CurrentChallenge: nil,
	}
}

// InitNvimClient initializes the Neovim RPC client
func (m *Model) InitNvimClient() {
	m.NvimClient = nvim.NewClient(m)
	m.NvimClient.Start()
	// Notify Neovim that game is ready
	m.NvimClient.SendGameReady()
}

// Handler interface implementation for nvim.Client

// HandleChallengeComplete processes challenge results from Neovim
func (m *Model) HandleChallengeComplete(result *nvim.ChallengeResult) {
	if result.RequestID != m.NvimChallengeID {
		return // Stale result, ignore
	}

	if result.Success {
		// Award gold - Neovim already calculated it
		gold := result.GoldEarned
		if gold < 1 {
			gold = 1
		}
		m.Game.AddChallengeGold(gold)
	}

	// Clear challenge state
	m.NvimChallengeID = ""
	m.Game.EndChallenge()
}

// HandleConfigUpdate processes config updates from Neovim
func (m *Model) HandleConfigUpdate(config *nvim.ConfigUpdate) {
	// Could apply difficulty settings etc.
}

// HandlePause pauses the game
func (m *Model) HandlePause() {
	if m.Game.State == engine.StatePlaying {
		m.Game.TogglePause()
	}
}

// HandleResume resumes the game
func (m *Model) HandleResume() {
	if m.Game.State == engine.StatePaused {
		m.Game.TogglePause()
	}
}

// HandleStartChallenge handles user-triggered challenge requests from Neovim
func (m *Model) HandleStartChallenge() {
	m.startChallenge()
}

// HandleRestart handles restart requests from Neovim
func (m *Model) HandleRestart() {
	m.Game = engine.NewGame(GridWidth, GridHeight)
	m.LastUpdate = time.Now()
	m.CurrentChallenge = nil
	m.VimEditor = nil
	m.NvimChallengeID = ""
	m.PrevGameState = engine.StatePlaying
}

// sendStateNotification sends game state notifications to Neovim
func (m *Model) sendStateNotification() {
	if m.NvimClient == nil {
		return
	}

	switch m.Game.State {
	case engine.StateGameOver:
		m.NvimClient.SendGameOver(
			m.Game.Wave,
			m.Game.Gold,
			len(m.Game.Towers),
			m.Game.Health,
		)
	case engine.StateVictory:
		m.NvimClient.SendVictory(
			m.Game.Wave,
			m.Game.Gold,
			len(m.Game.Towers),
			m.Game.Health,
		)
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

		// Store previous state before update
		prevState := m.Game.State
		m.Game.Update(dt)

		// Check for state changes and send notifications in nvim mode
		if m.NvimMode && m.NvimClient != nil && m.Game.State != prevState {
			m.sendStateNotification()
		}
		m.PrevGameState = m.Game.State

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
	case engine.StateChallengeActive:
		return m.handleChallengeKeys(msg)
	case engine.StateChallengeWaiting:
		return m.handleChallengeWaitingKeys(msg)
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

	// Challenge
	case "c":
		m.startChallenge()
	}

	return m, nil
}

// startChallenge starts a new challenge
func (m *Model) startChallenge() {
	if m.Game.ChallengeActive {
		return
	}

	// Get tower category for challenge selection
	tower := m.Game.GetTowerAt(m.Game.CursorX, m.Game.CursorY)
	var category string
	if tower != nil {
		info := tower.Info()
		category = info.Category
	}

	// In Neovim mode, delegate to Neovim via RPC
	if m.NvimMode && m.NvimClient != nil {
		m.NvimChallengeCount++
		m.NvimChallengeID = fmt.Sprintf("challenge_%d", m.NvimChallengeCount)

		// Request challenge from Neovim - game PAUSES while user edits
		difficulty := 1
		if m.Game.Wave >= 4 {
			difficulty = 2
		}
		if m.Game.Wave >= 7 {
			difficulty = 3
		}

		m.NvimClient.RequestChallenge(m.NvimChallengeID, category, difficulty)
		m.Game.StartChallengeWaiting() // Use waiting state (paused) for nvim mode
		return
	}

	// Standalone mode: use internal vim editor
	if m.ChallengeManager == nil {
		return
	}

	// Get a random challenge (matching tower category if available)
	challenge := m.ChallengeManager.GetRandomChallenge(category, m.Game.Wave)
	if challenge == nil {
		// Fallback to any challenge
		challenge = m.ChallengeManager.GetRandomChallenge("", 0)
	}
	if challenge == nil {
		return
	}

	m.CurrentChallenge = challenge

	// Initialize vim editor with challenge buffer
	m.VimEditor = vim.NewEditor(challenge.InitialBuffer)

	// Set initial cursor position if specified
	if len(challenge.CursorStart) == 2 {
		m.VimEditor.SetCursor(vim.Position{
			Line: challenge.CursorStart[0],
			Col:  challenge.CursorStart[1],
		})
	}

	m.Game.StartChallenge()
}

func (m Model) handlePausedKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "p", " ", "enter":
		m.Game.TogglePause()
	}
	return m, nil
}

// handleChallengeWaitingKeys handles input while waiting for nvim challenge result
// Game is paused - only allow cancel via Escape
func (m Model) handleChallengeWaitingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Only allow cancel via Escape (user cancelled in Neovim or wants to cancel here)
	if msg.String() == "esc" || msg.Type == tea.KeyEscape {
		m.NvimChallengeID = ""
		m.Game.EndChallenge()
	}
	return m, nil
}

func (m Model) handleChallengeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// This is for standalone mode only now
	// Nvim mode uses handleChallengeWaitingKeys
	if m.NvimMode {
		// Should not reach here in nvim mode, but handle gracefully
		return m, nil
	}

	// Standalone mode: use internal vim editor
	if m.VimEditor == nil {
		return m, nil
	}

	// Translate Bubbletea key to vim key string
	key := translateKey(msg)

	// Check for submit (Ctrl+S)
	if key == "ctrl+s" {
		m.submitChallenge()
		return m, nil
	}

	// Check for cancel (Escape in normal mode or Ctrl+C)
	if key == "ctrl+c" {
		m.completeChallenge(false)
		return m, nil
	}

	// In normal mode, Escape cancels the challenge
	if key == "Escape" && m.VimEditor.Mode == vim.ModeNormal {
		m.completeChallenge(false)
		return m, nil
	}

	// Forward key to vim editor
	m.VimEditor.HandleKey(key)

	return m, nil
}

// translateKey converts Bubbletea key to vim key string
func translateKey(msg tea.KeyMsg) string {
	switch msg.Type {
	case tea.KeyEscape:
		return "Escape"
	case tea.KeyEnter:
		return "Enter"
	case tea.KeyBackspace:
		return "Backspace"
	case tea.KeyDelete:
		return "Delete"
	case tea.KeyTab:
		return "Tab"
	case tea.KeySpace:
		return " "
	case tea.KeyUp:
		return "Up"
	case tea.KeyDown:
		return "Down"
	case tea.KeyLeft:
		return "Left"
	case tea.KeyRight:
		return "Right"
	case tea.KeyCtrlR:
		return "ctrl+r"
	case tea.KeyCtrlS:
		return "ctrl+s"
	case tea.KeyCtrlC:
		return "ctrl+c"
	default:
		return msg.String()
	}
}

// submitChallenge validates and completes the challenge
func (m *Model) submitChallenge() {
	if m.VimEditor == nil || m.CurrentChallenge == nil {
		return
	}

	spec := vim.ChallengeSpec{
		ValidationType:  m.CurrentChallenge.ValidationType,
		ExpectedBuffer:  m.CurrentChallenge.ExpectedBuffer,
		ExpectedContent: m.CurrentChallenge.ExpectedContent,
		ExpectedCursor:  m.CurrentChallenge.ExpectedCursor,
		InitialBuffer:   m.CurrentChallenge.InitialBuffer,
		ParKeystrokes:   m.CurrentChallenge.ParKeystrokes,
	}

	result := vim.Validate(m.VimEditor, spec)

	if result.Success {
		// Calculate gold based on efficiency
		gold := int(float64(m.CurrentChallenge.GoldBase) * result.Efficiency)
		if gold < 1 {
			gold = 1
		}
		m.Game.AddChallengeGold(gold)
	}

	m.VimEditor = nil
	m.CurrentChallenge = nil
	m.Game.EndChallenge()
}

// completeChallenge ends the current challenge
func (m *Model) completeChallenge(success bool) {
	if m.CurrentChallenge == nil {
		return
	}

	if success {
		// Award gold based on the challenge
		m.Game.AddChallengeGold(m.CurrentChallenge.GoldBase)
	}

	m.VimEditor = nil
	m.CurrentChallenge = nil
	m.Game.EndChallenge()
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
