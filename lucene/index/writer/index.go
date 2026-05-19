package writer

import (
	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/index"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
)

type IndexWriter struct {
	writer    *DocumentWriter
	directory directory.Directory
}

type IndexWriterConfig struct {
	analyzer *analyzer.PerFieldAnalyzer
}

func NewIndexWriter(Directory directory.Directory, config IndexWriterConfig) *IndexWriter {
	return &IndexWriter{
		writer: NewDocumentWriter(config.analyzer),
	}
}

func (w *IndexWriter) AddDocument(doc document.IndexableDocument) {
	w.writer.addDocuments([]document.IndexableDocument{doc})
}

func (w *IndexWriter) Flush() error {
	// Create SegmentWriteState
	state, err := index.CreateNewSegmentWriteState(w.directory)
	if err != nil {
		return err
	}

	return w.writer.flush(state)
}
