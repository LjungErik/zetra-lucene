package codecs

import (
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

type StoredFieldWriter interface {
	Write(docID int, field string, data string)
	Close()
}

type PostingsWriter interface {
	Init(termsOut internal.DataOutputStream, sws *segment.SegmentWriteState) error
	Write(term index.Term) BlockTermState
	EncodeTerm(out internal.DataOutputStream, state BlockTermState) error
	Close() error
}
