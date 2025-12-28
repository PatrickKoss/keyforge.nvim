package vim

// OperatorType represents a vim operator
type OperatorType int

const (
	OpNone OperatorType = iota
	OpDelete
	OpChange
	OpYank
)

// String returns the operator character
func (o OperatorType) String() string {
	switch o {
	case OpDelete:
		return "d"
	case OpChange:
		return "c"
	case OpYank:
		return "y"
	default:
		return ""
	}
}

// ExecuteOperator applies an operator to a range
func (e *Editor) ExecuteOperator(op OperatorType, r Range) {
	switch op {
	case OpDelete:
		e.executeDelete(r)
	case OpChange:
		e.executeChange(r)
	case OpYank:
		e.executeYank(r)
	}
}

func (e *Editor) executeDelete(r Range) {
	e.saveUndo()

	if r.Linewise {
		// Delete entire lines
		var deleted string
		for i := r.Start.Line; i <= r.End.Line; i++ {
			if i > r.Start.Line {
				deleted += "\n"
			}
			deleted += e.Buffer.GetLine(r.Start.Line)
			e.Buffer.DeleteLine(r.Start.Line)
		}
		e.Unnamed = deleted

		// Position cursor
		e.Cursor.Line = r.Start.Line
		if e.Cursor.Line >= e.Buffer.LineCount() {
			e.Cursor.Line = e.Buffer.LineCount() - 1
		}
		// Move to first non-blank
		line := e.Buffer.GetLine(e.Cursor.Line)
		e.Cursor.Col = 0
		for i, ch := range []rune(line) {
			if ch != ' ' && ch != '\t' {
				e.Cursor.Col = i
				break
			}
		}
	} else {
		// Character-wise delete
		deleted := e.Buffer.DeleteRange(r.Start, r.End)
		e.Unnamed = deleted
		e.Cursor = r.Start
		e.Cursor = e.clampPosition(e.Cursor)
	}
}

func (e *Editor) executeChange(r Range) {
	e.executeDelete(r)
	e.EnterInsertMode()
}

func (e *Editor) executeYank(r Range) {
	if r.Linewise {
		var yanked string
		for i := r.Start.Line; i <= r.End.Line; i++ {
			if i > r.Start.Line {
				yanked += "\n"
			}
			yanked += e.Buffer.GetLine(i)
		}
		e.Unnamed = yanked
	} else {
		e.Unnamed = e.Buffer.GetRange(r.Start, r.End)
	}
	e.StatusMessage = "yanked"
}

// DeleteChar deletes the character under the cursor (x command)
func (e *Editor) DeleteChar(count int) {
	if count <= 0 {
		count = 1
	}

	e.saveUndo()

	line := e.Buffer.GetLine(e.Cursor.Line)
	runes := []rune(line)

	if e.Cursor.Col >= len(runes) {
		return
	}

	deleted := e.Buffer.DeleteAt(e.Cursor.Line, e.Cursor.Col, count)
	e.Unnamed = deleted

	// Adjust cursor if needed
	e.Cursor = e.clampPosition(e.Cursor)
}

// DeleteCharBefore deletes the character before the cursor (X command)
func (e *Editor) DeleteCharBefore(count int) {
	if count <= 0 {
		count = 1
	}

	if e.Cursor.Col == 0 {
		return
	}

	e.saveUndo()

	start := e.Cursor.Col - count
	if start < 0 {
		count = e.Cursor.Col
		start = 0
	}

	deleted := e.Buffer.DeleteAt(e.Cursor.Line, start, count)
	e.Unnamed = deleted
	e.Cursor.Col = start
}

// DeleteLine deletes the current line (dd command)
func (e *Editor) DeleteLine(count int) {
	if count <= 0 {
		count = 1
	}

	e.saveUndo()

	endLine := e.Cursor.Line + count - 1
	if endLine >= e.Buffer.LineCount() {
		endLine = e.Buffer.LineCount() - 1
	}

	r := Range{
		Start:    Position{Line: e.Cursor.Line, Col: 0},
		End:      Position{Line: endLine, Col: 0},
		Linewise: true,
	}

	e.executeDelete(r)
}

// YankLine yanks the current line (yy command)
func (e *Editor) YankLine(count int) {
	if count <= 0 {
		count = 1
	}

	endLine := e.Cursor.Line + count - 1
	if endLine >= e.Buffer.LineCount() {
		endLine = e.Buffer.LineCount() - 1
	}

	r := Range{
		Start:    Position{Line: e.Cursor.Line, Col: 0},
		End:      Position{Line: endLine, Col: 0},
		Linewise: true,
	}

	e.executeYank(r)
}

// Paste pastes from the unnamed register
func (e *Editor) Paste(before bool) {
	if e.Unnamed == "" {
		return
	}

	e.saveUndo()

	// Check if yanked content is linewise (contains newline)
	isLinewise := false
	for _, r := range e.Unnamed {
		if r == '\n' {
			isLinewise = true
			break
		}
	}

	if isLinewise {
		if before {
			e.Buffer.InsertLine(e.Cursor.Line, e.Unnamed)
		} else {
			e.Buffer.InsertLine(e.Cursor.Line+1, e.Unnamed)
			e.Cursor.Line++
		}
		// Move to first non-blank
		line := e.Buffer.GetLine(e.Cursor.Line)
		e.Cursor.Col = 0
		for i, ch := range []rune(line) {
			if ch != ' ' && ch != '\t' {
				e.Cursor.Col = i
				break
			}
		}
	} else {
		if before {
			e.Buffer.InsertAt(e.Cursor.Line, e.Cursor.Col, e.Unnamed)
		} else {
			e.Buffer.InsertAt(e.Cursor.Line, e.Cursor.Col+1, e.Unnamed)
			e.Cursor.Col += len([]rune(e.Unnamed))
		}
	}
}

// JoinLines joins current line with next line (J command)
func (e *Editor) JoinLines(count int) {
	if count <= 1 {
		count = 2 // J joins current with next, so minimum 2
	}

	if e.Cursor.Line >= e.Buffer.LineCount()-1 {
		return
	}

	e.saveUndo()

	for i := 0; i < count-1 && e.Cursor.Line < e.Buffer.LineCount()-1; i++ {
		// Get position for cursor after join (end of current line)
		line := e.Buffer.GetLine(e.Cursor.Line)
		joinPos := len([]rune(line))

		// Trim leading whitespace from next line and add space
		nextLine := e.Buffer.GetLine(e.Cursor.Line + 1)
		nextRunes := []rune(nextLine)
		start := 0
		for start < len(nextRunes) && (nextRunes[start] == ' ' || nextRunes[start] == '\t') {
			start++
		}
		trimmed := string(nextRunes[start:])

		// Join with space
		newLine := line
		if len(trimmed) > 0 {
			if len(line) > 0 {
				newLine = line + " " + trimmed
			} else {
				newLine = trimmed
			}
		}

		e.Buffer.SetLine(e.Cursor.Line, newLine)
		e.Buffer.DeleteLine(e.Cursor.Line + 1)

		e.Cursor.Col = joinPos
	}
}
