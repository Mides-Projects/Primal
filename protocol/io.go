package protocol

type IO interface {
	Uint8(x *uint8)
	Int8(x *int8)
	Int16(x *int16)
	Bool(x *bool)
	StringUTF(x *string)
	String(x *string)
	ByteSlice(x *[]byte)
	Varint64(x *int64)
	Varuint64(x *uint64)
	Varint32(x *int32)
	Varuint32(x *uint32)
}
