// Package git TODO
package git

import (
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const remoteName = "origin"
const rootPwdDir = "./"

// Client is an abstract git client for identifying directories which have changes aka git diffs
type Client struct {
	repo         *git.Repository
	targetBranch string
}

// NewClient inits the Client
func NewClient(dir string, targetBranch string) (Client, error) {
	repo, err := git.PlainOpen(dir)

	if err != nil {
		return Client{}, err
	}

	return Client{
		repo:         repo,
		targetBranch: targetBranch,
	}, nil
}

// DirectoriesWithChanges will compare the current Client repository
// with the Client targetBranch and returns all direcories which have changes.
// The targetBranch name will be used to get the HEAD reference from the remote.
// Current default value for the remote name is 'origin'. If it is needed to change the remote name,
// it should be part of the Client struct.
func (c Client) DirectoriesWithChanges() ([]string, error) {

	diff, err := c.diffBranches()
	if err != nil {
		return nil, err
	}

	return c.extractDirectories(diff), nil
}

func (c Client) diffBranches() (*object.Patch, error) {
	targetRef, err := c.repo.Reference(plumbing.NewRemoteReferenceName(remoteName, c.targetBranch), true)
	if err != nil {
		return nil, err
	}

	targetRefCommit, err := c.repo.CommitObject(targetRef.Hash())
	if err != nil {
		return nil, err
	}

	currentRef, err := c.repo.Head()
	if err != nil {
		return nil, err
	}

	currentRefCommit, err := c.repo.CommitObject(currentRef.Hash())
	if err != nil {
		return nil, err
	}

	return currentRefCommit.Patch(targetRefCommit)
}

func (c Client) extractDirectories(diff *object.Patch) []string {
	if diff == nil || len(diff.FilePatches()) == 0 {
		return []string{}
	}

	paths := make(map[string]any, len(diff.FilePatches()))

	for _, change := range diff.FilePatches() {
		from, to := change.Files()

		if from != nil {
			paths[filepath.Dir(from.Path())] = rootPwdDir
		}

		if to != nil {
			paths[filepath.Dir(to.Path())] = rootPwdDir
		}

	}

	dirs := make([]string, len(paths))
	for path := range paths {
		dirs = append(dirs, path)
	}

	return dirs
}
