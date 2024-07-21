// helpers.go
package taskwrappr

import (
	"fmt"
	"strings"
)

func splitTopLevelArgs(argsString string) []string {
	var args []string
	var currentArg strings.Builder
	openBrackets := 0
	inQuotes := false
	escaped := false

	for _, ch := range argsString {
		switch ch {
		case ParenOpenSymbol:
			if !inQuotes && !escaped {
				openBrackets++
			}
		case ParenCloseSymbol:
			if !inQuotes && !escaped {
				openBrackets--
			}
		case StringSymbol:
			if !escaped {
				inQuotes = !inQuotes
			}
		case DelimiterSymbol:
			if openBrackets == 0 && !inQuotes && !escaped {
				args = append(args, strings.TrimSpace(currentArg.String()))
				currentArg.Reset()
				continue
			}
		}

		currentArg.WriteRune(ch)

		if ch == EscapeSymbol {
			escaped = !escaped
		} else {
			escaped = false
		}
	}

	if currentArg.Len() > 0 {
		args = append(args, strings.TrimSpace(currentArg.String()))
	}

	return args
}

func analyzeTopLevelLine(line string) (*Token, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return NewToken(IgnoreToken, line), nil
	}

	switch line[0] {
	case CommentSymbol:
		return NewToken(IgnoreToken, line), nil
	case CodeBlockOpenSymbol:
		return NewToken(CodeBlockOpenToken, line), nil
	case CodeBlockCloseSymbol:
		return NewToken(CodeBlockCloseToken, line), nil
	}

	if ActionCallPattern.MatchString(line) {
		return NewToken(ActionToken, line), nil
	}

	if AssignmentPattern.MatchString(line) {
		return NewToken(AssignmentToken, line), nil
	}

	// Other tokens like +=, -=, etc.

	return NewToken(InvalidToken, line), fmt.Errorf("invalid line: %s", line)
}

func (s *Script) parseActionToken(token *Token) (*Action, error) {
	match := ActionArgumentsPattern.FindStringSubmatch(token.Value)
	if len(match) != 3 {
		return nil, fmt.Errorf("invalid action call format: %s", token.Value)
	}

	actionName := match[1]
	argsString := match[2]

	actionFound := s.Memory.GetAction(actionName)
	if actionFound == nil {
		return nil, fmt.Errorf("undefined action: %s", actionName)
	}
	action := NewAction(actionFound.ExecuteFunc, actionFound.ValidateFunc)

	var parsedArgs []interface{}

	rawArgs := splitTopLevelArgs(argsString)

	for _, arg := range rawArgs {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}

		parsedArg, err := s.parseExpression(arg)
		if err != nil {
			return nil, err
		}
		
		parsedArgs = append(parsedArgs, parsedArg)
	}

	action.Arguments = parsedArgs

	return action, nil
}

func (s *Script) parseAssignmentToken(token *Token) (*Action, error) {
	return nil, nil
}

func (s *Script) parseContent() (*Block, error) {
	lines := strings.Split(s.CleanedContent, string(NewLineSymbol))

	blockStack := []*Block{}
	currentBlock := NewBlock()

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		token, err := analyzeTopLevelLine(line)
		if err != nil {
			return nil, fmt.Errorf("error analyzing line %d: %v", i+1, err)
		}

		switch token.Type {
		case ActionToken:
			action, err := s.parseActionToken(token)
			if err != nil {
				return nil, fmt.Errorf("error parsing action on line %d: %v", i+1, err)
			}

			currentBlock.Actions = append(currentBlock.Actions, action)
		case AssignmentToken:
			action, err := s.parseAssignmentToken(token)
			if err != nil {
				return nil, fmt.Errorf("error parsing assignment on line %d: %v", i+1, err)
			}

			currentBlock.Actions = append(currentBlock.Actions, action)
		case CodeBlockOpenToken:
			newBlock := NewBlock()
			blockStack = append(blockStack, currentBlock)
			currentBlock = newBlock
		case CodeBlockCloseToken:
			if len(blockStack) == 0 {
				return nil, fmt.Errorf("unmatched closing brace on line %d", i+1)
			}

			completedBlock := currentBlock
			currentBlock = blockStack[len(blockStack)-1]
			blockStack = blockStack[:len(blockStack)-1]

			if len(currentBlock.Actions) > 0 {
				lastAction := currentBlock.Actions[len(currentBlock.Actions)-1]
				lastAction.Block = completedBlock
			} else {
				return nil, fmt.Errorf("code block without preceding action on line %d", i+1)
			}
		}
	}

	if len(blockStack) > 0 {
		return nil, fmt.Errorf("unmatched opening brace")
	}

	return currentBlock, nil
}

func (s *Script) runBlock(b *Block) error {
	s.CurrentBlock = b
	for _, action := range b.Actions {
		if err := action.Validate(s); err != nil {
			return err
		}

		result, err := action.Execute(s)
		if err != nil {
			return err
		}
		b.LastResult = result

		if action.Block != nil {
			if resultBool, ok := result.(bool); ok {
				if resultBool {
					if err := s.runBlock(action.Block); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (s *Script) normalizeContent() (string, error) {
	var result strings.Builder
	inQuotes := false
	escaped := false
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
			case EscapeSymbol:
				if inQuotes && !escaped {
					escaped = true
					result.WriteByte(line[i])
				} else {
					result.WriteByte(line[i])
					escaped = false
				}
			case StringSymbol:
				result.WriteByte(line[i])
				if !escaped {
					inQuotes = !inQuotes
				}
			case CommentSymbol:
				if inQuotes || escaped {
					result.WriteByte(line[i])
				}
			case SpaceSymbol, TabSymbol, ReturnSymbol:
				if inQuotes || escaped {
					result.WriteByte(line[i])
				}
				escaped = false
			case CodeBlockOpenSymbol:
				if !inQuotes && !escaped {
					openCurlyCount++
					result.WriteByte(NewLineSymbol)
					result.WriteByte(line[i])
					result.WriteByte(NewLineSymbol)
				} else {
					result.WriteByte(line[i])
				}
				escaped = false
			case CodeBlockCloseSymbol:
				if !inQuotes && !escaped {
					openCurlyCount--
					result.WriteByte(NewLineSymbol)
					result.WriteByte(line[i])
					result.WriteByte(NewLineSymbol)
				} else {
					result.WriteByte(line[i])
				}
				escaped = false
			case ParenOpenSymbol:
				if !inQuotes && !escaped {
					openParenCount++
					result.WriteByte(line[i])
				} else {
					result.WriteByte(line[i])
				}
				escaped = false
			case ParenCloseSymbol:
				if !inQuotes && !escaped {
					openParenCount--
					result.WriteByte(line[i])
				} else {
					result.WriteByte(line[i])
				}
				escaped = false
			default:
				result.WriteByte(line[i])
				escaped = false
			}
			if line[i] == CommentSymbol && !inQuotes && !escaped {
				break
			}
		}
		if result.Len() > 0 && result.String()[result.Len()-1] != NewLineSymbol {
			result.WriteByte(NewLineSymbol)
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
			cleanedResult.WriteByte(NewLineSymbol)
		}
	}

	return strings.TrimSpace(cleanedResult.String()), nil
}
