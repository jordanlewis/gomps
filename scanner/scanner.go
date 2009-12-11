package scanner

import (
	"fmt";
	"gomps/token";
)

func reportError(err string, pos token.Position) {
	fmt.Printf("Error in scanning: %s.%s:%s): %s\n", pos.Filename,
		pos.Line, pos.Column, err)
}

func isLetter(char int) bool {
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') ||
		(char == '_') || (char == '.')
}

func isDigit(char int) bool	{ return '0' <= char && char <= '9' }

type Scanner struct {
	input	[]byte;
	err		ErrorHandler;

	pos	token.Position;
	offset	int;
	c	int;
}

func (S *Scanner) Init(filename string, input []byte, err ErrorHandler) {
	S.input = input;
	S.pos = token.Position{filename, 0, 1, 0};
	S.err = err;
	S.offset = 0;
	S.next();
}

func (S *Scanner) error(pos token.Position, str string) {
	S.err.Error(pos, str);
}

func (S *Scanner) scanIdentifier() token.Token {
	pos := S.pos;
	for isLetter(S.c) || isDigit(S.c) {
		S.next()
	}
	if S.c == ':' {
		S.next();
		return token.LABEL;
	}

	return token.Lookup(S.input[pos.Offset:S.pos.Offset]);
}

func (S *Scanner) scanNumber() token.Token {
	for isDigit(S.c) {
		S.next()
	}
	return token.INT;
}

func (S *Scanner) scanEscape() {
	S.next();
	switch S.c {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
	default:
		S.error(S.pos, "Illegal character escape")
	}
}

func (S *Scanner) scanString() {
	for S.c != '"' {
		if S.c == '\n' || S.c < 0 {
			S.error(S.pos, "String unterminated");
			break;
		}
		if S.c == '\\' {
			S.scanEscape()
		}
		S.next();
	}

	S.next();
}

func (S *Scanner) scanReg() {
	for isLetter(S.c) || isDigit(S.c) {
		S.next()
	}
}

func (S *Scanner) next() {
	if S.offset < len(S.input) {
		S.pos.Offset = S.offset;
		S.pos.Column++;
		r := int(S.input[S.offset]);
		if r == '\n' {
			S.pos.Line++;
			S.pos.Column = 0;
		}
		S.offset += 1;
		S.c = r;
	} else {
		S.pos.Offset = len(S.input);
		S.c = -1;
	}
}

func (S *Scanner) Scan() (pos token.Position, tok token.Token, word []byte) {
restart_scan:
	for S.c == ' ' || S.c == '\t' || S.c == '\n' || S.c == '\r' {
		S.next()
	}

	pos, tok = S.pos, token.ILLEGAL;

	switch c := S.c; {
	case isLetter(c):
		tok = S.scanIdentifier();
	case isDigit(c):
		tok = S.scanNumber()
	default:
		S.next();
		switch c {
		case -1:
			tok = token.EOF
		case '"':
			tok = token.STRING;
			S.scanString();
		case '$':
			tok = token.REG;
			S.scanReg();
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN
		case ',':
			tok = token.COMMA

		case '#':	// Found a comment; go to next line and restart scan
			for S.pos.Column != 0 {
				S.next()
			}
			goto restart_scan;
		default:
			S.error(S.pos, "Illegal character");
			tok = token.EOF;
		}
	}
	return pos, tok, S.input[pos.Offset:S.pos.Offset];
}
