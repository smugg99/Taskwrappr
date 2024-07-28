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
	TabSymbol		     = '\t'
	ReturnSymbol         = '\r'
	EscapeSymbol		 = '\\'
	StringSymbol         = '"'
	SpaceSymbol          = ' '
	CommentSymbol        = '#'
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	AssignmentSymbol     = '='
	InvalidTokenSymbol   = -1
)

const (
	ParenOpenSymbol      = '('
	ParenCloseSymbol     = ')'
	BracketOpenSymbol    = '['
	BracketCloseSymbol   = ']'
	DelimiterSymbol      = ','
	DecimalSymbol        = '.'
	AdditionSymbol       = '+'
	SubtractionSymbol    = '-'
	MultiplicationSymbol = '*'
	DivisionSymbol       = '/'
	ModulusSymbol        = '%'
	ExponentSymbol       = '^'
	SelfReferenceSymbol  = '~'
	DeclarationSymbol	 = ':'
)

const (
	ActionToken TokenType = iota
	AssignmentToken
	DeclarationToken
	AugmentedAssignmentToken
	VariableToken
	LiteralToken
	CodeBlockOpenToken
	CodeBlockCloseToken
	OperatorAddToken
	OperatorSubtractToken
	OperatorUnaryMinusToken
	OperatorExponentToken
	OperatorMultiplyToken
	OperatorDivideToken
	OperatorModuloToken
	ParenOpenToken
	ParenCloseToken
	DelimiterToken
	DecimalToken
	LogicalAndToken
    LogicalOrToken
    LogicalNotToken
    LogicalXorToken
    EqualityToken
    InequalityToken
    LessThanToken
    LessThanOrEqualToken
    GreaterThanToken
    GreaterThanOrEqualToken
	InvalidToken
	IgnoreToken
	NoToken TokenType = InvalidTokenSymbol
)

const (
	TrueString                    = "true"
	FalseString                   = "false"
	NilString                     = "nil"
	LogicalAndString              = "&&"
	LogicalOrString               = "||"
	LogicalNotString              = "!"
	LogicalXorString              = "^^"
	EqualityString                = "=="
	InequalityString              = "!="
	LessThanString                = "<"
	LessThanOrEqualString         = "<="
	GreaterThanString             = ">"
	GreaterThanOrEqualString      = ">="
	DeclarationString             = string(DeclarationSymbol) + string(AssignmentSymbol)
	AugmentedAdditionString       = string(AdditionSymbol) + string(AssignmentSymbol)
	AugmentedSubtractionString    = string(SubtractionSymbol) + string(AssignmentSymbol)
	AugmentedMultiplicationString = string(MultiplicationSymbol) + string(AssignmentSymbol)
	AugmentedDivisionString       = string(DivisionSymbol) + string(AssignmentSymbol)
	AugmentedModulusString        = string(ModulusSymbol) + string(AssignmentSymbol)
	AugmentedExponentString       = string(ExponentSymbol) + string(AssignmentSymbol)
)

var (
	ActionCallPattern             = regexp.MustCompile(fmt.Sprintf(`\w+\%c[^%c]*\%c`, ParenOpenSymbol, ParenCloseSymbol, ParenCloseSymbol))
	ActionArgumentsPattern        = regexp.MustCompile(fmt.Sprintf(`^(\w+)\%c(.*)\%c$`, ParenOpenSymbol, ParenCloseSymbol))
	AssignmentPattern             = regexp.MustCompile(fmt.Sprintf(`^\s*([a-zA-Z_]\w*)\s*%c\s*(.+)\s*$`, AssignmentSymbol))
	DeclarationPattern 		      = regexp.MustCompile(fmt.Sprintf(`^\s*([a-zA-Z_]\w*)\s*%s\s*(.+)\s*$`, DeclarationString))
	LegalNamePattern              = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	VariableNamePattern           = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z_][a-zA-Z0-9_]*[^%c]$`, ParenOpenSymbol))
	IntegerPattern                = regexp.MustCompile(`^-?\d+$`)
	FloatPattern                  = regexp.MustCompile(`^-?\d*\.\d+$`)
	BooleanPattern                = regexp.MustCompile(fmt.Sprintf(`^(%s|%s)$`, TrueString, FalseString))
	StringPattern                 = regexp.MustCompile(fmt.Sprintf(`^%c.*%c$`, StringSymbol, StringSymbol))
	AugmentedAssignementPattern   = regexp.MustCompile(fmt.Sprintf(
		`^\s*(\w+)\s*([\%c\%c\%c\%c\%c\%c]%c)\s*(.*)\s*$`,
		AdditionSymbol, SubtractionSymbol, MultiplicationSymbol, DivisionSymbol, ModulusSymbol, ExponentSymbol, AssignmentSymbol,
	))
	LogicalOperatorsPattern       = regexp.MustCompile(fmt.Sprintf(
		`^(%s|%s|%s|%s|%s|%s|%s|%s|%s|%s)$`,
		EqualityString, InequalityString, LessThanString, LessThanOrEqualString, GreaterThanString, GreaterThanOrEqualString, LogicalAndString, LogicalOrString, LogicalNotString, LogicalXorString,
	))
)

var Operators = string([]rune{
	AdditionSymbol,
	SubtractionSymbol,
	MultiplicationSymbol,
	DivisionSymbol,
	ModulusSymbol,
	ExponentSymbol,
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
	case DeclarationToken:
		return "DeclarationToken"
	case AugmentedAssignmentToken:
		return "AugmentedAssignmentToken"
	case VariableToken:
		return "VariableToken"
	case LiteralToken:
		return "LiteralToken"
	case ParenOpenToken:
		return "ParenOpenToken"
	case ParenCloseToken:
		return "ParenCloseToken"
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
	case LogicalAndToken:
        return "LogicalAndToken"
    case LogicalOrToken:
        return "LogicalOrToken"
    case LogicalNotToken:
        return "LogicalNotToken"
    case LogicalXorToken:
        return "LogicalXorToken"
    case EqualityToken:
        return "EqualityToken"
    case InequalityToken:
        return "InequalityToken"
    case LessThanToken:
        return "LessThanToken"
    case LessThanOrEqualToken:
        return "LessThanOrEqualToken"
    case GreaterThanToken:
        return "GreaterThanToken"
    case GreaterThanOrEqualToken:
        return "GreaterThanOrEqualToken"
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