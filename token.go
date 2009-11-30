package gomps

/*
import (
    "fmt";
    "strconv";
)
*/

type Token int

const (
	ILLEGAL	Token	= iota;
	EOF;

	INT;	// 345
	FLOAT;	// 4.34
	LABEL;	// foo:
	STRING;	// "foo"
	REG;	// $r3
	FPREG;	// $f3
	INSTR;	// abs.d
	DIRECTIVE;	//.align

	LPAREN;
	RPAREN;
)

type Position struct {
	Filename	string;
	Offset		int;
	Line		int;
	Column		int;
}

func tokToString(tok Token) string{
	var str string;
	switch tok {
	case ILLEGAL: str = "ILLEGAL";
	case EOF: str = "EOF";
	case INT: str = "INT";
	case FLOAT: str = "FLOAT";
	case LABEL: str = "LABEL";
	case STRING: str = "STRING";
	case REG: str = "REG";
	case FPREG: str = "FPREG";
	case INSTR: str = "INSTR";
	case DIRECTIVE: str = "DIRECTIVE";
	case LPAREN: str = "LPAREN";
	case RPAREN: str = "RPAREN";
	}
	return str;
}
