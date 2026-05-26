package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index"

type FieldsConsumer interface {
	Write(fields index.Fields) error
	Close() error
}

type FieldsProducer interface {
}
