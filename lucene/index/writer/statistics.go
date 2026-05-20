package writer

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/index"
)

type StatisticsTermWriter struct {
	fieldsMetadata map[string]*index.SegmentDocumentsMetadata
}

func NewStatisticsTermWriter() *StatisticsTermWriter {
	return &StatisticsTermWriter{
		fieldsMetadata: make(map[string]*index.SegmentDocumentsMetadata),
	}
}

func (w *StatisticsTermWriter) write(docID int, fieldName string, dataLength int) {
	if _, ok := w.fieldsMetadata[fieldName]; !ok {
		w.fieldsMetadata[fieldName] = &index.SegmentDocumentsMetadata{
			DocsLength: make(map[string]int),
		}
	}

	docIDStr := fmt.Sprintf("%d", docID)

	w.fieldsMetadata[fieldName].DocsLength[docIDStr] = dataLength
	w.fieldsMetadata[fieldName].AvgDocsLength = (w.fieldsMetadata[fieldName].AvgDocsLength + float64(dataLength)) / float64(len(w.fieldsMetadata[fieldName].DocsLength))
	w.fieldsMetadata[fieldName].DocumentCount = len(w.fieldsMetadata[fieldName].DocsLength)
}

func (w *StatisticsTermWriter) flush(sws *index.SegementWriteState) (int64, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), index.STATICS_FILE_EXTENSION)

	s, err := sws.Directory.OpenOutputStream(filename)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	err = json.NewEncoder(s).Encode(w.fieldsMetadata)
	if err != nil {
		return 0, err
	}

	return s.GetWrittenBytes(), nil
}
