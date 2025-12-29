package vim

import (
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer("hello\nworld")
	if buf.LineCount() != 2 {
		t.Errorf("expected 2 lines, got %d", buf.LineCount())
	}
	if buf.GetLine(0) != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.GetLine(0))
	}
	if buf.GetLine(1) != "world" {
		t.Errorf("expected 'world', got '%s'", buf.GetLine(1))
	}
}

func TestBufferInsertAt(t *testing.T) {
	buf := NewBuffer("hello world")
	buf.InsertAt(0, 5, " there")
	if buf.GetLine(0) != "hello there world" {
		t.Errorf("expected 'hello there world', got '%s'", buf.GetLine(0))
	}
}

func TestBufferDeleteAt(t *testing.T) {
	buf := NewBuffer("hello world")
	deleted := buf.DeleteAt(0, 5, 6)
	if deleted != " world" {
		t.Errorf("expected ' world', got '%s'", deleted)
	}
	if buf.GetLine(0) != "hello" {
		t.Errorf("expected 'hello', got '%s'", buf.GetLine(0))
	}
}

func TestNewEditor(t *testing.T) {
	e := NewEditor("hello world")
	if e.Mode != ModeNormal {
		t.Errorf("expected ModeNormal, got %v", e.Mode)
	}
	if e.Cursor.Line != 0 || e.Cursor.Col != 0 {
		t.Errorf("expected cursor at (0,0), got (%d,%d)", e.Cursor.Line, e.Cursor.Col)
	}
}

func TestMotionRight(t *testing.T) {
	e := NewEditor("hello")
	e.HandleKey("l")
	if e.Cursor.Col != 1 {
		t.Errorf("expected col 1, got %d", e.Cursor.Col)
	}
}

func TestMotionWordForward(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("w")
	if e.Cursor.Col != 6 {
		t.Errorf("expected col 6, got %d", e.Cursor.Col)
	}
}

func TestInsertMode(t *testing.T) {
	e := NewEditor("hello")
	e.HandleKey("i")
	if e.Mode != ModeInsert {
		t.Errorf("expected ModeInsert, got %v", e.Mode)
	}
	e.HandleKey("X")
	if e.Buffer.GetLine(0) != "Xhello" {
		t.Errorf("expected 'Xhello', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestDeleteWord(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("d")
	e.HandleKey("w")
	if e.Buffer.GetLine(0) != "world" {
		t.Errorf("expected 'world', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestDeleteLine(t *testing.T) {
	e := NewEditor("line1\nline2\nline3")
	e.HandleKey("d")
	e.HandleKey("d")
	if e.Buffer.LineCount() != 2 {
		t.Errorf("expected 2 lines, got %d", e.Buffer.LineCount())
	}
	if e.Buffer.GetLine(0) != "line2" {
		t.Errorf("expected 'line2', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestYankPaste(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("y")
	e.HandleKey("w")
	// Should have yanked "hello " (word + trailing space)
	if e.Unnamed != "hello " {
		t.Errorf("expected yanked 'hello ', got '%s'", e.Unnamed)
	}
	e.HandleKey("$")
	e.HandleKey("p")
	// Paste after last character
	if e.Buffer.GetLine(0) != "hello worldhello " {
		t.Errorf("expected 'hello worldhello ', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestUndo(t *testing.T) {
	e := NewEditor("hello")
	e.HandleKey("x")
	if e.Buffer.GetLine(0) != "ello" {
		t.Errorf("expected 'ello', got '%s'", e.Buffer.GetLine(0))
	}
	e.HandleKey("u")
	if e.Buffer.GetLine(0) != "hello" {
		t.Errorf("expected 'hello' after undo, got '%s'", e.Buffer.GetLine(0))
	}
}

func TestVisualMode(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("v")
	if e.Mode != ModeVisual {
		t.Errorf("expected ModeVisual, got %v", e.Mode)
	}
	e.HandleKey("e")
	e.HandleKey("d")
	if e.Buffer.GetLine(0) != " world" {
		t.Errorf("expected ' world', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestTextObjectWord(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("w") // Move to 'world'
	e.HandleKey("d")
	e.HandleKey("i")
	e.HandleKey("w")
	if e.Buffer.GetLine(0) != "hello " {
		t.Errorf("expected 'hello ', got '%s'", e.Buffer.GetLine(0))
	}
}

func TestTextObjectQuotes(t *testing.T) {
	e := NewEditor(`say "hello world"`)
	// Move into the quoted string
	e.HandleKey("f")
	e.HandleKey("h") // Move to 'h' in hello
	e.HandleKey("d")
	e.HandleKey("i")
	e.HandleKey("\"")
	if e.Buffer.GetLine(0) != `say ""` {
		t.Errorf(`expected 'say ""', got '%s'`, e.Buffer.GetLine(0))
	}
}

func TestCount(t *testing.T) {
	e := NewEditor("hello world test")
	e.HandleKey("2")
	e.HandleKey("w")
	if e.Cursor.Col != 12 {
		t.Errorf("expected col 12, got %d", e.Cursor.Col)
	}
}

func TestFindChar(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("f")
	e.HandleKey("o")
	if e.Cursor.Col != 4 {
		t.Errorf("expected col 4, got %d", e.Cursor.Col)
	}
}

func TestRenderState(t *testing.T) {
	e := NewEditor("hello\nworld")
	state := e.GetRenderState()
	if len(state.Lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(state.Lines))
	}
	if state.ModeString != "NORMAL" {
		t.Errorf("expected 'NORMAL', got '%s'", state.ModeString)
	}
}

func TestValidateCursorPosition(t *testing.T) {
	e := NewEditor("hello world")
	e.HandleKey("$") // Go to end

	spec := &ChallengeSpec{
		ValidationType: "cursor_position",
		ExpectedCursor: []int{0, 10},
	}

	result := Validate(e, spec)
	if !result.Success {
		t.Errorf("expected success, cursor at (%d,%d)", e.Cursor.Line, e.Cursor.Col)
	}
}

func TestValidateExactMatch(t *testing.T) {
	e := NewEditor("hello")
	e.HandleKey("A") // Append
	e.HandleKey(" ")
	e.HandleKey("w")
	e.HandleKey("o")
	e.HandleKey("r")
	e.HandleKey("l")
	e.HandleKey("d")
	e.HandleKey("Escape")

	spec := &ChallengeSpec{
		ValidationType: "exact_match",
		ExpectedBuffer: "hello world",
	}

	result := Validate(e, spec)
	if !result.Success {
		t.Errorf("expected success, buffer: '%s'", e.Buffer.String())
	}
}
