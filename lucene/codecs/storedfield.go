package codecs

import "github.com/LjungErik/zetra-lucene/lucene/index"

type StoredFieldWriter interface {
	Write(docID int, field string, data string)
	Flush(sws *index.SegementWriteState)
	Close()
}
