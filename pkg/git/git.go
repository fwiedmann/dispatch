// Package git TODO
package git

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
func (c Client) DirectoriesWithChanges() []string {

	targetRef, err := c.repo.Reference(plumbing.NewBranchReferenceName(c.targetBranch), true)
	if err != nil {
		panic(err)
	}

	targetRefCommit, err := c.repo.CommitObject(targetRef.Hash())
	if err != nil {
		panic(err)
	}

	currentRef, err := c.repo.Head()
	if err != nil {
		panic(err)
	}

	commit, err := c.repo.CommitObject(currentRef.Hash())
	if err != nil {
		panic(err)
	}

	p, err := commit.Patch(targetRefCommit)
	if err != nil {
		panic(err)
	}

	for _, diff := range p.FilePatches() {
		fmt.Printf("%+v\n", diff)
	}

	return []string{}
}
