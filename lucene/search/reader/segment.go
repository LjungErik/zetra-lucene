package reader

import (
	"encoding/json"
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
)

type SegmentReader struct {
	metadata     index.SegmentMetadata
	index        map[string]map[string][]index.TermCount
	docs         map[string]map[string]string
	docsMetadata map[string]*index.SegmentDocumentsMetadata
}

var _ context.SearchContext = (*SegmentReader)(nil)

func OpenSegmentReader(metadata index.SegmentMetadata, dir directory.Directory) (*SegmentReader, error) {
	var (
		indexFilename  = fmt.Sprintf("%s%s", metadata.SegmentName, index.TERM_FILE_EXTENSION)
		storedFileName = fmt.Sprintf("%s%s", metadata.SegmentName, index.STORED_FILE_EXTENSION)
		statsFilename  = fmt.Sprintf("%s%s", metadata.SegmentName, index.STATICS_FILE_EXTENSION)
	)

	r := &SegmentReader{
		metadata: metadata,
	}

	if err := readJson(dir, indexFilename, &r.index); err != nil {
		return nil, err
	}

	if err := readJson(dir, storedFileName, &r.docs); err != nil {
		return nil, err
	}

	if err := readJson(dir, statsFilename, &r.docsMetadata); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *SegmentReader) GetStatistic(fieldName string) context.SearchStatistics {
	stat := context.SearchStatistics{}

	stat.AverageDataLength = s.docsMetadata[fieldName].AvgDocsLength
	stat.DocumentCount = len(s.docs[fieldName])

	return stat
}

func (s *SegmentReader) GetTermCounts(fieldName, term string) []index.TermCount {
	return s.index[fieldName][term]
}

func (s *SegmentReader) GetDocLength(fieldName string, docID int) int {
	docIDStr := fmt.Sprintf("%d", docID)
	return s.docsMetadata[fieldName].DocsLength[docIDStr]
}

func (s *SegmentReader) GetSegmentID() int {
	return s.metadata.SegmentID
}

func readJson(dir directory.Directory, filename string, v any) error {
	s, err := dir.OpenInputStream(filename)
	if err != nil {
		return err
	}
	defer s.Close()

	if err = json.NewDecoder(s).Decode(v); err != nil {
		return err
	}

	return nil
}
