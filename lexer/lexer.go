package lexer

import "github.com/JunNishimura/go-lisp/token"

type Lexer struct {
	input    string
	curPos   int
	nextPos  int
	prevChar byte
	curChar  byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	// move to the next character
	if l.nextPos >= len(l.input) {
		l.curChar = 0
	} else {
		l.curChar = l.input[l.nextPos]
	}
	l.curPos = l.nextPos
	l.nextPos++

	// store the previous character
	if l.curPos >= len(l.input) {
		l.prevChar = l.input[len(l.input)-1]
	} else if l.curPos > 0 {
		l.prevChar = l.input[l.curPos-1]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.curChar {
	case '(':
		tok = newToken(token.LPAREN, l.curChar)
	case ')':
		tok = newToken(token.RPAREN, l.curChar)
	case '+':
		if isSymbol(l.prevChar) {
			tok = newToken(token.SYMBOL, l.curChar)
		} else {
			tok = newToken(token.PLUS, l.curChar)
		}
	case '-':
		if isSymbol(l.prevChar) {
			tok = newToken(token.SYMBOL, l.curChar)
		} else {
			tok = newToken(token.MINUS, l.curChar)
		}
	case '.':
		tok = newToken(token.DOT, l.curChar)
	case '\'':
		tok = newToken(token.QUOTE, l.curChar)
	case '`':
		tok = newToken(token.BACKQUOTE, l.curChar)
	case ',':
		tok = newToken(token.COMMA, l.curChar)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.curChar) || isSpecialChar(l.curChar) {
			tok.Literal = l.readString()
			tok.Type = token.LookupKeyword(tok.Literal)
			return tok
		} else if isDigit(l.curChar) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.curChar)
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.curChar == ' ' || l.curChar == '\t' || l.curChar == '\n' || l.curChar == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSpecialChar(ch byte) bool {
	return ch == '*' ||
		ch == '/'
}

func isSymbol(ch byte) bool {
	return string(ch) == token.LPAREN
}

func (l *Lexer) readString() string {
	startPos := l.curPos
	for isLetter(l.curChar) || isSpecialChar(l.curChar) {
		l.readChar()
	}
	return l.input[startPos:l.curPos]
}

func (l *Lexer) readNumber() string {
	startPos := l.curPos
	for isDigit(l.curChar) {
		l.readChar()
	}
	return l.input[startPos:l.curPos]
}
