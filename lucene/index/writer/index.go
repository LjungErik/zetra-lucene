package writer

import (
	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
	"github.com/LjungErik/zetra-lucene/lucene/index/directory"
	"github.com/LjungErik/zetra-lucene/lucene/index/segment"
)

type IndexWriter struct {
	writer    *DocumentWriter
	directory directory.Directory
}

type IndexWriterConfig struct {
	Analyzer *analyzer.PerFieldAnalyzer
}

func NewIndexWriter(dir directory.Directory, config IndexWriterConfig) *IndexWriter {
	return &IndexWriter{
		writer:    NewDocumentWriter(config.Analyzer),
		directory: dir,
	}
}

func (w *IndexWriter) AddDocument(doc *document.IndexableDocument) {
	w.writer.addDocuments([]*document.IndexableDocument{doc})
}

func (w *IndexWriter) Flush() error {
	state, err := segment.CreateNewSegmentWriteState(w.directory)
	if err != nil {
		return err
	}

	err = w.writer.flush(state)
	if err != nil {
		return err
	}

	err = state.Segments.FlushNextSegment(w.directory)
	if err != nil {
		return err
	}

	return nil
}
