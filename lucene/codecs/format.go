package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index/segment"

type PostingsFormat interface {
	GetFieldsConsumer(*segment.SegmentWriteState) FieldsConsumer
	GetFieldsProducer(*segment.SegmentReadState) FieldsProducer
}

type StoredFieldsFormat interface {
	GetFieldsConsumer() StoredFieldsConsumer
	GetFieldsProducer() StoredFieldsProducer
}
