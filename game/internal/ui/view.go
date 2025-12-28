package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
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

	// Shop
	b.WriteString(renderShop(m))
	b.WriteString("\n")

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
		status = WaveStyle.Render("  [CHALLENGE]")
	}

	hud := fmt.Sprintf("%s    %s    %s%s", waveInfo, goldInfo, healthInfo, status)
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

func renderHelp(m *Model) string {
	help := "[hjkl/arrows] Move  [space] Place tower  [p] Pause  [q] Quit"
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
