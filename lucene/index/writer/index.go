package writer

import (
	"github.com/LjungErik/zetra-lucene/lucene/analysis/analyzer"
	"github.com/LjungErik/zetra-lucene/lucene/document"
)

type IndexWriter struct {
	writer *DocumentWriter
}

type IndexWriterConfig struct {
	analyzer *analyzer.PerFieldAnalyzer
}

func NewIndexWriter(config IndexWriterConfig) *IndexWriter {
	return &IndexWriter{
		writer: NewDocumentWriter(config.analyzer),
	}
}

func (w *IndexWriter) AddDocument(doc document.IndexableDocument) {
	w.writer.addDocuments([]document.IndexableDocument{doc})
}

func (w *IndexWriter) Flush() {
	w.writer.flush()
}
