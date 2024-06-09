package parser

import (
	"fmt"
	"strconv"

	"github.com/JunNishimura/go-lisp/ast"
	"github.com/JunNishimura/go-lisp/lexer"
	"github.com/JunNishimura/go-lisp/token"
)

// const (
// 	_ int = iota
// 	LOWEST
// 	SUM     // +
// 	PRODUCT // *
// )

// type (
// 	prefixParseFn func() ast.Expression
// )

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

	// 	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// 	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	// 	p.registerPrefix(token.PLUS, p.parsePrefixExpression)
	// 	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	// 	p.registerPrefix(token.ASTERISK, p.parsePrefixExpression)
	// 	p.registerPrefix(token.SLASH, p.parsePrefixExpression)

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

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) expectCur(t token.TokenType) bool {
	if p.curTokenIs(t) {
		p.nextToken()
		return true
	}
	p.curError(t)
	return false
}

func (p *Parser) expectOperator() (token.Token, bool) {
	if !p.isOperator() {
		p.curError(token.PLUS, token.MINUS, token.ASTERISK, token.SLASH, token.IDENT)
		return token.Token{}, false
	}
	operator := p.curToken
	p.nextToken()
	return operator, true
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		SExpressions: []ast.Cell{},
	}

	for p.curToken.Type != token.EOF {
		sexpression := p.parseSExpression()
		if sexpression != nil {
			program.SExpressions = append(program.SExpressions, sexpression)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseSExpression() ast.Cell {
	switch p.curToken.Type {
	case token.LPAREN:
		return p.parseConsCell()
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseAtom() ast.Cell {
	var cell ast.Cell

	switch p.curToken.Type {
	case token.INT:
		cell = p.parseIntegerAtom()
	default:
		return nil
	}
	p.nextToken()

	return cell
}

func (p *Parser) parseConsCell() *ast.ConsCell {
	if !p.expectCur(token.LPAREN) {
		return nil
	}

	// parse operator
	operator, ok := p.expectOperator()
	if !ok {
		return nil
	}

	// parse car
	car := p.parseSExpression()

	// parse cdr
	cdr := p.parseSExpression()

	if !p.expectCur(token.RPAREN) {
		return nil
	}

	return &ast.ConsCell{
		Operator: operator,
		Car:      car,
		Cdr:      cdr,
	}
}

func (p *Parser) parseIntegerAtom() *ast.Atom[int64] {
	intValue, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	intAtom := &ast.Atom[int64]{
		Token: p.curToken,
		Value: intValue,
	}

	return intAtom
}

func (p *Parser) isOperator() bool {
	return p.curTokenIs(token.PLUS) ||
		p.curTokenIs(token.MINUS) ||
		p.curTokenIs(token.ASTERISK) ||
		p.curTokenIs(token.SLASH) ||
		p.curTokenIs(token.IDENT)
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
