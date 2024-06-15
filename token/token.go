package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	LPAREN = "("
	RPAREN = ")"

	// Keywords
	NIL  = "NIL"
	CONS = "CONS"
)

var keywords = map[string]TokenType{
	"nil":  NIL,
	"cons": CONS,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
