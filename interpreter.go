// interpreter.go
package taskwrappr

const (
	NewLineSymbol        = '\n'
	EscapeSymbol         = '\\'
	SemicolonSymbol      = ';'
	QuoteSymbol          = '"'
	CodeBlockOpenSymbol  = '{'
	CodeBlockCloseSymbol = '}'
	AssignmentSymbol     = '='
	UnderscoreSymbol     = '_'
	ParenOpenSymbol      = '('
	ParenCloseSymbol     = ')'
	BracketOpenSymbol    = '['
	BracketCloseSymbol   = ']'
	DelimiterSymbol      = ','
	DecimalSymbol        = '.'
	AdditionSymbol       = '+'
	SubtractionSymbol    = '-'
	MultiplicationSymbol = '*'
	DivisionSymbol       = '/'
	ModulusSymbol        = '%'
	ExponentSymbol       = '^'
	SelfReferenceSymbol  = '~'
	DeclarationSymbol    = ':'
	CommentSymbol        = '#'
)

const (
	TrueString                    = "true"
	FalseString                   = "false"
	NilString                     = "nil"
	LogicalAndString              = "&&"
	LogicalOrString               = "||"
	LogicalNotString              = "!"
	LogicalXorString              = "^^"
	EqualityString                = "=="
	InequalityString              = "!="
	LessThanString                = "<"
	LessThanOrEqualString         = "<="
	GreaterThanString             = ">"
	GreaterThanOrEqualString      = ">="
	DeclarationString             = string(DeclarationSymbol) + string(AssignmentSymbol)
	AugmentedAdditionString       = string(AdditionSymbol) + string(AssignmentSymbol)
	AugmentedSubtractionString    = string(SubtractionSymbol) + string(AssignmentSymbol)
	AugmentedMultiplicationString = string(MultiplicationSymbol) + string(AssignmentSymbol)
	AugmentedDivisionString       = string(DivisionSymbol) + string(AssignmentSymbol)
	AugmentedModulusString        = string(ModulusSymbol) + string(AssignmentSymbol)
	AugmentedExponentString       = string(ExponentSymbol) + string(AssignmentSymbol)
)

var Separators = string([]rune{
	SemicolonSymbol,
	NewLineSymbol,
})

var ArithmeticOperators = []string{
	string(AdditionSymbol),
	string(SubtractionSymbol),
	string(MultiplicationSymbol),
	string(DivisionSymbol),
	string(ModulusSymbol),
	string(ExponentSymbol),
}

var AugmentedOperators = []string{
	AugmentedAdditionString,
	AugmentedSubtractionString,
	AugmentedMultiplicationString,
	AugmentedDivisionString,
	AugmentedModulusString,
	AugmentedExponentString,
}

var ComparisonOperators = []string{
	EqualityString,
	InequalityString,
	LessThanString,
	LessThanOrEqualString,
	GreaterThanString,
	GreaterThanOrEqualString,
}

var LogicalOperators = []string{
	LogicalAndString,
	LogicalOrString,
	LogicalNotString,
	LogicalXorString,
}

var AssignmentOperators = []string{
	string(AssignmentSymbol),
	string(DecimalSymbol),
	DeclarationString,
}

var ReservedVariableNames = []string{
	TrueString,
	FalseString,
	NilString,
}

var ReservedVariablesTypes = map[string]LiteralType{
	TrueString:  TypeBool,
	FalseString: TypeBool,
	NilString:   TypeNil,
}

// Adjust this value if you add more operators
var MaxOperatorLength = 4

/*
Operators can only consist of symbols that are not letters
or digits so they can be easily distinguished from identifiers
*/
var Operators = append(
		append(
			append(
				append(
					ArithmeticOperators, AugmentedOperators...,
				), ComparisonOperators...,
			), LogicalOperators...,
		), AssignmentOperators...,
	)