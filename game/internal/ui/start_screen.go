package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
)

// Start screen styles.
var (
	LogoStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	MenuTitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true).
			Padding(0, 1)

	MenuItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	MenuItemSelectedStyle = lipgloss.NewStyle().
				Background(ColorPrimary).
				Foreground(lipgloss.Color("#000000")).
				Padding(0, 2)

	MenuItemDisabledStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Padding(0, 2)

	PreviewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	SettingLabelStyle = lipgloss.NewStyle().
				Foreground(ColorMuted).
				Width(20)

	SettingValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#ffffff")).
				Bold(true)

	SettingSelectedStyle = lipgloss.NewStyle().
				Background(ColorSecondary).
				Foreground(lipgloss.Color("#000000")).
				Padding(0, 1)

	SliderTrackStyle = lipgloss.NewStyle().
				Foreground(ColorMuted)

	SliderFillStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary)
)

// RenderStartScreen renders the level selection screen.
func RenderStartScreen(m *Model) string {
	var b strings.Builder

	// Logo
	b.WriteString(renderLogo())
	b.WriteString("\n\n")

	// Create two-column layout: level list on left, preview on right
	leftColumn := renderLevelList(m)
	rightColumn := renderLevelPreview(m)

	// Join columns side by side
	leftLines := strings.Split(leftColumn, "\n")
	rightLines := strings.Split(rightColumn, "\n")

	// Pad to same height
	maxLines := len(leftLines)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}
	for len(leftLines) < maxLines {
		leftLines = append(leftLines, strings.Repeat(" ", 30))
	}
	for len(rightLines) < maxLines {
		rightLines = append(rightLines, "")
	}

	for i := range maxLines {
		left := leftLines[i]
		right := rightLines[i]
		// Pad left column to fixed width
		leftWidth := lipgloss.Width(left)
		if leftWidth < 35 {
			left += strings.Repeat(" ", 35-leftWidth)
		}
		b.WriteString(left)
		b.WriteString("  ")
		b.WriteString(right)
		b.WriteString("\n")
	}

	// Help text
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render("[j/k] Select level  [Enter] Configure settings  [q] Quit"))

	return b.String()
}

// RenderSettingsScreen renders the game settings configuration screen.
func RenderSettingsScreen(m *Model) string {
	var b strings.Builder

	// Title
	b.WriteString(MenuTitleStyle.Render("Game Settings"))
	b.WriteString("\n\n")

	// Level info
	if m.SelectedLevel != nil {
		levelInfo := "Level: " + m.SelectedLevel.Name
		b.WriteString(HelpStyle.Render(levelInfo))
		b.WriteString("\n\n")
	}

	// Settings
	settings := []struct {
		label   string
		options []string
		current int
	}{
		{"Difficulty", []string{"Easy", "Normal", "Hard"}, difficultyIndex(m.Settings.Difficulty)},
		{"Game Speed", []string{"0.5x", "1x", "1.5x", "2x"}, speedIndex(m.Settings.GameSpeed)},
	}

	for i, setting := range settings {
		isSelected := m.SettingsMenuIndex == i

		// Label
		label := SettingLabelStyle.Render(setting.label + ":")

		// Options
		var opts []string
		for j, opt := range setting.options {
			if j == setting.current {
				if isSelected {
					opts = append(opts, SettingSelectedStyle.Render(opt))
				} else {
					opts = append(opts, SettingValueStyle.Render("["+opt+"]"))
				}
			} else {
				opts = append(opts, HelpStyle.Render(opt))
			}
		}

		b.WriteString(label)
		b.WriteString(strings.Join(opts, "  "))
		b.WriteString("\n")
	}

	// Sliders for gold and health
	b.WriteString("\n")

	// Starting Gold slider
	goldSelected := m.SettingsMenuIndex == 2
	goldLabel := SettingLabelStyle.Render("Starting Gold:")
	goldValue := strconv.Itoa(m.Settings.StartingGold)
	if goldSelected {
		goldValue = SettingSelectedStyle.Render(goldValue)
	} else {
		goldValue = SettingValueStyle.Render(goldValue)
	}
	goldSlider := renderSlider(m.Settings.StartingGold, 100, 500)
	b.WriteString(goldLabel + goldValue + "  " + goldSlider + "\n")

	// Starting Health slider
	healthSelected := m.SettingsMenuIndex == 3
	healthLabel := SettingLabelStyle.Render("Starting Health:")
	healthValue := strconv.Itoa(m.Settings.StartingHealth)
	if healthSelected {
		healthValue = SettingSelectedStyle.Render(healthValue)
	} else {
		healthValue = SettingValueStyle.Render(healthValue)
	}
	healthSlider := renderSlider(m.Settings.StartingHealth, 50, 200)
	b.WriteString(healthLabel + healthValue + "  " + healthSlider + "\n")

	// Start Game button
	b.WriteString("\n")
	startSelected := m.SettingsMenuIndex == 4
	if startSelected {
		b.WriteString(MenuItemSelectedStyle.Render("[ Start Game ]"))
	} else {
		b.WriteString(MenuItemStyle.Render("[ Start Game ]"))
	}

	// Help text
	b.WriteString("\n\n")
	b.WriteString(HelpStyle.Render("[j/k] Navigate  [h/l] Adjust  [Enter] Start  [Esc] Back"))

	return b.String()
}

func renderLogo() string {
	logo := `
 ██╗  ██╗███████╗██╗   ██╗███████╗ ██████╗ ██████╗  ██████╗ ███████╗
 ██║ ██╔╝██╔════╝╚██╗ ██╔╝██╔════╝██╔═══██╗██╔══██╗██╔════╝ ██╔════╝
 █████╔╝ █████╗   ╚████╔╝ █████╗  ██║   ██║██████╔╝██║  ███╗█████╗
 ██╔═██╗ ██╔══╝    ╚██╔╝  ██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝
 ██║  ██╗███████╗   ██║   ██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗
 ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝`
	return LogoStyle.Render(logo)
}

func renderLevelList(m *Model) string {
	var b strings.Builder

	b.WriteString(MenuTitleStyle.Render("Select Level"))
	b.WriteString("\n\n")

	levels := m.LevelRegistry.GetAll()
	for i := range levels {
		level := &levels[i]
		isSelected := i == m.LevelMenuIndex

		// Level name and difficulty
		diffIcon := difficultyIcon(level.Difficulty)
		text := fmt.Sprintf("%s %s", diffIcon, level.Name)

		if isSelected {
			b.WriteString(MenuItemSelectedStyle.Render(text))
		} else {
			b.WriteString(MenuItemStyle.Render(text))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderLevelPreview(m *Model) string {
	levels := m.LevelRegistry.GetAll()
	if m.LevelMenuIndex >= len(levels) {
		return ""
	}

	level := levels[m.LevelMenuIndex]

	var b strings.Builder

	// Level name and description
	b.WriteString(MenuTitleStyle.Render(level.Name))
	b.WriteString("\n")
	b.WriteString(HelpStyle.Render(level.Description))
	b.WriteString("\n\n")

	// Mini grid preview
	b.WriteString(renderMiniGrid(&level))
	b.WriteString("\n\n")

	// Stats
	b.WriteString(fmt.Sprintf("Waves: %d  Difficulty: %s\n",
		level.TotalWaves,
		string(level.Difficulty)))
	b.WriteString("\n")

	// Enemies
	b.WriteString(HelpStyle.Render("Enemies: "))
	var enemies []string
	for _, et := range level.EnemyTypes {
		info := entities.EnemyTypes[et]
		enemies = append(enemies, info.Symbol+" "+info.Name)
	}
	b.WriteString(strings.Join(enemies, ", "))
	b.WriteString("\n")

	// Towers
	b.WriteString(HelpStyle.Render("Towers:  "))
	var towers []string
	for _, tt := range level.AllowedTowers {
		info := entities.TowerTypes[tt]
		towers = append(towers, info.Symbol+" "+info.Name)
	}
	b.WriteString(strings.Join(towers, ", "))
	b.WriteString("\n")

	return PreviewBoxStyle.Render(b.String())
}

func renderMiniGrid(level *engine.Level) string {
	// Create a smaller preview grid (scale down)
	previewWidth := 20
	previewHeight := 7

	scaleX := float64(level.GridWidth) / float64(previewWidth)
	scaleY := float64(level.GridHeight) / float64(previewHeight)

	// Initialize grid
	grid := make([][]rune, previewHeight)
	for y := range previewHeight {
		grid[y] = make([]rune, previewWidth)
		for x := range previewWidth {
			grid[y][x] = '·'
		}
	}

	// Plot path (scaled)
	for _, pos := range level.Path {
		x := int(pos.X / scaleX)
		y := int(pos.Y / scaleY)
		if x >= 0 && x < previewWidth && y >= 0 && y < previewHeight {
			grid[y][x] = '░'
		}
	}

	// Build string
	var b strings.Builder
	for y := range previewHeight {
		for x := range previewWidth {
			if grid[y][x] == '░' {
				b.WriteString(PathCellStyle.Render(string(grid[y][x])))
			} else {
				b.WriteString(EmptyCellStyle.Render(string(grid[y][x])))
			}
		}
		if y < previewHeight-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

const sliderWidth = 20

func renderSlider(value, minVal, maxVal int) string {
	ratio := float64(value-minVal) / float64(maxVal-minVal)
	filled := int(ratio * float64(sliderWidth))

	var b strings.Builder
	b.WriteString("[")
	for i := range sliderWidth {
		if i < filled {
			b.WriteString(SliderFillStyle.Render("█"))
		} else {
			b.WriteString(SliderTrackStyle.Render("░"))
		}
	}
	b.WriteString("]")
	return b.String()
}

func difficultyIcon(d engine.LevelDifficulty) string {
	switch d {
	case engine.LevelDifficultyBeginner:
		return "★☆☆"
	case engine.LevelDifficultyIntermediate:
		return "★★☆"
	case engine.LevelDifficultyAdvanced:
		return "★★★"
	default:
		return "★☆☆"
	}
}

func difficultyIndex(d string) int {
	switch d {
	case engine.DifficultyEasy:
		return 0
	case engine.DifficultyNormal:
		return 1
	case engine.DifficultyHard:
		return 2
	default:
		return 1
	}
}

func speedIndex(s engine.GameSpeed) int {
	switch s {
	case engine.SpeedHalf:
		return 0
	case engine.SpeedNormal:
		return 1
	case engine.SpeedFast:
		return 2
	case engine.SpeedDouble:
		return 3
	default:
		return 1
	}
}
