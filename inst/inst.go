package inst

import (
	"gomps/token";
	"fmt";
)

type ArgType int;
type InstType int;

type RegT uint8;
type LabelT uint16;
type ImmT uint32;

/*
type Inst interface {
	getOpname() token.Token;
	getRS();
	getRT();
	getRD();
	getSA();
	getIMM();
	getTGT();
}
*/

type Inst struct {
	Opname token.Token;
	RS RegT;
	RT RegT;
	RD RegT;
	SA uint8;
	IMM uint32;
	TGT LabelT;
}

func (i *Inst) String() string {
	return fmt.Sprintf("inst(%s: %d %d %d %d %d %d)", i.Opname.String(), i.RS, i.RT, i.RD, i.SA, i.IMM, i.TGT);
}

type RInst struct {
	Opname token.Token;
	RS RegT;
	RT RegT;
	RD RegT;
	SA uint8;
}

type IInst struct {
	Opname token.Token;
	RS RegT;
	RT RegT;
	IMM ImmT;
}

type JInst struct {
	Opname token.Token;
	TGT LabelT;
}

func (x *RInst) getOpname() token.Token	{return x.Opname}
func (x *IInst) isInst()	{}
func (x *JInst) isInst()	{}


const ( // R = register, I = immed, A = address, L = label, S = string
	INVALID = iota;
	RRR;
	RRI;
	RRO;
	ROR;
	RRL;
	RIL;
	RR;
	RI;
	RA;
	RL;
	R;
	A;
	L;
	S;
	M;
	I;
	NONE;
)

const (
	ARITH = iota;
	LOSTO;
	BRANCH;
)

var IType = map[token.Token] InstType {
	token.ADD: ARITH,
	token.ADDI: ARITH,
	token.ADDIU: ARITH,
	token.ADDU: ARITH,
	token.LA: ARITH,
	token.LI: ARITH,
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
	token.B: BRANCH,
	token.BAL: BRANCH,
	token.BEQ: BRANCH,
	token.BGEZ: BRANCH,
	token.BGEZAL: BRANCH,
	token.BGTZ: BRANCH,
	token.BLEZ: BRANCH,
	token.BLTZ: BRANCH,
	token.BLTZAL: BRANCH,
	token.BNE: BRANCH,
	token.BNEZ: BRANCH,
	token.J: BRANCH,
	token.JAL: BRANCH,
	token.JALR: BRANCH,
	token.JR: BRANCH,

/* Load and Store */
	token.LB: LOSTO,
	token.LH: LOSTO,
	token.LW: LOSTO,
	token.SB: LOSTO,
	token.SH: LOSTO,
	token.SW: LOSTO,
}

var Instrs = map[token.Token] ArgType {
	token.D_ALIGN:	I,
	token.D_ASCIIZ:	S,
	token.D_BYTE:	M,
	token.D_DATA:	M,
	token.D_SPACE:	I,
	token.D_TEXT:	M,
	token.D_WORD:	M,

/* atoken.rith */
	token.ADD:		RRR,
	token.ADDI:		RRI,
	token.ADDIU:	RI,
	token.ADDU:		RRR,
	token.LA:		RL,
	token.LI:		RI,
	token.LUI:		RI,
	token.SUB:		RRR,
	token.SUBU:		RRR,

/* ltoken.ogical */
	token.AND:		RRR,
	token.ANDI:		RRI,
	token.NOP:		NONE,
	token.NOR:		RRR,
	token.OR:		RRR,
	token.ORI:		RRI,
	token.XOR:		RRR,
	token.XORI:		RRI,

/* mtoken.ul and div */
	token.DIV:		RR,
	token.DIVU:		RR,
	token.MADD:		RR,
	token.MADDU:	RR,
	token.MSUB:		RR,
	token.MSUBU:	RR,
	token.MUL:		RRR,
	token.MULT:		RR,
	token.MULTU:	RR,

/* atoken.ccumulators */
	token.MFHI:		R,
	token.MFLO:		R,
	token.MTHI:		R,
	token.MTLO:		R,

/* jtoken.umps and branches */
	token.B:		L,
	token.BAL:		L,
	token.BEQ:		RRL,
	token.BGEZ:		RL,
	token.BGEZAL:	RL,
	token.BGTZ:		RL,
	token.BLEZ:		RL,
	token.BLTZ:		RL,
	token.BLTZAL:	RL,
	token.BNE:		RRL,
	token.BNEZ:		RL,
	token.J:		A,
	token.JAL:		A,
	token.JALR:		RR,
	token.JR:		R,

/* ltoken.oad and store */
	token.LB:		RA,
	token.LH:		RA,
	token.LW:		RA,
	token.SB:		RA,
	token.SH:		RA,
	token.SW:		RA,
}
