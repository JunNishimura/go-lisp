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

func (p *Parser) curError(t token.TokenType) {
	msg := fmt.Sprintf("expected current token to be %s, got %s instead", t, p.curToken.Type)
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
	// case token.LPAREN:
	// 	return p.parseConsCell()
	default:
		return p.parseAtom()
	}
}

func (p *Parser) parseAtom() ast.Cell {
	switch p.curToken.Type {
	case token.INT:
		return p.parseIntegerAtom()
	default:
		return nil
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

// func (p *Parser) parseStatement() ast.Statement {
// 	return p.parseExpressionStatement()
// }

// func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
// 	stmt := &ast.ExpressionStatement{Token: p.curToken}

// 	if !p.curTokenIs(token.LPAREN) {
// 		p.curError(token.LPAREN)
// 		return nil
// 	}
// 	p.nextToken()

// 	stmt.Expression = p.parseExpression(LOWEST)

// 	if !p.curTokenIs(token.RPAREN) {
// 		p.curError(token.RPAREN)
// 		return nil
// 	}

// 	p.nextToken()

// 	return stmt
// }

// func (p *Parser) curTokenIs(t token.TokenType) bool {
// 	return p.curToken.Type == t
// }

// func (p *Parser) peekTokenIs(t token.TokenType) bool {
// 	return p.peekToken.Type == t
// }

// func (p *Parser) parseExpression(precedence int) ast.Expression {
// 	prefix := p.prefixParseFns[p.curToken.Type]
// 	if prefix == nil {
// 		return nil
// 	}
// 	leftExp := prefix()

// 	return leftExp
// }

// func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
// 	p.prefixParseFns[tokenType] = fn
// }

// func (p *Parser) parsePrefixExpression() ast.Expression {
// 	expression := &ast.PrefixExpression{
// 		Token:    p.curToken,
// 		Operator: p.curToken.Literal,
// 	}

// 	p.nextToken()

// 	for !p.curTokenIs(token.RPAREN) {
// 		operand := p.parseExpression
// 	}
// }

// func (p *Parser) parseIntegerLiteral() ast.Expression {
// 	lit := &ast.IntegerLiteral{Token: p.curToken}

// 	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
// 	if err != nil {
// 		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
// 		p.errors = append(p.errors, msg)
// 		return nil
// 	}

// 	lit.Value = value

// 	return lit
// }
