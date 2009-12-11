package eval

import (
	"testing";
	"gomps/token";
	"gomps/inst";
)

func TestEval(t *testing.T) {
	var c = new(CPU);
	c.Init();
	c.Regs[1] = 2;
	c.Regs[2] = 3;
	c.Regs[4] = 5;
	c.Regs[5] = 6;
	c.Instrs[0] = &inst.Inst{token.ADD, 1,2,3,0,0,0};
	c.Instrs[1] = &inst.Inst{token.ADD, 4,5,6,0,0,0};
	c.pipeline.Dispatch();
	c.pipeline.Step();
	c.pipeline.Dispatch();
	c.pipeline.Step();
	c.pipeline.Step();
	c.pipeline.Step();
	c.pipeline.Step();
}
