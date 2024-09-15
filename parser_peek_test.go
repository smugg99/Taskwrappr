package taskwrappr

import (
	"testing"
)

func TestParserPeekRune(t *testing.T) {
	filePath := "./scripts/parserTests.tw"
	tok := NewTokenizer(filePath)
	tokens, err := tok.Tokenize()
	if err != nil {
		t.Errorf("Error tokenizing: %v", err)
	}

	par := NewParser(tokens, filePath)
	if peek := par.peekToken(0); peek.Kind() != TokenIdentifier {
		t.Errorf("Expected identifier 'hello' got %v", peek)
	}

	if peek := par.peekToken(1); peek.Kind() != TokenOperation {
		t.Errorf("Expected operation ':=' got %v", peek)
	}

	if peek := par.peekToken(2); peek.Kind() != TokenLiteral {
		t.Errorf("Expected literal number 5 got %v", peek)
	}
}

func TestParserPeekUntilTokenKind(t *testing.T) {
	filePath := "./scripts/parserTests.tw"
	tok := NewTokenizer(filePath)
	tokens, err := tok.Tokenize()
	if err != nil {
		t.Errorf("Error tokenizing: %v", err)
	}

	if len(tokens) != 3 {
		t.Errorf("Expected 3 tokens got %v", tokens)
	}

	par := NewParser(tokens, filePath)
	ok, tokens := par.peekUntilTokenKind(TokenLiteral);
	if !ok {
		t.Errorf("Expected 1 token got %v", tokens)
	}

	lastIndex := len(tokens) - 1
	if tokens[lastIndex].Kind() != TokenLiteral {
		t.Errorf("Expected literal number 5 got %v", tokens[lastIndex])
	}

	// Test for non-existent token kind
	ok, tokens = par.peekUntilTokenKind(TokenIndexingDelimiter);
	if ok {
		t.Errorf("Expected 0 tokens got %v", tokens)
	}

	ok, tokens = par.peekUntilTokenKind(TokenOperation);
	if !ok {
		t.Errorf("Expected 1 token got %v", tokens)
	}

	lastIndex = len(tokens) - 1
	if tokens[lastIndex].Kind() != TokenOperation {
		t.Errorf("Expected operation ':=' got %v", tokens[lastIndex])
	}
}