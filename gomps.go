package main

import (
	"os";
	"flag";
	"fmt";
	"gomps/debug";
	"gomps/parser";
	"gomps/eval";
)

var verbose *bool = flag.Bool("v", false, "full verbosity");
var multissue *bool = flag.Bool("m", false, "use the pipeline");

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: gomps [-vm] input.s\n");
		os.Exit(1);
	}
	flag.Parse();
	

	debug.Verbosity = debug.PPRINT;
	if *verbose {
		debug.Verbosity = debug.DEBUG;
	}
	var c = new(eval.CPU);
	c.Init();
	p, err := parser.Parse(flag.Arg(0));
	if err != nil {
		fmt.Printf("%s\n", err.String());
	}
	c.Instrs = p.Instlist;
	c.Mem = p.Memory;
	c.Execute(*multissue);
}
