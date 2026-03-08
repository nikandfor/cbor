package cbor

import "fmt"

type (
	Tag byte

	Message struct {
		b    []byte
		root int
	}
)

const (
	Int Tag = iota << 5
	Neg
	Bytes
	String
	Array
	Map
	Label
	Simple

	TagMask = 0b1110_0000
	SubMask = 0b0001_1111
)

const (
	Len1 Tag = 24 + iota
	Len2
	Len4
	Len8

	LenBreak Tag = Break
)

const (
	False Tag = 20 + iota
	True
	Null
	Undefined

	Float8
	Float16
	Float32
	Float64

	None Tag = 0

	Break Tag = 31
)

func IsBool(raw Tag) bool {
	return raw == Simple|False || raw == Simple|True
}

func IsNum(raw Tag) bool {
	return IsInt(raw) || IsFloat(raw)
}

func IsInt(tag Tag) bool {
	t := tag & TagMask
	return t == Int || t == Neg
}

func IsFloat(raw Tag) bool {
	return raw >= Simple|Float8 && raw <= Simple|Float64
}

func (t Tag) String() string {
	switch t & TagMask {
	case Int, Neg:
		return "int"
	case Bytes:
		return "bytes"
	case String:
		return "string"
	case Array:
		return "array"
	case Map:
		return "map"
	case Label:
		return "labelled"
	case Simple:
	}

	switch t & SubMask {
	case False, True:
		return "bool"
	case Null:
		return "null"
	case Undefined:
		return "undefined"
	case Float8, Float16, Float32, Float64:
		return "float"
	case None:
		return "none"
	case Break:
		return "break"
	default:
		return fmt.Sprintf("%x", int(t))
	}
}
