package reader

import (
	"fmt"

	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/search"
	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

type SegmentReader struct {
	metadata     index.SegmentMetadata
	index        map[string]map[string][]index.TermCount
	docs         map[string]map[int]string
	docsMetadata map[string]*index.SegmentDocumentsMetadata
}

var _ search.SearchContext = (*SegmentReader)(nil)

func OpenSegmentReader(metadata index.SegmentMetadata, dir directory.Directory) (*SegmentReader, error) {
	var (
		indexFilename  = fmt.Sprintf("%s%s", metadata.SegmentName, index.TERM_FILE_EXTENSION)
		storedFileName = fmt.Sprintf("%s%s", metadata.SegmentName, index.STORED_FILE_EXTENSION)
		statsFilename  = fmt.Sprintf("%s%s", metadata.SegmentName, index.STATICS_FILE_EXTENSION)
	)

	r := &SegmentReader{
		metadata: metadata,
	}

	if err := utils.ReadJsonFile(indexFilename, &r.index); err != nil {
		return nil, err
	}

	if err := utils.ReadJsonFile(storedFileName, &r.docs); err != nil {
		return nil, err
	}

	if err := utils.ReadJsonFile(statsFilename, &r.docs); err != nil {
		return nil, err
	}

	return r, nil
}

func (s *SegmentReader) GetStatistic(fieldName string) search.SearchStatistics {
	stat := search.SearchStatistics{}

	stat.AverageDataLength = s.docsMetadata[fieldName].AvgDocsLength
	stat.DocumentCount = len(s.docs[fieldName])

	return stat
}

func (s *SegmentReader) GetTermCounts(fieldName, term string) []index.TermCount {
	return s.index[fieldName][term]
}

func (s *SegmentReader) GetDocLength(fieldName string, docId int) int {
	return s.docsMetadata[fieldName].DocsLength[docId]
}

func (s *SegmentReader) GetSegmentID() int {
	return s.metadata.SegmentID
}
