package scanner

import (
	"bytes";
	"./token";
	"strconv";
)

type Scanner struct {
	input	[]byte;

	pos		token.Position;
	offset	int;
	char	int;
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
		S.char = r;
	} else {
		S.pos.Offset = len(S.input);
		S.char = -1;
	}
}

func (S *Scanner) Scan() (pos token.Position, tok token.Token, word []byte) {
restart_scan:
	for S.char == ' ' || S.char == '\t' || S.char == '\n' || S.char == '\r' {
		S.next()
	}
	if S.ch == '#' { // Found a comment; go to next line and restart scan
		for S.pos.Column != 0 {
			S.next()
		}
		goto restart_scan
	}
}
