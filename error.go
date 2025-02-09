package cbor

import (
	"fmt"
)

type Error int

const (
	ErrOK = iota
	ErrShortBuffer
	ErrMalformed
	ErrUnexpectedEOF
	ErrOverflow

	codeWidth = 4
	codeMask  = 1<<codeWidth - 1
	errorMask = int(^uint(0) >> 1)
)

var errStrings = []string{
	"ok",
	"short buffer",
	"malformed",
	"unexpected eof",
}

func newError(code, index int) int {
	idx := index << codeWidth

	return -(idx&errorMask | code&codeMask)
}

func (e Error) Error() string {
	return fmt.Sprintf("at %d (%#[1]x): %v", e.Index(), errStrings[e.Code()])
}

func (e Error) Code() int {
	if e >= 0 {
		return 0
	}

	return int(-e & codeMask)
}

func (e Error) Index() int {
	if e >= 0 {
		return int(e)
	}

	return int(-e >> codeWidth)
}

func (e Error) CodeIndex() (code, index int) {
	return e.Code(), e.Index()
}
