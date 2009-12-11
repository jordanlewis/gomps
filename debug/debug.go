package debug

import (
	"fmt";
)

var Verbosity verbosity;

type verbosity int;

const (
	NONE = iota;
	PPRINT;
	DEBUG;
)


func Debug(format string, v ...) {
	if Verbosity >= DEBUG {
		fmt.Printf(format, v);
	}
}

func PPrint(format string, v ...) {
	if Verbosity >= PPRINT {
		fmt.Printf(format, v);
	}
}
