package taskwrappr

import "testing"

func TestParserPeekRune(t *testing.T) {
	tok := NewTokenizer("hello")
	tokens, err := tok.Tokenize()
	if err != nil {
		t.Errorf("Error tokenizing: %v", err)
	}

	par := NewParser(tokens, "hello")
	if peek := par.peekToken(0); peek.Kind() != TokenIdentifier {
		t.Errorf("Expected identifier 'hello' got %v", peek)
	}
}
