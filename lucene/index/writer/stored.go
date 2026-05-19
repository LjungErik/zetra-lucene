package writer

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/index"
)

const (
	storedFileExtension = ".data"
)

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

func (w *StoredWriter) flush(sws *index.SegementWriteState) (int64, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), storedFileExtension)

	s, err := sws.Directory.OpenOutputStream(filename)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	err = json.NewEncoder(s).Encode(w.fieldDocs)
	if err != nil {
		return 0, err
	}

	return s.GetWrittenBytes(), nil
}
