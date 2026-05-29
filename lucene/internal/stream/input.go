package stream

import (
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type InputStream struct {
	reader io.ReadCloser
}

var _ io.ReadCloser = (*InputStream)(nil)
var _ internal.DataInputStream = (*InputStream)(nil)

func NewInputStream(reader io.ReadCloser) *InputStream {
	return &InputStream{
		reader: reader,
	}
}

func (i *InputStream) Close() error {
	return i.reader.Close()
}

func (i *InputStream) Read(p []byte) (n int, err error) {
	return i.reader.Read(p)
}
