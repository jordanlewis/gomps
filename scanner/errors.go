package scanner

import (
	"fmt";
	"container/vector";
	"gomps/token";
)

type ErrorHandler interface {
	Error(pos token.Position, str string);
}

type ErrorList struct {
	errors vector.Vector;
}

type Error struct {
	Pos token.Position;
	Str string;
}

func (h *ErrorList) Init()	{ h.errors.Init(0) }

func (h *ErrorList) Error(pos token.Position, str string) {
	h.errors.Push(&Error{pos, str});
	fmt.Printf("%s\n", h.errors.Last());

}

func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Pos.String(), e.Str);
}
