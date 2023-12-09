// directory test
package directory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNode_with_relative_path(t *testing.T) {
	node := NewNode(RootPwdDir)

	assert.True(t, node.Parent() == nil)
	assert.Nil(t, node.children)
	assert.Equal(t, "./", node.name)
	assert.Equal(t, "./", node.path)
	assert.Equal(t, 1, node.depth)
	assert.False(t, node.fileFound)
}

func TestNewNode_with_absolute_path(t *testing.T) {
	dir := t.TempDir()
	node := NewNode(dir + "/")

	assert.True(t, node.Parent() == nil)
	assert.Nil(t, node.children)
	assert.Equal(t, dir, node.name)
	assert.Equal(t, dir, node.path)
	assert.Equal(t, 1, node.depth)
	assert.False(t, node.fileFound)
}

func TestAddDirectoryPath_should_add_child_nodes_when_root_path_is_relative(t *testing.T) {
	node := NewNode(RootPwdDir)

	testPath := "1/2/3"
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, err)
	assert.Equal(t, "3", childNode.name)
	assert.Equal(t, 4, childNode.depth)
	assert.Equal(t, "./"+testPath, childNode.path)
	assert.False(t, childNode.fileFound)

	childNode, err = node.AddDirectoryPath(testPath + "/4")

	assert.Nil(t, err)
	assert.Equal(t, "4", childNode.name)
	assert.Equal(t, 5, childNode.depth)
	assert.Equal(t, "./"+testPath+"/4", childNode.path)
	assert.False(t, childNode.fileFound)
}

func TestAddDirectoryPath_should_add_child_nodes_when_root_path_is_absolute(t *testing.T) {
	dir := t.TempDir()
	node := NewNode(dir)

	testPath := dir + "/1/2/3"
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, err)
	assert.Equal(t, "3", childNode.name)
	assert.Equal(t, 4, childNode.depth)
	assert.Equal(t, testPath, childNode.path)
	assert.False(t, childNode.fileFound)

	childNode, err = node.AddDirectoryPath(testPath + "/4")

	assert.Nil(t, err)
	assert.Equal(t, "4", childNode.name)
	assert.Equal(t, 5, childNode.depth)
	assert.Equal(t, testPath+"/4", childNode.path)
	assert.False(t, childNode.fileFound)
}

func TestAddDirectoryPath_should_fail_when_absolute_path_is_empty(t *testing.T) {
	node := NewNode(RootPwdDir)

	testPath := ""
	childNode, err := node.AddDirectoryPath(testPath)

	assert.Nil(t, childNode)
	assert.ErrorIs(t, err, ErrorEmptyPath)
}
func TestAddDirectoryPath_should_return_nil_when_path_is_already_present(t *testing.T) {
	node := NewNode(RootPwdDir)

	testPath := "1"

	childNode, err := node.AddDirectoryPath(testPath)
	assert.NotNil(t, childNode)
	assert.NoError(t, err)

	childNode, err = node.AddDirectoryPath(testPath)
	assert.Nil(t, childNode)
	assert.NoError(t, err)
}
