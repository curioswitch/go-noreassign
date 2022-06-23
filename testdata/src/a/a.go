package a

import (
	"b"

	"errors"
)
import cc "c"

import (
	"d"
	"d/e"
)

var st = struct {
	ErrSt error
}{}

func foo() {
	b.ErrB = nil // want "reassigning sentinel error"

	cc.ErrC = nil // want "reassigning sentinel error"

	d.ErrD = nil // want "reassigning sentinel error"

	e.ErrE = nil // want "reassigning sentinel error"

	st.ErrSt = errors.New("foo")
}
