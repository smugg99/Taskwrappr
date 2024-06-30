package taskwrappr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	NewLineSymbol        = '\n'
	StringSymbol         = '"'
	CommentSymbol        = '#'
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	ActionOpenSymbol     = '('
	ActionCloseSymbol    = ')'
	IfStatementSymbol    = "if"
	ElseStatementSymbol  = "else"
	UndefinedTokenSymbol = -1
)

type Token int

const (
	ActionToken         Token = iota
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
	ActionArgumentsPattern = `(\w+)\(([^()]*)\)`
	NestedActionPattern    = `(\w+)\((.*)\)`
)

var (
	IllegalNames = []string{IfStatementSymbol, ElseStatementSymbol}
)

type InterpreterFlags struct {
	LastActionReturned bool
	LastActionSuccess  bool
	Depth              uint
}

type ScriptRunner struct {
	Script *Script
	Memory *MemoryMap
	flags  *InterpreterFlags
}

func NewScriptRunner(script *Script, memory *MemoryMap) *ScriptRunner {
	return &ScriptRunner{
		Script: script,
		Memory: memory,
		flags:  &InterpreterFlags{},
	}
}

func (s *ScriptRunner) createActionCall(actionName string, args []interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		if action, ok := s.Memory.Actions[actionName]; ok {
			return action.Execute(s, args...)
		}
		return nil, fmt.Errorf("unknown action '%s'", actionName)
	}
}

func (s *ScriptRunner) NormalizeContent() error {
	var result strings.Builder
	inQuotes := false
	lines := strings.Split(s.Script.Content, string(NewLineSymbol))

	openCurlyCount := 0
	openParenCount := 0

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if !inQuotes && (strings.HasPrefix(trimmedLine, string(CommentSymbol)) || trimmedLine == "") {
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
			case CodeBlockOpenSymbol:
				if !inQuotes {
					openCurlyCount++
					result.WriteByte('\n')
					result.WriteByte(line[i])
					result.WriteByte('\n')
				} else {
					result.WriteByte(line[i])
				}
			case CodeBlockCloseSymbol:
				if !inQuotes {
					openCurlyCount--
					result.WriteByte('\n')
					result.WriteByte(line[i])
					result.WriteByte('\n')
				} else {
					result.WriteByte(line[i])
				}
			case ActionOpenSymbol:
				if !inQuotes {
					openParenCount++
					result.WriteByte(line[i])
				} else {
					result.WriteByte(line[i])
				}
			case ActionCloseSymbol:
				if !inQuotes {
					openParenCount--
					result.WriteByte(line[i])
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

	if inQuotes {
		return fmt.Errorf("unclosed string literal")
	}

	if openCurlyCount != 0 {
		return fmt.Errorf("unbalanced curly braces")
	}

	if openParenCount != 0 {
		return fmt.Errorf("unbalanced parentheses")
	}

	cleaned := strings.TrimSpace(result.String())
	lines = strings.Split(cleaned, string(NewLineSymbol))
	cleanedResult := strings.Builder{}

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" {
			cleanedResult.WriteString(trimmedLine)
			cleanedResult.WriteByte('\n')
		}
	}

	s.Script.CleanedContent = strings.TrimSpace(cleanedResult.String())

	return nil
}

func (s *ScriptRunner) AnalyzeLine(line string) (Token, Token) {
	if strings.HasPrefix(line, string(CodeBlockOpenSymbol)) {
		return CodeBlockOpenToken, NoToken
	}

	if strings.HasPrefix(line, string(CodeBlockCloseSymbol)) {
		return CodeBlockCloseToken, NoToken
	}

	actionCallPattern := regexp.MustCompile(ActionCallPattern)
	if actionCallPattern.MatchString(line) {
		return ActionToken, NoToken
	}

	return UndefinedToken, NoToken
}

func (s *ScriptRunner) ParseActionLine(line string) (string, []interface{}, error) {
	match := regexp.MustCompile(NestedActionPattern).FindStringSubmatch(line)
	if len(match) != 3 {
		return "", nil, fmt.Errorf("invalid function call format")
	}

	actionName := match[1]
	argString := match[2]

	rawArgs := strings.Split(argString, ",")

	for i := range rawArgs {
		rawArgs[i] = strings.TrimSpace(rawArgs[i])
	}

	var parsedArgs []interface{}
	for _, arg := range rawArgs {
		if arg == "" {
			continue
		}

		// Check if the argument is a nested action call
		if nestedMatch := regexp.MustCompile(NestedActionPattern).FindStringSubmatch(arg); len(nestedMatch) == 3 {
			nestedActionName, nestedActionArgs, err := s.ParseActionLine(arg)
			if err != nil {
				return "", nil, fmt.Errorf("error parsing nested action '%s': %v", arg, err)
			}
			nestedAction := s.createActionCall(nestedActionName, nestedActionArgs)
			result, err := nestedAction()
			if err != nil {
				return "", nil, err
			}
			parsedArgs = append(parsedArgs, result)
			continue
		}

		// Check if the argument is a string
		if strings.HasPrefix(arg, string(StringSymbol)) && strings.HasSuffix(arg, string(StringSymbol)) {
			parsedArgs = append(parsedArgs, strings.Trim(arg, string(StringSymbol)))
			continue
		}

		// Try to parse as integer
		if intValue, err := strconv.Atoi(arg); err == nil {
			parsedArgs = append(parsedArgs, intValue)
			continue
		}

		// Try to parse as float
		if floatValue, err := strconv.ParseFloat(arg, 64); err == nil {
			parsedArgs = append(parsedArgs, floatValue)
			continue
		}

		// Try to parse as boolean
		if boolValue, err := strconv.ParseBool(arg); err == nil {
			parsedArgs = append(parsedArgs, boolValue)
			continue
		}

		// Check if the argument is a variable
		if variable, ok := s.Memory.Variables[arg]; ok {
			parsedArgs = append(parsedArgs, variable.Value)
			continue
		}

		return "", nil, fmt.Errorf("invalid argument type: %s", arg)
	}

	return actionName, parsedArgs, nil
}

func (s *ScriptRunner) ExecuteActionLine(line string) (interface{}, error) {
	actionName, actionArgs, err := s.ParseActionLine(line)
	if err != nil {
		return nil, err
	}

	action := s.createActionCall(actionName, actionArgs)
	if result, err := action(); err != nil {
		return result, fmt.Errorf("error executing action '%s': %v", actionName, err)
	} else {
		return result, nil
	}
}

func (s *ScriptRunner) ExecuteLine(line string, previousToken Token) (Token, error) {
	token, _ := s.AnalyzeLine(line)
	// fmt.Printf("%s : %s --> %s, %d / %d\n", line, previousToken, token, s.flags.Depth, s.flags.MaxDepth)

	switch token {
	case CodeBlockOpenToken:
		if previousToken == ActionToken && s.flags.LastActionReturned {
			if s.flags.LastActionSuccess {
				s.flags.Depth++
			} else {
				return IgnoreToken, nil
			}
		} else {
			return IgnoreToken, nil
		}
		
		return token, nil
	case CodeBlockCloseToken:
		s.flags.Depth--
		return token, nil
	case ActionToken:
		result, err := s.ExecuteActionLine(line)

		if boolResult, ok := result.(bool); ok {
			s.flags.LastActionReturned = true
			s.flags.LastActionSuccess = boolResult
		} else {
			s.flags.LastActionReturned = false
		}

		return token, err
	case UndefinedToken:
		return NoToken, fmt.Errorf("unknown token in line: %s", line)
	}

	return NoToken, nil
}

func (s *ScriptRunner) JumpToCodeBlock(lines []string, startIndex int, forward bool) (int, error) {
	nestingLevel := 0
	step := 1
	if !forward {
		step = -1
	}

	for i := startIndex; i >= 0 && i < len(lines); i += step {
		token, _ := s.AnalyzeLine(lines[i])
		switch token {
		case CodeBlockOpenToken:
			nestingLevel++
		case CodeBlockCloseToken:
			nestingLevel--
			if nestingLevel == 0 {
				return i, nil
			}
		}
	}

	return -1, fmt.Errorf("matching code block not found")
}

func (s *ScriptRunner) Run() (bool, error) {
	if err := s.NormalizeContent(); err != nil {
		return false, err
	}

	var previousToken Token = NoToken

	lines := strings.Split(s.Script.CleanedContent, string(NewLineSymbol))
	for i := 0; i < len(lines); i++ {
		token, err := s.ExecuteLine(lines[i], previousToken)
		if err != nil {
			return false, err
		}

		if token == IgnoreToken {
			var jumpIndex int
			jumpIndex, err = s.JumpToCodeBlock(lines, i, true)
			if err != nil {
				return false, err
			}
			i = jumpIndex
			previousToken = CodeBlockCloseToken
		} else {
			previousToken = token
		}
	}

	return true, nil
}
