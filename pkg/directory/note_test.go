// directory test
package directory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode(t *testing.T) {
	node := NewNode()

	assert.True(t, node.Parent() == nil)
	assert.Nil(t, node.childs)
	assert.Equal(t, "./", node.name)
	assert.Equal(t, "./", node.relativePath)
	assert.Equal(t, 1, node.depth)
	assert.False(t, node.fileFound)
}

func TestAddDirectiryPath_should_add_child_nodes(t *testing.T) {
	node := NewNode()

	testPath := "1/2/3"
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, err)
	assert.Equal(t, "3", childNode.name)
	assert.Equal(t, 4, childNode.depth)
	assert.Equal(t, "./"+testPath, childNode.relativePath)
	assert.False(t, childNode.fileFound)

	childNode, err = node.AddDirectoryPath(testPath + "/4")

	assert.Nil(t, err)
	assert.Equal(t, "4", childNode.name)
	assert.Equal(t, 5, childNode.depth)
	assert.Equal(t, "./"+testPath+"/4", childNode.relativePath)
	assert.False(t, childNode.fileFound)
}

func TestAddDirectiryPath_should_fail_when_absolute_path_is_given(t *testing.T) {
	node := NewNode()

	testPath := "/1/2/3"
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, childNode)
	assert.ErrorIs(t, err, ErrorInvalidPath)
}

func TestAddDirectiryPath_should_fail_when_absolute_path_is_empty(t *testing.T) {
	node := NewNode()

	testPath := ""
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, childNode)
	assert.ErrorIs(t, err, ErrorEmptyPath)
}
