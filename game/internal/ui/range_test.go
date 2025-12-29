package ui

import (
	"testing"
)

func TestCellsInRange(t *testing.T) {
	tests := []struct {
		name       string
		centerX    int
		centerY    int
		rangeVal   float64
		width      int
		height     int
		wantCount  int
		wantCenter bool // Should center be included?
	}{
		{
			name:       "Range 1.0 from center",
			centerX:    5,
			centerY:    5,
			rangeVal:   1.0,
			width:      10,
			height:     10,
			wantCount:  5, // center + 4 adjacent
			wantCenter: true,
		},
		{
			name:       "Range 2.0 from center",
			centerX:    5,
			centerY:    5,
			rangeVal:   2.0,
			width:      10,
			height:     10,
			wantCount:  13, // Larger circle
			wantCenter: true,
		},
		{
			name:       "Range at edge - clips to bounds",
			centerX:    0,
			centerY:    0,
			rangeVal:   2.0,
			width:      10,
			height:     10,
			wantCount:  6, // Only cells in positive quadrant + center
			wantCenter: true,
		},
		{
			name:       "Range 0 - only center",
			centerX:    5,
			centerY:    5,
			rangeVal:   0.0,
			width:      10,
			height:     10,
			wantCount:  1, // Only center
			wantCenter: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cells := CellsInRange(tc.centerX, tc.centerY, tc.rangeVal, tc.width, tc.height)

			if len(cells) != tc.wantCount {
				t.Errorf("Expected %d cells, got %d", tc.wantCount, len(cells))
			}

			// Check if center is included
			hasCenter := false
			for _, cell := range cells {
				if cell.X == tc.centerX && cell.Y == tc.centerY {
					hasCenter = true
					break
				}
			}
			if hasCenter != tc.wantCenter {
				t.Errorf("Center inclusion: expected %v, got %v", tc.wantCenter, hasCenter)
			}
		})
	}
}

func TestCellsInRangeBounds(t *testing.T) {
	// Verify no cells are returned outside grid bounds
	width := 10
	height := 10
	centerX := 0
	centerY := 0
	rangeVal := 5.0

	cells := CellsInRange(centerX, centerY, rangeVal, width, height)

	for _, cell := range cells {
		if cell.X < 0 || cell.X >= width {
			t.Errorf("Cell X=%d out of bounds [0, %d)", cell.X, width)
		}
		if cell.Y < 0 || cell.Y >= height {
			t.Errorf("Cell Y=%d out of bounds [0, %d)", cell.Y, height)
		}
	}
}

func TestCellsInRangeSymmetry(t *testing.T) {
	// Range should be symmetric around center (when not near edges)
	centerX := 10
	centerY := 10
	rangeVal := 3.0
	width := 20
	height := 20

	cells := CellsInRange(centerX, centerY, rangeVal, width, height)

	// Count cells in each direction from center
	left, right, up, down := 0, 0, 0, 0
	for _, cell := range cells {
		if cell.X < centerX {
			left++
		} else if cell.X > centerX {
			right++
		}
		if cell.Y < centerY {
			up++
		} else if cell.Y > centerY {
			down++
		}
	}

	if left != right {
		t.Errorf("Horizontal asymmetry: left=%d, right=%d", left, right)
	}
	if up != down {
		t.Errorf("Vertical asymmetry: up=%d, down=%d", up, down)
	}
}

func TestCellsInRangeDistanceAccuracy(t *testing.T) {
	// Verify all returned cells are actually within range
	centerX := 5
	centerY := 5
	rangeVal := 2.5
	width := 10
	height := 10

	cells := CellsInRange(centerX, centerY, rangeVal, width, height)

	for _, cell := range cells {
		dx := float64(cell.X - centerX)
		dy := float64(cell.Y - centerY)
		distSquared := dx*dx + dy*dy
		rangeSquared := rangeVal * rangeVal

		if distSquared > rangeSquared {
			t.Errorf("Cell (%d,%d) is outside range: dist^2=%v, range^2=%v",
				cell.X, cell.Y, distSquared, rangeSquared)
		}
	}
}

func TestCellsInRangeTowerRanges(t *testing.T) {
	// Test with actual tower range values from the game
	towerRanges := []float64{2.5, 3.0, 5.0} // Arrow, Refactor, LSP ranges

	for _, rangeVal := range towerRanges {
		cells := CellsInRange(10, 10, rangeVal, 20, 20)
		if len(cells) == 0 {
			t.Errorf("No cells for range %v", rangeVal)
		}
		// Larger range should have more cells
		if rangeVal > 2.5 {
			prevCells := CellsInRange(10, 10, rangeVal-0.5, 20, 20)
			if len(cells) < len(prevCells) {
				t.Errorf("Larger range %v has fewer cells than smaller range", rangeVal)
			}
		}
	}
}
