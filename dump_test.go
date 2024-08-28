package cbor

import "testing"

func TestDump(tb *testing.T) {
	var e Encoder
	var b []byte

	b = e.AppendArray(b, -1)

	b = e.AppendArray(b, 2)
	b = e.AppendInt(b, 12)
	b = e.AppendInt(b, -12)

	b = e.AppendMap(b, 1)
	b = e.AppendTagString(b, Bytes, "bytes")
	b = e.AppendTagString(b, String, "bytes")

	b = e.AppendTag(b, String, -1)
	b = e.AppendString(b, "first")
	b = e.AppendTag(b, String, -1)
	b = e.AppendString(b, ", second")
	b = e.AppendString(b, ", third")
	b = e.AppendBreak(b)
	b = e.AppendBreak(b)

	b = e.AppendTag(b, Bytes, -1)
	b = e.AppendTagString(b, Bytes, "first")
	b = e.AppendTag(b, String, -1)
	b = e.AppendString(b, ", second")
	b = e.AppendString(b, ", third")
	b = e.AppendBreak(b)
	b = e.AppendBreak(b)

	b = e.AppendBreak(b)

	tb.Logf("dump\n%s", Dump(b))
}
