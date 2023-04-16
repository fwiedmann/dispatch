// Package directory tbd.
package directory

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

const rootPwdDir = "./"

// ErrorInvalidPath error
var ErrorInvalidPath = errors.New("given path is absolute. Only relatives are allowed")

// ErrorEmptyPath error
var ErrorEmptyPath = errors.New("given path input is empty")

// NewNode inits a new root node
func NewNode() *Node {
	return &Node{
		name:         rootPwdDir,
		relativePath: rootPwdDir,
		depth:        1,
	}
}

// Node represents a directory, relative from its parent directory
type Node struct {
	name string

	relativePath string

	depth int

	fileFound bool

	parent *Node
	childs []*Node
}

// Parent returns the parent Node of the current node.
// If Node is nil the rootPwdDir is reached.
func (n *Node) Parent() *Node {
	return n.parent
}

// FoundFile marks the Node that the searched file was found
// This can be helpful during the bottom up traversal.
// If a node is already marked, the traversal can be stopped.
func (n *Node) FoundFile() {
	n.fileFound = true
}

// AddDirectoryPath will update the Nodes tree for the given path.
// The returned nodes represents the last sub directory of the given path
func (n *Node) AddDirectoryPath(path string) (*Node, error) {
	err := validatePath(path)
	if err != nil {
		return nil, err
	}

	return n.add(strings.Split(path, "/"))
}

func validatePath(path string) error {
	if path == "" {
		return ErrorEmptyPath
	}

	if filepath.IsAbs(path) {
		return ErrorInvalidPath
	}
	return nil
}

func (n *Node) add(dirNames []string) (*Node, error) {

	isLastDirName := len(dirNames) == 1
	dirName := dirNames[0]

	if n.childs == nil {
		n.childs = make([]*Node, 0)
	}

	var foundNode *Node
	for _, presentNode := range n.childs {
		if presentNode.name == dirName {
			foundNode = presentNode
			break
		}
	}

	if foundNode != nil {
		if isLastDirName {
			return foundNode, nil
		}
		return foundNode.add(removeFirstDirName(dirNames))
	}

	newNode := &Node{
		name:         dirName,
		relativePath: builSubDirPath(n, dirName),
		parent:       n,
		depth:        n.depth + 1,
	}

	n.childs = append(n.childs, newNode)

	if isLastDirName {
		return newNode, nil
	}

	return newNode.add(removeFirstDirName(dirNames))
}

func removeFirstDirName(dirNames []string) []string {
	return dirNames[1:]
}

func builSubDirPath(n *Node, subDirName string) string {
	formatString := "%s/%s"
	if n.relativePath == rootPwdDir {
		formatString = "%s%s"
	}
	return fmt.Sprintf(formatString, n.relativePath, subDirName)
}
