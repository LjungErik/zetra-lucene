package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index"

type TermsWriter interface {
	Write(fieldName string, term string)
	Close()
	Flush(sws *index.SegementWriteState)
}
