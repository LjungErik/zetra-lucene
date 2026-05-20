package query

import (
	"github.com/LjungErik/zetra-lucene/lucene/search/collector"
	"github.com/LjungErik/zetra-lucene/lucene/search/context"
)

type Query interface {
	Execute(context.IndexReaderContext, collector.TopDocumentCollector)
}
