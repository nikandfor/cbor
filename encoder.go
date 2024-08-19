package cbor

import "math"

type (
	Encoder struct{}
)

// InsertLen inserts length l before value starting at st copying the value bytes forward if needed.
// It's needed to encode the value of unknown size.
// You first expect lenth to take one byte, encodes the value
// and then inserts length moving the value if encoded length size exceeds one byte.
//
//	b = e.AppendTag(b, tag, 0)
//	st := len(b)
//	b = append(b, ...) // arbitrary value of arbitrary size
//	l := len(b) - st // for string or bytes
//	// or l = array/map length
//	b = e.InsertLen(b, tag, st, l)
func (e Encoder) InsertLen(b []byte, tag byte, st, l int) []byte {
	if l < 0 {
		panic(l)
	}

	b[st-1] &= TagMask

	if l < Len1 {
		b[st-1] |= byte(l)

		return b
	}

	sz := e.TagSize(l) - 1

	b = append(b, "        "[:sz]...)
	copy(b[st+sz:], b[st:])

	_ = e.AppendTag(b[:st-1], tag, l)

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
		return e.AppendTag64(b, Neg, uint64(-v)+1)
	}

	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendUint(b []byte, v uint) []byte {
	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendInt64(b []byte, v int64) []byte {
	if v < 0 {
		return e.AppendTag64(b, Neg, uint64(-v)+1)
	}

	return e.AppendTag64(b, Int, uint64(v))
}

func (e Encoder) AppendUint64(b []byte, v uint64) []byte {
	return e.AppendTag64(b, Int, v)
}

func (e Encoder) AppendFloat32(b []byte, v float32) []byte {
	if q := int8(v); float32(q) == v {
		return append(b, Simple|Float8, byte(q))
	}

	r := math.Float32bits(v)

	return append(b, Simple|Float32, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
}

func (e Encoder) AppendFloat(b []byte, v float64) []byte {
	if q := int8(v); float64(q) == v {
		return append(b, Simple|Float8, byte(q))
	}

	if q := float32(v); float64(q) == v {
		r := math.Float32bits(q)

		return append(b, Simple|Float32, byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
	}

	r := math.Float64bits(v)

	return append(b, Simple|Float64, byte(r>>56), byte(r>>48), byte(r>>40), byte(r>>32), byte(r>>24), byte(r>>16), byte(r>>8), byte(r))
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

func (e Encoder) AppendSimple(b []byte, x byte) []byte {
	return append(b, Simple|x)
}

func (e Encoder) AppendBool(b []byte, v bool) []byte {
	if v {
		return append(b, Simple|True)
	}

	return append(b, Simple|False)
}

func (e Encoder) AppendNil(b []byte) []byte {
	return append(b, Simple|Nil)
}

func (e Encoder) AppendUndefined(b []byte) []byte {
	return append(b, Simple|Undefined)
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
