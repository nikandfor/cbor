package cbor

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
	Labeled
	Simple

	TagMask = 0b1110_0000
	SubMask = 0b0001_1111
)

const (
	Len1 = 24 + iota
	Len2
	Len4
	Len8

	LenBreak = Break
)

const (
	False = 20 + iota
	True
	Null
	Undefined

	Float8
	Float16
	Float32
	Float64

	None = 0

	Break = 31
)

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
