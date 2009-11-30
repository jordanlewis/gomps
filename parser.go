package gomps

import (
	"container/vector";
	"fmt";
)

/* The grammar looks something like this:
   Statement = [Label ":"] [(instruction | "." directive) Arglist]
   Arglist = Arg ["," Arglist]
   Arg = Register | Integer */

type State int

const (
	LAB_INST	State	= iota;
	INST;
	ARG;
)

type Stmt struct {
	HasLabel	bool;
	Label		[]byte;
	Type		Token;	// Legally either an INSTR or DIRECTIVE
	Arglist		*vector.Vector;
}

type StmtStream struct {
	List *vector.Vector;
	Cur *Stmt;
}

func (ss *StmtStream) Init()	{ ss.List = vector.New(0) }

func (ss *StmtStream) Push(stmt *Stmt)	{
	ss.List.Push(stmt);
	ss.Cur = stmt;
}

func Parse(filename string) {
	var stmtStream StmtStream;
	stmtStream.Init();
	tokStream := Tokenize(filename);
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
