package stream

import (
	"bytes"
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

type MemBufferDataOutput struct {
	internal.DataOutputStream
	buf *bytes.Buffer
}

var _ internal.DataOutputStream = (*MemBufferDataOutput)(nil)

func NewMemBufferDataOutput() *MemBufferDataOutput {
	stream := &MemBufferDataOutput{
		buf: &bytes.Buffer{},
	}

	stream.DataOutputStream = NewOutputStream(utils.NopWriterCloser(stream.buf))

	return stream
}

func (m *MemBufferDataOutput) Close() error { return nil }

func (m *MemBufferDataOutput) CopyTo(w io.Writer) error {
	if _, err := m.buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

func (m *MemBufferDataOutput) Reset() {
	m.buf.Reset()
}
