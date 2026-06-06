package stream

import "github.com/LjungErik/zetra-lucene/lucene/internal"

type FailFastStream struct {
	stream internal.DataOutputStream
	err    error
}

var _ internal.ErrStream = (*FailFastStream)(nil)

func NewFailFastStream(s internal.DataOutputStream) *FailFastStream {
	return &FailFastStream{
		stream: s,
		err:    nil,
	}
}

func (f *FailFastStream) Write(p []byte) {
	if f.err != nil {
		return
	}

	if _, err := f.stream.Write(p); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteByte(b byte) {
	if f.err != nil {
		return
	}

	if err := f.stream.WriteByte(b); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteInt(i int) {
	if f.err != nil {
		return
	}

	if err := writeInt(f.stream, i); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteInt64(i int64) {
	if f.err != nil {
		return
	}

	if err := writeInt64(f.stream, i); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteUInt64(i uint64) {
	if f.err != nil {
		return
	}

	if err := writeUInt64(f.stream, i); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteVInt(i int) {
	if f.err != nil {
		return
	}

	if err := writeVInt(f.stream, i); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteVUInt64(i uint64) {
	if f.err != nil {
		return
	}

	if err := writeVUInt64(f.stream, i); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) WriteFunc(writeFunc func(internal.DataOutputStream) error) {
	if f.err != nil {
		return
	}

	if err := writeFunc(f.stream); err != nil {
		f.err = err
	}
}

func (f *FailFastStream) Error() error {
	return f.err
}
