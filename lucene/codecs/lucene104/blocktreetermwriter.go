package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type Lucene104BlockTreeTermsWriter struct {
	sws    *segment.SegmentWriteState
	writer codecs.PostingsWriter
}

var _ codecs.FieldsConsumer = (*Lucene104BlockTreeTermsWriter)(nil)

func NewLucene104BlockTreeTermsWriter(sws *segment.SegmentWriteState, writer codecs.PostingsWriter) *Lucene104BlockTreeTermsWriter {
	return &Lucene104BlockTreeTermsWriter{
		sws:    sws,
		writer: writer,
	}
}

// Write implements [codecs.FieldsConsumer].
func (l *Lucene104BlockTreeTermsWriter) Write(fields index.Fields) error {
	panic("unimplemented")
}

func (l *Lucene104BlockTreeTermsWriter) Close() error {
	panic("unimplemented")
}
