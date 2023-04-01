package directory

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var InvalidPathError = errors.New("given path is absolute. Only relatives are allowed.")
var NotADirectoryError = errors.New("given path does not point to a directory")
var EmptyPathError = errors.New("given path input is empty")

func NewNode() *Node {
	return &Node{}
}

type Node struct {
	name string

	fileFound bool

	parent *Node
	childs []*Node
}

func (n *Node) Parent() *Node {
	if n.parent == nil {
		return nil
	}

	return n.parent
}

func (n *Node) foundFile() {
	n.fileFound = true
}

func (n *Node) AddDirectoryPath(path string) (*Node, error) {
	err := validatePath(path)
	if err != nil {
		return nil, err
	}

	return n.add(strings.Split(path, "/"))
}

func validatePath(path string) error {
	if path == "" {
		return EmptyPathError
	}

	if filepath.IsAbs(path) {
		return InvalidPathError
	}

	stat, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return NotADirectoryError
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
		name:   dirName,
		parent: n,
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
