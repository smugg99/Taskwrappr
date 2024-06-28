package taskwrappr

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	NewLineSymbol        = "\n"
	StringSymbol         = '"'
	CommentSymbol        = "#"
	VariableOpenSymbol   = "${{"
	VariableCloseSymbol  = "}}"
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	ActionOpenSymbol     = '('
	ActionCloseSymbol    = ')'
	IfStatementSymbol    = "if"
	ElseStatementSymbol  = "else"
)

type Token string

const (
	IfStatementToken    Token = "if"
	ElseStatementToken  Token = "else"
	ActionToken         Token = "()"
	CodeBlockOpenToken  Token = "{"
	CodeBlockCloseToken Token = "}"
	UndefinedToken      Token = "undefined"
)

const (
	ActionCallPattern = `\w+\([^)]*\)`
)

type ScriptRunner struct {
	Script *Script
	Memory *MemoryMap
}

func NewScriptRunner(script *Script, memory *MemoryMap) *ScriptRunner {
	return &ScriptRunner{
		Script: script,
		Memory: memory,
	}
}

func (s *ScriptRunner) NormalizeContent() {
	var result strings.Builder
	inQuotes := false
	lines := strings.Split(s.Script.Content, NewLineSymbol)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if !inQuotes && (strings.HasPrefix(trimmedLine, CommentSymbol) || trimmedLine == "") {
			continue
		}

		for i := 0; i < len(line); i++ {
			switch line[i] {
			case StringSymbol:
				result.WriteByte(line[i])
				inQuotes = !inQuotes
			case ' ', '\t':
				if inQuotes {
					result.WriteByte(line[i])
				}
			case CodeBlockOpenSymbol, CodeBlockCloseSymbol:
				if !inQuotes {
					result.WriteByte('\n')
					result.WriteByte(line[i])
					result.WriteByte('\n')
				} else {
					result.WriteByte(line[i])
				}
			default:
				result.WriteByte(line[i])
			}
		}
		if result.Len() > 0 && result.String()[result.Len()-1] != '\n' {
			result.WriteByte('\n')
		}
	}

	cleaned := strings.TrimSpace(result.String())
	lines = strings.Split(cleaned, NewLineSymbol)
	cleanedResult := strings.Builder{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			cleanedResult.WriteString(trimmedLine)
			cleanedResult.WriteByte('\n')
		}
	}

	s.Script.CleanedContent = strings.TrimSpace(cleanedResult.String())
}

func (s *ScriptRunner) AnalyzeLine(line string) Token {
	if strings.HasPrefix(line, IfStatementSymbol) {
		return IfStatementToken
	}

	if strings.HasPrefix(line, ElseStatementSymbol) {
		return ElseStatementToken
	}

	if strings.HasPrefix(line, string(CodeBlockOpenSymbol)) {
		return CodeBlockOpenToken
	}

	if strings.HasPrefix(line, string(CodeBlockCloseSymbol)) {
		return CodeBlockCloseToken
	}

	actionCallPattern := regexp.MustCompile(ActionCallPattern)
	if actionCallPattern.MatchString(line) {
		return ActionToken
	}

	return UndefinedToken
}

func (s *ScriptRunner) ExecuteLine(line string) error {
	token := s.AnalyzeLine(line)
	fmt.Printf("Token: %s", token)

	switch token {
	case IfStatementToken:
		// Execute if statement
	case ElseStatementToken:
		// Execute else statement
	case CodeBlockOpenToken:
		// Execute code block open
	case CodeBlockCloseToken:
		// Execute code block close
	case ActionToken:
		// Execute action
	case UndefinedToken:
		// Execute line as a simple statement
	}

	return nil
}

func (s *ScriptRunner) Run() (bool, error) {
	s.NormalizeContent()
	fmt.Println(s.Script.CleanedContent)

	lines := strings.Split(s.Script.CleanedContent, NewLineSymbol)
	for i := 0; i < len(lines); i++ {
		if err := s.ExecuteLine(lines[i]); err != nil {
			return false, err
		}
	}

	return true, nil
}
