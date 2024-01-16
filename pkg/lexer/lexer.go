package lexer

import (
	"fmt"

	"github.com/hculpan/kabkey/pkg/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	lineNo       int
	linePosition int
	ch           byte
	errors       []string
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, lineNo: 1, linePosition: 0}
	l.errors = []string{}
	l.readChar()
	return l
}

func (l *Lexer) Errors() []string {
	return l.errors
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok = newToken(token.EQ, '=', l.lineNo, l.linePosition)
			tok.Literal = "=="
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.ch, l.lineNo, l.linePosition)
		}
	case '|':
		if l.peekChar() == '|' {
			tok = newToken(token.OR, '|', l.lineNo, l.linePosition)
			tok.Literal = "||"
			l.readChar()
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.lineNo, l.linePosition)
		}
	case '&':
		if l.peekChar() == '&' {
			tok = newToken(token.AND, '&', l.lineNo, l.linePosition)
			tok.Literal = "&&"
			l.readChar()
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.lineNo, l.linePosition)
		}
	case '"':
		val := l.readString()
		tok = newToken(token.STRING, ' ', l.lineNo, l.linePosition-len(val)-1)
		tok.Literal = val
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, l.lineNo, l.linePosition)
	case '(':
		tok = newToken(token.LPAREN, l.ch, l.lineNo, l.linePosition)
	case ')':
		tok = newToken(token.RPAREN, l.ch, l.lineNo, l.linePosition)
	case ',':
		tok = newToken(token.COMMA, l.ch, l.lineNo, l.linePosition)
	case '+':
		tok = newToken(token.PLUS, l.ch, l.lineNo, l.linePosition)
	case '{':
		tok = newToken(token.LBRACE, l.ch, l.lineNo, l.linePosition)
	case '}':
		tok = newToken(token.RBRACE, l.ch, l.lineNo, l.linePosition)
	case '!':
		if l.peekChar() == '=' {
			tok = newToken(token.NOT_EQ, '=', l.lineNo, l.linePosition)
			tok.Literal = "!="
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.ch, l.lineNo, l.linePosition)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch, l.lineNo, l.linePosition)
	case '/':
		tok = newToken(token.SLASH, l.ch, l.lineNo, l.linePosition)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, l.lineNo, l.linePosition)
	case '<':
		if l.peekChar() == '=' {
			tok = newToken(token.LTE, l.ch, l.lineNo, l.linePosition)
			tok.Literal = "<="
			l.readChar()
		} else {
			tok = newToken(token.LT, l.ch, l.lineNo, l.linePosition)
		}
	case '>':
		if l.peekChar() == '=' {
			tok = newToken(token.GTE, l.ch, l.lineNo, l.linePosition)
			tok.Literal = ">="
			l.readChar()
		} else {
			tok = newToken(token.GT, l.ch, l.lineNo, l.linePosition)
		}
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
		tok.LineNo = l.lineNo
		tok.Position = l.linePosition
	default:
		if isLetter(l.ch) {
			tok.LineNo = l.lineNo
			tok.Position = l.linePosition
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.LineNo = l.lineNo
			tok.Position = l.linePosition
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch, l.lineNo, l.linePosition)
		}
	}

	l.readChar()

	return tok
}

func (l *Lexer) readString() string {
	result := ""

	l.readChar() // skip over current quote

	for l.ch != '"' {
		if l.ch == '\n' || l.ch == '\r' || l.position >= len(l.input) {
			l.addError(l.lineNo, l.linePosition, "string not terminated with closing quote")
			break
		}
		result += string(l.ch)
		l.readChar()
	}

	return result
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.lineNo += 1
			l.linePosition = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || (ch == '_')
}

func newToken(tokenType token.TokenType, ch byte, lineNo, position int) token.Token {
	return token.Token{
		Type:     tokenType,
		Literal:  string(ch),
		LineNo:   lineNo,
		Position: position,
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
	l.linePosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) addError(lineNo, position int, format string, a ...interface{}) {
	msg := fmt.Sprintf("[%d:%d] %s", lineNo, position, fmt.Sprintf(format, a...))
	l.errors = append(l.errors, msg)
}
