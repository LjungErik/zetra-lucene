package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index/segment"

type PostingsFormat interface {
	GetFieldsConsumer(*segment.SegmentWriteState) (FieldsConsumer, error)
	GetFieldsProducer(*segment.SegmentReadState) (FieldsProducer, error)
}

type StoredFieldsFormat interface {
	GetFieldsConsumer() StoredFieldsConsumer
	GetFieldsProducer() StoredFieldsProducer
}
