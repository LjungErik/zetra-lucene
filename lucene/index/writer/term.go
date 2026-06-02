package writer

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/LjungErik/zetra-lucene/lucene/analysis"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
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

func (w *TermWriter) flush(sws *segment.SegmentWriteState) (int, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), segment.TERM_FILE_EXTENSION)

	s, err := sws.Directory.OpenOutputStream(filename)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	fieldIndex := make(map[string]map[string][]index.TermCount, len(w.fieldsCount))

	for k, m := range w.fieldsCount {
		fieldIndex[k] = make(map[string][]index.TermCount, len(w.fieldsCount[k]))
		for text, counts := range m {
			termCounts := make([]index.TermCount, 0, len(counts))
			for docId, count := range counts {
				termCounts = append(termCounts, index.TermCount{
					DocumentID: docId,
					Count:      count,
				})
			}

			sort.SliceStable(termCounts, func(i, j int) bool {
				return termCounts[i].DocumentID < termCounts[j].DocumentID
			})

			fieldIndex[k][text] = termCounts
		}
	}

	err = json.NewEncoder(s).Encode(fieldIndex)
	if err != nil {
		return 0, err
	}

	return s.GetWrittenBytes(), nil
}
