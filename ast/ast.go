package ast

import (
	"gomps/token";
)

/*
Prog := (Exp)*

Exp := Opname Arglist

Arglist :=
		 | Arg ("," Arg)*
		 | Arg "," Arg "(" Arg ")"

Arg := Reg
	 | Immed
	 | Label

*/


type MIPSInst interface {
	
}
type Inst struct {
	Opcode token.Token;
	Label string;
	Reg1 uint32;
	Reg2 uint32;
	Reg3 uint32;
}


type Lit struct {
	Kind token.Token; // INT, FLOAT, or STRING
	Value []byte;
}

type Label struct {
	Pos token.Position;
	Name []byte;
}

type Arg interface {
	isArg();
}

type (
	Reg struct {
		Name []byte;
	};
	Immed struct {
		Value Lit;
	};
	Addr struct {
		Address Label;
	};
	IndAddr struct {
		Base Reg;
		Offset int;
	};
)

func (x *Reg) isArg()		{}
func (x *Immed) isArg()		{}
func (x *Addr) isArg()		{}
func (x *IndAddr) isArg()	{}

type Decl interface {
	isDecl();
}

type (
	Labeled struct {
		Label Label;
		Decl Decl;
	};
	Instr struct {
		Instr token.Token;
		Args []Arg;
	};
	DataDecl struct {
		StorageType	token.Token; // .word, .byte, .space, ...?
		ValueList	[]Lit;
	};
	BadDecl struct {
		token.Position;
	};
)
func (x *Labeled) isDecl()	{}
func (x *Instr) isDecl()	{}
func (x *DataDecl) isDecl()	{}
func (x *BadDecl) isDecl()	{}

type File struct {
	DataDecls	[]DataDecl;		// variable declarations
	Instrs		[]Instr;		// instructions
}
