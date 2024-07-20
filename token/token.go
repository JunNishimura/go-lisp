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

	// Special Form
	LAMBDA = "LAMBDA"
	QUOTE  = "'"

	PLUS  = "+"
	MINUS = "-"

	DOT       = "."
	BACKQUOTE = "`"
	COMMA     = ","

	// Delimiters
	LPAREN = "("
	RPAREN = ")"

	// Keywords
	NIL = "nil"
)

var keywords = map[string]TokenType{
	"nil":    NIL,
	"lambda": LAMBDA,
}

func LookupKeyword(symbol string) TokenType {
	if tok, ok := keywords[symbol]; ok {
		return tok
	}
	return SYMBOL
}
