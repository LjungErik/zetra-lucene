package analysis

type Term string

type Analyzer interface {
	Analyze(data string) []Term
}

type DefaultAnalyzer struct {
}
