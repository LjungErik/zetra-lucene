package stream

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type InputStream struct {
	reader    io.ReadCloser
	bufReader *bufio.Reader
}

var _ io.ReadCloser = (*InputStream)(nil)
var _ internal.DataInputStream = (*InputStream)(nil)

func NewInputStream(reader io.ReadCloser) *InputStream {
	return &InputStream{
		reader:    reader,
		bufReader: bufio.NewReader(reader),
	}
}

func (s *InputStream) Close() error {
	return s.reader.Close()
}

func (s *InputStream) Read(p []byte) (int, error) {
	return s.bufReader.Read(p)
}

func (s *InputStream) ReadByte() (byte, error) {
	return s.bufReader.ReadByte()
}

func (s *InputStream) ReadUInts(dst []uint32) error {
	var buf [4]byte
	for i := range dst {
		if _, err := io.ReadFull(s.bufReader, buf[:]); err != nil {
			return err
		}
		dst[i] = binary.LittleEndian.Uint32(buf[:])
	}

	return nil
}
