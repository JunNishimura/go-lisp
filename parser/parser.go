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
		p.nextToken()
	}

	return program
}

func (p *Parser) parseSExpression() ast.SExpression {
	switch p.curToken.Type {
	case token.LPAREN:
		return p.parseList()
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseAtom() ast.Atom {
	var cell ast.Atom

	switch p.curToken.Type {
	case token.PLUS, token.MINUS:
		cell = p.parsePrefixAtom()
	case token.INT:
		cell = p.parseIntegerLiteral()
	case token.SYMBOL:
		cell = p.parseSymbol()
		// case token.IDENT:
		// 	cell = p.parseIdentifier()
		// case token.NIL:
		// 	cell = p.parseNilLiteral()
	}
	p.nextToken()

	return cell
}

func (p *Parser) parsePrefixAtom() ast.Atom {
	prefixAtom := &ast.PrefixAtom{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	prefixAtom.Right = p.parseAtom()

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

func (p *Parser) parseList() *ast.List {
	if !p.expectCur(token.LPAREN) {
		return nil
	}

	car := p.parseSExpression()

	var cdr []ast.SExpression
	for !p.curTokenIs(token.RPAREN) {
		cdr = append(cdr, p.parseSExpression())
	}

	if !p.expectCur(token.RPAREN) {
		return nil
	}

	return &ast.List{
		Car: car,
		Cdr: cdr,
	}
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
