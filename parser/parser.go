package parser

import (
	"pina771/regex-eng/lexer"
	"pina771/regex-eng/nfa"
	"pina771/regex-eng/token"
)

type (
	prefixFn func() Expression
	infixFn  func(Expression) Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixFns map[token.TokenType]prefixFn
	infixFns  map[token.TokenType]infixFn
}

var precedences = map[token.TokenType]int{
	"LOWEST":       0,
	token.OR:       1,
	token.CHAR:     3, // A character has higher precedence than '|' in regular expressions
	token.LBRACKET: 4,

	token.STAR: 10,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.prefixFns = map[token.TokenType]prefixFn{}
	p.prefixFns[token.CHAR] = p.parseChar

	p.infixFns = map[token.TokenType]infixFn{}
	p.infixFns[token.CHAR] = p.parseConcatenation
	p.infixFns[token.OR] = p.parseAlternation
	p.infixFns[token.STAR] = p.parseStarExpression

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseExpression(precedence int) Expression {
	var left Expression
	left = p.parsePrefix()

	for precedence < p.peekPrecedence() {
		p.nextToken()
		left = p.parseInfix(left)
	}

	return left
}

func (p *Parser) parsePrefix() Expression {
	curToken := p.curToken
	return p.prefixFns[curToken.Type]()
}

func (p *Parser) parseInfix(ex Expression) Expression {
	curToken := p.curToken
	return p.infixFns[curToken.Type](ex)
}

func (p *Parser) parseChar() Expression {
	charExp := &CharExpression{
		char: p.curToken.Literal,
	}
	return charExp
}

func (p *Parser) parseConcatenation(lhs Expression) Expression {
	rhs := p.parseExpression(precedences[token.CHAR])
	return &ConcatExpression{lhs, rhs}
}

func (p *Parser) parseAlternation(lhs Expression) Expression {
	alterToken := p.curToken
	p.nextToken()
	rhs := p.parseExpression(precedences[token.OR])
	return &AlternationExpression{lhs, rhs, alterToken}
}

func (p *Parser) parseStarExpression(lhs Expression) Expression {
	starToken := p.curToken
	return &StarExpression{lhs, starToken}
}

func (p *Parser) peekPrecedence() int {
	return precedences[p.peekToken.Type]
}

func (p *Parser) ToNFA() *nfa.Fragment {
	exp := p.parseExpression(precedences["LOWEST"])
	frag := exp.toNfa()
	return frag
}

// AST =================================
type Expression interface {
	toNfa() *nfa.Fragment
	TokenLiteral() string
}

type CharExpression struct {
	char string
}

func (ce *CharExpression) toNfa() *nfa.Fragment {
	return nfa.SingleChar(ce.char)
}
func (ce *CharExpression) TokenLiteral() string {
	return ce.char
}

type ConcatExpression struct {
	lhs Expression
	rhs Expression
}

func (ce *ConcatExpression) toNfa() *nfa.Fragment {
	left := ce.lhs.toNfa()
	right := ce.rhs.toNfa()
	return nfa.Concat(left, right)
}
func (ce *ConcatExpression) TokenLiteral() string {
	return "(" + ce.lhs.TokenLiteral() + ce.rhs.TokenLiteral() + ")"
}

type AlternationExpression struct {
	lhs   Expression
	rhs   Expression
	token token.Token
}

func (ae *AlternationExpression) toNfa() *nfa.Fragment {
	left := ae.lhs.toNfa()
	right := ae.rhs.toNfa()
	return nfa.Alternate(left, right)
}
func (ae *AlternationExpression) TokenLiteral() string {
	return "(" + ae.lhs.TokenLiteral() + ae.token.Literal + ae.rhs.TokenLiteral() + ")"
}

type StarExpression struct {
	lhs   Expression
	token token.Token
}

func (se *StarExpression) toNfa() *nfa.Fragment {
	left := se.lhs.toNfa()
	return nfa.Star(left)
}

func (se *StarExpression) TokenLiteral() string {
	return "(" + se.lhs.TokenLiteral() + se.token.Literal + ")"
}
