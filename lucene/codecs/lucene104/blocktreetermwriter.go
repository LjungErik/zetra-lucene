package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
	"github.com/LjungErik/zetra-lucene/lucene/internal"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
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
			tw := newTermWriter(l)
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
	pending      []PendingEntry
	prevTerm     string
	parent       *Lucene104BlockTreeTermsWriter
	prefixStarts []int
}

func newTermWriter(parent *Lucene104BlockTreeTermsWriter) *termWriter {
	return &termWriter{
		pending:      make([]PendingEntry, 0, 10),
		parent:       parent,
		prevTerm:     "",
		prefixStarts: make([]int, 8),
	}
}

func (w *termWriter) write(term index.Term) {
	state := w.parent.pw.Write(term)
	w.pushTerm(term.Value())

	w.pending = append(w.pending, PendingEntry{
		TermBytes: []byte(term.Value()),
		State:     state,
	})
}

func (w *termWriter) pushTerm(term string) {
	prefixLength := utils.CommonPrefixLength(w.prevTerm, term)

	for i := len(w.prevTerm); i >= prefixLength; i-- {
		prefixTopSize := len(w.pending) - w.prefixStarts[i]
		if prefixTopSize >= w.parent.minItemsPerBlock {
			// writing of this block
			// reset the prefix start
			w.prefixStarts[i] -= prefixTopSize - 1
		}
	}

	if len(w.prefixStarts) < len(term) {
		w.prefixStarts = utils.Grow(w.prefixStarts, len(term))
	}

	for i := prefixLength; i < len(term); i++ {
		w.prefixStarts[i] = len(w.pending)
	}

	w.prevTerm = term
}

func (w *termWriter) flush() {

}

func (w *termWriter) writeBlocks(prefixLen, count int) {
	lastSuffixLeadLabel := -1

	start := len(w.pending) - count
	end := len(w.pending)
	nextBlockStart := start

	for i := start; i < end; i++ {
		p := w.pending[i]

		var suffixLeadLabel int

		if len(p.TermBytes) == prefixLen {
			suffixLeadLabel = -1
		} else {
			suffixLeadLabel = int(p.TermBytes[prefixLen] & 0xFF)
		}

		if suffixLeadLabel != lastSuffixLeadLabel && i > 0 {
			// Write a block of data
			w.writeBlock(prefixLen, nextBlockStart, i)
			nextBlockStart = i
		}
	}

	if nextBlockStart < count {
		// Still data left that needs to be written to a block
		w.writeBlock(prefixLen, nextBlockStart, end)
	}
}

func (w *termWriter) writeBlock(
	prefixLen int,
	start int,
	end int,
) {
	for i := start; i < end; i++ {
		// Get the suffix for each entry and write to termOut
	}
}

func shortestPrefixLength(prev, first, last string) int {
	minLen := utils.CommonPrefixLength(first, last) + 1
	prevLen := utils.CommonPrefixLength(first, prev) + 1

	if prevLen > minLen {
		return prevLen
	}

	return minLen
}
