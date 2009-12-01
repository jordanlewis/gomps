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
	if tok == token.DIRECTIVE {
		dirStr := fmt.Sprintf("%s", S.input[pos.Offset:S.pos.Offset]);
		if t, ok := token.Directives[dirStr]; ok { tok = t }
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
	for S.c == ' ' || S.c == ',' || S.c == '\t' || S.c == '\n' || S.c == '\r' {
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

func Tokenize(filename string) *TokenStream {
	var stream TokenStream;
	var s Scanner;
	var t *TokenData;
	var e ErrorList;
	e.Init();
	stream.Init();
	input, _ := io.ReadFile(filename);
	s.Init(filename, input, &e);
	tok := token.ILLEGAL;
	for tok != token.EOF {
		pos, tokn, word := s.Scan();
		tok = tokn;
		t = &TokenData{pos, tok, word};
		stream.Push(t);
		fmt.Printf("%s@%d(%d:%d) %s %s\n", pos.Filename, pos.Offset, pos.Line, pos.Column, tok.String(), word);
	}
	return &stream;
}
