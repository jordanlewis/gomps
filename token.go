package token

import (
    "fmt";
    "strconv";
)

type Token int

const (
	ILLEGAL	Token	= iota;
	EOF;
	COMMENT;

	LABEL;	// foo:
	STRING;	// "foo"
	COMMA;	// ,
	COLON;	// :
)

type Position struct {
	Filename	string;
	Offset		int;
	Line		int;
	Column		int;
}
