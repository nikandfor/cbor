package cbor

import "math"

type (
	Encoder struct {
		Flags EncoderFlags
	}

	EncoderFlags int
)

const (
	_ = 1 << iota
	EncoderFloat8Int
	EncoderFloat16

	EncoderDefault    = EncoderFloat8Int
	EncoderCompatible = EncoderFloat16
)

// InsertLen inserts length l before value starting at st copying the value bytes forward if needed.
// It's needed to encode the value of unknown size.
// You first expect lenth to take one byte, encodes the value
// and then inserts length moving the value if encoded length size exceeds one byte.
//
//	expectedLen := 0
//	b = e.AppendTag(b, tag, expectedLen)
//	st := len(b)
//	b = append(b, ...) // arbitrary value of arbitrary size
//	l := len(b) - st // for string or bytes
//	// or l = array/map length
//	b = e.InsertLen(b, tag, st, expectedLen, l)
func (e Encoder) InsertLen(b []byte, tag byte, st, l0, l int) []byte {
	if l < 0 {
		panic(l)
	}

	sz0 := e.TagSize(l0)
	sz := e.TagSize(l)
	newst := st - sz0 + sz

	if sz > sz0 {
		b = append(b, "        "[:sz-sz0]...)
	}

	if sz != sz0 {
		copy(b[newst:], b[st:])
		b = b[:newst+l]
	}

	_ = e.AppendTag(b[:newst-sz], tag, l)

	return b
}

func (e Encoder) AppendMap(b []byte, l int) []byte {
	return e.AppendTag(b, Map, l)
}

func (e Encoder) AppendArray(b []byte, l int) []byte {
	return e.AppendTag(b, Array, l)
}

func (e Encoder) AppendString(b []byte, s string) []byte {
	b = e.AppendTag(b, String, len(s))
	return append(b, s...)
}

func (e Encoder) AppendBytes(b, s []byte) []byte {
	b = e.AppendTag(b, Bytes, len(s))
	return append(b, s...)
}

func (e Encoder) AppendTagString(b []byte, tag byte, s string) []byte {
	b = e.AppendTag(b, tag, len(s))
	return append(b, s...)
}

func (e Encoder) AppendTagBytes(b []byte, tag byte, s []byte) []byte {
	b = e.AppendTag(b, tag, len(s))
	return append(b, s...)
}

func (e Encoder) AppendInt(b []byte, v int) []byte {
	if v < 0 {
		return e.AppendTag64(b, Neg, uint64(-v)-1)
	}

	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendUint(b []byte, v uint) []byte {
	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendInt64(b []byte, v int64) []byte {
	if v < 0 {
		return e.AppendTag64(b, Neg, uint64(-v)-1)
	}

	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendUint64(b []byte, v uint64) []byte {
	return e.AppendTag64(b, Int, v)
}

func (e Encoder) AppendNegUint64(b []byte, v uint64) []byte {
	return e.AppendTag64(b, Neg, v-1)
}

func (e Encoder) AppendFloat32(b []byte, v float32) []byte {
	if e.Flags.Is(EncoderFloat8Int) {
		if q := int8(v); float32(q) == v {
			return append(b, Simple|Float8, byte(q))
		}
	}

	return e.appendFloat32(b, v)
}

func (e Encoder) AppendFloat(b []byte, v float64) []byte {
	if e.Flags.Is(EncoderFloat8Int) {
		if q := int8(v); float64(q) == v {
			return append(b, Simple|Float8, byte(q))
		}
	}

	q := float32(v)

	if float64(q) == v || math.IsNaN(v) {
		return e.appendFloat32(b, q)
	}

	r := math.Float64bits(v)

	return append(b, Simple|Float64, byte(r>>56), byte(r>>48), byte(r>>40), byte(r>>32), byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func (e Encoder) appendFloat32(b []byte, v float32) []byte {
	r := math.Float32bits(v)

	if e.Flags.Is(EncoderFloat16) {
		if b, ok := e.appendFloat16(b, r); ok {
			return b
		}
	}

	return append(b, Simple|Float32, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func (e Encoder) appendFloat16(b []byte, r uint32) ([]byte, bool) {
	const (
		// 1 + 8 + 23
		sig  = 0b1_00000000_00000000000000000000000
		exp  = 0b0_11111111_00000000000000000000000
		manm = 0b0_00000000_11111111111111111111111
		manx = 0b0_00000000_11111111110000000000000

		exp32 = 0b0_11111_0000000000
	)

	var r16 uint32

	switch {
	case r&^sig == 0: // zero
		r16 = r >> 16
	case r&exp == exp && r&manm == 0: // inf
		r16 = r >> 16 & 0b1_11111_0000000000
	case r&exp == exp: // nan
		r16 = r >> 16 & 0b1_11111_0000000000
		r16 |= r&1 | r>>22&1<<9
	case r&manm&^manx == 0:
		e := r&exp>>23 - 127 + 15
		if e >= 32 {
			break
		}

		r16 = r&sig>>16 | e<<10 | r&manx>>13
	}

	if r16 == 0 && r&^sig != 0 {
		return b, false
	}

	return append(b, Simple|Float16, byte(r16>>8), byte(r16)), true
}

func (e Encoder) AppendTag(b []byte, tag byte, v int) []byte {
	switch {
	case v == -1:
		return append(b, tag|LenBreak)
	case v < Len1:
		return append(b, tag|byte(v))
	case v <= 0xff:
		return append(b, tag|Len1, byte(v))
	case v <= 0xffff:
		return append(b, tag|Len2, byte(v>>8), byte(v))
	case v <= 0xffff_ffff:
		return append(b, tag|Len4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	default:
		return append(b, tag|Len8, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

func (e Encoder) AppendTag64(b []byte, tag byte, v uint64) []byte {
	switch {
	case v < Len1:
		return append(b, tag|byte(v))
	case v <= 0xff:
		return append(b, tag|Len1, byte(v))
	case v <= 0xffff:
		return append(b, tag|Len2, byte(v>>8), byte(v))
	case v <= 0xffff_ffff:
		return append(b, tag|Len4, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	default:
		return append(b, tag|Len8, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32), byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
	}
}

func (e Encoder) AppendTagBreak(b []byte, tag byte) []byte {
	return append(b, tag|LenBreak)
}

func (e Encoder) AppendLabeled(b []byte, x int) []byte {
	return e.AppendTag(b, Labeled, x)
}

func (e Encoder) AppendSimple(b []byte, x int) []byte {
	return append(b, Simple|byte(x))
}

func (e Encoder) AppendBool(b []byte, v bool) []byte {
	if v {
		return append(b, Simple|True)
	}

	return append(b, Simple|False)
}

func (e Encoder) AppendNull(b []byte) []byte {
	return append(b, Simple|Null)
}

func (e Encoder) AppendUndefined(b []byte) []byte {
	return append(b, Simple|Undefined)
}

func (e Encoder) AppendNone(b []byte) []byte {
	return append(b, Simple|None)
}

func (e Encoder) AppendBreak(b []byte) []byte {
	return append(b, Simple|Break)
}

func (e Encoder) TagSize(v int) int {
	switch {
	case v == -1:
		return 1
	case v < Len1:
		return 1
	case v <= 0xff:
		return 1 + 1
	case v <= 0xffff:
		return 1 + 2
	case v <= 0xffff_ffff:
		return 1 + 4
	default:
		return 1 + 8
	}
}

func (e Encoder) Tag64Size(v int64) int {
	switch {
	case v < Len1:
		return 1
	case v <= 0xff:
		return 1 + 1
	case v <= 0xffff:
		return 1 + 2
	case v <= 0xffff_ffff:
		return 1 + 4
	default:
		return 1 + 8
	}
}

func (f EncoderFlags) Is(ff EncoderFlags) bool {
	return f&ff == ff
}
