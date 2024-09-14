// tokens.go
package taskwrappr

import "fmt"

// top-level building blocks of the language
type TokenKind int

const (
	TokenUndefined                     TokenKind = iota
	TokenEOF

	/*
	Just a name of a variable or action in this case
	e.g. foo, bar, someVar, someAction
	*/
	TokenIdentifier

	/*
	Any operation that can be performed on literals, identifiers, or expressions
	e.g. +, :=, =, +=, &&, ||, !, ^^, ==, !=, <, <=, >, >=, (), .
	*/
	TokenOperation

	/*
	Value that is fixed, used as is, can be of any type except action in this case
	e.g. nil, true, 1, 2.3, -3.14, 3, "foo", "bar", "baz"
	*/
	TokenLiteral

	/*
	Delimiter for separating identifiers in a binding or bindings
	e.g. var1, var2, var3 := 1, 2, 3 (it's the comma)
	*/
	TokenIdentifierDelimiter

	/*
	Delimiter for opening and closing a block
	e.g. {, }
	*/
	TokenBlockDelimiter

	/*
	Delimiter for dictating the order of operations in an expression
	e.g. (, )
	*/
	TokenExpressionDelimiter


	/*
	Delimiter for indexing an array or object
	e.g. [, ]
	*/
	TokenIndexingDelimiter
)

func (k TokenKind) String() string {
	switch k {
	case TokenUndefined:
		return "undefined"
	case TokenEOF:
		return "end of file"
	case TokenIdentifier:
		return "identifier"
	case TokenOperation:
		return "operation"
	case TokenLiteral:
		return "literal"
	case TokenIdentifierDelimiter:
		return "identifier delimiter"
	case TokenBlockDelimiter:
		return "block delimiter"
	case TokenExpressionDelimiter:
		return "expression delimiter"
	case TokenIndexingDelimiter:
		return "indexing delimiter"
	default:
		return "unknown"
	}
}

type LiteralType int

const (
	TypeUndefined LiteralType = iota
	TypeNil
	TypeBool
	TypeNumber
	TypeString
	TypeAction
	TypeArray
	TypeObject
)

func (t LiteralType) String() string {
	switch t {
	case TypeUndefined:
		return "undefined"
	case TypeNil:
		return "nil"
	case TypeBool:
		return "bool"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeAction:
		return "action"
	case TypeArray:
		return "array"
	case TypeObject:
		return "object"
	default:
		return "unknown"
	}
}

type Token interface {
	String() string
	Line() uint
	Index() uint
	Kind() TokenKind
}

type EOFToken struct {
	index uint
	line  uint
}

func (t EOFToken) String() string {
	return fmt.Sprintf("[%d:%d] %s", t.Line(), t.Index(), t.Kind())
}

func (t EOFToken) Line() uint {
	return t.line
}

func (t EOFToken) Index() uint {
	return t.index
}

func (t EOFToken) Kind() TokenKind {
	return TokenEOF
}

type IdentifierToken struct {
	Value string
	index uint
	line  uint
}

func (t IdentifierToken) String() string {
	return fmt.Sprintf("[%d:%d] %s -> value: %v", t.Line(), t.Index(), t.Kind(), t.Value)
}

func (t IdentifierToken) Line() uint {
	return t.line
}

func (t IdentifierToken) Index() uint {
	return t.index
}

func (t IdentifierToken) Kind() TokenKind {
	return TokenIdentifier
}

type OperationToken struct {
	Value string
	index uint
	line  uint
}

func (t OperationToken) String() string {
	return fmt.Sprintf("[%d:%d] %s -> value: %v", t.Line(), t.Index(), t.Kind(), t.Value)
}

func (t OperationToken) Line() uint {
	return t.line
}

func (t OperationToken) Index() uint {
	return t.index
}

func (t OperationToken) Kind() TokenKind {
	return TokenOperation
}

type LiteralToken struct {
	Value interface{}
	Type  LiteralType
	index uint
	line  uint
}

func (t LiteralToken) String() string {
	return fmt.Sprintf("[%d:%d] %s -> value: %v, type: %s", t.Line(), t.Index(), t.Kind(), t.Value, t.Type)
}

func (t LiteralToken) Line() uint {
	return t.line
}

func (t LiteralToken) Index() uint {
	return t.index
}

func (t LiteralToken) Kind() TokenKind {
	return TokenLiteral
}

type IdentifierDelimiterToken struct {
	index uint
	line  uint
}

func (t IdentifierDelimiterToken) String() string {
	return fmt.Sprintf("[%d:%d] %s", t.Line(), t.Index(), t.Kind())
}

func (t IdentifierDelimiterToken) Line() uint {
	return t.line
}

func (t IdentifierDelimiterToken) Index() uint {
	return t.index
}

func (t IdentifierDelimiterToken) Kind() TokenKind {
	return TokenIdentifierDelimiter
}

type BlockDelimiterToken struct {
	IsOpen bool
	index  uint
	line   uint
}

func (t BlockDelimiterToken) String() string {
	_char := CodeBlockCloseSymbol
	if t.IsOpen {
		_char = CodeBlockOpenSymbol
	}

	return fmt.Sprintf("[%d:%d] %s -> char: %c", t.Line(), t.Index(), t.Kind(), _char)
}

func (t BlockDelimiterToken) Line() uint {
	return t.line
}

func (t BlockDelimiterToken) Index() uint {
	return t.index
}

func (t BlockDelimiterToken) Kind() TokenKind {
	return TokenBlockDelimiter
}

type ExpressionDelimiterToken struct {
	IsOpen bool
	index  uint
	line   uint
}

func (t ExpressionDelimiterToken) String() string {
	_char := ParenCloseSymbol
	if t.IsOpen {
		_char = ParenOpenSymbol
	}

	return fmt.Sprintf("[%d:%d] %s -> char: %c", t.Line(), t.Index(), t.Kind(), _char)
}

func (t ExpressionDelimiterToken) Line() uint {
	return t.line
}

func (t ExpressionDelimiterToken) Index() uint {
	return t.index
}

func (t ExpressionDelimiterToken) Kind() TokenKind {
	return TokenExpressionDelimiter
}


type IndexingDelimiterToken struct {
	IsOpen bool
	index uint
	line  uint
}

func (t IndexingDelimiterToken) String() string {
	_char := ParenCloseSymbol
	if t.IsOpen {
		_char = ParenOpenSymbol
	}

	return fmt.Sprintf("[%d:%d] %s -> char: %c", t.Line(), t.Index(), t.Kind(), _char)
}

func (t IndexingDelimiterToken) Line() uint {
	return t.line
}

func (t IndexingDelimiterToken) Index() uint {
	return t.index
}

func (t IndexingDelimiterToken) Kind() TokenKind {
	return TokenIndexingDelimiter
}