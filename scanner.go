package gomps

/*
import (
	"bytes";
	"strconv";
)
*/
import (
	"fmt";
	"io";
)

func reportError(err string, pos Position) {
	fmt.Printf("Error in scanning: %s.%s:%s): %s\n", pos.Filename,
			   pos.Line, pos.Column, err);
}

func isLetter(char int) bool {
	return	('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') ||
			(char == '_') || (char == '.')
}

func isDigit(char int) bool {
	return '0' <= char && char <= '9'
}

type Scanner struct {
	input	[]byte;

	pos		Position;
	offset	int;
	c	int;
}

func (S *Scanner) Init(filename string, input []byte) {
	S.input = input;
	S.pos = Position{filename, 0, 1, 0};
	S.offset = 0;
	S.next();
}

func (S *Scanner) scanIdentifier() Token {
	tok := INSTR;
	if S.c == '.' {
		tok = DIRECTIVE;
	}
	for isLetter(S.c) || isDigit(S.c) {
		S.next()
	}
	if S.c == ':' {
		S.next();
		if tok == DIRECTIVE {
			tok = ILLEGAL;
		} else {
			tok = LABEL;
		}
	}
	return tok;
}

func (S *Scanner) scanNumber() Token {
	for isDigit(S.c) {
		S.next()
	}
	return INT;
}

func (S *Scanner) scanEscape() {
	S.next();
	switch S.c {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
	default:
		fmt.Printf("Illegal character escape\n");
	}
}

func (S *Scanner) scanString() {
	for S.c != '"' {
		if S.c == '\n' || S.c < 0 {
			fmt.Printf("String unterminated\n");
			break;
		}
		if S.c == '\\' {
			S.scanEscape();
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

func (S *Scanner) Scan() (pos Position, tok Token, word []byte) {
restart_scan:
	for S.c == ' ' || S.c == ',' || S.c == '\t' || S.c == '\n' || S.c == '\r' {
		S.next()
	}

	pos, tok = S.pos, ILLEGAL;

	switch c := S.c; {
	case isLetter(c):
		tok = S.scanIdentifier()
	case isDigit(c):
		tok = S.scanNumber()
	default:
		S.next();
		switch c {
		case -1:
			tok = EOF
		case '"':
			tok = STRING;
			S.scanString();
		case '$':
			tok = REG;
			S.scanReg();
		case '(':
			tok = LPAREN
		case ')':
			tok = RPAREN

		case '#': // Found a comment; go to next line and restart scan
			for S.pos.Column != 0 {
				S.next()
			}
			goto restart_scan
		default:
			fmt.Printf("Illegal char %c\n", c);
			tok = EOF;
		}
	}
	return pos, tok, S.input[pos.Offset:S.pos.Offset];
}

// func Tokenize(filename string, input []byte) int {
// 	var s Scanner;
// 	s.Init(filename, input);
// 	for f(s.Scan()) {
// 		
// 	}
// 	return 0;
// }
func Tokenize(filename string) {
	var s Scanner;
	input, _ := io.ReadFile(filename);
	s.Init(filename, input);
	token := ILLEGAL;
	for token != EOF {
		pos, tok, word := s.Scan();
		token = tok;
		fmt.Printf("%s@%d(%d:%d) %s %s\n", pos.Filename, pos.Offset, pos.Line, pos.Column, tokToString(token), word);
	}
}
