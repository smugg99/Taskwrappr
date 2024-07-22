// expressioner.go
package taskwrappr

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

func isPotentialNegNumber(runes []rune, i int) bool {
	return runes[i] == SubtractionSymbol && (i == 0 || isOperatorOrParen(runes[i-1])) && (i+1 < len(runes) && (unicode.IsDigit(runes[i+1]) || runes[i+1] == '.'))
}

func isOperator(r rune) bool {
	return strings.ContainsRune(Operators, r)
}

func isOperatorOrParen(r rune) bool {
	return isOperator(r) || r == ParenOpenSymbol || r == ParenCloseSymbol
}

func isLiteral(expr string) bool {
	switch {
	case IntegerPattern.MatchString(expr):
		return true
	case FloatPattern.MatchString(expr):
		return true
	case BooleanPattern.MatchString(expr):
		return true
	case StringPattern.MatchString(expr):
		return true
	case strings.HasPrefix(expr, string(SubtractionSymbol)) && len(expr) > 1 && unicode.IsDigit([]rune(expr)[1]):
		if IntegerPattern.MatchString(expr) || FloatPattern.MatchString(expr) {
			return true
		}
	default:
		return false
	}
	return false
}

func isVariable(expr string) bool {
	return VariableNamePattern.MatchString(expr)
}

func isAction(expr string) bool {
	return len(ActionArgumentsPattern.FindStringSubmatch(expr)) == 3
}

func getPrecedence(token *Token) int {
	switch token.Type {
	case OperatorAddToken, OperatorSubtractToken:
		return 1
	case OperatorMultiplyToken, OperatorDivideToken, OperatorModuloToken:
		return 2
	case OperatorUnaryMinusToken:
		return 3
	case OperatorExponentToken:
		return 4
	}
	return 0
}

func findClosingParen(runes []rune, start int) int {
	n := len(runes)
	open := 1
	for i := start + 1; i < n; i++ {
		switch runes[i] {
		case ParenOpenSymbol:
			open++
		case ParenCloseSymbol:
			open--
			if open == 0 {
				return i
			}
		}
	}
	return start
}

func parseLiteral(exprString string) (*Variable, error) {
	exprString = strings.TrimSpace(exprString)

	if IntegerPattern.MatchString(exprString) {
		if value, err := strconv.Atoi(exprString); err == nil {
			return NewVariable(value, IntegerType), nil
		}
	}

	if FloatPattern.MatchString(exprString) {
		if value, err := strconv.ParseFloat(exprString, 64); err == nil {
			return NewVariable(value, FloatType), nil
		}
	}

	if BooleanPattern.MatchString(exprString) {
		if exprString == TrueString {
			return NewVariable(true, BooleanType), nil
		}
		if exprString == FalseString {
			return NewVariable(false, BooleanType), nil
		}
	}

	if StringPattern.MatchString(exprString) {
		return NewVariable(exprString[1:len(exprString)-1], StringType), nil
	}

	return nil, fmt.Errorf("unable to parse literal: %s", exprString)
}

func parseExpression(exprString string) ([]string, error) {
    var elements []string
    var hasOperatorOrParen bool
    runes := []rune(exprString)
    n, i := len(runes), 0
    
    for i < n {
        switch {
        case unicode.IsDigit(runes[i]) || (runes[i] == SubtractionSymbol && isPotentialNegNumber(runes, i)):
            start := i
            i++
            for i < n && (unicode.IsDigit(runes[i]) || runes[i] == DecimalSymbol) {
                i++
            }
            elements = append(elements, string(runes[start:i]))
        case unicode.IsLetter(runes[i]):
            start := i
            i++
            for i < n && (unicode.IsLetter(runes[i]) || unicode.IsDigit(runes[i])) {
                i++
            }
            if i < n && runes[i] == ParenOpenSymbol {
                i = findClosingParen(runes, i) + 1
            }
            elements = append(elements, string(runes[start:i]))
        case runes[i] == ParenOpenSymbol || runes[i] == ParenCloseSymbol:
            hasOperatorOrParen = true
            elements = append(elements, string(runes[i]))
            i++
        case isOperator(runes[i]):
            hasOperatorOrParen = true
            elements = append(elements, string(runes[i]))
            i++
        case runes[i] == StringSymbol:
            start := i
            i++
            for i < n && (runes[i] != StringSymbol || (i > 0 && runes[i-1] == EscapeSymbol)) {
                i++
            }
            i++
            elements = append(elements, string(runes[start:i]))
        case unicode.IsSpace(runes[i]):
            i++
        default:
            return nil, fmt.Errorf("unknown character in expression: %c", runes[i])
        }
    }
    if !hasOperatorOrParen && len(elements) == 1 {
        return elements, nil
    }
    return elements, nil
}

func tokenizeExpression(exprs []string) ([]*Token, error) {
	var tokens []*Token
	for i, expr := range exprs {
		char := []rune(expr)[0]
		if isAction(expr) {
			tokens = append(tokens, NewToken(ActionToken, expr))
		} else if isLiteral(expr) {
			tokens = append(tokens, NewToken(LiteralToken, expr))
		} else if isVariable(expr) {
			tokens = append(tokens, NewToken(VariableToken, expr))
		} else if isOperator(char) {
			if char == SubtractionSymbol {
				if i == 0 || (i > 0 && (tokens[i-1].Type == ParenOpenToken || tokens[i-1].Type == OperatorAddToken || tokens[i-1].Type == OperatorSubtractToken || tokens[i-1].Type == OperatorMultiplyToken || tokens[i-1].Type == OperatorDivideToken || tokens[i-1].Type == OperatorModuloToken || tokens[i-1].Type == OperatorExponentToken)) {
					tokens = append(tokens, NewToken(OperatorUnaryMinusToken, expr))
				} else {
					tokens = append(tokens, NewToken(OperatorSubtractToken, expr))
				}
			} else {
				switch char {
				case AdditionSymbol:
					tokens = append(tokens, NewToken(OperatorAddToken, expr))
				case MultiplicationSymbol:
					tokens = append(tokens, NewToken(OperatorMultiplyToken, expr))
				case DivisionSymbol:
					tokens = append(tokens, NewToken(OperatorDivideToken, expr))
				case ModulusSymbol:
					tokens = append(tokens, NewToken(OperatorModuloToken, expr))
				case ExponentSymbol:
					tokens = append(tokens, NewToken(OperatorExponentToken, expr))
				}
			}
		} else if char == ParenOpenSymbol {
			tokens = append(tokens, NewToken(ParenOpenToken, expr))
		} else if char == ParenCloseSymbol {
			tokens = append(tokens, NewToken(ParenCloseToken, expr))
		} else {
			return nil, fmt.Errorf("unknown token found: %s", expr)
		}
	}
	return tokens, nil
}

func (s *Script) toRPN(tokens []*Token) []*Token {
	var output []*Token
	var operatorStack []*Token
	for _, token := range tokens {
		switch token.Type {
		case VariableToken, ActionToken, LiteralToken:
			output = append(output, token)
		case OperatorAddToken, OperatorSubtractToken, OperatorMultiplyToken, OperatorDivideToken, OperatorModuloToken, OperatorUnaryMinusToken, OperatorExponentToken:
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type == ParenOpenToken {
					break
				}
				if getPrecedence(top) >= getPrecedence(token) {
					output = append(output, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)
		case ParenOpenToken:
			operatorStack = append(operatorStack, token)
		case ParenCloseToken:
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				operatorStack = operatorStack[:len(operatorStack)-1]
				if top.Type == ParenOpenToken {
					break
				}
				output = append(output, top)
			}
		}
	}
	for len(operatorStack) > 0 {
		output = append(output, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}
	return output
}

func (s *Script) evaluateRPN(rpn []*Token) (*Variable, error) {
	var stack []float64
	for _, token := range rpn {
		switch token.Type {
		case ActionToken:
			action, err := s.parseActionToken(token)
			if err != nil {
				return nil, err
			}
			value, err := action.Execute(s)
			if err != nil {
				return nil, err
			}
			if values, ok := value.([]interface{}); ok {
				if len(values) > 1 {
					return nil, fmt.Errorf("action '%s' returned multiple values", token.Value)
				} else if len(values) == 1 {
					if floatValue, ok := values[0].(float64); ok {
						stack = append(stack, floatValue)
					} else {
						return nil, fmt.Errorf("action '%s' returned a non-float64 value", token.Value)
					}
				}
			} else {
				if floatValue, ok := value.(float64); ok {
					stack = append(stack, floatValue)
				} else {
					return nil, fmt.Errorf("action '%s' returned a non-float64 value", token.Value)
				}
			}
		case VariableToken:
			variable := s.Memory.GetVariable(token.Value)
			if variable == nil {
				return nil, fmt.Errorf("undefined variable: %s", token.Value)
			}
			if variable.Type != FloatType {
				fmt.Println(variable.Type, variable, variable.Value)
				if value, err := variable.toFloat(); err != nil {
					return nil, err
				} else {
					stack = append(stack, value)
				}
			} else {
				stack = append(stack, variable.Value.(float64))
			}
		case LiteralToken:
			value, err := strconv.ParseFloat(token.Value, 64)
			if err != nil {
				return nil, err
			}
			stack = append(stack, value)
		case OperatorUnaryMinusToken:
			if len(stack) < 1 {
				return nil, fmt.Errorf("insufficient values for unary operation")
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, -a)
		case OperatorAddToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a+b)
		case OperatorSubtractToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a-b)
		case OperatorMultiplyToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, a*b)
		case OperatorDivideToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			if b == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			stack = append(stack, a/b)
		case OperatorModuloToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, math.Mod(a, b))
		case OperatorExponentToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression")
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			stack = append(stack, math.Pow(a, b))
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression")
	}

	return NewVariable(stack[0], FloatType), nil
}

func (s *Script) parseExpression(exprString string) (*Variable, error) {
    exprs, err := parseExpression(exprString)
    if err != nil {
        return nil, err
    }

    if len(exprs) == 1 {
        expr := exprs[0]

        if isLiteral(expr) {
            literal, err := parseLiteral(expr)
            if err != nil {
                return nil, err
            }
            return literal, nil
        }

        if isVariable(expr) {
            variable := s.Memory.GetVariable(expr)
            if variable == nil {
                return nil, fmt.Errorf("undefined variable: %s", expr)
            }
            return variable, nil
        }

        return nil, fmt.Errorf("invalid or unrecognized single-token expression: %s", expr)
    }

    tokens, err := tokenizeExpression(exprs)
    if err != nil {
        return nil, err
    }
    rpn := s.toRPN(tokens)

    return s.evaluateRPN(rpn)
}
