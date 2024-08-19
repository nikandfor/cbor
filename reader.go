package cbor

import (
	"errors"
	"fmt"
	"io"
)

type (
	Reader struct {
		io.Reader

		b    []byte
		i    int
		boff int64
	}
)

func NewReader(r io.Reader) *Reader {
	return &Reader{
		Reader: r,
	}
}

func (r *Reader) Decode() (data []byte, err error) {
	end, err := r.skipRead()
	if err != nil {
		return nil, err
	}

	st := r.i
	r.i = end

	return r.b[st:end:end], nil
}

func (r *Reader) Read(p []byte) (n int, err error) {
	end, err := r.skipRead()
	if err != nil {
		return 0, err
	}

	if len(p) < end-r.i {
		return 0, Error(r.newError(ErrShortBuffer, r.i))
	}

	copy(p, r.b[r.i:end])
	r.i = end

	return len(p), nil
}

func (r *Reader) WriteTo(w io.Writer) (n int64, err error) {
	for {
		data, err := r.Decode()
		if errors.Is(err, io.EOF) {
			return n, nil
		}
		if err != nil {
			return n, fmt.Errorf("decode: %w", err)
		}

		m, err := w.Write(data)
		n += int64(m)
		if err != nil {
			return n, fmt.Errorf("write: %w", err)
		}
	}
}

func (r *Reader) skipRead() (end int, err error) {
	for {
		end = r.skip(r.i)
		//	println("skip", r.i, end)
		if end > 0 {
			return end, nil
		}

		if end < 0 {
			return 0, Error(end)
		}

		err = r.more()
		if err != nil {
			return 0, err
		}
	}
}

func (r *Reader) skip(st int) (i int) {
	tag, sub, i := readTag(r.b, st)
	//	println("tag", st, tag, sub, i)
	if i < 0 {
		return r.newError(-i, st)
	}

	switch tag {
	case Int, Neg:
		// already read
	case Bytes, String:
		i += int(sub)
	case Array, Map:
		for el := 0; sub == -1 || el < int(sub); el++ {
			if i == len(r.b) {
				return r.newError(ErrUnexpectedEOF, i)
			}
			if sub == -1 && r.b[i] == Simple|Break {
				i++
				break
			}

			if tag == Map {
				i = r.skip(i)
				if i < 0 {
					return i
				}
			}

			i = r.skip(i)
			if i < 0 {
				return i
			}
		}
	case Labeled:
		return r.skip(i)
	case Simple:
		switch sub {
		case False,
			True,
			Null,
			Undefined,
			Break:
		case Float8:
			i += 1
		case Float16:
			i += 2
		case Float32:
			i += 4
		case Float64:
			i += 8
		default:
			return r.newError(ErrMalformed, i)
		}
	}

	if i > len(r.b) {
		return r.newError(ErrUnexpectedEOF, i)
	}

	return i
}

func (r *Reader) more() (err error) {
	copy(r.b, r.b[r.i:])
	r.b = r.b[:len(r.b)-r.i]
	r.boff += int64(r.i)
	r.i = 0

	end := len(r.b)

	if len(r.b) == 0 {
		r.b = make([]byte, 1024)
	} else {
		r.b = append(r.b, 0, 0, 0, 0, 0, 0, 0, 0)
	}

	r.b = r.b[:cap(r.b)]

	n, err := r.Reader.Read(r.b[end:])
	//	println("more", r.i, end, end+n, n, len(r.b))
	r.b = r.b[:end+n]

	if n != 0 && errors.Is(err, io.EOF) {
		err = nil
	}

	return err
}

func readTag(b []byte, st int) (tag byte, sub int64, i int) {
	if st >= len(b) {
		return tag, sub, -ErrUnexpectedEOF
	}

	i = st

	tag = b[i] & TagMask
	sub = int64(b[i] & SubMask)
	i++

	if tag == Simple {
		return
	}

	if sub < Len1 {
		return
	}

	switch sub {
	case LenBreak:
		sub = -1
	case Len1:
		if i+1 > len(b) {
			return tag, sub, -ErrUnexpectedEOF
		}

		sub = int64(b[i])
		i++
	case Len2:
		if i+2 > len(b) {
			return tag, sub, -ErrUnexpectedEOF
		}

		sub = int64(b[i])<<8 | int64(b[i+1])
		i += 2
	case Len4:
		if i+4 > len(b) {
			return tag, sub, -ErrUnexpectedEOF
		}

		sub = int64(b[i])<<24 | int64(b[i+1])<<16 | int64(b[i+2])<<8 | int64(b[i+3])
		i += 4
	case Len8:
		if i+8 > len(b) {
			return tag, sub, -ErrUnexpectedEOF
		}

		sub = int64(b[i])<<56 | int64(b[i+1])<<48 | int64(b[i+2])<<40 | int64(b[i+3])<<32 |
			int64(b[i+4])<<24 | int64(b[i+5])<<16 | int64(b[i+6])<<8 | int64(b[i+7])
		i += 8
	default:
		return tag, sub, -ErrMalformed
	}

	return tag, sub, i
}

func (r *Reader) newError(code, index int) int {
	return newError(code, int(r.boff)+index)
}
