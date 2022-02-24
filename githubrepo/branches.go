package githubrepos

/// Creates a new branch
func (repo *GithubRepository) CreateBranch(name string, fromBranch string, fromSha string) error {

	_, err := repo.Client.CreateNewRef(repo.owner, repo.name, name, fromBranch, fromSha)

	return err

}

/// Deletes a branch
func (repo *GithubRepository) DeleteBranch(name string) error {

	err := repo.Client.RemoveRef(repo.owner, repo.name, name)

	return err

}
