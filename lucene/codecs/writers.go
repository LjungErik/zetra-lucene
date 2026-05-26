package codecs

type StoredFieldWriter interface {
	Write(docID int, field string, data string)
	Close()
}

type PostingsWriter interface {
	Write(fieldName string, term string, freqency int)
	Close()
}
