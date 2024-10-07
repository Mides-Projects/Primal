package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"unsafe"
)

type Reader struct {
	r interface {
		io.Reader
		io.ByteReader
	}
}

func NewReader(payload []byte) *Reader {
	buf := bytes.NewBuffer(payload)

	return &Reader{r: buf}
}

// Uint8 reads an uint8 from the underlying buffer.
func (r *Reader) Uint8(x *uint8) {
	var err error
	*x, err = r.r.ReadByte()
	if err != nil {
		panic(err)
	}
}

// Int8 reads an int8 from the underlying buffer.
func (r *Reader) Int8(x *int8) {
	var b uint8
	r.Uint8(&b)
	*x = int8(b)
}

// Int16 reads a little endian int16 from the underlying buffer.
func (r *Reader) Int16(x *int16) {
	b := make([]byte, 2)
	if _, err := r.r.Read(b); err != nil {
		panic(err)
	}
	*x = int16(binary.BigEndian.Uint16(b))
}

// Bool reads a bool from the underlying buffer.
func (r *Reader) Bool(x *bool) {
	u, err := r.r.ReadByte()
	if err != nil {
		panic(err)
	}
	*x = *(*bool)(unsafe.Pointer(&u))
}

// errStringTooLong is an error set if a string decoded using the String method has a length that is too long.
var errStringTooLong = errors.New("string length overflows a 32-bit integer")

// StringUTF ...
func (r *Reader) StringUTF(x *string) {
	var length int16
	r.Int16(&length)
	l := int(length)
	if l > math.MaxInt16 {
		panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// String reads a string from the underlying buffer.
func (r *Reader) String(x *string) {
	var length uint32
	r.Varuint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		panic(err)
	}
	*x = *(*string)(unsafe.Pointer(&data))
}

// ByteSlice reads a byte slice from the underlying buffer, similarly to String.
func (r *Reader) ByteSlice(x *[]byte) {
	var length uint32
	r.Varuint32(&length)
	l := int(length)
	if l > math.MaxInt32 {
		panic(errStringTooLong)
	}
	data := make([]byte, l)
	if _, err := r.r.Read(data); err != nil {
		panic(err)
	}
	*x = data
}

// errVarIntOverflow is an error set if one of the Varint methods encounters a varint that does not terminate
// after 5 or 10 bytes, depending on the data type read into.
var errVarIntOverflow = errors.New("varint overflows integer")

// Varint64 reads up to 10 bytes from the underlying buffer into an int64.
func (r *Reader) Varint64(x *int64) {
	var ux uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(err)
		}

		ux |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int64(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	panic(errVarIntOverflow)
}

// Varuint64 reads up to 10 bytes from the underlying buffer into an uint64.
func (r *Reader) Varuint64(x *uint64) {
	var v uint64
	for i := 0; i < 70; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(err)
		}

		v |= uint64(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(errVarIntOverflow)
}

// Varint32 reads up to 5 bytes from the underlying buffer into an int32.
func (r *Reader) Varint32(x *int32) {
	var ux uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(err)
		}

		ux |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = int32(ux >> 1)
			if ux&1 != 0 {
				*x = ^*x
			}
			return
		}
	}
	panic(errVarIntOverflow)
}

// Varuint32 reads up to 5 bytes from the underlying buffer into an uint32.
func (r *Reader) Varuint32(x *uint32) {
	var v uint32
	for i := 0; i < 35; i += 7 {
		b, err := r.r.ReadByte()
		if err != nil {
			panic(err)
		}

		v |= uint32(b&0x7f) << i
		if b&0x80 == 0 {
			*x = v
			return
		}
	}
	panic(errVarIntOverflow)
}
