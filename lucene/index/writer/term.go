package writer

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/index"
)

const (
	termFileExtension = ".term"
)

type TermWriter struct {
	fieldsCount map[string]map[string]map[int]int
}

func NewTermWriter() *TermWriter {
	return &TermWriter{
		fieldsCount: make(map[string]map[string]map[int]int),
	}
}

func (w *TermWriter) write(docID int, fieldName string, tokens []analysis.Token) {
	if _, ok := w.fieldsCount[fieldName]; !ok {
		w.fieldsCount[fieldName] = make(map[string]map[int]int)
	}

	for _, token := range tokens {
		if _, ok := w.fieldsCount[fieldName][token.Text]; !ok {
			w.fieldsCount[fieldName][token.Text] = make(map[int]int)
		}

		w.fieldsCount[fieldName][token.Text][docID]++
	}
}

func (w *TermWriter) flush(sws *index.SegementWriteState) (int64, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), termFileExtension)

	s, err := sws.Directory.OpenOutputStream(filename)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	err = json.NewEncoder(s).Encode(w.fieldsCount)
	if err != nil {
		return 0, err
	}

	return s.GetWrittenBytes(), nil
}
