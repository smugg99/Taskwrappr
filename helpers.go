// helpers.go
package taskwrappr

import (
	"fmt"
	"path/filepath"
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
		if string(r) == op.Value[:1] {
			return true
		}
	}
	return false
}

func isOperator(operator string) bool {
	for _, knownOp := range Operators {
		if operator == knownOp.Value {
			return true
		}
	}
	return false
}

func categorizeOperator(operator string) (OperatorType, error) {
	for _, knownOp := range Operators {
		if operator == knownOp.Value {
			return knownOp.Type, nil
		}
	}

	return OperatorUndefined, fmt.Errorf("unknown operator: %s", operator)
}

func isReservedVariableName(name string) (bool, LiteralType) {
	for _, reservedVar := range ReservedVariables {
		if name == reservedVar.Name {
			return true, reservedVar.Type
		}
	}
	return false, LiteralUndefined
}

func sanitizeFilePath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty file path")
	}

	absPath, err := filepath.Abs(path)
    if err != nil {
        return "", fmt.Errorf("could not determine absolute path: %w", err)
    }

	if filepath.Ext(absPath) != ".tw" {
        return "", fmt.Errorf("invalid file type: %s", absPath)
    }

	return path, nil
}