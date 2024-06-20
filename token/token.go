package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Symbols  + literals
	SYMBOL = "SYMBOL"
	INT    = "INT"

	// Operators
	PLUS  = "+"
	MINUS = "-"

	// Delimiters
	LPAREN = "("
	RPAREN = ")"
)
