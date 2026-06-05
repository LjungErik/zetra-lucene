package lucene104

import (
	"github.com/LjungErik/zetra-lucene/lucene/codecs"
	"github.com/LjungErik/zetra-lucene/lucene/codecs/lucene104/constants"
	codec_utils "github.com/LjungErik/zetra-lucene/lucene/codecs/utils"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/filenames"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
	"github.com/LjungErik/zetra-lucene/lucene/internal"
)

const (
	VersionStart   = 0
	VersionCurrent = VersionStart
)

const (
	BlockSize = 128
)

type Lucene104BlockTermState struct {
}

type Lucene104PostingsWriter struct {
	version int

	metaOut internal.DataOutputStream
	docOut  internal.DataOutputStream
	posOut  internal.DataOutputStream
	payOut  internal.DataOutputStream
}

var _ codecs.PostingsWriter = (*Lucene104PostingsWriter)(nil)

func NewLucene104PostingsWriter(sws *segment.SegmentWriteState) (*Lucene104PostingsWriter, error) {
	w := &Lucene104PostingsWriter{
		version: VersionCurrent,
	}

	metaFileName := filenames.SegmentFileName(
		sws.Segments.NextSegmentName(),
		sws.SegmentSuffix(),
		constants.MetaExtension)
	metaOut, err := sws.Directory.OpenOutputStream(metaFileName)
	if err != nil {
		return nil, err
	}

	docFileName := filenames.SegmentFileName(
		sws.Segments.NextSegmentName(),
		sws.SegmentSuffix(),
		constants.DocumentExtension)
	docOut, err := sws.Directory.OpenOutputStream(docFileName)
	if err != nil {
		return nil, err
	}

	if err := codec_utils.WriteIndexHeader(
		metaOut,
		constants.MetaCodec,
		w.version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil, err
	}

	if err := codec_utils.WriteIndexHeader(
		docOut,
		constants.DocumentCodec,
		w.version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil, err
	}

	if sws.FieldsInfo.HasProx() {
		posFileName := filenames.SegmentFileName(
			sws.Segments.NextSegmentName(),
			sws.SegmentSuffix(),
			constants.PosExtension,
		)
		posOut, err := sws.Directory.OpenOutputStream(posFileName)
		if err != nil {
			return nil, err
		}

		if err := codec_utils.WriteIndexHeader(
			posOut,
			constants.PosCodec,
			w.version,
			[]byte(sws.Segments.NextSegmentName()),
			sws.SegmentSuffix(),
		); err != nil {
			return nil, err
		}

		w.posOut = posOut

		// TODO: Add  payload buffer
		// TODO: Add offset buffers

		if sws.FieldsInfo.HasPayload() || sws.FieldsInfo.HasOffsets() {
			payFilename := filenames.SegmentFileName(
				sws.Segments.NextSegmentName(),
				sws.SegmentSuffix(),
				constants.PayExtension,
			)

			payOut, err := sws.Directory.OpenOutputStream(payFilename)
			if err != nil {
				return nil, err
			}

			if err := codec_utils.WriteIndexHeader(
				payOut,
				constants.PayCodec,
				w.version,
				[]byte(sws.Segments.NextSegmentName()),
				sws.SegmentSuffix(),
			); err != nil {
				return nil, err
			}

			w.payOut = payOut
		}
	}

	w.metaOut = metaOut
	w.docOut = docOut

	return w, nil
}

func (w *Lucene104PostingsWriter) Close() error {
	panic("unimplemented")
}

func (w *Lucene104PostingsWriter) EncodeTerm(out internal.DataOutputStream, state codecs.BlockTermState) error {
	panic("unimplemented")
}

func (w *Lucene104PostingsWriter) Init(termsOut internal.DataOutputStream, sws *segment.SegmentWriteState) error {
	if err := codec_utils.WriteIndexHeader(
		termsOut,
		constants.TermsCodec,
		w.version,
		[]byte(sws.Segments.NextSegmentName()),
		sws.SegmentSuffix(),
	); err != nil {
		return nil
	}

	if err := termsOut.WriteVInt(BlockSize); err != nil {
		return err
	}

	return nil
}

func (w *Lucene104PostingsWriter) Write(term index.Term) codecs.BlockTermState {
	panic("unimplemented")
}
