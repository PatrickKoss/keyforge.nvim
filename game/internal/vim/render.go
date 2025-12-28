package vim

// RenderState contains all info needed to render the editor
// This keeps the vim package decoupled from UI frameworks
type RenderState struct {
	Lines       []string
	CursorLine  int
	CursorCol   int
	Mode        Mode
	ModeString  string
	PendingCmd  string    // e.g., "d" when waiting for motion
	Count       string    // e.g., "23" when count is being entered
	VisualStart *Position // nil if not in visual mode
	VisualEnd   *Position
	StatusMsg   string
}

// GetRenderState returns the current state for rendering
func (e *Editor) GetRenderState() RenderState {
	state := RenderState{
		Lines:      make([]string, e.Buffer.LineCount()),
		CursorLine: e.Cursor.Line,
		CursorCol:  e.Cursor.Col,
		Mode:       e.Mode,
		ModeString: e.Mode.String(),
		StatusMsg:  e.StatusMessage,
	}

	for i := 0; i < e.Buffer.LineCount(); i++ {
		state.Lines[i] = e.Buffer.GetLine(i)
	}

	// Pending command display
	if e.PendingOp != OpNone {
		state.PendingCmd = e.PendingOp.String()
	}

	// Count display
	if e.Count > 0 {
		state.Count = string(rune('0' + e.Count%10))
		c := e.Count / 10
		for c > 0 {
			state.Count = string(rune('0'+c%10)) + state.Count
			c = c / 10
		}
	}

	// Visual mode positions
	if e.Mode == ModeVisual || e.Mode == ModeVisualLine {
		start := e.VisualStart
		end := e.Cursor
		state.VisualStart = &start
		state.VisualEnd = &end
	}

	return state
}

// IsInVisualSelection checks if a position is within the visual selection
func (s *RenderState) IsInVisualSelection(line, col int) bool {
	if s.VisualStart == nil || s.VisualEnd == nil {
		return false
	}

	start := *s.VisualStart
	end := *s.VisualEnd

	// Normalize so start <= end
	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		start, end = end, start
	}

	if s.Mode == ModeVisualLine {
		return line >= start.Line && line <= end.Line
	}

	// Character-wise
	if line < start.Line || line > end.Line {
		return false
	}
	if line == start.Line && line == end.Line {
		return col >= start.Col && col <= end.Col
	}
	if line == start.Line {
		return col >= start.Col
	}
	if line == end.Line {
		return col <= end.Col
	}
	return true
}
