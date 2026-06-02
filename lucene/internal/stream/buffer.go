package stream

import (
	"bytes"
	"io"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type MemBufferDataOutput struct {
	buf bytes.Buffer
}

var _ internal.DataOutputStream = (*MemBufferDataOutput)(nil)

func NewMemBufferDataOutput() *MemBufferDataOutput {
	return &MemBufferDataOutput{}
}

func (m *MemBufferDataOutput) Close() error { return nil }

func (m *MemBufferDataOutput) GetWrittenBytes() uint64 {
	return uint64(m.buf.Len())
}

func (m *MemBufferDataOutput) GetCheckSum() uint64 {
	return 0
}

// Write implements [DataOutputStream].
func (m *MemBufferDataOutput) Write(p []byte) (int, error) {
	return m.buf.Write(p)
}

// WriteByte implements [DataOutputStream].
func (m *MemBufferDataOutput) WriteByte(b byte) error {
	return m.buf.WriteByte(b)
}

func (m *MemBufferDataOutput) WriteVInt(i int) error {
	return writeVInt(m, i)
}

func (m *MemBufferDataOutput) WriteVUInt64(i uint64) error {
	return writeVUInt64(m, i)
}

func (m *MemBufferDataOutput) CopyTo(w io.Writer) error {
	if _, err := m.buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}

func (m *MemBufferDataOutput) Reset() {
	m.buf.Reset()
}
