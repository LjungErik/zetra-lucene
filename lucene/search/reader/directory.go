package reader

import (
	"sort"

	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/search/document"
	"github.com/LjungErik/zetra-lucene/lucene/search/query"
)

type DirectoryReader interface {
	Query(query.Query, int) []document.TopDoc
}

type StandardDirectoryReader struct {
	readers []*SegmentReader
}

var _ DirectoryReader = (*StandardDirectoryReader)(nil)

func OpenStandrardDirectoryReader(dir directory.Directory) (*StandardDirectoryReader, error) {
	metadata, err := index.GetNewestSegment(dir)
	if err != nil {
		return nil, err
	}

	readers := make([]*SegmentReader, len(metadata.Segments))
	for i, seg := range metadata.Segments {
		readers[i], err = OpenSegmentReader(seg, dir)
		if err != nil {
			return nil, err
		}
	}

	return &StandardDirectoryReader{
		readers: readers,
	}, nil
}

func (r *StandardDirectoryReader) Query(q query.Query, n int) []document.TopDoc {
	ret := make([]document.TopDoc, 0, n)

	for _, segReader := range r.readers {
		docs := q.Execute(segReader)
		if docs != nil {
			ret = append(ret, docs...)
		}
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].Score > ret[j].Score
	})

	if len(ret) < n {
		return ret
	}

	return ret[:n]
}
