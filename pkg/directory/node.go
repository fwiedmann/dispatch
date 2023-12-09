// Package directory tbd.
package directory

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const RootPwdDir = "./"

// ErrorEmptyPath error
var ErrorEmptyPath = errors.New("given path input is empty")

// NewNode inits a new root node
func NewNode(startingPath string) *Node {
	if startingPath != RootPwdDir {
		startingPath = strings.TrimSuffix(startingPath, "/")
	}

	return &Node{
		name:     startingPath,
		path:     startingPath,
		rootPath: startingPath,
		depth:    1,
	}
}

// Node represents a directory, relative from its parent directory
type Node struct {
	name string

	path     string
	rootPath string

	depth int

	fileFound bool

	parent   *Node
	children []*Node
}

// Parent returns the parent Node of the current node.
// If Node is nil the rootPwdDir is reached.
func (n *Node) Parent() *Node {
	return n.parent
}

// AddDirectoryPath will update the Nodes tree for the given path.
// The returned nodes represents the last sub directory of the given path
func (n *Node) AddDirectoryPath(path string) (*Node, error) {
	err := validatePath(path)
	if err != nil {
		return nil, err
	}

	if filepath.IsAbs(path) {
		path = strings.TrimPrefix(path, n.rootPath+"/")
	}

	return n.add(strings.Split(path, "/"))
}

func validatePath(path string) error {
	if path == "" {
		return ErrorEmptyPath
	}

	return nil
}

func (n *Node) add(dirNames []string) (*Node, error) {

	isLastDirName := len(dirNames) == 1
	dirName := dirNames[0]

	if n.children == nil {
		n.children = make([]*Node, 0)
	}

	var foundNode *Node
	for _, presentNode := range n.children {
		if presentNode.name == dirName {
			foundNode = presentNode
			break
		}
	}

	if foundNode != nil {
		if isLastDirName {
			// TODO: hier auch nill?
			return nil, nil
		}
		return foundNode.add(removeFirstDirName(dirNames))
	}

	newNode := &Node{
		name:     dirName,
		path:     buildSubDirPath(n, dirName),
		rootPath: n.rootPath,
		parent:   n,
		depth:    n.depth + 1,
	}

	n.children = append(n.children, newNode)

	if isLastDirName {
		return newNode, nil
	}

	return newNode.add(removeFirstDirName(dirNames))
}

func removeFirstDirName(dirNames []string) []string {
	return dirNames[1:]
}

func buildSubDirPath(n *Node, subDirName string) string {
	formatString := "%s/%s"
	if n.path == RootPwdDir {
		formatString = "%s%s"
	}
	return fmt.Sprintf(formatString, n.path, subDirName)
}
