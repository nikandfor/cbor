package cbor

import (
	"bytes"
	"testing"
)

func TestEncoderInsertLen(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	mkstr := func(n int) []byte {
		b := make([]byte, n)

		for i := 0; i < len(b); {
			i += copy(b[i:], "_123456789abcdef")
		}

		return b
	}

	exp := mkstr(5)
	b = e.AppendTag(b, String, len(exp))
	st := len(b)
	b = append(b, exp...)
	b = e.InsertLen(b, String, st, len(exp), len(b)-st)

	s, end := d.Bytes(b, i)
	if end != len(b) || !bytes.Equal(exp, s) {
		tb.Errorf("wanted (%s) (%d) got (%s) (%d)\ni %d  end %d\nbuf: % x", exp, len(b[i:]), s, end-i, i, end, b[i:])
	}

	//

	i = end

	exp = mkstr(25)
	expl := 20
	b = e.AppendTag(b, String, expl)
	st = len(b)
	b = append(b, exp...)
	b = e.InsertLen(b, String, st, expl, len(b)-st)

	s, end = d.Bytes(b, i)
	if end != len(b) || !bytes.Equal(exp, s) {
		tb.Errorf("wanted (%s) (%d) got (%s) (%d)\ni %d  end %d\nbuf: % x", exp, len(b[i:]), s, end-i, i, end, b[i:])
	}

	//

	i = end

	exp = mkstr(10)
	expl = 256
	b = e.AppendTag(b, String, expl)
	st = len(b)
	b = append(b, exp...)
	b = e.InsertLen(b, String, st, expl, len(b)-st)

	s, end = d.Bytes(b, i)
	if end != len(b) || !bytes.Equal(exp, s) {
		tb.Errorf("wanted (%s) (%d) got (%s) (%d)\ni %d  end %d\nbuf: % x", exp, len(b[i:]), s, end-i, i, end, b[i:])
	}

	tb.Logf("buf: % x", b)
}
