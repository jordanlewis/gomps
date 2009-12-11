package eval

import (
	t "gomps/token";
	"gomps/debug";
	"gomps/inst";
	"fmt";
)


type CPU struct {
	sp uint32;
	fp uint32;
	gp uint32;
	Regs [32]uint32;
	pc uint32;
	Instrs [1024]*inst.Inst;
	pipeline *Pipeline;
	mem [1024]uint32;
	alu_res uint32;
	mem_res uint32;
	ins_reg *inst.Inst;
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
		//p.ifid.full = true;
		//p.ifid.inst = p.cpu.get_inst_at_pc();
		//p.ifid.npc = p.cpu.pc + 1;
	}
	//debug.Debug("Dispatch: from pc = %d, next inst = %s\n", p.cpu.pc, p.ins_reg.String());
}

func (p *Pipeline) Step() {
	/* ifetch already happened due to dispatch. we start at the end of the
	   pipeline and work our way backwards, to avoid the usage of temporary
	   pipeline latches. */

	/* writeback:  */
	if p.memwb.full {
		debug.Debug("Writeback: ");
		switch inst.IType[p.memwb.inst.Opname] {
		case inst.ARITH:
			switch p.memwb.inst.Opname {
			case t.ADD, t.ADDU, t.AND, t.NOR, t.OR, t.SUB, t.SUBU, t.XOR:
				p.cpu.Regs[p.memwb.inst.RD] = p.memwb.alu_out;
				debug.Debug("Reg[%d] = ALU_OUT (%d)\n", p.memwb.inst.RT,p.memwb.alu_out)
			case t.ADDI, t.ANDI, t.ORI, t.XORI:
				p.cpu.Regs[p.memwb.inst.RT] = p.memwb.alu_out;
				debug.Debug("Reg[%d] = ALU_OUT (%d)\n", p.memwb.inst.RT,p.memwb.alu_out)
			}
		case inst.LOSTO:
			switch p.memwb.inst.Opname {
			case t.LB, t.LBU, t.LH, t.LHU, t.LUI, t.LW:
				p.cpu.Regs[p.memwb.inst.RT] = p.memwb.lmd;
				debug.Debug("Reg[%d] = LMD (%d)\n", p.memwb.inst.RT,p.memwb.lmd)
			default:
				debug.Debug("got a store, doing nothing.\n");
			}
		case inst.BRANCH:
			debug.Debug("got a branch, doing nothing\n");
		}
		p.memwb.full = false;
	}

	/* memory access */
	if p.exmem.full {
		debug.Debug("Memory access: ");
		switch inst.IType[p.exmem.inst.Opname] {
		case inst.ARITH:
			p.memwb.inst = p.exmem.inst;
			p.memwb.alu_out = p.exmem.alu_out;
			debug.Debug("ARITH inst: passing along inst and alu out\n");
			p.memwb.full = true;
		case inst.LOSTO:
			p.memwb.inst = p.exmem.inst;
			switch p.exmem.inst.Opname {
			case t.LB, t.LH, t.LW:
				p.memwb.lmd = p.cpu.mem[p.exmem.alu_out];
				debug.Debug("LOAD inst: setting LMD to mem[alu_out]\n");
			case t.SB, t.SH, t.SW:
				p.cpu.mem[p.exmem.alu_out] = p.exmem.reg_b;
				debug.Debug("STOR inst: setting mem[alu_out] to reg_b\n");
			}
			p.memwb.full = true;
		case inst.BRANCH:
			debug.Debug("BRANCH inst: ending instruction\n");
		}
		p.exmem.full = false;
	}

	/* exec */
	if p.idex.full {
		debug.Debug("Execute: ");
		switch inst.IType[p.idex.inst.Opname] {
		case inst.ARITH:
			p.exmem.inst = p.idex.inst;
			p.exmem.alu_out = do_arith(&p.idex);
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
		debug.Debug("Decode: ");
		p.idex.reg_a = p.cpu.Regs[p.ifid.inst.RS];
		p.idex.reg_b = p.cpu.Regs[p.ifid.inst.RT];
		p.idex.npc = p.ifid.npc;
		p.idex.inst = p.ifid.inst;
		p.idex.imm = p.ifid.inst.IMM; /* Sign extend? */
		debug.Debug("%s\n", p.idex.String());

		p.idex.full = true;
		p.ifid.full = false;
	}

	/* fetch: calcuate next PC, fetch operands */
	/* Our instruction memory is not in byte form like a real machine would be,
	   instead its just an array of pointers to instruction structs. So the
	   program counter is just incremented, not += 4. */
	if p.ready {
		debug.Debug("fetch: ");
		p.ifid.inst = p.cpu.get_inst_at_pc();
		if p.exmem.full &&
		   inst.IType[p.exmem.inst.Opname] == inst.BRANCH && p.exmem.cond {
			debug.Debug("foo\n");
			p.ifid.npc = p.exmem.alu_out;
		} else {
			p.ifid.npc = p.cpu.pc + 1;
		}
		p.cpu.pc = p.ifid.npc;
		debug.Debug("%s\n", p.ifid.String());

		p.ifid.full = true;
		p.ready = false;
	}
	/* Skip argument fetching thing for now */
}

func (c *CPU) Init() {
	c.pc = 0;
	c.pipeline = &Pipeline{c, false, ifid_buf{}, idex_buf{}, exmem_buf{}, memwb_buf{}}
}

func (c *CPU)get_inst_at_pc() *inst.Inst {
	return c.Instrs[c.pc];
}

func do_arith(idex *idex_buf) (ret uint32) {
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
//	case t.MULT:
//	case t.MULTU:

/* Accumulator access */
//	case t.MFHI:
//	case t.MFLO:
//	case t.MTHI:
//	case t.MTLO:
	}
	return;
}

//func (c *CPU) Step() {
//	var tmp *ast.Inst;
//	var tmp2 *ast.Inst;
//
//	debug.Debug("--------- PC = %d ---------\n", c.pc);
//	debug.Debug("IFETCH\tDECODE\tEXECUTE\tMEMORY\tWRITEBACK\n");
//
//	tmp2 = c.pipeline[1];
//	c.pipeline[1] = tmp;
//	tmp = tmp2;
//
//	c.pc += 1;
//
//	debug.Debug("DECODE: pc = %d, args = ...\n", c.pc);
//
//	/* Execute */
//	if c.pipeline[1] == nil {
//		debug.Debug("Omg\n");
//	}
//	switch c.pipeline[1].Opcode {
//	case t.ADD:
//		c.alu_res = c.Regs[c.pipeline[1].Reg2] + c.Regs[c.pipeline[1].Reg3];
//	//t.ADD;
//	//t.ADDI;
//	//t.ADDIU;
//
//	}
//	debug.Debug("EXEC(%s) = %d\n", c.pipeline[1].Opcode.String(), c.alu_res);
//
//	/* Mem access */
//	tmp2 = c.pipeline[2];
//	c.pipeline[2] = tmp;
//	tmp = tmp2;
//	switch c.pipeline[2].Opcode {
//	case t.LB,t.LBU,t.LH,t.LHU,t.LW,t.LWL,t.LWR:
//		c.mem_res = c.mem[c.pipeline[1].Reg2 + c.Regs[c.pipeline[1].Reg3]];
//		debug.Debug("MEM: load = %d\n", c.mem_res);
//	case t.SB,t.SH,t.SW,t.SWL,t.SWR,t.ULW,t.USW:
//		c.mem[c.pipeline[1].Reg2 + c.Regs[c.pipeline[1].Reg3]] =
//			c.Regs[c.pipeline[1].Reg1];
//		debug.Debug("MEM: store\n");
//	default:
//		c.mem_res = c.alu_res;
//		debug.Debug("MEM: pass (%d)\n", c.mem_res);
//	}
//
//	/* Writeback */
//	tmp2 = c.pipeline[3];
//	c.pipeline[3] = tmp;
//	tmp = tmp2;
//
//	switch c.pipeline[3].Opcode {
//	case t.SB,t.SH,t.SW,t.SWL,t.SWR,t.ULW,t.USW:
//		/* writeback does nothing with stores */
//		debug.Debug("WB: Doing nothing due to a store\n");
//	default:
//		c.Regs[c.pipeline[3].Reg1] = c.mem_res;
//		debug.Debug("WB: Reg %d <--- %d\n", c.pipeline[3].Reg1, c.mem_res)
//	}
//	return;
//
//}
