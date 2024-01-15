package lexer

import (
	"fmt"
	"testing"

	"github.com/hculpan/kabkey/pkg/token"
)

func TestString(t *testing.T) {
	input := `"hello world!"`
	l := NewLexer(input)
	tok := l.NextToken()
	if len(l.Errors()) != 0 {
		for _, e := range l.Errors() {
			fmt.Println(e)
		}
	}

	if tok.Type != token.STRING {
		t.Fatalf("expected type %q, got %q", token.STRING, tok.Type)
	}

	if tok.Literal != "hello world!" {
		t.Fatalf("expected literal %q, got %q", "hello world!", tok.Literal)
	}
}

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;
	
let add = fn(x, y) {
    x + y;
}
	
let result = add(five, ten);

!-/*5;

5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return    false;
}

10 == 10;

10 != 9;

let a = "This is a string!";

let i = 0;
while (i < 10) {
    let i = i + 1;
}
`

	tests := []struct {
		expectedType     token.TokenType
		expectedLiteral  string
		expectedLine     int
		expectedPosition int
	}{
		{token.LET, "let", 1, 1},
		{token.IDENT, "five", 1, 5},
		{token.ASSIGN, "=", 1, 10},
		{token.INT, "5", 1, 12},
		{token.SEMICOLON, ";", 1, 13},
		{token.LET, "let", 2, 1},
		{token.IDENT, "ten", 2, 5},
		{token.ASSIGN, "=", 2, 9},
		{token.INT, "10", 2, 11},
		{token.SEMICOLON, ";", 2, 13},
		{token.LET, "let", 4, 1},
		{token.IDENT, "add", 4, 5},
		{token.ASSIGN, "=", 4, 9},
		{token.FUNCTION, "fn", 4, 11},
		{token.LPAREN, "(", 4, 13},
		{token.IDENT, "x", 4, 14},
		{token.COMMA, ",", 4, 15},
		{token.IDENT, "y", 4, 17},
		{token.RPAREN, ")", 4, 18},
		{token.LBRACE, "{", 4, 20},
		{token.IDENT, "x", 5, 5},
		{token.PLUS, "+", 5, 7},
		{token.IDENT, "y", 5, 9},
		{token.SEMICOLON, ";", 5, 10},
		{token.RBRACE, "}", 6, 1},
		{token.LET, "let", 8, 1},
		{token.IDENT, "result", 8, 5},
		{token.ASSIGN, "=", 8, 12},
		{token.IDENT, "add", 8, 14},
		{token.LPAREN, "(", 8, 17},
		{token.IDENT, "five", 8, 18},
		{token.COMMA, ",", 8, 22},
		{token.IDENT, "ten", 8, 24},
		{token.RPAREN, ")", 8, 27},
		{token.SEMICOLON, ";", 8, 28},
		{token.BANG, "!", 10, 1},
		{token.MINUS, "-", 10, 2},
		{token.SLASH, "/", 10, 3},
		{token.ASTERISK, "*", 10, 4},
		{token.INT, "5", 10, 5},
		{token.SEMICOLON, ";", 10, 6},
		{token.INT, "5", 12, 1},
		{token.LT, "<", 12, 3},
		{token.INT, "10", 12, 5},
		{token.GT, ">", 12, 8},
		{token.INT, "5", 12, 10},
		{token.SEMICOLON, ";", 12, 11},
		{token.IF, "if", 14, 1},
		{token.LPAREN, "(", 14, 4},
		{token.INT, "5", 14, 5},
		{token.LT, "<", 14, 7},
		{token.INT, "10", 14, 9},
		{token.RPAREN, ")", 14, 11},
		{token.LBRACE, "{", 14, 13},
		{token.RETURN, "return", 15, 5},
		{token.TRUE, "true", 15, 12},
		{token.SEMICOLON, ";", 15, 16},
		{token.RBRACE, "}", 16, 1},
		{token.ELSE, "else", 16, 3},
		{token.LBRACE, "{", 16, 8},
		{token.RETURN, "return", 17, 5},
		{token.FALSE, "false", 17, 15},
		{token.SEMICOLON, ";", 17, 20},
		{token.RBRACE, "}", 18, 1},
		{token.INT, "10", 20, 1},
		{token.EQ, "==", 20, 4},
		{token.INT, "10", 20, 7},
		{token.SEMICOLON, ";", 20, 9},
		{token.INT, "10", 22, 1},
		{token.NOT_EQ, "!=", 22, 4},
		{token.INT, "9", 22, 7},
		{token.SEMICOLON, ";", 22, 8},
		{token.LET, "let", 24, 1},
		{token.IDENT, "a", 24, 5},
		{token.ASSIGN, "=", 24, 7},
		{token.STRING, "This is a string!", 24, 9},
		{token.SEMICOLON, ";", 24, 28},
		{token.LET, "let", 26, 1},
		{token.IDENT, "i", 26, 5},
		{token.ASSIGN, "=", 26, 7},
		{token.INT, "0", 26, 9},
		{token.SEMICOLON, ";", 26, 10},
		{token.WHILE, "while", 27, 1},
		{token.LPAREN, "(", 27, 7},
		{token.IDENT, "i", 27, 8},
		{token.LT, "<", 27, 10},
		{token.INT, "10", 27, 12},
		{token.RPAREN, ")", 27, 14},
		{token.LBRACE, "{", 27, 16},
		{token.LET, "let", 28, 5},
		{token.IDENT, "i", 28, 9},
		{token.ASSIGN, "=", 28, 11},
		{token.IDENT, "i", 28, 13},
		{token.PLUS, "+", 28, 15},
		{token.INT, "1", 28, 17},
		{token.SEMICOLON, ";", 28, 18},
		{token.RBRACE, "}", 29, 1},
		{token.EOF, "", 30, 1},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - token type wrong, expected %q, got %q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - token literal wrong, expected %q, got %q", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.LineNo != tt.expectedLine {
			t.Fatalf("tests[%d] - token line wrong, expected %d, got %d", i, tt.expectedLine, tok.LineNo)
		}
		if tok.Position != tt.expectedPosition {
			t.Fatalf("tests[%d] - token position wrong, expected %d, got %d", i, tt.expectedPosition, tok.Position)
		}
	}
}
