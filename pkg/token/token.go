package token

// TokenType is basically a string
type TokenType string

// Token struct to represent a token
type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookUpIdent(ident string) TokenType {
	if tok , ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
