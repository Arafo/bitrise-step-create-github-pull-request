package github

import (
	gogithub "github.com/google/go-github/v42/github"
)

// Fetches all pull requests for a repo to a target branch
func (c *GithubClient) FetchPullRequests(owner, repo string, toBranch string, state string) (prs []*gogithub.PullRequest, err error) {

	var filter = gogithub.PullRequestListOptions{
		Base:  toBranch,
		State: state,
	}

	prs, _, err = c.PullRequests.List(c.Context, owner, repo, &filter)

	return prs, err

}

// Fetches all pull requests for a repo to a target branch
func (c *GithubClient) ClosePullRequest(owner string, repo string, id int) error {

	targetState := "closed"

	var edit = gogithub.PullRequest{
		State: &targetState,
	}

	_, _, err := c.PullRequests.Edit(c.Context, owner, repo, id, &edit)

	return err

}

// Create a new Pull request
func (c *GithubClient) CreatePullRequest(owner, repo string, title string, commitBranch string, baseBranch string, description string) (*gogithub.PullRequest, error) {
	pullRequestToCreate := gogithub.NewPullRequest{
		Title:               gogithub.String(title),
		Head:                gogithub.String(commitBranch),
		Base:                gogithub.String(baseBranch),
		Body:                gogithub.String(description),
		MaintainerCanModify: gogithub.Bool(true),
	}

	pullRequest, _, err := c.PullRequests.Create(c.Context, owner, repo, &pullRequestToCreate)

	return pullRequest, err
}
