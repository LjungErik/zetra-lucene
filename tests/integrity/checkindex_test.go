package integrity

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/document/field/textfield"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/index/writer"
)

var sampleDocs = []struct {
	field string
	data  string
}{
	{"name", "magic document starts here"},
	{"name", "A fish has on average 124 bones"},
	{"name", "A human has 207 bones"},
	{"name", "Fish are great at flying but they don't really survive for long outside of water"},
	{"name", "one fish, two fish, three fish, gold fish"},
}

func writeSampleIndex(t *testing.T, dir string) {
	t.Helper()

	fsDir := directory.OpenFSDirectory(dir)
	indexWriter := writer.NewIndexWriter(fsDir, writer.IndexWriterConfig{
		Analyzer: analyzer.NewPerFieldAnalyzer(),
	})

	for _, d := range sampleDocs {
		doc := document.NewDocument()
		doc.Add(textfield.New(d.field, d.data, true))
		indexWriter.AddDocument(doc)
	}

	require.NoError(t, indexWriter.Flush(), "failed to flush index")
}

func TestCheckIndex(t *testing.T) {
	if !JavaAvailable() {
		t.Skip("java not found on PATH; Lucene 10.4 CheckIndex requires Java 21+ — skipping integration test")
	}

	indexDir := os.Getenv("INTEGRITY_INDEX_DIR")
	if indexDir == "" {
		indexDir = t.TempDir()
	} else {
		require.NoError(t, os.MkdirAll(indexDir, 0o755))
	}

	writeSampleIndex(t, indexDir)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	result, err := RunCheckIndex(ctx, indexDir)
	require.NoError(t, err, "failed to execute Lucene CheckIndex")

	t.Logf("Lucene CheckIndex (exit code %d) for %s:\n%s", result.ExitCode, indexDir, result.Output)

	require.True(t, result.OK,
		"Lucene CheckIndex reported integrity problems with the generated index (exit code %d); see logged output above",
		result.ExitCode)

	os.RemoveAll(indexDir)
}
