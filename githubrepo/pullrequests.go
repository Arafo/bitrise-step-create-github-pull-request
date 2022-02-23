package githubrepos

import (
	gogithub "github.com/google/go-github/v42/github"
)

/// Fetches open pull requests
func (repo *GithubRepository) FetchPullRequests(fromBranch string, toBranch string) ([]*gogithub.PullRequest, error) {

	prs, err := repo.Client.FetchPullRequests(repo.owner, repo.name, toBranch, "open")

	if err != nil {
		return nil, err
	}

	return prs, nil
}

// Closes a single pull request for an id
func (repo *GithubRepository) ClosePullRequest(id int) error {

	err := repo.Client.ClosePullRequest(repo.owner, repo.name, id)

	return err

}

/// Create a new pull request
func (repo *GithubRepository) CreatePullRequest(title string, description string, fromBranch string, toBranch string) (int, error) {

	pullRequest, err := repo.Client.CreatePullRequest(repo.owner, repo.name, title, fromBranch, toBranch, description)

	if err != nil {
		return -1, err
	}

	return pullRequest.GetNumber(), nil

}
