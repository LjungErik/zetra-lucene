package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type Lucene104PostingsFormat struct{}

var _ codecs.PostingsFormat = (*Lucene104PostingsFormat)(nil)

func (l *Lucene104PostingsFormat) GetFieldsConsumer(sws *segment.SegmentWriteState) codecs.FieldsConsumer {
	panic("unimplemented")
}

func (l *Lucene104PostingsFormat) GetFieldsProducer(srs *segment.SegmentReadState) codecs.FieldsProducer {
	panic("unimplemented")
}
