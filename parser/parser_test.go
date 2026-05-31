package parser

import (
	"fmt"
	"pina771/regex-eng/lexer"
	"pina771/regex-eng/token"
	"testing"
)

func TestInitialSetup(t *testing.T) {
	input := "ab*"

	lexer := lexer.New(input)
	parser := New(lexer)

	expectedTokens := []token.Token{
		{Type: token.CHAR, Literal: "a"},
		{Type: token.CHAR, Literal: "b"},
		{Type: token.STAR, Literal: "*"},
	}

	for i, tt := range expectedTokens {
		if parser.curToken.Literal != tt.Literal {
			t.Fatalf("test[%d] failed. expected literal=[%q], got=%q",
				i,
				tt.Literal,
				parser.curToken.Literal,
			)
		}

		if parser.curToken.Type != tt.Type {
			t.Fatalf("test[%d] failed. expected Token.Type=[%q], got=%q",
				i,
				tt.Type,
				parser.curToken.Type,
			)
		}

		parser.nextToken()
	}

	// if parser.curToken.Literal != expectedTokens[0].Literal {
	// 	t.Fatalf("iniital token literal wrong."+
	// 		" Expected=%q, got=%q", parser.curToken.Literal,
	// 		expectedTokens[0].Literal)
	// }
}

func TestCharParsing(t *testing.T) {
	input := "a"

	lexer := lexer.New(input)
	parser := New(lexer)
	expectedExpression := &CharExpression{
		char: "a",
	}
	actualExpression := parser.parseExpression(0)

	if expectedExpression.TokenLiteral() != actualExpression.TokenLiteral() {
		t.Fatalf("test failed. Expression TokenLiteral wrong."+
			"expected=%q, got=%q", expectedExpression.TokenLiteral(),
			actualExpression.TokenLiteral())
	}
}

func TestConcatenation(t *testing.T) {
	inputs := []string{"ab", "abc"}
	expectedExpressions := []*ConcatExpression{
		{
			lhs: &CharExpression{char: "a"},
			rhs: &CharExpression{char: "b"},
		},
		{
			lhs: &CharExpression{char: "a"},
			rhs: &ConcatExpression{
				lhs: &CharExpression{char: "b"},
				rhs: &CharExpression{char: "c"},
			},
		},
	}

	for i, input := range inputs {
		parser := New(lexer.New(input))
		parsedExpression := parser.parseExpression(0)
		t.Logf("%s -> %s", inputs[i], parsedExpression.TokenLiteral())
		if parsedExpression.TokenLiteral() != expectedExpressions[i].TokenLiteral() {
			t.Fatalf(
				"test[%d] failed. Expression TokenLiteral wrong: expected=%s, got=%s",
				i,
				expectedExpressions[i].TokenLiteral(),
				parsedExpression.TokenLiteral(),
			)
		}
	}
}

func TestAlternation(t *testing.T) {
	input := "a|b"

	parser := New(lexer.New(input))
	parsedExpr := parser.parseExpression(0)
	fmt.Println(parsedExpr.TokenLiteral())
	if parsedExpr.TokenLiteral() != "(a|b)" {
		t.Fatalf("test failed. Expression TokenLiteral wrong. Expected: %s, got=%q",
			"a|b",
			parsedExpr.TokenLiteral(),
		)
	}
}

func TestMixedInfix(t *testing.T) {
	inputs := []string{
		"ab|c",
		"ab|cd",
		"a|bcd",
		"abc|d",
		"acd|b|cd",
	}
	expectedLiterals := []string{
		"((ab)|c)",
		"((ab)|(cd))",
		"(a|(b(cd)))",
		"((a(bc))|d)",
		"(((a(cd))|b)|(cd))",
	}
	for i, inp := range inputs {
		parser := New(lexer.New(inp))
		parsedExpr := parser.parseExpression(0)

		t.Logf("input=%s\texpectedOutput=%s\tactualOutput= %s",
			inp,
			expectedLiterals[i],
			parsedExpr.TokenLiteral(),
		)

		if parsedExpr.TokenLiteral() != expectedLiterals[i] {
			t.Fatalf("fail[%d].Input: %s\tExpected TokenLiteral: %s\tActualOutput: %s",
				i,
				inp,
				expectedLiterals[i],
				parsedExpr.TokenLiteral(),
			)
		}
	}
}

func TestNfaConstruction(t *testing.T) {
	input := "ab|c"
	parser := New(lexer.New(input))

	nfa := parser.ToNFA()
	t.Log(nfa)
}

func TestStar(t *testing.T) {
	inputs := []string{
		"ab*",
		"abc*",
		"a*PT",
	}
	expectedLiterals := []string{
		"(a(b*))",
		"(a(b(c*)))",
		"((a*)(PT))",
	}

	for idx := range inputs {
		parser := New(lexer.New(inputs[idx]))
		parsedExpr := parser.parseExpression(0)
		if parsedExpr.TokenLiteral() != expectedLiterals[idx] {
			t.Fatalf("Expression TokenLiteral wrong. Expected: %s, got=%s",
				expectedLiterals[idx],
				parsedExpr.TokenLiteral(),
			)
		}
	}
}
