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
	WriteTerm(term index.Term) BlockTermState
	EncodeTerm(out internal.DataOutputStream, state BlockTermState) error
	Close() error
}

type BasePostingsWriter struct {
	encoder PostingsEncoder
}

// Close implements [PostingsWriter].

var _ PostingsWriter = (*BasePostingsWriter)(nil)

func NewBasePostingsWriter(encoder PostingsEncoder) *BasePostingsWriter {
	return &BasePostingsWriter{
		encoder: encoder,
	}
}

func (b *BasePostingsWriter) Close() error {
	return nil
}

func (b *BasePostingsWriter) EncodeTerm(out internal.DataOutputStream, state BlockTermState) error {
	return nil
}

func (b *BasePostingsWriter) Init(termsOut internal.DataOutputStream, sws *segment.SegmentWriteState) error {
	return nil
}

func (b *BasePostingsWriter) WriteTerm(term index.Term) BlockTermState {
	b.encoder.StartTerm()

	var docFreq uint32
	var totTermFreq uint64

	for posting := range term.Postings() {
		b.encoder.StartDoc(posting.DocID(), posting.Frequency())

		totTermFreq += uint64(posting.Frequency())
		docFreq++

		for _, pos := range posting.Positions() {
			// TODO: Add correct handling of add postion
			b.encoder.AddPosition(pos, nil)
		}

		b.encoder.FinishDoc()
	}

	b.encoder.FinishTerm()

	return BlockTermState{
		DocumentFrequency:  docFreq,
		TotalTermFrequency: totTermFreq,
	}
}
