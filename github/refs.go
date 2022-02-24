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
func (c *GithubClient) CreateNewRef(owner string, repo string, branchName string, atBranch string, atSha string) (ref *gogithub.Reference, err error) {

	newRefName := fmt.Sprintf("refs/heads/%s", branchName)

	reference := gogithub.Reference{
		Ref: &newRefName,
		Object: &gogithub.GitObject{
			SHA: &atSha,
		},
	}

	ref, _, err = c.Client.Git.CreateRef(c.Context, owner, repo, &reference)

	return ref, err

}

// Remove existing reference
func (c *GithubClient) RemoveRef(owner string, repo string, branchName string) error {

	refName := fmt.Sprintf("refs/heads/%s", branchName)

	_, err := c.Client.Git.DeleteRef(c.Context, owner, repo, refName)

	return err

}
