package ui

import (
	"math"

	"github.com/charmbracelet/lipgloss"

	"github.com/keyforge/keyforge/internal/engine"
	"github.com/keyforge/keyforge/internal/entities"
)

// RangeOverlayStyle is the style for range indicator cells.
var RangeOverlayStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#60a5fa")).
	Faint(true)

// RangeOverlayChar is the character used to indicate cells within range.
const RangeOverlayChar = "·"

// CellsInRange returns all grid cells within the given range from a center point.
// Uses squared distance for efficiency (no sqrt needed).
func CellsInRange(centerX, centerY int, rangeRadius float64, width, height int) []struct{ X, Y int } {
	var cells []struct{ X, Y int }

	// Calculate the bounding box to check
	rangeInt := int(math.Ceil(rangeRadius))

	for dy := -rangeInt; dy <= rangeInt; dy++ {
		for dx := -rangeInt; dx <= rangeInt; dx++ {
			x := centerX + dx
			y := centerY + dy

			// Skip if out of bounds
			if x < 0 || x >= width || y < 0 || y >= height {
				continue
			}

			// Check if cell center is within range (using squared distance)
			distSquared := float64(dx*dx + dy*dy)
			rangeSquared := rangeRadius * rangeRadius

			if distSquared <= rangeSquared {
				cells = append(cells, struct{ X, Y int }{x, y})
			}
		}
	}

	return cells
}

// renderRangeOverlay renders the range indicator onto the grid.
// It only overlays empty cells and path cells, not entities.
func renderRangeOverlay(grid [][]string, g *engine.Game, centerX, centerY int, rangeRadius float64) {
	cells := CellsInRange(centerX, centerY, rangeRadius, g.Width, g.Height)

	for _, cell := range cells {
		// Skip the center cell (where cursor/tower is)
		if cell.X == centerX && cell.Y == centerY {
			continue
		}

		// Only overlay empty cells or path cells
		currentCell := grid[cell.Y][cell.X]
		if currentCell == EmptyCellStyle.Render(EmptyCell) || currentCell == PathCellStyle.Render(PathChar) {
			grid[cell.Y][cell.X] = RangeOverlayStyle.Render(RangeOverlayChar)
		}
	}
}

// RenderRangeForPlacement renders the range overlay when placing a tower.
func RenderRangeForPlacement(grid [][]string, g *engine.Game) {
	if g.State != engine.StatePlaying {
		return
	}

	// Get the selected tower's range
	info := entities.TowerTypes[g.SelectedTower]
	renderRangeOverlay(grid, g, g.CursorX, g.CursorY, info.Range)
}

// RenderRangeForHover renders the range overlay when hovering over an existing tower.
func RenderRangeForHover(grid [][]string, g *engine.Game) {
	if g.State != engine.StatePlaying && g.State != engine.StatePaused {
		return
	}

	// Check if cursor is on an existing tower
	tower := g.GetTowerAt(g.CursorX, g.CursorY)
	if tower == nil {
		return
	}

	// Use the tower's current range (including upgrades)
	renderRangeOverlay(grid, g, g.CursorX, g.CursorY, tower.Range)
}
