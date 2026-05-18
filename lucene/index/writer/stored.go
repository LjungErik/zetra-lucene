package writer

import "github.com/LjungErik/zetra-lucene/lucene/document"

type StoredWriter struct {
	fieldDocs map[string]map[int]string
}

func NewStoredWriter() *StoredWriter {
	return &StoredWriter{
		fieldDocs: make(map[string]map[int]string),
	}
}

func (w *StoredWriter) write(docId int, field document.DocumentField) {
	w.fieldDocs[field.Name()][docId] = field.ValueAsString()
}

func (w *StoredWriter) flush() {

}
