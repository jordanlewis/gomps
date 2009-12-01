package scanner

import (
	"testing";
	"fmt";
)


func TestScanner(t *testing.T) {
	stream := Tokenize("../tmmult.s");
	var td *TokenData;
	for i := 0; i < stream.Len(); i++ {
		td = stream.Next();
		fmt.Printf("%s@%d(%d:%d) %s %s\n", td.pos.Filename, td.pos.Offset, td.pos.Line, td.pos.Column, td.tok.String(), td.str);
	}
}
