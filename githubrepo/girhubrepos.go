package githubrepos

import "github.com/Arafo/bitrise-step-create-github-pull-request/github"

type GithubRepository struct {
	Client *github.GithubClient
	name   string
	owner  string
}

func NewGithubRepository(client github.GithubClient, name string, owner string) *GithubRepository {
	return &GithubRepository{
		Client: &client,
		name:   name,
		owner:  owner,
	}
}
