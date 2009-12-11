package parser

import (
	"container/vector";
	"os";
	"io";
	"fmt";
	t "gomps/token";
	"gomps/scanner";
	"gomps/inst";
	"gomps/debug";
	"strconv";
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

	Instlist *vector.Vector;
	Datalist *vector.Vector;
	Memory [1024]uint32;
	CurMemPos uint16;

	Labels map[string] *inst.LabelT;

	pos t.Position;
	tok t.Token;
	str []byte;
}

func (p *parser) init(filename string, input []byte) {
	p.ErrorList.Init();
	p.trace = true;
	p.scanner.Init(filename, input, p);
	p.Instlist = vector.New(0);
	p.Datalist = vector.New(0);
	p.Labels = make(map[string] *inst.LabelT);
	p.next();
}

func (p *parser) next() {
	p.pos, p.tok, p.str = p.scanner.Scan();
	debug.Debug("next() returned %s: %s\n", p.tok.String(), p.str);
}

func (p *parser) expectFailed(pos t.Position, str string) {
	str = "expected " + str;
	p.Error(pos, str);
}

func (p *parser) expect(tok t.Token) {
	p.next();
	if p.tok != tok {
		p.expectFailed(p.pos, "'" + tok.String() + "'")
	}
}

func (p *parser) Parse() {
	for p.tok != t.EOF {
		switch p.tok {
		case t.D_DATA:
			p.section = 1;
			p.next();
		case t.D_TEXT:
			p.section = 0;
			p.next();
		case t.LABEL:
			fmt.Printf("label is %s %s\n",p.str, string(p.str[0:len(p.str)-1]));
			var lab *inst.LabelT;
			if p.section == 0 {
				lab = &inst.LabelT{0, uint16(p.Instlist.Len())};
			} else {
				lab = &inst.LabelT{1, p.CurMemPos};
			}
			p.Labels[string(p.str[0:len(p.str)-1])] = lab;
			p.next();
		case t.D_WORD:
			first_time := true;
			reading := true;

			for reading == true {
				p.next();
				if !first_time {
					if p.tok != t.COMMA {
						reading = false;
						break;
					}
					p.next();
				}
				if p.tok == t.INT {
					i, _ := strconv.Atoi(string(p.str));
					p.Memory[p.CurMemPos] = uint32(i);
					p.CurMemPos += 4;
				} else {
					p.Error(p.pos, "Nonint in .word");
				}
				first_time = false;
			}
		case t.D_ASCIIZ, t.D_SPACE:
			p.next();
			continue;
		default:
			var pushInst bool = true;
			newInst := &inst.Inst{p.tok,0,0,0,0,0,&inst.LabelT{}};
			debug.Debug("Trying to parse a %s: ", p.str);
			switch ty, _ := inst.AType[p.tok]; ty {
			case inst.RRR:
				p.expect(t.REG);
				newInst.RD = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
			case inst.RR:
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
			case inst.R:
				tok := p.tok;
				p.expect(t.REG);
				switch tok {
				case t.MFHI, t.MFLO:
					newInst.RD = inst.Regnum(p.str);
				case t.MTHI, t.MTLO, t.JR:
					newInst.RS = inst.Regnum(p.str);
				}
			case inst.RRI:
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.INT);
				newInst.IMM, _ = strconv.Atoi(string(p.str));
			case inst.RRL:
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.IDENT);
				newInst.TGT = p.Labels[string(p.str)];
			case inst.RL:
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.IDENT);
				newInst.TGT = p.Labels[string(p.str)];
			case inst.RIR:
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.INT);
				newInst.IMM, _ = strconv.Atoi(string(p.str));
				p.expect(t.LPAREN);
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.RPAREN);
			case inst.RI:
				p.expect(t.REG);
				newInst.RT = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.INT);
				newInst.IMM, _ = strconv.Atoi(string(p.str));
			case inst.L:
				p.expect(t.IDENT);
				newInst.TGT = p.Labels[string(p.str)];
			default:
				pushInst = false;
				fmt.Printf("Ignoring token %s\n",p.str); 
				p.next();
			}
			if pushInst {
				p.Instlist.Push(newInst);
			}
			p.next();
		}
	}
}


func Parse(filename string) (*parser, os.Error) {
	var p parser;
	input, err := io.ReadFile(filename);
	p.init(filename, input);
	p.Parse();
	return &p, err;
}
