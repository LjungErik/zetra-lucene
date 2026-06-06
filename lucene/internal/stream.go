package internal

type DataInputStream interface {
	Read([]byte) (int, error)
	Close() error
}

type DataOutputStream interface {
	Write(p []byte) (int, error)
	WriteVInt(i int) error
	WriteInt(i int) error
	WriteInt64(i int64) error
	WriteVUInt64(i uint64) error
	WriteByte(b byte) error
	Close() error

	GetWrittenBytes() uint64
	GetCheckSum() uint64
}
