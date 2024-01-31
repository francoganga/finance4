package parser

import (
	"fmt"
	"strconv"
)

type Parser struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Parser {
	l := &Parser{input: input}
	l.readChar()

	return l
}

func (p *Parser) readChar() {
	if p.readPosition >= len(p.input) {
		p.ch = 0
	} else {
		p.ch = p.input[p.readPosition]
	}

	p.position = p.readPosition
	p.readPosition += 1
}

func (p *Parser) skipWhitespace() {
	for p.ch == ' ' || p.ch == '\t' || p.ch == '\n' || p.ch == '\r' {
		p.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Reads N chars
func (p *Parser) readNChar(n int) {
	for i := 1; i <= n; i++ {
		p.readChar()
	}

}

func (p *Parser) readNumber() string {
	position := p.position
	for isDigit(p.ch) {
		p.readChar()
	}
	return p.input[position:p.position]
}

func (p *Parser) peekChar() byte {
	if p.readPosition >= len(p.input) {
		return 0
	} else {
		return p.input[p.readPosition]
	}
}

// Peeks char with offset from readPosition
func (p *Parser) peekCharAt(offset int) byte {

	if p.position+offset >= len(p.input) {
		return 0
	}

	return p.input[p.position+offset]
}

func (p *Parser) peekTill(till byte, fun func(ch byte) bool) bool {
	cv := p.position
	nv := p.position + 1

	for p.input[cv] != till {

		if fun(p.input[cv]) {
			return true
		}

		if nv >= len(p.input) {
			break
		}

		cv = nv
		nv += 1
	}

	return false
}

func (p *Parser) readDate() string {
	pos := p.position

	for isDigit(p.ch) || p.ch == '/' {
		p.readChar()
	}

	return p.input[pos:p.position]
}

func (p *Parser) ParseDate() string {
	p.skipWhitespace()

	if isDigit(p.ch) && p.peekCharAt(2) == '/' {
		return p.readDate()
	}

	return ""
}

func (p *Parser) ParseCode() string {
	p.skipWhitespace()

	if isDigit(p.ch) {
		pos := p.position

		for isDigit(p.ch) {
			p.readChar()
		}

		return p.input[pos:p.position]
	}

	return ""
}

func (p *Parser) ParseSentence() string {
	p.skipWhitespace()

	if isLetter(p.ch) {
		return p.readSentence()
	}

	return ""
}

func (p *Parser) readSentence() string {
	position := p.position

	for isLetter(p.ch) || p.ch == ' ' || isDigit(p.ch) {
		p.readChar()
		if p.ch == ' ' && p.peekChar() == ' ' {
			break
		}

		if !isLetter(p.peekChar()) && p.peekChar() != ' ' && p.peekChar() != 0 {
			p.readChar()
			if p.position < len(p.input) {
				p.readChar()
			}
		}
	}

	// fmt.Printf("start:=%d, end=%d\n", position, p.position)

	return p.input[position:p.position]

}

func (p *Parser) ParseAmount() (int, error) {
	p.skipWhitespace()

	if isDigit(p.ch) || p.ch == '-' || p.ch == '$' || p.ch == 'U' {
		return p.readAmount()
	}

	return 0, fmt.Errorf("Expected int but got %c", p.ch)

}

func (p *Parser) readAmount() (int, error) {
	val := ""

	for isDigit(p.ch) || p.ch == '-' || p.ch == '.' || p.ch == ',' || p.ch == '$' || p.ch == ' ' || p.ch == 'U' || p.ch == 'S' {

		if p.ch == ' ' && p.peekChar() == ' ' {
			break
		}

		if isDigit(p.ch) || p.ch == '-' {
			val += string(p.ch)
		}
		p.readChar()
	}

	return strconv.Atoi(val)
}

type Consumo struct {
	Date        string
	Code        string
	Description string
	Amount      int
	Balance     int
}

func (p *Parser) Parse() (Consumo, error) {
	consumo := Consumo{}

	consumo.Date = p.ParseDate()
	consumo.Code = p.ParseCode()
	consumo.Description = p.ParseSentence()

	amount, err := p.ParseAmount()

	if err != nil {
		return consumo, err
	}

	consumo.Amount = amount

	balance, err := p.ParseAmount()

	if err != nil {
		return consumo, err
	}

	consumo.Balance = balance

	consumo.Description += " " + p.ParseSentence()

	return consumo, nil
}

