// tokens.go
package taskwrappr

import "fmt"

// top-level building blocks of the language
type TokenKind int

const (
	TokenUndefined                     TokenKind = iota
	TokenEOF                           // end of file
	TokenIdentifier                    // e.g. foo, bar, someVar, someAction
	TokenOperation                     // e.g. +, :=, =, +=, &&, ||, !, ^^, ==, !=, <, <=, >, >=, ()
	TokenLiteral                       // e.g. nil, true, 1, 2.3, -3.14, 3, "foo", "bar", "baz"
	TokenExpression                    // e.g. 1 + 2, foo() + 2, bar(), someVar + 2, "foo" + "bar"
	TokenIdentifierDelimiter           // e.g. var1, var2, var3 := 1, 2, 3 (it's the comma)
	TokenBlockDelimiter                // e.g. {, }
	TokenExpressionDelimiter           // e.g. (, )
	TokenIndexingDelimiter             // e.g. [, ]
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
	case TokenExpression:
		return "expression"
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
	return fmt.Sprintf("[%d:%d] end of file", t.Line(), t.Index())
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
	return fmt.Sprintf("[%d:%d] identifier: %v", t.Line(), t.Index(), t.Value)
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
	return fmt.Sprintf("[%d:%d] operation: %v", t.Line(), t.Index(), t.Value)
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
	return fmt.Sprintf("[%d:%d] literal: %v, type: %s", t.Line(), t.Index(), t.Value, t.Type)
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

type ExpressionToken struct {
	Values []Token
	index  uint
	line   uint
}

func (t ExpressionToken) String() string {
	return fmt.Sprintf("[%d:%d] expression: %v", t.Line(), t.Index(), t.Values)
}

func (t ExpressionToken) Line() uint {
	return t.line
}

func (t ExpressionToken) Index() uint {
	return t.index
}

func (t ExpressionToken) Kind() TokenKind {
	return TokenExpression
}

type IdentifierDelimiterToken struct {
	index uint
	line  uint
}

func (t IdentifierDelimiterToken) String() string {
	return fmt.Sprintf("[%d:%d] identifier delimiter", t.Line(), t.Index())
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

	return fmt.Sprintf("[%d:%d] block delimiter: %c", t.Line(), t.Index(), _char)
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

	return fmt.Sprintf("[%d:%d] expression delimiter: %c", t.Line(), t.Index(), _char)
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

	return fmt.Sprintf("[%d:%d] indexing delimiter: %c", t.Line(), t.Index(), _char)
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