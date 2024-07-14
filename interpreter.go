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
	NoToken Token = UndefinedTokenSymbol
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
	MainBlock	   *Block
	CurrentBlock   *Block
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