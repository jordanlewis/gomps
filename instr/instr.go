package instr

const ( // R = register, I = immed, A = address, L = label
	RRR = iota;
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
	NONE;
)

var InstrType = map[string] int {
/* arith */
	"add":		RRR,
	"addi":		RRI,
	"addiu":	RRI,
	"addu":		RRR,
	"clo":		RR,
	"clz":		RR,
	"la":		RL,
	"li":		RI,
	"lui":		RI,
	"move":		RR,
	"negu":		RR,
	"seb":		RR,
	"seh":		RR,
	"sub":		RRR,
	"subu":		RRR,

/* logical */
	"and":		RRR,
	"andi":		RRI,
	"nop":		NONE,
	"nor":		RRR,
	"not":		RR,
	"or":		RRR,
	"ori":		RRI,
	"xor":		RRR,
	"xori":		RRI,

/* mul and div */
	"div":		RR,
	"divu":		RR,
	"madd":		RR,
	"maddu":	RR,
	"msub":		RR,
	"msubu":	RR,
	"mul":		RRR,
	"mult":		RR,
	"multu":	RR,

/* accumulator access */
	"mfhi":		R,
	"mflo":		R,
	"mthi":		R,
	"mtlo":		R,

/* jumps and branches */
	"b":		L,
	"bal":		L,
	"beq":		RRL,
	"bgez":		RL,
	"bgezal":	RL,
	"bgtz":		RL,
	"blez":		RL,
	"bltz":		RL,
	"bltzal":	RL,
	"bne":		RRL,
	"bnez":		RL,
	"j":		A,
	"jal":		A,
	"jalr":		RR,
	"jr":		R,

/* load and store */
	"lb":		RA,
	"lbu":		RA,
	"lh":		RA,
	"lhu":		RA,
	"lw":		RA,
	"lwl":		RA,
	"lwr":		RA,
	"sb":		RA,
	"sh":		RA,
	"sw":		RA,
	"swl":		RA,
	"swr":		RA,
	"ulw":		RA,
	"usw":		RA,
}
