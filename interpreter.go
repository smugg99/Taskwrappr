// interpreter.go
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

type InterpreterFlags struct {
	LastActionReturned bool
	LastActionSuccess  bool
	Depth 			   uint
}

type ScriptRunner struct {
	Script *Script
	Memory *MemoryMap
	flags  *InterpreterFlags
}

// Constructor for ScriptRunner
func NewScriptRunner(script *Script, memory *MemoryMap) *ScriptRunner {
	return &ScriptRunner{
		Script: script,
		Memory: memory,
		flags:  &InterpreterFlags{},
	}
}

// Parse action line
func (s *ScriptRunner) ParseActionLine(line string) (string, []interface{}, error) {
    match := regexp.MustCompile(ActionArgumentsPattern).FindStringSubmatch(line)
    if len(match) != 3 {
        return "", nil, fmt.Errorf("invalid function call format: %s", line)
    }

    actionName := match[1]
    argString := match[2]

    var parsedArgs []interface{}

    parseArg := func(arg string) (interface{}, error) {
        if nestedMatch := regexp.MustCompile(ActionArgumentsPattern).FindStringSubmatch(arg); len(nestedMatch) == 3 {
            nestedActionName, nestedActionArgs, err := s.ParseActionLine(nestedMatch[0])
            if err != nil {
                return nil, fmt.Errorf("error parsing nested action '%s': %v", arg, err)
            }

            nestedAction := s.createActionCall(nestedActionName, nestedActionArgs)
            result, err := nestedAction()
            if err != nil {
                return nil, err
            }

            return result, nil
        } else {
            if strings.HasPrefix(arg, string(StringSymbol)) && strings.HasSuffix(arg, string(StringSymbol)) {
                return strings.Trim(arg, string(StringSymbol)), nil
            } else if intValue, err := strconv.Atoi(arg); err == nil {
                return intValue, nil
            } else if floatValue, err := strconv.ParseFloat(arg, 64); err == nil {
                return floatValue, nil
            } else if boolValue, err := strconv.ParseBool(arg); err == nil {
                return boolValue, nil
            } else if variable, ok := s.Memory.Variables[arg]; ok {
                return variable.Value, nil
            } else {
                return nil, fmt.Errorf("invalid argument type: %s", arg)
            }
        }
    }

    rawArgs := splitTopLevelArgs(argString)
    for _, arg := range rawArgs {
        arg = strings.TrimSpace(arg)
        if arg == "" {
            continue
        }

        parsedArg, err := parseArg(arg)
        if err != nil {
            return "", nil, err
        }
        parsedArgs = append(parsedArgs, parsedArg)
    }

    return actionName, parsedArgs, nil
}

// Execute action line
func (s *ScriptRunner) ExecuteActionLine(line string) (interface{}, error) {
	actionName, actionArgs, err := s.ParseActionLine(line)
	if err != nil {
		return nil, err
	}

	action := s.createActionCall(actionName, actionArgs)
	result, err := action()
	if err != nil {
		return result, fmt.Errorf("error executing action '%s': %v", actionName, err)
	}
	return result, nil
}

// Execute line
func (s *ScriptRunner) ExecuteLine(line string, previousToken Token) (Token, error) {
	token, _ := analyzeLine(line)

	switch token {
	case CodeBlockOpenToken:
		if previousToken == ActionToken && s.flags.LastActionReturned {
			if !s.flags.LastActionSuccess {
				return IgnoreToken, nil
			}
		} else {
			return IgnoreToken, nil
		}
		return token, nil
	case CodeBlockCloseToken:
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

// Run the script
func (s *ScriptRunner) Run() (bool, error) {
	cleanedContent, err := normalizeContent(s.Script.Content)
	if err != nil {
		return false, err
	}
	s.Script.CleanedContent = cleanedContent

	var previousToken Token = NoToken
	lines := strings.Split(s.Script.CleanedContent, string(NewLineSymbol))
	for i := 0; i < len(lines); i++ {
		token, err := s.ExecuteLine(lines[i], previousToken)
		if err != nil {
			return false, err
		}

		if token == IgnoreToken {
			jumpIndex, err := s.jumpToCodeBlock(lines, i, true)
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
