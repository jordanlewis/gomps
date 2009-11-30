package token

import (
	"fmt";
)

type Token int

const (
	ILLEGAL	Token	= iota;
	EOF;

	INT;		// 345
	FLOAT;		// 4.34
	LABEL;		// foo:
	STRING;		// "foo"
	REG;		// $r3
	FPREG;		// $f3
	INSTR;		// abs.d
	DIRECTIVE;	//.align

	LPAREN;
	RPAREN;
)

var tokens = map[Token] string {
	ILLEGAL: "ILLEGAL",
	EOF: "EOF",

	INT: "INT",
	FLOAT: "FLOAT",
	STRING: "STRING",

	LABEL: "LABEL",
	INSTR: "INSTR",
	DIRECTIVE: "DIRECTIVE",
	REG: "REG",
	FPREG: "FPREG",

	LPAREN: "(",
	RPAREN: ")"
}


func (tok Token) String() string {
	if str, ok := tokens[tok]; ok {
		return str
	}
	return "unknown_tok"
}

type Position struct {
	Filename	string;
	Offset		int;
	Line		int;
	Column		int;
}

func (pos Position) String() string {
	return fmt.Sprintf("%s @ %d:%d", pos.Filename, pos.Line, pos.Column);
}
