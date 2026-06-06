package utils

import "io"

type nopWriterCloser struct {
	io.Writer
}

func NopWriterCloser(w io.Writer) io.WriteCloser {
	return nopWriterCloser{
		Writer: w,
	}
}

func (nopWriterCloser) Close() error { return nil }
