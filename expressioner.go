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

func isLogicalOperator(s string) bool {
	return LogicalOperatorsPattern.MatchString(s)
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
	case LogicalOrToken, LogicalXorToken:
		return 0
	case LogicalAndToken:
		return 1
	case LogicalNotToken:
		return 5
	case EqualityToken, InequalityToken, LessThanToken, LessThanOrEqualToken, GreaterThanToken, GreaterThanOrEqualToken:
		return 1
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

func castToFloat(v *Variable) (*Variable, error) {
	value, err := v.toFloat()
	if err != nil {
		return nil, fmt.Errorf("could not cast %v to float: %v", v.Value, err)
	}
	return NewVariable(value, FloatType), nil
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
			if i+1 < n && runes[i] == ExponentSymbol && runes[i+1] == ExponentSymbol {
				elements = append(elements, string(runes[i:i+2]))
				i += 2
			} else {
				elements = append(elements, string(runes[i]))
				i++
			}
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
			start := i
			for i < n && !unicode.IsDigit(runes[i]) && !unicode.IsLetter(runes[i]) && !unicode.IsSpace(runes[i]) && runes[i] != ParenOpenSymbol && runes[i] != ParenCloseSymbol {
				i++
			}
			opStr := string(runes[start:i])
			if isLogicalOperator(opStr) {
				hasOperatorOrParen = true
				elements = append(elements, opStr)
			} else {
				return nil, fmt.Errorf("unknown character in expression: %s", string(runes[start:i]))
			}
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
		switch expr {
		case LogicalAndString:
			tokens = append(tokens, NewToken(LogicalAndToken, expr))
		case LogicalOrString:
			tokens = append(tokens, NewToken(LogicalOrToken, expr))
		case LogicalNotString:
			tokens = append(tokens, NewToken(LogicalNotToken, expr))
		case LogicalXorString:
			tokens = append(tokens, NewToken(LogicalXorToken, expr))
		case EqualityString:
			tokens = append(tokens, NewToken(EqualityToken, expr))
		case InequalityString:
			tokens = append(tokens, NewToken(InequalityToken, expr))
		case LessThanString:
			tokens = append(tokens, NewToken(LessThanToken, expr))
		case LessThanOrEqualString:
			tokens = append(tokens, NewToken(LessThanOrEqualToken, expr))
		case GreaterThanString:
			tokens = append(tokens, NewToken(GreaterThanToken, expr))
		case GreaterThanOrEqualString:
			tokens = append(tokens, NewToken(GreaterThanOrEqualToken, expr))
		default:
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
					case SubtractionSymbol:
						tokens = append(tokens, NewToken(OperatorSubtractToken, expr))
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
	}
	return tokens, nil
}

func ensureCompatibleOperands(a, b *Variable) (*Variable, *Variable, error) {
	var err error
	if a.Type == NilType || b.Type == NilType {
		return nil, nil, fmt.Errorf("operands cannot be nil")
	}

	if a.Type == InvalidType || b.Type == InvalidType {
		return nil, nil, fmt.Errorf("operands cannot be invalid")
	}

	if a.Type == StringType || b.Type == StringType {
		a.Value, b.Value = fmt.Sprintf("%v", a.Value), fmt.Sprintf("%v", b.Value)
		a.Type, b.Type = StringType, StringType
	} else {
		if a.Type != FloatType {
			if a, err = castToFloat(a); err != nil {
				return nil, nil, fmt.Errorf("cannot cast '%v' of type %s to float: %v", a.Value, a.Type.String(), err)
			}
		}
		if b.Type != FloatType {
			if b, err = castToFloat(b); err != nil {
				return nil, nil, fmt.Errorf("cannot cast '%v' of type %s to float: %v", b.Value, b.Type.String(), err)
			}
		}
	}
	return a, b, nil
}

func (s *Script) toRPN(tokens []*Token) []*Token {
	var output []*Token
	var operatorStack []*Token
	for _, token := range tokens {
		switch token.Type {
		case VariableToken, ActionToken, LiteralToken:
			output = append(output, token)
		case OperatorAddToken, OperatorSubtractToken, OperatorMultiplyToken, OperatorDivideToken, OperatorModuloToken, OperatorUnaryMinusToken, OperatorExponentToken,
			LogicalAndToken, LogicalOrToken, LogicalXorToken, LogicalNotToken, EqualityToken, InequalityToken, LessThanToken,
			LessThanOrEqualToken, GreaterThanToken, GreaterThanOrEqualToken:
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
	var stack []*Variable
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
			if len(value) != 1 {
				return nil, fmt.Errorf("action returned multiple values: %v", value)
			}
			stack = append(stack, value[0])
		case VariableToken:
			variable := s.CurrentBlock.Memory.GetVariable(token.Value)
			if variable == nil {
				return nil, fmt.Errorf("undefined variable: %s", token.Value)
			}
			stack = append(stack, variable)
		case LiteralToken:
			variable, err := parseLiteral(token.Value)
			if err != nil {
				return nil, err
			}
			stack = append(stack, variable)
		case OperatorUnaryMinusToken:
			if len(stack) < 1 {
				return nil, fmt.Errorf("insufficient values for unary operation")
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if a.Type != FloatType && a.Type != IntegerType {
				return nil, fmt.Errorf("unary minus operand is not a number, got %s", a.Type.String())
			}
			castA, _ := a.toFloat()
			stack = append(stack, NewVariable(-castA, FloatType))
		case OperatorAddToken, OperatorSubtractToken, OperatorMultiplyToken, OperatorDivideToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression for %s", token.Type.String())
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			var err error
			a, b, err = ensureCompatibleOperands(a, b)
			if err != nil {
				return nil, fmt.Errorf("ensure compatible operands: %v", err)
			}
			if a.Type == StringType && b.Type == StringType {
				if token.Type == OperatorAddToken {
					result := a.Value.(string) + b.Value.(string)
					stack = append(stack, NewVariable(result, StringType))
				} else {
					return nil, fmt.Errorf("only addition '+' operator is supported for strings, got %s", token.Type.String())
				}
			} else if a.Type == FloatType && b.Type == FloatType {
				castA, _ := a.toFloat()
				castB, _ := b.toFloat()
				var result float64
				switch token.Type {
				case OperatorAddToken:
					result = castA + castB
				case OperatorSubtractToken:
					result = castA - castB
				case OperatorMultiplyToken:
					result = castA * castB
				case OperatorDivideToken:
					if castB == 0 {
						return nil, fmt.Errorf("division by zero")
					}
					result = castA / castB
				}
				stack = append(stack, NewVariable(result, FloatType))
			} else {
				return nil, fmt.Errorf("type mismatch between %v (type: %s) and %v (type: %s)", a.Value, a.Type.String(), b.Value, b.Type.String())
			}
		case OperatorModuloToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("%s", token.Type.String())
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			a, b, err := ensureCompatibleOperands(a, b)
			if err != nil {
				return nil, fmt.Errorf("ensure compatible operands: %v", err)
			}
			castA, _ := a.toFloat()
			castB, _ := b.toFloat()
			stack = append(stack, NewVariable(math.Mod(castA, castB), FloatType))
		case OperatorExponentToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values in expression for %s", token.Type.String())
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			a, b, err := ensureCompatibleOperands(a, b)
			if err != nil {
				return nil, fmt.Errorf("ensure compatible operands: %v", err)
			}
			castA, _ := a.toFloat()
			castB, _ := b.toFloat()
			stack = append(stack, NewVariable(math.Pow(castA, castB), FloatType))
		case LogicalAndToken, LogicalOrToken, LogicalXorToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values for %s", token.Type.String())
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			castA, _ := a.toBool()
			castB, _ := b.toBool()
			var result bool
			switch token.Type {
			case LogicalAndToken:
				result = castA && castB
			case LogicalOrToken:
				result = castA || castB
			case LogicalXorToken:
				result = (castA || castB) && !(castA && castB)
			}
			stack = append(stack, NewVariable(result, BooleanType))
		case LogicalNotToken:
			if len(stack) < 1 {
				return nil, fmt.Errorf("insufficient values for %s", token.Type.String())
			}
			a := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			castA, _ := a.toBool()
			stack = append(stack, NewVariable(!castA, BooleanType))
		case EqualityToken, InequalityToken, LessThanToken, LessThanOrEqualToken, GreaterThanToken, GreaterThanOrEqualToken:
			if len(stack) < 2 {
				return nil, fmt.Errorf("insufficient values for %s", token.Type.String())
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			a, b, err := ensureCompatibleOperands(a, b)
			if err != nil {
				return nil, err
			}
			var result bool
			if a.Type == StringType && b.Type == StringType {
				valA := a.Value.(string)
				valB := b.Value.(string)
				switch token.Type {
				case EqualityToken:
					result = (valA == valB)
				case InequalityToken:
					result = (valA != valB)
				}
			} else if a.Type == FloatType && b.Type == FloatType {
				valA, _ := a.toFloat()
				valB, _ := b.toFloat()
				switch token.Type {
				case EqualityToken:
					result = (valA == valB)
				case InequalityToken:
					result = (valA != valB)
				case LessThanToken:
					result = (valA < valB)
				case LessThanOrEqualToken:
					result = (valA <= valB)
				case GreaterThanToken:
					result = (valA > valB)
				case GreaterThanOrEqualToken:
					result = (valA >= valB)
				}
			} else {
				return nil, fmt.Errorf("type mismatch between %v and %v", a.Type.String(), b.Type.String())
			}
			stack = append(stack, NewVariable(result, BooleanType))
		default:
			return nil, fmt.Errorf("unknown token type: %v", token.Type)
		}
	}
	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression, final stack: %v", stack)
	}
	return stack[0], nil
}

func (s *Script) parseExpression(exprString string) (*Action, error) {
	exprs, err := parseExpression(exprString)
	if err != nil {
		return nil, err
	}

	expressionAction := func(s *Script, args ...*Variable) ([]*Variable, error) {
		tokens, err := tokenizeExpression(exprs)
		if err != nil {
			return nil, err
		}

		if len(tokens) == 1 {
			switch tokens[0].Type {
			case ActionToken:
				action, err := s.parseActionToken(tokens[0])
				if err != nil {
					return nil, err
				}

				value, err := action.Execute(s)
				if err != nil {
					return nil, err
				}

				return value, nil
			case VariableToken:
				variable := s.CurrentBlock.Memory.GetVariable(tokens[0].Value)
				if variable == nil {
					return nil, fmt.Errorf("undefined variable: %s", tokens[0].Value)
				}
				return []*Variable{variable}, nil
			}
		}

		rpn := s.toRPN(tokens)

		result, err := s.evaluateRPN(rpn)
		if err != nil {
			return nil, err
		}

		return []*Variable{result}, nil
	}

	return NewAction(expressionAction, nil), nil
}
