package cbor

import "testing"

func TestDecoder(tb *testing.T) {
	var d Decoder

	b := []byte{Int | 30, 0, 0, 0}

	tag, arg, i := d.Tag(b, 0)
	if i >= 0 {
		tb.Errorf("tag %x %x  i %d", tag, arg, i)
	} else {
		tb.Logf("err: %v", Error(i))
	}
}
