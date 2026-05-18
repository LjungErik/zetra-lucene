package analyzer

type PerFieldAnalyzer struct {
	analyzers       map[string]Analyzer
	defaultAnalyzer Analyzer
}

func NewPerFieldAnalyzer() *PerFieldAnalyzer {
	return &PerFieldAnalyzer{
		analyzers:       make(map[string]Analyzer),
		defaultAnalyzer: NewEnglishLanguageAnalyzer(),
	}
}

func (a *PerFieldAnalyzer) Set(fieldName string, analyzer Analyzer) {
	a.analyzers[fieldName] = analyzer
}

func (a *PerFieldAnalyzer) Get(fieldName string) Analyzer {
	if analyzer, ok := a.analyzers[fieldName]; ok {
		return analyzer
	}

	return a.defaultAnalyzer
}
