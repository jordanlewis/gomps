package token

import (
	"fmt";
	"strconv";
	"strings";
)

type Token int

const (
	ILLEGAL	Token	= iota;
	EOF;

	INT;		// 345
	FLOAT;		// 4.34
	LABEL;		// foo:
	STRING;		// "foo"
	IDENT;

	REG;		// $r3
	FPREG;		// $f3
	INSTR;		// abs.d

	LPAREN;
	RPAREN;
	COMMA;

	keyword_begin;
	D_ALIGN;
	D_ASCIIZ;
	D_BYTE;
	D_DATA;
	D_SPACE;
	D_TEXT;
	D_WORD;

/* arith */
	ADD;
	ADDI;
	ADDIU;
	ADDU;
	LA;
	LI;
	LUI;
	SUB;
	SUBU;

/* logical */
	AND;
	ANDI;
	NOP;
	NOR;
	OR;
	ORI;
	XOR;
	XORI;

/* mul and div */
	DIV;
	DIVU;
	MADD;
	MADDU;
	MSUB;
	MSUBU;
	MUL;
	MULT;
	MULTU;

/* accumulator access */
	MFHI;
	MFLO;
	MTHI;
	MTLO;

/* jumps and branches */
	B;
	BAL;
	BEQ;
	BGEZ;
	BGEZAL;
	BGTZ;
	BLEZ;
	BLTZ;
	BLTZAL;
	BNE;
	BNEZ;
	J;
	JAL;
	JALR;
	JR;

/* load and store */
	LB;
	LBU;
	LH;
	LHU;
	LW;
	LWL;
	LWR;
	SB;
	SH;
	SW;
	SWL;
	SWR;
	ULW;
	USW;
	keyword_end;
)

var tokens = map[Token] string {
	ILLEGAL:	"ILLEGAL",
	EOF:	"EOF",

	INT:	"INT",
	FLOAT:	"FLOAT",
	STRING:	"STRING",
	IDENT:	"IDENT",

	LABEL:	"LABEL",
	INSTR:	"INSTR",
	REG:	"REG",
	FPREG:	"FPREG",

	LPAREN:	"(",
	RPAREN:	")",
	COMMA:	",",

	D_ALIGN:	".align",
	D_ASCIIZ:	".asciiz",
	D_BYTE:	".byte",
	D_DATA:	".data",
	D_SPACE:	".space",
	D_TEXT:	".text",
	D_WORD:	".word",

	ADD:	"add",
	ADDI:	"addi",
	ADDIU:	"addiu",
	ADDU:	"addu",
	LA:		"la",
	LI:		"li",
	LUI:	"lui",
	SUB:	"sub",
	SUBU:	"subu",

/* Logical */
	AND:	"and",
	ANDI:	"andi",
	NOP:	"nop",
	NOR:	"nor",
	OR:		"or",
	ORI:	"ori",
	XOR:	"xor",
	XORI:	"xori",

/* Mul and dIv */
	DIV:	"div",
	DIVU:	"divu",
	MADD:	"madd",
	MADDU:	"maddu",
	MSUB:	"msub",
	MSUBU:	"msubu",
	MUL:	"mul",
	MULT:	"mult",
	MULTU:	"multu",

/* AccumulatOr access */
	MFHI:	"mfhi",
	MFLO:	"mflo",
	MTHI:	"mthi",
	MTLO:	"mtlo",

/* Jumps and branches */
	B:		"b",
	BAL:	"bal",
	BEQ:	"beq",
	BGEZ:	"bgez",
	BGEZAL:	"bgezal",
	BGTZ:	"bgtz",
	BLEZ:	"blez",
	BLTZ:	"bltz",
	BLTZAL:	"bltzal",
	BNE:	"bne",
	BNEZ:	"bnez",
	J:		"j",
	JAL:	"jal",
	JALR:	"jalr",
	JR:		"jr",

/* Load and Store */
	LB:		"lb",
	LH:		"lh",
	LW:		"lw",
	SB:		"sb",
	SH:		"sh",
	SW:		"sw",

}

func (tok Token) String() string {
	if str, ok := tokens[tok]; ok {
		return str
	}
	return "token(" + strconv.Itoa(int(tok)) + ")";
}

var keywords = make(map[string]Token);

func init() {
	keywords = make(map[string]Token);
	for i := keyword_begin + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(str []byte) Token {
	if tok, ok := keywords[strings.ToLower(string(str))]; ok {
		return tok
	}
	return IDENT;
}


type Position struct {
	Filename	string;
	Offset		int;
	Line		int;
	Column		int;
}

func (pos Position) String() string {
	return fmt.Sprintf("%s @ %d:%d", pos.Filename, pos.Line, pos.Column);
}
