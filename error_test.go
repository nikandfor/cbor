package cbor

import "testing"

func TestError(tb *testing.T) {
	for _, tc := range []struct {
		Code  int
		Index int
	}{
		{0, 0},
		{1, 0},
		{1, 10},
		{5, 1000000},
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
}
