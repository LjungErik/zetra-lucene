package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index"

type PostingsWriter interface {
	Write(fieldName string, term string, freqency int)
	Close()
	Flush(sws *index.SegementWriteState)
}
