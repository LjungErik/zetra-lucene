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
	SEGMENT_PREFIX         = "segments_"
	STORED_FILE_EXTENSION  = ".data"
	TERM_FILE_EXTENSION    = ".term"
	STATICS_FILE_EXTENSION = ".stcs"
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

type SegmentDocumentsMetadata struct {
	DocsLength    map[int]int `json:"docs_length"`
	AvgDocsLength float64     `json:"avg_docs_length"`
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

func (s *Segments) FlushNextSegment(dir directory.Directory) error {
	next := SegmentMetadata{
		SegmentName: s.NextSegmentName(),
		SegmentID:   s.NextSegment,
	}

	s.NextSegment += 1
	s.SegmentCount += 1
	s.Segments = append(s.Segments, next)

	filename := fmt.Sprintf("%s%d", SEGMENT_PREFIX, next.SegmentID)

	out, err := dir.OpenOutputStream(filename)
	if err != nil {
		return err
	}

	if err = json.NewEncoder(out).Encode(s); err != nil {
		return err
	}

	return nil
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

		genStr := strings.TrimPrefix(name, SEGMENT_PREFIX)
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
