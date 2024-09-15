package taskwrappr

import "testing"

func TestTokenizerPeekRune(t *testing.T) {
	tok := NewTokenizer("./scripts/tokenizerTests.tw")
	if peek := tok.peekRune(0); peek != 'h' {
		t.Errorf("Expected 'h' got %c", peek)
	}
}

func TestTokenizerPeekRuneWithSpecialChars(t *testing.T) {
	tok := NewTokenizer("./scripts/tokenizerTests.tw")
	if peek := tok.peekRune(0); peek != 'h' {
		t.Errorf("Expected 'h' got %c", peek)
	}
	if peek := tok.peekRune(4); peek != 'o' {
		t.Errorf("Expected 'o' got %c", peek)
	}
	if peek := tok.peekRune(5); peek != '!' {
		t.Errorf("Expected '!' got %c", peek)
	}
	if peek := tok.peekRune(6); peek != '@' {
		t.Errorf("Expected '@' got %c", peek)
	}
	if peek := tok.peekRune(7); peek != '#' {
		t.Errorf("Expected '#' got %c", peek)
	}

	if peek := tok.peekRune(8); peek != 'w' {
		t.Errorf("Expected 'w' got %c", peek)
	}
	if peek := tok.peekRune(9); peek != 'o' {
		t.Errorf("Expected 'o' got %c", peek)
	}
}

func TestTokenizerPeekRuneWithSpaces(t *testing.T) {
	tok := NewTokenizer("./scripts/tokenizerTests.tw")
	if peek := tok.peekRune(0); peek != 'h' {
		t.Errorf("Expected 'h' got %c", peek)
	}
	if peek := tok.peekRune(3); peek != 'l' {
		t.Errorf("Expected 'l' got %c", peek)
	}
	if peek := tok.peekRune(4); peek != 'o' {
		t.Errorf("Expected 'o' got %c", peek)
	}
	if peek := tok.peekRune(5); peek != '!' {
		t.Errorf("Expected '!' got %c", peek)
	}

	if peek := tok.peekRune(8); peek != 'w' {
		t.Errorf("Expected 'w' got %c", peek)
	}
}

func TestTokenizerPeekRuneOutOfBounds(t *testing.T) {
	tok := NewTokenizer("./scripts/tokenizerTests.tw")
	if peek := tok.peekRune(100); peek != 0 {
		t.Errorf("Expected 0 got %c", peek)
	}
}
