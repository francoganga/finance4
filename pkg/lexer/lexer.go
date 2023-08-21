package lexer

import "errors"

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	INT    = "INT"
	SLASH  = "SLASH"
	DOLLAR = "DOLLAR"
	MINUS  = "MINUS"
	DESC   = "DESC"
	COMMA  = "COMMA"
	DOT    = "DOT"
	EOF    = "EOF"
	USD    = "USD"
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()

	return l
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	if l.ch == 'U' {
		err, lit := l.maybeReadUSD()

		if err == nil {
			tok.Type = USD
			tok.Literal = lit
		}
	}

	switch l.ch {
	case '-':
		tok = newToken(MINUS, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case '$':
		tok = newToken(DOLLAR, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '.':
		tok = newToken(DOT, l.ch)

	case 0:
		tok.Literal = ""
		tok.Type = EOF

	default:

		if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else if isLetter(l.ch) {
			tok.Type = DESC
			tok.Literal = l.readSentence()
		}
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

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Reads N chars
func (l *Lexer) readNChar(n int) {
	for i := 1; i <= n; i++ {
		l.readChar()
	}

}

func (l *Lexer) maybeReadUSD() (error, string) {
	if l.ch == 'U' && l.peekCharAt(1) == '$' && l.peekCharAt(2) == 'S' {
		l.readNChar(3)

		return nil, "U$S"
	}

	return errors.New("Could not read USD token"), ""
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Peeks char with offset from readPosition
func (l *Lexer) peekCharAt(offset int) byte {

	if l.position+offset >= len(l.input) {
		return 0
	}

	return l.input[l.position+offset]
}

func (l *Lexer) peekTill(till byte, fun func(ch byte) bool) bool {
	cv := l.position
	nv := l.position + 1

	for l.input[cv] != till {

		if fun(l.input[cv]) {
			return true
		}

		if nv >= len(l.input) {
			break
		}

		cv = nv
		nv += 1
	}

	return false
}

func (l *Lexer) readSentence() string {
	position := l.position

	for isLetter(l.ch) || l.ch == ' ' || isDigit(l.ch) {
		l.readChar()
		if l.ch == ' ' && l.peekChar() == ' ' {
			break
		}

		if !isLetter(l.peekChar()) && l.peekChar() != ' ' && l.peekChar() != 0 {
			l.readChar()
			if l.position < len(l.input) {
				l.readChar()
			}
		}
	}

	// fmt.Printf("start:=%d, end=%d\n", position, l.position)

	return l.input[position:l.position]
}
