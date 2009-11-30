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
    spewTokens("tmmult.s");
}
