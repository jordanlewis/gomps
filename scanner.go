package scanner

import (
	"bytes";
	"./token";
	"strconv";
)

func isLetter(char int) bool {
	return	('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') ||
			(char == '_')
}

func isDigit(char int) bool {
	return '0' <= char && char <= '9'
}

type Scanner struct {
	input	[]byte;

	pos		token.Position;
	offset	int;
	c	int;
}

func (S *Scanner) Init(filename string, input []byte) {
	S.input = input;
	S.pos = token.Position{filename, 0, 1, 0};
	S.offset = 0;
	S.next();
}

func (S *Scanner) next() {
	if S.offset < len(S.input) {
		S.pos.Offset = S.offset;
		S.pos.Column++;
		r := int(S.input[S.offset])
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
	switch c := S.c {
	case == '#' { // Found a comment; go to next line and restart scan
		for S.pos.Column != 0 {
			S.next()
		}
		goto restart_scan
	}
}
