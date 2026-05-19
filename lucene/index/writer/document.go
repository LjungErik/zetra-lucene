package writer

import (
	"errors"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

var (
	ErrUnsupportedFieldType = errors.New("unsupported field type")
)

type DocumentWriter struct {
	stored   *StoredWriter
	term     *TermWriter
	stats    *StatisticsTermWriter
	counter  *utils.Counter
	analyzer *analyzer.PerFieldAnalyzer
}

func NewDocumentWriter(analyzer *analyzer.PerFieldAnalyzer) *DocumentWriter {
	return &DocumentWriter{
		stored:   &StoredWriter{},
		term:     &TermWriter{},
		stats:    &StatisticsTermWriter{},
		counter:  &utils.Counter{},
		analyzer: analyzer,
	}
}

func (w *DocumentWriter) addDocuments(docs []document.IndexableDocument) error {
	for _, doc := range docs {
		docId := w.counter.GetNextID()

		for field := range doc.Iter() {
			if field.Indexable() {
				switch field.Type() {

				case document.Text:
					tokens := w.analyzer.Get(field.Name()).Analyze(field.ValueAsString())
					w.term.write(docId, field.Name(), tokens)

					w.stats.write(docId, field.Name(), field.ValueLength())

				default:
					return ErrUnsupportedFieldType

				}
			}

			if field.Stored() {
				w.stored.write(docId, field)
			}
		}
	}

	return nil
}

func (w *DocumentWriter) flush(sws *index.SegementWriteState) error {
	var bytesWritten int64 = 0
	n, err := w.term.flush(sws)
	if err != nil {
		return err
	}
	bytesWritten += n

	n, err = w.stored.flush(sws)
	if err != nil {
		return err
	}
	bytesWritten += n

	n, err = w.stats.flush(sws)
	if err != nil {
		return err
	}
	bytesWritten += n

	fmt.Printf("Written %d bytes\n", bytesWritten)

	return nil
}
