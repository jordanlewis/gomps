package debug

import (
	"fmt";
)

func Debug(format string, v ...) {
	fmt.Printf(format, v);
}
