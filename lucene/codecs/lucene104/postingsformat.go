package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/codecs/lucene103/blocktree"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type Lucene104PostingsFormat struct{}

var _ codecs.PostingsFormat = (*Lucene104PostingsFormat)(nil)

func (l *Lucene104PostingsFormat) GetFieldsConsumer(sws *segment.SegmentWriteState) (codecs.FieldsConsumer, error) {
	postingsWriter, err := NewLucene104PostingsWriter(sws)
	if err != nil {
		return nil, err
	}

	return blocktree.NewLucene103BlockTreeTermsWriter(sws, postingsWriter)
}

func (l *Lucene104PostingsFormat) GetFieldsProducer(srs *segment.SegmentReadState) (codecs.FieldsProducer, error) {
	panic("unimplemented")
}
