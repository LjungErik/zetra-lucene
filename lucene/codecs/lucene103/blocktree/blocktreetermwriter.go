package blocktree

import (
	"bytes"

	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/codecs/lucene103/blocktree/compression"
	"github.com/LjungErik/zetra-lucene/lucene/codecs/lucene103/blocktree/constants"
	codec_utils "github.com/LjungErik/zetra-lucene/lucene/codecs/utils"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/filenames"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
	"github.com/LjungErik/zetra-lucene/lucene/internal"
	"github.com/LjungErik/zetra-lucene/lucene/internal/stream"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

const (
	minItemsPerBlock = 25
	maxItemsPerBlock = 48
)

type Lucene103BlockTreeTermsWriter struct {
	sws              *segment.SegmentWriteState
	pw               codecs.PostingsWriter
	termsOut         internal.DataOutputStream
	indexOut         internal.DataOutputStream
	metaOut          internal.DataOutputStream
	minItemsPerBlock int
	maxItemsPerBlock int
	closed           bool
	fields           []*stream.MemBufferDataOutput
}

var _ codecs.FieldsConsumer = (*Lucene103BlockTreeTermsWriter)(nil)

func NewLucene103BlockTreeTermsWriter(sws *segment.SegmentWriteState, writer codecs.PostingsWriter) (*Lucene103BlockTreeTermsWriter, error) {
	version := constants.VersionCurrent

	termsFilename := filenames.SegmentFileName(
		sws.Segments.NextSegmentName(),
		sws.SegmentSuffix(),
		constants.TermsExtension)
	termsOut, err := sws.Directory.OpenOutputStream(termsFilename)
	if err != nil {
		return nil, err
	}

	if err = codec_utils.WriteIndexHeader(
		termsOut,
		constants.TermsCodeName,
		version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil, err
	}

	indexFileName := filenames.SegmentFileName(
		sws.Segments.NextSegmentName(),
		sws.SegmentSuffix(),
		constants.TermsIndexExtension)
	indexOut, err := sws.Directory.OpenOutputStream(indexFileName)
	if err != nil {
		return nil, err
	}

	if err = codec_utils.WriteIndexHeader(
		indexOut,
		constants.TermsIndexCodecName,
		version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil, err
	}

	metaFileName := filenames.SegmentFileName(
		sws.Segments.NextSegmentName(),
		sws.SegmentSuffix(),
		constants.TermsMetaExtension)
	metaOut, err := sws.Directory.OpenOutputStream(metaFileName)
	if err != nil {
		return nil, err
	}

	if err = codec_utils.WriteIndexHeader(
		metaOut,
		constants.TermsMetaCodecName,
		version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil, err
	}

	if err := writer.Init(metaOut, sws); err != nil {
		return nil, err
	}

	return &Lucene103BlockTreeTermsWriter{
		sws:              sws,
		pw:               writer,
		termsOut:         termsOut,
		indexOut:         indexOut,
		metaOut:          metaOut,
		minItemsPerBlock: minItemsPerBlock,
		maxItemsPerBlock: maxItemsPerBlock,
		closed:           false,
		fields:           make([]*stream.MemBufferDataOutput, 0),
	}, nil
}

// Write implements [codecs.FieldsConsumer].
func (l *Lucene103BlockTreeTermsWriter) Write(fields index.Fields) error {
	// Go through each of the fields and for each field write the postings and also add it to a block.
	// If a block becomes full then write the block to FST
	for field := range fields.Iter() {
		terms := fields.Terms(field)
		tw := newTermWriter(l)

		for term := range terms.Terms() {

			tw.write(term)
		}

		tw.finish()
	}

	return nil
}

func (l *Lucene103BlockTreeTermsWriter) Close() error {
	if l.closed {
		return nil
	}
	l.closed = true

	if err := l.metaOut.WriteVInt(len(l.fields)); err != nil {
		return err
	}

	for _, fieldMeta := range l.fields {
		if err := fieldMeta.CopyTo(l.metaOut); err != nil {
			return err
		}
	}

	if err := codec_utils.WriteFooter(l.indexOut); err != nil {
		return err
	}

	if err := l.metaOut.WriteVUInt64(l.indexOut.GetWrittenBytes()); err != nil {
		return err
	}

	if err := codec_utils.WriteFooter(l.termsOut); err != nil {
		return err
	}

	if err := l.termsOut.WriteVUInt64(l.termsOut.GetWrittenBytes()); err != nil {
		return err
	}

	if err := codec_utils.WriteFooter(l.metaOut); err != nil {
		return err
	}

	if err := utils.CloseAll(l.metaOut, l.termsOut, l.indexOut, l.pw); err != nil {
		return err
	}

	return nil
}

type PendingEntry struct {
	TermBytes []byte
	State     codecs.BlockTermState
}

type termWriter struct {
	pending        []PendingEntry
	firstTermBytes []byte
	lastTermBytes  []byte
	lastTerm       string
	parent         *Lucene103BlockTreeTermsWriter
	prefixStarts   []int

	statsWriter        *stream.MemBufferDataOutput
	suffixLengthWriter *stream.MemBufferDataOutput
	suffixWriter       bytes.Buffer
	metaWriter         *stream.MemBufferDataOutput

	sumDocumentFrequency  uint64
	sumTotalTermFrequency uint64
	numTerms              uint64
}

func newTermWriter(parent *Lucene103BlockTreeTermsWriter) *termWriter {
	return &termWriter{
		pending:      make([]PendingEntry, 0, 10),
		parent:       parent,
		lastTerm:     "",
		prefixStarts: make([]int, 8),

		statsWriter:        stream.NewMemBufferDataOutput(),
		suffixLengthWriter: stream.NewMemBufferDataOutput(),
		metaWriter:         stream.NewMemBufferDataOutput(),
	}
}

func (w *termWriter) write(term index.Term) {
	state := w.parent.pw.WriteTerm(term)
	w.pushTerm(term.Value())

	entry := PendingEntry{
		TermBytes: []byte(term.Value()),
		State:     state,
	}

	w.pending = append(w.pending, entry)

	w.sumDocumentFrequency += uint64(state.DocumentFrequency)
	w.sumTotalTermFrequency += state.TotalTermFrequency
	w.numTerms++

	if w.firstTermBytes == nil {
		w.firstTermBytes = entry.TermBytes
	}
	w.lastTermBytes = entry.TermBytes
}

func (w *termWriter) pushTerm(term string) {
	prefixLength := utils.CommonPrefixLength(w.lastTerm, term)

	for i := len(w.lastTerm); i >= prefixLength; i-- {
		prefixTopSize := len(w.pending) - w.prefixStarts[i]
		if prefixTopSize >= w.parent.minItemsPerBlock {
			// writing of this block
			// reset the prefix start
			w.writeBlocks(i+1, prefixTopSize)
			w.prefixStarts[i] -= prefixTopSize - 1
		}
	}

	if len(w.prefixStarts) < len(term) {
		w.prefixStarts = utils.Grow(w.prefixStarts, len(term))
	}

	for i := prefixLength; i < len(term); i++ {
		w.prefixStarts[i] = len(w.pending)
	}

	w.lastTerm = term
}

func (w *termWriter) finish() error {
	if w.numTerms <= 0 {
		return nil
	}

	w.writeBlocks(0, len(w.pending))

	metaBuffOut := stream.NewMemBufferDataOutput()
	w.parent.fields = append(w.parent.fields, metaBuffOut)

	if err := metaBuffOut.WriteVUInt64(w.numTerms); err != nil {
		return err
	}

	// Add if statement to check index options
	if err := metaBuffOut.WriteVUInt64(w.sumTotalTermFrequency); err != nil {
		return err
	}

	if err := metaBuffOut.WriteVUInt64(w.sumDocumentFrequency); err != nil {
		return err
	}

	// TODO: write docs seen cardinality

	if err := writeBytesRef(metaBuffOut, w.firstTermBytes); err != nil {
		return err
	}

	if err := writeBytesRef(metaBuffOut, w.lastTermBytes); err != nil {
		return err
	}

	// Save this on the root index

	return nil
}

func (w *termWriter) writeBlocks(prefixLen, count int) error {
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
			if err := w.writeBlock(prefixLen, nextBlockStart, i); err != nil {
				return err
			}
			nextBlockStart = i
		}
	}

	if nextBlockStart < count {
		// Still data left that needs to be written to a block
		if err := w.writeBlock(prefixLen, nextBlockStart, end); err != nil {
			return err
		}
	}

	return nil
}

func (w *termWriter) writeBlock(
	prefixLen int,
	start int,
	end int,
) error {

	numEntries := end - start
	code := numEntries << 1
	if end == len(w.pending) {
		code |= 1
	}

	if err := w.parent.termsOut.WriteVInt(code); err != nil {
		return err
	}

	var statsWriter = NewStatsWriter(w.statsWriter, true)

	// Handle simple case where blocks only contain terms
	for i := start; i < end; i++ {
		// Get the suffix for each entry and write to termOut
		p := w.pending[i]

		suffixLen := len(p.TermBytes) - prefixLen
		if err := w.suffixLengthWriter.WriteVInt(suffixLen); err != nil {
			return err
		}

		if _, err := w.suffixWriter.Write(p.TermBytes[prefixLen:]); err != nil {
			return err
		}

		if err := statsWriter.Add(p.State.DocumentFrequency, p.State.TotalTermFrequency); err != nil {
			return err
		}

		if err := w.parent.pw.EncodeTerm(w.metaWriter, p.State); err != nil {
			return err
		}
	}

	if err := statsWriter.Finish(); err != nil {
		return err
	}

	compressionAlgo := compression.NoCompression

	var token uint64 = uint64(w.suffixWriter.Len()) << 3
	token |= 0x04
	token |= uint64(compressionAlgo)

	if err := w.parent.termsOut.WriteVUInt64(token); err != nil {
		return err
	}

	if err := w.writeSuffix(compressionAlgo); err != nil {
		return err
	}

	if err := w.writeSuffixLength(); err != nil {
		return err
	}

	if err := w.writeStats(); err != nil {
		return err
	}

	if err := w.writeMeta(); err != nil {
		return err
	}

	return nil
}

func (w *termWriter) writeSuffix(compressionAlgo compression.CompressionAlgorithm) error {
	if compressionAlgo == compression.NoCompression {
		if _, err := w.suffixWriter.WriteTo(w.parent.termsOut); err != nil {
			return err
		}
	} else {
		// TODO handle case where compression needs to happen
	}

	w.suffixWriter.Reset()

	return nil
}

func (w *termWriter) writeSuffixLength() error {
	numSuffixBytes := w.suffixLengthWriter.GetWrittenBytes()
	// TODO: Improve handling of spareBytes to avoid reallocation every time block is written
	var spareBytes bytes.Buffer

	if err := w.suffixLengthWriter.CopyTo(&spareBytes); err != nil {
		return err
	}

	w.suffixLengthWriter.Reset()

	// TODO: add Ignore all equal check to minimize write
	if err := w.parent.termsOut.WriteVInt(int(numSuffixBytes << 1)); err != nil {
		return err
	}

	if _, err := w.parent.termsOut.Write(spareBytes.Bytes()); err != nil {
		return err
	}

	return nil
}

func (w *termWriter) writeStats() error {
	numStatsBytes := w.statsWriter.GetWrittenBytes()
	if err := w.parent.indexOut.WriteVInt(int(numStatsBytes)); err != nil {
		return err
	}

	if err := w.statsWriter.CopyTo(w.parent.termsOut); err != nil {
		return err
	}

	w.statsWriter.Reset()

	return nil
}

func (w *termWriter) writeMeta() error {
	numMetaBytes := w.metaWriter.GetWrittenBytes()
	if err := w.parent.termsOut.WriteVInt(int(numMetaBytes)); err != nil {
		return err
	}

	if err := w.metaWriter.CopyTo(w.parent.termsOut); err != nil {
		return err
	}

	w.metaWriter.Reset()

	return nil
}

func writeBytesRef(out internal.DataOutputStream, data []byte) error {
	if err := out.WriteVInt(len(data)); err != nil {
		return err
	}

	if _, err := out.Write(data); err != nil {
		return err
	}

	return nil
}
