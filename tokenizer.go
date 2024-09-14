// tokenizer.go
package taskwrappr

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Tokenizer struct {
	Source    string // Source code to tokenize
	Index     uint   // Current index in the source code
	Line      uint   // Current line number
	Rune      rune   // Current rune being processed
	InString  bool   // Whether the tokenizer is currently processing a string
	InComment bool   // Whether the tokenizer is currently processing a comment
}

func (t *Tokenizer) String() string {
	return fmt.Sprintf("Tokenizer{Index: %d, Line: %d, Rune: %v, InString: %v, InComment: %v}", t.Index, t.Line, string(t.Rune), t.InString, t.InComment)
}

func NewTokenizer(source string) *Tokenizer {
	tokenizer := &Tokenizer{
		Source:   source,
		Index:    0,
		Line:     1,
		Rune:     0,
		InString: false,
	}

	tokenizer.readRune()
	return tokenizer
}

func (t *Tokenizer) readRune() {
	if t.Index >= uint(len(t.Source)) {
		t.Rune = 0
		return
	}

	t.Rune = rune(t.Source[t.Index])
	if t.Rune == NewLineSymbol {
		t.Line++
	}

	t.Index++
}

func (t *Tokenizer) peekRune(x uint) rune {
    index := t.Index - 1
    count := uint(0)

    // Traverse the source, skipping spaces
    for index < uint(len(t.Source)) {
        r := rune(t.Source[index])

        // Ignore spaces
        if !unicode.IsSpace(r) {
            // If x == 0, return the current non-space rune
            if count == x {
                return r
            }
            count++
        }
        index++
    }

    return 0 // Return 0 if we run out of characters
}

func (t *Tokenizer) nextToken() (Token, error) {
	for unicode.IsSpace(t.Rune) || isSeparator(t.Rune) {
		t.readRune()
	}

	// Handle single-character tokens
	switch t.Rune {
	case CodeBlockOpenSymbol:
		token := BlockDelimiterToken{IsOpen: true, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case CodeBlockCloseSymbol:
		token := BlockDelimiterToken{IsOpen: false, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case DelimiterSymbol:
		token := IdentifierDelimiterToken{index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case ParenOpenSymbol:
		token := ExpressionDelimiterToken{IsOpen: true, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case ParenCloseSymbol:
		token := ExpressionDelimiterToken{IsOpen: false, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case BracketOpenSymbol:
		token := IndexingDelimiterToken{IsOpen: true, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case BracketCloseSymbol:
		token := IndexingDelimiterToken{IsOpen: false, index: t.Index - 1, line: t.Line}
		t.readRune()
		return token, nil
	case CommentSymbol:
		t.skipComment()
		return nil, nil
	}

	// Handle string literals
	if isStringStart(t.Rune) {
		return t.handleStringLiteral()
	}

	// Handle numbers, including negative and float literals
	if isNumberStart(t.peekRune(0), t.peekRune(1), t.peekRune(2)) {
		return t.handleNumberLiteral()
	}

	// Handle identifiers (letters, digits, or underscores only)
	if isIdentifierStart(t.Rune) {
		return t.handleIdentifier()
	}

	// Handle operators (single or multi-character)
	if isOperatorStart(t.Rune) {
		return t.handleOperator()
	}

	// Handle end of file
	if t.Rune == 0 {
		return EOFToken{index: t.Index, line: t.Line}, nil
	}

	// Unknown token
	return nil, fmt.Errorf("unexpected token: %v %d", string(t.Rune), t.Rune)
}

// Handles single and multi-character operators
func (t *Tokenizer) handleOperator() (Token, error) {
	startIndex := t.Index
	startLine := t.Line
	possibleOperator := ""
	longestOperator := ""

	for i := 0; i < MaxOperatorLength && !(unicode.IsSpace(t.Rune) || unicode.IsLetter(t.Rune) || unicode.IsDigit(t.Rune)) && t.Rune != 0; i++ {
		possibleOperator += string(t.Rune)
		if isOperator(possibleOperator) {
			longestOperator = possibleOperator
		}
		t.readRune()
	}

	if longestOperator != "" {
		return OperationToken{Value: strings.TrimSpace(longestOperator), index: startIndex, line: startLine}, nil
	}

	if len(possibleOperator) > 0 && isOperator(possibleOperator[:1]) {
		operator := strings.TrimSpace(possibleOperator[:1])
		return OperationToken{Value: operator, index: startIndex, line: startLine}, nil
	}

	return nil, fmt.Errorf("invalid operator: %v", possibleOperator)
}

// Handle string literals, including escape sequences
func (t *Tokenizer) handleStringLiteral() (Token, error) {
	startIndex := t.Index
	startLine := t.Line
	var value strings.Builder

	for {
		t.readRune()
		// Handle escape sequences
		if t.Rune == EscapeSymbol && t.peekRune(0) == QuoteSymbol {
			t.readRune() // Consume the escape
			value.WriteRune(QuoteSymbol)
			continue
		}
		// End of string
		if t.Rune == QuoteSymbol || t.Rune == 0 {
			break
		}
		value.WriteRune(t.Rune)
	}
	// Consume closing quote
	if t.Rune == QuoteSymbol {
		t.readRune()
	}
	return LiteralToken{Value: value.String(), Type: TypeString, index: startIndex, line: startLine}, nil
}

// Handle numeric literals, including negative numbers and floats
func (t *Tokenizer) handleNumberLiteral() (Token, error) {
    startIndex := t.Index - 1
    startLine := t.Line
    hasDecimal := false
    var value strings.Builder

    // Handle optional sign
    if t.Rune == SubtractionSymbol {
        value.WriteRune(t.Rune)
        t.readRune()
    }

    // Collect digits and decimal point
    for unicode.IsDigit(t.Rune) || (!hasDecimal && t.Rune == DecimalSymbol) {
        if t.Rune == DecimalSymbol {
            hasDecimal = true
        }
        value.WriteRune(t.Rune)
        t.readRune()
    }

    // Skip trailing spaces after a number
    for unicode.IsSpace(t.Rune) {
        t.readRune()
    }

    // Attempt to parse the collected value as a float
    floatValue, err := strconv.ParseFloat(value.String(), 64)
    if err != nil {
        return nil, fmt.Errorf("invalid float literal: %v", value.String())
    }

    return LiteralToken{Value: floatValue, Type: TypeNumber, index: startIndex, line: startLine}, nil
}

// Handle identifiers, including reserved variable names
func (t *Tokenizer) handleIdentifier() (Token, error) {
    startIndex := t.Index
    startLine := t.Line
    var value strings.Builder

    // Collect valid identifier characters
    for unicode.IsLetter(t.Rune) || unicode.IsDigit(t.Rune) || t.Rune == UnderscoreSymbol {
        value.WriteRune(t.Rune)
		t.readRune()
    }

    idValue := strings.TrimSpace(value.String())
    // Check if the identifier is a reserved variable name
    if isReserved, varType := isReservedVariableName(idValue); isReserved {
        return LiteralToken{Value: idValue, Type: varType, index: startIndex, line: startLine}, nil
    }

    return IdentifierToken{Value: idValue, index: startIndex, line: startLine}, nil
}

// Skips over comments (until the end of the line or EOF)
func (t *Tokenizer) skipComment() {
	for t.Rune != NewLineSymbol && t.Rune != 0 {
		t.readRune()
	}
	// Consume newline, if any
	if t.Rune == NewLineSymbol {
		t.readRune()
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	var tokens []Token	
	for {
		token, err := t.nextToken()
		if err != nil {
			return nil, fmt.Errorf("[%d:%d] %v", t.Line, t.Index, err)
		}

		if token == nil {
			continue
		}

		if token.Kind() == TokenEOF {
			break
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}
