package storage

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/utils"
)

const (
	SEGMENT_PREFIX = "segments_"
)

var (
	ErrSegmentMetadataNotFound = errors.New("segment metadata file not found")
)

type DirectoryReader interface {
}

type StandardDirectoryReader struct {
	segments *Segments
	readers  []*SegmentReader
}

var _ DirectoryReader = (*StandardDirectoryReader)(nil)

func OpenStandrardDirectoryReader(directory string) (*StandardDirectoryReader, error) {
	segments, err := getNewestSegmentsMetadata(directory)
	if err != nil {
		return nil, err
	}

	// Create the individual readers for reading each segment
	readers := make([]*SegmentReader, len(segments.Segments))
	for i, seg := range segments.Segments {
		readers[i], err = OpenSegmentReader(seg, directory)
		if err != nil {
			return nil, err
		}
	}

	return &StandardDirectoryReader{
		segments: segments,
		readers:  readers,
	}, nil
}

func getNewestSegmentsMetadata(directory string) (*Segments, error) {
	entries, err := os.ReadDir(directory)
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
		return nil, ErrSegmentMetadataNotFound
	}

	var metadata = &Segments{}
	path := filepath.Join(directory, latestFile)
	if err := utils.ReadJsonFile(path, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}
