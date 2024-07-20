// interpreter.go
package taskwrappr

import (
	"fmt"
	"os"
	"regexp"
)

type TokenType    int

const (
	NewLineSymbol        = '\n'
	StringSymbol         = '"'
	SpaceSymbol          = ' '
	ControlSymbol		 = '\t'
	EscapeSymbol		 = '\\'
	CommentSymbol        = '#'
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	AssignmentSymbol     = '='
	InvalidTokenSymbol   = -1
)

const (
	ParenOpenSymbol      = '('
	ParenCloseSymbol     = ')'
	DelimiterSymbol      = ','
	DecimalSymbol        = '.'
	AdditionSymbol       = '+'
	SubtractionSymbol    = '-'
	MultiplicationSymbol = '*'
	DivisionSymbol       = '/'
	ModulusSymbol        = '%'
)

const (
	ActionToken TokenType = iota
	AssignmentToken
	VariableToken
	CodeBlockOpenToken
	CodeBlockCloseToken
	OperatorAddToken
	OperatorSubtractToken
	OperatorMultiplyToken
	OperatorDivideToken
	OperatorModuloToken
	ParenOpenToken
	ParenCloseToken
	DelimiterToken
	DecimalToken
	InvalidToken
	IgnoreToken
	NoToken TokenType = InvalidTokenSymbol
)

var (
	ActionCallPattern           = regexp.MustCompile(fmt.Sprintf(`\w+\%c[^%c]*\%c`, ParenOpenSymbol, ParenCloseSymbol, ParenCloseSymbol))
	ActionArgumentsPattern      = regexp.MustCompile(fmt.Sprintf(`^(\w+)\%c(.*)\%c$`, ParenOpenSymbol, ParenCloseSymbol))
	AssignmentPattern           = regexp.MustCompile(fmt.Sprintf(`^\s*\w+\s*\%c\s*.+\s*$`, AssignmentSymbol))
	AssignmentExpressionPattern = regexp.MustCompile(fmt.Sprintf(`^\s*(\w+)\s*\%c\s*(.+)\s*$`, AssignmentSymbol))
	LegalNamePattern            = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	VariableNamePattern         = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z_][a-zA-Z0-9_]*[^%c]$`, ParenOpenSymbol))
	IntegerPattern              = regexp.MustCompile(`^-?\d+$`)
	FloatPattern                = regexp.MustCompile(`^-?\d*\.\d+$`)
	BooleanPattern              = regexp.MustCompile(`^(true|false)$`)
	StringPattern               = regexp.MustCompile(`^".*"$`)
)

var Operators = string([]rune{
	AdditionSymbol,
	SubtractionSymbol,
	MultiplicationSymbol,
	DivisionSymbol,
	ModulusSymbol,
})

type Token struct {
	Type  TokenType
	Value string
}

type Script struct {
	Path           string
	Content        string
	CleanedContent string
	Memory 		   *MemoryMap
	MainBlock	   *Block
	CurrentBlock   *Block
}

func (t TokenType) String() string {
	switch t {
	case ActionToken:
		return "ActionToken"
	case CodeBlockOpenToken:
		return "CodeBlockOpenToken"
	case CodeBlockCloseToken:
		return "CodeBlockCloseToken"
	case AssignmentToken:
		return "AssignmentToken"
	case VariableToken:
		return "VariableToken"
	case OperatorAddToken:
		return "OperatorAddToken"
	case OperatorSubtractToken:
		return "OperatorSubtractToken"
	case OperatorMultiplyToken:
		return "OperatorMultiplyToken"
	case OperatorDivideToken:
		return "OperatorDivideToken"
	case OperatorModuloToken:
		return "OperatorModuloToken"
	case DelimiterToken:
		return "DelimiterToken"
	case DecimalToken:
		return "DecimalToken"
	case IgnoreToken:
		return "IgnoreToken"
	case NoToken:
		return "NoToken"
	default:
		return "InvalidToken"
	}
}

func NewToken(tokenType TokenType, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}

func NewScript(filePath string, memory *MemoryMap) (*Script, error) {
	return &Script{
		Path:           filePath,
		Memory:         memory,
	}, nil
}

func (s *Script) Run() (bool, error) {
	content, err := os.ReadFile(s.Path)
	if err != nil {
		return false, err
	}
	s.Content = string(content)

	cleanedContent, err := s.normalizeContent()
    if err != nil {
        return false, err
	}
	s.CleanedContent = cleanedContent

	parsedContent, err := s.parseContent()
	if err != nil {
		return false, err
	}
	s.MainBlock = parsedContent
	
	if err := s.runBlock(s.MainBlock); err != nil {
		return false, err
	}

    return true, nil
}