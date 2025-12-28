package vim

import (
	"strconv"
	"unicode"
)

// HandleKey processes a single keypress and returns if more input is needed
func (e *Editor) HandleKey(key string) bool {
	e.KeystrokeCount++
	e.StatusMessage = "" // Clear status on new key

	switch e.Mode {
	case ModeInsert:
		return e.handleInsertKey(key)
	case ModeNormal, ModeOperatorPending:
		return e.handleNormalKey(key)
	case ModeVisual, ModeVisualLine:
		return e.handleVisualKey(key)
	}
	return false
}

func (e *Editor) handleInsertKey(key string) bool {
	switch key {
	case "Escape", "ctrl+c", "ctrl+[":
		e.EnterNormalMode()
		// Move cursor back one if possible
		if e.Cursor.Col > 0 {
			e.Cursor.Col--
		}
		return false

	case "Backspace", "ctrl+h":
		if e.Cursor.Col > 0 {
			e.Buffer.DeleteAt(e.Cursor.Line, e.Cursor.Col-1, 1)
			e.Cursor.Col--
		} else if e.Cursor.Line > 0 {
			// Join with previous line
			prevLine := e.Buffer.GetLine(e.Cursor.Line - 1)
			e.Cursor.Col = len([]rune(prevLine))
			e.Buffer.JoinLines(e.Cursor.Line - 1)
			e.Cursor.Line--
		}
		return false

	case "Delete":
		line := e.Buffer.GetLine(e.Cursor.Line)
		if e.Cursor.Col < len([]rune(line)) {
			e.Buffer.DeleteAt(e.Cursor.Line, e.Cursor.Col, 1)
		} else if e.Cursor.Line < e.Buffer.LineCount()-1 {
			// Join with next line
			e.Buffer.JoinLines(e.Cursor.Line)
		}
		return false

	case "Enter":
		e.saveUndo()
		e.Buffer.SplitLine(e.Cursor.Line, e.Cursor.Col)
		e.Cursor.Line++
		e.Cursor.Col = 0
		return false

	case "Tab":
		e.Buffer.InsertAt(e.Cursor.Line, e.Cursor.Col, "\t")
		e.Cursor.Col++
		return false

	case "Left":
		if e.Cursor.Col > 0 {
			e.Cursor.Col--
		}
		return false

	case "Right":
		line := e.Buffer.GetLine(e.Cursor.Line)
		if e.Cursor.Col < len([]rune(line)) {
			e.Cursor.Col++
		}
		return false

	case "Up":
		if e.Cursor.Line > 0 {
			e.Cursor.Line--
			e.Cursor = e.clampPosition(e.Cursor)
		}
		return false

	case "Down":
		if e.Cursor.Line < e.Buffer.LineCount()-1 {
			e.Cursor.Line++
			e.Cursor = e.clampPosition(e.Cursor)
		}
		return false

	default:
		// Insert the character
		if len(key) == 1 {
			e.Buffer.InsertAt(e.Cursor.Line, e.Cursor.Col, key)
			e.Cursor.Col++
		}
		return false
	}
}

func (e *Editor) handleNormalKey(key string) bool {
	// Handle waiting states first
	switch e.WaitingFor {
	case WaitChar:
		// f/F/t/T waiting for character
		if len(key) == 1 {
			r := []rune(key)[0]
			newPos := e.ExecuteFindMotion(e.LastFind.Forward, e.LastFind.Till, r, e.getCount())
			if e.PendingOp != OpNone {
				// Operator pending
				rng := Range{Start: e.Cursor, End: newPos}
				if rng.Start.Col > rng.End.Col {
					rng.Start, rng.End = rng.End, rng.Start
				}
				rng.End.Col++ // Make end inclusive
				e.ExecuteOperator(e.PendingOp, rng)
			} else {
				e.Cursor = newPos
			}
		}
		e.resetCommandState()
		return false

	case WaitMotion:
		// After 'i' or 'a' for text object
		if key == "i" || key == "a" {
			e.WaitingFor = WaitNone
			// Will be handled below as text object prefix
		}
	}

	// Build count
	if len(key) == 1 {
		r := []rune(key)[0]
		if unicode.IsDigit(r) && (r != '0' || e.Count > 0) {
			digit, _ := strconv.Atoi(key)
			e.Count = e.Count*10 + digit
			return true
		}
	}

	// Check for operators
	switch key {
	case "d":
		if e.PendingOp == OpDelete {
			// dd - delete line
			e.DeleteLine(e.getCount())
			e.resetCommandState()
			return false
		}
		e.PendingOp = OpDelete
		e.Mode = ModeOperatorPending
		return true

	case "c":
		if e.PendingOp == OpChange {
			// cc - change line
			e.saveUndo()
			count := e.getCount()
			endLine := e.Cursor.Line + count - 1
			if endLine >= e.Buffer.LineCount() {
				endLine = e.Buffer.LineCount() - 1
			}
			// Delete line contents but keep the line
			for i := e.Cursor.Line; i <= endLine; i++ {
				e.Buffer.SetLine(e.Cursor.Line, "")
				if i < endLine {
					e.Buffer.DeleteLine(e.Cursor.Line + 1)
				}
			}
			e.Cursor.Col = 0
			e.EnterInsertMode()
			e.resetCommandState()
			return false
		}
		e.PendingOp = OpChange
		e.Mode = ModeOperatorPending
		return true

	case "y":
		if e.PendingOp == OpYank {
			// yy - yank line
			e.YankLine(e.getCount())
			e.resetCommandState()
			return false
		}
		e.PendingOp = OpYank
		e.Mode = ModeOperatorPending
		return true
	}

	// Text object prefixes (only valid after operator)
	if e.PendingOp != OpNone && e.WaitingFor == WaitNone && (key == "i" || key == "a") {
		e.WaitingFor = WaitMotion
		// Store 'i' or 'a' for next key
		if key == "i" {
			e.CountStack = []int{1} // marker for inner
		} else {
			e.CountStack = []int{0} // marker for around
		}
		return true
	}

	// Text object second character (after 'i' or 'a')
	if e.PendingOp != OpNone && e.WaitingFor == WaitMotion && len(e.CountStack) > 0 {
		inner := e.CountStack[0] == 1
		objType, _ := ParseTextObject(map[bool]string{true: "i", false: "a"}[inner], key)
		if objType != TextObjectNone {
			rng, ok := e.GetTextObjectRange(objType, inner)
			if ok {
				e.ExecuteOperator(e.PendingOp, rng)
			}
			e.resetCommandState()
			return false
		}
		// Not a valid text object, reset
		e.resetCommandState()
		return false
	}

	// Check for motions
	motion := e.keyToMotion(key)
	if motion != MotionNone {
		count := e.getCount()
		if e.PendingOp != OpNone {
			// Operator + motion
			start := e.Cursor
			end := e.ExecuteMotion(motion, count)

			// Normalize so start <= end
			if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
				start, end = end, start
			}

			// For word motions, include the character at end position
			if motion == MotionWordForward || motion == MotionWordEnd ||
				motion == MotionWORDForward || motion == MotionWORDEnd {
				// dw should delete up to (but not including) next word
				// The end position IS the start of next word, so we use it as-is
			} else {
				// For most motions, the end is inclusive
				end.Col++
			}

			rng := Range{Start: start, End: end, Linewise: false}
			if motion == MotionFileStart || motion == MotionFileEnd {
				rng.Linewise = true
			}
			e.ExecuteOperator(e.PendingOp, rng)
		} else {
			// Just motion
			e.Cursor = e.ExecuteMotion(motion, count)
		}
		e.resetCommandState()
		return false
	}

	// Find/till commands
	switch key {
	case "f":
		e.LastFind.Forward = true
		e.LastFind.Till = false
		e.WaitingFor = WaitChar
		return true
	case "F":
		e.LastFind.Forward = false
		e.LastFind.Till = false
		e.WaitingFor = WaitChar
		return true
	case "t":
		e.LastFind.Forward = true
		e.LastFind.Till = true
		e.WaitingFor = WaitChar
		return true
	case "T":
		e.LastFind.Forward = false
		e.LastFind.Till = true
		e.WaitingFor = WaitChar
		return true
	case ";":
		newPos := e.RepeatFind(false, e.getCount())
		if e.PendingOp != OpNone {
			rng := Range{Start: e.Cursor, End: Position{Line: newPos.Line, Col: newPos.Col + 1}}
			e.ExecuteOperator(e.PendingOp, rng)
		} else {
			e.Cursor = newPos
		}
		e.resetCommandState()
		return false
	case ",":
		newPos := e.RepeatFind(true, e.getCount())
		if e.PendingOp != OpNone {
			rng := Range{Start: e.Cursor, End: Position{Line: newPos.Line, Col: newPos.Col + 1}}
			e.ExecuteOperator(e.PendingOp, rng)
		} else {
			e.Cursor = newPos
		}
		e.resetCommandState()
		return false
	}

	// Simple commands
	switch key {
	case "x":
		e.DeleteChar(e.getCount())
		e.resetCommandState()
		return false

	case "X":
		e.DeleteCharBefore(e.getCount())
		e.resetCommandState()
		return false

	case "s":
		e.DeleteChar(e.getCount())
		e.EnterInsertMode()
		e.resetCommandState()
		return false

	case "S":
		e.saveUndo()
		e.Buffer.SetLine(e.Cursor.Line, "")
		e.Cursor.Col = 0
		e.EnterInsertMode()
		e.resetCommandState()
		return false

	case "r":
		// Replace character - need to wait for next char
		e.WaitingFor = WaitChar
		e.PendingOp = OpNone
		e.CountStack = []int{-1} // marker for replace
		return true

	case "i":
		e.EnterInsertMode()
		return false

	case "I":
		// Insert at first non-blank
		e.Cursor = e.ExecuteMotion(MotionFirstNonBlank, 1)
		e.EnterInsertMode()
		return false

	case "a":
		// Append after cursor
		line := e.Buffer.GetLine(e.Cursor.Line)
		if len([]rune(line)) > 0 {
			e.Cursor.Col++
		}
		e.EnterInsertMode()
		return false

	case "A":
		// Append at end of line
		e.Cursor.Col = e.Buffer.RuneCount(e.Cursor.Line)
		e.EnterInsertMode()
		return false

	case "o":
		// Open line below
		e.saveUndo()
		e.Buffer.InsertLine(e.Cursor.Line+1, "")
		e.Cursor.Line++
		e.Cursor.Col = 0
		e.EnterInsertMode()
		return false

	case "O":
		// Open line above
		e.saveUndo()
		e.Buffer.InsertLine(e.Cursor.Line, "")
		e.Cursor.Col = 0
		e.EnterInsertMode()
		return false

	case "p":
		e.Paste(false)
		e.resetCommandState()
		return false

	case "P":
		e.Paste(true)
		e.resetCommandState()
		return false

	case "u":
		e.Undo()
		e.resetCommandState()
		return false

	case "ctrl+r":
		e.Redo()
		e.resetCommandState()
		return false

	case "J":
		e.JoinLines(e.getCount() + 1)
		e.resetCommandState()
		return false

	case "v":
		e.EnterVisualMode(false)
		return false

	case "V":
		e.EnterVisualMode(true)
		return false

	case "g":
		// Waiting for second key (gg)
		e.CountStack = append(e.CountStack, -2) // marker for g prefix
		return true

	case "Escape":
		e.resetCommandState()
		return false
	}

	// Handle 'g' prefix commands
	if len(e.CountStack) > 0 && e.CountStack[len(e.CountStack)-1] == -2 {
		e.CountStack = e.CountStack[:len(e.CountStack)-1]
		switch key {
		case "g":
			// gg - go to start
			if e.PendingOp != OpNone {
				rng := e.GetMotionRange(MotionFileStart, 1)
				rng.Linewise = true
				e.ExecuteOperator(e.PendingOp, rng)
			} else {
				e.Cursor = e.ExecuteMotion(MotionFileStart, 1)
			}
		}
		e.resetCommandState()
		return false
	}

	// Handle replace character
	if len(e.CountStack) > 0 && e.CountStack[0] == -1 {
		if len(key) == 1 {
			e.saveUndo()
			count := e.getCount()
			line := e.Buffer.GetLine(e.Cursor.Line)
			runes := []rune(line)
			for i := 0; i < count && e.Cursor.Col+i < len(runes); i++ {
				runes[e.Cursor.Col+i] = []rune(key)[0]
			}
			e.Buffer.SetLine(e.Cursor.Line, string(runes))
		}
		e.resetCommandState()
		return false
	}

	e.resetCommandState()
	return false
}

func (e *Editor) handleVisualKey(key string) bool {
	// Build count
	if len(key) == 1 {
		r := []rune(key)[0]
		if unicode.IsDigit(r) && (r != '0' || e.Count > 0) {
			digit, _ := strconv.Atoi(key)
			e.Count = e.Count*10 + digit
			return true
		}
	}

	// Check for operators - apply to selection
	switch key {
	case "d", "x":
		rng := e.GetVisualRange()
		e.ExecuteOperator(OpDelete, rng)
		e.EnterNormalMode()
		return false

	case "c", "s":
		rng := e.GetVisualRange()
		e.ExecuteOperator(OpChange, rng)
		return false

	case "y":
		rng := e.GetVisualRange()
		e.ExecuteOperator(OpYank, rng)
		e.EnterNormalMode()
		return false

	case "Escape", "ctrl+c", "ctrl+[":
		e.EnterNormalMode()
		return false

	case "v":
		if e.Mode == ModeVisual {
			e.EnterNormalMode()
		} else {
			e.Mode = ModeVisual
		}
		return false

	case "V":
		if e.Mode == ModeVisualLine {
			e.EnterNormalMode()
		} else {
			e.Mode = ModeVisualLine
		}
		return false

	case "o":
		// Swap cursor and visual start
		e.Cursor, e.VisualStart = e.VisualStart, e.Cursor
		return false
	}

	// Motions extend selection
	motion := e.keyToMotion(key)
	if motion != MotionNone {
		e.Cursor = e.ExecuteMotion(motion, e.getCount())
		e.Count = 0
		return false
	}

	// Find commands
	switch key {
	case "f", "F", "t", "T":
		e.LastFind.Forward = key == "f" || key == "t"
		e.LastFind.Till = key == "t" || key == "T"
		e.WaitingFor = WaitChar
		return true
	}

	if e.WaitingFor == WaitChar && len(key) == 1 {
		r := []rune(key)[0]
		e.Cursor = e.ExecuteFindMotion(e.LastFind.Forward, e.LastFind.Till, r, e.getCount())
		e.WaitingFor = WaitNone
		e.Count = 0
		return false
	}

	e.Count = 0
	return false
}

// keyToMotion converts a key to a motion type
func (e *Editor) keyToMotion(key string) MotionType {
	switch key {
	case "h", "Left":
		return MotionLeft
	case "l", "Right":
		return MotionRight
	case "k", "Up":
		return MotionUp
	case "j", "Down":
		return MotionDown
	case "0":
		return MotionLineStart
	case "$":
		return MotionLineEnd
	case "^":
		return MotionFirstNonBlank
	case "w":
		return MotionWordForward
	case "b":
		return MotionWordBackward
	case "e":
		return MotionWordEnd
	case "W":
		return MotionWORDForward
	case "B":
		return MotionWORDBackward
	case "E":
		return MotionWORDEnd
	case "G":
		return MotionFileEnd
	case "%":
		return MotionMatchBracket
	}
	return MotionNone
}
