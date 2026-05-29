package stream

import (
	"bytes"

	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type MemBufferDataOutput struct {
	buf bytes.Buffer
}

var _ internal.DataOutputStream = (*MemBufferDataOutput)(nil)

func (m *MemBufferDataOutput) Close() error { return nil }

func (m *MemBufferDataOutput) GetWrittenBytes() int64 {
	return 0
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

// WriteVInt implements [DataOutputStream].
func (m *MemBufferDataOutput) WriteVInt(i int) error {
	panic("unimplemented")
}
