package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
)

const (
	SEGMENT_PREFIX = "segments_"
)

var (
	ErrSegmentMetadataNotFound = errors.New("segment metadata file not found")
)

type Segments struct {
	NextSegment  int               `json:"next_segment"`
	SegmentCount int               `json:"segment_count"`
	Segments     []SegmentMetadata `json:"segments"`
}

type SegmentMetadata struct {
	SegmentName string `json:"segment_name"`
	SegmentID   int    `json:"segment_id"`
}

type SegementWriteState struct {
	Segments  *Segments
	Directory directory.Directory
}

func CreateNewSegmentWriteState(dir directory.Directory) (*SegementWriteState, error) {
	s, err := getNewestSegmentsMetadata(dir)
	if err != nil {
		return nil, err
	}

	return &SegementWriteState{
		Directory: dir,
		Segments:  s,
	}, nil
}

func defaultSegment() *Segments {
	return &Segments{
		NextSegment:  0,
		SegmentCount: 0,
		Segments:     []SegmentMetadata{},
	}
}

func GetNewestSegment(dir directory.Directory) (*Segments, error) {
	return getNewestSegmentsMetadata(dir)
}

func (s *Segments) NextSegmentName() string {
	return fmt.Sprintf("_%d", s.NextSegment)
}

func (s *Segments) AddNextSegment() string {
	next := SegmentMetadata{
		SegmentName: s.NextSegmentName(),
		SegmentID:   s.NextSegment,
	}

	s.NextSegment += 1
	s.SegmentCount += 1
	s.Segments = append(s.Segments, next)

	return next.SegmentName
}

func getNewestSegmentsMetadata(dir directory.Directory) (*Segments, error) {
	entries, err := dir.GetEntries()
	if err != nil {
		return nil, err
	}

	var maxGen int64 = -1
	var latestFile string

	for _, entry := range entries {
		name := entry.Name()
		if !strings.HasPrefix(name, SEGMENT_PREFIX) {
			continue
		}

		genStr := strings.TrimPrefix(name, "segments_")
		gen, err := strconv.ParseInt(genStr, 10, 64)
		if err != nil {
			continue
		}

		if gen > maxGen {
			maxGen = gen
			latestFile = name
		}
	}

	if maxGen == -1 {
		return defaultSegment(), nil
	}

	var metadata = &Segments{}
	s, err := dir.OpenInputStream(latestFile)
	if err := json.NewDecoder(s).Decode(metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
