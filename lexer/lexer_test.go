package lexer

import (
	"pina771/regex-eng/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := "a*bc|d.{}"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.CHAR, "a"},
		{token.STAR, "*"},
		{token.CHAR, "b"},
		{token.CHAR, "c"},
		{token.OR, "|"},
		{token.CHAR, "d"},
		{token.CHAR, "."},
		{token.LBRACKET, "{"},
		{token.RBRACKET, "}"},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
