package scanner

import (
	"fmt";
	"io";
	"container/vector";
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


type TokenData struct {
	pos	token.Position;
	tok	token.Token;
	str	[]byte;
}

type TokenStream struct {
	list	*vector.Vector;
	curTok	int;
}

type Scanner struct {
	input	[]byte;

	pos	token.Position;
	offset	int;
	c	int;
}

func (T *TokenStream) Init() {
	T.list = vector.New(0);
	T.curTok = 0;
}

func (T *TokenStream) Push(td *TokenData)	{ T.list.Push(td) }

func (T *TokenStream) Len() int { return T.list.Len()}

func (T *TokenStream) Next() *TokenData {
	T.curTok++;
	return T.list.At(T.curTok - 1).(*TokenData);
}

func (S *Scanner) Init(filename string, input []byte) {
	S.input = input;
	S.pos = token.Position{filename, 0, 1, 0};
	S.offset = 0;
	S.next();
}

func (S *Scanner) scanIdentifier() token.Token {
	tok := token.INSTR;
	if S.c == '.' {
		tok = token.DIRECTIVE
	}
	for isLetter(S.c) || isDigit(S.c) {
		S.next()
	}
	if S.c == ':' {
		S.next();
		if tok == token.DIRECTIVE {
			tok = token.ILLEGAL
		} else {
			tok = token.LABEL
		}
	}
	return tok;
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
		fmt.Printf("Illegal character escape\n")
	}
}

func (S *Scanner) scanString() {
	for S.c != '"' {
		if S.c == '\n' || S.c < 0 {
			fmt.Printf("String unterminated\n");
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
	for S.c == ' ' || S.c == ',' || S.c == '\t' || S.c == '\n' || S.c == '\r' {
		S.next()
	}

	pos, tok = S.pos, token.ILLEGAL;
	startOffset, endOffset := 0, 0;

	switch c := S.c; {
	case isLetter(c):
		tok = S.scanIdentifier();
		if tok == token.LABEL {
			endOffset = -1; // Chop off the :
		}
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
			startOffset = 1; // Chop off the $
			S.scanReg();
		case '(':
			tok = token.LPAREN
		case ')':
			tok = token.RPAREN

		case '#':	// Found a comment; go to next line and restart scan
			for S.pos.Column != 0 {
				S.next()
			}
			goto restart_scan;
		default:
			fmt.Printf("Illegal char %c\n", c);
			tok = token.EOF;
		}
	}
	return pos, tok, S.input[pos.Offset + startOffset:S.pos.Offset + endOffset];
}

// func Tokenize(filename string, input []byte) int {
// 	var s Scanner;
// 	s.Init(filename, input);
// 	for f(s.Scan()) {
//
// 	}
// 	return 0;
// }
func Tokenize(filename string) *TokenStream {
	var stream TokenStream;
	var s Scanner;
	var t *TokenData;
	stream.Init();
	input, _ := io.ReadFile(filename);
	s.Init(filename, input);
	tok := token.ILLEGAL;
	for tok != token.EOF {
		pos, tok, word := s.Scan();
		t = &TokenData{pos, tok, word};
		stream.Push(t);
		fmt.Printf("%s@%d(%d:%d) %s %s\n", pos.Filename, pos.Offset, pos.Line, pos.Column, tok.String(), word);
	}
	return &stream;
}
