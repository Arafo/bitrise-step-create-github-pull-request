package githubrepos

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
)

// Returns the last commit on a branch
func (repo *GithubRepository) GetLastCommit(branch string) (string, error) {

	ref, err := repo.Client.GetRefToHead(repo.owner, repo.name, branch)

	if err != nil {
		return "", err
	}

	sha := ref.GetObject().GetSHA()

	return sha, nil

}

// Returns the last commit on a branch
func (repo *GithubRepository) CreateCommit(
	targetBranch string,
	commitBranch string,
	sourceFiles string,
	botName string,
	botEmail string,
	message string,
) error {

	ref, err := repo.Client.GetRefToHead(repo.owner, repo.name, commitBranch)
	if err != nil {
		log.Errorf("Unable to get/create the commit reference: %s\n", err)
		return err
	}

	tree, err := repo.Client.GenerateTree(repo.owner, repo.name, sourceFiles, ref)
	if err != nil {
		log.Errorf("Unable to create the tree based on the provided files: %s\n", err)
		return err
	}

	commit, err := repo.Client.PushCommit(repo.owner, repo.name, botName, botEmail, message, ref, tree)
	if err != nil {
		log.Errorf("Unable to create the commit: %s\n", err)
		return err
	}

	fmt.Println("Commit created SHA:", commit.GetSHA())

	return nil

}
