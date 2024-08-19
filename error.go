package cbor

import "fmt"

type Error int

const (
	ErrOK = iota
	ErrShortBuffer
	ErrMalformed
	ErrUnexpectedEOF
	ErrOverflow

	errorMask       = 0xff
	errorIndexShift = 8
)

var errStrings = []string{
	"",
	"short buffer",
	"malformed",
	"unexpected eof",
}

func newError(code, index int) int {
	return -(index<<errorIndexShift | code)
}

func (e Error) Error() string {
	return fmt.Sprintf("at %d (%#[1]x): %v", e.Index(), errStrings[e.Code()])
}

func (e Error) Code() int {
	if e >= 0 {
		return 0
	}

	return int(-e & errorMask)
}

func (e Error) Index() int {
	if e >= 0 {
		return int(e)
	}

	return int(-e >> errorIndexShift)
}

func (e Error) CodeIndex() (code, index int) {
	return e.Code(), e.Index()
}
