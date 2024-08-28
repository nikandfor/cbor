package cbor

import (
	"testing"
)

func TestDecoder(tb *testing.T) {
	var d Decoder

	b := []byte{Int | 30, 0, 0, 0}

	tag, arg, i := d.Tag(b, 0)
	if i >= 0 {
		tb.Errorf("tag %x %x  i %d", tag, arg, i)
	} else {
		tb.Logf("got expected err: %v", Error(i))
	}
}

func TestDecoderSkipNeg(tb *testing.T) {
	b := []byte{
		0x72, 0x74, 0x79, 0x05, 0x24, 0xfa, 0x3f, 0x80, 0x00, 0x00, 0xfa, 0xbf, 0x80, 0x00, 0x00, 0x8d,
	}

	var d Decoder

	st := 0x5
	tag, sub, i := d.SkipTag(b, st)
	if tag != Simple || sub != Float32 || i != st+5 {
		tb.Errorf("%x -> %x %x %x", st, tag, sub, i)
	}
}
