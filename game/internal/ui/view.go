package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
	"github.com/keyforge/keyforge/internal/vim"
)

// RenderGame renders the complete game view
func RenderGame(m *Model) string {
	var b strings.Builder

	// Title bar
	b.WriteString(renderTitle(m))
	b.WriteString("\n")

	// HUD
	b.WriteString(renderHUD(m))
	b.WriteString("\n")

	// Game grid
	b.WriteString(renderGrid(m))
	b.WriteString("\n")

	// Challenge display (when active)
	if m.Game.State == engine.StateChallengeActive && m.CurrentChallenge != nil {
		b.WriteString(renderChallenge(m))
		b.WriteString("\n")
	} else {
		// Shop
		b.WriteString(renderShop(m))
		b.WriteString("\n")
	}

	// Help
	b.WriteString(renderHelp(m))

	return b.String()
}

func renderTitle(m *Model) string {
	title := "═══════════════════════ KEYFORGE ═══════════════════════"
	return TitleStyle.Render(title)
}

func renderHUD(m *Model) string {
	g := m.Game

	// Wave info
	waveInfo := WaveStyle.Render(fmt.Sprintf("Wave: %d/%d", g.Wave, g.TotalWaves))

	// Gold
	goldInfo := GoldStyle.Render(fmt.Sprintf("💰 Gold: %d", g.Gold))

	// Health bar
	healthBar := RenderHealthBar(g.Health, g.MaxHealth, 10)
	healthInfo := HealthStyle.Render(fmt.Sprintf("❤️  Health: %d/%d ", g.Health, g.MaxHealth)) + healthBar

	// Status
	var status string
	switch g.State {
	case engine.StatePaused:
		status = PausedStyle.Render("  [PAUSED]")
	case engine.StateGameOver:
		status = GameOverStyle.Render("  [GAME OVER]")
	case engine.StateVictory:
		status = VictoryStyle.Render("  [VICTORY!]")
	case engine.StateChallengeActive:
		status = ChallengeStyle.Render("  [CHALLENGE ACTIVE - Game continues!]")
	case engine.StateChallengeWaiting:
		status = ChallengeStyle.Render("  [CHALLENGE IN PROGRESS - Game paused]")
	}

	// Challenge hint when not in challenge
	var challengeHint string
	if g.State == engine.StatePlaying && !g.ChallengeActive {
		challengeHint = HelpStyle.Render("  [Press c for challenge]")
	}

	hud := fmt.Sprintf("%s    %s    %s%s%s", waveInfo, goldInfo, healthInfo, status, challengeHint)
	return HUDStyle.Render(hud)
}

func renderGrid(m *Model) string {
	g := m.Game
	var b strings.Builder

	// Top border
	b.WriteString(BoxTopLeft)
	for i := 0; i < g.Width*2; i++ {
		b.WriteString(BoxHorizontal)
	}
	b.WriteString(BoxTopRight)
	b.WriteString("\n")

	// Create a 2D grid for rendering
	grid := make([][]string, g.Height)
	for y := 0; y < g.Height; y++ {
		grid[y] = make([]string, g.Width)
		for x := 0; x < g.Width; x++ {
			grid[y][x] = EmptyCellStyle.Render(EmptyCell)
		}
	}

	// Render path
	for _, p := range g.Path {
		x, y := p.IntPos()
		if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
			grid[y][x] = PathCellStyle.Render(PathChar)
		}
	}

	// Render towers
	for _, tower := range g.Towers {
		x, y := tower.Pos.IntPos()
		if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
			grid[y][x] = renderTower(tower)
		}
	}

	// Render enemies
	for _, enemy := range g.Enemies {
		if enemy.Dead {
			continue
		}
		x, y := enemy.Pos.IntPos()
		if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
			grid[y][x] = renderEnemy(enemy)
		}
	}

	// Render projectiles
	for _, proj := range g.Projectiles {
		if proj.Done {
			continue
		}
		x, y := proj.Pos.IntPos()
		if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
			grid[y][x] = ProjectileStyle.Render(ProjectileChar)
		}
	}

	// Render effects (on top of everything except cursor)
	for _, effect := range g.Effects.Effects {
		if effect.Done {
			continue
		}
		x, y := effect.Pos.IntPos()
		if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
			frame := effect.CurrentFrame()
			if frame != "" {
				info := effect.Info()
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(info.Color))
				grid[y][x] = style.Render(frame)
			}
		}
	}

	// Render cursor
	if g.State == engine.StatePlaying || g.State == engine.StatePaused {
		if g.CursorX >= 0 && g.CursorX < g.Width && g.CursorY >= 0 && g.CursorY < g.Height {
			// Show cursor with different style based on placement validity
			if g.CanPlaceTower(g.CursorX, g.CursorY) {
				grid[g.CursorY][g.CursorX] = CursorStyle.Render("▪")
			} else {
				grid[g.CursorY][g.CursorX] = lipgloss.NewStyle().
					Background(ColorDanger).
					Foreground(lipgloss.Color("#ffffff")).
					Render("✗")
			}
		}
	}

	// Build grid rows
	for y := 0; y < g.Height; y++ {
		b.WriteString(BoxVertical)
		for x := 0; x < g.Width; x++ {
			cell := grid[y][x]
			b.WriteString(cell)
			// Pad to fixed width (2 chars per cell)
			// Emojis are width 2, regular chars are width 1
			cellWidth := lipgloss.Width(cell)
			if cellWidth < 2 {
				b.WriteString(" ")
			}
		}
		b.WriteString(BoxVertical)
		b.WriteString("\n")
	}

	// Bottom border
	b.WriteString(BoxBottomLeft)
	for i := 0; i < g.Width*2; i++ {
		b.WriteString(BoxHorizontal)
	}
	b.WriteString(BoxBottomRight)

	return b.String()
}

func renderTower(tower *entities.Tower) string {
	info := tower.Info()
	var style lipgloss.Style
	var char string

	switch tower.Type {
	case entities.TowerArrow:
		style = TowerArrowStyle
		char = TowerArrowChar
	case entities.TowerLSP:
		style = TowerLSPStyle
		char = TowerLSPChar
	case entities.TowerRefactor:
		style = TowerRefactorStyle
		char = TowerRefactorChar
	default:
		style = TowerArrowStyle
		char = info.Symbol
	}

	// Add level indicator
	if tower.Level > 0 {
		char = fmt.Sprintf("%s", char) // could add level marker
	}

	return style.Render(char)
}

func renderEnemy(enemy *entities.Enemy) string {
	var style lipgloss.Style
	var char string

	switch enemy.Type {
	case entities.EnemyBug:
		style = EnemyBugStyle
		char = EnemyBugChar
	case entities.EnemyGremlin:
		style = EnemyGremlinStyle
		char = EnemyGremlinChar
	case entities.EnemyDaemon:
		style = EnemyDaemonStyle
		char = EnemyDaemonChar
	case entities.EnemyBoss:
		style = EnemyBossStyle
		char = EnemyBossChar
	default:
		style = EnemyBugStyle
		char = EnemyCharASCII
	}

	return style.Render(char)
}

func renderShop(m *Model) string {
	g := m.Game
	var items []string

	towers := []entities.TowerType{
		entities.TowerArrow,
		entities.TowerLSP,
		entities.TowerRefactor,
	}

	for i, towerType := range towers {
		info := entities.TowerTypes[towerType]
		canAfford := g.Gold >= info.Cost
		isSelected := g.SelectedTower == towerType

		text := fmt.Sprintf("[%d] %s %s (%dg)", i+1, info.Symbol, info.Name, info.Cost)

		var style lipgloss.Style
		if isSelected {
			style = SelectedStyle
		} else if canAfford {
			style = ShopItemStyle
		} else {
			style = ShopItemDisabledStyle
		}

		items = append(items, style.Render(text))
	}

	// Upgrade option if on tower
	tower := g.GetTowerAt(g.CursorX, g.CursorY)
	if tower != nil && tower.CanUpgrade() {
		cost := tower.UpgradeCost()
		canAfford := g.Gold >= cost
		text := fmt.Sprintf("[u] Upgrade (%dg)", cost)
		if canAfford {
			items = append(items, ShopItemStyle.Render(text))
		} else {
			items = append(items, ShopItemDisabledStyle.Render(text))
		}
	}

	return strings.Join(items, "  ")
}

func renderChallenge(m *Model) string {
	c := m.CurrentChallenge
	if c == nil {
		return ""
	}

	var b strings.Builder

	// Challenge header
	header := fmt.Sprintf("Challenge: %s (%s)", c.Name, c.Category)
	b.WriteString(ChallengeStyle.Render(header))
	b.WriteString("\n")

	// Description
	b.WriteString(HelpStyle.Render(c.Description))
	b.WriteString("\n\n")

	// Render vim editor if available, otherwise show static buffer
	if m.VimEditor != nil {
		b.WriteString(renderVimBuffer(m.VimEditor))
		b.WriteString("\n\n")

		// Mode line
		b.WriteString(renderModeLine(m.VimEditor, c))
	} else if c.InitialBuffer != "" {
		b.WriteString(HelpStyle.Render("Buffer:"))
		b.WriteString("\n")
		bufferStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#333333")).
			Padding(0, 1)
		b.WriteString(bufferStyle.Render(strings.TrimSpace(c.InitialBuffer)))
		b.WriteString("\n\n")

		// Gold reward
		goldInfo := fmt.Sprintf("Reward: %dg  |  Par: %d keystrokes", c.GoldBase, c.ParKeystrokes)
		b.WriteString(GoldStyle.Render(goldInfo))
	}

	return b.String()
}

// renderVimBuffer renders the vim buffer with cursor
func renderVimBuffer(e *vim.Editor) string {
	state := e.GetRenderState()
	var b strings.Builder

	bufferBg := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1)

	var lines []string
	for lineNum, line := range state.Lines {
		renderedLine := renderVimLine(line, lineNum, state)
		lines = append(lines, renderedLine)
	}

	content := strings.Join(lines, "\n")
	b.WriteString(bufferBg.Render(content))

	return b.String()
}

// renderVimLine renders a single line with cursor highlighting
func renderVimLine(line string, lineNum int, state vim.RenderState) string {
	runes := []rune(line)

	// Handle empty line
	if len(runes) == 0 {
		if lineNum == state.CursorLine {
			// Show cursor on empty line
			if state.Mode == vim.ModeInsert {
				return InsertCursorStyle.Render("|")
			}
			return NormalCursorStyle.Render(" ")
		}
		return " " // Empty line placeholder
	}

	// Not the cursor line - check for visual selection
	if lineNum != state.CursorLine {
		if state.IsInVisualSelection(lineNum, 0) {
			// Entire line is in selection
			return VisualSelectionStyle.Render(line)
		}
		return line
	}

	// This is the cursor line
	col := state.CursorCol
	if col >= len(runes) {
		col = len(runes) - 1
	}
	if col < 0 {
		col = 0
	}

	// Build line with cursor
	var result strings.Builder

	// Handle visual selection on cursor line
	if state.VisualStart != nil && state.VisualEnd != nil {
		selStart := state.VisualStart.Col
		selEnd := state.VisualEnd.Col
		if selStart > selEnd {
			selStart, selEnd = selEnd, selStart
		}

		for i, r := range runes {
			inSelection := i >= selStart && i <= selEnd
			isCursor := i == col

			char := string(r)
			if isCursor {
				if state.Mode == vim.ModeVisual || state.Mode == vim.ModeVisualLine {
					result.WriteString(VisualCursorStyle.Render(char))
				} else {
					result.WriteString(NormalCursorStyle.Render(char))
				}
			} else if inSelection {
				result.WriteString(VisualSelectionStyle.Render(char))
			} else {
				result.WriteString(char)
			}
		}
		return result.String()
	}

	// Normal cursor rendering
	before := string(runes[:col])
	cursorChar := string(runes[col])
	after := ""
	if col+1 < len(runes) {
		after = string(runes[col+1:])
	}

	result.WriteString(before)

	switch state.Mode {
	case vim.ModeInsert:
		result.WriteString(InsertCursorStyle.Render("|"))
		result.WriteString(cursorChar)
	case vim.ModeVisual, vim.ModeVisualLine:
		result.WriteString(VisualCursorStyle.Render(cursorChar))
	default:
		result.WriteString(NormalCursorStyle.Render(cursorChar))
	}

	result.WriteString(after)
	return result.String()
}

// renderModeLine renders the vim mode line
func renderModeLine(e *vim.Editor, c *engine.Challenge) string {
	state := e.GetRenderState()
	var parts []string

	// Mode indicator
	modeStyle := lipgloss.NewStyle().Bold(true)
	switch state.Mode {
	case vim.ModeInsert:
		modeStyle = modeStyle.Foreground(lipgloss.Color("#22c55e"))
	case vim.ModeVisual, vim.ModeVisualLine:
		modeStyle = modeStyle.Foreground(lipgloss.Color("#8b5cf6"))
	default:
		modeStyle = modeStyle.Foreground(lipgloss.Color("#60a5fa"))
	}
	parts = append(parts, modeStyle.Render(fmt.Sprintf("-- %s --", state.ModeString)))

	// Pending command
	if state.Count != "" || state.PendingCmd != "" {
		parts = append(parts, HelpStyle.Render(state.Count+state.PendingCmd))
	}

	// Keystroke count
	parts = append(parts, HelpStyle.Render(fmt.Sprintf("Keys: %d", e.KeystrokeCount)))

	// Reward and par
	parts = append(parts, GoldStyle.Render(fmt.Sprintf("Reward: %dg | Par: %d", c.GoldBase, c.ParKeystrokes)))

	return strings.Join(parts, "  ")
}

func renderHelp(m *Model) string {
	if m.Game.State == engine.StateChallengeActive {
		return HelpStyle.Render("[Ctrl+S] Submit  [Esc] Cancel  |  Use vim commands to edit")
	}
	help := "[hjkl/arrows] Move  [space] Place tower  [c] Challenge  [p] Pause  [q] Quit"
	return HelpStyle.Render(help)
}

// RenderGameOver renders the game over screen
func RenderGameOver(m *Model) string {
	var b strings.Builder

	b.WriteString("\n\n")
	b.WriteString(GameOverStyle.Render("  ██████╗  █████╗ ███╗   ███╗███████╗     ██████╗ ██╗   ██╗███████╗██████╗  \n"))
	b.WriteString(GameOverStyle.Render(" ██╔════╝ ██╔══██╗████╗ ████║██╔════╝    ██╔═══██╗██║   ██║██╔════╝██╔══██╗ \n"))
	b.WriteString(GameOverStyle.Render(" ██║  ███╗███████║██╔████╔██║█████╗      ██║   ██║██║   ██║█████╗  ██████╔╝ \n"))
	b.WriteString(GameOverStyle.Render(" ██║   ██║██╔══██║██║╚██╔╝██║██╔══╝      ██║   ██║╚██╗ ██╔╝██╔══╝  ██╔══██╗ \n"))
	b.WriteString(GameOverStyle.Render(" ╚██████╔╝██║  ██║██║ ╚═╝ ██║███████╗    ╚██████╔╝ ╚████╔╝ ███████╗██║  ██║ \n"))
	b.WriteString(GameOverStyle.Render("  ╚═════╝ ╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝     ╚═════╝   ╚═══╝  ╚══════╝╚═╝  ╚═╝ \n"))
	b.WriteString("\n")

	b.WriteString(fmt.Sprintf("  Wave reached: %d/%d\n", m.Game.Wave, m.Game.TotalWaves))
	b.WriteString(fmt.Sprintf("  Towers built: %d\n", len(m.Game.Towers)))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("  Press [r] to restart or [q] to quit\n"))

	return b.String()
}

// RenderVictory renders the victory screen
func RenderVictory(m *Model) string {
	var b strings.Builder

	b.WriteString("\n\n")
	b.WriteString(VictoryStyle.Render(" ██╗   ██╗██╗ ██████╗████████╗ ██████╗ ██████╗ ██╗   ██╗██╗\n"))
	b.WriteString(VictoryStyle.Render(" ██║   ██║██║██╔════╝╚══██╔══╝██╔═══██╗██╔══██╗╚██╗ ██╔╝██║\n"))
	b.WriteString(VictoryStyle.Render(" ██║   ██║██║██║        ██║   ██║   ██║██████╔╝ ╚████╔╝ ██║\n"))
	b.WriteString(VictoryStyle.Render(" ╚██╗ ██╔╝██║██║        ██║   ██║   ██║██╔══██╗  ╚██╔╝  ╚═╝\n"))
	b.WriteString(VictoryStyle.Render("  ╚████╔╝ ██║╚██████╗   ██║   ╚██████╔╝██║  ██║   ██║   ██╗\n"))
	b.WriteString(VictoryStyle.Render("   ╚═══╝  ╚═╝ ╚═════╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝   ╚═╝   ╚═╝\n"))
	b.WriteString("\n")

	b.WriteString(VictoryStyle.Render("  Congratulations! You defended against all waves!\n\n"))
	b.WriteString(fmt.Sprintf("  Final gold: %d\n", m.Game.Gold))
	b.WriteString(fmt.Sprintf("  Final health: %d/%d\n", m.Game.Health, m.Game.MaxHealth))
	b.WriteString(fmt.Sprintf("  Towers built: %d\n", len(m.Game.Towers)))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("  Press [r] to play again or [q] to quit\n"))

	return b.String()
}
