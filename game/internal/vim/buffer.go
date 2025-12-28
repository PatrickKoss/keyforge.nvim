package vim

import (
	"strings"
	"unicode/utf8"
)

// Buffer represents an editable text buffer as lines
type Buffer struct {
	lines []string
}

// NewBuffer creates a buffer from initial text
func NewBuffer(text string) *Buffer {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &Buffer{lines: lines}
}

// LineCount returns the number of lines in the buffer
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// GetLine returns the content of line n (0-based)
func (b *Buffer) GetLine(n int) string {
	if n < 0 || n >= len(b.lines) {
		return ""
	}
	return b.lines[n]
}

// SetLine sets the content of line n
func (b *Buffer) SetLine(n int, content string) {
	if n < 0 || n >= len(b.lines) {
		return
	}
	b.lines[n] = content
}

// InsertLine inserts a new line at position n
func (b *Buffer) InsertLine(n int, content string) {
	if n < 0 {
		n = 0
	}
	if n > len(b.lines) {
		n = len(b.lines)
	}
	b.lines = append(b.lines[:n], append([]string{content}, b.lines[n:]...)...)
}

// DeleteLine removes and returns the line at position n
func (b *Buffer) DeleteLine(n int) string {
	if n < 0 || n >= len(b.lines) {
		return ""
	}
	deleted := b.lines[n]
	b.lines = append(b.lines[:n], b.lines[n+1:]...)
	// Ensure at least one line
	if len(b.lines) == 0 {
		b.lines = []string{""}
	}
	return deleted
}

// SplitLine splits line at column, creating a new line below
func (b *Buffer) SplitLine(line, col int) {
	if line < 0 || line >= len(b.lines) {
		return
	}
	content := b.lines[line]
	runes := []rune(content)

	if col < 0 {
		col = 0
	}
	if col > len(runes) {
		col = len(runes)
	}

	before := string(runes[:col])
	after := string(runes[col:])

	b.lines[line] = before
	b.InsertLine(line+1, after)
}

// JoinLines joins line n with line n+1
func (b *Buffer) JoinLines(n int) {
	if n < 0 || n >= len(b.lines)-1 {
		return
	}
	b.lines[n] = b.lines[n] + b.lines[n+1]
	b.lines = append(b.lines[:n+1], b.lines[n+2:]...)
}

// InsertAt inserts text at the specified position
func (b *Buffer) InsertAt(line, col int, text string) {
	if line < 0 || line >= len(b.lines) {
		return
	}

	content := b.lines[line]
	runes := []rune(content)

	if col < 0 {
		col = 0
	}
	if col > len(runes) {
		col = len(runes)
	}

	newContent := string(runes[:col]) + text + string(runes[col:])
	b.lines[line] = newContent
}

// DeleteAt deletes count runes starting at position
func (b *Buffer) DeleteAt(line, col, count int) string {
	if line < 0 || line >= len(b.lines) || count <= 0 {
		return ""
	}

	content := b.lines[line]
	runes := []rune(content)

	if col < 0 || col >= len(runes) {
		return ""
	}

	end := col + count
	if end > len(runes) {
		end = len(runes)
	}

	deleted := string(runes[col:end])
	newContent := string(runes[:col]) + string(runes[end:])
	b.lines[line] = newContent

	return deleted
}

// DeleteRange deletes text between two positions (inclusive start, exclusive end)
func (b *Buffer) DeleteRange(start, end Position) string {
	if start.Line == end.Line {
		// Same line deletion
		return b.DeleteAt(start.Line, start.Col, end.Col-start.Col)
	}

	// Multi-line deletion
	var deleted strings.Builder

	// Delete from start to end of first line
	firstLine := b.GetLine(start.Line)
	firstRunes := []rune(firstLine)
	if start.Col < len(firstRunes) {
		deleted.WriteString(string(firstRunes[start.Col:]))
	}
	deleted.WriteString("\n")

	// Delete middle lines
	for i := start.Line + 1; i < end.Line; i++ {
		deleted.WriteString(b.GetLine(start.Line + 1))
		deleted.WriteString("\n")
		b.DeleteLine(start.Line + 1)
	}

	// Handle last line
	lastLine := b.GetLine(start.Line + 1)
	lastRunes := []rune(lastLine)
	if end.Col <= len(lastRunes) {
		deleted.WriteString(string(lastRunes[:end.Col]))
	}

	// Join the remaining parts
	newFirst := string(firstRunes[:start.Col]) + string(lastRunes[end.Col:])
	b.SetLine(start.Line, newFirst)
	b.DeleteLine(start.Line + 1)

	return deleted.String()
}

// GetRange returns text between two positions
func (b *Buffer) GetRange(start, end Position) string {
	if start.Line == end.Line {
		line := b.GetLine(start.Line)
		runes := []rune(line)
		if start.Col >= len(runes) {
			return ""
		}
		endCol := end.Col
		if endCol > len(runes) {
			endCol = len(runes)
		}
		return string(runes[start.Col:endCol])
	}

	var result strings.Builder

	// First line
	firstLine := b.GetLine(start.Line)
	firstRunes := []rune(firstLine)
	if start.Col < len(firstRunes) {
		result.WriteString(string(firstRunes[start.Col:]))
	}
	result.WriteString("\n")

	// Middle lines
	for i := start.Line + 1; i < end.Line; i++ {
		result.WriteString(b.GetLine(i))
		result.WriteString("\n")
	}

	// Last line
	lastLine := b.GetLine(end.Line)
	lastRunes := []rune(lastLine)
	endCol := end.Col
	if endCol > len(lastRunes) {
		endCol = len(lastRunes)
	}
	result.WriteString(string(lastRunes[:endCol]))

	return result.String()
}

// String returns the full buffer as a string
func (b *Buffer) String() string {
	return strings.Join(b.lines, "\n")
}

// Clone creates a deep copy of the buffer
func (b *Buffer) Clone() *Buffer {
	newLines := make([]string, len(b.lines))
	copy(newLines, b.lines)
	return &Buffer{lines: newLines}
}

// RuneCount returns the number of runes in line n
func (b *Buffer) RuneCount(line int) int {
	if line < 0 || line >= len(b.lines) {
		return 0
	}
	return utf8.RuneCountInString(b.lines[line])
}

// LastCol returns the last valid column for line n (0 for empty line)
func (b *Buffer) LastCol(line int) int {
	count := b.RuneCount(line)
	if count == 0 {
		return 0
	}
	return count - 1
}
