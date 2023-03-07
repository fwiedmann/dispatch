package git_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fwiedmann/dispatch/pkg/git"
	giclient "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
)

const sourceBranchName = "feat-187"
const targetBranchName = "master"

func TestNewClientWithError(t *testing.T) {
	// given when
	_, err := git.NewClient(t.TempDir(), targetBranchName)
	// then
	assert.Error(t, err)
}

func TestNewClient(t *testing.T) {
	// given
	testDir := t.TempDir()
	setUpTestRepo(t, testDir)

	// when
	_, err := git.NewClient(testDir, "")

	// then
	assert.Nil(t, err)
}

func TestDirectoriesWithChanges(t *testing.T) {
	// given
	testDir := t.TempDir()
	setUpTestRepo(t, testDir)

	c, err := git.NewClient(testDir, targetBranchName)

	// when

	c.DirectoriesWithChanges()

	// then
	assert.Nil(t, err)
}

func setUpTestRepo(t *testing.T, path string) {
	// git init, will create default branch called 'master'
	repo, err := giclient.PlainInit(path, false)
	assert.Nil(t, err)

	tree, err := repo.Worktree()
	assert.Nil(t, err)

	for _, dir := range testDirs {
		createDir(t, tree, path, dir)
		createFileAndCommit(t, tree, path, dir, testTargetBranchFileContent)
	}

	// create source/feature branch
	assert.Nil(t, tree.Checkout(&giclient.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(sourceBranchName),
		Create: true,
		Force:  true,
		Keep:   true,
	}))

	for _, dir := range testDirs {
		createFileAndCommit(t, tree, path, dir, testSourceBranchFileContent)
	}
}

func createDir(t *testing.T, tree *giclient.Worktree, testdir, subDir string) {
	absoluteDirPath := fmt.Sprintf("%s/%s", testdir, subDir)
	if err := os.Mkdir(absoluteDirPath, 0777); err != nil {
		panic(err)
	}
}

func createFileAndCommit(t *testing.T, tree *giclient.Worktree, testdir, subDir, content string) {
	absoluteFilePath := fmt.Sprintf("%s/%s/%s", testdir, subDir, testFileName)
	relativeFilePath := fmt.Sprintf("%s/%s", subDir, testFileName)

	assert.Nil(t, os.WriteFile(absoluteFilePath, []byte(content), 0777))

	_, err := tree.Add(relativeFilePath)
	assert.Nil(t, err)

	// git commit
	_, err = tree.Commit("feat(file): cool new feature", &giclient.CommitOptions{})
	assert.Nil(t, err)
}

const testFileName = "file.txt"
const testSourceBranchFileContent = "changed"
const testTargetBranchFileContent = "best feature ever"

var testDirs []string = []string{"foo", "bar"}
