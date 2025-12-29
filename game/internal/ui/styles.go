package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette.
var (
	ColorPrimary    = lipgloss.Color("#22c55e") // green
	ColorSecondary  = lipgloss.Color("#8b5cf6") // purple
	ColorWarning    = lipgloss.Color("#f59e0b") // amber
	ColorDanger     = lipgloss.Color("#ef4444") // red
	ColorSuccess    = lipgloss.Color("#22c55e") // green (same as primary)
	ColorMuted      = lipgloss.Color("#6b7280") // gray
	ColorBackground = lipgloss.Color("#1f2937") // dark gray
	ColorBorder     = lipgloss.Color("#374151") // medium gray
	ColorPath       = lipgloss.Color("#4b5563") // path color
	ColorCursor     = lipgloss.Color("#fbbf24") // yellow for cursor
	ColorGold       = lipgloss.Color("#fbbf24") // gold
	ColorHealth     = lipgloss.Color("#ef4444") // health red
)

// Styles.
var (
	// Box styles.
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

		// HUD styles.
	HUDStyle = lipgloss.NewStyle().
			Padding(0, 1)

	GoldStyle = lipgloss.NewStyle().
			Foreground(ColorGold).
			Bold(true)

	HealthStyle = lipgloss.NewStyle().
			Foreground(ColorHealth).
			Bold(true)

	WaveStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

		// Cell styles.
	EmptyCellStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	PathCellStyle = lipgloss.NewStyle().
			Foreground(ColorPath)

	CursorStyle = lipgloss.NewStyle().
			Background(ColorCursor).
			Foreground(lipgloss.Color("#000000"))

		// Tower styles.
	TowerArrowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#22c55e"))

	TowerLSPStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8b5cf6"))

	TowerRefactorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f59e0b"))

		// Enemy styles.
	EnemyMiteStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a3e635")) // Lime green

	EnemyBugStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444"))

	EnemyGremlinStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#f97316"))

	EnemyCrawlerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#78716c")) // Stone gray

	EnemySpecterStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#c4b5fd")) // Light purple

	EnemyDaemonStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#dc2626"))

	EnemyBossStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7c2d12")).
			Bold(true)

		// Projectile style.
	ProjectileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#fef08a"))

		// Status styles.
	PausedStyle = lipgloss.NewStyle().
			Foreground(ColorWarning).
			Bold(true)

	GameOverStyle = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	VictoryStyle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true)

	ChallengeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60a5fa")). // blue
			Bold(true)

		// Shop styles.
	ShopItemStyle = lipgloss.NewStyle().
			Padding(0, 1)

	ShopItemDisabledStyle = lipgloss.NewStyle().
				Padding(0, 1).
				Foreground(ColorMuted)

	SelectedStyle = lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(lipgloss.Color("#000000"))

		// Help text.
	HelpStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			Italic(true)

		// Vim editor cursor styles.
	NormalCursorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#fbbf24")).
				Foreground(lipgloss.Color("#000000"))

	InsertCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#22c55e")).
				Bold(true)

	VisualCursorStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#8b5cf6")).
				Foreground(lipgloss.Color("#ffffff"))

	VisualSelectionStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#4c1d95"))
)

// Characters for rendering.
const (
	// Box drawing.
	BoxTopLeft     = "╔"
	BoxTopRight    = "╗"
	BoxBottomLeft  = "╚"
	BoxBottomRight = "╝"
	BoxHorizontal  = "═"
	BoxVertical    = "║"

	// Path characters.
	PathHorizontal = "─"
	PathVertical   = "│"
	PathCornerTL   = "┌"
	PathCornerTR   = "┐"
	PathCornerBL   = "└"
	PathCornerBR   = "┘"

	// Entity characters (with emoji fallbacks).
	TowerArrowChar    = "🏹"
	TowerLSPChar      = "🔮"
	TowerRefactorChar = "⚡"

	EnemyMiteChar    = "🦟"
	EnemyBugChar     = "🐛"
	EnemyGremlinChar = "👹"
	EnemyCrawlerChar = "🐌"
	EnemySpecterChar = "👻"
	EnemyDaemonChar  = "👿"
	EnemyBossChar    = "💀"

	ProjectileChar = "•"

	// Fallback ASCII.
	TowerCharASCII      = "T"
	EnemyCharASCII      = "E"
	ProjectileCharASCII = "*"

	// UI characters.
	EmptyCell  = "·"
	CursorChar = "█"
	PathChar   = "░"

	// Health bar.
	HealthFull  = "█"
	HealthHalf  = "▓"
	HealthLow   = "▒"
	HealthEmpty = "░"
)

// RenderHealthBar creates a visual health bar.
func RenderHealthBar(current, maxHealth int, width int) string {
	if maxHealth <= 0 {
		return ""
	}
	ratio := float64(current) / float64(maxHealth)
	filled := int(ratio * float64(width))

	bar := ""
	var barSb194 strings.Builder
	for i := range width {
		if i < filled {
			if ratio > 0.6 {
				barSb194.WriteString(HealthStyle.Render(HealthFull))
			} else if ratio > 0.3 {
				bar += lipgloss.NewStyle().Foreground(ColorWarning).Render(HealthHalf)
			} else {
				bar += lipgloss.NewStyle().Foreground(ColorDanger).Render(HealthLow)
			}
		} else {
			barSb194.WriteString(EmptyCellStyle.Render(HealthEmpty))
		}
	}
	bar += barSb194.String()
	return bar
}
