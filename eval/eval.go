package eval

import (
	t "gomps/token";
	"gomps/debug";
	"gomps/inst";
	"os";
	"fmt";
	"container/vector";
)


type CPU struct {
	cycles int;
	Regs [32]uint32;
	pc uint32;
	Data *vector.Vector;
	Instrs *vector.Vector;
	Pipeline *Pipeline;
	Mem [1024]uint32;
	lo_reg uint32;
	hi_reg uint32;
}

func (c *CPU) printRegs() {
	for s, i := range inst.Regmap {
		fmt.Printf("%s (%d) = %d\n", s, i, c.Regs[int(i)]);
	}
}

func (c *CPU) printMem() {
	for i := 0; i < 30; i++ {
		fmt.Printf("mem[%d] = %d\n", i*4, c.Mem[i*4]);
	}
}

type Pipeline struct {
	cpu *CPU;
	ready bool;
	ifid ifid_buf;
	idex idex_buf;
	exmem exmem_buf;
	memwb memwb_buf;
}

type ifid_buf struct { /* Buffer between ifetch and decode */
	full bool;
	npc uint32;
	inst *inst.Inst;
}

func (i *ifid_buf) String() string {
	return fmt.Sprintf("ifid: npc = %d, inst = %s", i.npc, i.inst.String());
}

type idex_buf struct { /* Buffer between decode and execute */
	full bool;
	npc uint32;
	reg_a uint32;
	reg_b uint32;
	imm uint32;
	inst *inst.Inst;
}

func (i *idex_buf) String() string {
	return fmt.Sprintf("idex: npc = %d, reg_a = %d, reg_b = %d, imm = %d, inst = %s", i.npc, i.reg_a, i.reg_b, i.imm, i.inst.String());
}

type exmem_buf struct { /* Buffer between execute and memory access */
	full bool;
	npc uint32;
	cond bool;
	alu_out uint32;
	reg_b uint32;
	inst *inst.Inst;
}

func (e *exmem_buf) String() string {
	return fmt.Sprintf("exmem: npc = %d, cond = %d, alu_out = %d, reg_b = %d, inst = %s", e.npc, e.cond, e.alu_out, e.reg_b, e.inst.String());
}

type memwb_buf struct { /* Buffer between memory access and writeback */
	was_full bool; /* just for printing functions */
	full bool;
	lmd uint32;
	alu_out uint32;
	inst *inst.Inst;
}

func (m *memwb_buf) String() string {
	return fmt.Sprintf("memwb: lmd = %d, alu_out = %d, inst = %s", m.lmd, m.alu_out, m.inst.String());
}

func (p *Pipeline) Dispatch() {
	/*
	if p.ins_reg != nil {
		debug.Debug("Not dispatching when IFETCH is in progress\n");
		return;
	}
	*/
	if !p.ready {
		debug.Debug("Dispatching new instruction\n");
		p.ready = true;
	}
	//debug.Debug("Dispatch: from pc = %d, next inst = %s\n", p.cpu.pc, p.ins_reg.String());
}

func (p *Pipeline) Step() {
	/* ifetch already happened due to dispatch. we start at the end of the
	   pipeline and work our way backwards, to avoid the usage of temporary
	   pipeline latches. */

	/* writeback:  */
	if p.memwb.full {
		p.memwb.was_full = true;
		debug.Debug("WRITE: ");
		switch inst.IType[p.memwb.inst.Opname] {
		case inst.ARITH:
			switch p.memwb.inst.Opname {
			case t.ADD, t.ADDU, t.AND, t.NOR, t.OR, t.SUB, t.SUBU, t.XOR, t.MFLO, t.MFHI:
				p.cpu.Regs[p.memwb.inst.RD] = p.memwb.alu_out;
				debug.Debug("Reg[%d] <- ALU_OUT (%d): ", p.memwb.inst.RD,p.memwb.alu_out)
			case t.ADDI, t.ANDI, t.ORI, t.XORI:
				p.cpu.Regs[p.memwb.inst.RT] = p.memwb.alu_out;
				debug.Debug("Reg[%d] <- ALU_OUT (%d): ", p.memwb.inst.RT,p.memwb.alu_out)
			}
		case inst.LOSTO:
			switch p.memwb.inst.Opname {
			case t.LB, t.LH, t.LW:
				p.cpu.Regs[p.memwb.inst.RT] = p.memwb.lmd;
				debug.Debug("Reg[%d] <- LMD (%d): ",p.memwb.inst.RT,p.memwb.lmd)
			case t.SB, t.SH, t.SW:
				debug.Debug("got a store, doing nothing: ");
			case t.LA:
				p.cpu.Regs[p.memwb.inst.RS] = p.memwb.lmd;
				debug.Debug("Reg[%d] <- LMD (%d): ",p.memwb.inst.RS,p.memwb.lmd)
			}
		case inst.BRANCH:
			debug.Debug("got a branch, doing nothing: ");
		}
		debug.Debug("%s\n", p.memwb.String());
		p.memwb.full = false;
	}

	/* memory access */
	if p.exmem.full {
		debug.Debug("MEM: ");
		switch inst.IType[p.exmem.inst.Opname] {
		case inst.ARITH:
			p.memwb.inst = p.exmem.inst;
			p.memwb.alu_out = p.exmem.alu_out;
			debug.Debug("ARITH inst: forwarding inst, alu out: ");
			p.memwb.full = true;
		case inst.LOSTO:
			p.memwb.inst = p.exmem.inst;
			switch p.exmem.inst.Opname {
			case t.LB, t.LH, t.LW:
				p.memwb.lmd = p.cpu.Mem[p.exmem.alu_out];
				debug.Debug("LOAD inst: LMD <- mem[alu_out]: ");
			case t.SB, t.SH, t.SW:
				p.cpu.Mem[p.exmem.alu_out] = p.exmem.reg_b;
				debug.Debug("STOR inst: mem[alu_out] <- reg_b: ");
			case t.LA:
				if p.ifid.inst.TGT.Section == 0 {
					fmt.Printf("Loading from text section? Sorry dave...\n");
					os.Exit(-1);
				}
				p.memwb.lmd=uint32(p.ifid.inst.TGT.Offset);
			}
			p.memwb.full = true;
		case inst.BRANCH:
			debug.Debug("BRANCH inst: ending instruction: ");
		}
		debug.Debug("%s\n", p.memwb.String());
		p.exmem.full = false;
	}

	/* exec */
	if p.idex.full {
		debug.Debug("EXEC: ");
		switch inst.IType[p.idex.inst.Opname] {
		case inst.ARITH:
			p.exmem.inst = p.idex.inst;
			p.exmem.alu_out = do_arith(p.cpu, &p.idex);
		case inst.LOSTO:
			p.exmem.inst = p.idex.inst;
			p.exmem.alu_out = p.idex.reg_a + p.idex.imm;
			p.exmem.reg_b = p.idex.reg_b;
		case inst.BRANCH:
			p.exmem.alu_out = p.idex.npc + p.idex.imm; /* <<2? */
			p.exmem.cond = p.idex.reg_a == 0;
		}
		debug.Debug("%s\n", p.exmem.String());
		p.exmem.full = true;
		p.idex.full = false;
	}

	/* decode */
	if p.ifid.full {
		debug.Debug("DECODE: ");
		p.idex.reg_a = p.cpu.Regs[p.ifid.inst.RS];
		p.idex.reg_b = p.cpu.Regs[p.ifid.inst.RT];
		p.idex.npc = p.ifid.npc;
		p.idex.inst = p.ifid.inst;
		p.idex.imm = uint32(p.ifid.inst.IMM); /* Sign extend? */
		debug.Debug("%s\n", p.idex.String());

		p.idex.full = true;
		p.ifid.full = false;
	}

	/* fetch: calcuate next PC, fetch operands */
	/* Our instruction memory is not in byte form like a real machine would be,
	   instead its just an array of pointers to instruction structs. So the
	   program counter is just incremented, not += 4. */
	if p.ready {
		debug.Debug("IFETCH: ");
		p.ifid.inst = p.cpu.get_inst_at_pc();
		debug.Debug("%s\n", p.ifid.inst.String());
		var branch bool;
		if inst.IType[p.ifid.inst.Opname] == inst.BRANCH {
			op := p.ifid.inst;
			switch op.Opname {
			case t.BEQ:    branch = p.cpu.Regs[op.RT] == p.cpu.Regs[op.RS];
			case t.BGEZ:   branch = p.cpu.Regs[op.RS] >= 0;
			//case t.BGEZAL: branch = 
			case t.BGTZ:   branch = p.cpu.Regs[op.RS] > 0;
			case t.BLEZ:   branch = p.cpu.Regs[op.RS] <= 0;
			case t.BLT:    branch = p.cpu.Regs[op.RS] < p.cpu.Regs[op.RT];
			case t.BLTZ:   branch = p.cpu.Regs[op.RS] < 0;
			//case t.BLTZAL:
			case t.BNE:    branch = p.cpu.Regs[op.RT] != p.cpu.Regs[op.RS];
			case t.BNEZ:   branch = p.cpu.Regs[op.RS] != 0;
			}
		}
		if branch {
			if p.ifid.inst.TGT.Section == 1 {
				debug.Debug("Trying to branch to data section? Uhoh...\n");
			}
			p.ifid.npc = uint32(p.ifid.inst.TGT.Offset);
		} else {
			p.ifid.npc = p.cpu.pc + 1;
		}
		p.cpu.pc = p.ifid.npc;
		//debug.Debug("%s\n", p.ifid.String());

		if !branch {
			p.ifid.full = true;
		}
		p.ready = false;
	}

	if p.cpu.cycles % 20 == 0 {
		debug.PPrint("\tIFETCH\tDECODE\tEXEC\tMEM\tWB\n");
	}

	p.cpu.cycles += 1;

	debug.PPrint("%d:", p.cpu.cycles);
	//if p.ifid.full { debug.PPrint("\t%s %s\t", p.ifid.inst.Opname.String(), p.ifid.inst.String()) }
	if p.ifid.full { debug.PPrint("\t%s", p.ifid.inst.Opname.String()) }
	else { debug.PPrint ("\t\t") }
	if p.idex.full { debug.PPrint("%s\t", p.idex.inst.Opname.String()) }
	else { debug.PPrint ("\t") }
	if p.exmem.full { debug.PPrint("%s\t", p.exmem.inst.Opname.String()) }
	else { debug.PPrint ("\t") }
	if p.memwb.full { debug.PPrint("%s\t", p.memwb.inst.Opname.String()) }
	else { debug.PPrint ("\t") }
	if p.memwb.was_full { debug.PPrint("%s => alu:%d lmd:%d\t",
									   p.memwb.inst.Opname.String(),
									   p.memwb.alu_out, p.memwb.lmd) }
	else { debug.PPrint ("\t") }
	p.memwb.was_full = false;
	debug.PPrint("\n");
}

func (c *CPU) Init() {
	c.pc = 0;
	c.cycles = 0;
	c.Pipeline = &Pipeline{c, false, ifid_buf{}, idex_buf{}, exmem_buf{}, memwb_buf{}}
}

func (c *CPU) Execute() {
	for c.pc < uint32(c.Instrs.Len()){
		c.Pipeline.Dispatch();
		c.Pipeline.Step();
		c.Pipeline.Step();
		c.Pipeline.Step();
		c.Pipeline.Step();
		c.Pipeline.Step();
		//c.printMem();
		//c.printRegs();
	}
	//c.printRegs();
	c.printMem();
}

func (c *CPU)get_inst_at_pc() *inst.Inst {
	return c.Instrs.At(int(c.pc)).(*inst.Inst);
}

func do_arith(c *CPU, idex *idex_buf) (ret uint32) {
	switch idex.inst.Opname {
	case t.ADD, t.ADDU: ret = idex.reg_a + idex.reg_b;
	case t.ADDI, t.ADDIU: ret = idex.reg_a + idex.imm;
	case t.SUB, t.SUBU: ret = idex.reg_a - idex.reg_b;

/* Logical */
	case t.AND: ret = idex.reg_a & idex.reg_b;
	case t.ANDI: ret = idex.reg_a & idex.imm;
	// case t.NOP:
	case t.NOR: ret = ^ (idex.reg_a | idex.reg_b);
	case t.OR: ret = idex.reg_a | idex.reg_b;
	case t.ORI: ret = idex.reg_a | idex.imm;
	case t.XOR: ret = idex.reg_a ^ idex.reg_b;
	case t.XORI: ret = idex.reg_a ^ idex.imm;

/* Mul and div */
//	case t.DIV:
//	case t.DIVU:
//	case t.MADD:
//	case t.MADDU:
//	case t.MSUB:
//	case t.MSUBU:
//	case t.MUL:
	case t.MULT: c.lo_reg = idex.reg_a * idex.reg_b; ret = 0;
//	case t.MULTU:

/* Accumulator access */
	case t.MFHI: ret = c.hi_reg;
	case t.MFLO: ret = c.lo_reg;
//	case t.MTHI:
//	case t.MTLO:
	}
	return;
}

