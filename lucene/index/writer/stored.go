package writer

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type StoredWriter struct {
	fieldDocs map[string]map[string]string
}

func NewStoredWriter() *StoredWriter {
	return &StoredWriter{
		fieldDocs: make(map[string]map[string]string),
	}
}

func (w *StoredWriter) write(docID int, field document.DocumentField) {
	if _, ok := w.fieldDocs[field.Name()]; !ok {
		w.fieldDocs[field.Name()] = make(map[string]string)
	}

	docIDStr := fmt.Sprintf("%d", docID)

	w.fieldDocs[field.Name()][docIDStr] = field.ValueAsString()
}

func (w *StoredWriter) flush(sws *segment.SegmentWriteState) (uint64, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), segment.STORED_FILE_EXTENSION)

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
