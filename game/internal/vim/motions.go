package vim

import (
	"unicode"
)

// Motion types.
type MotionType int

const (
	MotionNone           MotionType = iota
	MotionLeft                      // h
	MotionRight                     // l
	MotionUp                        // k
	MotionDown                      // j
	MotionLineStart                 // 0
	MotionLineEnd                   // $
	MotionFirstNonBlank             // ^
	MotionWordForward               // w
	MotionWordBackward              // b
	MotionWordEnd                   // e
	MotionWORDForward               // W
	MotionWORDBackward              // B
	MotionWORDEnd                   // E
	MotionFileStart                 // gg
	MotionFileEnd                   // G
	MotionFindChar                  // f{char}
	MotionFindCharBack              // F{char}
	MotionTillChar                  // t{char}
	MotionTillCharBack              // T{char}
	MotionRepeatFind                // ;
	MotionRepeatFindBack            // ,
	MotionMatchBracket              // %
)

// ExecuteMotion moves the cursor based on the motion type.
func (e *Editor) ExecuteMotion(motion MotionType, count int) Position {
	if count <= 0 {
		count = 1
	}

	pos := e.Cursor
	buf := e.Buffer

	switch motion {
	case MotionLeft:
		pos.Col -= count
		if pos.Col < 0 {
			pos.Col = 0
		}

	case MotionRight:
		maxCol := buf.LastCol(pos.Line)
		pos.Col += count
		if pos.Col > maxCol {
			pos.Col = maxCol
		}

	case MotionUp:
		pos.Line -= count
		if pos.Line < 0 {
			pos.Line = 0
		}
		// Adjust column for new line length
		maxCol := buf.LastCol(pos.Line)
		if pos.Col > maxCol {
			pos.Col = maxCol
		}

	case MotionDown:
		pos.Line += count
		if pos.Line >= buf.LineCount() {
			pos.Line = buf.LineCount() - 1
		}
		// Adjust column for new line length
		maxCol := buf.LastCol(pos.Line)
		if pos.Col > maxCol {
			pos.Col = maxCol
		}

	case MotionLineStart:
		pos.Col = 0

	case MotionLineEnd:
		pos.Col = buf.LastCol(pos.Line)

	case MotionFirstNonBlank:
		line := buf.GetLine(pos.Line)
		runes := []rune(line)
		pos.Col = 0
		for i, r := range runes {
			if !unicode.IsSpace(r) {
				pos.Col = i
				break
			}
		}

	case MotionWordForward:
		for range count {
			pos = e.nextWord(pos)
		}

	case MotionWordBackward:
		for range count {
			pos = e.prevWord(pos)
		}

	case MotionWordEnd:
		for range count {
			pos = e.wordEnd(pos)
		}

	case MotionWORDForward:
		for range count {
			pos = e.nextWORD(pos)
		}

	case MotionWORDBackward:
		for range count {
			pos = e.prevWORD(pos)
		}

	case MotionWORDEnd:
		for range count {
			pos = e.WORDEnd(pos)
		}

	case MotionFileStart:
		pos.Line = 0
		pos.Col = 0
		// Find first non-blank
		line := buf.GetLine(0)
		for i, r := range []rune(line) {
			if !unicode.IsSpace(r) {
				pos.Col = i
				break
			}
		}

	case MotionFileEnd:
		targetLine := buf.LineCount() - 1
		if count > 1 {
			targetLine = count - 1
			if targetLine >= buf.LineCount() {
				targetLine = buf.LineCount() - 1
			}
		}
		pos.Line = targetLine
		// Find first non-blank
		line := buf.GetLine(pos.Line)
		pos.Col = 0
		for i, r := range []rune(line) {
			if !unicode.IsSpace(r) {
				pos.Col = i
				break
			}
		}

	case MotionMatchBracket:
		pos = e.matchBracket(pos)

	case MotionNone, MotionFindChar, MotionFindCharBack, MotionTillChar, MotionTillCharBack, MotionRepeatFind, MotionRepeatFindBack:
		// These motions are handled separately or are no-ops
	}

	return pos
}

// ExecuteFindMotion executes f/F/t/T motion with the given character.
func (e *Editor) ExecuteFindMotion(forward, till bool, char rune, count int) Position {
	if count <= 0 {
		count = 1
	}

	pos := e.Cursor
	line := e.Buffer.GetLine(pos.Line)
	runes := []rune(line)

	// Save for repeat
	e.LastFind = FindState{
		Char:    char,
		Forward: forward,
		Till:    till,
	}

	if forward {
		// Search forward
		found := 0
		for i := pos.Col + 1; i < len(runes); i++ {
			if runes[i] == char {
				found++
				if found == count {
					if till {
						pos.Col = i - 1
					} else {
						pos.Col = i
					}
					break
				}
			}
		}
	} else {
		// Search backward
		found := 0
		for i := pos.Col - 1; i >= 0; i-- {
			if runes[i] == char {
				found++
				if found == count {
					if till {
						pos.Col = i + 1
					} else {
						pos.Col = i
					}
					break
				}
			}
		}
	}

	return pos
}

// RepeatFind repeats the last f/F/t/T motion.
func (e *Editor) RepeatFind(reverse bool, count int) Position {
	if e.LastFind.Char == 0 {
		return e.Cursor
	}

	forward := e.LastFind.Forward
	if reverse {
		forward = !forward
	}

	return e.ExecuteFindMotion(forward, e.LastFind.Till, e.LastFind.Char, count)
}

// Helper functions for word motions

func (e *Editor) isWordChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

func (e *Editor) nextWord(pos Position) Position {
	buf := e.Buffer
	line := buf.GetLine(pos.Line)
	runes := []rune(line)

	// If at end of line, go to next line
	if pos.Col >= len(runes) {
		if pos.Line < buf.LineCount()-1 {
			pos.Line++
			pos.Col = 0
			// Skip to first non-space
			line = buf.GetLine(pos.Line)
			runes = []rune(line)
			for pos.Col < len(runes) && unicode.IsSpace(runes[pos.Col]) {
				pos.Col++
			}
		}
		return pos
	}

	// Skip current word
	if pos.Col < len(runes) {
		if e.isWordChar(runes[pos.Col]) {
			// In a word, skip to end
			for pos.Col < len(runes) && e.isWordChar(runes[pos.Col]) {
				pos.Col++
			}
		} else if !unicode.IsSpace(runes[pos.Col]) {
			// In punctuation, skip to end
			for pos.Col < len(runes) && !e.isWordChar(runes[pos.Col]) && !unicode.IsSpace(runes[pos.Col]) {
				pos.Col++
			}
		}
	}

	// Skip whitespace (possibly across lines)
	for {
		line = buf.GetLine(pos.Line)
		runes = []rune(line)

		for pos.Col < len(runes) && unicode.IsSpace(runes[pos.Col]) {
			pos.Col++
		}

		if pos.Col < len(runes) {
			break
		}

		// End of line, go to next
		if pos.Line < buf.LineCount()-1 {
			pos.Line++
			pos.Col = 0
		} else {
			break
		}
	}

	return pos
}

func (e *Editor) prevWord(pos Position) Position {
	buf := e.Buffer
	line := buf.GetLine(pos.Line)
	runes := []rune(line)

	// If at beginning of line, go to previous line
	if pos.Col == 0 {
		if pos.Line > 0 {
			pos.Line--
			line = buf.GetLine(pos.Line)
			runes = []rune(line)
			pos.Col = len(runes)
		} else {
			return pos
		}
	}

	// Skip whitespace backward
	for pos.Col > 0 && pos.Col <= len(runes) && unicode.IsSpace(runes[pos.Col-1]) {
		pos.Col--
	}

	// Handle crossing line boundary
	if pos.Col == 0 {
		if pos.Line > 0 {
			pos.Line--
			line = buf.GetLine(pos.Line)
			runes = []rune(line)
			pos.Col = len(runes)
			// Continue skipping whitespace
			for pos.Col > 0 && unicode.IsSpace(runes[pos.Col-1]) {
				pos.Col--
			}
		}
	}

	if pos.Col == 0 {
		return pos
	}

	// Skip to beginning of word
	if pos.Col > 0 && pos.Col <= len(runes) {
		if e.isWordChar(runes[pos.Col-1]) {
			for pos.Col > 0 && e.isWordChar(runes[pos.Col-1]) {
				pos.Col--
			}
		} else if !unicode.IsSpace(runes[pos.Col-1]) {
			for pos.Col > 0 && !e.isWordChar(runes[pos.Col-1]) && !unicode.IsSpace(runes[pos.Col-1]) {
				pos.Col--
			}
		}
	}

	return pos
}

func (e *Editor) wordEnd(pos Position) Position {
	buf := e.Buffer

	// Move at least one character
	pos.Col++

	// Skip whitespace
	var runes []rune
	for {
		line := buf.GetLine(pos.Line)
		runes = []rune(line)

		for pos.Col < len(runes) && unicode.IsSpace(runes[pos.Col]) {
			pos.Col++
		}

		if pos.Col < len(runes) {
			break
		}

		if pos.Line < buf.LineCount()-1 {
			pos.Line++
			pos.Col = 0
		} else {
			// End of file
			pos.Col = buf.LastCol(pos.Line)
			return pos
		}
	}

	// Move to end of word
	if pos.Col < len(runes) {
		if e.isWordChar(runes[pos.Col]) {
			for pos.Col < len(runes)-1 && e.isWordChar(runes[pos.Col+1]) {
				pos.Col++
			}
		} else {
			for pos.Col < len(runes)-1 && !e.isWordChar(runes[pos.Col+1]) && !unicode.IsSpace(runes[pos.Col+1]) {
				pos.Col++
			}
		}
	}

	return pos
}

// WORD motions (space-separated)

func (e *Editor) nextWORD(pos Position) Position {
	buf := e.Buffer

	// Skip current WORD (non-space characters)
	for {
		line := buf.GetLine(pos.Line)
		runes := []rune(line)

		for pos.Col < len(runes) && !unicode.IsSpace(runes[pos.Col]) {
			pos.Col++
		}

		// Skip whitespace
		for pos.Col < len(runes) && unicode.IsSpace(runes[pos.Col]) {
			pos.Col++
		}

		if pos.Col < len(runes) {
			break
		}

		if pos.Line < buf.LineCount()-1 {
			pos.Line++
			pos.Col = 0
		} else {
			break
		}
	}

	return pos
}

func (e *Editor) prevWORD(pos Position) Position {
	buf := e.Buffer

	// Handle beginning of line
	if pos.Col == 0 {
		if pos.Line > 0 {
			pos.Line--
			pos.Col = buf.RuneCount(pos.Line)
		} else {
			return pos
		}
	}

	line := buf.GetLine(pos.Line)
	runes := []rune(line)

	// Skip whitespace backward
	for pos.Col > 0 && unicode.IsSpace(runes[pos.Col-1]) {
		pos.Col--
	}

	// Skip WORD backward
	for pos.Col > 0 && !unicode.IsSpace(runes[pos.Col-1]) {
		pos.Col--
	}

	return pos
}

func (e *Editor) WORDEnd(pos Position) Position {
	buf := e.Buffer

	// Move at least one character
	pos.Col++

	// Skip whitespace
	for {
		line := buf.GetLine(pos.Line)
		runes := []rune(line)

		for pos.Col < len(runes) && unicode.IsSpace(runes[pos.Col]) {
			pos.Col++
		}

		if pos.Col < len(runes) {
			break
		}

		if pos.Line < buf.LineCount()-1 {
			pos.Line++
			pos.Col = 0
		} else {
			pos.Col = buf.LastCol(pos.Line)
			return pos
		}
	}

	// Move to end of WORD
	line := buf.GetLine(pos.Line)
	runes := []rune(line)
	for pos.Col < len(runes)-1 && !unicode.IsSpace(runes[pos.Col+1]) {
		pos.Col++
	}

	return pos
}

// matchBracket finds the matching bracket.
func (e *Editor) matchBracket(pos Position) Position {
	line := e.Buffer.GetLine(pos.Line)
	runes := []rune(line)

	if pos.Col >= len(runes) {
		return pos
	}

	char := runes[pos.Col]
	var match rune
	var forward bool

	switch char {
	case '(':
		match, forward = ')', true
	case ')':
		match, forward = '(', false
	case '[':
		match, forward = ']', true
	case ']':
		match, forward = '[', false
	case '{':
		match, forward = '}', true
	case '}':
		match, forward = '{', false
	default:
		// Not on a bracket, search forward for one
		for i := pos.Col; i < len(runes); i++ {
			switch runes[i] {
			case '(', '[', '{':
				pos.Col = i
				return e.matchBracket(pos)
			}
		}
		return pos
	}

	// Search for matching bracket
	depth := 1
	if forward {
		for lineNum := pos.Line; lineNum < e.Buffer.LineCount(); lineNum++ {
			l := e.Buffer.GetLine(lineNum)
			r := []rune(l)
			startCol := 0
			if lineNum == pos.Line {
				startCol = pos.Col + 1
			}
			for col := startCol; col < len(r); col++ {
				switch r[col] {
				case char:
					depth++
				case match:
					depth--
					if depth == 0 {
						return Position{Line: lineNum, Col: col}
					}
				}
			}
		}
	} else {
		for lineNum := pos.Line; lineNum >= 0; lineNum-- {
			l := e.Buffer.GetLine(lineNum)
			r := []rune(l)
			endCol := len(r) - 1
			if lineNum == pos.Line {
				endCol = pos.Col - 1
			}
			for col := endCol; col >= 0; col-- {
				switch r[col] {
				case char:
					depth++
				case match:
					depth--
					if depth == 0 {
						return Position{Line: lineNum, Col: col}
					}
				}
			}
		}
	}

	return pos
}

// GetMotionRange returns the range affected by a motion from current cursor.
func (e *Editor) GetMotionRange(motion MotionType, count int) Range {
	start := e.Cursor
	end := e.ExecuteMotion(motion, count)

	// Normalize start/end
	if start.Line > end.Line || (start.Line == end.Line && start.Col > end.Col) {
		start, end = end, start
	}

	linewise := motion == MotionFileStart || motion == MotionFileEnd

	return Range{
		Start:    start,
		End:      end,
		Linewise: linewise,
	}
}
