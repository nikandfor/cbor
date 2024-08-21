package cbor

import (
	"bytes"
	"math"
	"testing"
)

func TestInt(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte

	for _, tc := range []struct {
		Data    uint64
		Encoded []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{23, []byte{0x17}},
		{24, []byte{0x18, 0x18}},
		{25, []byte{0x18, 0x19}},
		{100, []byte{0x18, 0x64}},
		{1_000, []byte{0x19, 0x03, 0xe8}},
		{1_000_000, []byte{0x1a, 0x00, 0x0f, 0x42, 0x40}},
		{1_000_000_000_000, []byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00}},
		{18446744073709551615, []byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
	} {
		if tc.Data <= math.MaxInt {
			b = e.AppendInt(b[:0], int(tc.Data))

			if !bytes.Equal(b, tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> %#x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			dec, i := d.Signed(b, 0)
			if i != len(b) || int(dec) != int(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}

		if tc.Data <= math.MaxInt64 {
			b = e.AppendInt64(b[:0], int64(tc.Data))

			if !bytes.Equal(b, tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> %#x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			dec, i := d.Signed(b, 0)
			if i != len(b) || dec != int64(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}

		if tc.Data <= math.MaxUint {
			b = e.AppendUint(b[:0], uint(tc.Data))

			if !bytes.Equal(b, tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> %#x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			dec, i := d.Signed(b, 0)
			if i != len(b) || uint(dec) != uint(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}

		if true {
			b = e.AppendUint64(b[:0], tc.Data)

			if !bytes.Equal(b, tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> %#x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			dec, i := d.Unsigned(b, 0)
			if i != len(b) || uint64(dec) != uint64(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}
	}
}

func TestNeg(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	for _, tc := range []struct {
		Data    uint64
		Encoded []byte
	}{
		{1, []byte{0x20}},
		{10, []byte{0x29}},
		{100, []byte{0x38, 0x63}},
		{1_000, []byte{0x39, 0x03, 0xe7}},
		{18446744073709551615, []byte{0x3b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe}},
	} {
		var dec int64

		if tc.Data <= -math.MinInt {
			b = e.AppendInt(b, int(-tc.Data))

			if !bytes.Equal(b[i:], tc.Encoded) {
				tb.Errorf("%T(%[1]d) -> % #x, wanted %d", tc.Data, b, tc.Encoded)
			}

			dec, i = d.Signed(b, i)
			if i != len(b) || int(dec) != int(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}

		if tc.Data <= -math.MinInt64 {
			b = e.AppendInt64(b, int64(-tc.Data))

			if !bytes.Equal(b[i:], tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> % #x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			dec, i = d.Signed(b, i)
			if i != len(b) || dec != int64(tc.Data) {
				tb.Errorf("%d -> %d,  i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}

		if tc.Data <= math.MaxUint64 {
			b = e.AppendNegUint64(b, tc.Data)

			if !bytes.Equal(b[i:], tc.Encoded) {
				tb.Errorf("%T(%#[1]v) -> % #x, wanted %#x", tc.Data, b, tc.Encoded)
			}

			var dec uint64
			dec, i = d.Unsigned(b, i)
			if i != len(b) || dec != tc.Data {
				tb.Errorf("%d -> %d, i %#x / %#x", tc.Data, dec, i, len(b))
			}
		}
	}
}

func TestFloat(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	e.Flags = EncoderCompatible

	for _, tc := range []struct {
		Data    float64
		Encoded []byte
	}{
		{0, []byte{0xf9, 0x00, 0x00}},
		{math.Copysign(0, -1), []byte{0xf9, 0x80, 0x00}},
		{1, []byte{0xf9, 0x3c, 0x00}},
		{1.1, []byte{0xfb, 0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{1.5, []byte{0xf9, 0x3e, 0x00}},
		{65504, []byte{0xf9, 0x7b, 0xff}},
		{100000, []byte{0xfa, 0x47, 0xc3, 0x50, 0x00}},
		{3.4028234663852886e+38, []byte{0xfa, 0x7f, 0x7f, 0xff, 0xff}},
		{1.0e+300, []byte{0xfb, 0x7e, 0x37, 0xe4, 0x3c, 0x88, 0x00, 0x75, 0x9c}},
		//	{5.960464477539063e-8, []byte{0xf9, 0x00, 0x01}}, // wtf?
		{0.00006103515625, []byte{0xf9, 0x04, 0x00}},
		{-4, []byte{0xf9, 0xc4, 0x00}},
		{-4.1, []byte{0xfb, 0xc0, 0x10, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66}},
		{math.Inf(1), []byte{0xf9, 0x7c, 0x00}},
		{math.NaN(), []byte{0xf9, 0x7e, 0x00}},
		{math.Inf(-1), []byte{0xf9, 0xfc, 0x00}},
	} {
		b = e.AppendFloat(b, tc.Data)

		if !bytes.Equal(b[i:], tc.Encoded) {
			tb.Errorf("%T(%#24[1]v) -> % #x, wanted % #x", tc.Data, b, tc.Encoded)
		}

		var dec float64
		dec, i = d.Float(b, i)
		if i != len(b) || dec != tc.Data && !math.IsNaN(dec) && !math.IsNaN(tc.Data) {
			tb.Errorf("%T(%#24[1]v) -> %# x -> %#24v, i %#x / %#x", tc.Data, b[i:], dec, i, len(b))
		}
	}

	b = append(b[:0], 0, 1)

	tb.Logf("% x -> %f", b, d.float16(b, 0))

	check := func(e []byte) {
		if bytes.Equal(e, b[i:]) {
			i = len(b)
			return
		}

		tb.Errorf("% #x, wanted % #x", b[i:], e)

		_ = d
		i = len(b)
	}

	i = len(b)

	e.Flags = EncoderDefault

	b = e.AppendFloat(b, 0)
	check([]byte{0xf8, 0x0})

	b = e.AppendFloat(b, 127)
	check([]byte{0xf8, 127})

	b = e.AppendFloat(b, -128)
	check([]byte{0xf8, 128})
}

func TestSimple(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte

	b = e.AppendBool(b, false)
	b = e.AppendBool(b, true)
	b = e.AppendSimple(b, Null)
	b = e.AppendSimple(b, Undefined)
	b = e.AppendSimple(b, 16)

	if !bytes.Equal([]byte{
		0xf4,
		0xf5,
		0xf6,
		0xf7,
		0xf0,
	}, b) {
		tb.Errorf("% 02x", b)
	}

	var tag byte
	var arg int64

	i := 0

	for j, tc := range []struct {
		Simple int64
	}{
		{Simple: False},
		{Simple: True},
		{Simple: Null},
		{Simple: Undefined},
		{Simple: 16},
	} {
		tag, arg, i = d.Tag(b, i)
		if tag != Simple || arg != tc.Simple {
			tb.Errorf("j %d: %x %x", j, tag, arg)
		}
	}

	if i != len(b) {
		tb.Errorf("i %d: % #x", i, b)
	}
}

func TestString(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b, v []byte
	var i int

	for _, tc := range []struct {
		Tag     byte
		Data    string
		Encoded []byte
	}{
		{Bytes, "", []byte{0x40}},
		{Bytes, string([]byte{1, 2, 3, 4}), []byte{0x44, 1, 2, 3, 4}},
		{String, "", []byte{0x60}},
		{String, "a", []byte{0x61, 0x61}},
		{String, "IETF", []byte{0x64, 0x49, 0x45, 0x54, 0x46}},
		{String, `"\`, []byte{0x62, 0x22, 0x5c}},
		{String, "\u00fc", []byte{0x62, 0xc3, 0xbc}},
		{String, "\u6c34", []byte{0x63, 0xe6, 0xb0, 0xb4}},
		//	{string([]byte{0xd8, 0x00, 0xdd, 0x51}), []byte{0x64, 0xf0, 0x90, 0x85, 0x91}}, // TODO
	} {
		b = e.AppendTagString(b, tc.Tag, tc.Data)

		if !bytes.Equal(b[i:], tc.Encoded) {
			tb.Errorf("%#24v -> % #x, wanted % #x", tc.Data, b[i:], tc.Encoded)
		}

		tag := d.TagOnly(b, i)
		v, i = d.Bytes(b, i)
		if tag != tc.Tag || string(v) != tc.Data {
			tb.Errorf("%x : %#24v -> % #x -> %x : %#24v", tc.Tag, tc.Data, b[i:], tag, string(v))
		}
	}
}

func TestArray(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int
	var arr []byte

	check := func() {
		var tag byte
		var l int64

		tag, l, i = d.Tag(b, i)
		if tag != Array || int(l) != len(arr) {
			tb.Errorf("%x %x", tag, l)
		}
		if !bytes.Equal(arr, b[i:]) {
			tb.Errorf("% #x", b[i:])
		}

		i = len(b)
	}

	b = e.AppendArray(b, 0)

	check()

	arr = []byte{1, 2, 3}
	b = e.AppendArray(b, len(arr))
	b = append(b, arr...)

	check()

	b = e.AppendArray(b, 3)
	b = append(b, 1)

	b = e.AppendArray(b, 2)
	b = append(b, 2, 3)

	b = e.AppendArray(b, 2)
	b = append(b, 4, 5)

	if !bytes.Equal([]byte{0x83, 0x01, 0x82, 0x02, 0x03, 0x82, 0x04, 0x05}, b[i:]) {
		tb.Errorf("% #x", b[i:])
	}

	_ = d
	i = len(b)

	arr = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}

	b = e.AppendArray(b, len(arr))
	b = append(b, arr...)

	check()
}

func TestMap(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	check := func(e []byte) {
		if bytes.Equal(e, b[i:]) {
			i = len(b)
			return
		}

		tb.Errorf("% #x, wanted % #x", b[i:], e)

		_ = d
		i = len(b)
	}

	b = e.AppendMap(b, 0)
	check([]byte{0xa0})

	b = e.AppendMap(b, 2)
	b = append(b, 1, 2, 3, 4)

	check([]byte{0xa2, 1, 2, 3, 4})

	b = e.AppendMap(b, 2)

	b = e.AppendString(b, "a")
	b = append(b, 1)

	b = e.AppendString(b, "b")

	b = e.AppendArray(b, 2)
	b = append(b, 2, 3)

	check([]byte{0xa2, 0x61, 'a', 1, 0x61, 'b', 0x82, 2, 3})

	b = e.AppendArray(b, 2)
	b = e.AppendString(b, "a")

	b = e.AppendMap(b, 1)
	b = e.AppendString(b, "b")
	b = e.AppendString(b, "c")

	check([]byte{0x82, 0x61, 'a', 0xa1, 0x61, 'b', 0x61, 'c'})

	s := "abcde"

	b = e.AppendMap(b, len(s))

	for i := range s {
		b = e.AppendString(b, s[i:i+1])
		b = e.AppendString(b, string(s[i]-'a'+'A'))
	}

	check([]byte{
		0xa5,
		0x61, 'a', 0x61, 'A',
		0x61, 'b', 0x61, 'B',
		0x61, 'c', 0x61, 'C',
		0x61, 'd', 0x61, 'D',
		0x61, 'e', 0x61, 'E',
	})
}

func TestBreak(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	check := func(e []byte) {
		if bytes.Equal(e, b[i:]) {
			i = len(b)
			return
		}

		tb.Errorf("% #x, wanted % #x", b[i:], e)

		_ = d
		i = len(b)
	}

	b = e.AppendTag(b, Bytes, -1)
	b = e.AppendBytes(b, []byte{1, 2})
	b = e.AppendBytes(b, []byte{3, 4, 5})
	b = e.AppendBreak(b)

	check([]byte{0x5f, 0x42, 1, 2, 0x43, 3, 4, 5, 0xff})

	b = e.AppendTag(b, String, -1)
	b = e.AppendString(b, "strea")
	b = e.AppendString(b, "ming")
	b = e.AppendBreak(b)

	check([]byte{0x7f, 0x65, 's', 't', 'r', 'e', 'a', 0x64, 'm', 'i', 'n', 'g', 0xff})

	b = e.AppendArray(b, -1)
	b = e.AppendBreak(b)

	check([]byte{0x9f, 0xff})

	b = e.AppendArray(b, -1)
	b = append(b, 1)
	b = e.AppendArray(b, 2)
	b = append(b, 2, 3)
	b = e.AppendArray(b, -1)
	b = append(b, 4, 5)
	b = e.AppendBreak(b)
	b = e.AppendBreak(b)

	check([]byte{0x9f, 0x01, 0x82, 0x02, 0x03, 0x9f, 0x04, 0x05, 0xff, 0xff})

	b = e.AppendArray(b, -1)
	b = append(b, 1)
	b = e.AppendArray(b, 2)
	b = append(b, 2, 3)
	b = e.AppendArray(b, 2)
	b = append(b, 4, 5)
	b = e.AppendBreak(b)

	check([]byte{0x9f, 0x01, 0x82, 0x02, 0x03, 0x82, 0x04, 0x05, 0xff})

	b = e.AppendArray(b, 3)
	b = append(b, 1)
	b = e.AppendArray(b, 2)
	b = append(b, 2, 3)
	b = e.AppendArray(b, -1)
	b = append(b, 4, 5)
	b = e.AppendBreak(b)

	check([]byte{0x83, 0x01, 0x82, 0x02, 0x03, 0x9f, 0x04, 0x05, 0xff})

	b = e.AppendArray(b, 3)
	b = append(b, 1)
	b = e.AppendArray(b, -1)
	b = append(b, 2, 3)
	b = e.AppendBreak(b)
	b = e.AppendArray(b, 2)
	b = append(b, 4, 5)

	check([]byte{0x83, 0x01, 0x9f, 0x02, 0x03, 0xff, 0x82, 0x04, 0x05})

	exp := []byte{0x9f, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 0xff}

	b = e.AppendArray(b, -1)
	b = append(b, exp[1:len(exp)-1]...)
	b = e.AppendBreak(b)

	check(exp)

	b = e.AppendMap(b, -1)
	b = e.AppendString(b, "a")
	b = append(b, 1)
	b = e.AppendString(b, "b")
	b = e.AppendArray(b, -1)
	b = append(b, 2, 3)
	b = e.AppendBreak(b)
	b = e.AppendBreak(b)

	check([]byte{0xbf, 0x61, 'a', 1, 0x61, 'b', 0x9f, 0x02, 0x03, 0xff, 0xff})

	b = e.AppendArray(b, 2)
	b = e.AppendString(b, "a")
	b = e.AppendMap(b, -1)
	b = e.AppendString(b, "b")
	b = e.AppendString(b, "c")
	b = e.AppendBreak(b)

	check([]byte{0x82, 0x61, 'a', 0xbf, 0x61, 'b', 0x61, 'c', 0xff})

	b = e.AppendMap(b, -1)
	b = e.AppendString(b, "Fun")
	b = e.AppendBool(b, true)
	b = e.AppendString(b, "Amt")
	b = e.AppendInt(b, -1)
	b = e.AppendBreak(b)

	check([]byte{0xbf, 0x63, 'F', 'u', 'n', 0xf5, 0x63, 'A', 'm', 't', 0x20, 0xff})

	// SkipTag

	// tb.Logf("b: % #x", b)

	i = 0

	for i != len(b) {
		i = d.Skip(b, i)

		if i <= 0 || i > len(b) {
			tb.Errorf("i: %v", i)
			break
		}
	}
}

func TestLabeled(tb *testing.T) {
	var e Encoder
	var d Decoder
	var b []byte
	var i int

	check := func(e []byte) {
		if bytes.Equal(e, b[i:]) {
			i = len(b)
			return
		}

		tb.Errorf("% #x, wanted % #x", b[i:], e)

		_ = d
		i = len(b)
	}

	b = e.AppendLabeled(b, 0)
	b = e.AppendString(b, "abcd")

	check([]byte{0xc0, 0x64, 'a', 'b', 'c', 'd'})

	b = e.AppendLabeled(b, 23)
	b = append(b, 1)

	check([]byte{0xd7, 1})

	b = e.AppendLabeled(b, 24)
	b = append(b, 1)

	check([]byte{0xd8, 24, 1})

	b = e.AppendLabeled(b, 32)
	b = append(b, 1)

	check([]byte{0xd8, 32, 1})
}
