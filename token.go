package token

import (
    "fmt";
    "strconv";
)

type Token int

const (
	ILLEGAL	Token	= iota;
	EOF;

	INT;	// 345
	FLOAT:	// 4.34
	LABEL;	// foo:
	STRING;	// "foo"
	REG:	// $r3
	FPREG:	// $f3
)

type Position struct {
	Filename	string;
	Offset		int;
	Line		int;
	Column		int;
}
