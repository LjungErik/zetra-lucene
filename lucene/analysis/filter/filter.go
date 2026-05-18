package filter

import "github.com/LjungErik/zetra-lucene/lucene/analysis"

type Filter interface {
	Apply([]analysis.Token) []analysis.Token
}
