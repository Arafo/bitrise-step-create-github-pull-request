package githubrepos

/// Fetches open pull requests
func (repo *GithubRepository) CreateBranch(name string, fromBranch string, fromSha string) error {

	_, err := repo.Client.CreateNewRef(repo.owner, repo.name, name, fromBranch, fromSha)

	return err

}
