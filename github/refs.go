package github

import (
	"fmt"

	gogithub "github.com/google/go-github/v42/github"
)

// Return a reference to the head of the branch
func (c *GithubClient) GetRefToHead(owner string, repo string, branch string) (ref *gogithub.Reference, err error) {

	targetBranchRef := fmt.Sprintf("refs/heads/%s", branch)

	ref, _, err = c.Client.Git.GetRef(c.Context, owner, repo, targetBranchRef)

	return ref, err

}

// Creates a new reference
func (c *GithubClient) CreateNewRef(owner string, repo string, refName string, atBranch string, atSha string) (ref *gogithub.Reference, err error) {

	newRefName := fmt.Sprintf("refs/heads/%s", refName)

	reference := gogithub.Reference{
		Ref: &newRefName,
		Object: &gogithub.GitObject{
			SHA: &atSha,
		},
	}

	ref, _, err = c.Client.Git.CreateRef(c.Context, owner, repo, &reference)

	return ref, err

}
