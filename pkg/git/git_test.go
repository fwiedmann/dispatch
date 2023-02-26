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
const targetBranchName = "main"

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
	// git init
	repo, err := giclient.PlainInit(path, false)
	assert.Nil(t, err)

	// create test file
	assert.Nil(t, os.WriteFile(fmt.Sprintf("%s/%s", path, testFileName), []byte(testTargetBranchFileContent), 0777))

	tree, err := repo.Worktree()
	assert.Nil(t, err)

	// git add
	_, err = tree.Add(testFileName)
	assert.Nil(t, err)

	// git commit
	_, err = tree.Commit("test", &giclient.CommitOptions{})
	assert.Nil(t, err)

	// create target branch
	assert.Nil(t, tree.Checkout(&giclient.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("/refs/heads/%s", targetBranchName)),
		Create: true,
		Force:  true,
		Keep:   true,
	}))

	// create source branch
	assert.Nil(t, tree.Checkout(&giclient.CheckoutOptions{
		Branch: plumbing.ReferenceName(fmt.Sprintf("/refs/heads/%s", sourceBranchName)),
		Create: true,
		Force:  true,
		Keep:   true,
	}))

	// update file content
	assert.Nil(t, os.WriteFile(fmt.Sprintf("%s/%s", path, testFileName), []byte(testSourceBranchFileContent), 0777))

	// git add
	_, err = tree.Add(testFileName)
	assert.Nil(t, err)

	// git commit
	_, err = tree.Commit("test", &giclient.CommitOptions{})
	assert.Nil(t, err)
}

const testFileName = "file.txt"
const testSourceBranchFileContent = "changed"
const testTargetBranchFileContent = "best feature ever"
