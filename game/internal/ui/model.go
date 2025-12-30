package ui

import (
	"fmt"
	"strings"
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

	keyEnter  = "enter"
	keySpace  = " "
	keyCtrlC  = "ctrl+c"
	keyCtrlS  = "ctrl+s"
	keyDown   = "down"
	keyUp     = "up"
	keyEsc    = "esc"
	keyEscape = "Escape"
)

// TickMsg is sent on each frame update.
type TickMsg time.Time

// StartMenuSection indicates which section of the start menu is active.
type StartMenuSection int

const (
	SectionLevels StartMenuSection = iota
	SectionModes
)

// Notification represents a temporary notification to display.
type Notification struct {
	Message   string
	IsSuccess bool
	ShowUntil time.Time
}

// ChallengeFeedback holds feedback from the previous challenge to send with next request.
type ChallengeFeedback struct {
	Success bool
	Streak  int
	Gold    int
}

// Model is the bubbletea model for the game.
type Model struct {
	Game              *engine.Game
	LastUpdate        time.Time
	Width             int
	Height            int
	Quitting          bool
	ChallengeManager  *engine.ChallengeManager
	ChallengeSelector *engine.ChallengeSelector
	CurrentChallenge  *engine.Challenge
	VimEditor         *vim.Editor

	// Start screen state
	LevelRegistry     *engine.LevelRegistry
	SelectedLevel     *engine.Level
	Settings          engine.GameSettings
	LevelMenuIndex    int // Selected level in level browser
	SettingsMenuIndex int // Selected setting in settings menu

	// Start menu section navigation
	StartSection  StartMenuSection // Which section is active (levels or modes)
	ModeMenuIndex int              // Selected mode (0 = Challenge Mode, 1 = Challenge Selection)

	// Challenge mode state
	ChallengeModeStreak int           // Successful challenges in a row
	Notification        *Notification // Current notification to display

	// Challenge selection state
	ChallengeList          []engine.Challenge // All loaded challenges for selection
	ChallengeListIndex     int                // Currently hovered challenge in list
	ChallengeListOffset    int                // Scroll offset for long lists
	SelectedChallengeIndex int                // Which challenge is being practiced

	// Neovim integration
	NvimMode           bool
	NvimClient         *nvim.Client       // Legacy stdin/stderr RPC
	NvimSocket         *nvim.SocketServer // Unix socket RPC (preferred)
	NvimRPC            nvim.RPCClient     // Interface to whichever is active
	NvimChallengeID    string             // Current challenge request ID
	NvimChallengeCount int                // Counter for generating unique IDs
	PrevGameState      engine.GameState   // Track state changes for notifications
	// Pending feedback to send with next challenge request
	PendingFeedback *ChallengeFeedback

	// Channels for RPC commands (thread-safe communication with Update loop)
	ChallengeResultChan chan *nvim.ChallengeResult
	RestartChan         chan struct{}
	LevelSelectChan     chan struct{}
}

// NewModel creates a new game model with default settings.
func NewModel() Model {
	return NewModelWithSettings(engine.DefaultGameSettings())
}

// NewModelWithSettings creates a new game model with specified settings.
func NewModelWithSettings(settings engine.GameSettings) Model {
	cm, _ := engine.NewChallengeManager()
	cs := engine.NewChallengeSelector(cm)
	registry := engine.NewLevelRegistry()

	// Get the first level as default selection
	levels := registry.GetAll()
	var selectedLevel *engine.Level
	if len(levels) > 0 {
		selectedLevel = &levels[0]
	}

	// Create a placeholder game (will be replaced when starting from settings)
	// Start in level select state
	game := engine.NewGame(GridWidth, GridHeight)
	game.State = engine.StateLevelSelect

	// Build challenge list for selection mode
	var challengeList []engine.Challenge
	if cm != nil {
		challengeList = cm.GetAllChallenges()
	}

	return Model{
		Game:                game,
		LastUpdate:          time.Now(),
		Width:               GridWidth,
		Height:              GridHeight,
		Quitting:            false,
		ChallengeManager:    cm,
		ChallengeSelector:   cs,
		CurrentChallenge:    nil,
		LevelRegistry:       registry,
		SelectedLevel:       selectedLevel,
		Settings:            settings,
		LevelMenuIndex:      0,
		SettingsMenuIndex:   0,
		StartSection:        SectionLevels,
		ModeMenuIndex:       0,
		ChallengeModeStreak: 0,
		ChallengeList:       challengeList,
		ChallengeListIndex:  0,
		ChallengeListOffset: 0,
		ChallengeResultChan: make(chan *nvim.ChallengeResult, 10),
		RestartChan:         make(chan struct{}, 1),
		LevelSelectChan:     make(chan struct{}, 1),
	}
}

// InitNvimClient initializes the Neovim RPC client (legacy stdin/stderr).
func (m *Model) InitNvimClient() {
	m.NvimClient = nvim.NewClient(m)
	m.NvimClient.Start()
	m.NvimRPC = m.NvimClient // Use Client as the RPC interface
	// Notify Neovim that game is ready
	if err := m.NvimClient.SendGameReady(); err != nil {
		m.Game.SetStatusMessage("Failed to send game ready notification")
	}
}

// InitNvimSocket initializes the Neovim RPC via Unix socket.
func (m *Model) InitNvimSocket(socketPath string) {
	m.NvimSocket = nvim.NewSocketServer(socketPath, m)
	if err := m.NvimSocket.Start(); err != nil {
		// Fall back to no RPC - game still works standalone
		return
	}
	m.NvimRPC = m.NvimSocket // Use SocketServer as the RPC interface
}

// Handler interface implementation for nvim.Client

// HandleChallengeComplete processes challenge results from Neovim
// This is called from a goroutine, so we send to a channel for processing in the Update loop
// NOTE: This is called on the original model pointer, not Bubbletea's copies, so we cannot
// reliably check m.NvimChallengeID here. Instead, we pass the result to the channel and
// let the Update loop (which has the current state) decide whether to process it.
func (m *Model) HandleChallengeComplete(result *nvim.ChallengeResult) {
	// Send to channel for processing in the main Update loop (thread-safe)
	// The Update loop will check if this is still the active challenge
	select {
	case m.ChallengeResultChan <- result:
		// Successfully sent
	default:
		// Channel full, drop the result (shouldn't happen with buffered channel)
	}
}

// HandleConfigUpdate processes config updates from Neovim.
func (m *Model) HandleConfigUpdate(config *nvim.ConfigUpdate) {
	// Could apply difficulty settings etc.
}

// HandlePause pauses the game.
func (m *Model) HandlePause() {
	if m.Game.State == engine.StatePlaying {
		m.Game.TogglePause()
	}
}

// HandleResume resumes the game.
func (m *Model) HandleResume() {
	if m.Game.State == engine.StatePaused {
		m.Game.TogglePause()
	}
}

// HandleStartChallenge handles user-triggered challenge requests from Neovim.
func (m *Model) HandleStartChallenge() {
	m.startChallenge()
}

// HandleRestart handles restart requests from Neovim.
// Sends to channel for processing in Update loop (thread-safe with Bubbletea).
func (m *Model) HandleRestart() {
	select {
	case m.RestartChan <- struct{}{}:
	default:
	}
}

// HandleGoToLevelSelect handles level select requests from Neovim.
// Sends to channel for processing in Update loop (thread-safe with Bubbletea).
func (m *Model) HandleGoToLevelSelect() {
	select {
	case m.LevelSelectChan <- struct{}{}:
	default:
	}
}

// ShowNotification displays a notification for a short time.
func (m *Model) ShowNotification(message string, isSuccess bool) {
	m.Notification = &Notification{
		Message:   message,
		IsSuccess: isSuccess,
		ShowUntil: time.Now().Add(2 * time.Second),
	}
}

// ClearExpiredNotification clears the notification if it has expired.
func (m *Model) ClearExpiredNotification() {
	if m.Notification != nil && time.Now().After(m.Notification.ShowUntil) {
		m.Notification = nil
	}
}

// sendStateNotification sends game state notifications to Neovim.
func (m *Model) sendStateNotification() {
	if m.NvimRPC == nil {
		return
	}

	switch m.Game.State {
	case engine.StateGameOver:
		if err := m.NvimRPC.SendGameOver(
			m.Game.Wave,
			m.Game.Gold,
			len(m.Game.Towers),
			m.Game.Health,
		); err != nil {
			// RPC error - game continues but notification failed
			m.Game.SetStatusMessage("Failed to notify Neovim of game over")
		}
	case engine.StateVictory:
		if err := m.NvimRPC.SendVictory(
			m.Game.Wave,
			m.Game.Gold,
			len(m.Game.Towers),
			m.Game.Health,
		); err != nil {
			// RPC error - game continues but notification failed
			m.Game.SetStatusMessage("Failed to notify Neovim of victory")
		}
	default:
		// Other states don't need notifications
	}
}

// Init initializes the model.
// Bubbletea requires value receiver for this interface method.
func (m Model) Init() tea.Cmd { //nolint:gocritic // hugeParam: required by Bubbletea interface
	return tickCmd()
}

// tickCmd returns a command that sends tick messages at 60fps.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/TargetFPS, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Update handles messages and updates the model.
// Bubbletea requires value receiver for this interface method.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: required by Bubbletea interface
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case TickMsg:
		now := time.Time(msg)
		dt := now.Sub(m.LastUpdate).Seconds()
		m.LastUpdate = now

		// Process pending challenge results from RPC channel (non-blocking)
		select {
		case result := <-m.ChallengeResultChan:
			m.handleChallengeResult(result)
		default:
			// No pending result, continue
		}

		// Process restart commands from RPC (non-blocking)
		select {
		case <-m.RestartChan:
			// Restart with same level and settings
			if m.SelectedLevel != nil {
				m.Game = engine.NewGameFromLevelAndSettings(m.SelectedLevel, m.Settings)
			} else {
				m.Game = engine.NewGame(GridWidth, GridHeight)
			}
			m.LastUpdate = now
			m.CurrentChallenge = nil
			m.VimEditor = nil
			m.NvimChallengeID = ""
			m.PrevGameState = engine.StatePlaying
			if m.ChallengeSelector != nil {
				m.ChallengeSelector.Reset()
			}
		default:
		}

		// Process level select commands from RPC (non-blocking)
		select {
		case <-m.LevelSelectChan:
			m.Game.State = engine.StateLevelSelect
			m.SettingsMenuIndex = 0
			m.CurrentChallenge = nil
			m.VimEditor = nil
			m.NvimChallengeID = ""
		default:
		}

		// Clear expired notifications
		m.ClearExpiredNotification()

		// Store previous state before update
		prevState := m.Game.State
		m.Game.Update(dt)

		// Check for state changes and send notifications in nvim mode
		if m.NvimMode && m.NvimRPC != nil && m.Game.State != prevState {
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

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	// Global quit
	if msg.String() == keyCtrlC {
		m.Quitting = true
		return m, tea.Quit
	}

	// Game state specific keys
	switch m.Game.State {
	case engine.StateLevelSelect:
		return m.handleLevelSelectKeys(msg)
	case engine.StateSettings:
		return m.handleSettingsKeys(msg)
	case engine.StatePlaying:
		return m.handlePlayingKeys(msg)
	case engine.StatePaused:
		return m.handlePausedKeys(msg)
	case engine.StateChallengeActive:
		return m.handleChallengeKeys(msg)
	case engine.StateChallengeWaiting:
		return m.handleChallengeWaitingKeys(msg)
	case engine.StateGameOver, engine.StateVictory:
		return m.handleEndGameKeys(msg)
	case engine.StateChallengeMode:
		return m.handleChallengeModeKeys(msg)
	case engine.StateChallengeModePractice:
		return m.handleChallengeModePracticeKeys(msg)
	case engine.StateChallengeSelection:
		return m.handleChallengeSelectionKeys(msg)
	case engine.StateChallengeSelectionPractice:
		return m.handleChallengeSelectionPracticeKeys(msg)
	case engine.StateMenu, engine.StateWaveComplete:
		// Fallthrough to quit handling
	}

	// Quit handling for other states
	if msg.String() == "q" {
		m.Quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) handleLevelSelectKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	levels := m.LevelRegistry.GetAll()
	numModes := 2 // Challenge Mode and Challenge Selection

	switch msg.String() {
	case "j", keyDown:
		if m.StartSection == SectionLevels {
			if m.LevelMenuIndex < len(levels)-1 {
				m.LevelMenuIndex++
			} else {
				// Move to modes section
				m.StartSection = SectionModes
				m.ModeMenuIndex = 0
			}
		} else {
			// In modes section
			if m.ModeMenuIndex < numModes-1 {
				m.ModeMenuIndex++
			}
		}
	case "k", keyUp:
		if m.StartSection == SectionModes {
			if m.ModeMenuIndex > 0 {
				m.ModeMenuIndex--
			} else {
				// Move back to levels section
				m.StartSection = SectionLevels
				m.LevelMenuIndex = len(levels) - 1
			}
		} else {
			// In levels section
			if m.LevelMenuIndex > 0 {
				m.LevelMenuIndex--
			}
		}
	case keyEnter, keySpace:
		if m.StartSection == SectionLevels {
			// Select level and go to settings
			if m.LevelMenuIndex < len(levels) {
				m.SelectedLevel = &levels[m.LevelMenuIndex]
				m.Game.State = engine.StateSettings
				m.SettingsMenuIndex = 0
			}
		} else {
			// Select mode
			if m.ModeMenuIndex == 0 {
				// Challenge Mode
				m.Game.State = engine.StateChallengeMode
				m.ChallengeModeStreak = 0
				m.Notification = nil
				m.startChallengeModeChallenge()
			} else {
				// Challenge Selection
				m.Game.State = engine.StateChallengeSelection
				m.ChallengeListIndex = 0
				m.ChallengeListOffset = 0
			}
		}
	case "q":
		m.Quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) handleSettingsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	maxIndex := 4 // 0: difficulty, 1: speed, 2: gold, 3: health, 4: start button

	switch msg.String() {
	case "j", keyDown:
		if m.SettingsMenuIndex < maxIndex {
			m.SettingsMenuIndex++
		}
	case "k", keyUp:
		if m.SettingsMenuIndex > 0 {
			m.SettingsMenuIndex--
		}
	case "h", "left":
		m.adjustSetting(-1)
	case "l", "right":
		m.adjustSetting(1)
	case keyEnter, keySpace:
		if m.SettingsMenuIndex == 4 {
			// Start game
			m.startGameFromSettings()
		}
	case keyEsc:
		// Back to level select
		m.Game.State = engine.StateLevelSelect
	case "q":
		m.Quitting = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *Model) adjustSetting(delta int) {
	switch m.SettingsMenuIndex {
	case 0: // Difficulty
		difficulties := []string{engine.DifficultyEasy, engine.DifficultyNormal, engine.DifficultyHard}
		idx := 1 // default normal
		for i, d := range difficulties {
			if d == m.Settings.Difficulty {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = 0
		}
		if idx >= len(difficulties) {
			idx = len(difficulties) - 1
		}
		m.Settings.Difficulty = difficulties[idx]
	case 1: // Speed
		speeds := engine.GameSpeedOptions()
		idx := 1 // default 1x
		for i, s := range speeds {
			if s == m.Settings.GameSpeed {
				idx = i
				break
			}
		}
		idx += delta
		if idx < 0 {
			idx = 0
		}
		if idx >= len(speeds) {
			idx = len(speeds) - 1
		}
		m.Settings.GameSpeed = speeds[idx]
	case 2: // Starting Gold
		m.Settings.StartingGold += delta * 25
		if m.Settings.StartingGold < 100 {
			m.Settings.StartingGold = 100
		}
		if m.Settings.StartingGold > 500 {
			m.Settings.StartingGold = 500
		}
	case 3: // Starting Health
		m.Settings.StartingHealth += delta * 10
		if m.Settings.StartingHealth < 50 {
			m.Settings.StartingHealth = 50
		}
		if m.Settings.StartingHealth > 200 {
			m.Settings.StartingHealth = 200
		}
	}
}

func (m *Model) startGameFromSettings() {
	if m.SelectedLevel == nil {
		return
	}
	m.Game = engine.NewGameFromLevelAndSettings(m.SelectedLevel, m.Settings)
	m.LastUpdate = time.Now()
	m.CurrentChallenge = nil
	m.VimEditor = nil
	m.NvimChallengeID = ""
	m.PrevGameState = engine.StatePlaying
	if m.ChallengeSelector != nil {
		m.ChallengeSelector.Reset()
	}

	// Notify Neovim that game is ready (if in nvim mode)
	if m.NvimMode && m.NvimRPC != nil {
		_ = m.NvimRPC.SendGameReady()
	}
}

func (m Model) handleEndGameKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	switch msg.String() {
	case "r":
		// Restart with same level and settings
		if m.SelectedLevel != nil {
			m.Game = engine.NewGameFromLevelAndSettings(m.SelectedLevel, m.Settings)
			m.LastUpdate = time.Now()
			if m.ChallengeSelector != nil {
				m.ChallengeSelector.Reset()
			}
		}
	case "m":
		// Return to menu
		m.Game.State = engine.StateLevelSelect
		m.SettingsMenuIndex = 0
	case "q":
		m.Quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handlePlayingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
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

	// Quit to start screen
	case "q":
		m.Game.State = engine.StateLevelSelect
		m.SettingsMenuIndex = 0
	}

	return m, nil
}

// startChallenge starts a new challenge.
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
	if m.NvimMode && m.NvimRPC != nil {
		m.NvimChallengeCount++
		m.NvimChallengeID = fmt.Sprintf("challenge_%d", m.NvimChallengeCount)

		// Use selector to pick challenge (with variety and anti-repetition)
		difficulty := 1
		if m.Game.Wave >= 4 {
			difficulty = 2
		}
		if m.Game.Wave >= 7 {
			difficulty = 3
		}

		// Get challenge from selector
		var challengeData *nvim.ChallengeData
		if m.ChallengeSelector != nil {
			challenge := m.ChallengeSelector.GetChallenge(category, difficulty)
			if challenge == nil {
				challenge = m.ChallengeSelector.GetChallenge("", 0)
			}
			if challenge != nil {
				challengeData = &nvim.ChallengeData{
					ID:              challenge.ID,
					Name:            challenge.Name,
					Category:        challenge.Category,
					Difficulty:      challenge.Difficulty,
					Description:     challenge.Description,
					InitialBuffer:   challenge.InitialBuffer,
					ExpectedBuffer:  challenge.ExpectedBuffer,
					ValidationType:  challenge.ValidationType,
					ExpectedCursor:  challenge.ExpectedCursor,
					ExpectedContent: challenge.ExpectedContent,
					FunctionName:    challenge.FunctionName,
					CursorStart:     challenge.CursorStart,
					ParKeystrokes:   challenge.ParKeystrokes,
					GoldBase:        challenge.GoldBase,
					Filetype:        challenge.Filetype,
					HintAction:      challenge.HintAction,
					HintFallback:    challenge.HintFallback,
				}
			}
		}

		if err := m.NvimRPC.RequestChallenge(m.NvimChallengeID, challengeData); err != nil {
			// Failed to request challenge, don't enter waiting state
			m.NvimChallengeID = ""
			return
		}
		m.Game.StartChallenge() // Game continues during challenge for time pressure
		return
	}

	// Standalone mode: use internal vim editor
	if m.ChallengeSelector == nil {
		return
	}

	// Get a challenge using the selector (avoids repetition, ensures variety)
	challenge := m.ChallengeSelector.GetChallenge(category, m.Game.Wave)
	if challenge == nil {
		// Fallback to any challenge
		challenge = m.ChallengeSelector.GetChallenge("", 0)
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

func (m Model) handlePausedKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	switch msg.String() {
	case "p", " ", "enter":
		m.Game.TogglePause()
	case "q":
		m.Game.State = engine.StateLevelSelect
		m.SettingsMenuIndex = 0
	}
	return m, nil
}

// handleChallengeWaitingKeys handles input while waiting for nvim challenge result
// Game is paused - only allow cancel via Escape.
func (m Model) handleChallengeWaitingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	// Only allow cancel via Escape (user canceled in Neovim or wants to cancel here)
	if msg.String() == "esc" || msg.Type == tea.KeyEscape {
		m.NvimChallengeID = ""
		m.Game.EndChallenge()
	}
	return m, nil
}

func (m Model) handleChallengeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
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
	if key == keyCtrlS {
		m.submitChallenge()
		return m, nil
	}

	// Check for cancel (Escape in normal mode or Ctrl+C)
	if key == keyCtrlC {
		m.completeChallenge(false)
		return m, nil
	}

	// In normal mode, Escape cancels the challenge
	if key == keyEscape && m.VimEditor.Mode == vim.ModeNormal {
		m.completeChallenge(false)
		return m, nil
	}

	// Forward key to vim editor
	m.VimEditor.HandleKey(key)

	return m, nil
}

// translateKey converts Bubbletea key to vim key string.
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

// submitChallenge validates and completes the challenge.
func (m *Model) submitChallenge() {
	if m.VimEditor == nil || m.CurrentChallenge == nil {
		return
	}

	spec := &vim.ChallengeSpec{
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

// completeChallenge ends the current challenge.
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

// handleChallengeResult processes a challenge result from Neovim RPC.
func (m *Model) handleChallengeResult(result *nvim.ChallengeResult) {
	// Check if this is for the current challenge (ignore stale results)
	if result.RequestID != m.NvimChallengeID {
		return // Stale result, ignore
	}

	// Award gold if successful
	if result.Success {
		gold := result.GoldEarned
		if gold < 1 {
			gold = 1
		}
		m.Game.AddChallengeGold(gold)
	}

	// Handle based on mode (determined by challenge ID prefix)
	m.NvimChallengeID = ""

	if strings.HasPrefix(result.RequestID, "challenge_mode_") {
		// Challenge Mode: update streak and start next challenge
		if result.Success {
			m.ChallengeModeStreak++
			m.ShowNotification("Success!", true)
		} else {
			m.ChallengeModeStreak = 0
			m.ShowNotification("Try again!", false)
		}
		// Set feedback to send with next challenge
		m.PendingFeedback = &ChallengeFeedback{
			Success: result.Success,
			Streak:  m.ChallengeModeStreak,
			Gold:    result.GoldEarned,
		}
		m.CurrentChallenge = nil
		m.Game.State = engine.StateChallengeMode
		m.startChallengeModeChallenge()
	} else if strings.HasPrefix(result.RequestID, "challenge_selection_") {
		// Challenge Selection: show result and start next challenge
		if result.Success {
			m.ShowNotification("Success!", true)
		} else {
			m.ShowNotification("Try again!", false)
		}
		// Set feedback to send with next challenge
		m.PendingFeedback = &ChallengeFeedback{
			Success: result.Success,
			Gold:    result.GoldEarned,
		}
		m.CurrentChallenge = nil
		// Move to next challenge in list
		m.SelectedChallengeIndex++
		if m.SelectedChallengeIndex >= len(m.ChallengeList) {
			m.SelectedChallengeIndex = 0
		}
		m.ChallengeListIndex = m.SelectedChallengeIndex
		m.Game.State = engine.StateChallengeSelection
		m.startChallengeSelectionChallenge()
	} else {
		// Tower defense mode: return to playing
		m.Game.EndChallenge()
	}
}

// buildChallengeData creates a ChallengeData struct from a Challenge.
func buildChallengeData(challenge *engine.Challenge, mode string) *nvim.ChallengeData {
	return &nvim.ChallengeData{
		ID:              challenge.ID,
		Name:            challenge.Name,
		Category:        challenge.Category,
		Difficulty:      challenge.Difficulty,
		Description:     challenge.Description,
		InitialBuffer:   challenge.InitialBuffer,
		ExpectedBuffer:  challenge.ExpectedBuffer,
		ValidationType:  challenge.ValidationType,
		ExpectedCursor:  challenge.ExpectedCursor,
		ExpectedContent: challenge.ExpectedContent,
		FunctionName:    challenge.FunctionName,
		CursorStart:     challenge.CursorStart,
		ParKeystrokes:   challenge.ParKeystrokes,
		GoldBase:        challenge.GoldBase,
		Filetype:        challenge.Filetype,
		HintAction:      challenge.HintAction,
		HintFallback:    challenge.HintFallback,
		Mode:            mode,
	}
}

// initVimEditor initializes the vim editor with a challenge buffer.
func (m *Model) initVimEditor(challenge *engine.Challenge) {
	m.VimEditor = vim.NewEditor(challenge.InitialBuffer)
	if len(challenge.CursorStart) == 2 {
		m.VimEditor.SetCursor(vim.Position{
			Line: challenge.CursorStart[0],
			Col:  challenge.CursorStart[1],
		})
	}
}

// startChallengeModeChallenge starts a random challenge for challenge mode.
func (m *Model) startChallengeModeChallenge() {
	if m.ChallengeSelector == nil {
		return
	}

	// Get a random challenge
	challenge := m.ChallengeSelector.GetChallenge("", 0)
	if challenge == nil {
		return
	}

	m.CurrentChallenge = challenge
	m.Game.State = engine.StateChallengeModePractice

	// In Neovim mode, send challenge to Neovim
	if m.NvimMode && m.NvimRPC != nil {
		m.NvimChallengeCount++
		m.NvimChallengeID = fmt.Sprintf("challenge_mode_%d", m.NvimChallengeCount)

		challengeData := buildChallengeData(challenge, "challenge_mode")
		// Include feedback from previous challenge if available
		if m.PendingFeedback != nil {
			challengeData.PrevSuccess = &m.PendingFeedback.Success
			challengeData.PrevStreak = m.PendingFeedback.Streak
			challengeData.PrevGold = m.PendingFeedback.Gold
			m.PendingFeedback = nil // Clear after using
		}
		if err := m.NvimRPC.RequestChallenge(m.NvimChallengeID, challengeData); err != nil {
			m.NvimChallengeID = ""
			m.CurrentChallenge = nil
			m.Game.State = engine.StateChallengeMode
			return
		}
		return
	}

	// Standalone mode: use internal vim editor
	m.initVimEditor(challenge)
}

// startChallengeSelectionChallenge starts the selected challenge from the list.
func (m *Model) startChallengeSelectionChallenge() {
	if m.ChallengeListIndex >= len(m.ChallengeList) {
		return
	}

	challenge := &m.ChallengeList[m.ChallengeListIndex]
	m.CurrentChallenge = challenge
	m.SelectedChallengeIndex = m.ChallengeListIndex
	m.Game.State = engine.StateChallengeSelectionPractice

	// In Neovim mode, send challenge to Neovim
	if m.NvimMode && m.NvimRPC != nil {
		m.NvimChallengeCount++
		m.NvimChallengeID = fmt.Sprintf("challenge_selection_%d", m.NvimChallengeCount)

		challengeData := buildChallengeData(challenge, "challenge_selection")
		// Include feedback from previous challenge if available
		if m.PendingFeedback != nil {
			challengeData.PrevSuccess = &m.PendingFeedback.Success
			challengeData.PrevGold = m.PendingFeedback.Gold
			m.PendingFeedback = nil // Clear after using
		}
		if err := m.NvimRPC.RequestChallenge(m.NvimChallengeID, challengeData); err != nil {
			m.NvimChallengeID = ""
			m.CurrentChallenge = nil
			m.Game.State = engine.StateChallengeSelection
			return
		}
		return
	}

	// Standalone mode: use internal vim editor
	m.initVimEditor(challenge)
}

// handleChallengeModeKeys handles input in challenge mode.
func (m Model) handleChallengeModeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	switch msg.String() {
	case keyEsc, "q":
		// Return to level select
		m.Game.State = engine.StateLevelSelect
		m.StartSection = SectionModes
		m.ModeMenuIndex = 0
		m.CurrentChallenge = nil
		m.VimEditor = nil
		m.NvimChallengeID = ""
		m.Notification = nil
	}
	return m, nil
}

// handleChallengeModePracticeKeys handles input while practicing a challenge in challenge mode.
func (m Model) handleChallengeModePracticeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	// In Nvim mode, just handle escape to cancel
	if m.NvimMode {
		if msg.String() == keyEsc || msg.Type == tea.KeyEscape {
			m.NvimChallengeID = ""
			m.CurrentChallenge = nil
			m.Game.State = engine.StateChallengeMode
		}
		return m, nil
	}

	// Standalone mode: use internal vim editor
	if m.VimEditor == nil {
		return m, nil
	}

	key := translateKey(msg)

	// Submit challenge
	if key == keyCtrlS {
		m.submitChallengeModeChallenge()
		return m, nil
	}

	// Cancel challenge
	if key == keyCtrlC || (key == keyEscape && m.VimEditor.Mode == vim.ModeNormal) {
		m.CurrentChallenge = nil
		m.VimEditor = nil
		m.Game.State = engine.StateChallengeMode
		return m, nil
	}

	m.VimEditor.HandleKey(key)
	return m, nil
}

// submitChallengeModeChallenge validates and handles challenge completion in challenge mode.
func (m *Model) submitChallengeModeChallenge() {
	if m.VimEditor == nil || m.CurrentChallenge == nil {
		return
	}

	spec := &vim.ChallengeSpec{
		ValidationType:  m.CurrentChallenge.ValidationType,
		ExpectedBuffer:  m.CurrentChallenge.ExpectedBuffer,
		ExpectedContent: m.CurrentChallenge.ExpectedContent,
		ExpectedCursor:  m.CurrentChallenge.ExpectedCursor,
		InitialBuffer:   m.CurrentChallenge.InitialBuffer,
		ParKeystrokes:   m.CurrentChallenge.ParKeystrokes,
	}

	result := vim.Validate(m.VimEditor, spec)

	if result.Success {
		m.ChallengeModeStreak++
		m.ShowNotification("Success!", true)
	} else {
		m.ChallengeModeStreak = 0
		m.ShowNotification("Try again!", false)
	}

	m.VimEditor = nil
	m.CurrentChallenge = nil
	m.Game.State = engine.StateChallengeMode

	// Start next challenge automatically after a short delay (handled in view)
	m.startChallengeModeChallenge()
}

// handleChallengeSelectionKeys handles input in challenge selection mode.
func (m Model) handleChallengeSelectionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	maxVisible := 15 // Number of visible challenges in list

	switch msg.String() {
	case "j", keyDown:
		if m.ChallengeListIndex < len(m.ChallengeList)-1 {
			m.ChallengeListIndex++
			// Scroll if needed
			if m.ChallengeListIndex >= m.ChallengeListOffset+maxVisible {
				m.ChallengeListOffset = m.ChallengeListIndex - maxVisible + 1
			}
		}
	case "k", keyUp:
		if m.ChallengeListIndex > 0 {
			m.ChallengeListIndex--
			// Scroll if needed
			if m.ChallengeListIndex < m.ChallengeListOffset {
				m.ChallengeListOffset = m.ChallengeListIndex
			}
		}
	case keyEnter, keySpace:
		// Start the selected challenge
		m.startChallengeSelectionChallenge()
	case keyEsc, "q":
		// Return to level select
		m.Game.State = engine.StateLevelSelect
		m.StartSection = SectionModes
		m.ModeMenuIndex = 1
		m.Notification = nil
	}
	return m, nil
}

// handleChallengeSelectionPracticeKeys handles input while practicing a selected challenge.
func (m Model) handleChallengeSelectionPracticeKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) { //nolint:gocritic // hugeParam: returns modified model
	// In Nvim mode, handle escape to cancel or 'b' to go back
	if m.NvimMode {
		switch msg.String() {
		case keyEsc:
			m.NvimChallengeID = ""
			m.CurrentChallenge = nil
			m.Game.State = engine.StateChallengeSelection
		case "b":
			// Back to selection list
			m.NvimChallengeID = ""
			m.CurrentChallenge = nil
			m.Game.State = engine.StateChallengeSelection
		}
		return m, nil
	}

	// Standalone mode: use internal vim editor
	if m.VimEditor == nil {
		return m, nil
	}

	key := translateKey(msg)

	// Submit challenge
	if key == keyCtrlS {
		m.submitChallengeSelectionChallenge()
		return m, nil
	}

	// Back to selection
	if key == keyCtrlC || (key == keyEscape && m.VimEditor.Mode == vim.ModeNormal) {
		m.CurrentChallenge = nil
		m.VimEditor = nil
		m.Game.State = engine.StateChallengeSelection
		return m, nil
	}

	m.VimEditor.HandleKey(key)
	return m, nil
}

// submitChallengeSelectionChallenge validates and handles challenge completion in selection mode.
func (m *Model) submitChallengeSelectionChallenge() {
	if m.VimEditor == nil || m.CurrentChallenge == nil {
		return
	}

	spec := &vim.ChallengeSpec{
		ValidationType:  m.CurrentChallenge.ValidationType,
		ExpectedBuffer:  m.CurrentChallenge.ExpectedBuffer,
		ExpectedContent: m.CurrentChallenge.ExpectedContent,
		ExpectedCursor:  m.CurrentChallenge.ExpectedCursor,
		InitialBuffer:   m.CurrentChallenge.InitialBuffer,
		ParKeystrokes:   m.CurrentChallenge.ParKeystrokes,
	}

	result := vim.Validate(m.VimEditor, spec)

	if result.Success {
		m.ShowNotification("Success!", true)
	} else {
		m.ShowNotification("Try again!", false)
	}

	m.VimEditor = nil
	m.CurrentChallenge = nil

	// Move to next challenge in list
	m.SelectedChallengeIndex++
	if m.SelectedChallengeIndex >= len(m.ChallengeList) {
		m.SelectedChallengeIndex = 0 // Wrap around
	}
	m.ChallengeListIndex = m.SelectedChallengeIndex

	// Start next challenge
	m.startChallengeSelectionChallenge()
}

// View renders the model.
// Bubbletea requires value receiver for this interface method.
func (m Model) View() string { //nolint:gocritic // hugeParam: required by Bubbletea interface
	if m.Quitting {
		return "Thanks for playing Keyforge!\n"
	}

	switch m.Game.State {
	case engine.StateLevelSelect:
		return RenderStartScreen(&m)
	case engine.StateSettings:
		return RenderSettingsScreen(&m)
	case engine.StateGameOver:
		// In Nvim mode, Lua handles the popup overlay
		if m.NvimMode {
			return "\n"
		}
		return RenderGameOver(&m)
	case engine.StateVictory:
		// In Nvim mode, Lua handles the popup overlay
		if m.NvimMode {
			return "\n"
		}
		return RenderVictory(&m)
	case engine.StateChallengeMode, engine.StateChallengeModePractice:
		return RenderChallengeMode(&m)
	case engine.StateChallengeSelection:
		return RenderChallengeSelection(&m)
	case engine.StateChallengeSelectionPractice:
		return RenderChallengeSelectionPractice(&m)
	default:
		return RenderGame(&m)
	}
}
