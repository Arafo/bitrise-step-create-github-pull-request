package main

import (
	"fmt"

	"github.com/Arafo/bitrise-step-create-github-pull-request/github"
	"github.com/bitrise-io/go-utils/log"
	gogithub "github.com/google/go-github/v42/github"
)

func getLastCommitSHA(githubClient github.GithubClient, conf Config) (string, error) {

	owner, repo := ownerAndRepo(conf.RepositoryURL)

	targetBranchRef := fmt.Sprintf("refs/heads/%s", conf.TargetBranch)

	ref, _, err := githubClient.Git.GetRef(githubClient.Context, owner, repo, targetBranchRef)

	if err != nil {
		log.Errorf("Unable to get target brancg ref: %s\n", err)
		return "", err
	}

	sha := ref.GetObject().GetSHA()
	fmt.Println("Got head commit SHA", sha)

	return sha, nil

}

func createNewBranch(githubClient github.GithubClient, conf Config, sha string) error {

	owner, repo := ownerAndRepo(conf.RepositoryURL)

	newRefName := fmt.Sprintf("refs/heads/%s", conf.NewBranch)

	fmt.Println("Creating new Ref:", newRefName)

	reference := gogithub.Reference{
		Ref: &newRefName,
		Object: &gogithub.GitObject{
			SHA: &sha,
		},
	}

	_, _, err := githubClient.Git.CreateRef(githubClient.Context, owner, repo, &reference)

	if err != nil {
		log.Errorf("Unable to create new ref: %s\n", err)
		return err
	}

	return nil
}

// create the pr
func CreateNewPr(githubClient github.GithubClient, conf Config) error {

	owner, repo := ownerAndRepo(conf.RepositoryURL)
	targetBranch := conf.TargetBranch
	newBranch := conf.NewBranch
	sourceFiles := conf.SourceFiles
	pullRequestTitle := conf.PullRequestTitle
	pullRequestDescription := conf.PullRequestDescription

	fmt.Println("Getting last commit on Target branch...")
	sha, err := getLastCommitSHA(githubClient, conf)
	if err != nil {
		log.Errorf("Unable to get last commit SHA: %s\n", err)
		return err
	}

	fmt.Printf("Creating new branch (%s) from base (%s) at SHA (%s)\n", newBranch, targetBranch, sha)
	err = createNewBranch(githubClient, conf, sha)
	if err != nil {
		log.Errorf("Unable to create new branch: %s\n", err)
		return err
	}

	_, _, err = githubClient.Repositories.GetBranch(githubClient.Context, owner, repo, newBranch, true)
	if err != nil {
		log.Errorf("Unable to get branch: %s\n", err)
		return err
	}

	fmt.Println("Getting commit branch reference...")
	ref, err := githubClient.GetCommitBranchReference(owner, repo, targetBranch, newBranch)
	if err != nil {
		log.Errorf("Unable to get/create the commit reference: %s\n", err)
		return err
	}

	fmt.Println("Generating tree difference from sourcefiles...")
	tree, err := githubClient.GenerateTree(owner, repo, sourceFiles, ref)
	if err != nil {
		log.Errorf("Unable to create the tree based on the provided files: %s\n", err)
		return err
	}

	fmt.Println("Pushing new commit...")
	commit, err := githubClient.PushCommit(owner, repo, conf.BotName, conf.BotEmail, conf.CommitMessage, ref, tree)
	if err != nil {
		log.Errorf("Unable to create the commit: %s\n", err)
		return err
	} else {
		fmt.Println("New commit created SHA:", commit.GetSHA())
	}

	fmt.Println("Creating pull request...")
	pullRequest, err := githubClient.CreatePullRequest(owner, repo, pullRequestTitle, newBranch, targetBranch, pullRequestDescription)
	if err != nil {
		log.Errorf("Github API call failed when creating a Pull Request: %w\n", err)
		return err
	} else {
		fmt.Println("New pull request created id:", pullRequest.GetNumber())
	}

	return nil
}
