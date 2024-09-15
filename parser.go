// parser.go
package taskwrappr

import "fmt"

type Parser struct {
	FilePath       string
	Tokens         []Token
	Index          uint
	IndexSinceLine uint
	Line           uint
	Token          Token
	Node           Node
}

func (p *Parser) String() string {
	return fmt.Sprintf("Parser{FilePath: %s, Index: %d, Token: %v, Node: %v}", p.FilePath, p.Index, p.Token, p.Node)
}

func NewParser(tokens []Token, filePath string) *Parser {
	parser := &Parser{
		FilePath: filePath,
		Tokens:   tokens,
	}

	parser.readToken()
	return parser
}

func (p *Parser) readToken() {
	if p.Index >= uint(len(p.Tokens)) {
		p.Token = nil
		return
	}

	p.Token = p.Tokens[p.Index]
	p.Index++
}

func (p *Parser) peekToken(x uint) Token {
	index := p.Index - 1 + x
	if index >= uint(len(p.Tokens)) {
		return EOFToken{p.Index, p.IndexSinceLine, p.Line}
	}

	return p.Tokens[index]
}

func (p *Parser) nextNode() (Node, error) {
	if p.Token == nil {
		return nil, nil
	}

	switch p.Token.Kind() {
		case TokenIdentifier:
			// Check if the next token is an identifier delimiter
			if p.peekToken(1).Kind() == TokenIdentifierDelimiter && p.peekToken(2).Kind() == TokenIdentifier {
				return p.parseBindingsNode()
			}
			return p.parseBindingNode()
		case TokenEOF:
			return nil, nil
	}

	// Unknown node
	return nil, fmt.Errorf("unexpected token: %s", p.Token)
}

func (p *Parser) parseBindingsNode() (Node, error) {
	node := BindingsNode{
		line:           p.Token.Line(),
		index:          p.Token.Index(),
		indexSinceLine: p.Token.IndexSinceLine(),
	}

	p.readToken()
	return node, nil
}

func (p *Parser) parseBindingNode() (Node, error) {
	if ok := p.Token.Kind() == TokenIdentifier; !ok {
		return nil, fmt.Errorf("unexpected token: %s", p.Token)
	}

	identifierNode := IdentifierNode{
		BaseName: p.Token.(IdentifierToken).Value,
		line:  p.Token.Line(),
		index: p.Token.Index(),
		indexSinceLine: p.Token.IndexSinceLine(),
	}

	node := BindingNode{
		Identifier:     identifierNode,
		line:           p.Token.Line(),
		index:          p.Token.Index(),
		indexSinceLine: p.Token.IndexSinceLine(),
	}

	p.readToken()
	return node, nil
}


func (p *Parser) Parse() ([]Node, error) {
	var ast []Node

	for {
		node, err := p.nextNode()
		if err != nil {
			return nil, fmt.Errorf("[%s:%d:%d] error parsing: %v", p.FilePath, p.Token.Line(), p.Token.IndexSinceLine(), err)
		}
		
		if node == nil {
			break
		}

		ast = append(ast, node)
	}

	return ast, nil
}
