package internal

import "io"

type OutputStream struct {
	writtenBytes int64
	writer       io.WriteCloser
}

var _ io.WriteCloser = (*OutputStream)(nil)

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

func (s *OutputStream) Close() error {
	return s.writer.Close()
}

func (s *OutputStream) GetWrittenBytes() int64 {
	return s.writtenBytes
}

type InputStream struct {
	reader io.ReadCloser
}

var _ io.ReadCloser = (*InputStream)(nil)

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
