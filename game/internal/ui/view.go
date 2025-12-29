package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
	"github.com/keyforge/keyforge/internal/vim"
)

// RenderGame renders the complete game view.
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

func renderTitle(_ *Model) string {
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
	case engine.StateMenu, engine.StateLevelSelect, engine.StateSettings, engine.StatePlaying, engine.StateWaveComplete:
		// No special status display for these states
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
	for range g.Width * 2 {
		b.WriteString(BoxHorizontal)
	}
	b.WriteString(BoxTopRight)
	b.WriteString("\n")

	// Create and populate grid
	grid := initEmptyGrid(g.Width, g.Height)
	populateGridEntities(grid, g)

	// Render range overlay before cursor (so cursor is on top)
	// Show range when hovering over existing tower, or when placing new tower
	tower := g.GetTowerAt(g.CursorX, g.CursorY)
	if tower != nil {
		RenderRangeForHover(grid, g)
	} else if g.CanPlaceTower(g.CursorX, g.CursorY) {
		RenderRangeForPlacement(grid, g)
	}

	renderGridCursor(grid, g)

	// Build grid rows
	for y := range g.Height {
		b.WriteString(BoxVertical)
		for x := range g.Width {
			cell := grid[y][x]
			b.WriteString(cell)
			// Pad to fixed width (2 chars per cell)
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
	for range g.Width * 2 {
		b.WriteString(BoxHorizontal)
	}
	b.WriteString(BoxBottomRight)

	return b.String()
}

// initEmptyGrid creates a grid filled with empty cells.
func initEmptyGrid(width, height int) [][]string {
	grid := make([][]string, height)
	for y := range height {
		grid[y] = make([]string, width)
		for x := range width {
			grid[y][x] = EmptyCellStyle.Render(EmptyCell)
		}
	}
	return grid
}

// populateGridEntities renders path, towers, enemies, projectiles, and effects onto the grid.
func populateGridEntities(grid [][]string, g *engine.Game) {
	width := g.Width
	height := g.Height

	// Render path
	for _, p := range g.Path {
		x, y := p.IntPos()
		if isInBounds(x, y, width, height) {
			grid[y][x] = PathCellStyle.Render(PathChar)
		}
	}

	// Render towers
	for _, tower := range g.Towers {
		x, y := tower.Pos.IntPos()
		if isInBounds(x, y, width, height) {
			grid[y][x] = renderTower(tower)
		}
	}

	// Render enemies
	for _, enemy := range g.Enemies {
		if enemy.Dead {
			continue
		}
		x, y := enemy.Pos.IntPos()
		if isInBounds(x, y, width, height) {
			grid[y][x] = renderEnemy(enemy)
		}
	}

	// Render projectiles
	for _, proj := range g.Projectiles {
		if proj.Done {
			continue
		}
		x, y := proj.Pos.IntPos()
		if isInBounds(x, y, width, height) {
			grid[y][x] = ProjectileStyle.Render(ProjectileChar)
		}
	}

	// Render effects (on top of everything except cursor)
	for _, effect := range g.Effects.Effects {
		if effect.Done {
			continue
		}
		x, y := effect.Pos.IntPos()
		if isInBounds(x, y, width, height) {
			frame := effect.CurrentFrame()
			if frame != "" {
				info := effect.Info()
				style := lipgloss.NewStyle().Foreground(lipgloss.Color(info.Color))
				grid[y][x] = style.Render(frame)
			}
		}
	}
}

// renderGridCursor renders the cursor onto the grid if applicable.
func renderGridCursor(grid [][]string, g *engine.Game) {
	if g.State != engine.StatePlaying && g.State != engine.StatePaused {
		return
	}
	if !isInBounds(g.CursorX, g.CursorY, g.Width, g.Height) {
		return
	}

	if g.CanPlaceTower(g.CursorX, g.CursorY) {
		grid[g.CursorY][g.CursorX] = CursorStyle.Render("▪")
	} else {
		grid[g.CursorY][g.CursorX] = lipgloss.NewStyle().
			Background(ColorDanger).
			Foreground(lipgloss.Color("#ffffff")).
			Render("✗")
	}
}

// isInBounds checks if coordinates are within grid bounds.
func isInBounds(x, y, width, height int) bool {
	return x >= 0 && x < width && y >= 0 && y < height
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

	// TODO: Add level indicator when upgrade system is implemented
	_ = tower.Level // Reserved for future use

	return style.Render(char)
}

func renderEnemy(enemy *entities.Enemy) string {
	var style lipgloss.Style
	var char string

	switch enemy.Type {
	case entities.EnemyMite:
		style = EnemyMiteStyle
		char = EnemyMiteChar
	case entities.EnemyBug:
		style = EnemyBugStyle
		char = EnemyBugChar
	case entities.EnemyGremlin:
		style = EnemyGremlinStyle
		char = EnemyGremlinChar
	case entities.EnemyCrawler:
		style = EnemyCrawlerStyle
		char = EnemyCrawlerChar
	case entities.EnemySpecter:
		style = EnemySpecterStyle
		char = EnemySpecterChar
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

// renderVimBuffer renders the vim buffer with cursor.
func renderVimBuffer(e *vim.Editor) string {
	state := e.GetRenderState()
	var b strings.Builder

	bufferBg := lipgloss.NewStyle().
		Background(lipgloss.Color("#1a1a1a")).
		Padding(0, 1)

	var lines []string
	for lineNum, line := range state.Lines {
		renderedLine := renderVimLine(line, lineNum, &state)
		lines = append(lines, renderedLine)
	}

	content := strings.Join(lines, "\n")
	b.WriteString(bufferBg.Render(content))

	return b.String()
}

// renderVimLine renders a single line with cursor highlighting.
func renderVimLine(line string, lineNum int, state *vim.RenderState) string {
	runes := []rune(line)

	// Handle empty line
	if len(runes) == 0 {
		return renderEmptyLine(lineNum, state)
	}

	// Not the cursor line - check for visual selection
	if lineNum != state.CursorLine {
		return renderNonCursorLine(line, lineNum, state)
	}

	// This is the cursor line
	col := clampCursorCol(state.CursorCol, len(runes))

	// Handle visual selection on cursor line
	if state.VisualStart != nil && state.VisualEnd != nil {
		return renderVisualSelectionLine(runes, col, state)
	}

	// Normal cursor rendering
	return renderNormalCursorLine(runes, col, state.Mode)
}

// renderEmptyLine renders an empty line with cursor if applicable.
func renderEmptyLine(lineNum int, state *vim.RenderState) string {
	if lineNum == state.CursorLine {
		if state.Mode == vim.ModeInsert {
			return InsertCursorStyle.Render("|")
		}
		return NormalCursorStyle.Render(" ")
	}
	return " " // Empty line placeholder
}

// renderNonCursorLine renders a line that doesn't contain the cursor.
func renderNonCursorLine(line string, lineNum int, state *vim.RenderState) string {
	if state.IsInVisualSelection(lineNum, 0) {
		return VisualSelectionStyle.Render(line)
	}
	return line
}

// clampCursorCol ensures cursor column is within valid bounds.
func clampCursorCol(col, lineLen int) int {
	if col >= lineLen {
		col = lineLen - 1
	}
	if col < 0 {
		col = 0
	}
	return col
}

// renderVisualSelectionLine renders a line with visual selection highlighting.
func renderVisualSelectionLine(runes []rune, cursorCol int, state *vim.RenderState) string {
	selStart := state.VisualStart.Col
	selEnd := state.VisualEnd.Col
	if selStart > selEnd {
		selStart, selEnd = selEnd, selStart
	}

	var result strings.Builder
	for i, r := range runes {
		char := string(r)
		result.WriteString(renderVisualChar(char, i, cursorCol, selStart, selEnd, state.Mode))
	}
	return result.String()
}

// renderVisualChar renders a single character in visual selection mode.
func renderVisualChar(char string, idx, cursorCol, selStart, selEnd int, mode vim.Mode) string {
	inSelection := idx >= selStart && idx <= selEnd
	isCursor := idx == cursorCol

	if isCursor {
		if mode == vim.ModeVisual || mode == vim.ModeVisualLine {
			return VisualCursorStyle.Render(char)
		}
		return NormalCursorStyle.Render(char)
	}
	if inSelection {
		return VisualSelectionStyle.Render(char)
	}
	return char
}

// renderNormalCursorLine renders a line with normal cursor highlighting.
func renderNormalCursorLine(runes []rune, col int, mode vim.Mode) string {
	var result strings.Builder

	before := string(runes[:col])
	cursorChar := string(runes[col])
	after := ""
	if col+1 < len(runes) {
		after = string(runes[col+1:])
	}

	result.WriteString(before)
	result.WriteString(renderCursorChar(cursorChar, mode))
	result.WriteString(after)

	return result.String()
}

// renderCursorChar renders the character under the cursor based on mode.
func renderCursorChar(char string, mode vim.Mode) string {
	switch mode {
	case vim.ModeInsert:
		return InsertCursorStyle.Render("|") + char
	case vim.ModeVisual, vim.ModeVisualLine:
		return VisualCursorStyle.Render(char)
	default:
		return NormalCursorStyle.Render(char)
	}
}

// renderModeLine renders the vim mode line.
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

	// Keystroke count and reward/par
	parts = append(parts,
		HelpStyle.Render(fmt.Sprintf("Keys: %d", e.KeystrokeCount)),
		GoldStyle.Render(fmt.Sprintf("Reward: %dg | Par: %d", c.GoldBase, c.ParKeystrokes)),
	)

	return strings.Join(parts, "  ")
}

func renderHelp(m *Model) string {
	if m.Game.State == engine.StateChallengeActive {
		return HelpStyle.Render("[Ctrl+S] Submit  [Esc] Cancel  |  Use vim commands to edit")
	}
	help := "[hjkl/arrows] Move  [space] Place tower  [c] Challenge  [p] Pause  [q] Quit"
	return HelpStyle.Render(help)
}

// RenderGameOver renders the game over screen.
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
	b.WriteString(HelpStyle.Render("  Press [r] to restart, [m] for menu, or [q] to quit\n"))

	return b.String()
}

// RenderVictory renders the victory screen.
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
	b.WriteString(HelpStyle.Render("  Press [r] to play again, [m] for menu, or [q] to quit\n"))

	return b.String()
}
