package reader

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LjungErik/zetra-lucene/lucene/storage"
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
	segments *storage.Segments
}

var _ DirectoryReader = (*StandardDirectoryReader)(nil)

func Open(directory string) (*StandardDirectoryReader, error) {
	segments, err := getNewestSegmentsMetadata(directory)
	if err != nil {
		return nil, err
	}

	return &StandardDirectoryReader{
		segments: segments,
	}, nil
}

func getNewestSegmentsMetadata(directory string) (*storage.Segments, error) {
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

	var metadata = &storage.Segments{}
	path := filepath.Join(directory, latestFile)
	if err := readJsonFile(path, metadata); err != nil {
		return nil, err
	}

	return metadata, nil
}

func readJsonFile(filename string, v any) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	if err = json.NewDecoder(f).Decode(v); err != nil {
		return err
	}

	return nil
}
