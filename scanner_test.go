package gomps

import ("testing"; "io"; "fmt";)

func spewTokens(filename string) {
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
func TestScanner(t *testing.T) {
	l := Tokenize("tmmult.s");
	var td TokenData;
	for i := 0; i < l.Len(); i++ {
		td = l.At(i).(TokenData);
		fmt.Printf("%s@%d(%d:%d) %s %s\n", td.pos.Filename, td.pos.Offset, td.pos.Line, td.pos.Column, tokToString(td.tok), td.str);
	}
}
