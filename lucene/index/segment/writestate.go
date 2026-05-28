package segment

import "github.com/LjungErik/zetra-lucene/lucene/index/directory"

type SegmentWriteState struct {
	Segments      *Segments
	Directory     directory.Directory
	segmentSuffix string
}

func CreateNewSegmentWriteState(dir directory.Directory) (*SegmentWriteState, error) {
	s, err := getNewestSegmentsMetadata(dir)
	if err != nil {
		return nil, err
	}

	return &SegmentWriteState{
		Directory: dir,
		Segments:  s,
	}, nil
}

func (sws *SegmentWriteState) SegmentSuffix() string {
	return sws.segmentSuffix
}
