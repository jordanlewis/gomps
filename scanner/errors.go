package scanner

import (
	"fmt";
	"container/vector";
	"gomps/token";
)

type ErrorHandler struct {
	errors vector.Vector;
}

type Error struct {
	Pos token.Position;
	Str string;
}

func (h *ErrorHandler) Init()	{ h.errors.Init(0) }

func (h *ErrorHandler) Error(pos token.Position, str string) {
	h.errors.Push(&Error{pos, str})
}

func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Pos.String(), e.Str);
}
