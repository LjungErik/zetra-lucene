package stream

import (
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type OutputStream struct {
	writtenBytes int64
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
	s.writtenBytes += int64(n)
	return n, err
}

func (s *OutputStream) WriteVInt(i int) error {
	return writeVInt(s, i)
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

func (s *OutputStream) GetWrittenBytes() int64 {
	return s.writtenBytes
}

func (s *OutputStream) GetCheckSum() uint64 {
	return 0
}
