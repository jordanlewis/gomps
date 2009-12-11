package inst

import (
	"gomps/token";
	"fmt";
)

type ArgType int;
type InstType int;

type RegT uint8;
type LabelT struct {
	Section int;
	Offset uint16;
}
func (l *LabelT) String() string {
	var secstr string;
	if l.Section == 0 {
		secstr = "text"
	} else {
		secstr = "data"
	}
	return fmt.Sprintf("%s:%d", secstr, l.Offset);
}

type ImmT uint32;


type Inst struct {
	Opname token.Token;
	RS RegT;
	RT RegT;
	RD RegT;
	SA uint8;
	IMM int;
	TGT *LabelT;
}

func (i *Inst) String() string {
	return fmt.Sprintf("inst(%s: %d %d %d %d %d %s)", i.Opname.String(), i.RS, i.RT, i.RD, i.SA, i.IMM, i.TGT.String());
}

var Regmap = map[string] RegT {
	"$zero": 0,
	"$at":1,
	"$v0":2, "$v1":3,
	"$a0":4, "$a1":5, "$a2":6, "$a3":7,
	"$t0":8, "$t1":9, "$t2":10, "$t3":11, "$t4":12, "$t5":13, "$t6":14,"$t7":15,
	"$s0":16,"$s1":17,"$s2":18, "$s3":19, "$s4":20, "$s5":21, "$s6":22,"$s7":23,
	"$t8":24, "$t9":25,
	"$k0":26, "$k1":27,
	"$gp":28, "$sp":29, "$fp":30, "$ra":31
}


func Regnum(str []byte) RegT {
	return Regmap[string(str)];
}

const ( // R = register, I = immed, A = address, L = label, S = string
	INVALID = iota;

	//r_types
	RRR;
	RR;
	R;
	RRI;

	//i_types
	RRL;
	RL;
	RIR;
	RI;

	//j_types
	L;
)

const (
	ARITH = iota;
	LOSTO;
	BRANCH;
	DIRECT;
)

var IType = map[token.Token] InstType {
	token.ADD: ARITH,
	token.ADDI: ARITH,
	token.ADDIU: ARITH,
	token.ADDU: ARITH,
	token.LUI: ARITH,
	token.SUB: ARITH,
	token.SUBU: ARITH,

/* Logical */
	token.AND: ARITH,
	token.ANDI: ARITH,
	token.NOP: ARITH,
	token.NOR: ARITH,
	token.OR: ARITH,
	token.ORI: ARITH,
	token.XOR: ARITH,
	token.XORI: ARITH,

/* Mul and dIv */
	token.DIV: ARITH,
	token.DIVU: ARITH,
	token.MADD: ARITH,
	token.MADDU: ARITH,
	token.MSUB: ARITH,
	token.MSUBU: ARITH,
	token.MUL: ARITH,
	token.MULT: ARITH,
	token.MULTU: ARITH,

/* AccumulatOr access */
	token.MFHI: ARITH,
	token.MFLO: ARITH,
	token.MTHI: ARITH,
	token.MTLO: ARITH,

/* Jumps and branches */
	token.BEQ: BRANCH,
	token.BGEZ: BRANCH,
	token.BGEZAL: BRANCH,
	token.BGTZ: BRANCH,
	token.BLEZ: BRANCH,
	token.BLT: BRANCH,
	token.BLTZ: BRANCH,
	token.BLTZAL: BRANCH,
	token.BNE: BRANCH,
	token.BNEZ: BRANCH,
	token.J: BRANCH,
	token.JAL: BRANCH,
	token.JALR: BRANCH,
	token.JR: BRANCH,

/* Load and Store */
	token.LA: LOSTO,
	token.LB: LOSTO,
	token.LH: LOSTO,
	token.LW: LOSTO,
	token.SB: LOSTO,
	token.SH: LOSTO,
	token.SW: LOSTO,

	token.D_ALIGN:	DIRECT,
	token.D_ASCIIZ:	DIRECT,
	token.D_BYTE:	DIRECT,
	token.D_DATA:	DIRECT,
	token.D_SPACE:	DIRECT,
	token.D_TEXT:	DIRECT,
	token.D_WORD:	DIRECT,
}

var AType = map[token.Token] ArgType {
	//token.D_ALIGN:	I,
	//token.D_ASCIIZ:	S,
	//token.D_BYTE:	M,
	//token.D_DATA:	M,
	//token.D_SPACE:	I,
	//token.D_TEXT:	M,
	//token.D_WORD:	M,

/* arith */
	token.ADD:		RRR,
	token.ADDI:		RRI,
	token.ADDIU:	RI,
	token.ADDU:		RRR,
	token.LUI:		RI,
	token.SUB:		RRR,
	token.SUBU:		RRR,

/* logical */
	token.AND:		RRR,
	token.ANDI:		RRI,
	//token.NOP:		NONE,
	token.NOR:		RRR,
	token.OR:		RRR,
	token.ORI:		RRI,
	token.XOR:		RRR,
	token.XORI:		RRI,

/* mul and div */
	token.DIV:		RR,
	token.DIVU:		RR,
	token.MADD:		RR,
	token.MADDU:	RR,
	token.MSUB:		RR,
	token.MSUBU:	RR,
	token.MUL:		RRR,
	token.MULT:		RR,
	token.MULTU:	RR,

/* accumulators */
	token.MFHI:		R,
	token.MFLO:		R,
	token.MTHI:		R,
	token.MTLO:		R,

/* jumps and branches */
	token.BEQ:		RRL,
	token.BGEZ:		RL,
	token.BGEZAL:	RL,
	token.BGTZ:		RL,
	token.BLEZ:		RL,
	token.BLT:		RRL,
	token.BLTZ:		RL,
	token.BLTZAL:	RL,
	token.BNE:		RRL,
	token.BNEZ:		RL,
	token.J:		L,
	token.JAL:		L,
	token.JALR:		RR,
	token.JR:		R,

/* load and store */
	token.LA:		RL,
	token.LB:		RIR,
	token.LH:		RIR,
	token.LW:		RIR,
	token.SB:		RIR,
	token.SH:		RIR,
	token.SW:		RIR,
}
