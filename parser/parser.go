package parser

import (
	"container/vector";
	"io";
	"fmt";
	"gomps/ast";
	"gomps/token";
	"gomps/scanner";
)

/* The grammar looks something like this:
   Statement = [Label ":"] [(instruction | "." directive) Arglist]
   Arglist = Arg ["," Arglist]
   Arg = Register | Integer */

type parser struct {
	scanner.ErrorList;
	scanner scanner.Scanner;

	trace bool;
	indent uint;

	section uint;

	pos token.Position;
	tok token.Token;
	str []byte;
}

func (p *parser) init(filename string, input []byte) {
	p.ErrorList.Init();
	p.trace = true;
	p.scanner.Init(filename, input, p);
	p.next();
}


func (p *parser) printTrace(a ...) {
	const dots = ". . . . . . . . . . . . . . . . . . . . . . . . . . . . . .";
	const n = uint(len(dots));
	fmt.Printf("%5d:%3d: ", p.pos.Line, p.pos.Column);
	i := 2 * p.indent;
	for ; i > n; i -= n {
		fmt.Print(dots)
	}
	fmt.Print(dots[0:i]);
	fmt.Println(a);
}

func trace(p *parser, str string) *parser {
	p.printTrace(str, "(");
	p.indent++;
	return p;
}

func un (p *parser) {
	p.indent--;
	p.printTrace(")");
}


func (p *parser) next() {
	p.pos, p.tok, p.str = p.scanner.Scan();
	fmt.Printf("next() returned %s: %s\n", p.tok.String(), p.str);
}

func (p *parser) expectFailed(pos token.Position, str string) {
	str = "expected " + str;
	p.Error(pos, str);
}

func (p *parser) expect(tok token.Token) token.Position {
	pos := p.pos;
	if p.tok != tok {
		p.expectFailed(p.pos, "'" + tok.String() + "'")
	}
	p.next();
	return pos;
}

//func (p *parser) parseDataDecl() *ast.DataDecl {
//	var label ast.Label;
//	switch p.tok {
//	case token.LABEL:
//		label = ast.Label{p.pos, p.str};
//	//case token.D_ALIGN:
//	//	fmt.Printf("Got an align, skipping for now\n");
//	default:
//		p.Error(p.pos, "Expected a label for a data decl");
//	}
//	p.next();
//	tok := p.tok;
//	p.next();
//	var values []ast.Lit;
//	switch tok {
//	case token.D_ASCIIZ:
//		if p.tok == token.STRING {
//			values = make([]ast.Lit, 1);
//			values[0] = ast.Lit{token.STRING, p.str};
//			d := &ast.DataDecl{tok, label, values};
//			p.next();
//			return d;
//		}
//		p.expect(token.STRING);
//	case token.D_WORD:
//		list := vector.New(0);
//		for p.tok == token.COMMA {
//			p.next();
//			list.Push(&ast.Lit{token.INT, p.str});
//			p.next();
//		}
//		values := make([]ast.Lit, list.Len());
//		for i := 0; i < list.Len(); i++ {
//			values[i] = list.At(i).(ast.Lit)
//		}
//		//d := &ast.DataDecl(tok, label, values);
//	}
//	return &ast.DataDecl{tok, label, values};
//	//return d;
//
//}

func (p *parser) parseInstr() *ast.Instr {
	var instr = new(ast.Instr);

	switch p.tok {
	case token.LABEL:
	case token.INSTR:
	}
	return instr;
}

func (p *parser) parseLabelDecl() ast.Decl {
	if p.trace { defer un(trace(p, "LabelDeclaration"))}
	pos, str := p.pos, p.str;
	p.next();
	return &ast.Labeled{ast.Label{pos, str}, parseDecl()};
}

func (p *parser) parseDecl() ast.Decl {
	if p.trace { defer un(trace(p, "Declaration"))}
	
	switch p.tok {
	case token.LABEL:
		return parseLabelDecl();
	case token.D_DATA:
		p.section = 1;
		p.next();
		return p.parseDecl();
	case token.D_TEXT:
		p.section = 0;
		p.next();
		return p.parseDecl();
	case token.D_ALIGN:
		p.next();
		p.expect(token.INT);
		return p.parseDecl();
	default:
		p.expectFailed(p.pos, "declaration");
		p.next();
		return &ast.BadDecl{p.pos};
	}
	return &ast.BadDecl{p.pos};
}

func (p *parser) Parse() []ast.Decl {
	if p.trace { defer un(trace(p, "File"))}

	list := vector.New(0);
	for p.tok != token.EOF {
		decl := p.parseDecl();
		list.Push(decl)
	}
	decls := make([]ast.Decl, list.Len());
	for i := 0; i < list.Len(); i ++ {
		decls[i] = list.At(i).(ast.Decl)
	}
	return decls;
}
//	p.expect(token.D_DATA);
//	declList := vector.New(0);
//	instrList := vector.New(0);
//	for p.tok != token.D_TEXT {
//		if p.tok == token.EOF {
//			os.Exit(1);
//		}
//		decl := p.parseDataDecl();
//		fmt.Printf("Parsed decl %s %s %s\n", decl.StorageType.String(), decl.Label.Name, decl.ValueList[0].Value);
//		declList.Push(decl);
//	}
//	decls := make([]ast.DataDecl, declList.Len());
//	for i := 0; i < declList.Len(); i++ {
//		decls[i] = declList.At(i).(ast.DataDecl)
//	}
//	p.expect(token.D_TEXT);
//	for p.tok != token.EOF {
//		instr := p.parseInstr();
//		instrList.Push(instr);
//	}
//	instrs := make([]ast.Instr, instrList.Len());
//	for i := 0; i < instrList.Len(); i++ {
//		instrs[i] = instrList.At(i).(ast.Instr)
//	}
//	return &ast.File{decls, instrs};
//}


func Parse(filename string) []ast.Decl {
	var p parser;
	input, _ := io.ReadFile(filename);
	p.init(filename, input);
	return p.Parse();
	
}



/*
type State int;

const (
	LAB_INST	State	= iota;
	INST;
	ARG;
)

type Stmt struct {
	HasLabel	bool;
	Label		[]byte;
	Type		token.Token;	// Legally either an INSTR or DIRECTIVE
	Arglist		*vector.Vector;
}

type StmtStream struct {
	List *vector.Vector;
	Cur *Stmt;
}

func (ss *StmtStream) Init() {
	ss.List = vector.New(0);
}

func (ss *StmtStream) Push(stmt *Stmt)	{
	ss.List.Push(stmt);
	ss.Cur = stmt;
}

func Parse(filename string) {
	var stmtStream StmtStream;
	stmtStream.Init();
	tokStream := token.Tokenize(filename);
	state := LAB_INST;
	var curTok *TokenData;
	for {
		curTok = tokStream.Next();
		switch state {
		case LAB_INST:
			switch curTok.tok {
			case LABEL:
				stmtStream.Push(&Stmt{HasLabel: true, Label: curTok.str});
				state = INST;
			case INSTR, DIRECTIVE:
				stmtStream.Push(&Stmt{HasLabel: false, Type: curTok.tok});
				state = ARG;
			default:
				fmt.Printf("Expected label, instr or directive, got %s\n",
					tokToString(curTok.tok));
				break;
			}
		case INST:
			switch curTok.tok {
			case LABEL:
				stmtStream.Push(&Stmt{HasLabel: true, Label: curTok.str});
				state = INST;
			case INSTR, DIRECTIVE:
				stmtStream.Cur.Type = curTok.tok;
				state = ARG;
			default:
				fmt.Printf("Expected label, instr or directive, got %s\n",
					tokToString(curTok.tok));
				break;
			}
		case ARG:
			switch curTok.tok {
			case LABEL:
				stmtStream.Push(&Stmt{HasLabel: true, Label: curTok.str});
				state = INST;
			case INSTR, DIRECTIVE:
				stmtStream.Push(&Stmt{HasLabel: false, Type: curTok.tok});
				state = ARG;
			case REG:
				//stmtStream.Cur.PushArg(&Arg{Type: REG, Value: curTok.str});
			}

		}
	}

}
*/
