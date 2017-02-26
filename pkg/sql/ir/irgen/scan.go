package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"os"
)

type Scanner struct {
	s         scanner.Scanner
	convPos   func(pos token.Pos) string
	results   defList
	lastPos   string
	lastError string
}

// Error implements the irgenLexer interface.
func (s *Scanner) Error(str string) {
	s.lastError = str
}

// Lex implements the irgenLexer interface.
func (s *Scanner) Lex(lval *irgenSymType) int {
	var tok token.Token
	var pos token.Pos

	for {
		pos, tok, lval.str = s.s.Scan()
		s.lastPos = s.convPos(pos)
		if tok == token.SEMICOLON && lval.str == "\n" {
			// We don't want to see implicit semicolons. Just
			// ignore them.
			continue
		}
		break
	}

	switch tok {
	case token.STRING:
		return STR
	case token.SEMICOLON:
		return ';'
	case token.OR:
		return '|'
	case token.LBRACE:
		return '{'
	case token.RBRACE:
		return '}'
	case token.ASSIGN:
		return '='
	case token.MUL:
		return '*'
	case token.IDENT:
		switch lval.str {
		case "sum":
			return SUM
		case "def":
			return DEF
		case "sql":
			return SQL
		case "enum":
			return ENUM
		default:
			return IDENT
		}
	case token.EOF:
		return 0
	default:
		if t, ok := tokenmap[tok]; ok {
			lval.str = t
			return IDENT
		}
		fmt.Fprintf(os.Stderr, "WOO %q // %q\n", tok, lval.str)
		return ERROR
	}
}

var tokenmap = map[token.Token]string{
	token.BREAK:       "break",
	token.CASE:        "case",
	token.CHAN:        "chan",
	token.CONST:       "const",
	token.CONTINUE:    "continue",
	token.DEFAULT:     "default",
	token.DEFER:       "defer",
	token.ELSE:        "else",
	token.FALLTHROUGH: "fallthrough",
	token.FOR:         "for",
	token.FUNC:        "func",
	token.GO:          "go",
	token.GOTO:        "goto",
	token.IF:          "if",
	token.IMPORT:      "import",
	token.INTERFACE:   "interface",
	token.MAP:         "map",
	token.PACKAGE:     "package",
	token.RANGE:       "range",
	token.RETURN:      "return",
	token.SELECT:      "select",
	token.STRUCT:      "struct",
	token.SWITCH:      "switch",
	token.TYPE:        "type",
	token.VAR:         "var",
}
