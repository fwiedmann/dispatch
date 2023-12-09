package directory

import (
	"encoding/json"
	"fmt"
	"github.com/fwiedmann/dispatch/pkg/code"
	"os"
)

const defaultSearchedFileName = ".code_owners"

// Analyzer is capable of traversing the given input paths upwards until the current root directory.
type Analyzer struct {
	root             *Node
	childs           []*Node
	searchedFileName string
}

// NewAnalyzer initializes the analyzer and builds the node tree for the given paths.
// Overwrite the defaults with the option functions
func NewAnalyzer(rootPath string, paths []string, opts ...AnalyzerOption) (Analyzer, error) {
	rootNode := NewNode(rootPath)
	childNodes := make([]*Node, 0)

	for _, path := range paths {
		child, err := rootNode.AddDirectoryPath(path)
		if err != nil {
			return Analyzer{}, err
		}
		if child != nil {
			childNodes = append(childNodes, child)
		}
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

// Analyze with breadth first search for each child node (bottom -> up)
func (a *Analyzer) Analyze() ([]code.Info, error) {
	fileRefs := make([]string, 0)
	for _, child := range a.childs {
		file := a.analyzeChild(child)
		if file != nil {
			fileRefs = append(fileRefs, *file)
		}
	}
	return a.readFiles(fileRefs)
}

func (a *Analyzer) readFiles(fileRefs []string) ([]code.Info, error) {
	infos := make([]code.Info, 0)
	for _, ref := range fileRefs {
		content, err := os.ReadFile(ref)
		if err != nil {
			return nil, err
		}
		var info code.Info

		if err := json.Unmarshal(content, &info); err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func (a *Analyzer) analyzeChild(child *Node) *string {
	for child != nil {
		// already added, ignore
		if child.fileFound {
			return nil
		}
		filename := a.pathToSearchingFile(child.path)
		_, err := os.Lstat(filename)
		if err != nil {
			child = child.parent
			continue
		}
		// mark node, so it won't be added twice
		child.fileFound = true
		return &filename
	}
	return nil
}

func (a *Analyzer) pathToSearchingFile(path string) string {
	return fmt.Sprintf("%s/%s", path, a.searchedFileName)
}
