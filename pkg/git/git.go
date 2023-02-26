// Package git TODO
package git

import (
	"github.com/go-git/go-git/v5"
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
	return []string{}
}
