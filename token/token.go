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
	DIRECTIVE;	// one of the below .directives

	D_ALIGN;
	D_ASCIIZ;
	D_BYTE;
	D_DATA;
	D_SPACE;
	D_TEXT;
	D_WORD;

	LPAREN;
	RPAREN;
	COMMA;
)

var tokens = map[Token] string {
	ILLEGAL: "ILLEGAL",
	EOF: "EOF",

	INT: "INT",
	FLOAT: "FLOAT",
	STRING: "STRING",

	LABEL: "LABEL",
	INSTR: "INSTR",
	REG: "REG",
	FPREG: "FPREG",

	DIRECTIVE: "DIRECTIVE",
	D_ALIGN: ".ALIGN",
	D_ASCIIZ: ".ASCIIZ",
	D_BYTE: ".BYTE",
	D_DATA: ".DATA",
	D_SPACE: ".SPACE",
	D_TEXT: ".TEXT",
	D_WORD: ".WORD",

	LPAREN: "(",
	RPAREN: ")"
}

var Directives = map[string] Token {
	".align": D_ALIGN,
	".asciiz": D_ASCIIZ,
	".byte": D_BYTE,
	".data": D_DATA,
	".space": D_SPACE,
	".text": D_TEXT,
	".word": D_WORD,
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
