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
	datalist *vector.Vector;

	Instlabs map[string] inst.LabelT;
	Datalabs map[string] inst.LabelT;

	pos t.Position;
	tok t.Token;
	str []byte;
}

func (p *parser) init(filename string, input []byte) {
	p.ErrorList.Init();
	p.trace = true;
	p.scanner.Init(filename, input, p);
	p.Instlist = vector.New(0);
	p.datalist = vector.New(0);
	p.Instlabs = make(map[string] inst.LabelT);
	p.Datalabs = make(map[string] inst.LabelT);
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
			fmt.Printf("label is %s %s\n", p.str, string(p.str[0:len(p.str)-1]));
			if p.section == 0 {
				p.Instlabs[string(p.str[0:len(p.str)-1])] = inst.LabelT(p.Instlist.Len());
			} else {
				p.Datalabs[string(p.str[0:len(p.str)-1])] = inst.LabelT(p.datalist.Len());
			}
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
					p.datalist.Push(i);
				} else {
					p.Error(p.pos, "Nonint in .word");
				}
				first_time = false;
			}
		case t.D_ASCIIZ, t.D_SPACE:
			p.next();
			continue;
		default:
			newInst := &inst.Inst{p.tok, 0,0,0,0,0,0};
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
				p.expect(t.REG);
				switch p.tok {
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
				newInst.TGT = p.Instlabs[string(p.str)];
			case inst.RL:
				p.expect(t.REG);
				newInst.RS = inst.Regnum(p.str);
				p.expect(t.COMMA);
				p.expect(t.IDENT);
				newInst.TGT = p.Instlabs[string(p.str)];
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
				newInst.TGT = p.Instlabs[string(p.str)];
			default:
				fmt.Printf("Ignoring token %s\n",p.str); 
				p.next();
			}
			p.Instlist.Push(newInst);
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
