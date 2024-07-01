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

	PLUS  = "+"
	MINUS = "-"

	DOT = "."

	// Delimiters
	LPAREN = "("
	RPAREN = ")"

	// Keywords
	NIL    = "nil"
	LAMBDA = "lambda"
)

var keywords = map[string]TokenType{
	"nil":    NIL,
	"lambda": LAMBDA,
}

func LookupSymbol(symbol string) TokenType {
	if tok, ok := keywords[symbol]; ok {
		return tok
	}
	return SYMBOL
}
