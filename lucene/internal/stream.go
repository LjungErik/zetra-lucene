package internal

type DataInputStream interface {
	Read([]byte) (int, error)
	Close() error
}

type DataOutputStream interface {
	Write(p []byte) (int, error)
	Close() error

	WriteByte(p byte) error

	WriteVInt(i int) error
	WriteInt(i int) error
	WriteInt64(i int64) error
	WriteUInt64(i uint64) error
	WriteVUInt64(i uint64) error

	GetWrittenBytes() uint64
	GetCheckSum() uint64
}

type ErrStream interface {
	Write(p []byte)
	WriteVInt(i int)
	WriteInt(i int)
	WriteInt64(i int64)
	WriteUInt64(i uint64)
	WriteVUInt64(i uint64)
	WriteByte(b byte)

	WriteFunc(func(DataOutputStream) error)

	Error() error
}
