package lexer

import (
	"testing"

	"github.com/skatsuta/monkey-compiler/token"
)

func TestNextToken(t *testing.T) {
	input := `
	let five = 5;
	let ten = 10;
	let add = fn(x, y) {
		x + y;
	};
	let result = add(five, ten);
	!-/*0;
	2 < 10 > 7;

	if (5 < 10) {
		return true;
	} else {
		return false;
	}

	10 <= 11;
	10 >= 9;
	10 == 10;
	10 != 9;

	true && false;
	true || false;

	"foobar";
	"foo bar";

	[1, 2];

	{"foo": "bar"};

	# comment
	let a = 1; # inline comment

	let b = 123.45;
	let c = 0.678;
	let d = 9.0;

	a = 2;
	b = nil;
	c = 1;
	c += 2;
	c -= 3;
	c *= 4;
	c /= 5;

	macro(x, y) { x + y; };
	`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.Let, "let"},
		{token.Ident, "five"},
		{token.Assign, "="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "ten"},
		{token.Assign, "="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "add"},
		{token.Assign, "="},
		{token.Function, "fn"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "result"},
		{token.Assign, "="},
		{token.Ident, "add"},
		{token.LParen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.RParen, ")"},
		{token.Semicolon, ";"},
		{token.Bang, "!"},
		{token.Minus, "-"},
		{token.Slash, "/"},
		{token.Astarisk, "*"},
		{token.Int, "0"},
		{token.Semicolon, ";"},
		{token.Int, "2"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.GT, ">"},
		{token.Int, "7"},
		{token.Semicolon, ";"},
		{token.If, "if"},
		{token.LParen, "("},
		{token.Int, "5"},
		{token.LT, "<"},
		{token.Int, "10"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.True, "true"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Else, "else"},
		{token.LBrace, "{"},
		{token.Return, "return"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Int, "10"},
		{token.LE, "<="},
		{token.Int, "11"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.GE, ">="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.Eq, "=="},
		{token.Int, "10"},
		{token.Semicolon, ";"},
		{token.Int, "10"},
		{token.NEq, "!="},
		{token.Int, "9"},
		{token.Semicolon, ";"},
		{token.True, "true"},
		{token.And, "&&"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.True, "true"},
		{token.Or, "||"},
		{token.False, "false"},
		{token.Semicolon, ";"},
		{token.String, "foobar"},
		{token.Semicolon, ";"},
		{token.String, "foo bar"},
		{token.Semicolon, ";"},
		{token.LBracket, "["},
		{token.Int, "1"},
		{token.Comma, ","},
		{token.Int, "2"},
		{token.RBracket, "]"},
		{token.Semicolon, ";"},
		{token.LBrace, "{"},
		{token.String, "foo"},
		{token.Colon, ":"},
		{token.String, "bar"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "a"},
		{token.Assign, "="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "b"},
		{token.Assign, "="},
		{token.Float, "123.45"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "c"},
		{token.Assign, "="},
		{token.Float, "0.678"},
		{token.Semicolon, ";"},
		{token.Let, "let"},
		{token.Ident, "d"},
		{token.Assign, "="},
		{token.Float, "9.0"},
		{token.Semicolon, ";"},
		{token.Ident, "a"},
		{token.Assign, "="},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.Ident, "b"},
		{token.Assign, "="},
		{token.Nil, "nil"},
		{token.Semicolon, ";"},
		{token.Ident, "c"},
		{token.Assign, "="},
		{token.Int, "1"},
		{token.Semicolon, ";"},
		{token.Ident, "c"},
		{token.AddAssign, "+="},
		{token.Int, "2"},
		{token.Semicolon, ";"},
		{token.Ident, "c"},
		{token.SubAssign, "-="},
		{token.Int, "3"},
		{token.Semicolon, ";"},
		{token.Ident, "c"},
		{token.MulAssign, "*="},
		{token.Int, "4"},
		{token.Semicolon, ";"},
		{token.Ident, "c"},
		{token.DivAssign, "/="},
		{token.Int, "5"},
		{token.Semicolon, ";"},
		{token.Macro, "macro"},
		{token.LParen, "("},
		{token.Ident, "x"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.RParen, ")"},
		{token.LBrace, "{"},
		{token.Ident, "x"},
		{token.Plus, "+"},
		{token.Ident, "y"},
		{token.Semicolon, ";"},
		{token.RBrace, "}"},
		{token.Semicolon, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Logf("tests[%d] - tok: %#v", i, tok)
			t.Fatalf("tests[%d] - token type wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Logf("tests[%d] - tok: %#v", i, tok)
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
