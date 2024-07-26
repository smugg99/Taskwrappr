// helpers.go
package taskwrappr

import (
	"fmt"
	"math"
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

func tokenizeLine(line string) (*Token, error) {
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

	if AssignmentPattern.MatchString(line) {
		return NewToken(AssignmentToken, line), nil
	}

	if ActionCallPattern.MatchString(line) {
		return NewToken(ActionToken, line), nil
	}

	if AugmentedAssignementPattern.MatchString(line) {
		return NewToken(AugmentedAssignmentToken, line), nil
	}

	return NewToken(InvalidToken, line), fmt.Errorf("invalid line: %s", line)
}

func ensureArithmeticOperands(a, b *Variable) (*Variable, *Variable, error) {
    var err error
    switch a.Type {
    case IntegerType, FloatType:
    default:
        if a, err = castToFloat(a); err != nil {
            return nil, nil, err
        }
    }
    switch b.Type {
    case IntegerType, FloatType:
    default:
        if b, err = castToFloat(b); err != nil {
            return nil, nil, err
        }
    }
    return a, b, nil
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
	action := CloneAction(actionFound)

	var parsedArgs []*Action
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

	action.SetArguments(parsedArgs)

	return action, nil
}

func (s *Script) parseAssignmentToken(token *Token) (*Action, error) {
	match := AssignmentPattern.FindStringSubmatch(token.Value)
	if len(match) != 3 {
		return nil, fmt.Errorf("invalid assignment format: %s", token.Value)
	}

	varName := match[1]
	exprString := match[2]

	assignmentAction := func(s *Script, args ...*Variable) ([]*Variable, error) {
		parseExprAction, err := s.parseExpression(exprString)
		if err != nil {
			return nil, err
		}

		parsedExpr, err := parseExprAction.Execute(s)
		if err != nil {
			return nil, err
		}

		if len(parsedExpr) != 1 {
			return nil, fmt.Errorf("invalid assignment expression: %s", exprString)
		}

		exprVar := parsedExpr[0]
		variable := s.Memory.SetVariable(varName, exprVar.Value, exprVar.Type)

		return []*Variable{variable}, nil
	}

	return NewAction(assignmentAction, nil), nil
}

func (s *Script) parseAugmentedAssignmentToken(token *Token) (*Action, error) {
    match := AugmentedAssignementPattern.FindStringSubmatch(token.Value)
    if len(match) != 4 {
        return nil, fmt.Errorf("invalid augmented assignment format: %s", token.Value)
    }
    varName := match[1]
    augmentedOperator := match[2]
    exprString := match[3]
    
    variable := s.Memory.GetVariable(varName)
    if variable == nil {
        return nil, fmt.Errorf("undefined variable: %s", varName)
    }
    
    parseExprAction, err := s.parseExpression(exprString)
    if err != nil {
        return nil, err
    }
    
    assignmentAction := func(s *Script, args ...*Variable) ([]*Variable, error) {
        parsedExpr, err := parseExprAction.Execute(s)
        if err != nil {
            return nil, err
        }

		if len(parsedExpr) != 1 {
			return nil, fmt.Errorf("invalid assignment expression: %s", exprString)
		}

        var (
            result interface{}
            castA, castB float64
            castErr error
        )
        
		exprVar := parsedExpr[0]
        variable, exprVar, castErr = ensureArithmeticOperands(variable, exprVar)
        if castErr != nil {
            return nil, fmt.Errorf("type mismatch in augmented assignment: %v", castErr)
        }
        
        castA, castErr = variable.toFloat()
        if castErr != nil {
            return nil, fmt.Errorf("failed to cast %s to float: %v", variable.Type.String(), castErr)
        }
        
        castB, castErr = exprVar.toFloat()
        if castErr != nil {
            return nil, fmt.Errorf("failed to cast %s to float: %v", exprVar.Type.String(), castErr)
        }

        switch augmentedOperator {
        case AugmentedAdditionString:
            result = castA + castB
        case AugmentedSubtractionString:
            result = castA - castB
        case AugmentedMultiplicationString:
            result = castA * castB
        case AugmentedDivisionString:
            if castB == 0 {
                return nil, fmt.Errorf("division by zero")
            }
            result = castA / castB
        case AugmentedModulusString:
            result = math.Mod(castA, castB)
        default:
            return nil, fmt.Errorf("unsupported operator for augmented assignment: %s", augmentedOperator)
        }
        
        variable.Value = result
        variable.Type = FloatType

        return []*Variable{variable}, nil
    }
    
    action := NewAction(assignmentAction, nil)
    return action, nil
}

func (s *Script) parseContent() (*Block, error) {
	lines := strings.Split(s.CleanedContent, string(NewLineSymbol))

	blockStack := []*Block{}
	currentBlock := NewBlock()

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		token, err := tokenizeLine(line)
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
		case AugmentedAssignmentToken:
			action, err := s.parseAugmentedAssignmentToken(token)
			if err != nil {
				return nil, fmt.Errorf("error parsing augmented assignment on line %d: %v", i+1, err)
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

		var resultVar *Variable

		if len(result) > 0 {
			resultVar = result[0]
			b.LastResult = resultVar
		} else {
			b.LastResult = nil
		}

		if action.Block != nil && resultVar != nil {
			if len(result) != 1 {
				return fmt.Errorf("invalid result from action: %v", result)
			}
			if resultBool, err := resultVar.toBool(); err == nil {
				if resultBool {
					if err := s.runBlock(action.Block); err != nil {
						return err
					}
				}
			} else {
				return err
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
