package writer

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type StatisticsTermWriter struct {
	fieldsMetadata map[string]*segment.SegmentDocumentsMetadata
}

func NewStatisticsTermWriter() *StatisticsTermWriter {
	return &StatisticsTermWriter{
		fieldsMetadata: make(map[string]*segment.SegmentDocumentsMetadata),
	}
}

func (w *StatisticsTermWriter) write(docID int, fieldName string, dataLength int) {
	if _, ok := w.fieldsMetadata[fieldName]; !ok {
		w.fieldsMetadata[fieldName] = &segment.SegmentDocumentsMetadata{
			DocsLength: make(map[string]int),
		}
	}

	docIDStr := fmt.Sprintf("%d", docID)

	w.fieldsMetadata[fieldName].DocsLength[docIDStr] = dataLength
	w.fieldsMetadata[fieldName].AvgDocsLength = (w.fieldsMetadata[fieldName].AvgDocsLength + float64(dataLength)) / float64(len(w.fieldsMetadata[fieldName].DocsLength))
	w.fieldsMetadata[fieldName].DocumentCount = len(w.fieldsMetadata[fieldName].DocsLength)
}

func (w *StatisticsTermWriter) flush(sws *segment.SegmentWriteState) (uint64, error) {
	filename := fmt.Sprintf("%s%s", sws.Segments.NextSegmentName(), segment.STATICS_FILE_EXTENSION)

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
