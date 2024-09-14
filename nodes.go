// nodes.go
package taskwrappr

import "fmt"

type NodeKind int

const (
	NodeUndefined NodeKind = iota

	/*
	Binding of an identifier to an expression.
	Left side has to have exactly one identifier,
	right side has to have exactly one expression
	e.g. var1, var2, var3 := 1, 2, 3
	*/
	NodeBinding

	/*
	Binding of multiple identifiers to multiple expressions.
	Left side has to have the same number (but more than one)
	of identifiers as the right side has expressions, separated by delimiters
	e.g. var1, var2, var3 := 1, 2, 3 (aka tuple unpacking)
	*/
	NodeBindings

	/*
	Sequence of literals, identifiers, and operations that evaluate to a value
	e.g. 1 + 2, foo() + 2, bar(), someVar + 2, "foo" + "bar"
	*/
	NodeExpression

	/*
	Sequence of nodes that are executed in order
	e.g. { var1 := 1; someFunc(54 + 4); var3 = 3 + 34; }
	*/
	NodeBlock

	/*
	Name or a chained index of a variable or action
	e.g. foo, bar["key"], someVar[1], someAction, someObject.key
	*/
	NodeIdentifier

	/*
	Value that is fixed, used as is, can be of any type, even action
	e.g. nil, true, 1, 2.3, -3.14, 3, "foo", "bar", "baz", { return(5 + 2); }
	*/
	NodeLiteral
)

func (k NodeKind) String() string {
	switch k {
	case NodeUndefined:
		return "undefined"
	case NodeBinding:
		return "binding"
	case NodeBindings:
		return "bindings"
	case NodeExpression:
		return "expression"
	case NodeBlock:
		return "block"
	case NodeIdentifier:
		return "identifier"
	case NodeLiteral:
		return "literal"
	default:
		return "unknown"
	}
}

type Node interface {
	String() string
	Line() uint
	Index() uint
	Kind() NodeKind
}

type BindingNode struct {
	//Identifier []IdentifierNode
	//Expression []ExpressionNode
	line uint
	index uint
}

func (n BindingNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n BindingNode) Line() uint {
	return n.line
}

func (n BindingNode) Index() uint {
	return n.index
}

func (n BindingNode) Kind() NodeKind {
	return NodeBinding
}


type BindingsNode struct {
	//Identifiers []IdentifierNode
	//Expressions []ExpressionNode
	line uint
	index uint
}

func (n BindingsNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n BindingsNode) Line() uint {
	return n.line
}

func (n BindingsNode) Index() uint {
	return n.index
}

func (n BindingsNode) Kind() NodeKind {
	return NodeBinding
}


type ExpressionNode struct {
	//Nodes []Node
	line uint
	index uint
}

func (n ExpressionNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n ExpressionNode) Line() uint {
	return n.line
}

func (n ExpressionNode) Index() uint {
	return n.index
}

func (n ExpressionNode) Kind() NodeKind {
	return NodeExpression
}


type BlockNode struct {
	//Nodes []Node
	line uint
	index uint
}

func (n BlockNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n BlockNode) Line() uint {
	return n.line
}

func (n BlockNode) Index() uint {
	return n.index
}

func (n BlockNode) Kind() NodeKind {
	return NodeBlock
}


type IdentifierNode struct {
	//Value string
	line uint
	index uint
}

func (n IdentifierNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n IdentifierNode) Line() uint {
	return n.line
}

func (n IdentifierNode) Index() uint {
	return n.index
}

func (n IdentifierNode) Kind() NodeKind {
	return NodeIdentifier
}


type LiteralNode struct {
	//Value interface{}
	line uint
	index uint
}

func (n LiteralNode) String() string {
	return fmt.Sprintf("[%d:%d] %s", n.Line(), n.Index(), n.Kind())
}

func (n LiteralNode) Line() uint {
	return n.line
}

func (n LiteralNode) Index() uint {
	return n.index
}

func (n LiteralNode) Kind() NodeKind {
	return NodeLiteral
}
