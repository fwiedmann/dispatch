// Package git TODO
package git

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

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

// DirectoriesWithChanges will compare the current Client repository with the Client targetBranch and returns all direcories which have changes.
func (c Client) DirectoriesWithChanges() ([]string, error) {
	diff, err := c.diffBranches()
	if err != nil {
		return nil, err
	}

	return c.extractDirectories(diff), nil
}

func (c Client) diffBranches() (*object.Patch, error) {
	targetRef, err := c.repo.Reference(plumbing.NewBranchReferenceName(c.targetBranch), true)
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

// TODO
func (c Client) extractDirectories(diff *object.Patch) []string {
	return []string{}
}
