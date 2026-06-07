package stream

import (
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type OutputStream struct {
	writtenBytes uint64
	writer       io.WriteCloser
}

var _ io.WriteCloser = (*OutputStream)(nil)
var _ internal.DataOutputStream = (*OutputStream)(nil)

func NewOutputStream(writer io.WriteCloser) *OutputStream {
	return &OutputStream{
		writtenBytes: 0,
		writer:       writer,
	}
}

func (s *OutputStream) Write(p []byte) (int, error) {
	n, err := s.writer.Write(p)
	s.writtenBytes += uint64(n)
	return n, err
}

func (s *OutputStream) WriteVInt(i int) error {
	return writeVInt(s, i)
}

func (s *OutputStream) WriteInt(i int) error {
	return writeInt(s, i)
}

func (s *OutputStream) WriteInt64(i int64) error {
	return writeInt64(s, i)
}

func (s *OutputStream) WriteUInt64(i uint64) error {
	return writeUInt64(s, i)
}

func (s *OutputStream) WriteVUInt64(i uint64) error {
	return writeVUInt64(s, i)
}

func (s *OutputStream) WriteByte(b byte) error {
	_, err := s.writer.Write([]byte{b})
	if err != nil {
		return err
	}

	return nil
}

func (s *OutputStream) Close() error {
	return s.writer.Close()
}

func (s *OutputStream) GetFilePointer() uint64 {
	return s.writtenBytes
}

func (s *OutputStream) GetCheckSum() uint64 {
	return 0
}
