package scanner

import (
	"gomps/token";
	"testing";
	"fmt";
	"io";
)


func TestScanner(t *testing.T) {
	const filename = "../tmmult.s";
	var s Scanner;
	var e ErrorList;
	var tok token.Token;
	input, _ := io.ReadFile(filename);
	e.Init();
	s.Init(filename, input, &e);
	tok = token.ILLEGAL;
	for tok != token.EOF {
		//var pos token.Position;
		//var str []byte;
		_, tok, _ = s.Scan();
		fmt.Printf("%s ", tok.String());
		//fmt.Printf("%s@%d:%d %s %s\n", pos.Filename, pos.Line, pos.Column, tok.String(), str);
	}
	fmt.Printf("\n");
}
