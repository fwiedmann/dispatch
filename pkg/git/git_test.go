package git_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/fwiedmann/dispatch/pkg/git"
	giclient "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
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
	if err != nil {
		t.Fatal(err)
	}

	_, err = repo.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{"https://github.com/fwiedmann/does-not-exist"},
	})
	if err != nil {
		t.Fatal(err)
	}

	tree, err := repo.Worktree()
	if err != nil {
		t.Fatal(err)
	}

	for _, dir := range testDirs {
		createDir(t, tree, path, dir)
		createFileAndCommit(t, tree, path, dir, testTargetBranchFileContent)
	}

	createFakeRemote(t, repo, path)

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

	err := os.WriteFile(absoluteFilePath, []byte(content), 0777)
	if err != nil {
		t.Fatal(err)
	}

	_, err = tree.Add(relativeFilePath)
	if err != nil {
		t.Fatal(err)
	}

	// git commit
	_, err = tree.Commit("feat(file): cool new feature", &giclient.CommitOptions{})
	if err != nil {
		t.Fatal(err)
	}
}

func createFakeRemote(t *testing.T, repo *giclient.Repository, repoPath string) {
	head, err := repo.Head()
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(fmt.Sprintf("%s/.git/refs/remotes", repoPath), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Mkdir(fmt.Sprintf("%s/.git/refs/remotes/origin", repoPath), 0777)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/.git/refs/remotes/origin/%s", repoPath, targetBranchName), []byte(head.Hash().String()), 0777)

	if err != nil {
		t.Fatal(err)
	}
}

const testFileName = "file.txt"
const testSourceBranchFileContent = "changed"
const testTargetBranchFileContent = "best feature ever"

var testDirs []string = []string{"foo", "bar"}
