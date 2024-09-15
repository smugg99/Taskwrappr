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

type OperatorTemplate struct {
	Type  OperatorType
	Value string
}

var ArithmeticOperators = []OperatorTemplate{
	{Type: OperatorAddition,       Value: string(AdditionSymbol)},
	{Type: OperatorSubtraction,    Value: string(SubtractionSymbol)},
	{Type: OperatorMultiplication, Value: string(MultiplicationSymbol)},
	{Type: OperatorDivision,       Value: string(DivisionSymbol)},
	{Type: OperatorModulus,        Value: string(ModulusSymbol)},
	{Type: OperatorExponentiation, Value: string(ExponentSymbol)},
}

var AugmentedOperators = []OperatorTemplate{
	{Type: OperatorAdditionAssignment,       Value: AugmentedAdditionString},
	{Type: OperatorSubtractionAssignment,    Value: AugmentedSubtractionString},
	{Type: OperatorMultiplicationAssignment, Value: AugmentedMultiplicationString},
	{Type: OperatorDivisionAssignment,       Value: AugmentedDivisionString},
	{Type: OperatorModulusAssignment,        Value: AugmentedModulusString},
	{Type: OperatorExponentiationAssignment, Value: AugmentedExponentString},
}

var ComparisonOperators = []OperatorTemplate{
	{Type: OperatorEqual,              Value: EqualityString},
	{Type: OperatorNotEqual,           Value: InequalityString},
	{Type: OperatorLessThan,           Value: LessThanString},
	{Type: OperatorLessThanOrEqual,    Value: LessThanOrEqualString},
	{Type: OperatorGreaterThan,        Value: GreaterThanString},
	{Type: OperatorGreaterThanOrEqual, Value: GreaterThanOrEqualString},
}

var LogicalOperators = []OperatorTemplate{
	{Type: OperatorAnd, Value: LogicalAndString},
	{Type: OperatorOr,  Value: LogicalOrString},
	{Type: OperatorNot, Value: LogicalNotString},
	{Type: OperatorXor, Value: LogicalXorString},
}

var AccessOperators = []OperatorTemplate{
	{Type: OperatorAssignment,  Value: string(AssignmentSymbol)},
	{Type: OperatorIndexing,    Value: string(DecimalSymbol)},
	{Type: OperatorDeclaration, Value: DeclarationString},
}

type ReservedVariableTemplate struct {
	Name string
	Type LiteralType
}

var ReservedVariables = []ReservedVariableTemplate{
	{Name: TrueString,  Type: LiteralBool},
	{Name: FalseString, Type: LiteralBool},
	{Name: NilString,   Type: LiteralNil},
}

// Adjust this value if you add more operators
var MaxOperatorLength = 4

/*
Operators ideally, should only consist of symbols that are not letters
or digits so they can be easily distinguished from identifiers
*/
var Operators = append(
		append(
			append(
				append(
					ArithmeticOperators, AugmentedOperators...,
				), ComparisonOperators...,
			), LogicalOperators...,
		), AccessOperators...,
	)