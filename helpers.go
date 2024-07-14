// helpers.go
package taskwrappr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func analyzeLine(line string) (Token, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return IgnoreToken, nil
	}

	switch line[0] {
	case CommentSymbol:
		return IgnoreToken, nil
	case CodeBlockOpenSymbol:
		return CodeBlockOpenToken, nil
	case CodeBlockCloseSymbol:
		return CodeBlockCloseToken, nil
	}

	if regexp.MustCompile(ActionCallPattern).MatchString(line) {
		return ActionToken, nil
	}

	return UndefinedToken, fmt.Errorf("invalid line: %s", line)
}

func splitTopLevelArgs(argString string) []string {
	var args []string
	var currentArg strings.Builder
	openBrackets := 0
	for _, ch := range argString {
		if ch == '(' {
			openBrackets++
		}
		if ch == ')' {
			openBrackets--
		}
		if ch == ActionArgumentDelim && openBrackets == 0 {
			args = append(args, strings.TrimSpace(currentArg.String()))
			currentArg.Reset()
		} else {
			currentArg.WriteRune(ch)
		}
	}
	if currentArg.Len() > 0 {
		args = append(args, strings.TrimSpace(currentArg.String()))
	}
	return args
}

func parseNonActionArg(argString string) (interface{}, error) {
	if strings.HasPrefix(argString, string(StringSymbol)) && strings.HasSuffix(argString, string(StringSymbol)) {
		return strings.Trim(argString, string(StringSymbol)), nil
	} else if intValue, err := strconv.Atoi(argString); err == nil {
		return intValue, nil
	} else if floatValue, err := strconv.ParseFloat(argString, 64); err == nil {
		return floatValue, nil
	} else if boolValue, err := strconv.ParseBool(argString); err == nil {
		return boolValue, nil
	} else {
		return nil, fmt.Errorf("invalid argument type: %s", argString)
	}
}

func (s *Script) runBlock(b *Block) error {
	for _, action := range b.Actions {
		
		_, err := action.Execute(s)
		if err != nil {
			return err
		}

		if action.Block != nil {
			if err := s.runBlock(action.Block); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Script) normalizeContent() (string, error) {
	var result strings.Builder
	inQuotes := false
	lines := strings.Split(s.Content, string(NewLineSymbol))

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
		return "", fmt.Errorf("unclosed string literal")
	}

	if openCurlyCount != 0 {
		return "", fmt.Errorf("unbalanced curly braces")
	}

	if openParenCount != 0 {
		return "", fmt.Errorf("unbalanced parentheses")
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

	return strings.TrimSpace(cleanedResult.String()), nil
}

// TODO: Add support for non-global memory map
func (s *Script) parseActionLine(line string) (*Action, error) {
	match := regexp.MustCompile(ActionArgumentsPattern).FindStringSubmatch(line)
    if len(match) != 3 {
        return nil, fmt.Errorf("invalid function call format: %s", line)
    }

    actionName := match[1]
    argString := match[2]

	actionFound, ok := s.Memory.Actions[actionName];
	if !ok {
		return nil, fmt.Errorf("undefined action: %s", actionName)
	}

	action := NewAction(actionFound.ExecuteFunc)

	var parsedArgs []interface{}
	parseArg := func(arg string) (interface{}, error) {
        if nestedMatch := regexp.MustCompile(ActionArgumentsPattern).FindStringSubmatch(arg); len(nestedMatch) == 3 {
			nestedAction, err := s.parseActionLine(nestedMatch[0])
            if err != nil {
                return nil, fmt.Errorf("error parsing nested action '%s': %v", arg, err)
            }

            return nestedAction, nil
        } else {
            return parseNonActionArg(arg)
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
            return nil, err
        }
        parsedArgs = append(parsedArgs, parsedArg)
    }

	action.Arguments = parsedArgs

    return action, nil
}

func (s *Script) parseContent() (*Block, error) {
	lines := strings.Split(s.CleanedContent, string(NewLineSymbol))

	actions := []*Action{}
	memory := NewMemoryMap()

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		token, err := analyzeLine(line)
		if err != nil {
			return nil, fmt.Errorf("error analyzing line %d: %v", i+1, err)
		}

		switch token {
		case ActionToken:
			action, err := s.parseActionLine(line)
			if err != nil {
				return nil, fmt.Errorf("error parsing action on line %d: %v", i+1, err)
			}

			actions = append(actions, action)
		case CodeBlockOpenToken:
		case CodeBlockCloseToken:
		}
	}

	return &Block{
		Actions: actions,
		Memory:  memory,
	}, nil
}