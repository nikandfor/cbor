package cbor

import (
	"encoding/hex"
	"fmt"
)

func Dump(r []byte) (s string) {
	var w []byte

	defer func() {
		p := recover()
		if p == nil {
			return
		}
		defer panic(p)

		w = fmt.Appendf(w, "panic: %v\n", p)
		w = fmt.Appendf(w, "%s\n", hex.Dump(r))

		s = string(w)
	}()

	w, _ = dump(w, r, 0, 0)
	return string(w)
}

func dump(w, r []byte, st, depth int) (_ []byte, i1 int) {
	const spaces = "                                          "
	var d Decoder

	tag, sub, i := d.Tag(r, st)

	w = fmt.Appendf(w, "%4x%s  ", st, spaces[:2*depth])

	switch tag {
	case Int, Neg:
		v, i := d.Unsigned(r, st)

		w = fmt.Appendf(w, "% x  %s%d\n", r[st:i], csel(tag == Neg, "-", ""), v)
	case Bytes, String:
		if sub >= 0 {
			var v []byte
			v, i = d.Bytes(r, st)
			w = fmt.Appendf(w, "% x  %q\n", r[st:i], v)
			break
		}

		w = fmt.Appendf(w, "% x\n", r[st:i])

		l := int(sub)

		for j := 0; l < 0 || j < l; j++ {
			if l < 0 && d.Break(r, &i) {
				w, i = dump(w, r, i-1, depth+1)
				break
			}

			w, i = dump(w, r, i, depth+1)
		}
	case Array, Map:
		w = fmt.Appendf(w, "% x  %x\n", r[st:i], sub)

		l := int(sub)

		for j := 0; l < 0 || j < l; j++ {
			if l < 0 && d.Break(r, &i) {
				w, i = dump(w, r, i-1, depth+1)
				break
			}

			if tag == Map {
				w, i = dump(w, r, i, depth+1)
			}

			w, i = dump(w, r, i, depth+1)
		}
	case Labeled:
		w = fmt.Appendf(w, "% x\n", r[st:i])
		w, i = dump(w, r, i, depth+1)
	case Simple:
		switch {
		case sub < 0:
			w = fmt.Appendf(w, "% x  break\n", r[st:i])
		case sub < Float8:
			v := []string{
				None:      "none",
				False:     "false",
				True:      "true",
				Null:      "null",
				Undefined: "undefined",
				Float8:    "",
			}[sub]

			w = fmt.Appendf(w, "% x  %v\n", r[st:i], v)
		case sub <= Float64:
			v, _ := d.Float(r, st)
			w = fmt.Appendf(w, "% x  %v\n", r[st:i], v)
		default:
			w = fmt.Appendf(w, "% x\n", r[st:i])
		}
	}

	return w, i
}

func csel[T any](cond bool, t, f T) T {
	if cond {
		return t
	}

	return f
}
