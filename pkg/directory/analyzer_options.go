package directory

type AnalyzerOption func(*Analyzer)

var WithSearchedFileName = func(name string) AnalyzerOption {
	return func(a *Analyzer) {
		a.searchedFileName = name
	}
}
