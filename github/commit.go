package github

import (
	"time"

	gogithub "github.com/google/go-github/v42/github"
)

func (c *GithubClient) PushCommit(owner, repo string, authorName string, authorEmail string, commitMessage string, ref *gogithub.Reference, tree *gogithub.Tree) (commit *gogithub.Commit, err error) {
	parent, _, err := c.Repositories.GetCommit(c.Context, owner, repo, *ref.Object.SHA, nil)
	if err != nil {
		return nil, err
	}

	parent.Commit.SHA = parent.SHA

	date := time.Now()
	author := gogithub.CommitAuthor{Date: &date, Name: gogithub.String(authorName), Email: gogithub.String(authorEmail)}
	commitToPush := gogithub.Commit{Author: &author, Message: gogithub.String(commitMessage), Tree: tree, Parents: []*gogithub.Commit{parent.Commit}}
	newCommit, _, err := c.Git.CreateCommit(c.Context, owner, repo, &commitToPush)
	if err != nil {
		return nil, err
	}

	ref.Object.SHA = newCommit.SHA
	_, _, err = c.Git.UpdateRef(c.Context, owner, repo, ref, false)
	return newCommit, err
}
