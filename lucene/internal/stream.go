package internal

type DataInputStream interface {
	Read([]byte) (int, error)
	Close() error
}

type DataOutputStream interface {
	Write(p []byte) (int, error)
	WriteVInt(i int) error
	WriteVUInt64(i uint64) error
	WriteByte(b byte) error
	Close() error

	GetWrittenBytes() int
	GetCheckSum() uint64
}
