package vim

import (
	"unicode"
)

// TextObjectType represents different text object types.
type TextObjectType int

const (
	TextObjectNone        TextObjectType = iota
	TextObjectWord                       // iw, aw
	TextObjectWORD                       // iW, aW
	TextObjectSentence                   // is, as
	TextObjectParagraph                  // ip, ap
	TextObjectDoubleQuote                // i", a"
	TextObjectSingleQuote                // i', a'
	TextObjectBacktick                   // i`, a`
	TextObjectParen                      // i(, a(, i), a)
	TextObjectBracket                    // i[, a[, i], a]
	TextObjectBrace                      // i{, a{, i}, a}
	TextObjectAngle                      // i<, a<, i>, a>
)

// GetTextObjectRange returns the range for a text object.
func (e *Editor) GetTextObjectRange(objType TextObjectType, inner bool) (Range, bool) {
	switch objType {
	case TextObjectWord:
		return e.wordObjectRange(inner, false)
	case TextObjectWORD:
		return e.wordObjectRange(inner, true)
	case TextObjectDoubleQuote:
		return e.quoteObjectRange('"', inner)
	case TextObjectSingleQuote:
		return e.quoteObjectRange('\'', inner)
	case TextObjectBacktick:
		return e.quoteObjectRange('`', inner)
	case TextObjectParen:
		return e.pairObjectRange('(', ')', inner)
	case TextObjectBracket:
		return e.pairObjectRange('[', ']', inner)
	case TextObjectBrace:
		return e.pairObjectRange('{', '}', inner)
	case TextObjectAngle:
		return e.pairObjectRange('<', '>', inner)
	case TextObjectNone, TextObjectSentence, TextObjectParagraph:
		// Not implemented or no-op
		return Range{}, false
	}
	return Range{}, false
}

// wordObjectRange gets the range for iw/aw or iW/aW.
func (e *Editor) wordObjectRange(inner, bigWord bool) (Range, bool) {
	line := e.Buffer.GetLine(e.Cursor.Line)
	runes := []rune(line)
	col := e.Cursor.Col

	if len(runes) == 0 {
		return Range{}, false
	}

	if col >= len(runes) {
		col = len(runes) - 1
	}

	var isWordChar func(rune) bool
	if bigWord {
		isWordChar = func(r rune) bool { return !unicode.IsSpace(r) }
	} else {
		isWordChar = e.isWordChar
	}

	// Find word boundaries
	start := col
	end := col

	if isWordChar(runes[col]) {
		// On a word character
		for start > 0 && isWordChar(runes[start-1]) {
			start--
		}
		for end < len(runes)-1 && isWordChar(runes[end+1]) {
			end++
		}
	} else if !unicode.IsSpace(runes[col]) && !bigWord {
		// On punctuation (for small word)
		for start > 0 && !isWordChar(runes[start-1]) && !unicode.IsSpace(runes[start-1]) {
			start--
		}
		for end < len(runes)-1 && !isWordChar(runes[end+1]) && !unicode.IsSpace(runes[end+1]) {
			end++
		}
	} else {
		// On whitespace
		for start > 0 && unicode.IsSpace(runes[start-1]) {
			start--
		}
		for end < len(runes)-1 && unicode.IsSpace(runes[end+1]) {
			end++
		}
	}

	// For "around" word, include trailing whitespace
	if !inner {
		// Try to include trailing whitespace first
		trailingEnd := end
		for trailingEnd < len(runes)-1 && unicode.IsSpace(runes[trailingEnd+1]) {
			trailingEnd++
		}
		if trailingEnd > end {
			end = trailingEnd
		} else {
			// No trailing whitespace, include leading
			for start > 0 && unicode.IsSpace(runes[start-1]) {
				start--
			}
		}
	}

	return Range{
		Start:    Position{Line: e.Cursor.Line, Col: start},
		End:      Position{Line: e.Cursor.Line, Col: end + 1},
		Linewise: false,
	}, true
}

// quoteObjectRange gets the range for i"/a", i'/a', i`/a`.
func (e *Editor) quoteObjectRange(quote rune, inner bool) (Range, bool) {
	line := e.Buffer.GetLine(e.Cursor.Line)
	runes := []rune(line)
	col := e.Cursor.Col

	if len(runes) == 0 {
		return Range{}, false
	}

	// Find quote boundaries
	// Strategy: find the quote pair that contains or is closest to cursor
	start := -1
	end := -1

	// If cursor is on a quote, determine if it's start or end
	if col < len(runes) && runes[col] == quote {
		// Check if there's a matching quote before
		hasBefore := false
		for i := col - 1; i >= 0; i-- {
			if runes[i] == quote {
				hasBefore = true
				break
			}
		}
		if hasBefore {
			// This is the end quote, search backward
			end = col
			for i := col - 1; i >= 0; i-- {
				if runes[i] == quote {
					start = i
					break
				}
			}
		} else {
			// This is the start quote, search forward
			start = col
			for i := col + 1; i < len(runes); i++ {
				if runes[i] == quote {
					end = i
					break
				}
			}
		}
	} else {
		// Search outward for quotes
		// First, try to find quotes surrounding cursor
		for i := col; i >= 0; i-- {
			if runes[i] == quote {
				start = i
				break
			}
		}
		if start >= 0 {
			for i := col + 1; i < len(runes); i++ {
				if runes[i] == quote {
					end = i
					break
				}
			}
		}
		// If that didn't work, search forward
		if start < 0 || end < 0 {
			start = -1
			for i := col; i < len(runes); i++ {
				if runes[i] == quote {
					if start < 0 {
						start = i
					} else {
						end = i
						break
					}
				}
			}
		}
	}

	if start < 0 || end < 0 || start >= end {
		return Range{}, false
	}

	if inner {
		return Range{
			Start:    Position{Line: e.Cursor.Line, Col: start + 1},
			End:      Position{Line: e.Cursor.Line, Col: end},
			Linewise: false,
		}, true
	}

	return Range{
		Start:    Position{Line: e.Cursor.Line, Col: start},
		End:      Position{Line: e.Cursor.Line, Col: end + 1},
		Linewise: false,
	}, true
}

// pairObjectRange gets the range for i(/a(, i[/a[, i{/a{, i</a<.
func (e *Editor) pairObjectRange(open, closeBracket rune, inner bool) (Range, bool) {
	// Search for the innermost pair containing the cursor
	// This needs to handle multi-line pairs

	startPos := Position{Line: -1, Col: -1}
	endPos := Position{Line: -1, Col: -1}

	// Search backward for opening bracket
	depth := 0
	found := false

	// Start from cursor
	line := e.Cursor.Line
	col := e.Cursor.Col

	// If on the close bracket, include it
	currentLine := e.Buffer.GetLine(line)
	currentRunes := []rune(currentLine)
	if col < len(currentRunes) && currentRunes[col] == closeBracket {
		depth = 1
		endPos = Position{Line: line, Col: col}
	}

	// Search backward for open bracket
	for lineNum := line; lineNum >= 0 && !found; lineNum-- {
		lineContent := e.Buffer.GetLine(lineNum)
		runes := []rune(lineContent)

		startCol := len(runes) - 1
		if lineNum == line {
			startCol = col
			if depth == 1 {
				startCol = col - 1 // Don't count the close bracket we're on
			}
		}

		for c := startCol; c >= 0; c-- {
			if runes[c] == closeBracket {
				depth++
			} else if runes[c] == open {
				if depth > 0 {
					depth--
				}
				if depth == 0 {
					startPos = Position{Line: lineNum, Col: c}
					found = true
					break
				}
			}
		}
	}

	if !found {
		return Range{}, false
	}

	// Now search forward for the matching close bracket
	depth = 1
	found = false

	for lineNum := startPos.Line; lineNum < e.Buffer.LineCount() && !found; lineNum++ {
		lineContent := e.Buffer.GetLine(lineNum)
		runes := []rune(lineContent)

		startCol := 0
		if lineNum == startPos.Line {
			startCol = startPos.Col + 1
		}

		for c := startCol; c < len(runes); c++ {
			if runes[c] == open {
				depth++
			} else if runes[c] == closeBracket {
				depth--
				if depth == 0 {
					endPos = Position{Line: lineNum, Col: c}
					found = true
					break
				}
			}
		}
	}

	if !found {
		return Range{}, false
	}

	if inner {
		// Adjust for inner - exclude the brackets
		innerStart := Position{Line: startPos.Line, Col: startPos.Col + 1}
		innerEnd := endPos

		// If opening bracket is followed by newline, start from next line col 0
		startLine := e.Buffer.GetLine(startPos.Line)
		startRunes := []rune(startLine)
		if startPos.Col == len(startRunes)-1 && startPos.Line < endPos.Line {
			innerStart = Position{Line: startPos.Line + 1, Col: 0}
		}

		return Range{
			Start:    innerStart,
			End:      innerEnd,
			Linewise: false,
		}, true
	}

	return Range{
		Start:    startPos,
		End:      Position{Line: endPos.Line, Col: endPos.Col + 1},
		Linewise: false,
	}, true
}

// ParseTextObject parses text object keys like "iw", "a\"", "i(".
func ParseTextObject(key1, key2 string) (TextObjectType, bool) {
	inner := key1 == "i"
	if key1 != "i" && key1 != "a" {
		return TextObjectNone, false
	}

	var objType TextObjectType
	switch key2 {
	case "w":
		objType = TextObjectWord
	case "W":
		objType = TextObjectWORD
	case "\"":
		objType = TextObjectDoubleQuote
	case "'":
		objType = TextObjectSingleQuote
	case "`":
		objType = TextObjectBacktick
	case "(", ")", "b":
		objType = TextObjectParen
	case "[", "]":
		objType = TextObjectBracket
	case "{", "}", "B":
		objType = TextObjectBrace
	case "<", ">":
		objType = TextObjectAngle
	default:
		return TextObjectNone, false
	}

	return objType, inner
}
