package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
	"github.com/LjungErik/zetra-lucene/lucene/internal"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
	"github.com/LjungErik/zetra-lucene/lucene/utils/queue"
)

const (
	minItemsPerBlock = 25
	maxItemsPerBlock = 48
)

type Lucene104BlockTreeTermsWriter struct {
	sws              *segment.SegmentWriteState
	pw               codecs.PostingsWriter
	termsOut         *internal.OutputStream
	minItemsPerBlock int
	maxItemsPerBlock int
}

var _ codecs.FieldsConsumer = (*Lucene104BlockTreeTermsWriter)(nil)

func NewLucene104BlockTreeTermsWriter(sws *segment.SegmentWriteState, writer codecs.PostingsWriter) (*Lucene104BlockTreeTermsWriter, error) {
	termsOut, err := sws.Directory.OpenOutputStream("tes.ttxt")
	if err != nil {
		return nil, err
	}

	return &Lucene104BlockTreeTermsWriter{
		sws:              sws,
		pw:               writer,
		termsOut:         termsOut,
		minItemsPerBlock: minItemsPerBlock,
		maxItemsPerBlock: maxItemsPerBlock,
	}, nil
}

// Write implements [codecs.FieldsConsumer].
func (l *Lucene104BlockTreeTermsWriter) Write(fields index.Fields) error {
	// Go through each of the fields and for each field write the postings and also add it to a block.
	// If a block becomes full then write the block to FST
	for field := range fields.Iter() {
		terms := fields.Terms(field)

		for term := range terms.Terms() {
			tw := newTermWriter(l.termsOut, l.minItemsPerBlock, l.pw)
			tw.write(term)
		}
	}

	return nil
}

func (l *Lucene104BlockTreeTermsWriter) Close() error {
	l.pw.Close()

	return nil
}

type PendingEntry struct {
	TermBytes []byte
	State     codecs.BlockTermState
}

type termWriter struct {
	pending      *queue.Queue[PendingEntry]
	prevLastTerm string
	parent       *Lucene104BlockTreeTermsWriter
}

func newTermWriter(parent *Lucene104BlockTreeTermsWriter) *termWriter {
	return &termWriter{
		pending:      queue.NewQueue[PendingEntry](),
		parent:       parent,
		prevLastTerm: "",
	}
}

func (w *termWriter) write(term index.Term) {
	state := w.parent.pw.Write(term)

	w.pending.Push(PendingEntry{
		TermBytes: []byte(term.Value()),
		State:     state,
	})

	if w.pending.Len() >= w.parent.minItemsPerBlock {
		first := w.pending.Peak()
		for pe, ok := w.pending.Pop(); ok; {
			// Write the actual term to the termsOut
			w.parent.termsOut.Write(pe.TermBytes)
			w.parent.pw.EncodeTerm(w.parent.termsOut, pe.State)
		}

		prefixLen := shortestPrefixLength(w.prevLastTerm, string(first.TermBytes), term.Value())

		// Write the prefix to the FST

		w.prevLastTerm = term.Value()

		// Write the FST with the smallest prefix between first and previous
		// Save the last term name as previous term
	}
}

func (w *termWriter) flush() {

}

func shortestPrefixLength(prev, first, last string) int {
	minLen := utils.CommonPrefixLength(first, last) + 1
	prevLen := utils.CommonPrefixLength(first, prev) + 1

	if prevLen > minLen {
		return prevLen
	}

	return minLen
}
