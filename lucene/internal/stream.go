package internal

type DataInputStream interface {
	Read([]byte) (int, error)
	Close() error
}

type DataOutputStream interface {
	Write(p []byte) (int, error)
	WriteVInt(i int) error
	WriteByte(b byte) error
	Close() error

	GetWrittenBytes() int64
	GetCheckSum() uint64
}
