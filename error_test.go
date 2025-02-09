package cbor

import (
	"math/bits"
	"testing"
)

func TestError(tb *testing.T) {
	for _, tc := range []struct {
		Code  int
		Index int
	}{
		{0, 0},
		{1, 0},
		{1, 10},
		{5, 1000000},
		{1<<codeWidth - 1, 10000},
	} {
		e := newError(tc.Code, tc.Index)

		code, index := Error(e).CodeIndex()
		if code != tc.Code || index != tc.Index {
			tb.Errorf("code-index %#x %#x -> error %#x -> code-index %#x %#x", tc.Code, tc.Index, e, code, index)
		}
	}

	for _, e := range []int{0, 10, 1000000} {
		code, index := Error(e).CodeIndex()

		if code != ErrOK || index != e {
			tb.Errorf("positive index %#x -> %#x %#x", e, code, index)
		}
	}

	{
		over := 1 << (bits.UintSize - codeWidth - 1)
		index := 1<<(bits.UintSize-codeWidth-2) | 1<<(bits.UintSize-8-1) | 1<<8 | 1
		code := 1<<codeWidth - 2

		e := newError(code, over|index)

		c, i := Error(e).CodeIndex()

		if e > 0 || c != code || i != index {
			tb.Errorf("index overflow: %x %x -> %x -> %x %x", over|index, code, e, i, c)

			tb.Logf("code w %d: masks %x %x", codeWidth, errorMask, codeMask)
		}
	}
}
