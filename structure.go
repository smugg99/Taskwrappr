// structure.go
package taskwrappr

import (
    "os"
)

type Token struct {
	Type  TokenType
	Value string
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

type Script struct {
    Path         string
    Content      string
    MainBlock    *Block
    CurrentBlock *Block
}

func NewScript(filePath string, memory *MemoryMap) (*Script, error) {
    mainBlock := NewBlock(memory)
    return &Script{
        Path:         filePath,
        MainBlock:    mainBlock,
        CurrentBlock: mainBlock,
    }, nil
}

func (s *Script) Run() error {
    content, err := os.ReadFile(s.Path)
    if err != nil {
        return err
    }

    cleanedContent, err := normalizeContent(string(content))
    if err != nil {
        return err
    }
    s.Content = cleanedContent

    mainBlock, err := s.parseContent()
    if err != nil {
        return err
    }
    s.MainBlock = mainBlock

    if err := s.runBlock(s.MainBlock); err != nil {
        return err
    }

    return nil
}

type Block struct {
    Actions    []*Action
	Executed   bool
    Memory     *MemoryMap
	LastResult *Variable
}

func NewBlock(parentMemory *MemoryMap) *Block {
    mem := NewMemoryMap(parentMemory)
    return &Block{
        Actions: []*Action{},
        Memory:  mem,
    }
}