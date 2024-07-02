// helpers.go
package taskwrappr

import (
	"fmt"
	"regexp"
	"strings"
)

func splitTopLevelArgs(argString string) []string {
	var args []string
	var currentArg strings.Builder
	depth := 0

	for _, r := range argString {
		if r == ActionOpenSymbol {
			depth++
		} else if r == ActionCloseSymbol {
			depth--
		} else if r == ActionArgumentDelim && depth == 0 {
			args = append(args, currentArg.String())
			currentArg.Reset()
			continue
		}
		currentArg.WriteRune(r)
	}
	args = append(args, currentArg.String())

	return args
}

func normalizeContent(contentString string) (string, error) {
	var result strings.Builder
	inQuotes := false
	lines := strings.Split(contentString, string(NewLineSymbol))

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

func analyzeLine(line string) (Token, Token) {
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

func (s *ScriptRunner) createActionCall(actionName string, args []interface{}) func() (interface{}, error) {
	return func() (interface{}, error) {
		if action, ok := s.Memory.Actions[actionName]; ok {
			return action.Execute(s, args...)
		}
		return nil, fmt.Errorf("unknown action '%s'", actionName)
	}
}

// Jump to matching code block
func (s *ScriptRunner) jumpToCodeBlock(lines []string, startIndex int, forward bool) (int, error) {
	nestingLevel := 0
	step := 1
	if !forward {
		step = -1
	}

	for i := startIndex; i >= 0 && i < len(lines); i += step {
		token, _ := analyzeLine(lines[i])
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