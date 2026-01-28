package lexer

import "pina771/regex-eng/token"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '*':
		tok.Literal = string(l.ch)
		tok.Type = token.STAR

	case '|':
		tok.Literal = string(l.ch)
		tok.Type = token.OR

	case '{':
		tok.Literal = string(l.ch)
		tok.Type = token.LBRACKET
	case '}':
		tok.Literal = string(l.ch)
		tok.Type = token.RBRACKET

	case 0:
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		tok.Literal = string(l.ch)
		tok.Type = token.CHAR
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
