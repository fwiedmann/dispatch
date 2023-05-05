package directory

const defaultSearchedFileName = ".code_owners"

// Analyzer is capable of traversing the given input paths upwards until the current root directory.
type Analyzer struct {
	root             *Node
	childs           []*Node
	searchedFileName string
}

// NewAnalyzer initializes the analyzer and builds the node tree for the given paths.
// Overwrite the defaults with the option functions
func NewAnalyzer(paths []string, opts ...AnalyzerOption) (Analyzer, error) {
	rootNode := NewNode()
	childNodes := make([]*Node, 0)

	for _, path := range paths {
		child, err := rootNode.AddDirectoryPath(path)
		if err != nil {
			return Analyzer{}, err
		}
		childNodes = append(childNodes, child)
	}

	analyzer := &Analyzer{
		root:             rootNode,
		childs:           childNodes,
		searchedFileName: defaultSearchedFileName,
	}

	for _, opt := range opts {
		opt(analyzer)
	}

	return *analyzer, nil
}
