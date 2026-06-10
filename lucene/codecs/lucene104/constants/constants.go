package constants

import "math"

const (
	MetaExtension     = "psm"
	DocumentExtension = "doc"
	PosExtension      = "pos"
	PayExtension      = "pay"
)

const (
	TermsCodec    = "Lucene104PostingsWriterTerms"
	MetaCodec     = "Lucene104PostingsWriterMeta"
	DocumentCodec = "Lucene104PostingsWriterDoc"
	PosCodec      = "Lucene104PostingsWriterPos"
	PayCodec      = "Lucene104PostingsWriterPay"
)

const (
	MaxPosition = math.MaxInt32 - 128
)
