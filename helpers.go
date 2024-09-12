// helpers.go
package taskwrappr

import (
	"unicode"
)

func isSeparator(r rune) bool {
	for _, sep := range Separators {
		if r == sep {
			return true
		}
	}

	return false
}

func isIdentifierStart(r rune) bool {
	return unicode.IsLetter(r) || r == UnderscoreSymbol
}

func isStringStart(r rune) bool {
	return r == QuoteSymbol
}

func isNumberStart(r1, r2, r3 rune) bool {
	if r1 == 0 || r2 == 0 || r3 == 0 {
		return false
	}

	return unicode.IsDigit(r1) ||
		(r1 == SubtractionSymbol && (unicode.IsDigit(r2) ||
			(r2 == DecimalSymbol && unicode.IsDigit(r3)))) ||
		(r1 == DecimalSymbol && unicode.IsDigit(r2))
}

func isOperatorStart(r rune) bool {
	for _, op := range Operators {
		if string(r) == op[:1] {
			return true
		}
	}
	return false
}

func isOperator(operator string) bool {
	for _, knownOp := range Operators {
		if operator == knownOp {
			return true
		}
	}
	return false
}

func isReservedVariableName(name string) (bool, LiteralType) {
	for _, reserved := range ReservedVariableNames {
		if name == reserved {
			if varType, ok := ReservedVariablesTypes[name]; ok {
				return true, varType
			}
			return true, TypeUndefined
		}
	}
	return false, TypeUndefined
}
