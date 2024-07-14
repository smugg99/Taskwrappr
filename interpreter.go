// interpreter.go
package taskwrappr

import (
	"os"
)

const (
	NewLineSymbol        = '\n'
	StringSymbol         = '"'
	CommentSymbol        = '#'
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	ActionOpenSymbol     = '('
	ActionCloseSymbol    = ')'
	ActionArgumentDelim  = ','
	UndefinedTokenSymbol = -1
)

type Token int

const (
	ActionToken Token = iota
	CodeBlockOpenToken
	CodeBlockCloseToken
	UndefinedToken
	IgnoreToken
	NoToken = UndefinedTokenSymbol
)

func (t Token) String() string {
	switch t {
	case ActionToken:
		return "ActionToken"
	case CodeBlockOpenToken:
		return "CodeBlockOpenToken"
	case CodeBlockCloseToken:
		return "CodeBlockCloseToken"
	case UndefinedToken:
		return "UndefinedToken"
	case NoToken:
		return "NoToken"
	case IgnoreToken:
		return "IgnoreToken"
	default:
		return "UnknownToken"
	}
}

const (
	ActionCallPattern      = `\w+\([^()]*\)`
	ActionArgumentsPattern = `^(\w+)\((.*)\)$`
)

type Script struct {
	Path           string
	Content        string
	CleanedContent string
	Memory 		   *MemoryMap
	Block		   *Block
}

func NewScript(filePath string, memory *MemoryMap) (*Script, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return &Script{
		Path:           filePath,
		Content:        string(content),
		Memory:         memory,
	}, nil
}

func (s *Script) Run() (bool, error) {
	cleanedContent, err := s.normalizeContent()
    if err != nil {
        return false, err
	}
	s.CleanedContent = cleanedContent

	parsedContent, err := s.parseContent()
	if err != nil {
		return false, err
	}
	s.Block = parsedContent
	
	if err := s.runBlock(s.Block); err != nil {
		return false, err
	}

    return true, nil
}