package vim

// Mode represents vim editing mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeInsert
	ModeVisual
	ModeVisualLine
	ModeOperatorPending
)

// String returns the mode name.
func (m Mode) String() string {
	switch m {
	case ModeNormal:
		return "NORMAL"
	case ModeInsert:
		return "INSERT"
	case ModeVisual:
		return "VISUAL"
	case ModeVisualLine:
		return "V-LINE"
	case ModeOperatorPending:
		return "NORMAL"
	default:
		return "UNKNOWN"
	}
}

// Position represents a cursor position.
type Position struct {
	Line int // 0-based line number
	Col  int // 0-based column (rune index, not byte)
}

// Range represents a text range.
type Range struct {
	Start    Position
	End      Position
	Linewise bool // For dd, yy operations (whole lines)
}

// WaitState represents what the editor is waiting for.
type WaitState int

const (
	WaitNone     WaitState = iota
	WaitMotion             // After operator, waiting for motion
	WaitChar               // After f/F/t/T, waiting for character
	WaitRegister           // After ", waiting for register name
)

// FindState stores the last find command for ; and ,.
type FindState struct {
	Char    rune
	Forward bool // f/t vs F/T
	Till    bool // t/T vs f/F
}

// Snapshot stores buffer state for undo/redo.
type Snapshot struct {
	Buffer *Buffer
	Cursor Position
}

// Editor is the main vim editor state.
type Editor struct {
	Buffer *Buffer
	Cursor Position
	Mode   Mode

	// Command state
	PendingOp  OperatorType // Operator waiting for motion
	Count      int          // Numeric prefix (0 = no count)
	CountStack []int        // For operator + count + motion
	WaitingFor WaitState    // What we're waiting for next

	// Find/till state for f, F, t, T and ; ,
	LastFind FindState

	// Registers
	Unnamed   string          // Default register ""
	Registers map[rune]string // Named registers a-z

	// Undo/Redo
	UndoStack []Snapshot
	RedoStack []Snapshot

	// Visual mode
	VisualStart Position

	// Keystroke tracking for efficiency scoring
	KeystrokeCount int

	// Status message for display
	StatusMessage string
}

// NewEditor creates a new editor with the given initial text.
func NewEditor(text string) *Editor {
	buf := NewBuffer(text)
	return &Editor{
		Buffer:    buf,
		Cursor:    Position{Line: 0, Col: 0},
		Mode:      ModeNormal,
		PendingOp: OpNone,
		Registers: make(map[rune]string),
		UndoStack: make([]Snapshot, 0),
		RedoStack: make([]Snapshot, 0),
	}
}

// SetCursor sets the cursor position, clamping to valid bounds.
func (e *Editor) SetCursor(pos Position) {
	e.Cursor = e.clampPosition(pos)
}

// clampPosition ensures position is within buffer bounds.
func (e *Editor) clampPosition(pos Position) Position {
	// Clamp line
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Line >= e.Buffer.LineCount() {
		pos.Line = e.Buffer.LineCount() - 1
	}

	// Clamp column based on mode
	maxCol := e.Buffer.RuneCount(pos.Line)
	if e.Mode == ModeInsert {
		// In insert mode, can be at end of line
		if pos.Col > maxCol {
			pos.Col = maxCol
		}
	} else {
		// In normal mode, can't be past last character
		if maxCol > 0 {
			if pos.Col >= maxCol {
				pos.Col = maxCol - 1
			}
		} else {
			pos.Col = 0
		}
	}
	if pos.Col < 0 {
		pos.Col = 0
	}

	return pos
}

// saveUndo saves current state to undo stack.
func (e *Editor) saveUndo() {
	snapshot := Snapshot{
		Buffer: e.Buffer.Clone(),
		Cursor: e.Cursor,
	}
	e.UndoStack = append(e.UndoStack, snapshot)
	e.RedoStack = nil // Clear redo on new change
}

// Undo restores the previous state.
func (e *Editor) Undo() bool {
	if len(e.UndoStack) == 0 {
		e.StatusMessage = "Already at oldest change"
		return false
	}

	// Save current state to redo
	e.RedoStack = append(e.RedoStack, Snapshot{
		Buffer: e.Buffer.Clone(),
		Cursor: e.Cursor,
	})

	// Restore previous state
	prev := e.UndoStack[len(e.UndoStack)-1]
	e.UndoStack = e.UndoStack[:len(e.UndoStack)-1]
	e.Buffer = prev.Buffer
	e.Cursor = prev.Cursor

	return true
}

// Redo restores the next state.
func (e *Editor) Redo() bool {
	if len(e.RedoStack) == 0 {
		e.StatusMessage = "Already at newest change"
		return false
	}

	// Save current to undo
	e.UndoStack = append(e.UndoStack, Snapshot{
		Buffer: e.Buffer.Clone(),
		Cursor: e.Cursor,
	})

	// Restore redo state
	next := e.RedoStack[len(e.RedoStack)-1]
	e.RedoStack = e.RedoStack[:len(e.RedoStack)-1]
	e.Buffer = next.Buffer
	e.Cursor = next.Cursor

	return true
}

// resetCommandState resets the command parsing state.
func (e *Editor) resetCommandState() {
	e.PendingOp = OpNone
	e.Count = 0
	e.CountStack = nil
	e.WaitingFor = WaitNone
	if e.Mode == ModeOperatorPending {
		e.Mode = ModeNormal
	}
}

// getCount returns the effective count (1 if no count specified).
func (e *Editor) getCount() int {
	if e.Count == 0 {
		return 1
	}
	return e.Count
}

// EnterInsertMode switches to insert mode.
func (e *Editor) EnterInsertMode() {
	e.Mode = ModeInsert
	e.resetCommandState()
}

// EnterNormalMode switches to normal mode.
func (e *Editor) EnterNormalMode() {
	e.Mode = ModeNormal
	e.resetCommandState()
	// Adjust cursor - in normal mode, cursor should be on a character
	e.Cursor = e.clampPosition(e.Cursor)
}

// EnterVisualMode switches to visual mode.
func (e *Editor) EnterVisualMode(lineMode bool) {
	if lineMode {
		e.Mode = ModeVisualLine
	} else {
		e.Mode = ModeVisual
	}
	e.VisualStart = e.Cursor
	e.resetCommandState()
}

// GetVisualRange returns the selected range in visual mode.
func (e *Editor) GetVisualRange() Range {
	start := e.VisualStart
	end := e.Cursor

	// Normalize so start <= end
	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		start, end = end, start
	}

	if e.Mode == ModeVisualLine {
		return Range{
			Start:    Position{Line: start.Line, Col: 0},
			End:      Position{Line: end.Line, Col: e.Buffer.RuneCount(end.Line)},
			Linewise: true,
		}
	}

	// Character-wise visual mode: end is inclusive, so add 1 to col
	return Range{
		Start:    start,
		End:      Position{Line: end.Line, Col: end.Col + 1},
		Linewise: false,
	}
}
