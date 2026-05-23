package reader

import (
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
)

type DirectoryReader interface {
	GetSegments() []context.LeafReaderContext
}

type StandardDirectoryReader struct {
	readers []context.LeafReaderContext
}

var _ DirectoryReader = (*StandardDirectoryReader)(nil)

func OpenStandrardDirectoryReader(dir directory.Directory) (*StandardDirectoryReader, error) {
	metadata, err := index.GetNewestSegment(dir)
	if err != nil {
		return nil, err
	}

	readers := make([]context.LeafReaderContext, len(metadata.Segments))
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

func (r *StandardDirectoryReader) GetSegments() []context.LeafReaderContext {
	return r.readers
}
