package cbor

const (
	Int = iota << 5
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
	Float8 = 24 + iota
	Float16
	Float32
	Float64

	False     = 20
	True      = 21
	Null      = 22
	Undefined = 23

	None = 0

	Break = 31
)

type (
	Tag = byte

	Message struct {
		b    []byte
		root int
	}
)
