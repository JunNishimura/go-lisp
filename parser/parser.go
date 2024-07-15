package parser

import (
	"fmt"
	"strconv"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) curError(types ...token.TokenType) {
	expected := ""
	for i, t := range types {
		if i > 0 {
			expected += " or "
		}
		expected += string(t)
	}
	msg := fmt.Sprintf("expected token to be %s, got %s instead", expected, p.curToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectCur(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.curError(t)
	return false
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Expressions: []ast.SExpression{},
	}

	for p.curToken.Type != token.EOF {
		sexpression := p.parseSExpression()
		if sexpression != nil {
			program.Expressions = append(program.Expressions, sexpression)
		}
	}

	return program
}

func (p *Parser) parseSExpression() ast.SExpression {
	switch p.curToken.Type {
	case token.QUOTE, token.BACKQUOTE, token.COMMA:
		return p.parseDataMode()
	default:
		return p.parseCodeMode()
	}
}

func (p *Parser) parseList() ast.List {
	p.nextToken()

	// treat empty list as nil
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}}
	}

	// parse car
	car := p.parseSExpression()

	// if list is composed of only one element
	// treat it as a ConsCell with cdr being nil
	if p.curTokenIs(token.RPAREN) {
		p.nextToken()
		return &ast.ConsCell{
			CarField: car,
			CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
		}
	}

	var consCell *ast.ConsCell
	if p.curTokenIs(token.DOT) {
		// parse list defined below
		// "(" <s-expression> "." <s-expression> ")"
		p.nextToken()
		consCell = &ast.ConsCell{
			CarField: car,
			CdrField: p.parseSExpression(),
		}
	} else {
		// parse list defined below
		// "(" <s-expression> <s-expression> ... ")"
		consCell = &ast.ConsCell{
			CarField: car,
			CdrField: p.parseContinuousSExpression(),
		}
	}

	if !p.expectCur(token.RPAREN) {
		return nil
	}

	return consCell
}

func (p *Parser) parseCodeMode() ast.SExpression {
	switch p.curToken.Type {
	case token.LPAREN:
		return p.parseList()
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseDataMode() ast.SExpression {
	var car ast.SExpression
	switch p.curToken.Type {
	case token.QUOTE:
		car = &ast.Symbol{Token: p.curToken, Value: "quote"}
	case token.BACKQUOTE:
		car = &ast.Symbol{Token: p.curToken, Value: "backquote"}
	case token.COMMA:
		car = &ast.Symbol{Token: p.curToken, Value: "unquote"}
	}

	p.nextToken()

	sexpression := p.parseSExpression()

	return &ast.ConsCell{
		CarField: car,
		CdrField: &ast.ConsCell{
			CarField: sexpression,
			CdrField: &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}},
		},
	}
}

func (p *Parser) parseAtom() ast.Atom {
	atom := p.parseAtomByType()

	p.nextToken()

	return atom
}

func (p *Parser) parseAtomByType() ast.Atom {
	switch p.curToken.Type {
	case token.PLUS, token.MINUS:
		return p.parsePrefixAtom()
	case token.INT:
		return p.parseIntegerLiteral()
	case token.SYMBOL:
		return p.parseSymbol()
	case token.NIL:
		return &ast.Nil{Token: p.curToken}
	}

	msg := fmt.Sprintf("could not parse %q as atom", p.curToken.Literal)
	p.errors = append(p.errors, msg)
	return nil
}

func (p *Parser) parsePrefixAtom() *ast.PrefixAtom {
	prefixAtom := &ast.PrefixAtom{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	prefixAtom.Right = p.parseAtomByType()

	return prefixAtom
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	intValue, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.curToken,
		Value: intValue,
	}
}

func (p *Parser) parseSymbol() *ast.Symbol {
	return &ast.Symbol{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseContinuousSExpression() ast.SExpression {
	if p.curTokenIs(token.RPAREN) {
		return &ast.Nil{Token: token.Token{Type: token.NIL, Literal: "nil"}}
	}

	return &ast.ConsCell{
		CarField: p.parseSExpression(),
		CdrField: p.parseContinuousSExpression(),
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
