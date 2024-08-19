package cbor

import "math"

type (
	Decoder struct{}
)

func (d Decoder) Skip(b []byte, st int) (i int) {
	_, _, i = d.SkipTag(b, st)
	return
}

func (d Decoder) SkipTag(b []byte, st int) (tag byte, sub int64, i int) {
	tag, sub, i = d.Tag(b, st)

	//	println(fmt.Sprintf("Skip %x  tag %x %x %x  data % x", st, tag, sub, i, b[st:]))

	switch tag {
	case Int, Neg:
		_, i = d.Unsigned(b, st)
	case String, Bytes:
		_, i = d.Bytes(b, st)
	case Array, Map:
		for el := 0; sub == -1 || el < int(sub); el++ {
			if sub == -1 && d.Break(b, &i) {
				break
			}

			if tag == Map {
				i = d.Skip(b, i)
			}

			i = d.Skip(b, i)
		}
	case Labeled:
		i = d.Skip(b, i)
	case Simple:
		switch sub {
		case False,
			True,
			Nil,
			Undefined,
			Break:
		case Float8, Float16, Float32, Float64:
			i += 1 << (int(sub) - Float8)
		default:
			return tag, sub, newError(ErrMalformed, i)
		}
	}

	return
}

func (d Decoder) Break(b []byte, i *int) bool {
	if b[*i] != Simple|Break {
		return false
	}

	*i++

	return true
}

func (d Decoder) Bytes(b []byte, st int) (v []byte, i int) {
	_, l, i := d.Tag(b, st)

	return b[i : i+int(l)], i + int(l)
}

func (d Decoder) TagOnly(b []byte, st int) (tag byte) {
	return b[st] & TagMask
}

func (d Decoder) Tag(b []byte, st int) (tag byte, sub int64, i int) {
	i = st

	tag = b[i] & TagMask
	sub = int64(b[i] & SubMask)
	i++

	if tag == Simple {
		return
	}

	switch {
	case sub < Len1:
		// we are ok
	case sub == LenBreak:
		sub = -1
	case sub == Len1:
		sub = int64(b[i])
		i++
	case sub == Len2:
		sub = int64(b[i])<<8 | int64(b[i+1])
		i += 2
	case sub == Len4:
		sub = int64(b[i])<<24 | int64(b[i+1])<<16 | int64(b[i+2])<<8 | int64(b[i+3])
		i += 4
	case sub == Len8:
		sub = int64(b[i])<<56 | int64(b[i+1])<<48 | int64(b[i+2])<<40 | int64(b[i+3])<<32 |
			int64(b[i+4])<<24 | int64(b[i+5])<<16 | int64(b[i+6])<<8 | int64(b[i+7])
		i += 8
	default:
		return tag, sub, newError(ErrMalformed, i)
	}

	return
}

func (d Decoder) Signed(b []byte, st int) (v int64, i int) {
	u, i := d.Unsigned(b, st)

	if b[st]&TagMask == Int {
		return int64(u), i
	}

	return 1 - int64(u), i
}

func (d Decoder) Unsigned(b []byte, st int) (v uint64, i int) {
	i = st

	v = uint64(b[i]) & SubMask
	i++

	switch {
	case v < Len1:
		// we are ok
	case v == Len1:
		v = uint64(b[i])
		i++
	case v == Len2:
		v = uint64(b[i])<<8 | uint64(b[i+1])
		i += 2
	case v == Len4:
		v = uint64(b[i])<<24 | uint64(b[i+1])<<16 | uint64(b[i+2])<<8 | uint64(b[i+3])
		i += 4
	case v == Len8:
		v = uint64(b[i])<<56 | uint64(b[i+1])<<48 | uint64(b[i+2])<<40 | uint64(b[i+3])<<32 |
			uint64(b[i+4])<<24 | uint64(b[i+5])<<16 | uint64(b[i+6])<<8 | uint64(b[i+7])
		i += 8
	default:
		return 0, newError(ErrMalformed, i)
	}

	return
}

func (d Decoder) Float32(b []byte, st int) (v float64, i int) {
	i = st

	st = int(b[i]) & SubMask
	i++

	switch {
	case st == Float8:
		v = float64(int8(b[i]))
		i++
	case st == Float32:
		v = float64(math.Float32frombits(
			uint32(b[i])<<24 | uint32(b[i+1])<<16 | uint32(b[i+2])<<8 | uint32(b[i+3]),
		))

		i += 4
	case st == Float64:
		return 0, newError(ErrOverflow, i)
	default:
		return 0, newError(ErrMalformed, i)
	}

	return
}

func (d Decoder) Float(b []byte, st int) (v float64, i int) {
	i = st

	st = int(b[i]) & SubMask
	i++

	switch {
	case st == Float8:
		v = float64(int8(b[i]))
		i++
	case st == Float32:
		v = float64(math.Float32frombits(
			uint32(b[i])<<24 | uint32(b[i+1])<<16 | uint32(b[i+2])<<8 | uint32(b[i+3]),
		))

		i += 4
	case st == Float64:
		v = math.Float64frombits(
			uint64(b[i])<<56 | uint64(b[i+1])<<48 | uint64(b[i+2])<<40 | uint64(b[i+3])<<32 |
				uint64(b[i+4])<<24 | uint64(b[i+5])<<16 | uint64(b[i+6])<<8 | uint64(b[i+7]),
		)

		i += 8
	default:
		return 0, newError(ErrMalformed, i)
	}

	return
}
