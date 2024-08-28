package cbor

import "math"

type (
	Decoder struct {
		Flags FeatureFlags
	}
)

func (d Decoder) Skip(b []byte, st int) (i int) {
	_, _, i = d.SkipTag(b, st)
	return
}

func (d Decoder) SkipTag(b []byte, st int) (tag byte, sub int64, i int) {
	tag, sub, i = d.Tag(b, st)

	//	println(fmt.Sprintf("Skip %x  tag %x %x  i %x  data % x", st, tag, sub, i, b[st:]))

	switch tag {
	case Int, Neg:
	case String, Bytes:
		if sub >= 0 {
			_, i = d.Bytes(b, st)
			break
		}

		for !d.Break(b, &i) {
			i = d.Skip(b, i)
		}
	case Array, Map:
		for el := 0; sub == -1 && !d.Break(b, &i) || el < int(sub); el++ {
			if tag == Map {
				i = d.Skip(b, i)
			}

			i = d.Skip(b, i)
		}
	case Labeled:
		i = d.Skip(b, i)
	case Simple:
	}

	return
}

func (d Decoder) Raw(b []byte, st int) ([]byte, int) {
	i := d.Skip(b, st)

	return b[st:i], i
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
		return tag, sub, newError(ErrMalformed, st)
	}

	return
}

func (d Decoder) Signed(b []byte, st int) (v int64, i int) {
	tag, v, i := d.Tag(b, st)
	if tag == Neg {
		v++
	}

	return v, i
}

func (d Decoder) Unsigned(b []byte, st int) (v uint64, i int) {
	tag, x, i := d.Tag(b, st)
	if tag == Neg {
		x++
	}

	return uint64(x), i
}

func (d Decoder) Float32(b []byte, st int) (v float32, i int) {
	_, x, i := d.Tag(b, st)
	sub := b[st] & SubMask

	switch sub {
	case Float8:
		return float32(x), i
	case Float16:
		return d.float16(b, st+1), i
	case Float32:
		return math.Float32frombits(uint32(x)), i
	case Float64:
		return 0, newError(ErrOverflow, st)
	default:
		return 0, newError(ErrMalformed, st)
	}
}

func (d Decoder) Float(b []byte, st int) (v float64, i int) {
	_, x, i := d.Tag(b, st)
	sub := b[st] & SubMask

	switch sub {
	case Float8:
		return float64(x), i
	case Float16:
		return float64(d.float16(b, st+1)), i
	case Float32:
		return float64(math.Float32frombits(uint32(x))), i
	case Float64:
		return math.Float64frombits(uint64(x)), i
	default:
		return 0, newError(ErrMalformed, st)
	}
}

func (d Decoder) float16(b []byte, i int) float32 {
	const sig = 0b1_00000_0000000000
	const exp = 0b0_11111_0000000000
	const man = 0b0_00000_1111111111

	const exp32 = 0b0_11111111_00000000000000000000000

	r := uint32(b[i])<<8 | uint32(b[i+1])

	switch {
	case r&^sig == 0:
		return math.Float32frombits(r << 16)
	case r&exp == exp && r&man == 0:
		return math.Float32frombits(r<<16 | exp32)
	case r&exp == exp:
		return math.Float32frombits(r<<16 | exp32 | 1)
	}

	e := r&exp>>10 - 15 + 127

	r32 := r>>15<<31 | e<<23 | r&man<<13

	return float32(math.Float32frombits(r32))
}
