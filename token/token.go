package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	EOF  = "EOF"
	INT  = "INT"
	CHAR = "CHAR"

	LBRACKET = "{"
	RBRACKET = "}"

	STAR = "*"
	OR   = "|"
)
